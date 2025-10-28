import SceneRunner from "./core/SceneRunner.js";
import SocketManager from "./network/SocketManager.js";
import { WalkerScene } from "./game/scenes/walker/WalkerScene.js";
import { config } from "./game/scenes/walker/config.js";
import eventBus from "./core/EventSystem.js";
import { uuid } from "./utils/uuid.js";
import { DOMController } from "./ui/DOMController.js";

const ownerId = uuid();
const scene = new WalkerScene(config);
const sceneRunner = new SceneRunner(ownerId);
const socketManager = new SocketManager(sceneRunner, {
  host: "localhost",
  port: "8000",
  ownerId,
});

const domController = new DOMController(socketManager, ownerId);

// Initialize the application
function init() {
  domController.init();
  registerGameEventHandlers();
  domController.renderConnectionStatus(false);

  // Set and initialize scene
  sceneRunner.setScene(scene).init();
}

// Game event handlers
function handleSceneKeyPress(event) {
  const keyMap = {
    w: { type: "input", up: true },
    d: { type: "input", right: true },
    s: { type: "input", down: true },
    a: { type: "input", left: true },
  };

  const inputData = keyMap[event.data.key];
  if (inputData) {
    socketManager.send({
      id: ownerId,
      msg: inputData,
    });
  }
}

function registerGameEventHandlers() {
  eventBus.subscribe("sceneKeyPressed", handleSceneKeyPress);
  eventBus.subscribe("serverconnected", (event) =>
    domController.handleServerConnChange(event.data)
  );
  eventBus.subscribe("serverdisconnected", (event) =>
    domController.handleServerConnChange(event.data)
  );
}

// Start the application
init();
