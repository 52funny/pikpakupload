package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/52funny/pikpakupload/conf"
	"github.com/52funny/pikpakupload/model"
	"github.com/52funny/pikpakupload/utils"
	"github.com/sirupsen/logrus"
)

type Exn []string

func (i *Exn) String() string {
	return fmt.Sprint(*i)
}
func (i *Exn) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var parentPath = ""
var exn = Exn{}
var defaultRegexp []*regexp.Regexp = []*regexp.Regexp{
	regexp.MustCompile(`^\..+`),
}

var parentId = ""
var sync bool

func init() {
	if os.Getenv("PIKPAKUPLOAD_DEBUG") != "" {
		logrus.SetLevel(logrus.DebugLevel)
	}
	err := conf.InitConfig()
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
	flag.StringVar(&parentPath, "p", "", "parent path")
	flag.Var(&exn, "exn", "exclude path")
	flag.BoolVar(&sync, "sync", false, "sync")
	flag.Parse()
	for _, r := range exn {
		reg, err := regexp.Compile(r)
		if err != nil {
			logrus.Error(err)
			os.Exit(1)
		}
		defaultRegexp = append(defaultRegexp, reg)
	}
}

func main() {
	path := os.Args[len(os.Args)-1]

	uploadFilePath := utils.GetUploadFilePath(path, defaultRegexp)
	var f *os.File

	// sync mode
	if sync {
		file, err := os.OpenFile(filepath.Join(".", ".pikpaksync.txt"), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
		if err != nil {
			logrus.Error(err)
			os.Exit(1)
		}
		f = file
		bs, err := ioutil.ReadAll(f)
		if err != nil {
			logrus.Error("read file error: ", err)
			os.Exit(1)
		}
		alreadySyncFiles := strings.Split(string(bs), "\n")
		files := make([]string, 0)
		for _, f := range uploadFilePath {
			if !utils.Contains(alreadySyncFiles, f) {
				files = append(files, f)
			}
		}
		uploadFilePath = files
	}

	logrus.Infoln("upload file list:")
	for _, f := range uploadFilePath {
		logrus.Infoln(f)
	}

	p := model.NewPikPak(conf.Config.Username, conf.Config.Password)

	err := p.Login()
	if err != nil {
		logrus.Error(err)
	}
	err = p.AuthCaptchaToken("POST:/drive/v1/files")
	if err != nil {
		logrus.Error(err)
	}

	go func() {
		ticker := time.NewTicker(time.Second * 7200 * 3 / 4)
		defer ticker.Stop()
		for range ticker.C {
			err := p.RefreshToken()
			if err != nil {
				logrus.Warn(err)
				continue
			}
		}
	}()

	if parentPath != "" {
		parentPathS := strings.Split(parentPath, "/")
		for i, v := range parentPathS {
			if v == "." {
				parentPathS = append(parentPathS[:i], parentPathS[i+1:]...)
			}
		}
		id, err := p.GetDeepParentOrCreateId(parentId, parentPathS)
		if err != nil {
			logrus.Error(err)
			os.Exit(-1)
		} else {
			parentId = id
		}
	}
	logrus.Debug("parentPath: ", parentPath, " parentId: ", parentId)

	for _, v := range uploadFilePath {
		if strings.Contains(v, "/") {
			basePath := filepath.Dir(v)
			basePathS := strings.Split(basePath, "/")
			id, err := p.GetDeepParentOrCreateId(parentId, basePathS)
			if err != nil {
				logrus.Error(err)
			}
			err = p.UploadFile(id, filepath.Join(path, v))
			if err != nil {
				logrus.Error(err)
			}
			if sync {
				f.WriteString(v + "\n")
			}
			logrus.Infof("%s upload completed!\n", v)
		} else {
			err = p.UploadFile(parentId, filepath.Join(path, v))
			if err != nil {
				logrus.Error(err)
			}
			if sync {
				f.WriteString(v + "\n")
			}
		}
	}
}
