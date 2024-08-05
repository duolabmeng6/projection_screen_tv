package myModel

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alexballas/go2tv/devices"
	"github.com/alexballas/go2tv/soapcalls"
	"github.com/alexballas/go2tv/soapcalls/utils"
	"os"
	"path/filepath"
	"sort"
	"strconv"
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
	ctx         context.Context
	fn          func(eventName, jsonstring string)
	播放状态        chan PlaybackEvent
	当前播放状态      string
}

func New投屏模块() *E投屏模块 {
	m := new(E投屏模块)
	m.serverStart = make(chan struct{})
	m.ctx = context.Background()
	m.播放状态 = make(chan PlaybackEvent)
	m.FileServer = New文件服务器("12306", m.播放状态)
	m.监听状态()
	return m
}
func (this *E投屏模块) 监听状态() {
	go func() {
		for {
			select {
			case msg := <-this.播放状态:
				fmt.Println("Received from channel1:", msg.Status, msg.UUID)
				data := map[string]string{}
				data["Status"] = msg.Status
				data["UUID"] = msg.UUID
				data["event"] = "playStatus"
				marshal, err := json.Marshal(data)
				this.当前播放状态 = msg.Status
				if err != nil {
					return
				}
				this.fn("playStatus", string(marshal))
			}
		}

	}()
	//每秒检查一次
	go func() {
		for {
			time.Sleep(1 * time.Second)
			//println("每秒检查一次", this.当前播放状态)

			if this.当前播放状态 != "PLAYING" {
				continue
			}
			当前位置, 总长度, 当前时间, 总时间 := this.E取当前播放位置()
			data := map[string]string{}
			data["currentPosition"] = strconv.Itoa(当前位置)
			data["overallLength"] = strconv.Itoa(总长度)
			data["currentTime"] = 当前时间
			data["totalEvent"] = 总时间
			data["event"] = "playPosition"
			marshal, err := json.Marshal(data)
			if err != nil {
				return
			}
			this.fn("playPosition", string(marshal))
		}
	}()
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

// 返回值 当前播放位置,总长度,当前时间,总时间
func (this *E投屏模块) E取当前播放位置() (int, int, string, string) {
	if this.tvData == nil {
		return 0, 0, "", ""
	}
	getPos, err := this.tvData.GetPositionInfo()
	if err != nil {
		fmt.Println(err)
		return 0, 0, "", ""
	}

	total, err := utils.ClockTimeToSeconds(getPos[0])
	if err != nil {
		fmt.Println(err)
		return 0, 0, "", ""

	}

	current, err := utils.ClockTimeToSeconds(getPos[1])
	if err != nil {
		fmt.Println(err)
		return 0, 0, "", ""

	}
	//println("total", total)
	//println("current", current)
	return current, total, getPos[1], getPos[0]
	//

	//valueToSet := float64(current) * s.SlideBar.Max / float64(total)
	//if !math.IsNaN(valueToSet) {
	//	println("valueToSet", valueToSet)
	//
	//	end, err := utils.FormatClockTime(getPos[0])
	//	if err != nil {
	//		return
	//	}
	//
	//	current, err := utils.FormatClockTime(getPos[1])
	//	if err != nil {
	//		return
	//	}
	//
	//	println("current", current)
	//	println("end", end)
	//}

}
func (this *E投屏模块) E设置播放进度(cur int) error {
	if this.tvData == nil {
		return errors.New("Haven't played yet")
	}
	//getPos, err := this.tvData.GetPositionInfo()
	//if err != nil {
	//	return
	//}
	//
	//total, err := utils.ClockTimeToSeconds(getPos[0])
	//if err != nil {
	//	return
	//}

	roundedInt := cur
	reltime, err := utils.SecondsToClockTime(roundedInt)
	if err != nil {
		return err
	}

	//end, err := utils.FormatClockTime(getPos[0])
	//if err != nil {
	//	return
	//}

	//println("reltime", reltime)
	//println("end", end)
	//println("total", total)

	if err := this.tvData.SeekSoapCall(reltime); err != nil {
		return err
	}
	return nil
}

func (this *E投屏模块) E投递视频文件(设备URL string, 文件路径 string, 字幕文件路径 string) error {
	//获取文件名
	fileName := filepath.Base(文件路径)

	播放文件URL := this.FileServer.E写文件名与路径(fileName, 文件路径)
	fileName = filepath.Base(字幕文件路径)
	字幕文件URL := this.FileServer.E写文件名与路径(fileName, 字幕文件路径)

	// 获取设备的 UPnP DLNA 媒体渲染器信息
	upnpServicesURLs, err := soapcalls.DMRextractor(this.ctx, 设备URL)
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
	//创建 TVPayload 结构体，包含投射需要的信息
	this.tvData = &soapcalls.TVPayload{
		ControlURL:                  upnpServicesURLs.AvtransportControlURL,
		EventURL:                    upnpServicesURLs.AvtransportEventSubURL,
		RenderingControlURL:         upnpServicesURLs.RenderingControlURL,
		ConnectionManagerURL:        upnpServicesURLs.ConnectionManagerURL,
		CallbackURL:                 this.FileServer.E取回调地址(),
		MediaURL:                    播放文件URL,
		SubtitlesURL:                字幕文件URL,
		MediaType:                   mediaType,
		CurrentTimers:               make(map[string]*time.Timer),
		MediaRenderersStates:        make(map[string]*soapcalls.States),
		InitialMediaRenderersStates: make(map[string]bool),
		//Transcode:                   *文件转码路径,
		//Seekable:                    isSeek,
		Transcode: false,
		Seekable:  false,
	}
	//this.tvData, err = soapcalls.NewTVPayload(&soapcalls.Options{
	//	DMR:       设备URL,
	//	Media:     播放文件URL,
	//	Subs:      字幕文件URL,
	//	Mtype:     mediaType,
	//	Logging:   nil,
	//	Transcode: false,
	//	Seek:      false,
	//})

	// 发送播放命令到 TV
	err = this.tvData.SendtoTV("Play1")
	if err != nil {
		return err
	}

	return nil
}
func (this *E投屏模块) E暂停播放() error {
	this.当前播放状态 = "PAUSED_PLAYBACK"
	if this.tvData == nil {
		return errors.New("Haven't played yet")
	}
	err := this.tvData.PlayPauseStopSoapCall("Pause")
	if err != nil {
		return err
	}

	return nil
}
func (this *E投屏模块) E停止播放() error {
	this.当前播放状态 = "STOPPED"
	if this.tvData == nil {
		return errors.New("Haven't played yet")
	}
	err := this.tvData.SendtoTV("Stop")
	if err != nil {
		return err
	}
	return nil
}
func (this *E投屏模块) E继续播放() error {
	if this.tvData == nil {
		return errors.New("Haven't played yet")
	}
	err := this.tvData.PlayPauseStopSoapCall("Play")
	if err != nil {
		return err
	}
	return nil
}

func (this *E投屏模块) E音量(up bool) error {
	if this.tvData == nil {
		return errors.New("Haven't played yet")
	}
	currentVolume, err := this.tvData.GetVolumeSoapCall()
	if err != nil {
		return err
	}
	setVolume := currentVolume - 1
	if up {
		setVolume = currentVolume + 1
	}

	if setVolume < 0 {
		setVolume = 0
	}

	stringVolume := strconv.Itoa(setVolume)

	if err := this.tvData.SetVolumeSoapCall(stringVolume); err != nil {
		return err
	}

	return nil
}
func (this *E投屏模块) E静音() error {
	if this.tvData == nil {
		return errors.New("Haven't played yet")
	}
	if err := this.tvData.SetMuteSoapCall("1"); err != nil {
		return err
	}
	return nil
}
func (this *E投屏模块) E取消静音() error {
	if this.tvData == nil {
		return errors.New("Haven't played yet")
	}
	if err := this.tvData.SetMuteSoapCall("0"); err != nil {
		return err
	}
	return nil
}

func (this *E投屏模块) E回调信息(回调信息 func(eventName, jsonstring string)) {
	this.fn = 回调信息
}
