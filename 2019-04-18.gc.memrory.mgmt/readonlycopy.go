package main

type verybig [512 * 1024 * 1024]byte

func test(data verybig) {
	println(len(data))
}

func main() {
	test(verybig{})
}
