package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	url2 "net/url"

	"github.com/fibbery/douban_zufang/g"
	"github.com/fibbery/douban_zufang/model"
	"github.com/fibbery/douban_zufang/resp"
	"github.com/fibbery/douban_zufang/utils"
	"github.com/gocolly/colly"

	"github.com/fibbery/douban_zufang/logger"

	log "github.com/sirupsen/logrus"
)

var (
	conf   *string
	numReg = regexp.MustCompile("\\d+")
)

func init() {
	conf = flag.String("conf", "conf.toml", "configuration file")
	flag.Parse()

	aconf()
	pconf()
	logger.ParseConfig(nil)
	g.InitDB()
}

func main() {

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
				log.Infof("unmarshal json error , url is %s, error is %v", response.Request.URL, err)
			} else {
				for i := 0; i < len(page.Items); i++ {
					err := c.Visit(page.Items[i].List.URL)
					if err != nil {
						log.Errorf("action: request doulist html page error, url: %s, error: %+v", page.Items[i].List.URL, err)
					}
				}
			}
		}
	})

	startUrl := fmt.Sprintf(g.DouListUrl, g.Config.User.DouList, 0)
	err := c.Visit(startUrl)
	if err != nil {
		log.Fatalf("visit start url : %s error: %+v", startUrl, err)
	}

	// stop
	c.Wait()
	log.Info("colly stop!!!")
}

/**
 * 访问豆列，从而获取豆列中的文章
 */
func vistDoulist(c *colly.Collector, doc *colly.HTMLElement) {
	size, _ := strconv.Atoi(strings.TrimFunc(doc.ChildText("div.doulist-filter > a.active > span"), func(r rune) bool {
		return r == '(' || r == ')'
	}))

	//判断当前页的offset
	current := 0
	url, _ := url2.Parse(doc.Request.URL.String())
	douListId := numReg.FindString(url.Path)
	if url.RawQuery != "" && strings.Contains(url.RawQuery, "start") {
		current, _ = strconv.Atoi(numReg.FindString(url.RawQuery))
	}

	doc.ForEach("div.bd.doulist-note", func(_ int, element *colly.HTMLElement) {
		//目标主题
		href := element.ChildAttr("div.title > a", "href")
		topicId := numReg.FindString(href)
		url := fmt.Sprintf(g.TopicUrl, topicId)
		err := c.Visit(url)
		if err != nil {
			log.Errorf("action: request topic info html page error,  url : %s, error : %+v", url, err)
		}
	})

	//判断是否需要获取下一页doulist
	if current+g.DouListPageSize < size {
		err := c.Visit(fmt.Sprintf(g.DouListUrl, douListId, current+g.DouListPageSize))
		if err != nil {
			log.Errorf("action: request doulist html page error,  url : %s, error : %+v", url, err)
		}
	}
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
	if createTime.After(time.Now().AddDate(0, 0, -g.Config.User.ExpireDay)) {
		log.Warnf("topic [%+v] has expire, will not store", topic)
	} else {
		g.DB.Table("TopicInfo").Save(&topic)
		log.Infof("topic [%+v] is new, will store", topic)
	}

	//继续访问收藏该文章的豆列
	url := fmt.Sprintf(g.TopicCollectUrl, topic.ID)
	err := c.Visit(url)
	if err != nil {
		log.Errorf("action: request topic collect doulist json error,  url : %s, error : %+v", url, err)
	}

}

func pconf() {
	if err := g.Parse(*conf); err != nil {
		log.Infof("parse configuration file error, %v", err)
		os.Exit(1)
	}
}

func aconf() {
	if *conf != "" && utils.IsExsit(*conf) {
		return
	}
	log.Info("configuration file is not exist!!!")
	os.Exit(1)
}
