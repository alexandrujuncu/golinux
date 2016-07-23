package parser

import (
	"fmt"
	"strings"
)

type ParserToken struct {
	command  string
	operator OP
	left     *ParserToken
	right    *ParserToken
}

type OP int

const (
	OP_NONE OP = iota
	OP_SEQ
	OP_PAR
	OP_CONDZ
	OP_CONDNZ
	OP_PIPE
)

func ParseSimpleCommand(token *ParserToken) {

}

func ParseCommandLine(command string, c chan *ParserToken) {
	operator := OP_NONE
	var left, right string

	token := new(ParserToken)

	/* Search for possible operators. */
	for i := 0; i < len(command) && operator == OP_NONE; i++ {
		switch string(command[i]) {
		/* Possible operators: ;, &, &&, |, ||. */
		case ";":
			operator = OP_SEQ
			left = command[:i]
			right = command[i+1:]
		case "&":
			/* Check if & is actually &&. */
			if i+1 < len(command) {
				if string(command[i+1]) == "&" {
					operator = OP_CONDNZ
					left = command[:i]
					right = command[i+3:]
				}
				/* Implement background run separately.*/
			}
		case "|":
			/* Check if | is actually ||. */
			if i+1 < len(command) {
				if string(command[i+1]) == "|" {
					operator = OP_CONDZ
					left = command[:i]
					right = command[i+3:]
				} else {
					operator = OP_PIPE
					left = command[:i]
					right = command[i+1:]
				}
			}
		}
	}
	left = strings.TrimSpace(left)
	right = strings.TrimSpace(right)

	if operator != OP_NONE {
		/* Parse subcomponents. */
		left_c := make(chan *ParserToken, 1)
		right_c := make(chan *ParserToken, 1)

		go ParseCommandLine(left, left_c)
		go ParseCommandLine(right, right_c)

		token.left = <-left_c
		token.right = <-right_c
		token.operator = operator

	} else {
		token.command = strings.TrimSpace(command)
		token.operator = operator
	}
	c <- token
}

func DebugParserTree(token *ParserToken) string {
	if token.operator == OP_NONE {
		return fmt.Sprintf("[%v]", token.command)
	} else {
		return fmt.Sprintf("(%v %v %v)", DebugParserTree(token.left), token.operator, DebugParserTree(token.right))
	}

}

func Test() {
	fmt.Println("hello world!")

	var cmd [6]string

	cmd[0] = "vim  &       "
	cmd[1] = "ls -l; mkdir dir1 || mkdir dir2 || mkdir dir3; touch dir1/file && chmod 600 dir1/file; ls -l"
	cmd[2] = "cat /etc/passwd|grep root"
	cmd[3] = "echo test >/tmp/gost_test; cat /gost_test"
	cmd[4] = "true && false"
	cmd[5] = "true || false"

	for i := 0; i < len(cmd); i++ {
		c := make(chan *ParserToken, 1)
		fmt.Println("testing", cmd[i])
		ParseCommandLine(cmd[i], c)
		token := <-c
		fmt.Println("Processed", DebugParserTree(token))
	}
}
