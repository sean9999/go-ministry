import { SoccerMessage, SoccerMessageHandler } from './msg';

//  Soccer is a simple wrapper around WebSocket
class Soccer {
    url : string;
    ws : WebSocket;
    messageHandler : SoccerMessageHandler
    connectionAttempts : number;
    constructor(url : string){
        this.url = url;
        this.init();
        this.connectionAttempts = 0;
    }
    init(){
        this.connectionAttempts++;
        this.ws = new WebSocket( this.url );
        this.ws.addEventListener("message", (ev : MessageEvent) => {
            const msg : SoccerMessage = JSON.parse(ev.data);
            if (this.messageHandler !== undefined) {
                this.messageHandler(msg);
            }
        });
        this.ws.addEventListener("close", console.info.bind(this.ws));
        this.ws.addEventListener("error", console.error.bind(this.ws));
        return this.ws.readyState;
    }
    send(msg : SoccerMessage) {
        this.ws.send( msg.serialize() );
    }
    retry(msg : SoccerMessage) {
        this.reconnect();
        this.send(msg);
    }
    reconnect() {
        console.info("reconnecting");
        return this.init();
    }
    onMessage(fn : SoccerMessageHandler) {
        this.messageHandler = fn;
    }
}

export default Soccer;




// ws.addEventListener("message", (ev : MessageEvent) => {
//     const msg = <Message>ev.data;
//     console.log("receive",{msg});
// });

// ws.addEventListener("open", (ev) => {
//     console.log('connected',{ev});
// });

// ws.addEventListener("close", (evt) => {
//     console.log({evt}, "close");
// });

// ws.addEventListener("error", (err) => {
//     console.log({err}, "error");
// });

// setInterval(() => {
//     const msg = {
//         "hello": 7,
//         "why": false
//     };
//     console.log("sending", msg);
//     ws.send(JSON.stringify(msg));
// }, 5555);
