import { UiController } from "../classes/ui/UiController.js";
import { TimerUi } from "../classes/ui/TimerUi.js";
import { FrameRateUi } from "../classes/ui/FrameRateUi.js";
import { initGlobals } from "../globals.js";
import { initDom } from "../dom.js";

export class BaseScene {
  uiController = new UiController([]);
  _config = null;
  _players = {};
  _agents = [];
  _isSetupComplete = false;
  _worldDataReceived = false;

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

  constructor(config) {
    this._config = config;
  }

  handleWorldSnapshot(worldData) {
    // Update world properties
    this.config.canvas.background =
      worldData.worldbg || this.config.canvas.background;
    this.config.canvas.height =
      worldData.worldheight || this.config.canvas.height;
    this.config.canvas.width = worldData.worldwidth || this.config.canvas.width;

    this._worldDataReceived = true;

    // If we haven't set up the canvas yet, do it now
    if (!this._isSetupComplete) {
      this.performSetup();
    }
  }

  handlePlayerSnapshot(snapshot) {
    for (const id in snapshot.players) {
      this.players[id].setPosition(snapshot.players[id].pos);
    }
  }

  getPlayer(playerId) {
    return this.players[playerId];
  }

  getPlayers() {
    return this.players;
  }

  addPlayer(player) {
    this.players[player._id] = player;
  }

  removePlayer(playerId) {
    console.log("Attempting to remove player:", playerId);
    console.log("Players before deletion:", this.players);
    console.log("Player exists?", playerId in this.players);

    delete this.players[playerId];

    console.log("Players after deletion:", this.players);
    console.log("Player still exists?", playerId in this.players);
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
    // Only perform initial setup if we already have world data
    if (this._worldDataReceived) {
      this.performSetup();
    }
  }

  performSetup() {
    if (this._isSetupComplete) return; // Prevent double setup

    const canvasEl = document.getElementById(this.config.canvas.domId);
    createCanvas(this.config.canvas.width, this.config.canvas.height, canvasEl);
    background(this.background.r, this.background.g, this.background.b);

    this.initUiComponents();
    this.uiController.setup();
    this._isSetupComplete = true;
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
    // Only draw if setup is complete
    if (!this._isSetupComplete) {
      // Show a loading message or just return
      if (window.fill && window.text) {
        fill(255);
        text("Waiting for world data from server...", 20, 30);
      }
      return;
    }

    background(this.background.r, this.background.g, this.background.b);
    this.uiController.draw();

    for (const id in this.players) {
      this.players[id].draw();
    }
  }
}
