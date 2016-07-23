package goinit

import (
	"fmt"
	"os"
	"syscall"
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

func earlyMounts() {
	var err error

	/*
	        Procfs is used by several tools to get information about the
		system. On Linux system is de facto mandatory.
	*/
	err = mountProcfs()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error mounting /proc:", err)
	}

	/* Sysfs is rather new, but most kernels should support it.*/
	err = mountSysfs()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error mounting /sys:", err)
	}

	/*
	        /dev should always be provided by the kernel, but not populated.
		We don't have udev in init so we try to use devtmpfs is availabe
		in the kernel.
	*/
	err = mountDevtmpfs()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error mounting /dev:", err)
	}
}
