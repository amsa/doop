package main

import (
	"fmt"

	"github.com/amsa/doop-core/lib"
)

func main() {
	result := doop.Query("SELECT * FROM Towns;")
	fmt.Println(result)
}
