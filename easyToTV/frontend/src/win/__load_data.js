import {ref} from 'vue'
import {defineStore} from 'pinia'
import {WindowSetSize, WindowSetTitle} from "../../wailsjs/runtime/runtime"; // 根据实际文件路径进行修改
import {BindWindowEvent} from '@/win/event'

export const __load_data = defineStore('window_data', {
    state: () => {
        let data = {}
        data.list = ref([])
        data.comps = {}
        return data
    },
    actions: {
        init() {

            BindWindowEvent(this, this.comps)
            try {
                if (this.comps.Win.hasOwnProperty("event_created")) {
                    this.WinCreated()
                }
            } catch (e) {
                console.log("WinCreated To be defined")
            }

            const dthis = this

            try {
                if (dthis.comps.Win.width.includes('v') || dthis.comps.Win.width.includes('%')) {
                    return;
                }
                if (dthis.comps.Win.height.includes('v') || dthis.comps.Win.height.includes('%')) {
                    return;
                }
                WindowSetSize(parseInt(dthis.comps.Win.width), parseInt(dthis.comps.Win.height))
                //Recalculate the width and height of the client area
                setTimeout(function () {
                    var WidthFix = parseInt(dthis.comps.Win.width) - window.innerWidth
                    var HeightFix = parseInt(dthis.comps.Win.height) - window.innerHeight
                    WindowSetSize(parseInt(dthis.comps.Win.width) + WidthFix, parseInt(dthis.comps.Win.height) + HeightFix)
                    document.body.style.overflow = 'auto'
                }, 1)
                WindowSetTitle(dthis.comps.Win.text)
            } catch (e) {
                console.error("Error initializing win size", e)
            }
        },
        handleAllEvents(el, e, item, callFuncName) {
            try {
                var dynamicFunction = undefined
                eval(`dynamicFunction = this.${callFuncName}`)
                dynamicFunction(e, item)
            } catch (e) {
                console.log("Function call error", callFuncName, "dynamicFunction", dynamicFunction, "e", e)
            }
        },

    },
})