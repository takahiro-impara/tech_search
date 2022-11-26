package main

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/gomodule/redigo/redis"
)

type articleInfo struct {
	Title   string
	URL     string
	Date    string
	Company string
}

type articles []*articleInfo

func scrapeMercari(campany string) []articleInfo {
	MERCARI_ENDPOINT := os.Getenv("MERCARI_ENDPOINT")
	MERCARI_BASEURL := os.Getenv("MERCARI_BASEURL")

	articles := make([]articleInfo, 0)
	c := colly.NewCollector()
	c.OnHTML(".post-list__item", func(e *colly.HTMLElement) {
		article := articleInfo{}
		d := e.DOM.Find("div > time").Text()
		article.Date = strings.TrimSpace(d)
		article.Title = e.DOM.Find("div > h3").Text()
		l, _ := e.DOM.Find("a").Attr("href")
		article.URL = MERCARI_BASEURL + l
		article.Company = campany

		articles = append(articles, article)
	})

	c.Visit(MERCARI_ENDPOINT)
	writeToRedis(articles)
	return articles
}

func scrapeClassmethod(campany string) []articleInfo {
	CLASSMETHOD_ENDPOINT := os.Getenv("CLASSMETHOD_ENDPOINT")

	articles := make([]articleInfo, 0)
	c := colly.NewCollector()
	c.OnHTML(".post-container", func(e *colly.HTMLElement) {
		article := articleInfo{}
		d := e.DOM.Find("div > p").Text()
		article.Date = strings.Replace(strings.TrimSpace(d), ".", "/", -1)
		article.Title = strings.TrimSpace(e.DOM.Find("div > h3").Text())
		l, _ := e.DOM.Find("a").Attr("href")
		article.URL = CLASSMETHOD_ENDPOINT + l
		article.Company = campany

		articles = append(articles, article)
	})

	c.Visit(CLASSMETHOD_ENDPOINT)
	writeToRedis(articles)
	return articles
}

func writeToRedis(articles []articleInfo) {
	REDIS_TTL, _ := strconv.Atoi(os.Getenv("REDIS_TTL"))
	REDIS_ENDPOINT := os.Getenv("REDIS_ENDPOINT")

	pool := newPool(REDIS_ENDPOINT)
	conn := pool.Get()
	defer conn.Close()

	for _, article := range articles {
		key := article.Company + ";" + article.URL
		_, err := conn.Do("HSET", key, "title", article.Title)
		_, err = conn.Do("HSET", key, "url", article.URL)
		_, err = conn.Do("HSET", key, "date", article.Date)
		_, err = conn.Do("HSET", key, "company", article.Company)
		if err != nil {
			panic(err)
		}
		Expire(key, REDIS_TTL, conn)
	}

}

func newPool(addr string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		MaxActive:   0,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", addr) },
	}
}

func Expire(key string, ttl int, c redis.Conn) {
	c.Do("EXPIRE", key, ttl)
}

func main() {
	scrapeMercari("Mercari")
	scrapeClassmethod("Classmethod")
}
