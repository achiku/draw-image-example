package main

import (
	"flag"
	"fmt"
)

var (
	src = flag.String("src", "", "source image")
	dst = flag.String("dst", "", "destination image")
)

func main() {
	fmt.Println("drawimg")
}
