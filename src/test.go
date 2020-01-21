package main

import (
	"fmt"

	"github.com/gogf/gf/frame/g"

	"github.com/gogf/gf"
)

func main() {
	fmt.Println("hello GF", gf.VERSION)
	g.Server().Run()
}
