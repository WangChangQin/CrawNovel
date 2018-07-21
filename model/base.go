package model

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"os"
	"github.com/crawnovel/util"
	"github.com/crawnovel/config"
)

type NovelConfig struct {
	config.NovelConfig
}

func (n NovelConfig) QExist(novelName string) (bool, string) {
	url := fmt.Sprintf(n.SearchUrl, novelName)
	//fmt.Println(url)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		fmt.Println(err)
		return false, ""
	}
	var listUrl []string
	doc.Find(n.ExistSel).Each(func(i int, sel *goquery.Selection) {
		//fmt.Println(sel.Attr("href"))
		//fmt.Println(sel.Text())
		if sel.Text() == novelName {
			url, _ = sel.Attr("href")
			var prefix string
			if !util.HasHttpPrefix(url) {
				prefix = n.Host
			}
			listUrl = append(listUrl, prefix+url)
		}
	})

	if len(listUrl) > 0 {
		return true, listUrl[0]
	}
	return false, ""
}

func (n NovelConfig) QList(novelName string) string {
	b, url := n.QExist(novelName)
	if !b {
		return ""
	}
	fmt.Println("小说地址:", url)
	return url
}

func (n NovelConfig) QItem(novelName string) []string {
	listUrl := n.QList(novelName)

	if listUrl == "" {
		return nil
	}
	var sections []string
	doc, err := goquery.NewDocument(listUrl)
	if err != nil {
		return nil
	}
	doc.Find(n.ListSel).Each(func(i int, sel *goquery.Selection) {
		url, _ := sel.Attr("href")
		var prefix string
		if !util.HasHttpPrefix(url) {
			prefix = listUrl
		}
		sections = append(sections, prefix+url)
	})
	if len(sections) > 0 {
		return sections
	}
	return nil
}

func (n NovelConfig) QContent(sectionUrl string) string {
	doc, err := goquery.NewDocument(sectionUrl)
	if err != nil {
		return ""
	}
	novel, err := doc.Find(n.ContentSel).Html()
	title := doc.Find(n.TitleSel).Text()
	novel = title + "\n" + novel
	if err != nil {
		fmt.Fprintf(os.Stderr, "find: %v\n", err)
		return ""
	}
	if n.ConvertCode {
		novel = util.ConvertToString(novel, "gbk", "utf-8")
	}
	novel = util.ReplaceNovelStr(novel)
	return novel
}
