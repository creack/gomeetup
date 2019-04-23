package main

func heap() *int {
	x := 42
	return &x
}

func main() {
	println(*heap())
}
