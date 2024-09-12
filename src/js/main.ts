import { WEBSOCKET_URL } from './env';
import init, { Link } from "./graph";
import { SoccerMessage, SoccerMessageHandler } from './msg';
import { emitParticle, pulseEdge, pulseNode } from './viz';
const ta : HTMLTextAreaElement = <HTMLTextAreaElement>document.getElementById("t");
//const logTempl : HTMLTemplateElement = <HTMLTemplateElement>document.getElementById('log');
const logPoint : HTMLDivElement = <HTMLDivElement>document.getElementById("logs");
const graphContainer : HTMLDivElement = <HTMLDivElement>document.getElementById('sigma-container');
const btnAddNode : HTMLButtonElement = <HTMLButtonElement>document.getElementById('aan');

const f1 : HTMLSelectElement = <HTMLSelectElement>document.getElementById('f1');
const f2 : HTMLSelectElement = <HTMLSelectElement>document.getElementById('f2');

const ws = new WebSocket(WEBSOCKET_URL);

const {registry} = init(graphContainer, ws);

btnAddNode.addEventListener("click", ev => {
    const msg = new SoccerMessage("please/addNode");
    console.log(ev,msg);
    ws.send(msg.serialize());
});

const rebuildFriends = () => {
    for (const child of f1.children) {
        f1.removeChild(child)
    }
    for (const child of f2.children) {
        f2.removeChild(child)
    }
    registry.graph.nodes().forEach(nodeId => {
        let opt = document.createElement("option")
        opt.value = nodeId;
        opt.innerText = nodeId;
        let opt2 = opt.cloneNode(true);
        f1.append(opt);
        f2.append(opt2);
    });
};

document.getElementById('makefriends').addEventListener("click", ev => {
    ev.preventDefault();
    const msg = new SoccerMessage("please/addRelationship", [f1.value, f2.value]);
    ws.send(msg.serialize());
});

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

    console.log("receiving soccer message", msg.record);

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
        case "command/addNode":
            const attrs = msg.record.payload;
            const coords = randomCoordinates();
            attrs.id = attrs.nick;
            attrs.x = coords.x;
            attrs.y = coords.y;
            attrs.size = 9;
            attrs.color = "green";
            registry.addNode(attrs);
            rebuildFriends();
        break;
        case "command/addEdge":
        case "command/addRelationship":
            const [from, to] = msg.record.payload;
            registry.addEdge(from, to);
        break;
        case "command/removeRelationship":
            console.info("remove link", msg.record);
            let lnk = Link(msg.record.payload[0], msg.record.payload[1]);
            registry.removeLink(lnk);
        break;
        case "command/pulseNode":
            pulseNode(registry.graph, msg.record.payload);
        break;
        case "command/pulseEdge":
            pulseEdge(registry, msg.record.payload[0], msg.record.payload[1])
        break;
        case "command/updateNode":
            // sendMessage(registry, msg.record.from, msg.record.to).then(() => {
            //     return registry.updateNode(msg.record.to, msg.record.payload);
            // }).then(console.log);
            registry.updateNode(msg.record.to, msg.record.payload);

        case "command/passItOn":
            emitParticle(registry, msg.record.from, msg.record.to)
            .then(() => {
                return registry.updateNode(msg.record.to, msg.record.payload);
            })
            .then(console.log)
            .then(console.error);
        break;

        default:
            console.log("soccer mesage", "unhandled subject", msg.record);
    }
}


ws.addEventListener("message", ev => {
    const msg = SoccerMessage.deserialize(ev.data);
    logit("receive", msg);
    handleMessage(msg);
});
ws.addEventListener("open", console.info);
ws.addEventListener("error", console.error);
ws.addEventListener("close", console.debug);

const awaken = () => {
    const msg = new SoccerMessage("hello/imAwake");
    sendAndLog(ws, "send", msg);
};

setTimeout(awaken, 997);