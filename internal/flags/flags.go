package flags

import (
	"flag"
	"fmt"
	"os"
)

const (
	configDefault      string = ""
	configDescription  string = "Config file path."
	debugDefault       bool   = false
	debugDescription   string = "Debug logging."
	helpDefault        bool   = false
	helpDescription    string = "Show help."
	logDefault         string = ""
	logDescription     string = "Log file path."
	versionDefault     bool   = false
	versionDescription string = "Show package version."
)

var config = flag.String("config", configDefault, configDescription)
var debug = flag.Bool("debug", debugDefault, debugDescription)
var help = flag.Bool("help", helpDefault, helpDescription)
var log = flag.String("log", logDefault, logDescription)
var versionFlag = flag.Bool("version", versionDefault, versionDescription)

func init() {
	flag.StringVar(config, "c", configDefault, configDescription)
	flag.BoolVar(debug, "d", debugDefault, debugDescription)
	flag.BoolVar(help, "h", helpDefault, helpDescription)
	flag.StringVar(log, "l", logDefault, *log)
	flag.BoolVar(versionFlag, "v", versionDefault, versionDescription)
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

// Parse parses the command line arguments and returns the options
func Parse(version string) Options {
	flag.Usage = usage
	flag.Parse()
	if *help {
		flag.Usage()
		os.Exit(0)
	}
	if *versionFlag {
		fmt.Println(version)
		os.Exit(0)
	}
	return Options{
		Config: *config,
		Debug:  *debug,
		Log:    *log,
	}
}

// Options is the command line options
type Options struct {
	Config string
	Debug  bool
	Log    string
}
