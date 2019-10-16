package main

import (
	"app"
	"assembly"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func processFile(filename string, f *ast.File, master *assembly.MasterTypeCollection, verbose bool) {
	if verbose {
		fmt.Println("----------------------------------------------------------------------------")
		fmt.Printf("file: %s\n", filename)
	}

	packageName := f.Name.Name

	for _, dec := range f.Decls {

		switch declartion := dec.(type) {
		case *ast.GenDecl:
			for _, spec := range declartion.Specs {
				switch specification := spec.(type) {
				case *ast.ImportSpec:
					//
					// dont need to handle import specifications right now
					//
				case *ast.TypeSpec:
					master.UpdateFromTypeSpec(packageName, specification)
				case *ast.ValueSpec:
					//
					// dont need to handle value specifications right now
					//
				default:
					fmt.Printf("ERROR: unknown general declaration %s\n", specification)
				}
			}
		case *ast.FuncDecl:
			master.UpdateFromFuncDecl(packageName, declartion)
		default:
			fmt.Printf("ERROR: unknown top declaration %s\n", dec)
		}
	}
}

func main() {

	outputFlag := flag.String("output", "", "the output folder for the PlantUML txt files")
	verboseFlag := flag.Bool("verbose", false, "indicates if to be verbose with the output while analyzing")
	rootFlag := flag.String("root", "", "where the program should look for the top of the Go source tree")

	flag.Parse()

	if len(*outputFlag) <= 0 {
		panic("no output path defined")
	}
	if len(*rootFlag) <= 0 {
		panic("no root path defined")
	}

	//outputFolder := "C:/Users/drp21/Documents/dev/gocodevis/output/"
	outputFolder := *outputFlag
	//root := "C:/Users/drp21/Documents/dev/couchbase/sync_gateway/"
	//root := "C:/Users/drp21/Documents/dev/gocodevis/src/"
	root := *rootFlag
	verbose := *verboseFlag

	master := assembly.NewMasterTypeCollection()
	master.SetVerbose(verbose)

	files := app.FindGoFilesToProcess(root, verbose)

	for len(files) > 0 {
		filename := files[0]
		files = files[1:]
		fset := token.NewFileSet()
		f, parseErr := parser.ParseFile(fset, filename, nil, 0)
		if parseErr != nil {
			panic(parseErr)
		}

		processFile(filename, f, master, verbose)
	}

	app.GeneratePlantUMLFiles(master, outputFolder, verbose)
}
