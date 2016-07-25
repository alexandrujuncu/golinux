package gosh

/*
	Implementation of gosh internal commands
	(commands which would never be independent binaries).
*/

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var internalCommands = map[string]func(args []string) error{
	"cd":         run_cd,
	"pwd":        run_pwd,
	"env":        run_env,
	"exit":       run_exit,
	"tellmemore": run_tellmemore}

func run_cd(args []string) error {
	var dir string
	var err error
	oldDir, _ := os.Getwd()
	if len(args) > 1 {
		if strings.Compare(args[1], "-") == 0 {
			/* Change to previous directory. */
			dir = os.Getenv("OLDPWD")
			if strings.Compare(dir, "") == 0 {
				err = errors.New("OLDPWD not set")
			}
		} else {
			dir = args[1]
		}
	} else {
		/* Change into home directory (if env set). */
		dir = os.Getenv("HOME")
		if strings.Compare(dir, "") == 0 {
			err = errors.New("HOME not set")
		}
	}
	if err == nil {
		/* If no other error, try to change dir.*/
		err = os.Chdir(dir)
		if err == nil {
			/* If successful, update env variables. */
			newDir, _ := os.Getwd()
			os.Setenv("PWD", newDir)
			os.Setenv("OLDPWD", oldDir)
		}
	}
	return err
}

func run_pwd(args []string) error {
	dir, err := os.Getwd()
	if err == nil {
		fmt.Println(dir)
	}
	return err
}

func run_env(args []string) error {
	for _, e := range os.Environ() {
		fmt.Println(e)
	}
	return nil
}

func run_exit(args []string) error {
	if len(args) > 1 {
		code, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			return err
		} else {
			os.Exit(int(code))
		}
	} else {
		os.Exit(0)
	}
	/* This will actually never be reached. */
	return nil
}

func run_tellmemore(args []string) error {
	fmt.Println("PID:", os.Getpid())
	fmt.Println("PPID:", os.Getppid())
	fmt.Println("User ID:", os.Getuid())
	fmt.Println("Group ID:", os.Getgid())
	fmt.Println("Effective User ID:", os.Geteuid())
	fmt.Println("Effective Group ID:", os.Getegid())

	return nil
}

func isInternalCommand(command string) bool {
	_, ok := internalCommands[command]
	return ok
}

func runInternalCommand(args []string) (int, error) {
	fmt.Println("[DEBUG] Running internal command ", args[0], "with arguments", args[1:])
	if len(args) > 0 {
		if isInternalCommand(args[0]) {
			err := internalCommands[args[0]](args)
			if err != nil {
				return 1, err
			} else {
				return 0, err
			}
		} else {
			return 1, errors.New("Not internal command")
		}
	} else {
		return 1, errors.New("Null command")
	}
}
