package gosh

import "syscall"

func run_reboot(args []string) (int, error) {
	/* Sync first to prevent data loss.*/
	syscall.Sync()
	err := syscall.Reboot(syscall.LINUX_REBOOT_CMD_RESTART)
	if err != nil {
		return 1, err
	} else {
		return 0, err
	}
}

func run_halt(args []string) (int, error) {
	/* Sync first to prevent data loss.*/
	syscall.Sync()
	err := syscall.Reboot(syscall.LINUX_REBOOT_CMD_HALT)
	if err != nil {
		return 1, err
	} else {
		return 0, err
	}
}

func run_poweroff(args []string) (int, error) {
	/* Sync first to prevent data loss.*/
	syscall.Sync()
	err := syscall.Reboot(syscall.LINUX_REBOOT_CMD_POWER_OFF)
	if err != nil {
		return 1, err
	} else {
		return 0, err
	}
}
