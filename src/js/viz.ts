import Graph from "graphology";
import { Registry } from "./graph";

const TICK = 250;

const waitFor = (ms : number) => {
    return new Promise(resolve => {
        window.setTimeout(resolve, ms);
    });
};

const pulseNode = async (g : Graph, id : string) => {
    const originalColour = g.getNodeAttribute(id, "_originalColor");
    const originalSize = g.getNodeAttribute(id, "_originalSize");
    g.setNodeAttribute(id, "size", 17);
    g.setNodeAttribute(id, "color", "purple");
    await waitFor(TICK);
    g.setNodeAttribute(id, "size", originalSize);
    g.setNodeAttribute(id, "color", originalColour);
    return Promise.resolve(true);
}

const pulseEdge = async (reg : Registry, fromId : string, toId : string) => {
    const e = reg.getLink(fromId, toId);
    const graph = reg.graph;
    const originalColor = graph.getEdgeAttribute(e, "color");
    const originalSize = graph.getEdgeAttribute(e, "size");
    graph.setEdgeAttribute(e, "size", 7);
    graph.setEdgeAttribute(e, "color", "purple");
    await waitFor(TICK);
    graph.setEdgeAttribute(e, "size", originalSize);
    graph.setEdgeAttribute(e, "color", originalColor);
    return Promise.resolve(true);    
}

const emitParticle = async (reg : Registry, fromId : string, toId : string) => {
    const e = reg.getLink(fromId, toId);
    const graph = reg.graph;
    if (!e) {
        console.warn({fromId, toId, reg});
        return Promise.reject(e);
    }
    await pulseNode(graph, fromId);
    await pulseEdge(reg, fromId, toId);
    await pulseNode(graph, toId);
    return Promise.resolve(e);
}

export { pulseEdge, pulseNode, emitParticle as sendMessage };

