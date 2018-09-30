package main

import (
	"fmt"
	"os"

	"github.com/ansemjo/aenker/ae"
)

func main() {

	a := ae.Aenker2{}

	err := a.OpenHeader(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	fmt.Println(a)

}
