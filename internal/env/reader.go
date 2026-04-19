package env

import (
	"bufio"
	"os"
	"strings"
)

// Reader reads an existing .env file into a map.
type Reader struct {
	path string
}

// NewReader creates a Reader for the given file path.
func NewReader(path string) *Reader {
	return &Reader{path: path}
}

// Read parses the .env file and returns a map of key-value pairs.
// Lines starting with '#' and empty lines are ignored.
func (r *Reader) Read() (map[string]string, error) {
	f, err := os.Open(r.path)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]string{}, nil
		}
		return nil, err
	}
	defer f.Close()

	result := make(map[string]string)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.Trim(strings.TrimSpace(parts[1]), "\"")
		result[key] = val
	}
	return result, scanner.Err()
}
