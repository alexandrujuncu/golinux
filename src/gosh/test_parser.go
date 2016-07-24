package gosh

import (
	"fmt"
)

func DebugParserTree(token *ParserToken) string {
	if token.operator == OP_NONE {
		return fmt.Sprintf("{%v|IN=%v|OUT=%v|ERR=%v|BG=%v}", token.args, token.inFile, token.outFile, token.errFile, token.inBg)
	} else {
		return fmt.Sprintf("(%v %v %v)", DebugParserTree(token.left), token.operator, DebugParserTree(token.right))
	}

}

func TestParser() {
	fmt.Println("hello world!")

	var cmd [11]string

	cmd[0] = "vim  &       "
	cmd[1] = "ls -l; mkdir dir1 || mkdir dir2 || mkdir dir3; touch dir1/file && chmod 600 dir1/file; ls -l"
	cmd[2] = "cat /etc/passwd|grep root"
	cmd[3] = "echo test >/tmp/gost_test; cat /gost_test"
	cmd[4] = "true && false"
	cmd[5] = "true || false"
	cmd[6] = "echo test >/tmp/out"
	cmd[7] = "echo test 1>/tmp/outagain"
	cmd[8] = "echo test 2>/tmp/err"
	cmd[9] = "echo test &>/tmp/both"
	cmd[10] = "cat </tmp/src >/tmp/dst"

	for i := 0; i < len(cmd); i++ {
		c := make(chan *ParserToken, 1)
		fmt.Println("testing {", cmd[i], "}")
		doParseCommandLine(cmd[i], c)
		token := <-c
		fmt.Println("Processed", DebugParserTree(token))
	}
}
