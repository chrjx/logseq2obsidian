package scanner

type property struct {
	name  string
	value string
}

type block struct {
	level    int
	parent   []block
	children []block
}

type Page struct {
	properties []*property
	blocks     []*block
	parent     *Page
	children   []*Page
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
	return false
}

// ReadPage
// ReadPage By Line
func ReadPage(lines []string) *Page {
	// p := Page{}
	// var curr *block
	// curr = nil
	// stack

	// for _, l := range lines {
	// if isProperty(l) {
	// 	// if on blocks read, then the properties are block

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
