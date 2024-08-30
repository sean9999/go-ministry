import { SoccerMessage } from './msg';

const ws = new WebSocket("ws://localhost:8282/ws");

ws.addEventListener("message", ev => {
    const msg = SoccerMessage.deserialize(ev.data);
    console.log("msg", msg.record);
});
ws.addEventListener("open", console.info);
ws.addEventListener("error", console.error);
ws.addEventListener("close", console.debug);

setInterval(() => {
    const msg = new SoccerMessage("hello", Math.floor(Math.random()*10000));
    ws.send( msg.serialize() );
}, 30511);


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

// setInterval(() => {
//     const msg = new SoccerMessage("hello", Math.floor(Math.random()*10000));
//     console.log("hello", msg.record);
//     ws.send(msg);
// }, 7511);
