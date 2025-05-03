package repl

import (
	"bufio"
	"fmt"
	"github.com/kkumar-gcc/enumgen/src/parser"
	"io"
	"strings"

	"github.com/kkumar-gcc/enumgen/src/lexer"
)

const PROMPT = ">> "

// Start runs a multi-line REPL. Input is collected until an empty line is entered.
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	for {
		// collect lines
		var lines []string
		for {
			fmt.Print(PROMPT)
			if !scanner.Scan() {
				return // EOF
			}
			line := scanner.Text()

			if strings.TrimSpace(line) == "" {
				break
			}
			lines = append(lines, line)
		}

		src := strings.Join(lines, "\n")
		l := lexer.New([]byte(src), lexer.CommentMode)
		//for {
		//	pos, tok, lit := l.Lex()
		//	if tok == token.EOF {
		//		break
		//	}
		//	fmt.Printf("%s\t%s\t%s\n", pos, tok, lit)
		//}
		p := parser.New(l)

		program := p.Parse()
		if p.Errors().Len() != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		io.WriteString(out, program.String())
		io.WriteString(out, "\n")
	}
}

func printParserErrors(out io.Writer, errors lexer.ErrorList) {
	io.WriteString(out, "Woops! We ran into some parser errors!\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg.Error()+"\n")
	}
}
