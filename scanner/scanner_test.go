package scanner

import (
	"fmt"
	"github.com/go-test/deep"
	"os"
	"reflect"
	"testing"
)

const (
	aLogseqPageStr = "title:: <section>\n\n- Represents a section of a document. Each `section` has a heading tag (`h1`-`h6`), then the section body.\n- Example:\n\t- ```\n\t  <section>\n\t      <h2>A section of the page</h2>\n\t      <p>...</p>\n\t      <img ...>\n\t  </section>\n\t  ```\n\t- LaTeX Equation\n\t\t- $$\n\t\t  C = e_k (M) = M^e \\pmod{n}\n\t\t  $$\n- It's useful to break a long article into different sections.\n- Shouldn't be used as a generic container element. [[<div>]] is made for this."
)

var (
	title  = property{name: "title", value: "<section>"}
	block1 = block{
		level:   0,
		content: "Represents a section of a document. Each `section` has a heading tag (`h1`-`h6`), then the section body.",
	}
	block2 = block{
		level:   0,
		content: "Example:",
	}
	block3 = block{
		level:   1,
		content: "```\n<section>\n    <h2>A section of the page</h2>\n    <p>...</p>\n    <img ...>\n</section>\n```",
	}
	block4 = block{
		level:   1,
		content: "LaTeX Equation",
	}
	block5 = block{
		level:   2,
		content: "$$\nC = e_k (M) = M^e \\pmod{n}\n$$",
	}
	block6 = block{
		level:   0,
		content: "It's useful to break a long article into different sections.",
	}
	block7 = block{
		level:   0,
		content: "Shouldn't be used as a generic container element. [[<div>]] is made for this.",
	}

	aLogseqPage = Page{
		title: "test.md",
		props: []property{title},
		blocks: []*block{
			&block1,
			&block2,
			&block3,
			&block4,
			&block5,
			&block6,
			&block7,
		},
	}
)

func TestAppendContent(t *testing.T) {
	b := block{content: "foo"}
	b.appendContent("bar")
	if b.content != "foobar" {
		t.Errorf("append error: %s", b.content)
	}
}

func TestIsBlock(t *testing.T) {
	t.Run("standard block", func(t *testing.T) {
		if b, l := isBlock("- foobar"); !b || l != 0 {
			t.Errorf("pasring error: {boolean: %t, level: %d}", b, l)
		}
	})
	t.Run("one-level block", func(t *testing.T) {
		if b, l := isBlock("\t- foobar"); !b || l != 1 {
			t.Errorf("pasring error: {boolean: %t, level: %d}", b, l)
		}
	})
	t.Run("two-level block", func(t *testing.T) {
		if b, l := isBlock("\t\t- foo-bar"); !b || l != 2 {
			t.Errorf("pasring error: {boolean: %t, level: %d}", b, l)
		}
	})
	t.Run("random line", func(t *testing.T) {
		if b, l := isBlock("231rdfda34242d12334r"); b && l != -1 {
			t.Errorf("pasring error: {boolean: %t, level: %d}", b, l)
		}
	})
	t.Run("strings following the dash without a space", func(t *testing.T) {
		if b, l := isBlock("-foobar"); b && l != -1 {
			t.Errorf("pasring error: {boolean: %t, level: %d}", b, l)
		}
	})
	t.Run("wrong format", func(t *testing.T) {
		if b, l := isBlock("\t-foobar"); b && l != -1 {
			t.Errorf("pasring error: {boolean: %t, level: %d}", b, l)
		}
	})
	t.Run("dash in the line", func(t *testing.T) {
		if b, l := isBlock("- foo-bar"); !b || l != 0 {
			t.Errorf("pasring error: {boolean: %t, level: %d}", b, l)
		}
	})
}

func TestIsProperty(t *testing.T) {
	t.Run("standard", func(t *testing.T) {
		if !isProperty("key:: val") {
			t.Errorf("parsing error")
		}
	})
	t.Run("after white space", func(t *testing.T) {
		if !isProperty("\t key:: val") {
			t.Errorf("parsing error")
		}
	})
	t.Run("single column", func(t *testing.T) {
		if isProperty("key: val") {
			t.Errorf("parsing error")
		}
	})
	t.Run("columns with no space", func(t *testing.T) {
		if isProperty("key::val") {
			t.Errorf("parsing error")
		}
	})
	t.Run("\":: \" in the middle", func(t *testing.T) {
		if isProperty("asd key:: val") {
			t.Errorf("parsing error")
		}
	})
}

func TestGetPageTitle(t *testing.T) {
	t.Run("with /", func(t *testing.T) {
		get := "Section1%2FSection2%2FSection3%2F"
		want := "Section1/Section2/Section3/"
		if got := GetPageTitle(get); got != want {
			t.Errorf("Get: %s, Want: %s, Got:%s", get, want, got)
		}
	})
	t.Run("some/dir/URL%2FEncoding", func(t *testing.T) {
		get := "some/dir/URL%2FEncoding"
		want := "URL/Encoding"
		if got := GetPageTitle(get); got != want {
			t.Errorf("Get: %s, Want: %s, Got:%s", get, want, got)
		}
	})
}

func TestParseBlock(t *testing.T) {
	t.Run("standard block", func(t *testing.T) {
		if b := parseBlock("- foobar"); b.content != "foobar" {
			t.Errorf("pasring error, got: %s", b.content)
		}
	})
	t.Run("one-level block", func(t *testing.T) {
		if b := parseBlock("\t- foobar"); b.content != "foobar" {
			t.Errorf("pasring error, got: %s", b.content)
		}
	})
	t.Run("two-level block", func(t *testing.T) {
		if b := parseBlock("\t\t- foo-bar"); b.content != "foo-bar" {
			t.Errorf("pasring error, got: %s", b.content)
		}
	})
}

func TestParseProperty(t *testing.T) {
	t.Run("standard", func(t *testing.T) {
		prop := "\t\t\trefs:: Les Miserable, The Great Gatsby, Illusion Perdues"
		get := parseProperty(prop)
		want := &property{
			name:  "refs",
			value: "Les Miserable, The Great Gatsby, Illusion Perdues",
		}
		if reflect.DeepEqual(get, want) {
			t.Errorf("Unequal")
		}
	})
}

func TestParsePage(t *testing.T) {
	// Write to file
	path := "test.md"
	f, _ := os.Create(path)
	_, err := f.WriteString(aLogseqPageStr)
	defer f.Close()
	if err != nil {
		fmt.Println(err)
	}

	got := ParsePage(path)
	deep.CompareUnexportedFields = true
	if diff := deep.Equal(got, &aLogseqPage); diff != nil {
		t.Error(diff)
	}
	os.Remove(path)
}

func TestPage_WriteInObsidian(t *testing.T) {
	path := "foo%2Fbar%2Ftest.md"
	f, _ := os.Create(path)
	_, err := f.WriteString(aLogseqPageStr)
	defer f.Close()
	if err != nil {
		t.Error(err)
	}

	p := ParsePage(path)
	err = p.WriteInObsidian(".")
	if err != nil {
		t.Error(err)
	}
	// check if file in the right path foo/bar/test.md
	content, err := os.ReadFile("foo/bar/test.md")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(content))
	os.Remove(path)
	os.RemoveAll("foo/")
}
