import { ElMessage, ElMessageBox } from "element-plus";
import * as systemFc from "../../wailsjs/runtime";
import * as goFc from "../../wailsjs/go/main/App";

export function BindWindowEvent(c, comps) {
    c.WinCreated = async function () {
        console.log("WinCreated")
        comps.Win.text = "多多投屏" + await goFc.GetVersion();
        comps.Button_停止播放.disable = true;
        comps.按钮_音量加.disable = true;
        comps.按钮_音量减.disable = true;
        comps.按钮_静音.disable = true;

        c.刷新设备列表();
        systemFc.EventsOn("playStatus", function (data) {
            let jsondata = JSON.parse(data);
            console.log("playStatus", jsondata)
            // {"Status":"STOPPED","UUID":"9c39e130-6ca7-4bc4-b199-60c5acc261d4","event":"playStatus"}
            // Status PLAYING PAUSED_PLAYBACK STOPPED
            let status = {
                "PLAYING" :"暂停播放",
                "PAUSED_PLAYBACK" :"继续播放",
                "STOPPED" :"开始播放"
            }
            comps.Button_开始播放.text = status[jsondata.Status];

        });
        systemFc.EventsOn("playPosition", function (data) {
            let jsondata = JSON.parse(data);
            console.log("playPosition", jsondata)

            // {"currentPosition":"3673","currentTime":"1:59:09","event":"playPosition","overallLength":"7149","totalEvent":"1:01:13"}
            comps.Label_开始时间.text = jsondata.currentTime;
            comps.Label_结束时间.text = jsondata.totalEvent;
            comps.进度条_时间轴.max = parseInt(jsondata.overallLength);
            comps.进度条_时间轴.n = parseInt(jsondata.currentPosition);

        });
    }
    c.刷新设备列表 = async function () {
        console.log("刷新设备列表")
        let 设备列表Json = await goFc.E获取设备列表();
        console.log(设备列表Json)
        // json文本转换为对象
        let 设备列表 = [];
        try{
            设备列表 = JSON.parse(设备列表Json);
        }catch (e) {
            ElMessage.error('你的局域网内没有发现设备,多刷新几次看看');
            return
        }
        console.log(设备列表)
        // [{"Model":"华为智慧屏 S65","URL":"http://192.168.100.204:25826/description.xml"},{"Model":"奇异果极速投屏-华为(204)","URL":"http://192.168.100.204:39620/description.xml"},{"Model":"MacBook Pro","URL":"http://192.168.10scription.xml"}]
        // 转换为 Model转换为label url转换为value
        设备列表 = 设备列表.map((item) => {
            return {
                label: item.Model,
                value: item.URL,
            }
        })

        // comps.Select_设备列表.options = [
        //     {"label": "华为智慧屏1", "value": "华为智慧屏1"},
        //     {"label": "华为智慧屏2", "value": "华为智慧屏2"},
        // ]
        comps.Select_设备列表.options = 设备列表;


        // comps.Select_设备列表.value = "华为智慧屏1"

        // 用 ElMessage 弹出提示 刷新成功
        ElMessage.success('设备列表刷新成功');

        if (设备列表.length >= 1) {
            comps.Select_设备列表.value = 设备列表[0].value
        }

    }

    c.Button_刷新Click = function () {
        console.log("Button_刷新Click")
        c.刷新设备列表();

    }

    c.Button_选择文件Click = async function () {
        console.log("Button_选择文件Click")
        let 文件路径 = await goFc.OpenFileDialog();
        console.log(文件路径)
        comps.TextEdit_URL.text = 文件路径;

    }

    c.Button_检查更新Click = function () {
        console.log("Button_检查更新Click")

        ElMessage.success('暂时没有开发');
    }

    c.Button_开始播放Click = function () {
        console.log("Button_开始播放Click")
        if (comps.Select_设备列表.value == "") {
            ElMessage.error('请选择设备');
            return
        }
        if (comps.TextEdit_URL.text == "") {
            ElMessage.error('请选择视频文件');
            return
        }

        if(comps.Button_开始播放.text == "开始播放"){
            //字幕文件路径设置为当前的路径 comps.TextEdit_URL.text 修改后缀为 srt
            let 字幕文件路径 = comps.TextEdit_URL.text.replace(".mp4", ".srt");
            goFc.E投递视频文件(comps.Select_设备列表.value, comps.TextEdit_URL.text, 字幕文件路径);
            ElMessage.success('已发送播放指令');
        }
        if(comps.Button_开始播放.text == "暂停播放"){
            goFc.E暂停播放()
            ElMessage.success('已发送暂停播放指令');
        }
        if(comps.Button_开始播放.text == "继续播放"){
            goFc.E继续播放()
            ElMessage.success('已发送继续播放指令');
        }

        comps.Button_停止播放.disable = false;
        comps.按钮_音量加.disable = false;
        comps.按钮_音量减.disable = false;
        comps.按钮_静音.disable = false;
    }

    c.Button_停止播放Click = function () {
        console.log("Button_停止播放Click")
        goFc.E停止播放()
        ElMessage.success('已发送停止播放指令');
        comps.Button_开始播放.text = "开始播放"
        comps.Label_开始时间.text = "00:00:00";
        comps.Label_结束时间.text = "00:00:00";
        comps.进度条_时间轴.max = parseInt(60);
        comps.进度条_时间轴.n = parseInt(0);
    }
    
    c.Button_检查更新被单击 = function () {
        console.log("Button_检查更新被单击")
        ElMessage.success('等待开发');

    }
    

    c.进度条_时间轴鼠标左键被放开 = async function () {
        console.log("进度条_时间轴鼠标左键被放开")
        let msg = await goFc.E设置播放进度(comps.进度条_时间轴.n);
        ElMessage.success(msg);

    }

    c.按钮_静音被单击 = async function () {
        console.log("按钮_静音被单击")
        if (comps.按钮_静音.text == "取消静音") {
            let msg = await  goFc.E取消静音()
            ElMessage.success(msg);
            comps.按钮_静音.text = "静音"
        } else {
            let msg = await goFc.E静音()
            ElMessage.success(msg);

            comps.按钮_静音.text = "取消静音"
        }

    }

    c.按钮_音量加被单击 = async function () {
        console.log("按钮_音量加被单击")
        let msg = await goFc.E音量('+')
        ElMessage.success(msg);

    }

    c.按钮_音量减被单击 = async function () {
        console.log("按钮_音量减被单击")
        let msg = await goFc.E音量('-')
        ElMessage.success(msg);

    }
//Don't delete the event function flag
}