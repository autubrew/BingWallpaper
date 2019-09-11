package ui

import (
	"github.com/lxn/win"
	"github.com/sciter-sdk/go-sciter"
	"github.com/sciter-sdk/go-sciter/window"
	"os"
	"syscall"
	"config"
)

//“设置”窗口初始化工作
func winInit(w *window.Window) error {

	conf, _ := config.ReadConfigFile()

	//显示输入框内默认值
	root, _ := w.GetRootElement()
	input1, _ := root.SelectById("wpdir")
	err := input1.SetValue(sciter.NewValue(conf.Wpdir))
	if err != nil {
		return err
	}
	input2, _ := root.SelectById("likedir")
	err = input2.SetValue(sciter.NewValue(conf.Likedir))
	if err != nil {
		return err
	}
	return nil
}

//更新配置目录
func updateConfigDir(dirs... string) error {

	conf, _ := config.ReadConfigFile()
	if conf.Wpdir == dirs[0] && conf.Likedir == dirs[1] {
		return nil
	} else {
		conf.Wpdir = dirs[0]
		conf.Likedir = dirs[1]

		//生成新目录
		err := os.MkdirAll(dirs[0], os.ModePerm)
		if err != nil {
			//TODO:判断路径是否正确，错误则给出提示
			return err
		}
		err = os.MkdirAll(dirs[1], os.ModePerm)
		if err != nil {
			//TODO:判断路径是否正确，错误则给出提示
			return err
		}

		err = config.WriteConfigFile(conf)
		if err != nil {
			return err
		}

		return nil
	}
}

func createWinSetting() error {

	screenWidth, _, _ := syscall.NewLazyDLL(`User32.dll`).NewProc(`GetSystemMetrics`).Call(uintptr(0))
	screenHeight, _, _ := syscall.NewLazyDLL(`User32.dll`).NewProc(`GetSystemMetrics`).Call(uintptr(1))
	winWidth := 460
	winHeight := 180

	rect := sciter.NewRect((int(screenHeight)-winHeight)/2, (int(screenWidth)-winWidth)/2, winWidth, winHeight)
	w, err := window.New(
		//sciter.SW_TITLEBAR |		//顶级窗口，有标题栏
		//sciter.SW_GLASSY |		//可调整大小
		sciter.SW_CONTROLS| //有最大、最小按钮
			sciter.SW_MAIN, //应用程序主窗口，关闭后其他所有窗口也会关闭
		rect) //创建窗口的矩型
	if err != nil {
		return err
	}

	defer win.OleInitialize()

	w.SetTitle(APP_NAME)
	err = w.LoadFile(dir + "\\res\\setting.html") //html格式必须为带BOM的UTF8，否则中文会出现乱码
	if err != nil {
		return err
	}

	winInit(w)

	defFunc(w)

	w.Show()

	w.Run()

	return nil
}
