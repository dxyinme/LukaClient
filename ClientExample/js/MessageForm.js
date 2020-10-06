const LukaText = 1
const LukaImg = 2
const LukaGroup = 1
const LukaSingle = 2


function utf8_to_b64( str ) {
    return window.btoa(unescape(encodeURIComponent( str )));
}

function b64_to_utf8( str ) {
    return decodeURIComponent(escape(window.atob( str )));
}


function encodeLukaMsg(from,target,msg){
    return JSON.stringify(
        {
            from: from,
            target: target,
            // 转换成base64
            content: utf8_to_b64(msg)
        })
}