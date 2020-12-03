const LukaText = 1;
const LukaImg = 2;
const LukaSingle = 1;
const LukaGroup = 2;

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
            content: utf8_to_b64(msg),
            msgType: LukaSingle,
            msgContentType: LukaText
        })
}



class IpcMsg {

    constructor() {
        this.TypeNothing = -1;
        this.TypeErr = 0;
        this.TypeLogin = 1;
        this.TypeMessage = 2;
        this.TypeLoginFinished = 3;
        this.TypeMessageRequired = 6;
    }
    unifiedIpcMsg(TypeName, ContextJson) {
        return JSON.stringify({
            Type: TypeName,
            ContextByte: utf8_to_b64(JSON.stringify(ContextJson))
        })
    }
    LoginMsg(name, password) {
        return this.unifiedIpcMsg(this.TypeLogin, {
            Name: name,
            Password: password
        })
    }
    IMMsg(from, target, msg) {
        return this.unifiedIpcMsg(this.TypeMessage, {
            from: from,
            target: target,
            content: utf8_to_b64(msg),
            msgType: LukaSingle,
            msgContentType: LukaText
        })
    }
    MessageRequiredMsg(target, msgType) {
        return this.unifiedIpcMsg(this.TypeMessageRequired, {
            From: target,
            MsgType: msgType
        })
    }
    ShowErr(msg) {
        console.log(msg)
    }
}