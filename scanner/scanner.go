package scanner

import (
	"bufio"
	"net/url"
	"os"
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
	//parent   *block
	//children []block
}

func (b *block) appendContent(val string) {
	b.content += val
}

type Page struct {
	title      string
	properties []property
	blocks     []*block
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

func GetPageTitle(title string) string {
	fileName := filepath.Base(title)
	decoded, _ := url.QueryUnescape(fileName)
	return decoded
}

	// 	// if cursor holding a block, append the property to the block
	// }
	//  if boolean, level := isBlock(l); boolean {
	// 	// level > curr.level, a child block of the current block

	// 	// level = curr level, append a new block with same level

	// 	// level < curr.level, exit current level,
	// } else {

	// }

	// }
	return nil
}
