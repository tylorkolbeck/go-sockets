import Player from "../game/entities/Player.js";
import eventBus from "./EventSystem.js";

export default class SceneRunner {
  scene;
  _ownerId;

  constructor(ownerId) {
    this._ownerId = ownerId;
  }

  get ownerId() {
    return this._ownerId;
  }

  setScene(scene) {
    this.scene = scene;
    return this;
  }

  init() {
    if (!this.scene) {
      throw new Error("Scene has not been set");
    }
    p5.disableFriendlyErrors = true;
    this.startP5Scene();
  }

  handleWorldSnapshot(worldData) {
    this.scene.handleWorldSnapshot(worldData);
  }

  handlePlayerListUpdate(playerIds) {
    playerIds.forEach((pid) => {
      if (!this.scene.getPlayer(pid)) {
        this.playerJoin({ id: pid });
      }
    });
  }

  playerJoin(playerConfig, isOwner = false) {
    if (!this.scene.getPlayer(playerConfig.id)) {
      this.scene.addPlayer(new Player(playerConfig.id, isOwner));
    }
  }

  playerLeft(data) {
    this.scene.removePlayer(data.id);
  }

  handlePlayerSnapshot(snapshot) {
    this.scene.handlePlayerSnapshot(snapshot);
  }

  startP5Scene() {
    // Make functions global for p5.js
    window.preload = () => {
      this.scene.preload();
    };

    window.keyPressed = () => {
      eventBus.dispatch("sceneKeyPressed", {
        key,
        keyCode,
      });
    };

    window.setup = () => {
      this.scene.setup();
    };

    window.draw = () => {
      this.scene.draw();
      this.checkHeldKeys();
    };
  }

  checkHeldKeys() {
    const keyMappings = [
      { keyCode: 87, key: "w" }, // W
      { keyCode: 83, key: "s" }, // S
      { keyCode: 65, key: "a" }, // A
      { keyCode: 68, key: "d" }, // D
    ];

    keyMappings.forEach(({ keyCode, key }) => {
      if (keyIsDown(keyCode)) {
        eventBus.dispatch("sceneKeyPressed", { key, keyCode });
      }
    });
  }
}
