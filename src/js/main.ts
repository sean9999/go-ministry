import { WEBSOCKET_URL } from './env';
import { SoccerMessage, SoccerMessageHandler } from './msg';

const handleMessage : SoccerMessageHandler = (msg : SoccerMessage) => {
    switch (msg.subject) {
        case "hello":
        case "goodbye":
            console.log("soccer mesage", "hello / goodbye", msg.record);
        break;
        case "marco":
        case "polo":
            console.log("soccer mesage", "marco polo", msg.record);
        break;
        default:
            console.log("soccer mesage", "unhandled subject", msg.record);
    }
}

const ws = new WebSocket(WEBSOCKET_URL);

ws.addEventListener("message", ev => {
    const msg = SoccerMessage.deserialize(ev.data);
    handleMessage(msg);
});
ws.addEventListener("open", console.info);
ws.addEventListener("error", console.error);
ws.addEventListener("close", console.debug);

setInterval(() => {
    const msg = new SoccerMessage("hello", Math.floor(Math.random()*10000));
    ws.send( msg.serialize() );
}, 30511);
