package main

import (
	"fmt"
	"reflect"
)

func main() {
	type A struct {
		N int
	}

	a := A{N: 100}

	rv := reflect.ValueOf(&a).Elem()

	nv := rv.FieldByName("N")

	var n int = 300
	newV := interface{}(n)
	newVr := reflect.ValueOf(newV)

	nv.Set(newVr.Convert(nv.Type()))

	fmt.Printf("%v\n", a)
}
