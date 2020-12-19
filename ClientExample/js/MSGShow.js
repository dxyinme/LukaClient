// MSG Show
class MSGShow {

    constructor(){

    }

    MsgBox(from, target, msg, isSelf) {
        let temp = '<div class="{{boxClass}}">\n' +
        '{{content}}\n' +
        '<button class="{{iconClass}}">\n' +
        '<img class="chat-round-icon" src="{{icon-src}}"/>\n' +
        '</button>\n' +
        '</div>';
        let boxClass, iconClass, chatRoundIcon;
        chatRoundIcon = "chat-round-icon";
        if(isSelf === true) {
            boxClass = "message-self";
            iconClass = chatRoundIcon + " chat-round-icon-self";
        } else {
            boxClass = "message-opposite";
            iconClass = chatRoundIcon + " chat-round-icon-opposite";
        }
        temp =  temp.replace('{{boxClass}}', boxClass)
                    .replace('{{iconClass}}', iconClass)
                    .replace('{{icon-src}}',
                "https://avatars1.githubusercontent.com/u/32793868?s=400&u=2e2f7b7637470fff15afe82927532efb1a1fde5a&v=4")
                    .replace('{{content}}', msg);
        console.log(temp);
        return temp;
    }
}