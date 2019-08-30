package ui

import (
	"bing"
	"io"
	"os"
	"os/exec"
	"syscall"
	"time"
	"util"
)

//壁纸添加收藏
func addLike(conf util.Configuration) error {
	src, err := os.Open(conf.Wpdir + conf.Updatedate + "_" + conf.Bing.Discription + bing.IMG_FMT)
	if err != nil {
		return err
	}
	defer src.Close()
	dst, err := os.OpenFile(conf.Likedir+conf.Updatedate+"_"+conf.Bing.Discription+bing.IMG_FMT,
		os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer dst.Close()
	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	} else {
		conf.Likedate = time.Now().Format("20060102")
		return util.WriteConfiguration(conf)
	}
}

//取消壁纸收藏
func cancelLike(conf util.Configuration) error {
	err := os.Remove(conf.Likedir + conf.Updatedate + "_" + conf.Bing.Discription + bing.IMG_FMT)
	if err != nil {
		return err
	} else {
		err = util.WriteConfiguration(conf)
		if err != nil {
			return err
		}
		return nil
	}
}

func isLike(conf util.Configuration) bool {
	if conf.Likedate != time.Now().Format("20060102") {
		return false
	} else {
		return true
	}
}

//打开壁纸夹
func openWpDir(conf util.Configuration) error {
	cmd := exec.Command(`cmd`, `/c`, `explorer`, conf.Wpdir)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	err := cmd.Start()
	if err != nil {
		return err
	} else {
		return nil
	}
}

//打开收藏夹
func openLikeDir(conf util.Configuration) error {
	cmd := exec.Command(`cmd`, `/c`, `explorer`, conf.Likedir)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	err := cmd.Start()
	if err != nil {
		return err
	} else {
		return nil
	}
}
