package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/nosixtools/solarlunar"
	"github.com/shawnfeng/sutil/slog"
	"github.com/shawnfeng/sutil/snetutil"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
	"time"
)

func dingding(dings []string, msg string) {
	fun := "dingding -->"

	for _, it := range dings {

		data, _ := json.Marshal(map[string]interface{}{
			"msgtype": "text",
			"text": map[string]interface{}{
				"content": msg,
			},
			/*
				"at": map[string]interface{}{
					"atMobiles": []string{
						phone,
					},
					"isAtAll": false,
				},
			*/
		})

		var header = make(map[string]string)
		header["Content-Type"] = "application/json"

		_, err := snetutil.HttpReqWithHeadOk(it, "POST", header, data, time.Second*10)
		if err != nil {
			slog.Warnf("%s alarm aid:%s err:%s", fun, msg, err)
		}

	}

}

type SignDate struct {
	// 12-01
	Date string `yaml:"date"`
	// is solar
	Solar   bool   `yaml:"solar"`
	Content string `yaml:"content"`
}

type NotifyInfo struct {
	Dings []string    `yaml:"dings"`
	Dates []*SignDate `yaml:"dates"`
}

type Config struct {
	Notifys []*NotifyInfo `yaml:"notifys"`
}

var configPath string

func main() {
	flag.StringVar(&configPath, "config", "", "yaml config file")
	flag.Parse()
	if len(configPath) == 0 {
		slog.Errorf("config file empty")
		return
	}

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		slog.Errorf("config file read err:" + err.Error())
		return
	}

	var config Config
	err = yaml.Unmarshal([]byte(data), &config)
	if err != nil {
		slog.Errorf("config file format error")
		return
	}

	solarDate := time.Now().Format("2006-01-02")
	lunarDate := solarlunar.SolarToSimpleLuanr(solarDate)
	cnDate := solarlunar.SolarToChineseLuanr(solarDate)

	solarMd := solarDate[len("2006-"):]

	lunarMd := strings.Replace(lunarDate, "年", "-", 1)
	lunarMd = strings.Replace(lunarMd, "月", "-", 1)
	lunarMd = strings.Replace(lunarMd, "日", "", 1)

	if len(lunarMd) == len("2006-01-02") {
		lunarMd = lunarMd[len("2006-"):]
	}

	slog.Infof("now solar:%s lunar:%s cnlunar:%s lunarmd:%s solarmd:%s", solarDate, lunarDate, cnDate, lunarMd, solarMd)

	for i, it := range config.Notifys {
		dings := it.Dings
		for j, jt := range it.Dates {
			slog.Infof("%d %d ding:%s date:%s solar:%v content:%s", i, j, dings, jt.Date, jt.Solar, jt.Content)
			if !jt.Solar && jt.Date == lunarMd {
				dingding(dings, fmt.Sprintf(jt.Content, cnDate))
			}

			if jt.Solar && jt.Date == solarMd {
				dingding(dings, fmt.Sprintf(jt.Content, solarDate))
			}
		}
	}

}
