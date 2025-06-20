package main

import (
    "fmt"
)

// Version is the current version of the provider
var Version = "1.0.0"

// PrintVersion prints the version information
func PrintVersion() {
    fmt.Printf("terraform-provider-oneuptime v%s\n", Version)
}
