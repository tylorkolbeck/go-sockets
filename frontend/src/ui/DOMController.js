export class DOMController {
  socketManager;
  ownerId;
  elements = {};

  constructor(socketManager, ownerId) {
    this.socketManager = socketManager;
    this.ownerId = ownerId;
  }

  init() {
    this.getElements();
    this.registerDomEventHandlers();
  }

  getElements() {
    const el2 = (id) => document.getElementById(id);

    this.elements = {
      disconnectBtn: el2("disconnectBtn"),
      // wsInput: el2("wsInput"),
      // sendBtn: el2("sendBtn"),
      isConnected: el2("connectStatusP"),
      connectBtn: el2("connectBtn"),
      userId: el2("userId"),
    };

    this.elements.userId.innerText = this.ownerId;
  }

  registerDomEventHandlers() {
    this.elements.disconnectBtn.addEventListener("click", () => {
      this.socketManager.close(() => this.renderConnectionStatus(false));
    });

    // this.elements.sendBtn.addEventListener("click", () => {
    //   if (this.socketManager.isConnected) {
    //     // Handle custom message sending if needed
    //     this.elements.wsInput.value = "";
    //   }
    // });

    this.elements.connectBtn.addEventListener("click", () => {
      this.socketManager.connectToServer();
    });
  }

  handleServerConnChange(isConnected) {
    this.renderConnectionStatus(isConnected);
  }

  renderConnectionStatus(connected) {
    this.elements.isConnected.innerText = connected ? "✅" : "❌";
  }
}
