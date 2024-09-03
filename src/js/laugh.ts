import { EdgeCurvedArrowProgram } from "@sigma/edge-curve";
import chroma from "chroma-js";
import Graph from "graphology";
import ForceSupervisor from "graphology-layout-force/worker";
import { Attributes } from "graphology-types";
import Sigma from "sigma";
import { uuidv7 } from "uuidv7";
import { sendMessage } from "./viz";

export type Link = {
	From : string;
	To : string;
}

export type EdgeMap = Map<Link, string>

export class Registry{
	graph : Graph
	edges : EdgeMap
	nodes : Set<string>
	constructor(g : Graph){
		this.graph = g;
		this.edges = new Map();
		this.nodes = new Set();
	}
	getLink(From : string, To : string) : string {
		const lnk : Link = {
			From,
			To
		};
		return this.edges.get(lnk);
	}
	addNode(attrs? : Attributes) : string {
		const label = uuidv7(); 
		this.graph.addNode(label, attrs);
		this.nodes.add(label);
		return label;
	}
	addEdge(From : string, To : string, attrs? : Attributes) : string {
		const lnk : Link = {
			From,
			To
		};
		const id = this.graph.addEdge(From, To);
		this.edges.set(lnk, id);
		return id;
	}
}

export default () => {

	let edgeThickness = 0;
	const grow = (ts : number) => {
		edgeThickness++;
		graph.forEachEdge(e => {
			graph.setEdgeAttribute(e, "size", edgeThickness);
		});
	};

	const changeThings = () => {
		console.log("change things");
		graph.forEachNode(nd => {
			graph.setNodeAttribute(nd, "size", 19);
		});
		graph.forEachInEdge(edg => {
			//graph.setEdgeAttribute(edg, "type", "curvedArrow");
			requestAnimationFrame(grow);
		});
	};
	document.getElementById("ct").addEventListener("click", changeThings);


	// Retrieve the html document for sigma container
	const container = document.getElementById("sigma-container") as HTMLElement;

	// Create a sample graph
	const graph = new Graph();
	const registry = new Registry(graph);
	const n1 = registry.addNode({ x: 0, y: 0, size: 10, color: chroma.random().hex(), "label": "bob" });
	const n2 = registry.addNode({ x: -5, y: 5, size: 10, color: chroma.random().hex() });
	const n3 = registry.addNode({ x: 5, y: 5, size: 10, color: chroma.random().hex() });
	const n4 = registry.addNode({ x: 0, y: 10, size: 10, color: chroma.random().hex() });
	registry.addEdge(n1, n2, {"color": "#00FF00", "type": "line", "label": "sleeps with"});
	registry.addEdge(n2, n4);
	registry.addEdge(n4, n3, {"color": "#9900FF"});
	registry.addEdge(n3, n1);

	// Create the spring layout and start it
	const layout = new ForceSupervisor(graph, { isNodeFixed: (_, attr) => attr.highlighted });
	layout.start();

	// Create the sigma
	const renderer = new Sigma(graph, container,{
		defaultEdgeType: "curve",
		labelSize: 15,
		edgeProgramClasses: {
			curve: EdgeCurvedArrowProgram,
		},
	});

	//
	// Drag'n'drop feature
	// ~~~~~~~~~~~~~~~~~~~
	//

	// State for drag'n'drop
	let draggedNode: string | null = null;
	let isDragging = false;

	// On mouse down on a node
	//  - we enable the drag mode
	//  - save in the dragged node in the state
	//  - highlight the node
	//  - disable the camera so its state is not updated
	renderer.on("downNode", (e) => {
		isDragging = true;
		draggedNode = e.node;
		graph.setNodeAttribute(draggedNode, "highlighted", true);
	});

	// On mouse move, if the drag mode is enabled, we change the position of the draggedNode
	renderer.getMouseCaptor().on("mousemovebody", (e) => {
		if (!isDragging || !draggedNode) return;

		// Get new position of node
		const pos = renderer.viewportToGraph(e);

		graph.setNodeAttribute(draggedNode, "x", pos.x);
		graph.setNodeAttribute(draggedNode, "y", pos.y);

		// Prevent sigma to move camera:
		e.preventSigmaDefault();
		e.original.preventDefault();
		e.original.stopPropagation();
	});

	// On mouse up, we reset the autoscale and the dragging mode
	renderer.getMouseCaptor().on("mouseup", () => {
		if (draggedNode) {
			graph.removeNodeAttribute(draggedNode, "highlighted");
		}
		isDragging = false;
		draggedNode = null;
	});

	// Disable the autoscale at the first down interaction
	renderer.getMouseCaptor().on("mousedown", () => {
		if (!renderer.getCustomBBox()) renderer.setCustomBBox(renderer.getBBox());
	});

	//
	// Create node (and edge) by click
	// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
	//

	// When clicking on the stage, we add a new node and connect it to the closest node
	renderer.on("clickStage", ({ event }: { event: { x: number; y: number } }) => {
		// Sigma (ie. graph) and screen (viewport) coordinates are not the same.
		// So we need to translate the screen x & y coordinates to the graph one by calling the sigma helper `viewportToGraph`
		const coordForGraph = renderer.viewportToGraph({ x: event.x, y: event.y });

		// We create a new node
		const node = {
			...coordForGraph,
			size: 10,
			color: chroma.random().hex(),
		};

		// Searching the two closest nodes to auto-create an edge to it
		const closestNodes = graph
			.nodes()
			.map((nodeId) => {
				const attrs = graph.getNodeAttributes(nodeId);
				const distance = Math.pow(node.x - attrs.x, 2) + Math.pow(node.y - attrs.y, 2);
				return { nodeId, distance };
			})
			.sort((a, b) => a.distance - b.distance)
			.slice(0, 2);

		// We register the new node into graphology instance
		//const id = uuidv7();
		const id = registry.addNode(node);

		// We create the edges
		closestNodes.forEach(ee => {
			registry.addEdge(id, ee.nodeId)
		});

		//	animate
		sendMessage(registry, id, closestNodes[0].nodeId).then(() => {
			return sendMessage(registry, id, closestNodes[1].nodeId);
		})
		.then(console.log)
		.catch(console.error);

	});

};

