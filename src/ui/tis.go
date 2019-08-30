package ui

import (
	"fmt"
	"github.com/lxn/win"
	"github.com/sciter-sdk/go-sciter"
	"github.com/sciter-sdk/go-sciter/window"
	"syscall"
	"util"
)

func defFunc(w *window.Window)  {

	//点击浏览打开文件管理器选择指定目录
	w.DefineFunction("selectCustomDir", func(args ...*sciter.Value) *sciter.Value {

		conf, _ := util.ReadConfiguration()

		var path string
		id := args[0].String()
		if id == "wpdir" {
			path = conf.Wpdir
		} else if id == "likedir" {
			path = conf.Likedir
		}
		fmt.Println(path)

		bi := win.BROWSEINFO{
			HwndOwner:      0,
			PidlRoot:       0,	//TODO:自定义打开目录
			PszDisplayName: nil,
			LpszTitle:      nil,
			UlFlags:        0,
			Lpfn:           0,
			LParam:         0,
			IImage:         0,
		}
		pidl := win.SHBrowseForFolder(&bi)
		defer win.CoTaskMemFree(pidl)

		//获取选择的文件夹路径
		if pidl != 0 {
			var pszPath [win.MAX_PATH]uint16
			win.SHGetPathFromIDList(pidl, &pszPath[0])
			selectedPath := syscall.UTF16ToString(pszPath[:]) + "\\"
			fmt.Println(selectedPath)
			root, _ := w.GetRootElement()
			tag, _ := root.SelectById(args[0].String())
			tag.SetValue(sciter.NewValue(selectedPath))
		}
		return sciter.NullValue()
	})

	//“确定”按钮，更新配置信息
	w.DefineFunction("ok", func(args ...*sciter.Value) *sciter.Value {
		updateConfigDir(args[0].String(), args[1].String())
		return sciter.NullValue()
	})

}
