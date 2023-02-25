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
