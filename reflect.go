package main

import (
	"fmt"
	"reflect"
	"time"
)

func main() {
	var t = time.Now()
	typ := reflect.TypeOf(t)
	fmt.Printf("%v\n", typ)

	v := reflect.New(typ).Interface()

	if t1, ok := v.(time.Time); ok {
		fmt.Printf("can cast time.Time %v\n", t1)
	}

	if t2, ok := v.(*time.Time); ok {
		fmt.Printf("can cast *time.Time %v\n", t2)
	}
}
