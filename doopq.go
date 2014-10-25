//file: doopq.go

//Doop SQL command executor
//Sample:
//	doopq new_latte_branch@school_db "SELECT * FROM teachers WHERE name='Bob'"

package doopq

import (
	"fmt"
	"github.com/amsa/doop-core/core"
	"log"
	"os"
)
