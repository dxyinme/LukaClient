const LukaText = 1
const LukaImg = 2
const LukaGroup = 1
const LukaSingle = 2


function encodeLukaMsg(from,target,msg){
    return JSON.stringify(
        {
            from: from,
            target: target,
            // 转换成base64
            content: window.btoa(msg)
        })
}