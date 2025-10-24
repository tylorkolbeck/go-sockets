import Player from "./classes/Agent/Player.js";

export default class Canvas {
  scene;

  constructor(scene) {
    this.scene = scene;
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
      if (this.scene.keyPressed) {
        this.scene.keyPressed();
      }
    };

    window.setup = () => {
      this.scene.setup();
    };

    window.draw = () => {
      this.scene.draw();

      if (this.scene.keyDown) {
        this.checkHeldKeys();
      }
    };
  }

  checkHeldKeys() {
    if (keyIsDown(87)) {
      this.scene.keyDown("w", 87);
    }
    if (keyIsDown(83)) {
      this.scene.keyDown("s", 83);
    }
    if (keyIsDown(65)) {
      this.scene.keyDown("a", 65);
    }
    if (keyIsDown(68)) {
      this.scene.keyDown("d", 68);
    }
  }
}
