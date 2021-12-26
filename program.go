package main

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	xmldom "github.com/GameWorkstore/html-includer/xmldom"
	"github.com/otiai10/copy"
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
			if err := fileOrDirExists(source); err != nil {
				return false, "", nil
			}
			continue
		case 2:
			destiny = arg
			if err := fileOrDirExists(destiny); err != nil {
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

func fileOrDirExists(filepath string) error {
	_, err := os.Stat(filepath)
	return err
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

func patchHTMLFile(fpath string) {
	println("Compiling ", fpath)

	content, err := os.ReadFile(fpath)
	if err != nil {
		log.Fatal(err)
	}

	parentPath := filepath.Dir(fpath)

	doc := xmldom.Must(xmldom.ParseXML(string(content)))
	if err != nil {
		panic(err)
	}

	var htmlRoot = doc.Root.FindOneByName("html")
	patchHTMLNodeRecursive(parentPath, htmlRoot)

	result := doc.XMLPretty()
	result = strings.ReplaceAll(result, "<script src=\"scripts/html-include.js\"></script>", "")
	result = strings.ReplaceAll(result, "<script src=\"scripts/html-include.js\"/>", "")
	result = strings.ReplaceAll(result, "<script>HtmlInclude();</script>", "")

	data := []byte(result)
	os.WriteFile(fpath, data, 0644)
	println("Done")
}

func patchHTMLNodeRecursive(parentPath string, htmlNode *xmldom.Node) {
	att := htmlNode.GetAttribute("html-include")
	if att != nil {
		includer := parentPath
		if strings.HasPrefix(att.Value, "/") {
			includer += att.Value
		} else {
			includer += "/" + att.Value
		}

		if err := fileOrDirExists(includer); err == nil {
			println("Node ", htmlNode.Name, " includes ", includer)

			content, err := os.ReadFile(includer)
			if err != nil {
				log.Fatal(err)
			}

			doc := xmldom.Must(xmldom.ParseXML(string(content)))
			if doc.Root.Name == "html-include" {
				for _, node := range doc.Root.Children {
					htmlNode.AppendChild(node)
				}
			} else {
				htmlNode.AppendChild(doc.Root)
			}
		} else {
			println("ERROR: Node ", htmlNode.Name, " includes ", includer, "doesn't exists.")
		}
		htmlNode.RemoveAttribute("include-html")
	}
	for _, node := range htmlNode.Children {
		patchHTMLNodeRecursive(parentPath, node)
	}
}
