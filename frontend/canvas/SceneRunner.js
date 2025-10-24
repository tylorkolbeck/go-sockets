import Player from "./classes/Agent/Player.js";

export default class SceneRunner {
  scene;
  _ownerId;
  eventBus;

  constructor(eventBus, scene, ownerId) {
    this.scene = scene;
    this._ownerId = ownerId;
    this.eventBus = eventBus;
  }

  get ownerId() {
    return this._ownerId;
  }

  init() {
    p5.disableFriendlyErrors = true;
    this.bindP5Functions();
  }

  playerJoin(playerConfig, isOwner = false) {
    this.scene.addPlayer(new Player(playerConfig.id, isOwner));
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
      this.eventBus.dispatch("sceneKeyPressed", {
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
      this.eventBus.dispatch("sceneKeyPressed", {
        key: "w",
        keyCode: 87,
      });
    }
    if (keyIsDown(83)) {
      this.eventBus.dispatch("sceneKeyPressed", {
        key: "s",
        keyCode: 83,
      });
    }
    if (keyIsDown(65)) {
      this.eventBus.dispatch("sceneKeyPressed", {
        key: "a",
        keyCode: 65,
      });
    }
    if (keyIsDown(68)) {
      this.eventBus.dispatch("sceneKeyPressed", {
        key: "d",
        keyCode: 68,
      });
    }
  }
}
