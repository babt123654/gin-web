package v1

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jordan-wright/email"
	"google.golang.org/api/gmail/v1"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// 工具：err处理
func CheckErr(err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
	//var f string
	//switch format.(type) {
	//case string:
	//	f = format.(string)
	//case error:
	//	f = fmt.Sprintf("%v", format.(error))
	//}
	//if f != "" {
	//	FailWithMsg(f, a...)
	//}
}

// 工具：判断两个[]byte对应的json值转换成数字是否相等 相等返回true
func JsonEqual(a, b []byte) bool {
	var va, vb interface{}
	if err := json.Unmarshal(a, &va); err != nil {
		return false
	}
	if err := json.Unmarshal(b, &vb); err != nil {
		return false
	}
	return reflect.DeepEqual(va, vb)
}

// 工具：邮件发送
func SendEmail(fromAddress string, subject string, body string, toAddress string) gmail.Message {
	e := email.NewEmail()
	e.From = fromAddress
	//"1871437892@qq.com"
	e.To = []string{toAddress}
	e.Subject = subject
	e.Text = []byte(body)
	e.Send("smtp.qq.com:25", smtp.PlainAuth("", "1334642655@qq.com", "cmhtvatsahmvbafg", "smtp.qq.com"))
	return gmail.Message{}
}

// 工具：获取json中的字段
func GetFieldValueFromJSON(data []byte, field string) (string, []byte, error) {
	var jsonData map[string]interface{}
	err := json.Unmarshal(data, &jsonData)
	CheckErr(err)
	fieldValue, ok := jsonData[field]
	if !ok {
		return "", nil, fmt.Errorf("field '%s' not found in JSON", field)
	}
	fieldValueBytes, err := json.Marshal(fieldValue)
	CheckErr(err)
	return field, fieldValueBytes, nil
}

// 模拟浏览器访问网站
func ReqwebByUrl(url string) (io.ReadCloser, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	return response.Body, err
}
func ReqByUrl(url string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf(url), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	//fmt.Println(string(body))
	return body, err
}

// 根据正则表达式裁剪字符串
func ApplyRegexString(targetString, regexPattern string) string {
	r := regexp.MustCompile(regexPattern)
	match := r.FindStringSubmatch(targetString)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

// []byte转float64
func Byte2float(bytes []byte) float64 {
	str := string(bytes)
	floatValue, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0
	}
	return floatValue
}

// sendMessage TG机器人推送 是一个发送消息的方法
func SendMessage(bot *tgbotapi.BotAPI, chatID int64, message string) error {
	// 创建一个HTTP客户端
	client := &http.Client{}

	// 构建Telegram Bot API的URL
	apiURL := "https://api.telegram.org/bot" + bot.Token + "/sendMessage"

	// 构建POST请求的表单数据
	formData := url.Values{}
	formData.Set("chat_id", string(chatID))
	formData.Set("text", message)

	// 发送POST请求
	req, err := http.NewRequest("POST", apiURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("请求失败，状态码：%d", resp.StatusCode)
	}

	log.Println("消息已成功发送！")
	return nil
}
