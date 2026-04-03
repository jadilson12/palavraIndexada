package main

import (
	"fmt"
	"html"
	"os"
	"path/filepath"
	"strings"
)

// Entidades não cobertas por html.UnescapeString
var extraEntities = map[string]string{
	"&nbsp;":   "\u00A0",
	"&laquo;":  "\u00AB",
	"&raquo;":  "\u00BB",
	"&bull;":   "\u2022",
	"&middot;": "\u00B7",
	"&frac12;": "\u00BD",
	"&frac14;": "\u00BC",
	"&frac34;": "\u00BE",
}

func fix(s string) string {
	for e, c := range extraEntities {
		s = strings.ReplaceAll(s, e, c)
	}
	return html.UnescapeString(s)
}

func main() {
	root := "."
	if len(os.Args) > 1 {
		root = os.Args[1]
	}

	filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".md") {
			return err
		}
		data, _ := os.ReadFile(path)
		orig := string(data)
		if fixed := fix(orig); fixed != orig {
			fmt.Println("fixing:", path)
			os.WriteFile(path, []byte(fixed), 0644)
		}
		return nil
	})

	fmt.Println("done.")
}
