import { WEBSOCKET_URL } from './env';
import laugh from "./laugh";
import { SoccerMessage, SoccerMessageHandler } from './msg';

const reg = laugh();

//  output
const ta : HTMLTextAreaElement = <HTMLTextAreaElement>document.getElementById("t");
const logTempl : HTMLTemplateElement = <HTMLTemplateElement>document.getElementById('log');
const logPoint : HTMLDivElement = <HTMLDivElement>document.getElementById("logs");

document.getElementById('f').addEventListener("submit", ev => {
    ev.preventDefault();
    const msg = new SoccerMessage("loose text", ta.value);
    sendAndLog(ws, "send", msg);
});

const sendAndLog = (ws : WebSocket, subject : string, msg : SoccerMessage) => {
    ws.send(msg.serialize());
    logit(subject, msg);
}

const logit = (subject : string, msg : SoccerMessage) => {
    const pretty = JSON.stringify(msg.record,null,"\t");
    const html = `<div><h4>${subject}</h4><pre>${pretty}</pre></div>`;
    const parser = new DOMParser();
    const doc =  parser.parseFromString(html, 'text/html');
    const element = doc.body.firstChild;
    logPoint.appendChild(element);
}

type coords = {
    x : number;
    y : number;
}

const randomCoordinates = () : coords => {
    const x = Math.floor(Math.random()*20) - 10;
    const y = Math.floor(Math.random()*20) - 10;
    return {x,y};
}

const handleMessage : SoccerMessageHandler = (msg : SoccerMessage) => {

    //  handle it
    switch (msg.subject) {
        case "hello":
        case "goodbye":
            console.log("soccer mesage", "hello / goodbye", msg.record);
        break;
        case "marco":
        case "polo":
            if (msg.record.payload < 11) {
                let retort = msg.reply(msg.record.payload+1);
                if (msg.subject === "marco") {
                    retort.record.subject = "polo";
                } else {
                    retort.record.subject = "marco";
                }
                sendAndLog(ws, "send", retort);
            }
            console.log("soccer mesage", "marco polo", msg.record);
        break;
        case "command/addPeer":
            console.log("add this peer", msg.record);
            const attrs = msg.record.payload;
            const coords = randomCoordinates();
            attrs.id = attrs.nick;
            attrs.x = coords.x;
            attrs.y = coords.y;
            attrs.size = 9;
            attrs.color = "green";
            const xx = reg.addNode(attrs);
            console.log({attrs, xx});
        break;
        default:
            console.log("soccer mesage", "unhandled subject", msg.record);
    }
}

const ws = new WebSocket(WEBSOCKET_URL);

ws.addEventListener("message", ev => {
    const msg = SoccerMessage.deserialize(ev.data);
    logit("receive", msg);
    handleMessage(msg);
});
ws.addEventListener("open", console.info);
ws.addEventListener("error", console.error);
ws.addEventListener("close", console.debug);

const msg = new SoccerMessage("hello", Math.floor(Math.random()*10000));
sendAndLog(ws, "hello", msg);
