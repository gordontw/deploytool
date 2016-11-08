package main

import (
	"fmt"
	"github.com/fatih/color"
	"golang.org/x/crypto/ssh/terminal"
	"os/exec"
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
	cmdline := fmt.Sprintf("sshpass -p %s ssh %s '%s'", passwd, myhost, command)
	out, _ := exec.Command("sh", "-c", cmdline).Output()
	colorMsg(">OUTPUT:", color.FgHiGreen)
	fmt.Printf("%s\n", out)
}

func doTask(myhost string, mytask string) {
}
