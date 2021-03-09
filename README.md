# douban_zufang
douban_zufang is a simple Go scrapy aimed at scraping douban zufang info, built on colly.
General procedure steps:
1. give a topic id as the initial url, base your interest
2. scraping topic info, store it base on it's create time
3. scraping dou_list what has collected this topic
4. scraping topic info where is collected by dou_list at step 3
4. cycle step 2

## How-To-Use
1. clone this repository
```shell script
git clone git@github.com:fibbery/douban_zufang.git
```
2. init the database
in the work directory, use the mysql client connect to the mysql server, execute sql like below:
```mysql
source db.sql
```
this operate will create database, table, and user.
3. start program
```shell script
 go run main.go 
```
or 
```shell script
go run main.go -conf conf.toml
```

## About Conf
The configuration file is like below:
```toml
[user]
topic = 213336517  # your interest topic as initial url 
expire_day = 10  # expire day whether to store

[http]
agents = [
    "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.95 Safari/537.36",
    "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:34.0) Gecko/20100101 Firefox/34.0",
    "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/534.57.2 (KHTML, like Gecko) Version/5.1.7 Safari/534.57.2",
    "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36",
    "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.64 Safari/537.11",
    "Mozilla/5.0 (Windows; U; Windows NT 6.1; en-US) AppleWebKit/534.16 (KHTML, like Gecko) Chrome/10.0.648.133 Safari/534.16",
    "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/534.57.2 (KHTML, like Gecko) Version/5.1.7 Safari/534.57.2",
] # brower user-agents
interval = 1000  # http request interval: ms

[db]
dsn = "root@tcp(127.0.0.1:3306)/douban_zufang?charset=utf8&parseTime=true" 
```
