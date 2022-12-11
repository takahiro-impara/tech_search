package main

import (
	"log"
	"os"
	"regexp"
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

func validCheck(a *articleInfo) bool {
	isValid := true
	urlFormat := regexp.MustCompile(`https*`)
	dateFormat := regexp.MustCompile(`^\d{4}\/\d{2}\/\d{2}$`)

	if len(a.Title) == 0 {
		isValid = false
		log.Printf("[WARNING] %s's Title is not valid: %s", a.Company, a.Title)
	} else if urlFormat.MatchString(a.URL) == false {
		isValid = false
		log.Printf("[WARNING] %s's URL is not valid: %s", a.Company, a.URL)
	} else if dateFormat.MatchString(a.Date) == false {
		isValid = false
		log.Printf("[WARNING] %s's Date is not valid: %s", a.Company, a.Date)
	}
	return isValid
}

func scrapeMercari(campany string) []articleInfo {
	log.Printf("start to scrape %s", campany)
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
	log.Printf("complete %s process", campany)
	return articles
}

func scrapeClassmethod(campany string) []articleInfo {
	log.Printf("start to scrape %s", campany)
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
	log.Printf("complete %s process", campany)
	return articles
}

func scrapeZozo(campany string) []articleInfo {
	log.Printf("start to scrape %s", campany)
	ZOZO_ENDPOINT := os.Getenv("ZOZO_ENDPOINT")

	articles := make([]articleInfo, 0)
	c := colly.NewCollector()
	c.OnHTML(".archive-entry-header", func(e *colly.HTMLElement) {
		article := articleInfo{}
		d, _ := e.DOM.Find("a > time").Attr("datetime")
		article.Date = strings.Replace(strings.TrimSpace(d), "-", "/", -1)
		article.Title = strings.TrimSpace(e.DOM.Find("div > h1").Text())
		article.URL, _ = e.DOM.Find("h1 > a").Attr("href")
		article.Company = campany

		articles = append(articles, article)
	})

	c.Visit(ZOZO_ENDPOINT)
	writeToRedis(articles)
	log.Printf("complete %s process", campany)
	return articles
}

func scrapeDeNA(campany string) []articleInfo {
	log.Printf("start to scrape %s", campany)
	DeNA_ENDPOINT := os.Getenv("DeNA_ENDPOINT")
	DeNA_BASEURL := os.Getenv("DeNA_BASEURL")

	articles := make([]articleInfo, 0)
	c := colly.NewCollector()
	c.OnHTML(".justify-items-start", func(e *colly.HTMLElement) {
		article := articleInfo{}
		reg := "\r\n|\n"
		d := regexp.MustCompile(reg).Split(strings.TrimSpace(e.DOM.Find("p > span").Text()), -1)[0]
		article.Date = strings.Replace(d, ".", "/", -1)
		article.Title = strings.TrimSpace(e.DOM.Find("section > a").Text())
		l, _ := e.DOM.Find("section > a").Attr("href")
		article.URL = DeNA_BASEURL + l
		article.Company = campany

		articles = append(articles, article)
	})

	c.Visit(DeNA_ENDPOINT)
	writeToRedis(articles)
	log.Printf("complete %s process", campany)
	return articles
}

func writeToRedis(articles []articleInfo) {
	REDIS_TTL, _ := strconv.Atoi(os.Getenv("REDIS_TTL"))
	REDIS_ENDPOINT := os.Getenv("REDIS_ENDPOINT")

	pool := newPool(REDIS_ENDPOINT)
	conn := pool.Get()
	log.Printf("connect to redis: endpoint: %s", REDIS_ENDPOINT)

	defer conn.Close()

	for _, article := range articles {
		if validCheck(&article) == false {
			continue
		}
		key := article.Company + ";" + article.URL
		_, err := conn.Do("HSET", key, "title", article.Title)
		_, err = conn.Do("HSET", key, "url", article.URL)
		_, err = conn.Do("HSET", key, "date", article.Date)
		_, err = conn.Do("HSET", key, "company", article.Company)
		if err != nil {
			log.Panic(err)
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
	_, err := c.Do("EXPIRE", key, ttl)
	if err != nil {
		log.Panic(err)
	}
}

func main() {
	SERVICE := os.Getenv("SERVICE")
	log.SetPrefix(SERVICE + ": ")

	log.Println("start batch process")
	scrapeMercari("Mercari")
	scrapeClassmethod("Classmethod")
	scrapeZozo("ZOZO")
	scrapeDeNA("DeNA")
	log.Println("[OK] completed batch process ")
}
