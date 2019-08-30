package main

import (
	"bing"
	"time"
	"util"
)

func checkDate() bool {
	conf, _ := util.ReadConfiguration()
	return conf.Updatedate == time.Now().Format("20060102")
}

func UpdateTask() {

	if !checkDate() {	//启动时确定是否需要更新一次
		_ = bing.Update()
	}
	for { //每天0点触发
		h1, m1, s1 := time.Now().Clock()
		t := (23-h1)*3600 + (59-m1)*60 + (60 - s1) + 2 //加2s的延迟误差，保证在0点后更新
		time.Sleep(time.Duration(t) * time.Second)
		_ = bing.Update()
	}
}
