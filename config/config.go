package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"fmt"
)

type QNovelConfig struct {
	Version    string                 `yaml:"version"`
	MaxWorks   int                    `yaml:"max-goroutine"`
	KindleMail string                 `yaml:"kindle-mail"`
	UserMail   bool                 `yaml:"use-mail"`
	Engine     map[string]NovelConfig `yaml:"engines"`
}
type NovelConfig struct {
	Host          string `yaml:"host"`
	SearchUrl     string `yaml:"searchUrl"`
	ExistSel      string `yaml:"existSel"`
	DetailListSel string `yaml:"detailListSel"`
	ListSel       string `yaml:"listSel"`
	TitleSel      string `yaml:"titleSel"`
	ContentSel    string `yaml:"contentSel"`
	ConvertCode   bool   `yaml:"convertCode"`
}

var qConfig *QNovelConfig

func init() {
	var c QNovelConfig
	qConfig = c.getConf()
}
func GetQNovelConfig() *QNovelConfig {
	return qConfig
}
func (conf *QNovelConfig) getConf() *QNovelConfig {
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		fmt.Println(err.Error())
	}
	err = yaml.UnmarshalStrict(yamlFile, conf)

	if err != nil {
		fmt.Println(err.Error())
	}

	return conf
}
