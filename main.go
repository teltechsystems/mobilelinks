package main

import (
	"flag"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"net/http"
	"strings"
)

var (
	listenAddr = flag.String("listen", ":8000", "Host:port to expose this service on")
	publicUrl  = flag.String("public_url", "http://localhost:8000", "The public url that hashes will be appended to")
)

func main() {
	flag.Parse()

	redisPool := redis.NewPool(func() (redis.Conn, error) {
		return redis.Dial("tcp", "127.0.0.1:6379")
	}, 10)

	redisLinker := &RedisLinker{redisPool}

	http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")

		r.ParseForm()
		link, err := redisLinker.CreateLink(Link{
			Default: r.Form.Get("default"),
			Android: r.Form.Get("android"),
			IOS:     r.Form.Get("ios"),
		})

		if err != nil {
			fmt.Fprintf(w, "Failed to generate your link: %s", err)
			return
		}

		fmt.Fprintf(w, "%s/%s", *publicUrl, GetHashFromId(link.Id))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		link, err := redisLinker.FindLinkById(GetIdFromHash(r.RequestURI[1:]))
		if err != nil {
			fmt.Fprintf(w, "Failed to discover your link!")
			return
		}

		userAgent := r.Header.Get("User-Agent")

		if strings.Contains(userAgent, "Android") && link.Android != "" {
			http.Redirect(w, r, link.Android, http.StatusFound)
		} else if strings.Contains(userAgent, "iPhone") && link.IOS != "" {
			http.Redirect(w, r, link.IOS, http.StatusFound)
		} else if strings.Contains(userAgent, "iPad") && link.IOS != "" {
			http.Redirect(w, r, link.IOS, http.StatusFound)
		} else if strings.Contains(userAgent, "iPod") && link.IOS != "" {
			http.Redirect(w, r, link.IOS, http.StatusFound)
		} else {
			http.Redirect(w, r, link.Default, http.StatusFound)
		}
	})

	http.ListenAndServe(*listenAddr, nil)
}
