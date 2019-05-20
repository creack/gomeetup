package main

func foo(x *int) {
	println(*x)
}

func stack() {
	x := 42
	foo(&x)
}

func main() {
	stack()
}
