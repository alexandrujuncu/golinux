package gosh

import (
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
