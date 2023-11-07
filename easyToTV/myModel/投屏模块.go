package myModel

import (
	"errors"
	"fmt"
	"github.com/alexballas/go2tv/devices"
	"github.com/alexballas/go2tv/soapcalls"
	"github.com/alexballas/go2tv/utils"
	"os"
	"path"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

type E投屏模块接口 interface {
	E获取设备列表() string
	E投递视频文件(设备名称 string, 文件路径 string) string
	E暂停播放(设备名称 string) string
	E停止播放(设备名称 string) string
}

type E投屏模块 struct {
	E投屏模块接口
	serverStart chan struct{}
	tvData      *soapcalls.TVPayload
	FileServer  *E文件服务器
}

func New投屏模块() *E投屏模块 {
	m := new(E投屏模块)
	m.serverStart = make(chan struct{})
	m.FileServer = New文件服务器("12306")

	return m
}
func (this *E投屏模块) E获取设备列表() ([]map[string]string, error) {

	deviceList, err := devices.LoadSSDPservices(1)
	if err != nil {
		return nil, errors.New("failed to list devices")
	}

	keys := make([]string, 0)
	for k := range deviceList {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	//for q, k := range keys {
	//	fmt.Printf("%sDevice %v%s\n", q+1)
	//	fmt.Printf("%s--------%s\n")
	//	fmt.Printf("%sModel:%s %s\n", k)
	//	fmt.Printf("%sURL:%s   %s\n", deviceList[k])
	//	fmt.Println()
	//}
	//定义 变量m 为 map 作为数组 插入 Model URL 的数据
	var m []map[string]string
	//遍历 deviceList
	for _, k := range keys {
		//定义 变量n 为 map 作为数组 插入 Model URL 的数据
		var n = map[string]string{}
		//插入数据
		n["Model"] = k
		n["URL"] = deviceList[k]
		//插入数组
		m = append(m, n)
	}

	return m, nil
}

func (this *E投屏模块) E投递视频文件(设备URL string, 文件路径 string, 字幕文件路径 string) error {
	//获取文件名
	fileName := path.Base(文件路径)
	播放文件URL := this.FileServer.E写文件名与路径(fileName, 文件路径)
	fileName = path.Base(字幕文件路径)
	字幕文件URL := this.FileServer.E写文件名与路径(fileName, 字幕文件路径)

	// 获取设备的 UPnP DLNA 媒体渲染器信息
	upnpServicesURLs, err := soapcalls.DMRextractor(设备URL)
	if err != nil {
		return fmt.Errorf("获取设备信息时发生错误：%w", err)
	}

	// 解析设备的监听地址和端口
	//whereToListen, err := utils.URLtoListenIPandPort(设备URL)
	//if err != nil {
	//	return fmt.Errorf("解析设备监听地址时发生错误：%w", err)
	//}

	// 获取文件的绝对路径
	absFilePath, err := filepath.Abs(文件路径)
	if err != nil {
		return fmt.Errorf("获取文件绝对路径时发生错误：%w", err)
	}

	// 打开文件
	file, err := os.Open(absFilePath)
	if err != nil {
		return fmt.Errorf("打开文件时发生错误：%w", err)
	}
	defer file.Close()

	// 获取文件的 MIME 类型
	mediaType, err := utils.GetMimeDetailsFromFile(file)
	if err != nil {
		return fmt.Errorf("获取文件 MIME 类型时发生错误：%w", err)
	}
	//callbackPath, err := utils.RandomString()
	println(播放文件URL)
	// 创建 TVPayload 结构体，包含投射需要的信息
	this.tvData = &soapcalls.TVPayload{
		ControlURL:                  upnpServicesURLs.AvtransportControlURL,
		EventURL:                    upnpServicesURLs.AvtransportEventSubURL,
		RenderingControlURL:         upnpServicesURLs.RenderingControlURL,
		CallbackURL:                 this.FileServer.E取回调地址(),
		MediaURL:                    播放文件URL,
		SubtitlesURL:                字幕文件URL,
		MediaType:                   mediaType,
		CurrentTimers:               make(map[string]*time.Timer),
		MediaRenderersStates:        make(map[string]*soapcalls.States),
		InitialMediaRenderersStates: make(map[string]bool),
		RWMutex:                     &sync.RWMutex{},
		//Transcode:                   *文件转码路径,
		//Seekable:                    isSeek,
	}

	// 发送播放命令到 TV
	if err := this.tvData.SendtoTV("Play1"); err != nil {
		return fmt.Errorf("发送播放命令时发生错误：%w", err)
	}

	return nil
}
func (this *E投屏模块) E暂停播放() error {
	err := this.tvData.AVTransportActionSoapCall("Pause")
	if err != nil {
		return err
	}

	return nil
}
func (this *E投屏模块) E停止播放() error {
	err := this.tvData.SendtoTV("Stop")
	if err != nil {
		return err
	}
	return nil
}
func (this *E投屏模块) E跳转时间() error {

	return nil
}
func (this *E投屏模块) 加音量() error {

	return nil
}
func (this *E投屏模块) 减音量() error {

	return nil
}

func (this *E投屏模块) 回调信息() error {

	return nil
}
