package gosh

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

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
		run_reboot(args)
	case "halt":
		run_halt(args)
	case "poweroff":
		run_poweroff(args)
	case "exit":
		run_exit(args)
	case "tellmemore":
		run_tellmemore()
	case "env":
		run_env(args)
	case "pwd":
		run_pwd(args)
	default:
		run_external(args)
	}
}

func getPrompt() string {
	if os.Getuid() == 0 {
		return "# "
	} else {
		return "$ "
	}
}

func Shell() {
	fmt.Println("Welcome to Gosh!")

	/* Change into home directory, if set. */
	dir := os.Getenv("HOME")
	if strings.Compare(dir, "") == 0 {
		err := os.Chdir(dir)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	} else {
		/* If, by change, PWD was set by parent, use it. */
		dir = os.Getenv("PWD")
		if strings.Compare(dir, "") == 0 {
			err := os.Chdir(dir)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		} else {
			/* User / as current dir. */
			dir = "/"
			err := os.Chdir(dir)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}
	}
	/* Set PWD environment variable as directory above. */
	os.Setenv("PWD", dir)

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(getPrompt())

		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)

		if line != "" {
			args := strings.Split(line, " ")
			run(args)
		}
	}
}
