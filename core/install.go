package core

import (
	"fmt"
	"os"

	. "github.com/amsa/doop/common"
)

func (doop *Doop) install() {
	if _, err := os.Stat(doop.homeDir); err == nil {
		fmt.Println("Doop is ready! You can create a new Doop project.")
		return
	}

	Debug("Doop directory does not exist. Creating doop directory at %s...", doop.homeDir)
	HandleError(os.Mkdir(doop.homeDir, 0755)) // Create Doop home directory
}
