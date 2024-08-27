const ws = new WebSocket("ws://localhost:8282/ws");

ws.addEventListener("message", (msg) => {
    console.log("receive",{msg});
});

ws.addEventListener("open", (ev) => {
    console.log('connected',{ev});
});

ws.addEventListener("close", (evt) => {
    console.log({evt}, "close");
});

ws.addEventListener("error", (err) => {
    console.log({err}, "error");
});

setInterval(() => {
    const msg = {
        "hello": 7,
        "why": false
    };
    console.log("sending", msg);
    ws.send(JSON.stringify(msg));
}, 5555);
