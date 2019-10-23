package main

import (
	"bing"
	"config"
	"time"
)


func UpdateTask() {

	conf, _ := config.ReadConfigFile()

	//启动时更新一次，载入壁纸数据
	//TODO：考虑无网络的情况，也需要载入数据
	updatedate, err := bing.Update(conf.Wpdir)
	if err == nil {
		conf.Updatedate = updatedate
		_ = config.WriteConfigFile(conf)
	}

	var interval = 60	//检查更新的间隔
	for {
		conf, _ := config.ReadConfigFile()
		if conf.Updatedate != time.Now().Format("20060102") {
			updatedate, err := bing.Update(conf.Wpdir)
			if err == nil {
				conf.Updatedate = updatedate
				_ = config.WriteConfigFile(conf)
			}
		}
		time.Sleep(time.Duration(interval) * time.Second)
	}
}
