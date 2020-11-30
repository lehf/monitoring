package main

import (
	"fmt"
	"os"
	"strconv"
	"net/http"
	"path/filepath"
	"io/ioutil"
  	"github.com/PuerkitoBio/goquery"
	"time"
	. "github.com/lxn/walk/declarative"
)

var time1 int



func main() {

	path, _ := os.Executable()
	dir := filepath.Dir(path)

	time_bytes, time_text_err := ioutil.ReadFile(dir+"\\time.txt")

	if time_text_err == nil{
	    time1, _ = strconv.Atoi(string(time_bytes))
	}else{
	 	time1 = 10
	}


	var ch chan int
	//定时任务
	ticker := time.NewTicker(time.Second * time.Duration(time1))
	go func() {
			for range ticker.C {
				opFile()
			}
			ch <- 1
	}()
	<-ch

}	
func opFile(){

	// 获取执行文件当前文件夹路径
	path, _ := os.Executable()

	dir := filepath.Dir(path)

	older1_path := dir+"\\older1.txt"
	older2_path := dir+"\\older2.txt"
	new_path := dir+"\\new.txt"


	_, older1_err := os.Stat(older1_path)
	if older1_err != nil {
			fmt.Println(older1_err)
	}

	if older1_err == nil {
		os.Remove(older1_path)
	}

	_, older2_err := os.Stat(older2_path)
	if older2_err != nil {
		fmt.Println(older2_err)
}

	if older2_err == nil {
		os.Rename(older2_path, older1_path)
	}

	_, err := os.Stat(new_path)
	if err != nil {
		fmt.Println(err)
}
	if err == nil{
		re_err := os.Rename(new_path, older2_path)
		if re_err !=nil{
			fmt.Println(re_err)
		}
	}else{
		fmt.Println(err)

	}	
	f, _ := os.Create(dir+"\\new.txt")
	f.Close()
	spider(new_path)

	if compare(new_path, older2_path) == true {
		fmt.Println("数据有更新！请查看")
		message()
	}else{
		fmt.Println("数据暂无更新")
	}
}

// 对比是否有更新
func compare(spath, dpath string) bool {

	text1, err1 := ioutil.ReadFile(spath)
	text2, err2 := ioutil.ReadFile(dpath)

    if err1 != nil {
        fmt.Println("error : %s", err1)
	}
	if err2 != nil {
        fmt.Println("error : %s", err2)
    }
	
	
	if string(text1) == string(text2) {
		return false
	}else{
		return true
	}

}

// 爬取
func spider(path string) {
	var url string
	current_path, _ := os.Executable()
	dir := filepath.Dir(current_path)
	url_bytes, url_err := ioutil.ReadFile(dir+"\\url.txt")

	if url_err == nil{
	    url = string(url_bytes)
	}else{
		url = ""
	}
	if url== ""{
		fmt.Println("请输入要监测的网址！！！！")
	}

	result, get_err := HttpGet(url)
	if get_err != nil {
		fmt.Printf("抓取错误：%s\n", get_err )
	}else{
		fmt.Printf("%s抓取\n",time.Now()   )

	}
	f, _ := os.OpenFile(path, os.O_RDWR, 6)
	f.WriteString(result)
	f.Close()
}

func HttpGet(url string) (result string, err error) {
	resp, err1 := http.Get(url) //发送get请求

	if err1 != nil {
			err = err1
			return
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)

	var rule string
	path, _ := os.Executable()
	dir := filepath.Dir(path)
	rule_bytes, rule_text_err := ioutil.ReadFile(dir+"\\rule.txt")

	if rule_text_err == nil{
	    rule = string(rule_bytes)
	}else{
	 	rule = ""
	}
	if rule == ""{
		fmt.Println("请输入要监测的规则！！！！")
	}
	doc.Find(rule).Each(func(i int, s *goquery.Selection) {
    // For each item found, get the band and title
    band := s.Find("a").Text()
		result += string(band)
	})
	// 读取网页内容
	return
}


// 弹窗提醒
func message() {
	var LableHello=Label{
		Text: "网页有更新，请及时查看!", 
	  }
	  var widget=[]Widget{
		 LableHello,
	  }
	  var mainWindow=MainWindow{
		Title:"网页检测",
		Size:Size{Width: 300, Height: 100},
		Layout:VBox{}, 
		Children:widget,
	  }
	  mainWindow.Run()
}