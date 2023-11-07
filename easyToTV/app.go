package main

import (
	"changeme/myModel"
	"context"
	"encoding/json"
	"fmt"
	"github.com/duolabmeng6/goefun/ecore"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx  context.Context
	投屏模块 *myModel.E投屏模块
}

// NewApp creates a new App application struct
func NewApp() *App {
	a := &App{}
	a.投屏模块 = myModel.New投屏模块()

	return a
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	//隐藏窗口
	//runtime.WindowShow(a.ctx)
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	println("收到js的调用信息")
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) E获取设备列表() string {
	设备列表, err := a.投屏模块.E获取设备列表()
	if err != nil {
		return err.Error()
	}
	//[{"Model":"华为智慧屏 S65","URL":"http://192.168.100.204:25826/description.xml"},{"Model":"奇异果极速投屏-华为(204)","URL":"http://192.168.100.204:39620/description.xml"},{"Model":"MacBook Pro","URL":"http://192.168.10scription.xml"}]

	json, err := json.Marshal(设备列表)
	println(string(json))
	return string(json)
}

func (a *App) E获取系统时间() string {
	println("E获取系统时间")
	t := ecore.E取现行时间().E时间到文本("Y-m-d H:i:s")
	return t
}
func (a *App) OpenFileDialog() string {
	文件路径, _ := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "选择文件",
	})
	return 文件路径
}
func (a *App) E投递视频文件(设备URL string, 媒体文件路径 string, 字幕文件 string) string {
	println("E投递视频文件", 设备URL, 媒体文件路径, 字幕文件)
	err := a.投屏模块.E投递视频文件(设备URL, 媒体文件路径, 字幕文件)
	if err != nil {
		return err.Error()
	}
	return "ok"
}

func (a *App) E停止播放() string {
	err := a.投屏模块.E停止播放()
	if err != nil {
		return err.Error()
	}
	return "ok"
}

func (a *App) E暂停播放() string {
	err := a.投屏模块.E暂停播放()
	if err != nil {
		return err.Error()
	}
	return "ok"
}

func (a *App) GetVersion() string {
	println("GetVersion", myModel.Version)
	return myModel.Version
}
