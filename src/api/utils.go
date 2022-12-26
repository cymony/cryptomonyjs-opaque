package main

import (
	"fmt"
	"strings"
	"syscall/js"

	"github.com/cymony/cryptomony/opaque"
)

func rejectErr(reject js.Value, err error) {
	reject.Invoke(fmt.Sprintf("cryptomonyjs-opaque: %s", err.Error()))
}

func promiser(runner func(resolve js.Value, reject js.Value)) js.Value {
	handler := js.FuncOf(func(this js.Value, args []js.Value) any {
		resolve := args[0]
		reject := args[1]

		go runner(resolve, reject)

		return nil
	})

	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}

func checkIsString(input js.Value, argName string) error {
	if input.Type() != js.TypeString {
		return fmt.Errorf("%s argument must be string", argName)
	}
	return nil
}

func checkInputLen(inputs []js.Value, want int) error {
	if len(inputs) != want {
		return fmt.Errorf("inputs must be %d of length", want)
	}
	return nil
}

func copyBytesToGo(arr js.Value, argName string) ([]byte, error) {
	if err := checkArrType(arr, "Uint8Array", argName); err != nil {
		return nil, err
	}

	arrLen := arr.Get("length").Int()
	res := make([]byte, arrLen)
	js.CopyBytesToGo(res, arr)
	return res, nil
}

func copyBytesToJS(data []byte) js.Value {
	arrConstructor := js.Global().Get("Uint8Array")
	dataJS := arrConstructor.New(len(data))
	js.CopyBytesToJS(dataJS, data)
	return dataJS
}

func checkArrType(arr js.Value, typeStr string, argName string) error {
	typeDef := js.Global().Get("Object").Get("prototype").Get("toString").Call("call", arr)
	if !strings.Contains(typeDef.String(), typeStr) {
		return fmt.Errorf("%s argument must be %s", argName, typeStr)
	}
	return nil
}

func strToSuite(suiteStr string) (opaque.Identifier, error) {
	var s opaque.Identifier

	switch suiteStr {
	case string(ristretto255Suite):
		s = opaque.Ristretto255Suite
	case string(p256Suite):
		s = opaque.P256Suite
	default:
		return s, fmt.Errorf("first argument must be one of '%s' or '%s'", ristretto255Suite, p256Suite)
	}
	return s, nil
}
