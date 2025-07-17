package main

import "github.com/spf13/pflag"

var Version = ""

func main() {
	pflag.Parse()
	pflag.Usage = func() {}

	if Version == "" {
		Version = "unknown"
	}
}
