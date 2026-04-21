package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDedupe_NoDuplicates(t *testing.T) {
	pairs := []string{"FOO=bar", "BAZ=qux"}
	res := Dedupe(pairs, DedupeKeepFirst)
	assert.Equal(t, map[string]string{"FOO": "bar", "BAZ": "qux"}, res.Secrets)
	assert.Empty(t, res.Duplicates)
}

func TestDedupe_KeepFirst(t *testing.T) {
	pairs := []string{"FOO=first", "BAR=baz", "FOO=second"}
	res := Dedupe(pairs, DedupeKeepFirst)
	assert.Equal(t, "first", res.Secrets["FOO"])
	assert.Contains(t, res.Duplicates, "FOO")
}

func TestDedupe_KeepLast(t *testing.T) {
	pairs := []string{"FOO=first", "BAR=baz", "FOO=second"}
	res := Dedupe(pairs, DedupeKeepLast)
	assert.Equal(t, "second", res.Secrets["FOO"])
	assert.Contains(t, res.Duplicates, "FOO")
}

func TestDedupe_SkipsInvalidPairs(t *testing.T) {
	pairs := []string{"NODEQUALS", "KEY=val"}
	res := Dedupe(pairs, DedupeKeepFirst)
	assert.Equal(t, map[string]string{"KEY": "val"}, res.Secrets)
	assert.Empty(t, res.Duplicates)
}

func TestDedupe_EmptyInput(t *testing.T) {
	res := Dedupe([]string{}, DedupeKeepFirst)
	assert.Empty(t, res.Secrets)
	assert.Empty(t, res.Duplicates)
}

func TestDedupe_MultipleDuplicates(t *testing.T) {
	pairs := []string{"A=1", "A=2", "A=3", "B=x", "B=y"}
	res := Dedupe(pairs, DedupeKeepFirst)
	assert.Equal(t, "1", res.Secrets["A"])
	assert.Equal(t, "x", res.Secrets["B"])
	// Two extra occurrences of A and one extra of B.
	count := 0
	for _, d := range res.Duplicates {
		if d == "A" {
			count++
		}
	}
	assert.Equal(t, 2, count)
}

func TestDedupe_ValueWithEqualsSign(t *testing.T) {
	pairs := []string{"TOKEN=abc=def=ghi"}
	res := Dedupe(pairs, DedupeKeepFirst)
	assert.Equal(t, "abc=def=ghi", res.Secrets["TOKEN"])
}
