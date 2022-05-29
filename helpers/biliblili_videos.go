package helpers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type RespBodyDataListVListData struct {
	Title string
	BVid  string
}

type RespBodyDataList struct {
	VList []RespBodyDataListVListData
}

type RespBodyData struct {
	List RespBodyDataList
}

type RespBody struct {
	Data RespBodyData
}

func GetVideos() []string {
	logger := log.Default()
	resp, err := http.Get("https://api.bilibili.com/x/space/arc/search?mid=627888730&ps=30&tid=0&pn=1&keyword=&order=pubdate&jsonp=jsonp")
	if err != nil {
		logger.Fatal(err)
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	var bodyJson RespBody
	_err := json.Unmarshal(bodyBytes, &bodyJson)
	if _err != nil {
		logger.Fatal(err)
	}
	VList := bodyJson.Data.List.VList
	var UrlList []string

	for _, Video := range VList {
		UrlList = append(UrlList, "https://www.bilibili.com/video/"+Video.BVid)
	}
	return UrlList
}
