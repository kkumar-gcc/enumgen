package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/kkumar-gcc/enumgen/src/ast"
	"github.com/kkumar-gcc/enumgen/src/lexer"
	"github.com/kkumar-gcc/enumgen/src/parser"
)

func main() {
	var (
		filePath string
		drawAST  bool
	)

	flag.StringVar(&filePath, "file", "example.edl", "Input EDL file")
	flag.BoolVar(&drawAST, "ast", false, "Generate AST dot and PNG")
	flag.Parse()

	src, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	l := lexer.New(src, lexer.CommentMode)
	p := parser.New(l)
	parsed := p.Parse()

	if len(p.Errors()) > 0 {
		for _, err := range p.Errors() {
			fmt.Println("Parse error:", err.Error())
		}
		return
	}

	if drawAST {
		dot := toDOT(parsed)
		err := os.WriteFile("ast.dot", []byte(dot), 0644)
		if err != nil {
			fmt.Println("Failed to write ast.dot:", err)
			return
		}

		cmd := exec.Command("dot", "-Tpng", "ast.dot", "-o", "ast.png")
		if err := cmd.Run(); err != nil {
			fmt.Println("Graphviz error:", err)
			fmt.Println("Make sure Graphviz is installed and `dot` is in your PATH.")
			return
		}
		fmt.Println("Generated ast.dot and ast.png")
	}
}

func toDOT(node ast.Node) string {
	var sb strings.Builder
	sb.WriteString("digraph AST {\n")
	sb.WriteString("  node [shape=box, fontname=\"Courier\"];\n")

	nodeID := 0
	getNodeID := func() int {
		nodeID++
		return nodeID
	}

	var traverse func(ast.Node) int
	var writeEdge func(parentID, childID int)

	writeEdge = func(parentID, childID int) {
		if childID != 0 {
			sb.WriteString(fmt.Sprintf("  %d -> %d;\n", parentID, childID))
		}
	}

	traverse = func(node ast.Node) int {
		if node == nil {
			return 0
		}
		id := getNodeID()

		switch n := node.(type) {
		case *ast.File:
			sb.WriteString(fmt.Sprintf("  %d [label=%s];\n", id, strconv.Quote("File")))
			for _, decl := range n.Declarations {
				child := traverse(decl)
				writeEdge(id, child)
			}

		case *ast.EnumDefinition:
			sb.WriteString(fmt.Sprintf("  %d [label=%s];\n", id, strconv.Quote("Enum: "+n.Name.Name)))
			writeEdge(id, traverse(&n.Name))
			if n.TypeSpec != nil {
				writeEdge(id, traverse(n.TypeSpec))
			}
			for _, m := range n.Members {
				writeEdge(id, traverse(m))
			}

		case *ast.TypeSpec:
			sb.WriteString(fmt.Sprintf("  %d [label=%s];\n", id, strconv.Quote("TypeSpec")))
			for _, t := range n.Types {
				writeEdge(id, traverse(t))
			}

		case *ast.TypeRef:
			sb.WriteString(fmt.Sprintf("  %d [label=%s];\n", id, strconv.Quote("TypeRef ("+n.Name.Name+")")))
			if n.Package != nil {
				writeEdge(id, traverse(n.Package))
			}
			writeEdge(id, traverse(&n.Name))

		case *ast.MemberDefinition:
			sb.WriteString(fmt.Sprintf("  %d [label=%s];\n", id, strconv.Quote("Member")))
			writeEdge(id, traverse(&n.Name))
			writeEdge(id, traverse(n.Value))

		case *ast.KeyValueExpr:
			sb.WriteString(fmt.Sprintf("  %d [label=%s];\n", id, strconv.Quote("KeyValueExpr")))
			writeEdge(id, traverse(n.Key))
			writeEdge(id, traverse(n.Value))

		case *ast.BasicLit:
			sb.WriteString(fmt.Sprintf("  %d [label=%s];\n", id, strconv.Quote("Literal: "+n.Value)))

		case *ast.Ident:
			sb.WriteString(fmt.Sprintf("  %d [label=%s];\n", id, strconv.Quote("Ident: "+n.Name)))

		default:
			sb.WriteString(fmt.Sprintf("  %d [label=%s];\n", id, strconv.Quote("Unknown")))
		}

		return id
	}

	traverse(node)
	sb.WriteString("}\n")
	return sb.String()
}
