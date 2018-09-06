package main

import (
	"syscall/js"
)

func add(i []js.Value) {
	result := js.ValueOf(i[0].Int() + i[1].Int())
	i[2].Invoke(result)
}

func subtract(i []js.Value) {
	result := js.ValueOf(i[0].Int() - i[1].Int())
	i[2].Invoke(result)
}

func registerCallbacks() {
	js.Global().Set("add", js.NewCallback(add))
	js.Global().Set("subtract", js.NewCallback(subtract))
}

func main() {
	c := make(chan struct{}, 0)

	println("WASM Go Initialized")
	// register functions
	registerCallbacks()
	<-c
}
