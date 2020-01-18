package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

func downloadImg(k int, v string, di chan int, dirPath string) {

	oneImgUrl := "https://xgao5.com/media/photos/" + v + ".jpg"
	fmt.Println(oneImgUrl)
	imgData := httpGet(oneImgUrl)
	fmt.Println("获取数据")
	fileName := dirPath + "/" + v + ".jpg"
	fName, err1 := os.Create(fileName)
	if err1 != nil {
		fmt.Println("文件创建错误", err1)
		return
	}
	defer fName.Close()
	_, err := fName.Write([]byte(imgData))
	if err != nil {
		fmt.Println("图片下载错误=", err)
		return
	}
	di <- k
}

func statWork(title string) {
	var dirPath string
	di := make(chan int, 100)
	imgUrl := "https://xgao5.com/album/" + title
	fmt.Println("拿到url=", imgUrl)
	dirName := getFileName(httpGet(imgUrl))
	for _, v := range dirName {
		fmt.Println("name=", v[1])
		err := os.Mkdir(v[1], os.ModePerm)
		if err != nil {
			fmt.Println("创建文件夹错误=", err)
		}
		dirPath = v[1]
		fmt.Println("dirPath=", dirPath)
		path, err := os.Getwd()
		if err != nil {
			fmt.Println("获取当前路径失败：", err)
			return
		}
		dirPath = path + "/" + v[1]
	}
	fmt.Println("imgPath", dirPath)

	pagNum := getOnePageUrl(httpGet(imgUrl))
	for _, v := range pagNum {
		pageUrl := "https://xgao5.com/album/"
		Num := num(httpGet(imgUrl))
		for _, v := range Num {
			pageUrl = "https://xgao5.com/album/" + v[1] + "/?page="
		}
		fmt.Printf("开始第%s\n", v[1])
		pageUrl = pageUrl + v[1]
		fmt.Println("完整url", pageUrl)
		imgNum := regImg(httpGet(pageUrl))
		for k, v := range imgNum {
			go downloadImg(k, v[1], di, dirPath)
		}
		for i := 0; i < len(imgNum); i++ {
			fmt.Printf("第%d张图片下载完成\n", <-di)
		}
		fmt.Println("单个网页里面的所有url下载完成")
	}
	fmt.Println("内循环完成")
}

func num(urlString string) (regResult [][]string) {
	reg := "/album/slideshow/([0-9]*)\">"
	r := regexp.MustCompile(reg)
	regResult = r.FindAllStringSubmatch(urlString, -1)
	return
}

func getOnePageUrl(urlString string) (regResult [][]string) {
	reg := "\">([0-9])</"
	r := regexp.MustCompile(reg)
	regResult = r.FindAllStringSubmatch(urlString, -1)
	return
}

func getFileName(urlString string) (regResult [][]string) {
	reg := "<title>(.*?) "
	r := regexp.MustCompile(reg)
	regResult = r.FindAllStringSubmatch(urlString, -1)
	return
}

func regImg(urlString string) (regResult [][]string) {
	reg := "/media/photos/([0-9]*)"
	r := regexp.MustCompile(reg)
	regResult = r.FindAllStringSubmatch(urlString, -1)
	return
}

func regTitle(urlString string) (regResult [][]string) {
	reg := "<a href=\"/album/(.*?)\">"
	r := regexp.MustCompile(reg)
	regResult = r.FindAllStringSubmatch(urlString, -1)
	return
}

func httpGet(url string) (result string) {
	buf := make([]byte, 4096)
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("网页获取错误")
		return
	}
	defer res.Body.Close()
	for {
		n, err := res.Body.Read(buf)
		if n == 0 {
			break
		}
		if err != io.EOF && err != nil {
			return
		}
		result = result + string(buf[:n])
	}
	return
}

func main() {
	for {
		i := 1
		url := "https://xgao5.com/albums?page=" + strconv.Itoa(i)
		title := regTitle(httpGet(url)) //获取单个相册的名字
		for _, v := range title {
			fmt.Println("main里面的=", v[1])
			statWork(v[1])
			fmt.Println("下载完成名叫", v[1])
		}
	}
}
