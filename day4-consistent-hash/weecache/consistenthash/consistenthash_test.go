package consistenthash

import (
	"strconv"
	"testing"
)

func TestHash(t *testing.T) {
	hash := New(3, func(key []byte) uint32 {
		i, _ := strconv.Atoi(string(key))
		return uint32(i)
	})

	hash.Add("6", "4", "2")
	TestCase := map[string]string{
		"2":  "2",
		"11": "2",
		"23": "4",
		"27": "2",
	}

	for k, v := range TestCase {
		if hash.Get(k) != v {
			t.Errorf("Asking %s, should have yielded %s", k, v)
		}
	}

	hash.Add("8")
	TestCase["27"] = "8"

	for k, v := range TestCase {
		if hash.Get(k) != v {
			t.Errorf("Asking %s, should have yielded %s", k, v)
		}
	}
}
