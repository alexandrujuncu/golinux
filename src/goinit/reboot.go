package goinit

import (
	"syscall"
)

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
