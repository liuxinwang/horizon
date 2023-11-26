package utils

import (
	"flag"
	"github.com/go-demo/version"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

type Help struct {
	printVersion bool
	ConfigFile   string
}

func HelpInit() *Help {
	var help Help
	flag.StringVar(&help.ConfigFile, "config", "", "horizon config file")
	flag.BoolVar(&help.printVersion, "version", false, "print program build version")
	flag.Parse()
	// 这个需要放在第一个判断
	if help.printVersion {
		version.PrintVersion()
		os.Exit(0)
	}
	if help.ConfigFile == "" {
		log.Infof("-config param does not exist!")
		os.Exit(0)
	} else {
		abs, err := filepath.Abs(help.ConfigFile)
		if err != nil {
			log.Fatal("-config abs error: ", err.Error())
		}
		help.ConfigFile = abs
	}
	return &help
}
