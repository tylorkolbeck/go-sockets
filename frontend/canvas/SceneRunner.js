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

  playerJoin(playerConfig) {
    this.scene.addPlayer(new Player());
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
    };
  }
}
