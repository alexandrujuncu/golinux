package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"syscall"
	"time"
)

func mount_procfs() {
	var err error

	err = os.Mkdir("/proc", 0555)
	if err != nil {
		fmt.Println("Error creating /proc.", err)
	}
	err = syscall.Mount("/proc", "proc", "proc", 0, "")
	if err != nil {
		fmt.Println("Error mounting procfs.", err)
	}
}

func mount_sysfs() {
	var err error

	err = os.Mkdir("/sys", 0555)
	if err != nil {
		fmt.Println("Error creating /sys.", err)
	}
	err = syscall.Mount("/sys", "sys", "sysfs", 0, "")
	if err != nil {
		fmt.Println("Error mounting sysfs.", err)
	}
}

func run_reboot() {
	/* Sync first to prevent data loss.*/
	syscall.Sync()
	syscall.Reboot(syscall.LINUX_REBOOT_CMD_RESTART)
}

func run_halt() {
	/* Sync first to prevent data loss.*/
	syscall.Sync()
	syscall.Reboot(syscall.LINUX_REBOOT_CMD_HALT)
}

func run_poweroff() {
	/* Sync first to prevent data loss.*/
	syscall.Sync()
	syscall.Reboot(syscall.LINUX_REBOOT_CMD_POWER_OFF)
}

func run_sysc() {
	syscall.Sync()
}

func run_cd(args []string) {
	if len(args) > 1 {
		os.Chdir(args[1])
	}
}

func run_ls(args []string) {
	var dir string
	if len(args) > 1 {
		dir = args[1]
	} else {
		dir = "."
	}
	files, _ := ioutil.ReadDir(dir)
	for _, f := range files {
		fmt.Println(f.Name())
	}
}

func run_cat(args []string) {
	if len(args) == 1 {
		_, err := io.Copy(os.Stdout, os.Stdin)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	} else {
		for _, file_name := range args[1:] {
			fd, err := os.Open(file_name)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			_, err = io.Copy(os.Stdout, fd)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}
	}
}

func run(args []string) {
	fmt.Println("Running", args[0], "with arguments", args[1:])
	/* Check for internal commands */
	switch args[0] {
	case "cd":
		run_cd(args)
	case "ls":
		run_ls(args)
	case "cat":
		run_cat(args)
	case "reboot":
		run_reboot()
	case "halt":
		run_halt()
	case "poweroff":
		run_poweroff()
	default:
		fmt.Println("External command")
	}
}

func shell() {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("> ")

		line, _ := reader.ReadString('\n')
		line = strings.Trim(line, "\n")

		args := strings.Split(line, " ")
		if len(args) > 0 {
			run(args)
		} else {
			fmt.Println("null-command")
		}
	}
}

func nothing() {
	for {
		time.Sleep(time.Second)
	}
}

func main() {
	if os.Getpid() != 1 {
		fmt.Fprintln(os.Stderr, "Only the kernel can summon me!")
		nothing()
	}

	fmt.Println("Welcome to GoLinux!")

	/* Mount filesystems used to talk to the kernel. */
	mount_procfs()
	mount_sysfs()

	/* Start a minimal shell to interact with the kernel. */
	shell()

	/*
		This point should never be reached since it would cause
		the kernel to panic.
	*/
	fmt.Println("Oh noes!")
	nothing()
}
