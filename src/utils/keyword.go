package utils

import (
	"HANG-backend/src/global"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"strings"
)

var API_KEY string
var SECRET_KEY string

func InitKey() {
	API_KEY = viper.GetString("baidu.api_key")
	SECRET_KEY = viper.GetString("baidu.secret_key")
}

func SetKeyword(word string) {
	url := "https://aip.baidubce.com/rpc/2.0/nlp/v1/txt_keywords_extraction?access_token=" + GetAccessToken()
	payload := strings.NewReader(fmt.Sprintf(`{"text":["%s"]}`, word))

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	words := ExtractList(string(body))
	_ = UpsertKeywords(words)
}

/**
 * 使用 AK，SK 生成鉴权签名（Access Token）
 * @return string 鉴权签名信息（Access Token）
 */
func GetAccessToken() string {
	url := "https://aip.baidubce.com/oauth/2.0/token"
	postData := fmt.Sprintf("grant_type=client_credentials&client_id=%s&client_secret=%s", API_KEY, SECRET_KEY)
	resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(postData))
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	accessTokenObj := map[string]any{}
	_ = json.Unmarshal([]byte(body), &accessTokenObj)
	return accessTokenObj["access_token"].(string)
}

type Result struct {
	Word string `json:"word"`
}

type Data struct {
	Results []Result `json:"results"`
}

func ExtractList(jsonStr string) []string {
	var data Data
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		fmt.Println("解析JSON出错:", err)
		return []string{}
	}
	words := make([]string, len(data.Results))
	for i, r := range data.Results {
		words[i] = r.Word
	}
	return words
}

func UpsertKeywords(keywords []string) error {
	for _, keyword := range keywords {
		var existingKeyword Keyword
		result := global.RDB.Where("name =?", keyword).First(&existingKeyword)
		if result.Error == gorm.ErrRecordNotFound {
			// 如果没找到，就创建新的关键词记录
			newKeyword := Keyword{
				Name:  keyword,
				Count: 1,
			}
			err := global.RDB.Create(&newKeyword).Error
			if err != nil {
				return fmt.Errorf("创建关键词失败: %v", err)
			}
		} else if result.Error != nil {
			return fmt.Errorf("查询关键词时出错: %v", result.Error)
		} else {
			// 如果找到了，就把Count加1并更新
			existingKeyword.Count++
			err := global.RDB.Save(&existingKeyword).Error
			if err != nil {
				return fmt.Errorf("更新关键词Count失败: %v", err)
			}
		}
	}
	return nil
}

type Keyword struct {
	ID    uint   `gorm:"primaryKey;autoIncrement;not null" json:"id"`
	Name  string `gorm:"not null;index" json:"name"`
	Count int    `gorm:"default:0" json:"count"`
}

func ListAllKeywords() ([]Keyword, error) {
	var keywords []Keyword
	result := global.RDB.Find(&keywords)
	if result.Error != nil {
		return nil, result.Error
	}
	return keywords, nil
}
