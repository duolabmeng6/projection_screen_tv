import { ElMessage, ElMessageBox } from "element-plus";
import * as systemFc from "../../wailsjs/runtime";
import * as goFc from "../../wailsjs/go/main/App";

export function BindWindowEvent(c, comps) {
    c.WinCreated = async function () {
        console.log("WinCreated")
        comps.Win.text = "多多投屏" + await goFc.GetVersion();
        c.刷新设备列表();
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
        //字幕文件路径设置为当前的路径 comps.TextEdit_URL.text 修改后缀为 srt 
        let 字幕文件路径 = comps.TextEdit_URL.text.replace(".mp4", ".srt");

        goFc.E投递视频文件(comps.Select_设备列表.value, comps.TextEdit_URL.text, 字幕文件路径);
        ElMessage.success('已发送播放指令');

    }

    c.Button_停止播放Click = function () {
        console.log("Button_停止播放Click")
        goFc.E停止播放()
        ElMessage.success('已发送停止播放指令');

    }
    
    c.Button_检查更新被单击 = function () {
        console.log("Button_检查更新被单击")
    }
//Don't delete the event function flag
}