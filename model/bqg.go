package model

type BqgCfg struct {
	Cfg NovelConfig
}

func (b BqgCfg) QExist(novelName string) (bool, string) {
	return b.Cfg.QExist(novelName)
}

func (b BqgCfg) QList(novelName string) string {
	return b.Cfg.QList(novelName)
}

func (b BqgCfg) QItem(novelName string) []string {
	return b.Cfg.QItem(novelName)
}

func (b BqgCfg) QContent(sectionUrl string) string {
	return b.Cfg.QContent(sectionUrl)
}

