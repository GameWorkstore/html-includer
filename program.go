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
	ok, source, destiny, ignoreFolders := getArguments()
	if !ok {
		printExplanation()
		return
	}
	printArgs(source, destiny, ignoreFolders)
	patchHTMLFiles(destiny, ignoreFolders)
}

func getArguments() (bool, string, string, []string) {
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
				return false, "", "", nil
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
	return true, source, destiny, ignoreFolders
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

	sContent := PreprocessFile(string(content))

	doc := xmldom.Must(xmldom.ParseXML(sContent))
	if err != nil {
		panic(err)
	}

	var htmlRoot = doc.Root.FindOneByName("html")
	patchHTMLNodeRecursive(parentPath, htmlRoot)
	RemoveHTMLIncluderJavascripts(htmlRoot)

	result := doc.XMLPretty()
	result = PostprocessFile(result)

	data := []byte(result)
	os.WriteFile(fpath, data, 0644)
	println("Done")
}

const attName = "html-include"

func patchHTMLNodeRecursive(parentPath string, htmlNode *xmldom.Node) {
	att := htmlNode.GetAttribute(attName)
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

			sContent := PreprocessFile(string(content))

			doc := xmldom.Must(xmldom.ParseXML(sContent))
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
		htmlNode.RemoveAttribute(attName)
	}

	for _, node := range htmlNode.Children {
		patchHTMLNodeRecursive(parentPath, node)
	}
}

func RemoveHTMLIncluderJavascripts(htmlNode *xmldom.Node) bool {
	if htmlNode.Name == "script" {
		if htmlNode.Text == "HtmlInclude();" {
			return true
		}
		if att := htmlNode.GetAttribute("src"); att != nil {
			if strings.HasSuffix(att.Value, "html-include.js") {
				return true
			}
		}
	}

	var nodes []*xmldom.Node
	for _, node := range htmlNode.Children {
		if RemoveHTMLIncluderJavascripts(node) {
			continue
		}
		nodes = append(nodes, node)
	}
	htmlNode.Children = nodes
	return false
}

const br_tag = "<br>"
const br_tag_alt = "$$$br$$$"
const brb_tag = "<br/>"
const brb_tag_alt = "$$$brb$$$"

// replace <br>, <br/> tags to something else.
func PreprocessFile(fileContent string) string {
	fileContent = strings.ReplaceAll(fileContent, br_tag, br_tag_alt)
	fileContent = strings.ReplaceAll(fileContent, brb_tag, brb_tag_alt)
	return fileContent
}

// restore <br>, <br/> tags.
func PostprocessFile(fileContent string) string {
	fileContent = strings.ReplaceAll(fileContent, br_tag_alt, br_tag)
	fileContent = strings.ReplaceAll(fileContent, brb_tag_alt, brb_tag)
	return fileContent
}
