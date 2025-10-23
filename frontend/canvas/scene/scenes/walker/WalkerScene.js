import { BaseScene } from "../../BaseScene.js";
import { Walker } from "../../../classes/Walker.js";
import { AccelVisualizer } from "../../../classes/AccelVisualizer.js";
import { vec } from "../../../globals.js";

export class WalkerScene extends BaseScene {
  paused = false;
  agents = [];
  accelVisualizer;

  constructor(config, eventBus) {
    super(config, eventBus);
  }

  preload() {
    super.preload();
    this.accelVisualizer = new AccelVisualizer(vec(100, 100), 10);
  }

  setup() {
    super.setup();

    this.agents.push(new Walker(vec(this.width / 2, this.height / 2)));
  }

  keyPressed() {
    super.keyPressed();
    if (key === "p") this.paused = !this.paused;
    if (key === "w") {
    }
  }

  draw() {
    super.draw();
    if (this.paused) return;

    // const accelVal = vec(random(-1, 1), random(-1, 1));
    // this.agents.forEach((w) => {
    //   w.applyForce(accelVal);
    //   w.draw();
    // });

    // this.accelVisualizer.applyAccel(accelVal);
    // this.accelVisualizer.draw();
  }
}
