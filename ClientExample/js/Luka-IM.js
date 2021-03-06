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
        this.TypeNewWindow = 7;
        this.TypeGroupOperator = 8;
        this.TypeSecret = 9;
        this.TypeFile = 10;

        this.TypeWindowGroupWindow = "group";


        this.JoinGroup = "JOIN GROUP";
        this.CreateGroup = "CREATE GROUP";
        this.DeleteGroup = "DELETE GROUP";
        this.LeaveGroup  = "LEAVE GROUP";

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

    IMFile(from, target, filePath) {
        return this.unifiedIpcMsg(this.TypeFile, {
            From: from,
            Target: target,
            Path: filePath
        })
    };

    IMMsg(from, target, msg, msgType) {
        let msgTypeNow , targetNow = "", groupNameNow = "";
        if(msgType === "single") {
            msgTypeNow = LukaSingle;
            targetNow = target;
        } else {
            msgTypeNow = LukaGroup;
            groupNameNow = target;
        }
        return this.unifiedIpcMsg(this.TypeMessage, {
            from: from,
            target: targetNow,
            groupName: groupNameNow,
            content: utf8_to_b64(msg),
            msgType: msgTypeNow,
            msgContentType: LukaText
        })
    }
    IMMsgS(from, target, msg, msgType) {
        let msgTypeNow , targetNow = "", groupNameNow = "";
        if(msgType === "single") {
            msgTypeNow = LukaSingle;
            targetNow = target;
        } else {
            msgTypeNow = LukaGroup;
            groupNameNow = target;
        }
        return this.unifiedIpcMsg(this.TypeMessage, {
            from: from,
            target: targetNow,
            groupName: groupNameNow,
            content: utf8_to_b64(msg),
            msgType: msgTypeNow,
            msgContentType: LukaText,
            secretLevel: 1
        })
    }
    MessageRequiredMsg(target, msgType) {
        return this.unifiedIpcMsg(this.TypeMessageRequired, {
            From: target,
            MsgType: msgType
        })
    }

    OpenNewWindowMsg(windowType) {
        return this.unifiedIpcMsg(this.TypeNewWindow, {
            WindowType: windowType
        });
    }

    GroupOperator(uid, groupName, groupOp) {
        return this.unifiedIpcMsg(this.TypeGroupOperator, {
            GroupOp: groupOp,
            GroupName: groupName,
            Uid: uid
        });
    }

    ShowErr(msg) {
        console.log(msg)
    }

    randomString(len) {
        let chars = 'ABCDEFGHJKLMNOPQRSTUVWXYZabcdefhijklmnoprstuvwxyz12345678';
        let maxPos = chars.length;
        let i, randStr = '';
        for (i = 0; i < len; i++) {
            randStr += chars.charAt(Math.floor(Math.random() * maxPos));
        }
        return randStr;
    }

    GoSecret(from, target) {
        return this.unifiedIpcMsg(this.TypeSecret, {
            From: from,
            Target: target,
            AESKey: utf8_to_b64(this.randomString(16)),
        })
    }
}