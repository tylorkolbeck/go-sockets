import { UiController } from "../classes/ui/UiController.js";
import { TimerUi } from "../classes/ui/TimerUi.js";
import { FrameRateUi } from "../classes/ui/FrameRateUi.js";
import { initGlobals } from "../globals.js";
import { initDom } from "../dom.js";

export class BaseScene {
  uiController = new UiController([]);
  _config = null;
  _players = [];
  _agents = [];

  get players() {
    return this._players;
  }
  set players(players) {
    this._players = players;
  }
  get config() {
    return this._config;
  }
  get width() {
    return this.config.canvas.width;
  }
  get height() {
    return this.config.canvas.height;
  }
  set background(rgb) {
    this.config.canvas.background = rgb;
  }
  get background() {
    return this.config.canvas.background;
  }

  constructor(config, eventBus) {
    this._config = config;
    this.eventBus = eventBus;
  }

  getPlayer(playerId) {
    return this.players.find((p) => p.id === playerId);
  }

  addPlayer(player) {
    this.players.push(player);
  }

  addUiControl(control) {
    this.uiController.addControl(control);
  }

  setBackground(rgb) {
    this.background = rgb;
  }

  setCanvasDimensions(width, height) {
    this.config.canvas.height = height;
    this.config.canvas.width = width;
  }

  preload() {
    initGlobals();
    initDom();
  }

  setup() {
    const canvasEl = document.getElementById(this.config.canvas.domId);
    createCanvas(this.config.canvas.width, this.config.canvas.height, canvasEl);
    background(this.background.r, this.background.g, this.background.b);

    this.initUiComponents();
    this.uiController.setup();
  }

  initUiComponents() {
    for (const key in this.config.ui.controls) {
      switch (key) {
        case "timerConfig":
          this.addUiControl(
            new TimerUi(this.config.ui.controls.timerConfig, {
              x: this.config.ui.x,
              y: this.config.ui.y,
            })
          );
        case "frameRateConfig":
          this.addUiControl(
            new FrameRateUi(this.config.ui.controls.frameRateConfig, {
              x: this.config.ui.x,
              y: this.config.ui.y,
            })
          );
      }
    }
  }

  draw() {
    background(this.background.r, this.background.g, this.background.b);
    this.uiController.draw();

    this.players.forEach((p) => {
      p.draw();
    });
  }

  keyPressed() {
    this.eventBus.dispatch("sceneKeyPressed", {
      key,
      keyCode,
    });
  }
}
