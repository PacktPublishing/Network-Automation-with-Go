package main

import (
	"fmt"
)

func main() {
	d := "interpreted\nliteral"

	e := `raw
literal`

	fmt.Println(d)
	fmt.Println(e)
}
