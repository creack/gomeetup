package main

type tiny [8]byte

func test(data tiny) {
	println(len(data))
}

func main() {
	test(tiny{})
}
