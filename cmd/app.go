package cmd

import (
	"os"
	"fmt"
	"github.com/urfave/cli"
	"sync"
	"sort"
	"time"
	"strings"
	"strconv"
	"math"
	"io"
	"github.com/crawnovel/config"
	"github.com/crawnovel/model"
	"github.com/crawnovel/util"
)

type CliInfo struct {
	NovelName  string
	Outputfile string
	Piece      string
	Source     string
	Web        string
}
type novelPiece struct {
	num     int
	content []string
	err     []error
}
type novels struct {
	fullnovel []novelPiece
}

func (n novels) Len() int {
	return len(n.fullnovel)
}

func (n novels) Less(i, j int) bool {
	return n.fullnovel[i].num < n.fullnovel[j].num
}

func (n novels) Swap(i, j int) {
	n.fullnovel[i], n.fullnovel[j] = n.fullnovel[j], n.fullnovel[i]
}
func Run() {
	app := cli.NewApp()
	app.Name = "novel-cli"
	app.Usage = "小说爬取"
	app.Version = rootCfg.Version

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "name,n",
			Value: "善良的死神",
			Usage: "指定下载小说的名字",
		},
		cli.StringFlag{
			Name:  "outputfile,o",
			Value: "小说名.txt",
			Usage: "可选输出的文件名",
		},
		cli.StringFlag{
			Name:  "piece,p",
			Value: "-1",
			Usage: "可选下载部分章节，例如200-500",
		},
		cli.StringFlag{
			Name:  "source,s",
			Value: "3",
			Usage: "可选:\n" +
				"\tbqg2\t1\t" +
				"\tbqg\t2\t" +
				"\t88dushu\t3",
		},
		cli.StringFlag{
			Name:  "web,w",
			Value: "www",
			Usage: "网址",
		},
	}

	app.Action = func(c *cli.Context) error {
		var cliInfo CliInfo
		cliInfo.NovelName = c.String("name")
		cliInfo.Outputfile = c.String("outputfile")
		cliInfo.Piece = c.String("piece")
		cliInfo.Web = c.String("web")
		source := c.String("source")
		if source == "1" {
			cliInfo.Source = "bqg2"
		} else if source == "2" {
			cliInfo.Source = "bqg"
		} else if source == "3" {
			cliInfo.Source = "88dushu"
		} else if source == "4" {
			cliInfo.Source = "sxyj"
		}
		if cliInfo.NovelName == "善良的死神" {
			fmt.Fprintln(os.Stdout, "请输入要下载的文件名")
			return nil
		}
		spider(cliInfo)
		return nil
	}
	app.Run(os.Args)
}

var (
	rootCfg   *config.QNovelConfig
	cfg       model.NovelConfig
	cfgs      map[string]config.NovelConfig
	per_works = 10
)

func init() {
	rootCfg = config.GetQNovelConfig()
	per_works = rootCfg.MaxWorks
}
func spider(cliInfo CliInfo) {
	var novelUrl string
	if rootCfg != nil && len(rootCfg.Engine) > 0 {
		cfgs = rootCfg.Engine
		c := cfgs[cliInfo.Source]
		cfg = model.NovelConfig{c}
	}

	useWebUrl := cliInfo.Web != "www"
	if !useWebUrl {
		exist, url := handleSwitchConfig(cliInfo) //"https://www.88dushu.com/xiaoshuo/61/61992/"
		if !exist {
			fmt.Fprintln(os.Stdout, "小说不存在")
			os.Exit(1)
		}
		novelUrl = url
	} else {
		novelUrl = cliInfo.Web
	}

	fmt.Fprintln(os.Stdout, "小说存在,即将下载 listurl", novelUrl)
	t1 := time.Now()

	var sections []string
	var novelPieces = make(chan novelPiece, 10)
	var myNovels novels
	var length int
	sections = cfg.QItem(cliInfo.NovelName)
	if cliInfo.Piece != "-1" {
		tb := strings.Split(cliInfo.Piece, "-")
		top, _ := strconv.Atoi(tb[0])
		bottom, _ := strconv.Atoi(tb[1])
		if bottom >= len(sections) {
			bottom = len(sections)
		}
		sections = sections[top:bottom]
	}
	//sections = sections[0:50]
	length = len(sections)
	gorunnum := int(math.Ceil(float64(length * 1.0 / per_works)))
	fmt.Printf("有%d个章节,共分配%d个任务\n", length, gorunnum)
	var wg sync.WaitGroup
	for i := 0; i < gorunnum; i++ {
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			fmt.Println("创建一个gorunnale id = ", num)
			var n novelPiece
			n.num = num
			for _, sectionUrl := range sections[num*per_works : (num+1)*per_works] {
				fmt.Println("处理的url是", sectionUrl)
				content := cfg.QContent(sectionUrl)
				if len(content) > 0 {
					n.content = append(n.content, content)
				}
			}
			novelPieces <- n

		}(i)
	}
	go func() {
		wg.Wait()
		close(novelPieces)
	}()
	for p := range novelPieces {
		if p.err != nil {
			fmt.Fprintf(os.Stdout, "发现异常%v\n", p.err)
			continue
		}
		myNovels.fullnovel = append(myNovels.fullnovel, p)
		fmt.Println("收集了一段小说")
	}
	sort.Sort(myNovels)
	var saveName string
	for _, v := range myNovels.fullnovel {
		saveName =  handleSave(cliInfo, v.content...)
	}
	fmt.Println("下载小说成功")
	elapsed := time.Since(t1)
	if rootCfg.UserMail {
		util.SendMail(rootCfg.KindleMail,saveName)
	}
	fmt.Println("App elapsed: ", elapsed)
}

func handleSwitchConfig(cliInfo CliInfo) (bool, string) {
	fmt.Println("当前小说源==>", cliInfo.Source)
	exist, url := cfg.QExist(cliInfo.NovelName)
	if exist {
		return exist, url
	}
	for k, c := range cfgs {
		cfg = model.NovelConfig{c}
		if k == cliInfo.Source {
			continue
		}
		fmt.Println("不存在，切换小说源==>", k)
		exist, url = cfg.QExist(cliInfo.NovelName)
		if exist {
			break
		}
	}
	return exist, url
}
func handleSave(clinfo CliInfo, res ...string)string{
	saveName := clinfo.Outputfile
	if saveName == "小说名.txt" || saveName == "" {
		saveName = clinfo.NovelName + ".txt"
	} else {
		saveName = fmt.Sprintf("%s.txt", saveName)
	}

	f, err := os.OpenFile(saveName, os.O_RDWR|os.O_APPEND|os.O_SYNC|os.O_CREATE, 0755)

	if err != nil {
		fmt.Fprintln(os.Stderr, "access file fail")
		os.Exit(1)
	}
	if len(res) > 0 {
		for _, s := range res {
			io.WriteString(f, s)
		}
	}
	return saveName
}
