import { vec } from "../globals.js";

export class Walker {
  position;
  size;
  speed;

  stepList = [];

  constructor(loc, size) {
    this.position = loc || vec(0, 0);
    this.size = size || 5;
  }

  setSize(value) {
    this.size = value;
  }

  applyForce(vec) {
    this.position.add(vec);
  }

  replay() {}

  getStepList() {
    return this.stepList;
  }

  draw() {
    noStroke();
    fill(0, 204, 0);
    circle(this.position.x, this.position.y, this.size);
    this.stepList.push({
      x: this.position.x,
      y: this.position.y,
      z: this.position.y,
    });

    this.stepList.forEach((s) => circle(s.x, s.y, this.size));
  }
}
