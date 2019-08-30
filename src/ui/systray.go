package ui

import (
	"bing"
	"github.com/getlantern/systray"
	"github.com/go-toast/toast"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"time"
	"util"
)

const (
	APP_NAME = "BingWallpaper"
)

type updateSigns struct {
	TrayTooltip  chan bool
	Notification chan bool
}

var (
	dir, _      = filepath.Abs(".")
	IconOk      = dir + "\\res\\imgs\\ok_16px.png"
	IconSystray = dir + "\\res\\imgs\\icon_32px.ico"
	IconMsg     = dir + "\\res\\imgs\\icon_40px.png"
	sign        updateSigns
)

func tooltip() string {
	conf, _ := util.ReadConfiguration()
	return conf.Bing.Discription + "\n" + conf.Bing.Copyright
}

func loadIcon(path string) []byte {
	byteValue, _ := ioutil.ReadFile(path)
	return byteValue
}

//更新监测
func updateMonitors() {

	sign.Notification = make(chan bool, 1)
	sign.TrayTooltip = make(chan bool, 1)

	//监听更新信号
	go func() {
		for {
			<-bing.HasUpdated
			sign.Notification <- true
			sign.TrayTooltip <- true
		}
	}()

	//通知更新
	go func() {
		for {
			<-sign.Notification
			notification := toast.Notification{
				AppID:               APP_NAME,
				Title:               strconv.Itoa(time.Now().Year()) + "年" + strconv.Itoa(int(time.Now().Month())) + "月" + strconv.Itoa(time.Now().Day()) + "日",
				Message:             tooltip(),
				Icon:                IconMsg,
				ActivationType:      "",
				ActivationArguments: "",
				Actions:             nil,
				Audio:               "",
				Loop:                false,
				Duration:            "short",
			}
			_ = notification.Push()
		}
	}()

	//托盘信息更新
	go func() {
		for {
			systray.SetTooltip(tooltip())
			<-sign.TrayTooltip
		}
	}()
}

func onReady() {

	conf, _ := util.ReadConfiguration()

	systray.SetIcon(loadIcon(IconSystray))

	//更新域
	mUpdate := systray.AddMenuItem("更新壁纸", "")
	mWpDir := systray.AddMenuItem("打开壁纸夹", "")
	systray.AddSeparator()

	//收藏域
	mLike := systray.AddMenuItem("添加收藏", "")
	if isLike(conf) {
		mLike.Check()
		mLike.SetIcon(loadIcon(IconOk))
	} else {
		mLike.Uncheck()
	}
	mOpenLikeDir := systray.AddMenuItem("打开收藏夹", "123")
	systray.AddSeparator()

	//设置
	mSetting := systray.AddMenuItem("设置", "")
	systray.AddSeparator()

	//退出
	mQuit := systray.AddMenuItem("退出", "")

	updateMonitors()

	//TODO:开机自启可选关闭

	go func() {
		for {
			select {
			case <-mUpdate.ClickedCh:
				bing.Update()
			case <-mWpDir.ClickedCh:
				openWpDir(conf)
			case <-mLike.ClickedCh:
				if mLike.Checked() {
					cancelLike(conf)
					mLike.Uncheck()
					mLike.SetIcon([]byte{})
				} else {
					addLike(conf)
					mLike.Check()
					mLike.SetIcon(loadIcon(IconOk))
				}
			case <-mOpenLikeDir.ClickedCh:
				openLikeDir(conf)
			case <-mSetting.ClickedCh:
				createWinSetting()
			case <-mQuit.ClickedCh:
				systray.Quit()
			}
		}
	}()

}

func onExit() {
	close(bing.HasUpdated)
	close(sign.Notification)
	close(sign.TrayTooltip)
	//TODO:若设置窗口打开，则一并关闭
}

func Run() {
	systray.Run(onReady, onExit)
}
