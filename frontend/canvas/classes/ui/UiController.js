export class UiController {
  controls = [];
  postition;

  constructor(controls, position) {
    if (!controls || !Array.isArray(controls)) {
      throw new Error(
        "Controls needs to be an empty array or array of controls"
      );
    }
    this.controls = controls;
    this.position = position || { x: 0, y: 0 };
  }

  addControl(control) {
    if (Array.isArray(control)) {
      this.controls = [...this.controls, ...control];
    } else {
      this.controls.push(control);
    }
  }

  setup() {
    this.controls.forEach((c) => c.setup());
  }

  draw() {
    this.controls.forEach((c) => c.draw());
  }
}
