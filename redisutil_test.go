package redisutil

import "testing"

func TestDropKey(t *testing.T) {
	v := DropKey(3, "a")
	if !v {
		t.Error("失败")
	}
}

func TestGetStrValue(t *testing.T) {
	v, err := GetStrValue(0, 13)
	if err != nil {
		t.Error(err.Error())
	}
	println(v)
}

func TestHGetKeyFieldStrValue(t *testing.T) {

	res := HSetKeyFieldValue(4, 32, 34, 55)
	if !res {
		t.Error("设置值失败!")
	}
	v, err := HGetKeyFieldStrValue(4, 32, 34)
	if err != nil {
		t.Error(err.Error())
	}
	println(v)
}