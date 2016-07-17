package goinit

import (
	"fmt"
	"goinit/login"
	"os"
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
	mount_procfs()
	mount_sysfs()

	/* Provide multiuser login. */
	provideLoginPrompt()

	/*
		This point should never be reached since it would cause
		the kernel to panic.
	*/
	fmt.Println("Oh noes!")
	nothing()
}
