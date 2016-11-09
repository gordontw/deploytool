package main

import (
	"fmt"
	"github.com/fatih/color"
	"golang.org/x/crypto/ssh/terminal"
	"os/exec"
	"regexp"
	"strconv"
	"syscall"
)

var (
	passwd string
)

func runAction() {
	if yamlhost == false {
		fmt.Print("Enter Password: ")
		bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
		passwd = string(bytePassword)
	}
	for i := 0; i < len(deployHost); i++ {
		colorMsg(fmt.Sprintf("\n[Deploy %s]", deployHost[i]), color.FgHiBlue)
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
	colorMsg(fmt.Sprintf(">RUN: cmd( %s )", command), color.FgHiGreen)
	cmdline := ""
	if myhost == "localhost" || myhost == "127.0.0.1" {
		cmdline = fmt.Sprintf("%s", command)
	} else {
		cmdline = fmt.Sprintf("sshpass -p '%s' ssh -o ConnectTimeout=3 %s '%s'", passwd, myhost, command)
	}
	out, _ := exec.Command("sh", "-c", cmdline).Output()
	colorMsg(">OUTPUT:", color.FgHiGreen)
	fmt.Printf("%s\n", out)
}

func doTask(myhost string, mytask string) {
	colorMsg(fmt.Sprintf(">RUN: task( %s )", mytask), color.FgHiGreen)
	cmdline := ""
	re := regexp.MustCompile("^([a-zA-Z]+).")
	switch fmt.Sprintf("%s", re.FindStringSubmatch(mytask)[1]) {
	case "local":
		cmdline = fmt.Sprintf("%s", deployTask[mytask])
	case "remote":
		cmdline = fmt.Sprintf("sshpass -p '%s' ssh -o ConnectTimeout=3 %s '%s'", passwd, myhost, deployTask[mytask])
	default:
		colorMsg(fmt.Sprintf(">ERROR: No Task( %s )", mytask), color.FgHiRed)
		return
	}
	if debugMode {
		colorMsg(fmt.Sprintf("-- %s --", deployTask[mytask]), color.FgYellow)
	}
	out, err := exec.Command("sh", "-c", cmdline).Output()
	colorMsg(">OUTPUT:", color.FgHiGreen)
	fmt.Printf("%s\n", out)
	if err != nil {
		fmt.Println("ERR" + fmt.Sprint(err))
	}
}
