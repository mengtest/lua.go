package state

import "fmt"
import . "luago/api"

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_pushnil
func (self *luaState) PushNil() {
	self.stack.push(nil)
}

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_pushboolean
func (self *luaState) PushBoolean(b bool) {
	self.stack.push(b)
}

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_pushinteger
func (self *luaState) PushInteger(n int64) {
	self.stack.push(n)
}

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_pushnumber
func (self *luaState) PushNumber(n float64) {
	self.stack.push(n)
}

// [-0, +1, m]
// http://www.lua.org/manual/5.3/manual.html#lua_pushstring
func (self *luaState) PushString(s string) {
	self.stack.push(s)
}

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_pushcfunction
func (self *luaState) PushGoFunction(f GoFunction) {
	self.stack.push(f)
}

// [-n, +1, m]
// http://www.lua.org/manual/5.3/manual.html#lua_pushcclosure
func (self *luaState) PushGoClosure(f GoFunction, n int) {
	if n == 0 {
		self.stack.push(f)
	} else { // closure
		vals := self.stack.popN(n)
		closure := &goClosure{f, vals}
		self.stack.push(closure)
	}
}

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_pushlightuserdata
func (self *luaState) PushUserData(d UserData) {
	ud := &userData{data: d}
	self.stack.push(ud)
}

// [-0, +1, e]
// http://www.lua.org/manual/5.3/manual.html#lua_pushfstring
func (self *luaState) PushFString(fmtStr string, a ...interface{}) {
	str := fmt.Sprintf(fmtStr, a...)
	self.stack.push(str)
}

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_pushglobaltable
func (self *luaState) PushGlobalTable() {
	global := self.registry.get(LUA_RIDX_GLOBALS)
	self.stack.push(global)
}

// [-0, +1, –]
// http://www.lua.org/manual/5.3/manual.html#lua_pushthread
func (self *luaState) PushThread(ls LuaState) bool {
	self.stack.push(ls)
	return false // todo
}
