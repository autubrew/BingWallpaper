package bing

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"regexp"
	"syscall"
	"time"
	"unsafe"
)

type Bing struct {
	Images [] struct {
		Startdate     string   `json:"startdate"`
		Fullstartdate string   `json:"fullstartdate"`
		Enddate       string   `json:"enddate"`
		Url           string   `json:"url"`
		Urlbase       string   `json:"urlbase"`
		Copyright     string   `json:"copyright"`
		Copyrightlink string   `json:"copyrightlink"`
		Title         string   `json:"title"`
		Quiz          string   `json:"quiz"`
		Wp            bool     `json:"wp"`
		Hsh           string   `json:"hsh"`
		Drk           int      `json:"drk"`
		Top           int      `json:"top"`
		Bot           int      `json:"bot"`
		Hs            []string `json:"hs"`
	} `json:"images"`
	Tooltips struct {
		Loading  string `json:"loading"`
		Previous string `json:"previous"`
		Next     string `json:"next"`
		Walle    string `json:"walle"`
		Walls    string `json:"walls"`
	} `json:"tooltips"`
}

var (
	bing       Bing                 //壁纸相关信息
	wpname     string               //生成的壁纸名。格式：日期_内容
	wpinfo     string               //壁纸的信息。包括内容和copyright
	updating   = false              //判断程序是否正在更新壁纸，确保同时只有一个更新实例在运行
	SigUpdated = make(chan bool, 1) //更新成功信号
)

const (
	IMG_FMT = ".jpg"
)

//返回生成的图片名
func GetWallpaperName() string {
	return wpname
}

//返回壁纸信息
func GetWallpaperInfo() string {
	return wpinfo
}

//根据返回的JSON生成文件名，格式为：日期_壁纸描述（备选：Hsh_壁纸描述）
func generateWallpaperName() string {
	reg := regexp.MustCompile(`.* \(©`)
	description := string(reg.Find([]byte(bing.Images[0].Copyright)))
	description = description[0 : len(description)-4] //删除尾部“ (©”
	if description != "" {
		return bing.Images[0].Enddate + "_" + description + IMG_FMT
	}
	return bing.Images[0].Enddate + "_" + bing.Images[0].Hsh + IMG_FMT
}

//下载壁纸保存到指定目录
//savedir对应于配置文件的wpdir项，会转成绝对路径
func downloadWallpaper(savedir string) error {

	//根据bing开放api请求返回得到JSON
	resp, err := http.Get("https://cn.bing.com/HPImageArchive.aspx?format=js&idx=-1&n=1&mkt=zh-CN")
	if err != nil {
		return err
	}
	err = json.NewDecoder(resp.Body).Decode(&bing)
	if err != nil {
		return err
	}

	//解析图片真实地址
	url := "https://cn.bing.com" + bing.Images[0].Url
	//fmt.Println(url)

	//生成文件路径，文件名格式为：日期_壁纸内容描述
	if savedir[len(savedir)-1] != '\\' { //壁纸保存的目录，以“\”结尾
		savedir += "\\"
	}
	wpname = generateWallpaperName()
	wpabspath := savedir + wpname

	resp, err = http.Get(url)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	reader := bufio.NewReaderSize(resp.Body, 32*1024)

	file, err := os.Create(wpabspath)
	defer file.Close() //及时释放，否则后面systemParametersInfo的fWinIni会因为文件占用而无法SPIF_SENDCHANGE
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(file)
	_, err = io.Copy(writer, reader)
	if err != nil {
		return err
	}

	return nil
}

func setWallpaper(wpabspath string) bool {

	user32 := syscall.NewLazyDLL("user32.dll")
	systemParametersInfo := user32.NewProc("SystemParametersInfoW")
	filenameUTF16, _ := syscall.UTF16PtrFromString(wpabspath)

	ret, _, _ := systemParametersInfo.Call(
		uintptr(0x0014),                        //uiAction = SPI_SETDESKWALLPAPER
		uintptr(0x0000),                        //uiparam = 0
		uintptr(unsafe.Pointer(filenameUTF16)), //pvParam 指向壁纸文件
		uintptr(0x01|0x02),                     //fWinIni = SPIF_UPDATEINIFILE | SPIF_SENDCHANGE
	)

	return ret != 0
}

//更新壁纸
//只允许同时有一个Update实例进行更新
//若更新成功则返回更新时间和nil，同时设置桌面壁纸；否则返回空时间且给出error，说明正在有其他实例进行更新
//建议使用使用goroutine
func Update(savedir string) (string, error) {

	if !updating {
		updating = true
	} else {
		return "", errors.New("there's an update going on")
	}

	//下载壁纸
	for { //遇网络问题则循环等待
		err := downloadWallpaper(savedir)
		if err == nil {
			break
		} else {
			time.Sleep(5 * time.Second)
		}
	}

	//设置壁纸
	setWallpaper(savedir + wpname)

	reg := regexp.MustCompile(`.* \(©`)
	description := string(reg.Find([]byte(bing.Images[0].Copyright)))
	description = description[0 : len(description)-4] //删除尾部“ (©”
	copyright := bing.Images[0].Copyright[len(description)+2 : len(bing.Images[0].Copyright)-1]
	wpinfo = description + "\n" + copyright

	SigUpdated <- true
	updating = false

	return bing.Images[0].Enddate, nil
}