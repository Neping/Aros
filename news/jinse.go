package news

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"chos/utils"
)

func HttpPostJson(url string, json string) string {
    jsonStr :=[]byte(json)
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
    req.Header.Set("Content-Type", "application/json")
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        fmt.Print(err);
		return "";
    }
    defer resp.Body.Close()
    //statuscode := resp.StatusCode
    //hea := resp.Header
    body, _ := ioutil.ReadAll(resp.Body)
    //fmt.Println(string(body))
    //fmt.Println(statuscode)
    //fmt.Println(hea)

    return string(body)
}

const DING_ROBOT_URL = "https://oapi.dingtalk.com"
const DING_ROBOT_SEND_URI = "/robot/send?access_token="
func MarkdownData(title string, text string, atMobiles string, isAtAll string) string  {
	json	:= `{
		 "msgtype": "markdown",
		 "markdown": {
			 "title":"{{title}}",
			 "text": "{{text}}"
		 },
		  "at": {
			  "atMobiles": [
				  {{atMobiles}}
			  ],
			  "isAtAll": {{isAtAll}}
		  }
	 }`

	reqDt 	:= strings.Replace(json,"{{title}}", title, -1)
	reqDt = strings.Replace(reqDt,"{{text}}", text, -1)
	reqDt = strings.Replace(reqDt,"{{atMobiles}}", atMobiles, -1)
	reqDt = strings.Replace(reqDt,"{{isAtAll}}", isAtAll, -1)

	return reqDt
}

func TextData(content string, atMobiles string, isAtAll string) string {
	json := `
		{
			"msgtype": "text", 
			"text": {
				"content": "{{content}}"
			}, 
			"at": {
				"atMobiles": [
				  {{atMobiles}}
			  	],
			  	"isAtAll": {{isAtAll}}
			}
		}`

	reqDt 	:= strings.Replace(json,"{{content}}", content, -1)
	reqDt = strings.Replace(reqDt,"{{atMobiles}}", atMobiles, -1)
	reqDt = strings.Replace(reqDt,"{{isAtAll}}", isAtAll, -1)

	return reqDt
}

/**
 * link类型
 */
func LinkData(text string, title string, picUrl string, messageUrl string) string {
	json := `{
			"msgtype": "link", 
			"link": {
				"text": "{{text}}", 
				"title": "{{title}}", 
				"picUrl": "{{picUrl}}", 
				"messageUrl": "{{messageUrl}}"
			}
		}`

	reqDt 	:= strings.Replace(json,"{{text}}", text, -1)
	reqDt = strings.Replace(reqDt,"{{title}}", title, -1)
	reqDt = strings.Replace(reqDt,"{{picUrl}}", picUrl, -1)
	reqDt = strings.Replace(reqDt,"{{messageUrl}}", messageUrl, -1)

	return reqDt
}

/**
 * 整体跳转ActionCard类型
 */
func ActionCard1Data(title string, text string, btnOrientation string, singleTitle string, singleURL string) string {
	json := `{
			"actionCard": {
				"title": "{{title}}", 
				"text": "{{text}}", 
				"btnOrientation": "{{btnOrientation}}", 
				"singleTitle" : "{{singleTitle}}",
				"singleURL" : "{{singleURL}}"
			}, 
			"msgtype": "actionCard"
		}`

	reqDt 	:= strings.Replace(json,"{{text}}", text, -1)
	reqDt = strings.Replace(reqDt,"{{title}}", title, -1)
	reqDt = strings.Replace(reqDt,"{{btnOrientation}}", btnOrientation, -1)
	reqDt = strings.Replace(reqDt,"{{singleTitle}}", singleTitle, -1)
	reqDt = strings.Replace(reqDt,"{{singleURL}}", singleURL, -1)

	return reqDt
}

/**
 * 独立跳转ActionCard类型@todo
 "btns": [
            {
                "title": "内容不错",
                "actionURL": "https://www.dingtalk.com/"
            },
        ]
 */
func ActionCard2Data() string {
	json := `{
		"actionCard": {
			"title": "{{title}}", 
			"text": "{{text}}", 
			"btnOrientation": "{{btnOrientation}}", 
			"btns": [
				{{btns}}
			]
		}, 
		"msgtype": "actionCard"
	}`

	return json
}

/**
 * FeedCard类型@todo
 * [{
       "title": "时代的火车向前开2",
       "messageURL": "https://www.dingtalk.com/s",
       "picURL": "https://gw.alicdn.com/tfs/TB1ayl9mpYqK1RjSZLeXXbXppXa-170-62.png"
   },]
 */
func FeedCardData() string {
	json := `{
		"feedCard": {
			"links": [
				{{links}}
			]
		}, 
		"msgtype": "feedCard"
	}`

	return json
}


