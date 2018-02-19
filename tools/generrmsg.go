package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-clang/v3.4/clang"
)

func main() {
	idx := clang.NewIndex(0, 1)
	defer idx.Dispose()

	tuArgs := []string{}
	if len(flag.Args()) > 0 && flag.Args()[0] == "-" {
		tuArgs = make([]string, len(flag.Args()[1:]))
		copy(tuArgs, flag.Args()[1:])
	}
	_ = tuArgs
	cmdArgs := []string{
		"-std=c99",
		"-I/usr/include/",
		"-I/usr/lib/gcc/x86_64-pc-linux-gnu/7.2.0/include"}
	tu := idx.ParseTranslationUnit("/usr/include/tox/tox.h", cmdArgs, nil, 0)
	defer tu.Dispose()

	for _, d := range tu.Diagnostics() {
		log.Println("PROBLEM:", d.Spelling())
	}

	fmt.Fprintf(os.Stdout, "package tox\n")
	fmt.Fprintf(os.Stdout, "/*\n")
	fmt.Fprintf(os.Stdout, "#include \"tox/tox.h\"\n")
	fmt.Fprintf(os.Stdout, "*/\n")
	fmt.Fprintf(os.Stdout, "import \"C\"\n")

	var enumDecls = make(map[uint32]bool)
	var enumConstDecls = make(map[uint32]bool)
	var curEnumName string

	tuc := tu.TranslationUnitCursor()
	tuc.Visit(func(cursor, parent clang.Cursor) (status clang.ChildVisitResult) {
		switch cursor.Kind() {
		case clang.Cursor_EnumDecl:
			//  log.Println(cursor.BriefCommentText())
			if !strings.HasPrefix(cursor.Spelling(), "TOX_ERR_") {
				break
			}
			if _, ok := enumDecls[cursor.HashCursor()]; ok {
				break
			}
			enumDecls[cursor.HashCursor()] = true
			curEnumName = cursor.Spelling()[3:]
			log.Println(cursor.Type().Kind().String(), cursor.Type().Spelling())
			fmt.Fprintf(os.Stdout, "\nvar %sS = make(map[int]string)\n", cursor.Spelling()[3:])
			fmt.Fprintf(os.Stdout, "func init(){%sS[%d] = \"TE%02d: %s\"}\n", curEnumName, -1, -1, curEnumName)

		case clang.Cursor_EnumConstantDecl:
			if !strings.HasPrefix(cursor.Spelling(), "TOX_ERR_") {
				break
			}
			if _, ok := enumConstDecls[cursor.HashCursor()]; ok {
				break
			}
			enumConstDecls[cursor.HashCursor()] = true

			// log.Println(cursor.BriefCommentText())
			// log.Println(cursor.Kind().String(), cursor.Type().Kind().String(), cursor.Type().Spelling(), cursor.Spelling())
			constComment := cursor.BriefCommentText()
			constComment = strings.Replace(constComment, "\"", "\\\"", -1)
			constName := cursor.Spelling()
			enumNameFields := strings.Split(cursor.Spelling(), "_")
			enumName := strings.Join(enumNameFields[1:len(enumNameFields)-1], "_")
			enumName = curEnumName
			constValue := cursor.EnumConstantDeclValue()
			if false {
				log.Println(enumName, constName, constValue, constComment)
			}
			fmt.Fprintf(os.Stdout, "const %s = int(C.%s) // %d\n", constName[4:], constName, constValue)
			fmt.Fprintf(os.Stdout, "func init(){%sS[%s] = \"TE%02d: %s\"}\n", enumName, constName[4:], constValue, constComment)

		default:
			// log.Println(cursor.Kind().String(), cursor.Type().Kind().String(), cursor.Type().Spelling(), cursor.Spelling())
		}

		return clang.ChildVisit_Recurse
	})
}

/*
结构：
var _enumName_S map[int]string
*/
