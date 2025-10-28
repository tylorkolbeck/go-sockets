export default class Player {
  _id = "";
  position;
  size = 10;
  color;

  get id() {
    return this._id;
  }
  set id(value) {
    this._id = value;
  }

  constructor(id) {
    this.id = id;
    this.position = { x: 0, y: 0, z: 0 };
    this.color = { r: random(0, 255), g: random(0, 255), b: random(0, 255) };
    console.log(`[INITIALIZED PLAYER] - ID: ${this.id}`);
  }

  setPosition(vec3) {
    this.position = vec3;
  }

  update() {}

  draw() {
    noStroke();
    fill(this.color.r, this.color.g, this.color.b);
    circle(this.position.x, this.position.y, this.size);
  }
}
