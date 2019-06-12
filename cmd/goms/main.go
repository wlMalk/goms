package main

import (
	"fmt"
	goParser "go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/wlMalk/goms"
	"github.com/wlMalk/goms/generator"
	"github.com/wlMalk/goms/parser"
	"github.com/wlMalk/goms/parser/types"

	"github.com/vetcher/go-astra"
)

func main() {
	var service *types.Service
	var err error
	defer func() {
		if err == nil {
			fmt.Printf("GoMS v%s\n", goms.VERSION)
			fmt.Printf("All files are successfully generated for '%s' service\n", service.Name)
		}
	}()
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	version, err := parser.ParseVersion(filepath.Base(currentDir))
	if err != nil {
		log.Fatalln(err)
	}
	goPath := os.Getenv("GOPATH")
	if strings.TrimSpace(goPath) == "" {
		log.Fatalln("GOPATH is not defined")
	}
	if !filepath.HasPrefix(currentDir, goPath) {
		log.Fatalln("service has to be located inside GOPATH")
	}
	importPath, err := filepath.Rel(filepath.Join(os.Getenv("GOPATH"), "./src/"), currentDir)
	if err != nil {
		log.Fatalln(err)
	}
	path := filepath.Join(currentDir, "./service.go")
	fset := token.NewFileSet()
	f, err := goParser.ParseFile(fset, path, nil, goParser.ParseComments|goParser.AllErrors)
	if err != nil {
		log.Fatalln(fmt.Errorf("error when parse file: %v", err))
	}
	file, err := astra.ParseAstFile(f)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatalln(fmt.Errorf("%s file does not exist", path))
		}
		log.Fatalln(err)
	}

	service, err = parser.Parse(file)
	if err != nil {
		log.Fatalln(err)
	}
	service.Version = *version
	service.Path = currentDir
	service.ImportPath = importPath
	files, err := generator.GenerateService(service)
	if err != nil {
		log.Fatalln(err)
	}
	err = files.Save()
	if err != nil {
		log.Fatalln(err)
	}
	// t, err := json.MarshalIndent(service, "", " ")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(string(t))
	// fmt.Println(parser.CleanComments(file.Interfaces[0].Methods[0].Docs))
}
