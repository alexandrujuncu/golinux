package gosh

/*
	Linkings to builtin commands implementations
	(commands which could be independent binaries but
	are compiled into gosh binary).
*/

import (
	"errors"
	"fmt"
)

var builtinCommands = map[string]func(args []string) (int, error){
	"ls":       run_ls,
	"cat":      run_cat,
	"reboot":   run_reboot,
	"halt":     run_halt,
	"poweroff": run_poweroff,
	"mkdir":    run_mkdir,
	"mount":    run_mount}

func isBuiltinCommand(command string) bool {
	_, ok := builtinCommands[command]
	return ok
}

func runBuiltinCommand(args []string) (int, error) {
	fmt.Println("[DEBUG] Running builtin command ", args[0], "with arguments", args[1:])
	if len(args) > 0 {
		if isBuiltinCommand(args[0]) {
			returnCode, err := builtinCommands[args[0]](args)
			return returnCode, err
		} else {
			return 1, errors.New("Not builtin command")
		}
	} else {
		return 1, errors.New("Null command")
	}
}