// 金色财经
const JINSE_NEWSFLASH = "https://api.jinse.com/v4/live/list?limit=20&reading=false&source=web&sort=&flag=down&id=0"
type Jinse struct {
	News int `json:"news"`
	Count int `json:"count"`
	Total int `json:"total"`
	TopID int `json:"top_id"`
	BottomID int `json:"bottom_id"`
	List []Jinse_List `json:"list"`
	DefaultShareImg string `json:"default_share_img"`
	PrefixLink string `json:"prefix_link"`
}
type Jinse_Data struct {
	Symbol string `json:"symbol"`
	Slug string `json:"slug"`
	Change24H float64 `json:"change_24h"`
}
type WordBlocks struct {
	Type string `json:"type"`
	Data Jinse_Data `json:"data"`
}
type Jinse_Lives struct {
	ID int `json:"id"`
	Content string `json:"content"`
	ContentPrefix interface{} `json:"content_prefix"`
	LinkName string `json:"link_name"`
	Link string `json:"link"`
	Grade int `json:"grade"`
	Sort string `json:"sort"`
	HighlightColor string `json:"highlight_color"`
	Images []interface{} `json:"images"`
	CreatedAt int `json:"created_at"`
	Attribute string `json:"attribute"`
	UpCounts int `json:"up_counts"`
	DownCounts int `json:"down_counts"`
	ZanStatus string `json:"zan_status"`
	Readings []interface{} `json:"readings"`
	ExtraType int `json:"extra_type"`
	Extra interface{} `json:"extra"`
	Prev interface{} `json:"prev"`
	Next interface{} `json:"next"`
	WordBlocks []WordBlocks `json:"word_blocks"`
	IsShowComment int `json:"is_show_comment"`
	IsForbidComment int `json:"is_forbid_comment"`
	CommentCount int `json:"comment_count"`
	AnalystUser interface{} `json:"analyst_user"`
}
type Jinse_List struct {
	Date string `json:"date"`
	Lives []Jinse_Lives `json:"lives"`
}

func (news *Jinse)Newsflash()  {
	resp, _ := http.Get(JINSE_NEWSFLASH)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(body))
	err := json.Unmarshal(body, &news)
    if err != nil {
      log.Fatal(err)
    }
}

func WriteFile(filename , data string) {
    fileObj,err := os.Create(filename)
    if nil != err {
        log.Fatal(err)
    }
    writer := bufio.NewWriter(fileObj)
    defer writer.Flush()
    writer.WriteString(data)
}

func ReadFile(filename string) string {
    var (
        err error
        content []byte
    )
    content,err = ioutil.ReadFile(filename)
    if err != nil {
        WriteFile(filename, "")
    	return "0"
    }

    return string(content)
}

const NEWS_INDEX_HEIGHT_FILE = "./runtime/cache/jinse-push-height.log"
func GetLastIndex() int {
	readIndex := ReadFile(NEWS_INDEX_HEIGHT_FILE)
	oldIndex, err := strconv.Atoi(readIndex)
    if err != nil {
    	//log.Fatal(err)
    	return 0
    }

    return oldIndex
}

const DING_ROBOT_CONF = "./.config/dingding.conf"
type DingConf struct {
	RobotName string `json:"robotName"`
	AccessToken string `json:"accessToken"`
	Keywords string `json:"keywords"`
}
func GetAccessToken() DingConf {
	var conf DingConf
	str := ReadFile(DING_ROBOT_CONF)
	err := json.Unmarshal([]byte(str), &conf)
    if err != nil {
      log.Fatal(err)
    }

    return conf
}

// 金色财经快讯推送入口
func JinsePush()  {
	var jinse = Jinse{}
	jinse.Newsflash()
	dingConf := GetAccessToken()
	maxIndex 	:= 0;
	gradeLevel 	:= 5;
	for _, row := range jinse.List {
		for _, news := range row.Lives {
			if gradeLevel == news.Grade {
				lastIndex := GetLastIndex()
				maxIndex   = utils.IntMax(maxIndex, news.ID)
				if maxIndex>lastIndex {
					reqUrl 		:= DING_ROBOT_URL + DING_ROBOT_SEND_URI + dingConf.AccessToken
					title 		:= dingConf.Keywords
					text 		:= "> " + news.Content
					atMobiles 	:= ""
					isAtAll 	:= "false"
					reqDt  		:= MarkdownData(title, text, atMobiles, isAtAll)
					HttpPostJson(reqUrl, reqDt)
					// 更新index height
					WriteFile(NEWS_INDEX_HEIGHT_FILE, strconv.Itoa(maxIndex))
					fmt.Println("Dingtalk Robot", dingConf.RobotName, "Send Success, Jinse News ID: ", news.ID)
					goto FLAG
				}
			}

		}
	}
	FLAG:
}

func JinseTest()  {
	var jinse = Jinse{}
	jinse.Newsflash()
	for _, row := range jinse.List {
		for _, news := range row.Lives {
			fmt.Println(news.Content)
		}
	}
}

func Run()  {
	for  {
		JinsePush()
		time.Sleep(time.Duration(10)*time.Second)
	}
}