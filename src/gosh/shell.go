package gosh

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

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

func run_exit(args []string) {
	if len(args) > 1 {
		code, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		} else {
			os.Exit(int(code))
		}
	} else {
		os.Exit(0)
	}
}

func run_tellmemore() {
	fmt.Println("PID:", os.Getpid())
	fmt.Println("PPID:", os.Getppid())
	fmt.Println("User ID:", os.Getuid())
	fmt.Println("Group ID:", os.Getgid())
	fmt.Println("Effective User ID:", os.Geteuid())
	fmt.Println("Effective Group ID:", os.Getegid())
}

func run_external(args []string) {
	cmd := exec.Command(args[0])
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("Forking into", args[0])
	err := cmd.Start()
	if err != nil {
		fmt.Println(err)
	} else {
		err = cmd.Wait()
		if err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println("Back to gosh")
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
	case "exit":
		run_exit(args)
	case "tellmemore":
		run_tellmemore()
	default:
		run_external(args)
	}
}

func Shell() {
	fmt.Println("Welcome to Gosh!")
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
