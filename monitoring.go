package main

import (
	"fmt"
	"os"
	"net/http"
	"path/filepath"
	"io"
	"bytes"
  	"github.com/PuerkitoBio/goquery"
	"time"
	// "github.com/lxn/walk"
	// . "github.com/lxn/walk/declarative"
)


func main() {
	var ch chan int
	//定时任务
	ticker := time.NewTicker(time.Second * 5)
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
	fmt.Println(older1_err)

	if older1_err == nil {
		os.Remove(older1_path)
	}

	_, older2_err := os.Stat(older2_path)
	fmt.Println(older2_err)

	if older2_err == nil {
		os.Rename(older2_path, older1_path)
	}

	_, err := os.Stat(new_path)
	fmt.Println(err)
	if err == nil{
		re_err := os.Rename(new_path, older2_path)
		fmt.Println(re_err)

	}		
	os.Create(dir+"\\new.txt")
	spider(new_path)

	if compare(new_path, older2_path) != true {
		fmt.Println("数据有更新！请查看")
		
	}
}

func compareLines(line1, line2 string) {
	sign := "o"
	if line1 != line2 {
		sign = "x"
	}
	fmt.Printf("%s | %s | %s n", sign, line1, line2)
}

// 对比是否有更新
func compare(spath, dpath string) bool {
	sinfo, err := os.Lstat(spath)
	if err != nil {
			return false
	}
	dinfo, err := os.Lstat(dpath)
	if err != nil {
			return false
	}
	if sinfo.Size() != dinfo.Size() || !sinfo.ModTime().Equal(dinfo.ModTime()) {
			return false
	}
	return comparefile(spath, dpath)
}

func comparefile(spath, dpath string) bool {
	sFile, err := os.Open(spath)
	if err != nil {
			return false
	}
	dFile, err := os.Open(dpath)
	if err != nil {
			return false
	}
	b := comparebyte(sFile, dFile)
	sFile.Close()
	dFile.Close()
	return b
}
//下面可以代替md5比较.
func comparebyte(sfile *os.File, dfile *os.File) bool {
	var sbyte []byte = make([]byte, 512)
	var dbyte []byte = make([]byte, 512)
	var serr, derr error
	for {
			_, serr = sfile.Read(sbyte)
			_, derr = dfile.Read(dbyte)
			if serr != nil || derr != nil {
					if serr != derr {
							return false
					}
					if serr == io.EOF {
							break
					}
			}
			if bytes.Equal(sbyte, dbyte) {
					continue
			}
			return false
	}
	return true
}

// 爬取
func spider(path string) {
	url := "https://tieba.baidu.com/f?kw=%E6%AD%A6%E6%B1%89"
	result, _ := HttpGet(url)
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
	
	doc.Find(".threadlist_title.pull_left.j_th_tit").Each(func(i int, s *goquery.Selection) {
    // For each item found, get the band and title
    band := s.Find("a").Text()
		result += string(band)
	})
	// 读取网页内容
	return
}


// 弹窗提醒
