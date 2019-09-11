package main

import (
	"bing"
	"time"
	"config"
)


func UpdateTask() {

	conf, _ := config.ReadConfigFile()

	if conf.Updatedate != time.Now().Format("20060102") {	//启动时确定是否需要更新一次
		updatedate, err := bing.Update(conf.Wpdir)
		if err == nil {
			conf.Updatedate = updatedate
			config.WriteConfigFile(conf)
		}
	}
	for { //每天0点触发
		h1, m1, s1 := time.Now().Clock()
		t := (23-h1)*3600 + (59-m1)*60 + (60 - s1) + 2 //加2s的延迟误差，保证在0点后更新
		time.Sleep(time.Duration(t) * time.Second)
		updatedate, err := bing.Update(conf.Wpdir)
		if err == nil {
			conf.Updatedate = updatedate
			config.WriteConfigFile(conf)
		}
	}
}
