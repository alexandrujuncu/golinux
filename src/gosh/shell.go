package gosh

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func runExternal(args []string) (int, error) {
	fmt.Println("[DEBUG] Running external command ", args[0], "with arguments", args[1:])
	cmd := exec.Command(args[0])
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("[DEBUG] Forking into", args[0])
	err := cmd.Start()
	if err != nil {
		fmt.Println(err)
	} else {
		err = cmd.Wait()
		if err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println("[DEBUG] Back to gosh")
	if err != nil {
		return 1, err
	} else {
		return 0, err
	}
}

func run(args []string) (int, error) {
	/* Check if internal or built in command. Anything else is external.*/
	if isInternalCommand(args[0]) {
		return runInternalCommand(args)
	} else if isBuiltinCommand(args[0]) {
		return runBuiltinCommand(args)
	} else {
		return runExternal(args)
	}
}

func runCommandLine(token *ParserToken) {

}

func getPrompt() string {
	if os.Getuid() == 0 {
		return "# "
	} else {
		return "$ "
	}
}

func Shell() {
	fmt.Println("Welcome to Gosh!")

	/* Change into home directory, if set. */
	dir := os.Getenv("HOME")
	if strings.Compare(dir, "") == 0 {
		err := os.Chdir(dir)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	} else {
		/* If, by chance, PWD was set by parent, use it. */
		dir = os.Getenv("PWD")
		if strings.Compare(dir, "") == 0 {
			err := os.Chdir(dir)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		} else {
			/* Use / as current dir. */
			dir = "/"
			err := os.Chdir(dir)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}
	}
	/* Set PWD environment variable as directory above. */
	os.Setenv("PWD", dir)

	/* Set a default PATH value. */
	os.Setenv("PATH", "/bin")

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(getPrompt())

		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)

		if line != "" {
			token, _ := ParseCommandLine(line)
			fmt.Println("Processed", DebugParserTree(token))
			runCommandLine(token)
		}
	}
}
