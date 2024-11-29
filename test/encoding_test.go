package test

import (
	"github.com/adnpa/go-utils/pkg/basic"
	"testing"
)

type MyStruct struct {
	Name  string
	value string
}

func (m *MyStruct) GetValue() string {
	return m.value
}

func (m *MyStruct) SetValue(val string) {
	m.value = val
}

func TestJson(t *testing.T) {

}

func TestYaml(t *testing.T) {

}

func TestToml(t *testing.T) {

}

func TestGob(t *testing.T) {
	s := MyStruct{"abc", "abc"}
	bytes, err := basic.MarshalGob(s)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(bytes))

	var result MyStruct
	err = basic.UnmarshalGob(bytes, &result)
	if err != nil {
		return
	}
	t.Log(result)
}

func TestMsPack(t *testing.T) {
	s := MyStruct{"abc", "abc"}
	bytes, err := basic.MarshalMsPack(s)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(bytes))

	var result MyStruct
	err = basic.UnmarshalMsPack(bytes, &result)
	if err != nil {
		return
	}
	t.Log(result)
	t.Log(result.GetValue())

}
