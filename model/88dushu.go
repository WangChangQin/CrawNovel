package model


type Dushu88Cfg struct {
	Cfg NovelConfig
}

func (d Dushu88Cfg) QExist(novelName string) (bool, string) {
	return d.Cfg.QExist(novelName)
}

func (d Dushu88Cfg) QList(novelName string) string {
	return d.Cfg.QList(novelName)
}

func (d Dushu88Cfg) QItem(novelName string) []string {
	return d.Cfg.QItem(novelName)
}

func (d Dushu88Cfg) QContent(sectionUrl string) string {
	return d.Cfg.QContent(sectionUrl)
}

