import SceneRunner from "./canvas/SceneRunner.js";
import { WalkerScene } from "./canvas/scene/scenes/walker/WalkerScene.js";
import { config } from "./canvas/scene/scenes/walker/config.js";
import { EventEmitter } from "./canvas/classes/EventEmitter/EventEmitter.js";
import { uuid } from "./canvas/util/uuid.js";
import { handleOnMessage } from "./canvas/handlers/socket.js";

const gameEventBus = new EventEmitter();

const wsConfig = {
  autoConnect: false,
};

let ws;
let sceneRunner;

let scene = new WalkerScene(config, gameEventBus);

const el2 = (id) => document.getElementById(id);

let connectStatus = false;

const el = {
  disconnectBtn: null,
  wsInput: null,
  sendBtn: null,
  isConnected: null,
  connectBtn: null,
  userId: null,
};

getElements();
registerEventHandlers();
renderConnectionStatus();

sceneRunner = new SceneRunner(scene);

sceneRunner.init();

if (wsConfig.autoConnect) {
  const userId = getUserId();
  if (userId) connectToSocket(userId);
}

function setIsConnected(connected) {
  connectStatus = connected;
  renderConnectionStatus();
}

function renderConnectionStatus() {
  el.isConnected.innerText = connectStatus ? "✅" : "❌";
}

function getElements() {
  el.disconnectBtn = el2("disconnectBtn");
  el.wsInput = el2("wsInput");
  el.sendBtn = el2("sendBtn");
  el.isConnected = el2("connectStatusP");
  el.connectBtn = el2("connectBtn");
  el.userId = el2("userId");
  el.userId.value = uuid();
}

function handleSceneKeyPress(event) {
  switch (event.data.key) {
    case "w":
      ws.send(
        JSON.stringify({
          id: getUserId(),
          msg: {
            type: "input",
            up: true,
          },
        })
      );
      break;
    case "d":
      ws.send(
        JSON.stringify({
          id: getUserId(),
          msg: {
            type: "input",
            right: true,
          },
        })
      );
      break;
    case "s":
      ws.send(
        JSON.stringify({
          id: getUserId(),
          msg: {
            type: "input",
            down: true,
          },
        })
      );
      break;
    case "a":
      ws.send(
        JSON.stringify({
          id: getUserId(),
          msg: {
            type: "input",
            left: true,
          },
        })
      );
      break;
    default:
      break;
  }
}

function registerEventHandlers() {
  gameEventBus.subscribe("sceneKeyPressed", handleSceneKeyPress);

  el.disconnectBtn.addEventListener("click", () => {
    ws.close();
  });

  el.sendBtn.addEventListener("click", () => {
    if (connectStatus) {
      ws.send(wsInput.value);
      wsInput.value = "";
    }
  });

  el.connectBtn.addEventListener("click", () => {
    const userId = getUserId();
    if (userId) connectToSocket(userId);
  });
}

function connectToSocket(userId) {
  const wsUri = `ws://localhost:8000/connect?id=${userId}`;
  ws = new WebSocket(wsUri);
  ws.onopen = function (event) {
    console.log("Socket connection established: ", event);
    setIsConnected(true);
    sceneRunner.playerJoin({
      id: getUserId(),
      isOwner: true,
    });
  };

  ws.onmessage = function (event) {
    handleOnMessage(event, sceneRunner, getUserId());
    // const data = JSON.parse(event.data);
    // if (data.type === "snapshot") {
    //   // Update players
    //   sceneRunner.handlePlayerSnapshot(data);
    // }

    // if (data.type === "updatedplayerlist") {
    //   data.playerIds.forEach((id) => {
    //     if (id != getUserId) {
    //       sceneRunner.playerJoin({
    //         id: id,
    //       });
    //     }
    //   });
    // }
  };

  ws.onclose = function (event) {
    console.log("Disconnected", event);
    setIsConnected(false);
  };
}

function getUserId() {
  return el.userId.value;
}
