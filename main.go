// QiuBaiSpider project main.go
package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

//网页请求函数
func GetHtml(pagenum int) string {
	client := &http.Client{}
	url := "https://www.qiushibaike.com/hot/page/" + strconv.Itoa(pagenum) + "/"
	request, err := http.NewRequest("GET", url, nil)
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36")
	if err != nil {
		fmt.Println(err)
	}
	response, err := client.Do(request)
	if response.StatusCode != 200 {
		response.Body.Close()
		fmt.Println("QiuBaiSpider request failed!", err)
	}
	//fmt.Println(reflect.TypeOf(response))
	html, _ := ioutil.ReadAll(response.Body)
	return string(html)
}

//正则解析
func GetContent(html string) {
	pattern := regexp.MustCompile(`<div class="article block untagged mb15[\s\S]*?">[\s\S]*?<div class="author clearfix">[\s\S]*?<h2>([\s\S]*?)</h2>[\s\S]*?</a>[\s\S]*?<div class="content">[\s\S]*?<span>([\s\S]*?)</span>[\s\S]*?</div>[\s\S]*?</a>[\s\S]*?<div class="stats">[\s\S]*?<i class="number">([\s\S]*?)</i>[\s\S]*?<span class="stats-comments">[\s\S]*?<a [\s\S]*?>[\s\S]*?<i class="number">([\s\S]*?)</i>[\s\S]*?</div>`)
	results := pattern.FindAllStringSubmatch(html, -1)
	//写入csv
	for _, value := range results {
		SaveCsv(value)
		fmt.Println("作者：")
		fmt.Println(strings.Replace(value[1], "\n", "", -1))
		fmt.Println("段子内容：")
		fmt.Println(strings.Replace(value[2], "\n", "", -1))
		fmt.Println("点赞数量：")
		fmt.Println(value[3])
		fmt.Println("评论数量：")
		fmt.Println(value[4])
		fmt.Println("-------------------------------------")
	}
}

//存入CSV
func SaveCsv(value []string) {
	_, err := os.Stat("qiubai.csv") //如果文件存在返回nil
	// 如果存在则打开csv文件，不存在则创建一个csv文件
	if err == nil {
		file, err := os.OpenFile("qiubai.csv", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0640)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		write := csv.NewWriter(file)
		write.Write([]string{strings.Replace(value[1], "\n", "", -1), strings.Replace(value[2], "\n", "", -1), value[3], value[4]})
		write.Flush()
	} else {
		file, err := os.Create("qiubai.csv")
		if err != nil {
			panic(err)
		}
		defer file.Close()
		//写入utf-8 BOM防止中文乱码
		file.WriteString("\xEF\xBB\xBF")
		write := csv.NewWriter(file)
		write.Write([]string{strings.Replace(value[1], "\n", "", -1), strings.Replace(value[2], "\n", "", -1), value[3], value[4]})
		write.Flush()
	}

}

//抓取13个页面
func Run() {
	for i := 1; i <= 13; i++ {
		pagenum := i
		html := GetHtml(pagenum)
		GetContent(html)
	}
}

func main() {
	Run()
}
