import { BaseUi } from "./BaseUi.js";

export class TimerUi extends BaseUi {
  elapsedTime = 0;

  constructor(config, parentPos) {
    super(config, parentPos);
  }

  start() {
    setInterval(() => {
      this.elapsedTime += 1;
    }, 1000);
  }

  setup() {
    super.setup();
    this.start();
  }

  draw() {
    super.draw();
    text(`${this.config.label} ${this.formatElapsedTime()}`, this.x, this.y);
  }

  formatElapsedTime() {
    if (this.elapsedTime < 60) return `${this.elapsedTime} secs`;

    const mins = Math.floor(this.elapsedTime / 60);
    const remainingSecs = this.elapsedTime % 60;

    const formattedMins = String(mins).padStart(2, "0");
    const formattedSecs = String(remainingSecs).padStart(2, "0");

    return `${formattedMins}:${formattedSecs}`;
  }
}
