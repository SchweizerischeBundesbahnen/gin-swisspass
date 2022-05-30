// +build ignore

package main

import (
	"github.com/magefile/mage/mage"
	"os"
)

// this is a starter for the mage-targets.
// It parses all files with "+build mage" in the header and adds their targets to the list.
// It also generates an executable that can be reused. By default it's saved to $HOME/magefile. Use the MAGEFILE_CACHE env var to define your preferred location.
func main() { os.Exit(mage.Main()) }
