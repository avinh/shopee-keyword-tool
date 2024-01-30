package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type ConfigData struct {
	Cookie string `json:"cookie"`
}

type Api struct {
	Cookie string `json:"cookie"`
}

type ItemListProduct struct {
	ItemId int64 `json:"itemid"`
}

type DataListProduct struct {
	ItemList []ItemListProduct `json:"item_list"`
}

type ResponseListProduct struct {
	Message string          `json:"message"`
	Code    int             `json:"code"`
	Data    DataListProduct `json:"data"`
}

func (api *Api) GetListProducts() []ItemListProduct {
	url := "https://banhang.shopee.vn/api/marketing/v3/pas/product_selector/?offset=0&limit=50&is_ads=1&campaign_type=search"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return nil
	}
	req.Header.Add("Cookie", api.Cookie)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	var response ResponseListProduct
	json.Unmarshal(body, &response)

	return response.Data.ItemList
}

type SearchQuery struct {
	IsRecommended bool   `json:"is_recommended"`
	Keyword       string `json:"keyword"`
}

type ResponseListKeywordHint struct {
	Data []SearchQuery `json:"data"`
}

func (api *Api) GetListKeywordHint(keyword string, itemID int64) []SearchQuery {
	url := "https://banhang.shopee.vn/api/pas/v1/setup_helper/list_keyword_hint/"
	method := "POST"

	payload := strings.NewReader(fmt.Sprintf(`{"keyword":"%s","item_id":%d}`, keyword, itemID))
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return nil
	}
	req.Header.Add("Cookie", api.Cookie)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	var response ResponseListKeywordHint
	json.Unmarshal(body, &response)

	return response.Data
}

type KeywordData struct {
	IsRecommended    bool   `json:"is_recommended"`
	Keyword          string `json:"keyword"`
	RecommendedPrice int64  `json:"recommended_price"`
	Relevance        int    `json:"relevance"`
	SearchVolume     int    `json:"search_volume"`
	State            string `json:"state"`
}

type ResponseListKeywordData struct {
	Data []KeywordData `json:"data"`
}

func (api *Api) GetListKeywordData(keyword string, itemID int64) []KeywordData {
	url := "https://banhang.shopee.vn/api/pas/v1/setup_helper/search_keyword/"
	method := "POST"

	payload := strings.NewReader(fmt.Sprintf(`{"keyword":"%s","item_id":%d,"placement":0,"suggest_log_data":{"page":"suggest_creation"}}`, keyword, itemID))
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return nil
	}
	req.Header.Add("Cookie", api.Cookie)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	var response ResponseListKeywordData
	json.Unmarshal(body, &response)

	return response.Data
}

func loadConfigData() ConfigData {
	file, err := os.Open("config.json")
	if err != nil {
		fmt.Println("Không thể mở file cấu hình:", err)
		return ConfigData{}
	}
	defer file.Close()

	var config ConfigData
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println("Lỗi khi đọc cấu hình:", err)
		return ConfigData{}
	}

	return config
}

func appendOutputData(value string) {
	f, err := os.OpenFile("output.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Không thể mở file output:", err)
		return
	}
	defer f.Close()

	if _, err := f.WriteString(value + "\n"); err != nil {
		fmt.Println("Lỗi khi ghi file output:", err)
		return
	}
}

func clearFileOutput() {
	f, err := os.OpenFile("output.csv", os.O_TRUNC, 0644)
	if err != nil {
		// create file if not exist
		f, _ = os.Create("output.csv")
		f.WriteString("keyword,recommended_price\n")
		return
	}
	f.WriteString("keyword,recommended_price\n")
	defer f.Close()
}

func main() {
	configData := loadConfigData()
	clearFileOutput()

	api := Api(configData)

	fmt.Print("Vui lòng nhập keyword: ")
	reader := bufio.NewReader(os.Stdin)
	keyword, _ := reader.ReadString('\n')

	keyword = strings.TrimSpace(keyword)

	fmt.Println("Đang xử lý..." + keyword)

	listProducts := api.GetListProducts()

	if len(listProducts) == 0 {
		fmt.Println("Không tìm thấy sản phẩm nào!!!")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		return
	}

	listKeywordHint := api.GetListKeywordHint(keyword, listProducts[0].ItemId)

	if len(listKeywordHint) == 0 {
		fmt.Println("Không tìm thấy keyword nào phù hợp!!!")
		// wait for user press enter
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		return
	}

	for _, keywordHint := range listKeywordHint {
		dataKeywords := api.GetListKeywordData(keywordHint.Keyword, listProducts[0].ItemId)
		for _, dataKeyword := range dataKeywords {
			appendOutputData(dataKeyword.Keyword + "," + fmt.Sprintf("%d", dataKeyword.RecommendedPrice))
		}
	}

	fmt.Println("Đã xử lý xong!!!")
	// wait for user press enter
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
