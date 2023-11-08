package myModel

import (
	"encoding/json"
	"testing"
	"time"
)

func TestE投屏模块实现_获取设备列表(t *testing.T) {
	m := New投屏模块()
	m.E回调信息(func(eventName, jsonstring string) {
		if eventName == "playStatus" {
			println(jsonstring)
		}
		if eventName == "playPosition" {
			println(jsonstring)
		}
		// 播放进度 当前位置,总长度,当前时间,总时间
		// {"Status":"STOPPED","UUID":"9c39e130-6ca7-4bc4-b199-60c5acc261d4","event":"playStatus"}
		// Status PLAYING PAUSED_PLAYBACK STOPPED

		// {"currentPosition":"3673","currentTime":"1:59:09","event":"playPosition","overallLength":"7149","totalEvent":"1:01:13"}

	})
	// 播放状态 开始,暂停,停止
	设备列表, err := m.E获取设备列表()

	////设备列表转换为json
	json, err := json.Marshal(设备列表)
	println(string(json))
	设备URL := "http://192.168.100.234:57873/description.xml"
	err = m.E投递视频文件(设备URL, "/Users/ll/Downloads/2004-哈爾移動城堡【粤语】.mkv", "")
	if err != nil {
		println(err.Error())
	}

	//go func() {
	//	for {
	//		time.Sleep(1 * time.Second)
	//		m.E取当前播放位置()
	//	}
	//}()
	//time.Sleep(10 * time.Second)

	//m.E设置播放进度(1000)

	time.Sleep(300 * time.Second)

	m.E暂停播放()
	time.Sleep(10 * time.Second)
	m.E停止播放()

	err = m.E投递视频文件(设备URL, "/Users/ll/Downloads/1988-龍貓CD1-国语.mp4", "")
	if err != nil {
		println(err.Error())
	}
	// 10秒以后停止播放
	time.Sleep(10 * time.Second)
	m.E停止播放()

	time.Sleep(60 * time.Second)

}
