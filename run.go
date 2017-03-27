package main

import (
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"golang.org/x/crypto/ssh/terminal"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"syscall"
)

var (
	localtask map[string]bool //local task run one time only
	passwd    string
)

func runAction() {
	//if yamlhost == false {
	colorMsg("Password (enter if no passwd): ", color.FgHiYellow)
	bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
	passwd = string(bytePassword)
	//}
	localtask = make(map[string]bool)
	for i := 0; i < len(deployHost); i++ {
		colorMsg(fmt.Sprintf("\n[%s]\n", deployHost[i]), color.FgHiBlue)
		doAction(deployHost[i])
	}
}

func doAction(myhost string) {
	switch action {
	case "cmd":
		doCommand(myhost)
	case "task":
		doTask(myhost, command)
	case "deploy":
		n, _ := strconv.Atoi(deploySetting["Deployflow.size"])
		for i := 0; i < n; i++ {
			doTask(myhost, deploySetting[fmt.Sprintf("Deployflow.%d", i)])
		}
	}
}

func doCommand(myhost string) {
	colorMsg(fmt.Sprintf(">RUN: cmd( %s )\n", command), color.FgHiGreen)
	cmd := ""
	cmdline := ""
	if myhost == "localhost" || myhost == "127.0.0.1" {
		cmd = fmt.Sprintf("%s", command)
	} else {
		cmd = fmt.Sprintf("sshpass -p '%s' ssh -o ConnectTimeout=3 %s '%s'", passwd, myhost, command)
	}
	sudocmd := ""
	sudocmd = fmt.Sprintf("echo %s|sudo -kS", passwd)
	cmdline = strings.Replace(cmd, "sudo", sudocmd, -1)
	run(myhost, cmdline)
}

func doTask(myhost string, mytask string) {
	colorMsg(fmt.Sprintf(">RUN: task( %s )\n", mytask), color.FgHiGreen)
	cmdline := ""
	re := regexp.MustCompile("^([a-zA-Z]+).")
	node := fmt.Sprintf("%s", re.FindStringSubmatch(mytask)[1])
	switch node {
	case "local":
		if localtask[mytask] == true {
			colorMsg(fmt.Sprintf(">INFO: Ran Task( %s )\n", mytask), color.FgHiGreen)
			return
		}
		re = regexp.MustCompile("{!REMOTE!}")
		if re.MatchString(deployTask[mytask]) == false {
			localtask[mytask] = true
		}
		cmdline = strings.Replace(deployTask[mytask], "{!REMOTE!}", myhost, -1)
	case "remote":
		cmdline = fmt.Sprintf("sshpass -p '%s' ssh -o ConnectTimeout=3 %s '%s'", passwd, myhost, cmdline)
	default:
		colorMsg(fmt.Sprintf(">ERROR: No Task( %s )\n", mytask), color.FgHiRed)
		return
	}
	if debugMode {
		if node == "local" {
			colorMsg(fmt.Sprintf("-- %s --\n", cmdline), color.FgYellow)
		} else {
			colorMsg(fmt.Sprintf("-- %s --\n", deployTask[mytask]), color.FgYellow)
		}
	}
	run(myhost, cmdline)
}

func run(myhost string, cmdline string) {
	if cmdline == "" {
		colorMsg("ERROR: empty command!\n", color.FgHiRed)
		return
	}

	cmd := exec.Command("sh", "-c", cmdline)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		colorMsg(fmt.Sprintf("ERROR: %s\n", stderr.String()), color.FgHiRed)
		return
	}
	colorMsg("Result:\n", color.FgHiGreen)
	fmt.Println(out.String())
}
