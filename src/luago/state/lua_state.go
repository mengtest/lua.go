package state

import . "luago/api"

type luaState struct {
	/* global state */
	panicf     GoFunction
	registry   *luaTable
	mtOfNil    *luaTable // ?
	mtOfBool   *luaTable
	mtOfNumber *luaTable
	mtOfString *luaTable
	mtOfFunc   *luaTable
	mtOfThread *luaTable
	/* stack */
	stack     *luaStack
	callDepth int
	/* coroutine */
	coStatus ThreadStatus
	coCaller *luaState
	coChan   chan int
}

func New() LuaState {
	registry := newLuaTable(8, 0)
	registry.put(LUA_RIDX_MAINTHREAD, nil) // todo
	registry.put(LUA_RIDX_GLOBALS, newLuaTable(8, 0))

	ls := &luaState{registry: registry}
	ls.pushLuaStack(newLuaStack(LUA_MINSTACK, ls))
	return ls
}

func (self *luaState) pushLuaStack(stack *luaStack) {
	stack.prev = self.stack
	self.stack = stack
	self.callDepth++
}

func (self *luaState) popLuaStack() {
	stack := self.stack
	self.stack = stack.prev
	stack.prev = nil
	self.callDepth--
}

// debug
func (self *luaState) String() string {
	return stackToString(self.stack)
}
