package test

import (
	"fmt"
	"reflect"
	"testing"
)

func Test(t *testing.T) {
	var a = MyStruct{Name: "abc", value: "def"}
	//var b MyStruct
	//var m interface{} = a
	//b = m.(MyStruct)
	//fmt.Println(b)

	typ := reflect.TypeOf(a)
	val := reflect.ValueOf(a)
	fmt.Println(typ)
	fmt.Println(val)
	fmt.Println(typ.String())
	fmt.Println(typ.Kind())
	fmt.Println(val.Kind())
	//fmt.Println(typ.Elem())
	//fmt.Println(typ.Key())

}
