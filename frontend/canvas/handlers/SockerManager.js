export default class SocketManager {
  sceneRunner;
  wsUri;
  ws;
  isConnected = false;

  constructor(sceneRunner, wsUri) {
    this.wsUri = wsUri;
    if (!sceneRunner)
      throw new Error("Scene runner not provided to Message Handler");
    this.sceneRunner = sceneRunner;
  }

  send(data) {
    if (this.ws) {
      const stringData = JSON.stringify(data);
      this.ws.send(stringData);
    }
  }

  connectToServer(cb) {
    this.ws = new WebSocket(this.wsUri + `?id=${this.sceneRunner.ownerId}`);
    this._registerSockerHandlers(cb);
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
    }
  }

  handleSnapshot(data) {
    this.sceneRunner.handlePlayerSnapshot(data);
  }

  handleDisconnect() {}

  handleUpdatedPlayerList(data) {
    data.playerIds.forEach((id) => {
      if (id != this.sceneRunner.ownerId) {
        this.sceneRunner.playerJoin({
          id: id,
        });
      }
    });
  }

  handleOnClose(event) {
    console.log("Disconnected", event);
    this.isConnected = false;
  }

  handleOnOpen(event, cb) {
    console.log("Socket connection established: ", event);
    this.isConnected = true;
    cb(true);
    this.sceneRunner.playerJoin({
      id: this.sceneRunner.ownerId,
      isOwner: true,
    });
  }

  close(cb) {
    this.ws.close();
    this.isConnected = false;
    if (cb) cb();
  }

  _registerSockerHandlers(cb) {
    this.ws.onmessage = (event) => this.handleOnMessage(event);
    this.ws.onopen = (event) => this.handleOnOpen(event, cb);
    this.ws.onclose = (event) => this.handleOnClose(event);
  }
}
