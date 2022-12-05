package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/gomodule/redigo/redis"
	gintrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gin-gonic/gin"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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

func getBlogsFromRedis() string {
	REDIS_ENDPOINT := os.Getenv("REDIS_ENDPOINT")

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
	tracer.Start()
	defer tracer.Stop()

	r := gin.New()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"*",
		},
		AllowMethods: []string{
			"POST",
			"GET",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Access-Control-Allow-Credentials",
			"Access-Control-Allow-Headers",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"Authorization",
		},
		AllowCredentials: true,
		MaxAge:           24 * time.Hour,
	}))

	SERVICE := os.Getenv("SERVICE")
	r.Use(gintrace.Middleware(SERVICE))

	PORT := os.Getenv("BACKENDPORT")
	SEARCH_ENDPOINT_V1 := os.Getenv("SEARCH_ENDPOINT_V1")

	r.GET(SEARCH_ENDPOINT_V1, func(c *gin.Context) {
		blogs := getBlogsFromRedis()
		c.String(200, blogs)
	})
	r.Run(":" + PORT)
}
