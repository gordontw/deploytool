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
	fmt.Println(msg)
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
	flag.Parse()

	if (host == "") || (action != "deploy" && command == "") {
		flag.PrintDefaults()
		os.Exit(0)
	}
	conf = parseConfig.ParseYML(inputYAML)

	// if not make inital, will "panic: assignment to entry in nil map"
	deploySetting = make(map[string]string)
	deployTask = make(map[string]string)

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

	if len(deployHost) == 0 {
		hmap := strings.Split(host, ",")
		for i := 0; i < len(hmap); i++ {
			processHost(hmap[i])
		}
		yamlhost = false
	}

	if debugMode {
		colorMsg("-----HOST-----", color.FgHiRed)
		for i := 0; i < len(deployHost); i++ {
			fmt.Printf("%s\n", deployHost[i])
		}
		printLoop("-----SETTING-----", deploySetting)
		printLoop("-----TASKS-----", deployTask)
	}
	runAction()
}
