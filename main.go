// QiuBaiSpider project main.go
package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
)

//网页请求函数
func GetHtml() string {
	client := &http.Client{}
	url := "https://www.qiushibaike.com/hot/"
	request, err := http.NewRequest("GET", url, nil)
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36")
	if err != nil {
		fmt.Println(err)
	}
	response, err := client.Do(request)
	if response.StatusCode != 200 {
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
	file, err := os.Create("qiubai.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file.WriteString("\xEF\xBB\xBF")
	write := csv.NewWriter(file)
	for _, value := range results {
		write.Write([]string{value[1], value[2], value[3], value[4]})
		write.Flush()
		fmt.Println("作者：")
		fmt.Println(value[1])
		fmt.Println("段子内容：")
		fmt.Println(value[2])
		fmt.Println("点赞数量：")
		fmt.Println(value[3])
		fmt.Println("评论数量：")
		fmt.Println(value[4])
		fmt.Println("-------------------------------------")
	}
}

func main() {
	html := GetHtml()
	GetContent(html)
}
