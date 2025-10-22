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

  constructor() {
    this.id = uuid();
    this.position = { x: 50, y: 50, z: 0 };
    console.log(`[INITIALIZED PLAYER] - ID: ${this.id}`);
  }

  update() {}

  draw() {
    noStroke();
    fill(255, 0, 0);
    circle(this.position.x, this.position.y, this.size);
  }
}
