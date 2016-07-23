package gosh

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

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

func run_mkdir(args []string) (int, error) {
	if len(args) == 1 {
		err := errors.New("missing operand")
		return 1, err
	} else {
		for i := 1; i < len(args); i++ {
			err := os.Mkdir(args[i], 0777)
			if err != nil {
				return 1, err
			}
		}
	}
	return 0, nil
}
