package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"time"
)

const (
	contentDir  = "content"
	indexFile   = contentDir + "/_index.md"
	archivesDir = contentDir + "/archives"
	archivesFile = archivesDir + "/_index.md"
)

var cutoffYear = time.Now().Year() - 1

type post struct {
	title string
	url   string
	date  time.Time
}

type monthKey struct{ year, month int }

func escapeMarkdown(s string) string {
	s = strings.ReplaceAll(s, "[", "\\[")
	return strings.ReplaceAll(s, "]", "\\]")
}

// parseFrontmatter extrai title e date do bloco --- ... --- sem dependência externa.
func parseFrontmatter(content string) (title, dateStr string, ok bool) {
	if !strings.HasPrefix(content, "---\n") {
		return
	}
	end := strings.Index(content[4:], "\n---\n")
	if end == -1 {
		return
	}
	block := content[4 : end+4]
	scanner := bufio.NewScanner(strings.NewReader(block))
	for scanner.Scan() {
		line := scanner.Text()
		k, v, found := strings.Cut(line, ":")
		if !found {
			continue
		}
		v = strings.TrimSpace(v)
		switch strings.TrimSpace(k) {
		case "title":
			title = strings.Trim(v, `"'`)
		case "date":
			dateStr = strings.Trim(v, `"'`)
		}
	}
	ok = title != "" && dateStr != ""
	return
}

var dateFormats = []string{
	time.RFC3339,
	"2006-01-02T15:04:05",
	"2006-01-02",
}

func parseDate(s string) (time.Time, error) {
	for _, f := range dateFormats {
		if t, err := time.Parse(f, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("formato de data não reconhecido: %s", s)
}

func parsePost(path string) (*post, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	title, dateStr, ok := parseFrontmatter(string(data))
	if !ok {
		return nil, nil
	}
	date, err := parseDate(dateStr)
	if err != nil {
		return nil, err
	}
	url := strings.TrimPrefix(path, contentDir+"/")
	url = strings.TrimSuffix(url, "/index.md") + "/"
	return &post{title: title, url: url, date: date}, nil
}

func collectPosts(includeFuture bool) []post {
	now := time.Now()
	var posts []post
	filepath.WalkDir(contentDir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		if !strings.HasSuffix(path, "/index.md") {
			return nil
		}
		if path == contentDir+"/index.md" || path == contentDir+"/_index.md" {
			return nil
		}
		if strings.HasPrefix(path, archivesDir+"/") {
			return nil
		}
		p, err := parsePost(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "erro em %s: %v\n", path, err)
			return nil
		}
		if p != nil && (includeFuture || !p.date.After(now)) {
			posts = append(posts, *p)
		}
		return nil
	})
	return posts
}

func groupByMonth(posts []post) map[monthKey][]post {
	sort.Slice(posts, func(i, j int) bool { return posts[i].date.After(posts[j].date) })
	groups := map[monthKey][]post{}
	for _, p := range posts {
		k := monthKey{p.date.Year(), int(p.date.Month())}
		groups[k] = append(groups[k], p)
	}
	return groups
}

func sortedKeys(groups map[monthKey][]post) []monthKey {
	keys := make([]monthKey, 0, len(groups))
	for k := range groups {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].year != keys[j].year {
			return keys[i].year > keys[j].year
		}
		return keys[i].month > keys[j].month
	})
	return keys
}

func renderMonths(groups map[monthKey][]post, urlPrefix string) []string {
	var lines []string
	for _, k := range sortedKeys(groups) {
		month := time.Month(k.month)
		lines = append(lines, fmt.Sprintf("## %d - %s\n", k.year, month))
		for _, p := range groups[k] {
			lines = append(lines, fmt.Sprintf("- [%s](%s%s)", escapeMarkdown(p.title), urlPrefix, p.url))
		}
		lines = append(lines, "")
	}
	return lines
}

func writeFile(path, content string) {
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		fmt.Fprintln(os.Stderr, "erro ao gravar:", err)
		os.Exit(1)
	}
}

func countPosts(groups map[monthKey][]post) int {
	n := 0
	for _, v := range groups {
		n += len(v)
	}
	return n
}

func main() {
	includeFuture := slices.Contains(os.Args[1:], "--future")
	posts := collectPosts(includeFuture)
	groups := groupByMonth(posts)

	recent := map[monthKey][]post{}
	archived := map[monthKey][]post{}
	for k, v := range groups {
		if k.year >= cutoffYear {
			recent[k] = v
		} else {
			archived[k] = v
		}
	}

	os.MkdirAll(archivesDir, 0755)

	futureNote := ""
	if includeFuture {
		futureNote = " (including future posts)"
	}

	idx := strings.Join(append(
		[]string{"---\ntitle: Palavra Indexada\n---\n"},
		append(renderMonths(recent, ""), "[Arquivo completo →](/archives/)\n")...,
	), "\n")
	writeFile(indexFile, idx)
	fmt.Printf("Generated %s with %d posts (%d+).%s\n", indexFile, countPosts(recent), cutoffYear, futureNote)

	arch := strings.Join(append(
		[]string{"---\ntitle: Palavra Indexada - Arquivo\n---\n"},
		renderMonths(archived, "/")...,
	), "\n")
	writeFile(archivesFile, arch)
	fmt.Printf("Generated %s with %d posts (before %d).%s\n", archivesFile, countPosts(archived), cutoffYear, futureNote)
}
