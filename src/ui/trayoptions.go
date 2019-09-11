package ui

import (
	"bing"
	"errors"
	"io"
	"os"
	"os/exec"
	"syscall"
	"time"
	"config"
)

//壁纸添加收藏
func addLike(conf config.Configuration) error {
	if conf.Updatedate != time.Now().Format("20060102") { //核实今日是否更新，不更新无法收藏
		return errors.New("not updated")
	}
	src, err := os.Open(conf.Wpdir + bing.GetWallpaperName())
	if err != nil {
		return err
	}
	defer src.Close()
	dst, err := os.OpenFile(conf.Likedir + bing.GetWallpaperName(), os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer dst.Close()
	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}
	return nil
}

//取消壁纸收藏
func cancelLike(conf config.Configuration) error {
	err := os.Remove(conf.Likedir + bing.GetWallpaperName())
	if err != nil {
		return err
	} else {
		return nil
	}
}

//判断当日的壁纸是否被收藏
func isLike(conf config.Configuration) bool {
	if conf.Updatedate != time.Now().Format("20060102") { //未更新，则今日必未收藏
		return false
	}
	fileinfo, err := os.Stat(conf.Likedir + bing.GetWallpaperName())
	if err == nil && !fileinfo.IsDir() {
		return true
	} else {
		return false
	}
}

//打开壁纸夹
func openWpDir(conf config.Configuration) error {
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
func openLikeDir(conf config.Configuration) error {
	cmd := exec.Command(`cmd`, `/c`, `explorer`, conf.Likedir)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	err := cmd.Start()
	if err != nil {
		return err
	} else {
		return nil
	}
}
