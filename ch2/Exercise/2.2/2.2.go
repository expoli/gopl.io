package main

import (
	"fmt"
	"gopl.io/ch2/Exercise/lengthconv"
	"os"
	"strconv"
)

func main() {
	for _, length := range os.Args[1:] {
		l, err := strconv.ParseFloat(length, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "length: %v\n", err)
			os.Exit(1)
		}
		f := lengthconv.Feet(l)
		m := lengthconv.Meters(l)

		fmt.Printf("%s = %s, %s = %s\n",
			f, lengthconv.FToM(f), m, lengthconv.MToF(m))
	}
}
