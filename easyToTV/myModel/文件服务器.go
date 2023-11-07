package myModel

import (
	"github.com/alexballas/go2tv/soapcalls"
	"github.com/gin-gonic/gin"
	"html"
	"io"
	"net"
	"net/http"
	"net/url"
	"path"
	"strings"
)

type E文件服务器接口 interface {
	E获取设备列表() string
	E写文件名与路径(文件名, 路径 string)
	E取路径(文件名 string) string
	E清空()
}

type E文件服务器 struct {
	E文件服务器接口
	文件与路径      map[string]string
	router     *gin.Engine
	port       string
	serverAddr string
}

func New文件服务器(port string) *E文件服务器 {
	m := new(E文件服务器)
	m.port = port
	m.文件与路径 = make(map[string]string)
	m.初始化()

	return m
}

func (this *E文件服务器) 初始化() {
	this.router = gin.Default()

	this.router.NoRoute(func(c *gin.Context) {
		//判断是不是 callback
		if c.Request.URL.Path != "/callback" {
			c.String(http.StatusNotFound, "404 page not found")
			return
		}

		// 读取请求体
		reqParsed, _ := io.ReadAll(c.Request.Body)

		// 获取并检查请求头中的 Sid
		sidVal, sidExists := c.Request.Header["Sid"]
		if !sidExists || sidVal[0] == "" {
			c.String(http.StatusNotFound, "Sid header not found or empty")
			return
		}

		// 提取 uuid，去除 "uuid:" 前缀
		uuid := strings.TrimPrefix(sidVal[0], "uuid:")

		// 对请求体进行 HTML 反转义，并解析为 UPnP 事件
		reqParsedUnescape := html.UnescapeString(string(reqParsed))
		previousstate, newstate, err := soapcalls.EventNotifyParser(reqParsedUnescape)
		if err != nil {
			c.String(http.StatusNotFound, "Error parsing UPnP event")
			return
		}
		println("uuid", uuid, "previousstate", previousstate, "newstate", newstate)

		if newstate == "STOPPED" {
			//tv.SetProcessStopTrue(uuid)
			c.String(http.StatusOK, "OK\n")
			return
		}

		// 根据新状态执行相应的操作
		switch newstate {
		case "PLAYING":
			// 执行播放时的操作
			println("PLAYING")
		case "PAUSED_PLAYBACK":
			// 执行暂停时的操作
			println("PAUSED_PLAYBACK")
		case "STOPPED":
			// 执行停止时的操作
			println("STOPPED")
		}
	})
	this.router.GET("/file/:urlPath", func(c *gin.Context) {
		urlPath := c.Param("urlPath")
		urlPath, _ = url.QueryUnescape(urlPath)

		filePath := this.E取路径(urlPath)

		// Check if the file path exists
		if filePath == "" {
			c.String(http.StatusNotFound, "文件不存在")
			return
		}

		fileExt := path.Ext(filePath)
		contentType := "application/octet-stream" // Default content type

		switch fileExt {
		case ".mp4":
			contentType = "video/mp4"
		case ".mkv":
			contentType = "video/x-matroska"
		case ".avi":
			contentType = "video/x-msvideo"
		case ".mov":
			contentType = "video/quicktime"
		case ".srt":
			contentType = "text/plain; charset=utf-8"
		case ".ass":
			contentType = "text/plain; charset=utf-8"
			// Add more formats as needed
		}

		c.Header("Content-Disposition", "inline")
		c.Header("Content-Type", contentType)

		c.File(filePath)
	})

	// 获取本地IP地址
	localIP, err := 获取本地IP()
	if err == nil {
		this.serverAddr = localIP + ":" + this.port
	}

	// Run the server on 0.0.0.0:6161
	go func() {
		this.router.Run(":" + this.port)
	}()
}
func (this *E文件服务器) E写文件名与路径(文件名, 路径 string) string {
	this.文件与路径[文件名] = 路径
	return "http://" + this.serverAddr + "/file/" + 文件名
}

func (this *E文件服务器) E取路径(文件名 string) string {
	return this.文件与路径[文件名]

}
func (this *E文件服务器) E清空() {
	this.文件与路径 = make(map[string]string)

}

func (this *E文件服务器) E取回调地址() string {
	return "http://" + this.serverAddr + "/callback"

}

func 获取本地IP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			return ipNet.IP.String(), nil
		}
	}

	return "", nil
}
