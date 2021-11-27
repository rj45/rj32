package xform

import (
	"go/constant"
	"go/types"
	"log"

	"github.com/rj45/rj32/gorj/ir"
	"github.com/rj45/rj32/gorj/ir/op"
	"github.com/rj45/rj32/gorj/ir/reg"
	"github.com/rj45/rj32/gorj/sizes"
)

func calls(val *ir.Value) int {
	if val.Op != op.Call {
		return 0
	}

	changes := 0

	fnType := val.Arg(0).Type.(*types.Signature)

	if fnType.Results().Len() == 1 && val.Reg != reg.A0 {
		changes++
		val.Reg = reg.A0
		copy := val.Block().InsertCopy(val.Index()+1, val, reg.None)
		val.ReplaceOtherUsesWith(copy)
	} else if fnType.Results().Len() > 1 {
		if fnType.Results().Len() > 3 {
			log.Panicln("greater than 3 returned values not supported yet in function", val.Arg(0))
		}
		for i := 0; i < val.NumArgUses(); i++ {
			ext := val.ArgUse(i)
			if ext.Reg != reg.None {
				continue
			}
			changes++
			num, _ := constant.Int64Val(ext.Value)
			ext.Reg = reg.ArgRegs[num]
			repl := ir.BuildAfter(ext).Op(op.Copy, ext.Type, ext).PrevVal()
			ext.ReplaceOtherUsesWith(repl)
		}
	}

	if val.NumArgs() > 1 {
		if val.Arg(1).Reg != reg.A1 {
			changes++
			val.ReplaceArg(1, val.Block().InsertCopy(val.Index(), val.Arg(1), reg.A1))
		}

		if val.NumArgs() > 2 && val.Arg(2).Reg != reg.A2 {
			changes++
			val.ReplaceArg(2, val.Block().InsertCopy(val.Index(), val.Arg(2), reg.A2))
		}

		slots := val.NumArgs() - 3
		if slots >= 0 {
			if val.Func().ArgSlots < slots {
				val.Func().ArgSlots = slots
			}

			for i := 0; i < slots; i++ {
				if val.Arg(i+3).Op != op.Store {
					changes++
					val.ReplaceArg(i+3, val.Block().InsertCopy(val.Index(), val.Arg(i+3), reg.StackSlot(i)))
				}
			}
		}
	}

	return changes
}

var _ = addToPass(Elaboration, calls)

func builtinCalls(val *ir.Value) int {
	if val.Op != op.CallBuiltin {
		return 0
	}

	fn := val.Block().Func()
	name := constant.StringVal(val.Value)
	otherFn := fn.Pkg.LookupFunc(name)

	if otherFn == nil {
		if name == "builtin__len" {
			// length is always at memory location of string/slice plus 1
			ir.BuildReplacement(val).Op(op.Load, val.Type, val.Arg(1), 1)
			val.Value = nil
			val.Func().NumCalls--
			return 1
		}

		log.Fatalf("builtin %s not loaded!", name)
	}

	val.Value = nil

	arg := fn.Const(otherFn.Type, constant.MakeString(name))
	arg.Op = op.Func
	val.Type = otherFn.Type

	val.InsertArg(0, arg)

	val.Op = op.Call

	return 1
}

var _ = addToPass(Elaboration, builtinCalls)

// indexAddrs converts `IndexAddr` instructions into a `mul` and `add` instruction
// The `mul` is by a constant which can be optimized into shifts and adds by
// some other piece of code.
func indexAddrs(val *ir.Value) int {
	if val.Op != op.IndexAddr {
		return 0
	}

	elem := val.Type.(*types.Pointer).Elem()

	size := sizes.Sizeof(elem)

	sizeval := val.Func().IntConst(size)

	mulval := val.Func().NewValue(op.Mul, types.Typ[types.Int], val.Arg(1), sizeval)
	val.Block().InsertInstr(val.Index(), mulval)

	val.Op = op.Add
	val.ReplaceArg(1, mulval)

	return 1
}

var _ = addToPass(Elaboration, indexAddrs)

