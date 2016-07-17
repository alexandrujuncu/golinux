package gosh

import "syscall"

func run_reboot(args []string) error {
	/* Sync first to prevent data loss.*/
	syscall.Sync()
	return syscall.Reboot(syscall.LINUX_REBOOT_CMD_RESTART)
}

func run_halt(args []string) error {
	/* Sync first to prevent data loss.*/
	syscall.Sync()
	return syscall.Reboot(syscall.LINUX_REBOOT_CMD_HALT)
}

func run_poweroff(args []string) error {
	/* Sync first to prevent data loss.*/
	syscall.Sync()
	return syscall.Reboot(syscall.LINUX_REBOOT_CMD_POWER_OFF)
}
