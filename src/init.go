package main

import (
	"config"
	"golang.org/x/sys/windows/registry"
	"os"
	"path/filepath"
	"ui"
)

func path2Abs(pathname string) string {
	if filepath.IsAbs(pathname) {
		return pathname
	} else {
		dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		return dir + "\\" + pathname + "\\"
	}
}

//注册表开机自启
func autoSatrtup() error {
	key, err := registry.OpenKey(registry.CURRENT_USER, "Software\\Microsoft\\Windows\\CurrentVersion\\Run", registry.ALL_ACCESS)
	if err != nil {
		return err
	}
	defer key.Close()
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	err = key.SetStringValue(ui.APP_NAME, `"` + dir + "\\" + ui.APP_NAME + `.exe"`)
	if err != nil {
		return err
	}
	return nil
}

func Init() error {

	//读取配置文件
	var conf config.Configuration
	conf, err := config.ReadConfigFile()
	//fmt.Println(bing.Conf)
	if err != nil {
		return err
	}

	//将壁纸目录修改为绝对路径
	conf.Likedir = path2Abs(conf.Likedir)
	conf.Wpdir = path2Abs(conf.Wpdir)

	//创建文件夹
	err = os.MkdirAll(conf.Likedir, os.ModePerm)
	if err != nil {
		return err
	}
	err = os.MkdirAll(conf.Wpdir, os.ModePerm)
	if err != nil {
		return err
	}

	err = config.WriteConfigFile(conf)
	if err != nil {
		return err
	}

	//添加启动项
	//autoSatrtup()

	return nil
}
