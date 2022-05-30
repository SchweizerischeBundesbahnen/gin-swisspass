// +build mage

package main

import (
	"fmt"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Default target to run when none is specified
// If not set, running mage will list available targets
//var Default = Build

// runs tests, generates code and builds the executable
func Cibuild() {
	mg.Deps(Test)
	// we only need tests, as this project is only used as a library, there is no executable to build
}

// runs all tests and creates a coverage report
func Test() error {
	fmt.Println("Start tests...")
	defer fmt.Println("Finished tests")
	return sh.RunV("go", "test", "-coverprofile=../coverage.out", "../...")
}
