package gosh

import (
	"strings"
)

type ParserToken struct {
	args     []string
	operator OP
	left     *ParserToken
	right    *ParserToken
	inFile   string
	outFile  string
	errFile  string
	inBg     bool
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

func ParseSingleCommand(token *ParserToken, command string) {
	/* Split command into individual (sub)tokens separated by space.*/
	subtokens := strings.Split(command, " ")
	for i := range subtokens {
		subtoken := strings.TrimSpace(subtokens[i])
		/*
			Check for command operators. We assume that the
			operators will be at begining of the subtoken and
			filename imediately after.
		*/
		if index := strings.Index(subtoken, "<"); index == 0 {
			token.inFile = subtoken[index+1:]
		} else if index := strings.Index(subtoken, "1>"); index == 0 {
			token.outFile = subtoken[index+2:]
		} else if index := strings.Index(subtoken, "2>"); index == 0 {
			token.errFile = subtoken[index+2:]
		} else if index := strings.Index(subtoken, "&>"); index == 0 {
			token.outFile = subtoken[index+2:]
			token.errFile = subtoken[index+2:]
		} else if index := strings.Index(subtoken, ">"); index == 0 {
			token.outFile = subtoken[index+1:]
		} else if strings.Compare(subtoken, "&") == 0 {
			/* If no file redirection operatators, check to
			see if there is a background operator.
			*/
			token.inBg = true
		} else {
			/*
				If no argument in subtoken, treat as command argument.
			*/
			token.args = append(token.args, subtoken)
		}
	}
}

func doParseCommandLine(command string, c chan *ParserToken) {
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

		go doParseCommandLine(left, left_c)
		go doParseCommandLine(right, right_c)

		token.left = <-left_c
		token.right = <-right_c
		token.operator = operator

	} else {
		token.operator = OP_NONE
		/* Parse single command options. */
		command = strings.TrimSpace(command)
		ParseSingleCommand(token, command)
	}
	c <- token
}

func ParseCommandLine(command string) (*ParserToken, error) {
	c := make(chan *ParserToken, 1)
	doParseCommandLine(command, c)
	token := <-c
	return token, nil
}
