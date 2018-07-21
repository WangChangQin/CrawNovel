package _interface


type QNovel interface {
	QExist(novelName string) (bool,string)
	QList(novelName string) string
	QItem(novelName string) []string
	QContent(sectionUrl string) string
}
