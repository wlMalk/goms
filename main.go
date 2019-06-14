package main

import (
	"fmt"
	goParser "go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/wlMalk/goms/generator"
	"github.com/wlMalk/goms/parser"
	"github.com/wlMalk/goms/version"

	"github.com/gookit/color"
	"github.com/vetcher/go-astra"
)

func main() {
	color.New(color.FgBlack, color.BgWhite, color.Bold).Printf("  GoMS  ")
	fmt.Printf(" v%s", version.VERSION)
	fmt.Println("")
	var err error
	defer func() {
		if err == nil {
			success(fmt.Sprintf("All files are successfully generated"))
		}
	}()
	defer func() {
		if err := recover(); err != nil {
			fail(fmt.Errorf("%s", err))
		}
	}()
	currentDir, err := os.Getwd()
	if err != nil {
		fail(err)
	}
	goPath := os.Getenv("GOPATH")
	if strings.TrimSpace(goPath) == "" {
		fail(fmt.Errorf("GOPATH is not defined"))
	}
	if !filepath.HasPrefix(currentDir, goPath) {
		fail(fmt.Errorf("service has to be located inside GOPATH"))
	}
	importPath, err := filepath.Rel(filepath.Join(os.Getenv("GOPATH"), "./src/"), currentDir)
	if err != nil {
		fail(err)
	}
	path := filepath.Join(currentDir, "./service.go")
	fset := token.NewFileSet()
	f, err := goParser.ParseFile(fset, path, nil, goParser.ParseComments|goParser.AllErrors)
	if err != nil {
		fail(fmt.Errorf("error when parse file: %v", err))
	}
	file, err := astra.ParseAstFile(f)
	if err != nil {
		if os.IsNotExist(err) {
			fail(fmt.Errorf("%s file does not exist", path))
		}
		fail(err)
	}

	services, err := parser.Parse(file)
	if err != nil {
		fail(err)
	}
	for _, service := range services {
		service.Path = filepath.Join(currentDir, "v"+service.Version.FullStringSpecial("."))
		service.ImportPath = filepath.ToSlash(filepath.Join(importPath, "v"+service.Version.FullStringSpecial(".")))
		files, err := generator.GenerateService(service)
		if err != nil {
			fail(err)
		}
		err = files.Save()
		if err != nil {
			fail(err)
		}
	}
}

func success(s string) {
	color.BgGreen.Print("  DONE  ")
	fmt.Print(" ")
	color.Green.Printf("%s\n", s)
}

func fail(err error) {
	color.BgRed.Print("  FAIL  ")
	fmt.Print(" ")
	color.Red.Printf("%s\n", err.Error())
	os.Exit(2)
}
