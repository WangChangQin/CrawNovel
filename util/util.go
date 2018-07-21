package util

import (
	"os"
	"github.com/axgle/mahonia"
	"strings"
	"regexp"
	"net/textproto"
	"net/smtp"
	"github.com/jordan-wright/email"
	"fmt"
)

/**
 * 判断文件是否存在  存在返回 true 不存在返回false
 */
func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}
func HasHttpPrefix(url string) bool {
	return strings.HasPrefix(url, "http")
}
func ConvertToString(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}
func ReplaceNovelStr(novel string) string {
	novel = strings.Replace(novel, "<br/>", "\n", -1)
	novel = strings.Replace(novel, "聽", "", -1)
	novel = strings.Replace(novel, "&nbsp", " ", -1)
	reg := regexp.MustCompile("<.*>")
	reg2 := regexp.MustCompile("&lt;(p|/p)?&gt;")
	novel = reg.ReplaceAllString(novel, "")
	novel = reg2.ReplaceAllString(novel, "\n")
	return novel
}
func SendMail(kindleMail, saveName string){
	e := &email.Email{
		To:      []string{kindleMail},
		From:    "yourmail",
		Subject: "convert",
		//Text: []byte("Text Body is, of course, supported!"),
		//HTML: []byte("<h1>Fancy HTML is supported, too!</h1>"),
		Headers: textproto.MIMEHeader{},
	}
	e.AttachFile(saveName)

	err := e.Send("smtp.qq.com:587", smtp.PlainAuth("", "yourmail", "pass", "smtp.qq.com"))

	if err != nil {
		fmt.Println("发送Kindle成功")
	} else {
		fmt.Println("发送Kindle失败")
	}
}
