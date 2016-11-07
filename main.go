package main

import (
	"./parseConfig"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"os"
	"regexp"
	"strconv"
)

// input
var (
	conf          map[string]string
	host          string
	task          string
	inputYAML     string
	debugMode     bool
	defaultConfig = "deploy.yml"
)

// deploy value
var (
	deployHost    map[string]string
	deploySetting map[string]string
	deployTask    map[string]string
)

// error msg
func colorMsg(msg string, c color.Attribute) {
	color.Set(c)
	fmt.Println(msg)
	color.Unset()
}

func printLoop(note string, loop map[string]string) {
	colorMsg(note, color.FgHiGreen)
	for k, v := range loop {
		fmt.Printf("%s=>%s\n", k, v)
	}
}

func init() {
	flag.StringVar(&host, "h", "", "Host Group")
	flag.StringVar(&task, "t", "", "run task")
	flag.BoolVar(&debugMode, "d", false, "Debug mode")
	flag.StringVar(&inputYAML, "f", defaultConfig, "YAML file/directory")
}

func main() {
	flag.Parse()

	if host == "" {
		colorMsg("ERROR: host group", color.FgHiRed)
		os.Exit(0)
	}
	conf = parseConfig.ParseYML(inputYAML)

	// if not make inital, will "panic: assignment to entry in nil map"
	deployHost = make(map[string]string)
	deploySetting = make(map[string]string)
	deployTask = make(map[string]string)

	for k, v := range conf {
		re := regexp.MustCompile("^([a-zA-Z]+).(.*)")
		key := fmt.Sprintf("%s", re.FindStringSubmatch(k)[1])
		element := fmt.Sprintf("%s", re.FindStringSubmatch(k)[2])

		if key == host {
			reg := regexp.MustCompile("Nodes.*")
			if reg.FindStringIndex(element) == nil { // setting
				deploySetting[element] = v
			} else { // hosts
				deployHost[element] = v
			}
		}
		if key == "TASKS" { // task
			deployTask[element] = v
		}
	}

	if len(deployHost) == 0 {
		deployHost["Nodes.size"] = fmt.Sprintf("%d", 1)
		deployHost["Nodes.0"] = host
	}

	if debugMode {
		colorMsg("-----HOST-----", color.FgHiGreen)
		s, _ := strconv.Atoi(deployHost["Nodes.size"])
		for i := 0; i < s; i++ {
			fmt.Println(deployHost["Nodes."+strconv.Itoa(i)])
		}
		printLoop("-----SETTING-----", deploySetting)
		printLoop("-----TASKS-----", deployTask)
	}
}
