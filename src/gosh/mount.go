package gosh

import (
	"errors"
	"fmt"
	"io"
	"os"
	"syscall"
)

func getMounts() error {
	fd, err := os.Open("/proc/mounts")
	if err != nil {
		return err
	}
	_, err = io.Copy(os.Stdout, fd)
	if err != nil {
		return err
	}
	return nil
}

func run_mount(args []string) (int, error) {
	if len(args) == 1 {
		/* If no arguments given, print current mounts.
		Currently, just print contents of /proc/mounts.*/
		err := getMounts()
		if err != nil {
			return 1, err
		}
	} else if len(args) == 4 {
		err := syscall.Mount(args[1], args[2], args[3], 0, "")
		fmt.Println(err)
		return 1, err
	} else {
		err := errors.New("format not suported")
		fmt.Println(err)
		return 1, err
	}
	return 0, nil
}
