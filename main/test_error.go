package main

import (
	"errors"
	"fmt"
)

//
//
// 自己定义的 err 值类型
type errorString string

func (e errorString) Error() string {
	return string(e)
}
func New(text string) error {
	return errorString(text)
}

// 仿 error 库， 但不是返回地址
// 用于比较时，会展开struct 字段，比较字段，字段同则相同
type errorString2 struct {
	s string
}

func (e errorString2) Error() string {
	return e.s
}

func New2(text string) error {
	return errorString2{text}
}

// 自定义的 string 别名 的 error
var ErrNamedType = New("EOF")

// 自定义的值类型2
var ErrNamedStructType = New2("EOF")

// errors 库的 struct error, 返回的是地址
// 比较时，比较地址
var ErrStructType = errors.New("EOF")

type Value struct {
	st int
}

type TestValue struct {
	st int
	p  *Value
}

func main() {
	// 值相等，却不是同一个 err 了, 但被判定为相等
	if ErrNamedType == New("EOF") {
		// printed
		fmt.Println("Named Type Error")
	}
	// wraped struct 值类型, 值相同，但不是同一个 error,被判定为相等
	if ErrNamedStructType == New2("EOF") {
		// printed
		fmt.Println("ErrNamedStructType type error")
	}
	// 地址引用，不会相等
	if ErrStructType == errors.New("EOF") {
		// will not print
		fmt.Println("struct type error")
	}

	tv1 := TestValue{st: 1}
	tv2 := TestValue{st: 1}
	fmt.Println("当前情况：")
	fmt.Println(tv1, tv2)
	fmt.Println(&tv1, &tv2)
	if tv1 == tv2 {
		// 打印了， 值类型判断是每个字段判断相等
		fmt.Println("tv1 == tv2, 值类型判断是每个字段判断相等")
	}

	fmt.Println(&tv1, &tv2, "值虽同，地址不同，下个判断没打印")
	if &tv1 == &tv2 {
		// 没打印， 引用判断是地址
		fmt.Println("&tv1 == &tv2,没打印， 引用判断是地址")
	}
	// 值类型 内部有地址引用的时候
	v1 := Value{1}
	v2 := Value{2}
	tv3 := TestValue{st: 1, p: &v1}
	tv4 := TestValue{st: 1, p: &v2}
	fmt.Println("当前情况：")
	fmt.Println(tv3, tv4)
	fmt.Println(&tv3, &tv4)
	fmt.Println("值相等判断，已经不相等")
	if tv3 == tv4 {
		// 没打印， 有值不等
		fmt.Println("&tv3 == &tv4,没打印， 引用判断是地址")
	}
	fmt.Println("地址判断更不用说了，不相等")
	if &tv3 == &tv4 {
		// 没打印， 引用判断是地址
		fmt.Println("&tv3 == &tv4,没打印， 引用判断是地址")
	}
	// 不同的 struct 即使 字段全部相同，无法比较，会 panic
}
