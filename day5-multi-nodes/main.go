package main

import (
	"flag"
	"fmt"
	"log"
	"multi-nodes/weecache"
	"net/http"
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

func startCacheServer(addr string, addrs []string, wee *weecache.Group) {
	peers := weecache.NewHTTPPool(addr)
	peers.Set(addrs...)
	wee.RegisterPeers(peers)
	log.Println("weecache is running at", addr)
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}

func startAPIServer(apiAddr string, wee *weecache.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			view, err := wee.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(view.ByteSlice())
		}))
	log.Println("fontend server running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))
}

func main() {
	var api bool
	var port int
	flag.IntVar(&port, "port", 8001, "Geecache server port")
	flag.BoolVar(&api, "api", false, "Start a api server?")
	flag.Parse()

	apiAddr := "http://localhost:9999"

	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}

	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}
	wee := createGroup()
	if api {
		go startAPIServer(apiAddr, wee)
	}
	startCacheServer(addrMap[port], []string(addrs), wee)
}
