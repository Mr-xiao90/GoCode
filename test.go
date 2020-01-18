package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	result := make([]byte, 4096)
	buf := make([]byte, 4096)
	url := "https://xgao5.com/album/59979/%E7%A7%81%E4%BA%BA%E7%8E%A9%E7%89%A9%E5%B0%91%E5%A5%B3%E8%89%B2cos%E5%9F%83%E7%BD%97%E8%8A%92%E9%98%BF%E8%80%81%E5%B8%88"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("htt.get err =", err)
		return
	}
	defer resp.Body.Close()
	for {
		Len, err1 := resp.Body.Read(buf)
		if err1 != nil && err1 != io.EOF {
			fmt.Println("body.read err =", err1)
			return
		}
		if Len == 0 {
			break
		}
		result = append(result, buf[:Len]...)
	}
	fmt.Printf("result%s", result)

	//fmt.Println("多少页", pageNum)
}
