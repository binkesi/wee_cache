package main

import (
	"fmt"
	"http-server/weecache"
	"log"
	"net/http"
)

var db = map[string]string{
	"sungn":    "0528",
	"wumx":     "0722",
	"together": "1018",
}

func main() {
	weecache.NewGroup("scores", 2<<10, weecache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
	addr := "localhost:9999"
	peers := weecache.NewHTTPPool(addr)
	log.Println("geecache is running at", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}
