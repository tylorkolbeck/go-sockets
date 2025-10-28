import { vec } from "../globals.js";

export class AccelVisualizer {
  position;
  acceleration = vec(0, 0, 0);

  constructor(position, size) {
    this.position = position;
    this.size = size;
  }

  applyAccel(vec) {
    this.acceleration = vec.mult(5);
  }

  draw() {
    circle(
      this.position.x + this.acceleration.x,
      this.position.y + this.acceleration.y,
      this.size
    );
    this.acceleration = vec(0, 0, 0);
  }
}
