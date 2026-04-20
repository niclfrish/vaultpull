package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilter_NoCriteria_ReturnsAll(t *testing.T) {
	f := NewFilter(nil, nil)
	input := map[string]string{
		"DB_HOST": "localhost",
		"API_KEY": "secret",
		"PORT":    "5432",
	}
	out := f.Apply(input)
	assert.Equal(t, input, out)
}

func TestFilter_PrefixFilter_KeepsMatching(t *testing.T) {
	f := NewFilter([]string{"DB_"}, nil)
	input := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"API_KEY": "secret",
	}
	out := f.Apply(input)
	assert.Equal(t, map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	}, out)
}

func TestFilter_ExcludeFilter_RemovesKeys(t *testing.T) {
	f := NewFilter(nil, []string{"API_KEY"})
	input := map[string]string{
		"DB_HOST": "localhost",
		"API_KEY": "secret",
	}
	out := f.Apply(input)
	assert.Equal(t, map[string]string{
		"DB_HOST": "localhost",
	}, out)
}

func TestFilter_PrefixAndExclude_Combined(t *testing.T) {
	f := NewFilter([]string{"DB_"}, []string{"DB_PASSWORD"})
	input := map[string]string{
		"DB_HOST":     "localhost",
		"DB_PASSWORD": "s3cr3t",
		"API_KEY":     "key",
	}
	out := f.Apply(input)
	assert.Equal(t, map[string]string{
		"DB_HOST": "localhost",
	}, out)
}

func TestFilter_EmptyInput_ReturnsEmpty(t *testing.T) {
	f := NewFilter([]string{"DB_"}, nil)
	out := f.Apply(map[string]string{})
	assert.Empty(t, out)
}

func TestFilter_CaseInsensitivePrefixMatch(t *testing.T) {
	f := NewFilter([]string{"db_"}, nil)
	input := map[string]string{
		"DB_HOST": "localhost",
		"API_KEY": "key",
	}
	out := f.Apply(input)
	assert.Equal(t, map[string]string{
		"DB_HOST": "localhost",
	}, out)
}
