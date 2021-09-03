package main

import (
	"fmt"
	"log"
	"multi-nodes/weecache"
)

var db = map[string]string{
	"sungn": "0528",
	"wumx":  "0722",
	"toge":  "1018",
}

func createGroup() *weecache.Group {
	return weecache.NewGroup("scores", 2<<10, weecache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
}
