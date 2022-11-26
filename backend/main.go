package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/gomodule/redigo/redis"
)

type blog struct {
	Title   string `json:title`
	Url     string `json:url`
	Date    string `json:date`
	Company string `json:company`
}

type Blogs []blog

func (b Blogs) Len() int {
	return len(b)
}

func (b Blogs) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b Blogs) Less(i, j int) bool {
	return b[i].Date < b[j].Date
}

const PORT = "8080"
const SEARCH_ENDPOINT_V1 = "/techsearch/v1/blogs"
const REDIS_ENDPOINT = "localhost:6379"

func getblogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	blogs := getBlogsFromRedis()
	fmt.Fprintf(w, blogs)
}

func getBlogsFromRedis() string {
	pool := newPool(REDIS_ENDPOINT)
	conn := pool.Get()
	defer conn.Close()

	keys := getAllKeys(conn)
	b := getBlogsFromKeys(keys, conn)
	blogs, err := json.Marshal(b)
	if err != nil {
		fmt.Println(err)
	}
	return string(blogs)
}

func getAllKeys(c redis.Conn) []string {
	keys, err := redis.Strings(c.Do("KEYS", "*"))
	if err != nil {
		fmt.Println(err)
	}
	return keys
}

func getBlogsFromKeys(keys []string, c redis.Conn) []blog {
	var fields [4]string = [4]string{"title", "url", "date", "company"}
	var blogs Blogs = []blog{}
	for _, key := range keys {
		//url := strings.Split(key, ";")[1]

		blog := blog{}
		for _, field := range fields {
			value, err := redis.String(c.Do("HGET", key, field))
			if err != nil {
				fmt.Println(err)
			}
			switch field {
			case "title":
				blog.Title = value
			case "url":
				blog.Url = value
			case "date":
				blog.Date = value
			case "company":
				blog.Company = value
			}
		}
		blogs = append(blogs, blog)
	}
	sort.Sort(sort.Reverse(blogs))
	return blogs
}

func newPool(addr string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		MaxActive:   0,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", addr) },
	}
}

func main() {
	http.HandleFunc(SEARCH_ENDPOINT_V1, getblogs)
	http.ListenAndServe(":"+PORT, nil)
}