func lookups(val *ir.Value) int {
	if val.Op != op.Lookup {
		return 0
	}

	if val.NumArgs() != 2 {
		return 0
	}

	arg0 := val.Arg(0)

	basic, ok := arg0.Type.(*types.Basic)
	if !ok {
		return 0
	}

	if basic.Kind() != types.String {
		return 0
	}

	indexArg := val.Arg(1)

	// a string is a tuple of an address and length
	bd := ir.BuildBefore(val)

	if arg0.Op == op.Global {
		bd = bd.Op(op.Add, arg0.Type, arg0, reg.GP)
		arg0 = bd.PrevVal()
	}

	// load the address
	bd = bd.Op(op.Load, arg0.Type, arg0, 0)
	address := bd.PrevVal()

	if indexArg.Op.IsConst() {
		i, _ := constant.Int64Val(indexArg.Value)
		load := ir.BuildBefore(val).Op(op.Load, val.Type, address, i>>1).PrevVal()
		if i&1 == 0 {
			ir.BuildReplacement(val).
				Op(op.ShiftRight, val.Type, load, 8)
		} else {
			ir.BuildReplacement(val).
				Op(op.And, val.Type, load, 0xFF)
		}
		return 1
	}

	bd = bd.
		Op(op.ShiftRight, indexArg.Type, indexArg, 1).
		Op(op.Add, address.Type, address, ir.PrevBuildVal()).
		Op(op.Load, val.Type, ir.PrevBuildVal(), 0)

	and := bd.Op(op.And, indexArg.Type, indexArg, 1).PrevVal()

	ir.BuildReplacement(val).
		Op(op.StringShift, val.Type, bd.PrevVal(), and)

	return 1
}

var _ = addToPass(Elaboration, lookups)

func fieldAddrs(val *ir.Value) int {
	if val.Op != op.FieldAddr {
		return 0
	}

	field, ok := constant.Int64Val(val.Value)
	if !ok {
		panic("expected int constant")
	}

	elem := val.Arg(0).Type.(*types.Pointer).Elem()
	strct := elem.Underlying().(*types.Struct)

	fields := sizes.Fieldsof(strct)
	offsets := sizes.Offsetsof(fields)
	offset := offsets[field]

	if offset == 0 {
		// would just be adding zero, so this instruction can just be removed
		val.ReplaceWith(val.Arg(0))
		return 1
	}

	val.Op = op.Add
	val.Value = nil
	val.InsertArg(-1, val.Func().Const(val.Type, constant.MakeInt64(offset)))

	return 1
}

var _ = addToPass(Elaboration, fieldAddrs)

func gpAdjustLoadStores(val *ir.Value) int {
	if val.Op != op.Load && val.Op != op.Store {
		return 0
	}

	if val.Op == op.Load && val.NumArgs() == 2 {
		return 0
	}

	if val.Op == op.Store && val.NumArgs() == 3 {
		return 0
	}

	// if storing a constant directly to memory
	if val.Op == op.Store && val.Arg(val.NumArgs()-1).Op.IsConst() {
		con := val.Arg(val.NumArgs() - 1)
		cp := val.Block().InsertCopy(val.Index(), con, con.Reg)
		val.ReplaceArg(val.NumArgs()-1, cp)
		return 1
	}

	if val.Arg(0).Op.IsConst() {
		arg := val.RemoveArg(0)
		val.InsertArg(0, val.Func().FixedReg(reg.GP))
		val.InsertArg(1, arg)
		return 1
	}

	val.InsertArg(1, val.Func().IntConst(0))

	// need to add GP to the global value
	// so we check if there is a path to a global, then
	// add a add to GP to the function prologue, then
	// go through and replace each occurance of the global
	// with a reference to that add

	path := val.FindPathTo(func(v *ir.Value) bool {
		return v.Op == op.Global
	})

	if path != nil {
		global := path[len(path)-1]
		user := path[len(path)-2]
		if user.Op != op.Add || user.NumArgs() != 2 || user.Arg(0).Reg != reg.GP || user.Arg(1) != global {
			entry := user.Func().Blocks()[0]

			var addGP *ir.Value
			for i := 0; i < entry.NumInstrs(); i++ {
				instr := entry.Instr(i)
				if instr.Op == op.Add && instr.NumArgs() == 2 && instr.Arg(0).Reg == reg.GP && instr.Arg(1) == global {
					addGP = instr
				}
			}

			if addGP == nil {
				addGP = user.Func().NewValue(op.Add, global.Type, user.Func().FixedReg(reg.GP), global)
				entry.InsertInstr(0, addGP)
			}

			user.ReplaceArg(user.ArgIndex(global), addGP)
			return 1
		}
	}

	return 1
}

var _ = addToPass(Elaboration, gpAdjustLoadStores)

func fixupConverts(val *ir.Value) int {
	if val.Op != op.Convert {
		return 0
	}

	if sizes.Sizeof(val.Arg(0).Type) != sizes.Sizeof(val.Type) {
		log.Fatalf("Unable to convert %#v to %#v", val.Arg(0).Type, val.Type)
	}

	val.ReplaceWith(val.Arg(0))

	return 1
}

var _ = addToPass(Elaboration, fixupConverts)

