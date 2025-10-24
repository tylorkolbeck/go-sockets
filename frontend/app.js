import SceneRunner from "./canvas/SceneRunner.js";
import SocketManager from "./canvas/handlers/SockerManager.js";
import { WalkerScene } from "./canvas/scene/scenes/walker/WalkerScene.js";
import { config } from "./canvas/scene/scenes/walker/config.js";
import eventBus from "./canvas/classes/Event/EventSystem.js";
import { uuid } from "./canvas/util/uuid.js";

const ownerId = uuid();
const scene = new WalkerScene(config);
const sceneRunner = new SceneRunner(ownerId);
const socketManager = new SocketManager(sceneRunner, {
  host: "localhost",
  port: "8000",
  ownerId,
});

const el2 = (id) => document.getElementById(id);

const el = {
  disconnectBtn: null,
  wsInput: null,
  sendBtn: null,
  isConnected: null,
  connectBtn: null,
  userId: null,
};

initDom();
registerGameEventHandlers();
renderConnectionStatus(false);

// SET AND INIT SCENE
sceneRunner.setScene(scene).init();

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
  eventBus.subscribe("sceneKeyPressed", handleSceneKeyPress);
  eventBus.subscribe("serverconnected", handleServerConnChange);
  eventBus.subscribe("serverdisconnected", handleServerConnChange);
}

// TODO: pass is connected as an app level event
function connectToServerHandler() {
  socketManager.connectToServer();
}

// DOM
function initDom() {
  getElements();
  registerDomEventHandlers();
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

function handleServerConnChange(event) {
  renderConnectionStatus(event?.data);
}

// DOM
function renderConnectionStatus(connected) {
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
  el.userId.innerText = ownerId;
}
