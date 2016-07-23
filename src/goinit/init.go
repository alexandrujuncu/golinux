package goinit

import (
	"fmt"
	"goinit/login"
	"os"
	"time"
)

func nothing() {
	for {
		time.Sleep(time.Second)
	}
}

func provideLoginPrompt() {
	/* Endless loop providing login mechanism. */
	for {
		login.GetPrompt()
	}
}

func Init() {
	if os.Getpid() != 1 {
		fmt.Fprintln(os.Stderr, "Only the kernel can summon me!")
		os.Exit(1)
	}

	fmt.Println("Welcome to GoLinux!")

	/* Mount filesystems used to talk to the kernel. */

	earlyMounts()

	/* Provide multiuser login. */
	provideLoginPrompt()

	/*
		This point should never be reached since it would cause
		the kernel to panic.
	*/
	fmt.Println("Oh noes!")
	nothing()
}
