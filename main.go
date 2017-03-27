package main

import (
	"./parseConfig"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// input
var (
	conf          map[string]string
	host          string
	action        string
	command       string
	inputYAML     string
	debugMode     bool
	defaultConfig = "deploy.yml"
)

// deploy value
var (
	yamlhost      = true
	deployHost    []string
	deploySetting map[string]string
	deployTask    map[string]string
)

// print colorize msg
func colorMsg(msg string, c color.Attribute) {
	color.Set(c)
	fmt.Print(msg)
	color.Unset()
}

func printLoop(note string, loop map[string]string) {
	colorMsg(note, color.FgHiRed)
	for k, v := range loop {
		fmt.Printf("%s=>%s\n", k, v)
	}
}

func processHost(inhost string) {
	reg := regexp.MustCompile("\\[" + "([0-9]+)-([0-9]+)" + "\\]")
	if reg.MatchString(inhost) {
		hmap := reg.FindStringSubmatch(inhost)
		start, _ := strconv.Atoi(hmap[1])
		end, _ := strconv.Atoi(hmap[2])
		for iter := start; iter <= end; iter++ {
			ihost := strings.Replace(inhost, hmap[0], fmt.Sprintf("%d", iter), 1)
			deployHost = append(deployHost, ihost)
		}
	} else {
		deployHost = append(deployHost, inhost)
	}
}

func init() {
	flag.StringVar(&host, "h", "", "[Require!]Yaml Host Group or input host ex. m10[1-3].11.2.3,127.0.0.1")
	flag.StringVar(&action, "a", "deploy", "action [deploy|task|cmd]")
	flag.StringVar(&command, "i", "", "run item")
	flag.BoolVar(&debugMode, "d", false, "Debug mode")
	flag.StringVar(&inputYAML, "f", defaultConfig, "YAML file/directory")
}

func main() {
	runtime := time.Now().Format("20060102150405")

	flag.Parse()
	if host == "" {
		flag.PrintDefaults()
		os.Exit(0)
	}
	if command == "" {
		action = "deploy"
	} else {
		mytype := regexp.MustCompile("^([a-zA-Z]+).").FindStringSubmatch(command)[1]
		if mytype == "local" || mytype == "remote" {
			action = "task"
		} else {
			action = "cmd"
		}
	}
	conf = parseConfig.ParseYML(inputYAML)

	// if not make inital, will "panic: assignment to entry in nil map"
	deploySetting = make(map[string]string)
	deployTask = make(map[string]string)

	deploySetting["Runtime"] = runtime
	for k, v := range conf {
		re := regexp.MustCompile("^([a-zA-Z]+).(.*)")
		key := fmt.Sprintf("%s", re.FindStringSubmatch(k)[1])
		element := fmt.Sprintf("%s", re.FindStringSubmatch(k)[2])

		if key == host {
			re = regexp.MustCompile("Nodes.*")
			if re.FindStringIndex(element) == nil { // setting
				deploySetting[element] = v
			} else if element != "Nodes.size" { // hosts
				deployHost = append(deployHost, v)
			}
		}
		if key == "TASKS" { // task
			deployTask[element] = v
		}
	}
	for k, v := range deployTask {
		re := regexp.MustCompile("{{[a-zA-Z]+}}")
		regmap := re.FindAllString(v, -1)
		for i := 0; i < len(regmap); i++ {
			v = strings.Replace(v, regmap[i], deploySetting[regmap[i][2:len(regmap[i])-2]], -1)
		}
		deployTask[k] = v
	}

	if len(deployHost) == 0 {
		hmap := strings.Split(host, ",")
		for i := 0; i < len(hmap); i++ {
			processHost(hmap[i])
		}
		yamlhost = false
	}

	if debugMode {
		colorMsg("-----HOST-----\n", color.FgHiRed)
		for i := 0; i < len(deployHost); i++ {
			fmt.Printf("%s\n", deployHost[i])
		}
		printLoop("-----SETTING-----\n", deploySetting)
		printLoop("-----TASKS-----\n", deployTask)
	}
	runAction()
}
