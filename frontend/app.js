import SceneRunner from "./canvas/SceneRunner.js";
import { WalkerScene } from "./canvas/scene/scenes/walker/WalkerScene.js";
import { config } from "./canvas/scene/scenes/walker/config.js";

const wsUri = "ws://localhost:8000/connect";

const wsConfig = {
  autoConnect: true,
};

let ws;
let board;

let scene = new WalkerScene(config);

const el2 = (id) => document.getElementById(id);

let connectStatus = false;

const EL = {
  disconnectBtn: null,
  wsInput: null,
  sendBtn: null,
  isConnected: null,
  connectBtn: null,
};

getElements();
registerEventHandlers();
renderConnectionStatus();

board = new SceneRunner(scene);

board.init();

board.playerJoin({});

if (wsConfig.autoConnect) {
  connectToSocket();
}

function setIsConnected(connected) {
  connectStatus = connected;
  renderConnectionStatus();
}

function renderConnectionStatus() {
  EL.isConnected.innerText = connectStatus ? "✅" : "❌";
}

function getElements() {
  EL.disconnectBtn = el2("disconnectBtn");
  EL.wsInput = el2("wsInput");
  EL.sendBtn = el2("sendBtn");
  EL.isConnected = el2("connectStatusP");
  EL.connectBtn = el2("connectBtn");
}

function registerEventHandlers() {
  EL.disconnectBtn.addEventListener("click", () => {
    ws.close();
  });

  EL.sendBtn.addEventListener("click", () => {
    if (connectStatus) {
      ws.send(wsInput.value);
      wsInput.value = "";
    }
  });

  EL.connectBtn.addEventListener("click", () => {});
}

function connectToSocket() {
  ws = new WebSocket(wsUri);
  ws.onopen = function (event) {
    console.log("Socket connection established: ", event);
    setIsConnected(true);

    ws.send("Hello World");
  };

  ws.onmessage = function (event) {
    console.log("Message: ", event.data);
  };

  ws.onclose = function (event) {
    console.log("Disconnected", event);
    setIsConnected(false);
  };
}
