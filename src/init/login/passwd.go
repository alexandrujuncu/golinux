package login

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

/* Locations of passwd file. */
const PASSWD_FILE = "/etc/passwd"

/* User instance based on common Unix pattern. See the following for detals:
http://man7.org/linux/man-pages/man5/passwd.5.html
*/
type User struct {
	name     string
	password string
	uid      uint32
	gid      uint32
	gecos    string
	homeDir  string
	shell    string
}

func parsePasswdFile(c chan User) {
	inFile, err := os.Open(PASSWD_FILE)
	defer inFile.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		userEntryLine := strings.Split(line, ":")
		uid64, _ := strconv.ParseUint(userEntryLine[2], 10, 64)
		uid := uint32(uid64)
		gid64, _ := strconv.ParseUint(userEntryLine[3], 10, 64)
		gid := uint32(gid64)
		userEntry := User{
			userEntryLine[0],
			userEntryLine[1],
			uid,
			gid,
			userEntryLine[4],
			userEntryLine[5],
			userEntryLine[6],
		}
		c <- userEntry
	}
	close(c)
}

func getUserEntryByName(name string) (User, error) {
	c := make(chan User)
	var userEntry User

	go parsePasswdFile(c)

	for userEntry = range c {
		if strings.Compare(userEntry.name, name) == 0 {
			return userEntry, nil
		}
	}
	return userEntry, errors.New("Username does not exist.")
}

func getUserEntryByUid(uid uint32) (User, error) {
	c := make(chan User)
	var userEntry User

	go parsePasswdFile(c)

	for userEntry = range c {
		if userEntry.uid == uid {
			return userEntry, nil
		}
	}
	return userEntry, errors.New("UID does not exist.")
}
