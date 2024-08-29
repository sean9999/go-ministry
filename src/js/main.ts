import { SoccerMessage } from './msg';
import Soccer from './soccer';

const ws = new Soccer("ws://localhost:8282/ws");

ws.onMessage(msg => {
    console.log(msg);
});

ws.init();

// const hello = new SoccerMessage("hello");

// console.log({hello});

// ws.send(hello);

// ws.addEventListener("message", (ev : MessageEvent) => {
//     const msg = <SoccerMessage>ev.data;
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

setInterval(() => {
    const msg = new SoccerMessage("hello", Math.floor(Math.random()*10000));
    console.log("sending", {msg});
    ws.send(msg);
}, 5555);
