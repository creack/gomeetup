package main

import (
	"runtime"
	"unsafe"
)

var addr uintptr

type big [512 * 1024]byte

func test() {
	x := big{}
	addr = uintptr((unsafe.Pointer)(&x))
}

func test1() {
	test()
	println(len(*(*big)(unsafe.Pointer(addr))))
}

func main() {
	test1()
	runtime.GC()
	println(len(*(*big)(unsafe.Pointer(addr))))
}
