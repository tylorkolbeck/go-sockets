import eventBus from "../classes/Event/EventSystem.js";

export default class SocketManager {
  sceneRunner;
  ws;
  isConnected = false;
  config;

  constructor(sceneRunner, config) {
    this.config = config;
    if (!sceneRunner)
      throw new Error("Scene runner not provided to Message Handler");
    this.sceneRunner = sceneRunner;
  }

  send(data) {
    if (this.isConnected) {
      const stringData = JSON.stringify(data);
      this.ws.send(stringData);
    } else {
      console.warn("There is no connection to the server");
    }
  }

  connectToServer() {
    this.ws = new WebSocket(
      `ws://${this.config.host}:${this.config.port}/connect?id=${this.config.ownerId}`
    );
    this.registerSockerHandlers();
  }

  handleOnMessage(event) {
    const data = JSON.parse(event.data);

    switch (data.type) {
      case "snapshot":
        this.handleSnapshot(data);
        break;
      case "updatedplayerlist":
        this.handleUpdatedPlayerList(data);
        break;
      case "playerjoined":
        this.sceneRunner.playerJoin({ id: data.id });
        break;
      case "playerleft":
        this.sceneRunner.playerLeft(data);
        break;
      case "worldsettings":
        this.sceneRunner.handleWorldSnapshot(data);
        break;
    }
  }

  handleSnapshot(data) {
    this.sceneRunner.handlePlayerSnapshot(data);
  }

  handleDisconnect() {
    this.isConnected = false;
  }

  handleUpdatedPlayerList(playerList) {
    this.sceneRunner.handlePlayerListUpdate(playerList.playerIds);
  }

  handleOnClose() {
    this.isConnected = false;
    eventBus.dispatch("serverdisconnected", false);
  }

  handleOnOpen() {
    this.isConnected = true;
    eventBus.dispatch("serverconnected", true);
  }

  close(cb) {
    this.ws.close();
    this.isConnected = false;
    if (cb) cb();
  }

  registerSockerHandlers() {
    this.ws.onmessage = (event) => this.handleOnMessage(event);
    this.ws.onopen = (event) => this.handleOnOpen(event);
    this.ws.onclose = (event) => this.handleOnClose(event);
  }
}
