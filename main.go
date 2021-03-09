package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/fibbery/douban_zufang/model"
	"github.com/fibbery/douban_zufang/resp"
	"github.com/gocolly/colly"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/fibbery/douban_zufang/g"
	"github.com/fibbery/douban_zufang/utils"
)

var (
	conf   *string
	numReg = regexp.MustCompile("\\d+")
)

func init() {
	conf = flag.String("conf", "conf.toml", "configuration file")
	flag.Parse()
}

func main() {
	aconf()
	pconf()
	g.InitDB()

	// 生成colly采集器
	c := colly.NewCollector(
		colly.Async(true),
	)

	c.OnRequest(func(r *colly.Request) {
		//随机user-agents
		r.Headers.Set("User-Agent", g.Config.Http.Agents[rand.Intn(len(g.Config.Http.Agents))])
	})

	//限制并发
	_ = c.Limit(&colly.LimitRule{DomainGlob: "*.douban.*", Parallelism: 5, RandomDelay: 5 * time.Second})

	//主题的访问
	c.OnHTML(".article", func(element *colly.HTMLElement) {
		url := element.Request.URL.String()
		if strings.Contains(url, "topic") {
			vistTopic(c, element)
		} else if strings.Contains(url, "doulist") {
			vistDoulist(c, element)
		}

	})

	c.OnResponse(func(response *colly.Response) {
		// 硬编码处理，访问豆列请求体
		if strings.Contains(response.Headers.Get("Content-Type"), "application/json") {
			page := &resp.Page{}
			if err := json.Unmarshal(response.Body, page); err != nil {
				log.Printf("unmarshal json error , url is %s, error is %v", response.Request.URL, err)
			} else {
				for i := 0; i < len(page.Items); i++ {
					err := c.Visit(page.Items[i].List.URL)
					if err != nil {
						log.Printf("vist url error: %+v\n", err)
					}
				}
			}
		}
	})

	startUrl := fmt.Sprintf(g.TopicUrl, g.Config.User.Topic)
	err := c.Visit(startUrl)
	if err != nil {
		log.Fatalf("visit start url : %s error: %+v", startUrl, err)
	}

	// stop
	c.Wait()
	log.Println("colly stop!!!")
}

/**
 * 访问豆列，从而获取豆列中的文章
 */
func vistDoulist(c *colly.Collector, doc *colly.HTMLElement) {
	doc.ForEach("div.bd.doulist-note", func(_ int, element *colly.HTMLElement) {
		//目标主题
		href := element.ChildAttr(".title > a", "href")
		topicId := numReg.FindString(href)
		_ = c.Visit(fmt.Sprintf(g.TopicUrl, topicId))
	})
}

/**
 * 访问文章主题，获取详细信息并存储
 */
func vistTopic(c *colly.Collector, doc *colly.HTMLElement) {
	createTime, _ := time.Parse("2006-01-02 15:04:05", doc.ChildText(".create-time"))
	topic := &model.TopicInfo{
		ID:         numReg.FindString(doc.Request.URL.Path),
		Link:       doc.Request.URL.String(),
		Title:      doc.ChildText("h1"),
		Createtime: createTime,
	}
	if createTime.After(time.Now().AddDate(0, -g.Config.User.ExpireDay, 0)) {
		log.Printf("topic [%+v] has expire, will not store\n", topic)
	} else {
		g.DB.Table("TopicInfo").Save(&topic)
		log.Printf("topic [%+v] is new, will store\n", topic)
	}

	//继续访问收藏改文章的豆列
	_ = c.Visit(fmt.Sprintf(g.TopicCollectUrl, topic.ID))
}

func pconf() {
	if err := g.Parse(*conf); err != nil {
		log.Printf("parse configuration file error, %v", err)
		os.Exit(1)
	}
}

func aconf() {
	if *conf != "" && utils.IsExsit(*conf) {
		return
	}
	log.Println("configuration file is not exist!!!")
	os.Exit(1)
}
