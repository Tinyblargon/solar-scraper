package flags

import (
	"flag"
	"fmt"
	"os"
)

const (
	configDefault     string = ""
	configDescription string = "Config file path."
	debugDefault      bool   = false
	debugDescription  string = "Debug logging."
	helpDefault       bool   = false
	helpDescription   string = "Show help."
	logDefault        string = ""
	logDescription    string = "Log file path."
)

var config = flag.String("config", configDefault, configDescription)
var debug = flag.Bool("debug", debugDefault, debugDescription)
var help = flag.Bool("help", helpDefault, helpDescription)
var log = flag.String("log", logDefault, logDescription)

func init() {
	flag.StringVar(config, "c", configDefault, configDescription)
	flag.BoolVar(debug, "d", debugDefault, debugDescription)
	flag.BoolVar(help, "h", helpDefault, helpDescription)
	flag.StringVar(log, "l", logDefault, *log)
}

func usage() {
	fmt.Println("Usage: solar-scraper [options]")
	fmt.Println()
	fmt.Println("Options:")
	println("c", "config", configDescription)
	println("d", "debug", debugDescription)
	println("h", "help", helpDescription)
	println("l", "log", logDescription)
}

func println(shortFlag, longFlag, Description string) {
	fmt.Println("-" + shortFlag + ", --" + longFlag + "\t\t" + Description)
}

func Parse() Options {
	flag.Usage = usage
	flag.Parse()
	if *help {
		flag.Usage()
		os.Exit(0)
	}
	return Options{
		Config: *config,
		Debug:  *debug,
		Log:    *log,
	}
}

type Options struct {
	Config string
	Debug  bool
	Log    string
}
