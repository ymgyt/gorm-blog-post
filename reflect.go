package main

import (
	"fmt"
	"reflect"
)

func main() {
	var i int = 10

	rv := reflect.New(reflect.PtrTo(reflect.TypeOf(i)))
	fmt.Printf("%v %v\n", rv.Type(), rv.Kind())

	var ii int = 20
	rv.Elem().Set(reflect.ValueOf(&ii).Elem().Addr())

	v := rv.Interface().(**int)
	fmt.Printf("%#v\n", v)
	g := 200
	*v = &g

	fmt.Println(i)
}
