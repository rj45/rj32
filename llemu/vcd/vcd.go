package vcd

import (
	"fmt"
	"io"
	"strconv"
	"time"
)

// VCD produces a Value Change Dump for use
// with gtkWave.
// To use this, just create a new VCD with
// NewVCD, Add* any variables to track,
// WriteHeader() when all variables are defined
// and initialized, then periodically call Step()
// to record any changes.
type VCD struct {
	vars   []vcdVar
	time   uint64
	nextID uint64

	wr io.Writer
}

// NewVCD creates a new VCD writing to wr
func NewVCD(wr io.Writer) *VCD {
	return &VCD{wr: wr}
}

type varType int

const (
	// a single bit
	bitVar varType = iota

	// a collection of bits
	wireVar

	// an integer of some sort
	integerVar

	// a string value
	strVar
)

type vcdVar struct {
	module string
	name   string
	id     string
	typ    varType
	width  int

	bit     *bool
	wire    *uint32
	integer *uint32
	str     *string

	// previous values for diffing
	pbit     bool
	pwire    uint32
	pinteger uint32
	pstr     string
}

func (v *vcdVar) write(wr io.Writer) {
	switch v.typ {
	case bitVar:
		val := 0
		if *v.bit {
			val = 1
		}
		fmt.Fprintf(wr, "%d%s\n", val, v.id)
	case wireVar:
		fmt.Fprintf(wr, "b%b %s\n", *v.wire, v.id)
	case integerVar:
		fmt.Fprintf(wr, "r%d %s\n", *v.integer, v.id)
	case strVar:
		fmt.Fprintf(wr, "s%s %s\n", *v.str, v.id)
	default:
		panic(fmt.Errorf("unknown var typ: %d", v.typ))
	}
}

// AddBit adds a vcd var to track changes in a single bit, useful for control signals
func (vcd *VCD) AddBit(module, name string, bit *bool) {
	vcd.vars = append(vcd.vars, vcdVar{
		module: module,
		name:   name,
		id:     vcd.genID(),
		typ:    bitVar,
		bit:    bit,
	})
}

// AddWire adds a vcd var to track changes in a bus/register
func (vcd *VCD) AddWire(module, name string, width int, wire *uint32) {
	vcd.vars = append(vcd.vars, vcdVar{
		module: module,
		name:   name,
		id:     vcd.genID(),
		typ:    wireVar,
		width:  width,
		wire:   wire,
	})
}

// AddInt adds a vcd var to track changes in a number
func (vcd *VCD) AddInt(module, name string, width int, integer *uint32) {
	vcd.vars = append(vcd.vars, vcdVar{
		module:  module,
		name:    name,
		id:      vcd.genID(),
		typ:     integerVar,
		integer: integer,
	})
}

// AddStr adds a vcd var to track changes in a string value
// NOTE: this is non-standard, but gtkWave supports it
func (vcd *VCD) AddStr(module, name string, str *string) {
	vcd.vars = append(vcd.vars, vcdVar{
		module: module,
		name:   name,
		id:     vcd.genID(),
		typ:    strVar,
		str:    str,
	})
}

func (vcd *VCD) genID() string {
	vcd.nextID++
	return strconv.FormatUint(vcd.nextID, 36)
}

// WriteHeader writes the VCD header. Call this *after*
// calling any Add* methods to register variables.
func (vcd *VCD) WriteHeader(timescale time.Duration, version string) {
	fmt.Fprintf(vcd.wr, "$version %s $end\n", version)
	fmt.Fprintf(vcd.wr, "$date %s\n $end\n", time.Now().Format(time.ANSIC))
	fmt.Fprintf(vcd.wr, "$timescale %s $end\n", timescale.String())
	fmt.Fprintf(vcd.wr, "$scope module TOP $end\n")

	modules := make(map[string][]*vcdVar)
	for i := range vcd.vars {
		v := &vcd.vars[i]
		modules[v.module] = append(modules[v.module], v)
	}
	for module, list := range modules {
		fmt.Fprintf(vcd.wr, "  $scope module %s $end\n", module)
		for _, v := range list {
			switch v.typ {
			case bitVar:
				fmt.Fprintf(vcd.wr, "    $var wire 1 %s %s $end\n", v.id, v.name)
			case wireVar:
				fmt.Fprintf(vcd.wr, "    $var wire %d %s %s $end\n", v.width, v.id, v.name)
			case integerVar:
				fmt.Fprintf(vcd.wr, "    $var integer %d %s %s $end\n", v.width, v.id, v.name)
			case strVar:
				fmt.Fprintf(vcd.wr, "    $var string 24 %s %s $end\n", v.id, v.name)
			default:
				panic(fmt.Errorf("unknown var typ: %d", v.typ))
			}
		}
		fmt.Fprintf(vcd.wr, "  $upscope $end\n")
	}

	fmt.Fprintf(vcd.wr, "$upscope $end\n")
	fmt.Fprintf(vcd.wr, "#0\n")
	for i := range vcd.vars {
		v := &vcd.vars[i]
		v.write(vcd.wr)
	}
}

// Step increments the time by deltaTime and
// emits all the changes to the defined variables
func (vcd *VCD) Step(deltaTime uint64) {
	vcd.time += deltaTime

	emittedTime := false

	for i := range vcd.vars {
		v := &vcd.vars[i]
		changed := false

		switch v.typ {
		case bitVar:
			if v.pbit != *v.bit {
				v.pbit = *v.bit
				changed = true
			}
		case wireVar:
			if v.pwire != *v.wire {
				v.pwire = *v.wire
				changed = true
			}
		case integerVar:
			if v.pinteger != *v.integer {
				v.pinteger = *v.integer
				changed = true
			}
		case strVar:
			if v.pstr != *v.str {
				v.pstr = *v.str
				changed = true
			}
		default:
			panic(fmt.Errorf("unknown var typ: %d", v.typ))
		}

		if changed {
			if !emittedTime {
				emittedTime = true
				fmt.Fprintf(vcd.wr, "#%d\n", vcd.time)
			}
			v.write(vcd.wr)
		}
	}
}