func useParameterRegisters(val *ir.Value) int {
	if val.Op != op.Parameter {
		return 0
	}
	index := -1
	for j, p := range val.Func().Params {
		if p == val {
			index = j
		}
	}
	if index < 0 {
		log.Panicln("could not find parameter", val)
	}

	switch index {
	case 0:
		val.Op = op.Copy
		val.InsertArg(-1, val.Func().FixedReg(reg.A1))
		val.Value = nil
	case 1:
		val.Op = op.Copy
		val.InsertArg(-1, val.Func().FixedReg(reg.A2))
		val.Value = nil
	default:
		// we don't know yet what the stack frame size will be
		// so, leave for the prologue code to convert this to a load

		return 0
	}

	return 1
}

var _ = addToPass(Elaboration, useParameterRegisters)

func allocIterators(val *ir.Value) int {
	if val.Op != op.Range {
		return 0
	}

	if val.NumBlockUses() > 0 {
		log.Panicln("expecting unexpected block use of", val.ShortString(), "in function", val.Func())
	}

	var nextFunc types.Type
	for i := 0; i < val.NumArgUses(); i++ {
		next := val.ArgUse(i)

		if next.Op == op.CallBuiltin {
			// wait for the builtin to be converted to a call
			return 0
		}

		if next.Op != op.Call {
			log.Panicln("expecting all uses of a range iter to be next calls, but found", next.ShortString(), "in function", next.Func())
		}

		// insert the original string/slice as an argument to the next call
		next.InsertArg(1, val.Arg(0))
		nextFunc = next.Arg(0).Type
	}

	sig := nextFunc.(*types.Signature)

	// pointer to iterator
	itPtr := sig.Params().At(1).Type()

	// add a local to function entry block
	local := ir.BuildAt(val.Func().Blocks()[0], 0).Op(op.Local, itPtr).PrevVal()

	// convert the current instruction into a store 0 into iterator
	ir.BuildBefore(val).Op(op.Copy, itPtr, 0).
		Op(op.Store, itPtr, local, 0, ir.PrevBuildVal())

	val.ReplaceWith(local)

	return 1
}

var _ = addToPass(Elaboration, allocIterators)

// func allocateLocals(val *ir.Value) int {
// 	if val.Op != op.Local {
// 		return 0
// 	}

// 	size := int64(1)
// 	if val.Type.String() != "iter" {
// 		size = sizes.Sizeof(val.Type)
// 	}

// 	fn := val.Func()

// 	// todo: somehow mark these slots as in use
// 	startSlot := fn.SpillSlots

// 	fn.SpillSlots += size

// 	val.Op = op.Add
// 	val.InsertArg(0, fn.FixedReg(reg.SP), startSlot + fn.ArgSlots)

// 	return 1
// }

// var _ = addToPass(Elaboration, allocateLocals)

func rollupCompareToBlockOp(val *ir.Value) int {
	if !val.Op.IsCompare() {
		return 0
	}

	blk := val.Block()

	if blk.Op != op.If {
		return 0
	}

	if blk.Control(0) != val {
		return 0
	}

	vop := val.Op

	arg0 := val.Arg(0)
	arg1 := val.Arg(1)

	// swap if the constant somehow ended up on the left
	if arg0.Op.IsConst() && !arg1.Op.IsConst() {
		if vop != op.Equal && vop != op.NotEqual {
			vop = vop.Opposite()
		}
		arg0, arg1 = arg1, arg0
	}

	blk.ReplaceControl(0, arg0)
	blk.InsertControl(-1, arg1)
	val.Remove()

	switch vop {
	case op.Less:
		blk.Op = op.IfLess
	case op.Greater:
		blk.Op = op.IfGreater
	case op.LessEqual:
		blk.Op = op.IfLessEqual
	case op.GreaterEqual:
		blk.Op = op.IfGreaterEqual
	case op.Equal:
		blk.Op = op.IfEqual
	case op.NotEqual:
		blk.Op = op.IfNotEqual
	}

	return 1
}

var _ = addToPass(Elaboration, rollupCompareToBlockOp)

func AddReturnMoves(fn *ir.Func) {
	blks := fn.Blocks()
	blk := blks[len(blks)-1]
	if blk.NumControls() > 3 {
		log.Panicf("Returning more than 3 values in %s is not yet supported", blk.Func())
	}
	for i := 0; i < blk.NumControls(); i++ {
		ctrl := blk.Control(i)
		if ctrl.Reg != reg.ArgRegs[i] {
			val := ir.BuildAt(blk, -1).Op(op.Copy, ctrl.Type, ctrl).PrevVal()
			val.Reg = reg.ArgRegs[i]
		}
	}
}
