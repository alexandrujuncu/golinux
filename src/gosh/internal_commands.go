package gosh

/* Implementation of gosh internal commands
(commands which would never be independent binaries).
*/

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

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
		return err
	} else {
		return err
	}
}

func run_pwd(args []string) error {
	dir, err := os.Getwd()
	fmt.Println(dir)
	return err
}

func run_env(args []string) {
	for _, e := range os.Environ() {
		fmt.Println(e)
	}
}
