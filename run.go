package main

import (
	"fmt"
	"github.com/fatih/color"
	"os/exec"
)

func runAction() {
	for i := 0; i < len(deployHost); i++ {
		colorMsg(fmt.Sprintf("Deploy %s", deployHost[i]), color.FgHiBlue)
		doAction(deployHost[i])
	}
}

func doAction(myhost string) {
	switch action {
	case "cmd":
		colorMsg(fmt.Sprintf(">RUN: %s", command), color.FgHiGreen)
		doCommand(myhost)
	case "task":
	case "deploy":
	}
}

func doCommand(myhost string) {
	cmdline := fmt.Sprintf("ssh %s '%s'", myhost, command)
	out, _ := exec.Command("sh", "-c", cmdline).Output()
	colorMsg(">OUTPUT:", color.FgHiGreen)
	fmt.Printf("%s\n", out)
}

func doTask(myhost string, mytask string) {
}
