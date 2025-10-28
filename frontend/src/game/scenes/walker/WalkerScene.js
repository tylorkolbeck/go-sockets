import { BaseScene } from "../BaseScene.js";

export class WalkerScene extends BaseScene {
  paused = false;
  agents = [];

  constructor(config) {
    super(config);
  }

  preload() {
    super.preload();
  }

  setup() {
    super.setup();
  }

  keyPressed() {
    if (key === "p") this.paused = !this.paused;
  }

  draw() {
    super.draw();
    if (this.paused) return;
  }
}
