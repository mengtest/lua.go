package state

import . "luago/api"

// [-0, +1, m]
// http://www.lua.org/manual/5.3/manual.html#lua_newthread
// lua-5.3.4/src/lstate.c#lua_newthread()
func (self *luaState) NewThread() LuaState {
	t := &luaState{registry: self.registry}
	t.pushLuaStack(newLuaStack(LUA_MINSTACK, t))
	self.PushThread(t)
	return t
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_status
// lua-5.3.4/src/lapi.c#lua_status()
func (self *luaState) Status() ThreadStatus {
	return self.coStatus
}

// [-?, +?, –]
// http://www.lua.org/manual/5.3/manual.html#lua_resume
func (self *luaState) Resume(from LuaState, nArgs int) ThreadStatus {
	lsFrom := from.(*luaState)
	if lsFrom.coChan == nil {
		lsFrom.coChan = make(chan int)
	}

	if self.coChan == nil {
		// start coroutine
		self.coChan = make(chan int)
		self.coCaller = lsFrom
		go func() {
			self.Call(nArgs, -1)
			lsFrom.coChan <- 1
		}()
	} else {
		// resume coroutine
		self.coStatus = LUA_OK
		self.coChan <- 1
	}

	<-lsFrom.coChan // wait coroutine to finish or yield
	return LUA_OK
}

// [-?, +?, e]
// http://www.lua.org/manual/5.3/manual.html#lua_yield
func (self *luaState) Yield(nResults int) int {
	self.coStatus = LUA_YIELD
	self.coCaller.coChan <- 1
	<-self.coChan
	return self.GetTop()
}

// [-?, +?, e]
// http://www.lua.org/manual/5.3/manual.html#lua_yieldk
func (self *luaState) YieldK() {
	panic("todo!")
}

// [-0, +0, –]
// http://www.lua.org/manual/5.3/manual.html#lua_isyieldable
func (self *luaState) IsYieldable() bool {
	panic("todo!")
}
