package state

import "math"
import . "luago/api"
import "luago/number"

type operator struct {
	metamethod  string
	integerFunc func(int64, int64) int64
	floatFunc   func(float64, float64) float64
}

var (
	iadd = func(a, b int64) int64 { return a + b }
	fadd = func(a, b float64) float64 { return a + b }
	isub = func(a, b int64) int64 { return a - b }
	fsub = func(a, b float64) float64 { return a - b }
	imul = func(a, b int64) int64 { return a * b }
	fmul = func(a, b float64) float64 { return a * b }
	fdiv = func(a, b float64) float64 { return a / b }
	iunm = func(a, _ int64) int64 { return -a }
	funm = func(a, _ float64) float64 { return -a }
	band = func(a, b int64) int64 { return a & b }
	bor  = func(a, b int64) int64 { return a | b }
	bxor = func(a, b int64) int64 { return a ^ b }
	bnot = func(a, _ int64) int64 { return ^a }
)

var operators = []operator{
	operator{"__add", iadd, fadd},
	operator{"__sub", isub, fsub},
	operator{"__mul", imul, fmul},
	operator{"__mod", number.IMod, number.FMod},
	operator{"__pow", nil, math.Pow},
	operator{"__div", nil, fdiv},
	operator{"__idiv", number.IFloorDiv, number.FFloorDiv},
	operator{"__band", band, nil},
	operator{"__bor", bor, nil},
	operator{"__bxor", bxor, nil},
	operator{"__shl", number.ShiftLeft, nil},
	operator{"__shr", number.ShiftRight, nil},
	operator{"__unm", iunm, funm},
	operator{"__bnot", bnot, nil},
}

// [-(2|1), +1, e]
// http://www.lua.org/manual/5.3/manual.html#lua_arith
func (self *luaState) Arith(op ArithOp) {
	var a, b luaValue
	if op != LUA_OPUNM && op != LUA_OPBNOT {
		b = self.stack.pop()
	} else {
		b = int64(0)
	}
	a = self.stack.pop()

	operator := operators[op]

	if operator.floatFunc == nil { // bitwise
		if x, y, ok := _convertToIntegers(a, b); ok {
			f := operator.integerFunc
			self.stack.push(f(x, y))
			return
		}
	} else {
		if f := operator.integerFunc; f != nil {
			if x, y, ok := _castToIntegers(a, b); ok {
				self.stack.push(f(x, y))
				return
			}
		}
		if x, y, ok := _convertToFloats(a, b); ok {
			f := operator.floatFunc
			self.stack.push(f(x, y))
			return
		}
	}

	if result, ok := callMetamethod(a, b, operator.metamethod, self); ok {
		self.stack.push(result)
		return
	}

	panic("todo: " + operator.metamethod)
}

/* integer or float */
/* float */
/* bitwise */

/* helper */

func _castToIntegers(a, b luaValue) (int64, int64, bool) {
	if x, ok := a.(int64); ok {
		if y, ok := b.(int64); ok {
			return x, y, true
		}
	}
	return 0, 0, false
}

func _convertToIntegers(a, b luaValue) (int64, int64, bool) {
	if x, ok := convertToInteger(a); ok {
		if y, ok := convertToInteger(b); ok {
			return x, y, true
		}
	}
	return 0, 0, false
}

func _convertToFloats(a, b luaValue) (float64, float64, bool) {
	if x, ok := convertToFloat(a); ok {
		if y, ok := convertToFloat(b); ok {
			return x, y, true
		}
	}
	return 0, 0, false
}
