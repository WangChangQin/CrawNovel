package main

import (
	"testing"
	"fmt"
	"github.com/crawnovel/config"
	"github.com/crawnovel/model"
)

var (
	rootCfg *config.QNovelConfig
	c       model.NovelConfig
	cfgs    map[string]config.NovelConfig
)

func init() {
	rootCfg = config.GetQNovelConfig()
}

func TestExistNovel(t *testing.T) {
	if rootCfg != nil && len(rootCfg.Engine) > 0 {
		cfgs = rootCfg.Engine
		c = model.NovelConfig{cfgs["88dushu"]}
	}
	n := "两界搬运工"

	b, s := c.QExist(n)
	fmt.Println(b, s)
}
func TestListNovel(t *testing.T) {
	if rootCfg != nil && len(rootCfg.Engine) > 0 {
		cfgs = rootCfg.Engine
		c = model.NovelConfig{cfgs["88dushu"]}
	}
	//n := "两界搬运工"
}
func TestItemNovel(t *testing.T) {
	if rootCfg != nil && len(rootCfg.Engine) > 0 {
		cfgs = rootCfg.Engine
		c = model.NovelConfig{cfgs["88dushu"]}
	}
	n := "两界搬运工"

	s := c.QItem(n)
	i := 0
	for _, ss := range s {
		if i > 5 {
			break
		}
		fmt.Println(ss)
		i++
	}
}
func TestItemContent(t *testing.T) {
	if rootCfg != nil && len(rootCfg.Engine) > 0 {
		cfgs = rootCfg.Engine
		c = model.NovelConfig{cfgs["88dushu"]}
	}
	n := "英雄联盟之传奇正盛"
	sec := c.QItem(n)
	fmt.Println(sec[0])
	fmt.Println(c.QContent(sec[0]))
}
