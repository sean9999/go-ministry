import { Registry } from "./laugh";

const waitFor = (ms : number) => {
    return new Promise(resolve => {
        window.setTimeout(resolve, ms);
    });
};

const randomColour = () : string => {
    let r = "#";
    for (let i=0;i<3;i++) {
        let letter = Math.floor(Math.random()*16).toString(16); 
        r += `${letter}${letter}`;
    }
    return r
}

const sendMessage = async (reg : Registry, fromId : string, toId : string) => {
    //const e  = graph.edge(fromId, toId);
    
    const e = reg.getLink(fromId, toId);
    const graph = reg.graph;
    if (!e) {
        console.log({fromId, toId, reg});
        return Promise.reject(e);
    }
    graph.setNodeAttribute(fromId, "color", "purple");
    graph.setNodeAttribute(toId, "color", "purple");
    graph.setEdgeAttribute(e, "size", 7);
    await waitFor(100);
    for (let i=0;i<10;i++) {
        await waitFor(100);
        graph.setEdgeAttribute(e, "color", randomColour());
    }
    graph.setEdgeAttribute(e, "color", "gray");
    graph.setEdgeAttribute(e, "size", 2);
    return Promise.resolve(e);
}

export { sendMessage };
