import { mergeObject } from "../../util/objectMerge.js";

export class BaseUi {
  config = {
    color: "red",
    x: 0,
    y: 0,
    fontSize: 23,
    label: "Elapsed Time:",
  };

  parentX = 0;
  parentY = 0;

  get x() {
    return this.config.x + this.parentX;
  }

  get y() {
    return this.config.y + this.parentY;
  }

  constructor(config, parent) {
    if (parent) {
      this.parentX = parent.x;
      this.parentY = parent.y;
    }
    this.config = mergeObject(this.config, config);
  }

  setup() {}

  draw() {
    fill(this.config.color);
  }
}
