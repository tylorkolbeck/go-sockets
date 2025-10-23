import { uuid } from "../../util/uuid.js";

export default class Player {
  _id = "";
  position;
  size = 10;

  get id() {
    return this._id;
  }
  set id(value) {
    this._id = value;
  }

  constructor(id) {
    this.id = id;
    this.position = { x: 0, y: 0, z: 0 };
    console.log(`[INITIALIZED PLAYER] - ID: ${this.id}`);
  }

  setPosition(vec3) {
    this.position = vec3;
  }

  update() {}

  draw() {
    noStroke();
    fill(255, 0, 0);
    circle(this.position.x, this.position.y, this.size);
  }
}
