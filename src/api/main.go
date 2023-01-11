package main

import (
	"fmt"
	"syscall/js"
)

const (
	rootEl = "__cryptomonyjsopaque__"
)

func ExportMe(this js.Value, args []js.Value) any {
	fmt.Println("I am writing from go !!!")
	return nil
}

func main() {
	done := make(chan bool)

	clMan := newClientManager()
	// cl := newClient()
	// sv := newServer()

	js.Global().Set(rootEl, make(map[string]interface{}))
	rootModule := js.Global().Get(rootEl)

	clMan.exposeToJS(rootModule)
	// sv.exposeServer(rootModule)

	<-done
}
