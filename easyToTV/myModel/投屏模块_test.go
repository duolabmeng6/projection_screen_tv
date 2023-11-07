package myModel

import (
	"encoding/json"
	"testing"
	"time"
)

func TestE投屏模块实现_获取设备列表(t *testing.T) {
	m := New投屏模块()
	设备列表, err := m.E获取设备列表()
	//println(设备列表)
	println(err)
	////设备列表转换为json
	json, err := json.Marshal(设备列表)
	println(string(json))
	设备URL := "http://192.168.100.234:57873/description.xml"
	err = m.E投递视频文件(设备URL, "/Users/ll/Downloads/2004-哈爾移動城堡【粤语】.mkv", "")
	if err != nil {
		println(err.Error())
	}
	time.Sleep(10 * time.Second)
	m.E暂停播放(设备URL)
	time.Sleep(10 * time.Second)
	m.E停止播放(设备URL)

	err = m.E投递视频文件(设备URL, "/Users/ll/Downloads/1988-龍貓CD1-国语.mp4", "")
	if err != nil {
		println(err.Error())
	}
	// 10秒以后停止播放
	time.Sleep(10 * time.Second)
	m.E停止播放(设备URL)

	time.Sleep(60 * time.Second)

}
