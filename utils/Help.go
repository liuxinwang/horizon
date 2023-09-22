package utils

import (
	"flag"
	"github.com/go-demo/version"
	"os"
)

type Help struct {
	printVersion bool
}

func HelpInit() *Help {
	var help Help
	flag.BoolVar(&help.printVersion, "version", false, "print program build version")
	flag.Parse()
	// 这个需要放在第一个判断
	if help.printVersion {
		version.PrintVersion()
		os.Exit(0)
	}
	return &help
}
