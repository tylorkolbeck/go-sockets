import Player from "./classes/Agent/Player.js";
import eventBus from "./classes/Event/EventSystem.js";

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
    this.bindP5Functions();
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

  bindP5Functions() {
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
    if (keyIsDown(87)) {
      eventBus.dispatch("sceneKeyPressed", {
        key: "w",
        keyCode: 87,
      });
    }
    if (keyIsDown(83)) {
      eventBus.dispatch("sceneKeyPressed", {
        key: "s",
        keyCode: 83,
      });
    }
    if (keyIsDown(65)) {
      eventBus.dispatch("sceneKeyPressed", {
        key: "a",
        keyCode: 65,
      });
    }
    if (keyIsDown(68)) {
      eventBus.dispatch("sceneKeyPressed", {
        key: "d",
        keyCode: 68,
      });
    }
  }
}
