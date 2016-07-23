package goinit

import (
	"fmt"
	"goinit/login"
	"os"
	"syscall"
	"time"
)

func mountProcfs() error {
	/* Check if /proc exists and, if not, create it. */
	if _, err := os.Stat("/proc"); err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir("/proc", 0755)
			if err != nil {
				return err
			}
		}
	}

	return syscall.Mount("/proc", "proc", "proc", 0, "")
}

func mountSysfs() error {
	/* Check if /sys exists and, if not, create it. */
	if _, err := os.Stat("/sys"); err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir("/sys", 0755)
			if err != nil {
				return err
			}
		}
	}

	return syscall.Mount("/sys", "sys", "sysfs", 0, "")
}

func mountDevtmpfs() error {
	/* /dev should always be there, but just in case it's not, create it. */
	if _, err := os.Stat("/dev"); err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir("/dev", 0755)
			if err != nil {
				return err
			}
		}
	}

	return syscall.Mount("/dev", "dev", "devtmpfs", 0, "")
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
	var err error

	if os.Getpid() != 1 {
		fmt.Fprintln(os.Stderr, "Only the kernel can summon me!")
		os.Exit(1)
	}

	fmt.Println("Welcome to GoLinux!")

	/* Mount filesystems used to talk to the kernel. */

	/*
        Procfs is used by several tools to get information about the
	system. On Linux system is de facto mandatory.
        */
	err = mountProcfs()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	/* Sysfs is rather new, but most kernels should support it.*/
	err = mountSysfs()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	/*
        /dev should always be provided by the kernel, but not populated.
	We don't have udev in init so we try to use devtmpfs is availabe
	in the kernel.
	*/
	err = mountDevtmpfs()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	/* Provide multiuser login. */
	provideLoginPrompt()

	/*
		This point should never be reached since it would cause
		the kernel to panic.
	*/
	fmt.Println("Oh noes!")
	nothing()
}
