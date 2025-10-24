import SceneRunner from "./canvas/SceneRunner.js";
import SocketManager from "./canvas/handlers/SockerManager.js";
import { WalkerScene } from "./canvas/scene/scenes/walker/WalkerScene.js";
import { config } from "./canvas/scene/scenes/walker/config.js";
import { EventEmitter } from "./canvas/classes/EventEmitter/EventEmitter.js";
import { uuid } from "./canvas/util/uuid.js";

const ownerId = uuid();
const gameEventBus = new EventEmitter();
const scene = new WalkerScene(config, gameEventBus);
const sceneRunner = new SceneRunner(gameEventBus, scene, ownerId);
const socketManager = new SocketManager(
  sceneRunner,
  "ws://localhost:8000/connect"
);

const el2 = (id) => document.getElementById(id);

const el = {
  disconnectBtn: null,
  wsInput: null,
  sendBtn: null,
  isConnected: null,
  connectBtn: null,
  userId: null,
};

getElements();
registerDomEventHandlers();
registerGameEventHandlers();
renderConnectionStatus();

sceneRunner.init();

// GAME EVENTS
function handleSceneKeyPress(event) {
  switch (event.data.key) {
    case "w":
      socketManager.send({
        id: ownerId,
        msg: {
          type: "input",
          up: true,
        },
      });
      break;
    case "d":
      socketManager.send({
        id: ownerId,
        msg: {
          type: "input",
          right: true,
        },
      });
      break;
    case "s":
      socketManager.send({
        id: ownerId,
        msg: {
          type: "input",
          down: true,
        },
      });
      break;
    case "a":
      socketManager.send({
        id: ownerId,
        msg: {
          type: "input",
          left: true,
        },
      });
      break;
    default:
      break;
  }
}

// GAME EVENTS
function registerGameEventHandlers() {
  gameEventBus.subscribe("sceneKeyPressed", handleSceneKeyPress);
}

// TODO: pass is connected as an app level event
function connectToServerHandler() {
  socketManager.connectToServer((isConnected) =>
    renderConnectionStatus(isConnected)
  );
}

// DOM
function registerDomEventHandlers() {
  el.disconnectBtn.addEventListener("click", () => {
    socketManager.close(renderConnectionStatus);
  });

  el.sendBtn.addEventListener("click", () => {
    if (connectStatus) {
      ws.send(wsInput.value);
      wsInput.value = "";
    }
  });

  el.connectBtn.addEventListener("click", connectToServerHandler);
}

// DOM
function renderConnectionStatus(connected = false) {
  el.isConnected.innerText = connected ? "✅" : "❌";
}

// DOM
function getElements() {
  el.disconnectBtn = el2("disconnectBtn");
  el.wsInput = el2("wsInput");
  el.sendBtn = el2("sendBtn");
  el.isConnected = el2("connectStatusP");
  el.connectBtn = el2("connectBtn");
  el.userId = el2("userId");
  el.userId.innerText = uuid();
}
