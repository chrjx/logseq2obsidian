package scanner

import "testing"

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
