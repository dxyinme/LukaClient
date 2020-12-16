// use jquery.
class MSGShow {

    constructor(){

    }

    MsgBox(from, target, msg) {
        return $("<div/>", {
            "class": "chat-msg-box"
        }).append(
            $("<a/>", {
                "text": msg
            })
        );
    }
}