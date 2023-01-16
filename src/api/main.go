package main

import (
	"syscall/js"
)

const (
	rootEl = "__cryptomonyjsopaque__"
)

func main() {
	done := make(chan bool)

	clMan := newClientManager()
	svMan := newServerManager()

	js.Global().Set(rootEl, make(map[string]interface{}))
	rootModule := js.Global().Get(rootEl)

	clMan.exposeToJS(rootModule)
	svMan.exposeServer(rootModule)

	<-done
}
