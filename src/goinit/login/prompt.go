package login

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func userShell(userEntry *User) {
	cmd := exec.Command(userEntry.shell, "")

	/* Set owner or child process. */
	cmd.SysProcAttr = &syscall.SysProcAttr{}
	cmd.SysProcAttr.Credential = &syscall.Credential{
		Uid: userEntry.uid,
		Gid: userEntry.gid}
	/* Enharit input, output, error console.*/
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	/*Setup environment. */
	os.Clearenv()
	env := os.Environ()
	env = append(env, fmt.Sprintf("USER=%s", userEntry.name))
	env = append(env, fmt.Sprintf("USERNAME=%s", userEntry.name))
	env = append(env, fmt.Sprintf("LOGNAME=%s", userEntry.name))
	env = append(env, fmt.Sprintf("HOME=%s", userEntry.homeDir))
	env = append(env, fmt.Sprintf("PWD=%s", userEntry.homeDir))
	env = append(env, fmt.Sprintf("SHELL=%s", userEntry.shell))
	cmd.Env = env

	cmd.Dir = userEntry.homeDir

	/* Process will wait for command to return.
	This is intentional for now, since there is no other job to be done.
	*/
	err := cmd.Start()
	if err != nil {
		fmt.Println(err)
	} else {
		err = cmd.Wait()
		if err != nil {
			fmt.Println(err)
		}
	}
}

func GetPrompt() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("login: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)
	fmt.Print("password: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	userEntry, err := getUserEntryByName(username)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else {
		if strings.Compare(userEntry.password, password) == 0 {
			userShell(&userEntry)
		} else {
			fmt.Fprintln(os.Stderr, "Authentication failure")
		}
	}
}
