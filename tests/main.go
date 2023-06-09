package main

import (
	"flag"
	"fmt"
	"time"
)

//go:noinline
func recursion(level, maxLevel int) int {
	if level > maxLevel {
		return level
	}
	return recursion(level+1, maxLevel)
}

//go:noinline
func NewTestFunc() int {
	//nothing
	print("NewTestFunc\n")
	return 100
}

// uretprobe挂载的目标函数
//
//go:noinline
func CountCC(maxLevel int) (a int) {
	a = NewTestFunc()
	fmt.Println(a)
	if a > 100 {
		return a
	}

	a = recursion(0, maxLevel)
	fmt.Printf("CountCC return :%d\n", a)
	return a
}

func main() {
	var maxLevel = flag.Int("l", 100, "max recursion level")
	flag.Parse()
	for {
		go CountCC(*maxLevel)
		time.Sleep(time.Second)
	}
}
