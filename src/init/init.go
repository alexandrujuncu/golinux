package main

import (
	"fmt"
	"os"
	"os/exec"
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

func shell() {
	cmd := exec.Command("/bin/gosh", "")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	/* Process will wait for command to return.
	This is intentional for now, since there is no other job to be done.
	*/
	fmt.Println("Starting gosh")
	err := cmd.Start()
	if err != nil {
		fmt.Println(err)
	} else {
		err = cmd.Wait()
		if err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println("Back to init")
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
