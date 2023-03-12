package scanner

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

type property struct {
	name  string
	value string
}

type block struct {
	level   int
	content string
	props   []property
	//parent   *block
	//children []block
}

func (b *block) appendContent(val string) {
	b.content += val
}

type Page struct {
	title  string
	props  []property
	blocks []*block
	// parent     *Page
	// children   []*Page
}

func isBlock(line string) (r bool, level int) {
	level = 0
	for i, ch := range line {
		if ch == '\t' {
			level += 1
		}
		if ch == '-' && i+1 < len(line) && line[i+1] == ' ' {
			return true, level
		}
	}
	return false, -1
}

func isProperty(line string) bool {
	propName := false
	for i, ch := range line {
		if ch == ':' && i+2 < len(line) && line[i:i+3] == ":: " && propName {
			return true
		}
		if unicode.IsSpace(ch) {
			if propName {
				return false
			}
		} else {
			propName = true
			continue
		}
	}
	// matched, _ := regexp.Match("(:: ).", []byte(line))
	return false
}

func parseBlock(line string) block {
	s := strings.TrimSpace(line)
	return block{content: s[2:]}
}

func parseProperty(line string) *property {
	if !isProperty(line) {
		return nil
	}
	prop := strings.TrimSpace(line)
	var name, vals string
	for i, ch := range prop {
		if ch == ':' && i+2 < len(line) && line[i:i+3] == ":: " {
			name = prop[:i]
			vals = prop[i+3:]
			break
		}
	}
	p := &property{
		name:  name,
		value: vals,
	}
	return p
}

func GetPageTitle(title string) string {
	fileName := filepath.Base(title)
	decoded, _ := url.QueryUnescape(fileName)
	return decoded
}

func ParsePage(path string) *Page {
	lines, _ := readFileToLines(path)
	p := Page{title: GetPageTitle(path)}
	var curr *block = nil
	for _, l := range lines {
		if l == "" {
			continue
		}
		if isProperty(l) {
			prop := parseProperty(l)
			if curr == nil {
				// if no blocks have been read, the props area page properties
				p.props = append(p.props, *prop)
				continue
			} else {
				// if cursor holding a block, append the property to the block
				curr.props = append(p.props, *prop)
				continue
			}
		}
		// the line is a new block head
		if b, level := isBlock(l); b {
			b := parseBlock(l)
			b.level = level
			p.blocks = append(p.blocks, &b)
			curr = &b
			continue
		} else {
			// the line is neither a property nor a start of page ""append the line to current block
			// trim tabs
			l = strings.TrimLeft(l, strings.Repeat("\t", curr.level))
			l = l[2:]
			curr.appendContent("\n" + l)
		}
	}
	return &p
}

func readFileToLines(path string) ([]string, error) {
	lines := make([]string, 0)
	readFile, err := os.Open(path)
	defer readFile.Close()

	if err != nil {
		return nil, err
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		lines = append(lines, fileScanner.Text())
	}

	return lines, nil
}

func (p *Page) WriteInObsidian(obsidianPath string) error {
	// check last position if '/'
	if obsidianPath[len(obsidianPath)-1] != '/' {
		obsidianPath += "/"
	}
	path := obsidianPath + p.title
	create := func(path string) (*os.File, error) {
		if err := os.MkdirAll(filepath.Dir(path), 0770); err != nil {
			return nil, err
		}
		return os.Create(path)
	}
	f, err := create(path)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}

	for _, b := range p.blocks {
		if _, err := f.WriteString(b.content); err != nil {
			fmt.Errorf(err.Error())
		}
		f.WriteString("\n\n")
	}
	return nil
}
