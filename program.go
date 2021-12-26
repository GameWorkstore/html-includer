package main

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/otiai10/copy"
	"github.com/subchen/go-xmldom"
)

func main() {
	ok, destiny, ignoreFolders := getArguments()
	if !ok {
		printExplanation()
		return
	}
	patchHTMLFiles(destiny, ignoreFolders)
}

func getArguments() (bool, string, []string) {
	var source string
	var destiny string
	var ignoreFolders []string
	for i, arg := range os.Args {
		switch i {
		case 0:
			continue
		case 1:
			source = arg
			if has, _ := directoryExists(source); !has {
				return false, "", nil
			}
			continue
		case 2:
			destiny = arg
			if has, _ := directoryExists(destiny); has {
				err := os.RemoveAll(destiny)
				if err != nil {
					log.Fatal(err)
				}
			}
			if err := copy.Copy(source, destiny); err != nil {
				log.Fatal(err)
			}
			continue
		default:
			ignoreFolders = append(ignoreFolders, arg)
			continue
		}
	}
	return true, destiny, ignoreFolders
}

func directoryExists(filepath string) (bool, error) {
	_, err := os.Stat(filepath)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func find(root, ext string) []string {
	var a []string
	filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if filepath.Ext(d.Name()) == ext {
			a = append(a, s)
		}
		return nil
	})
	return a
}

func patchHTMLFiles(destiny string, folders []string) {
	for _, s := range find(destiny, ".html") {
		ignored := false
		for _, f := range folders {
			if strings.HasPrefix(s, f) {
				ignored = true
				break
			}
		}
		if ignored {
			continue
		}
		patchHTMLFile(s)
	}
}

func patchHTMLFile(filepath string) {
	println("Compiling ", filepath)

	content, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}
	doc := xmldom.Must(xmldom.ParseXML(string(content)))
	if err != nil {
		panic(err)
	}

	print(doc.XMLPretty())
	println("Done")
}
