package sizes

import "go/types"

var basicSizes = [...]byte{
	types.Bool:       1,
	types.Int8:       1,
	types.Int16:      1,
	types.Int32:      2,
	types.Int64:      3,
	types.Uint8:      1,
	types.Uint16:     1,
	types.Uint32:     2,
	types.Uint64:     4,
	types.Float32:    2,
	types.Float64:    4,
	types.Complex64:  4,
	types.Complex128: 8,
}

func Sizeof(T types.Type) int64 {
	switch t := T.Underlying().(type) {
	case *types.Basic:
		k := t.Kind()
		if int(k) < len(basicSizes) {
			if s := basicSizes[k]; s > 0 {
				return int64(s)
			}
		}
		if k == types.String {
			return 2
		}
	case *types.Array:
		n := t.Len()
		if n <= 0 {
			return 0
		}
		// n > 0
		z := Sizeof(t.Elem())
		return z * n
	case *types.Slice:
		return 3
	case *types.Struct:
		fields := Fieldsof(t)
		n := len(fields)
		if n == 0 {
			return 0
		}
		offsets := Offsetsof(fields)
		return offsets[n-1] + Sizeof(fields[n-1].Type())

	case *types.Interface:
		return 2
	}
	return 1 // catch-all
}

func Offsetsof(fields []*types.Var) []int64 {
	offsets := make([]int64, len(fields))
	var o int64
	for i, f := range fields {
		offsets[i] = o
		o += Sizeof(f.Type())
	}
	return offsets
}

func Fieldsof(t *types.Struct) []*types.Var {
	n := t.NumFields()
	if n == 0 {
		return nil
	}
	fields := make([]*types.Var, n)
	for i := 0; i < n; i++ {
		fields[i] = t.Field(i)
	}
	return fields
}
