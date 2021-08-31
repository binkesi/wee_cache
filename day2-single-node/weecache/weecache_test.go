package weecache

import (
	"reflect"
	"testing"
)

func TestGetter(t *testing.T) {
	var g Getter = GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})
	expect := []byte("key")
	if v, _ := g.Get("key"); !reflect.DeepEqual(v, expect) {
		t.Errorf("callback failed!")
	}
}
