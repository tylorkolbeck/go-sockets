import { BaseUi } from "./BaseUi.js";

export class FrameRateUi extends BaseUi {
  frameCount = 0;
  frameRate = 0;

  constructor(config, parentPos) {
    super(config, parentPos);
  }

  draw() {
    super.draw();
    this.frameCount++;
    if (this.frameCount % 10 === 0) {
      this.frameRate = Math.trunc(frameRate());
    }
    text(`${this.config.label} ${this.frameRate}`, this.x, this.y);
  }
}
