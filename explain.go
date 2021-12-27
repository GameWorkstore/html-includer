package main

import "fmt"

func printExplanation() {
	fmt.Println("htmlincluder program")
	fmt.Println("")
	fmt.Println("htmlincluder requires a source directory and destiny directory as arguments")
	fmt.Println("htmlincluder \\\n" +
		"\tabsolute/path/to/source \\\n" +
		"\tabsolute/path/to/destiny \\\n" +
		"\tabsolute/path/to/ignore1 \\\n" +
		"\tabsolute/path/to/ignore2 \\\n" +
		"\tabsolute/path/to/ignore3 ...")
}

func printArgs(source string, destiny string, ignoreFolders []string) {
	fmt.Println("Html Includer program")
	fmt.Println("SOURCE:", source)
	fmt.Println("DESTINY:", destiny)
	fmt.Println("IGNORED FOLDERS:")
	for _, f := range ignoreFolders {
		fmt.Println("\t", f)
	}
}
