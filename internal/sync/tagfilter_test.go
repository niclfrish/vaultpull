package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTagFilter_NoRequiredTags_ReturnsAll(t *testing.T) {
	f := NewTagFilter(nil)
	secrets := map[string]string{
		"KEY_A#prod": "val1",
		"KEY_B":      "val2",
	}
	out := f.Apply(secrets)
	assert.Equal(t, secrets, out)
}

func TestTagFilter_MatchingSingleTag(t *testing.T) {
	f := NewTagFilter([]string{"prod"})
	secrets := map[string]string{
		"DB_URL#prod,staging": "postgres://prod",
		"API_KEY#staging":     "stg-key",
		"PLAIN_KEY":           "plain",
	}
	out := f.Apply(secrets)
	assert.Equal(t, map[string]string{
		"DB_URL": "postgres://prod",
	}, out)
}

func TestTagFilter_MatchesAnyTag_OrSemantics(t *testing.T) {
	f := NewTagFilter([]string{"prod", "staging"})
	secrets := map[string]string{
		"DB_URL#prod":     "postgres://prod",
		"API_KEY#staging": "stg-key",
		"SECRET#dev":      "dev-secret",
	}
	out := f.Apply(secrets)
	assert.Len(t, out, 2)
	assert.Equal(t, "postgres://prod", out["DB_URL"])
	assert.Equal(t, "stg-key", out["API_KEY"])
}

func TestTagFilter_StripsTags_FromOutputKeys(t *testing.T) {
	f := NewTagFilter([]string{"prod"})
	secrets := map[string]string{
		"MY_VAR#prod": "value",
	}
	out := f.Apply(secrets)
	_, hasRaw := out["MY_VAR#prod"]
	_, hasClean := out["MY_VAR"]
	assert.False(t, hasRaw, "raw annotated key should not appear in output")
	assert.True(t, hasClean, "clean key should appear in output")
}

func TestTagFilter_EmptyInput_ReturnsEmpty(t *testing.T) {
	f := NewTagFilter([]string{"prod"})
	out := f.Apply(map[string]string{})
	assert.Empty(t, out)
}

func TestSplitKeyTags_NoHash(t *testing.T) {
	base, tags := splitKeyTags("PLAIN_KEY")
	assert.Equal(t, "PLAIN_KEY", base)
	assert.Empty(t, tags)
}

func TestSplitKeyTags_SingleTag(t *testing.T) {
	base, tags := splitKeyTags("MY_KEY#prod")
	assert.Equal(t, "MY_KEY", base)
	assert.Equal(t, []string{"prod"}, tags)
}

func TestSplitKeyTags_MultipleTags(t *testing.T) {
	base, tags := splitKeyTags("MY_KEY#prod,staging, dev ")
	assert.Equal(t, "MY_KEY", base)
	assert.Equal(t, []string{"prod", "staging", "dev"}, tags)
}
