const wsUri = "ws://localhost:8000/connect";

let ws;

const el = (id) => document.getElementById(id);

let connectStatus = false;

const EL = {
  disconnectBtn: null,
  wsInput: null,
  sendBtn: null,
  connectStatusP: null,
  connectBtn: null,
};

getElements();
registerEventHandlers();
styleConnectStatus();

function setIsConnected(connected) {
  connectStatus = connected;
  styleConnectStatus();
}

function styleConnectStatus() {
  EL.connectStatusP.style.padding = "6px";
  EL.connectStatusP.style.borderRadius = "4px";
  EL.connectStatusP.style.display = "inline-block";
  EL.connectStatusP.style.marginTop = "20px";

  if (connectStatus) {
    EL.connectStatusP.innerText = "Connected";
    EL.connectStatusP.style.background = "green";
  } else {
    EL.connectStatusP.innerText = "Disconnected";
    EL.connectStatusP.style.background = "red";
  }
}

function getElements() {
  EL.disconnectBtn = el("disconnectBtn");
  EL.wsInput = el("wsInput");
  EL.sendBtn = el("sendBtn");
  EL.connectStatusP = el("connectStatusP");
  EL.connectBtn = el("connectBtn");
}

function registerEventHandlers() {
  EL.disconnectBtn.addEventListener("click", () => {
    ws.close();
  });

  EL.sendBtn.addEventListener("click", () => {
    if (connectStatus) {
      ws.send(wsInput.value);
      wsInput.value = ""
    }
  });

  EL.connectBtn.addEventListener("click", () => {
    ws = new WebSocket(wsUri);
    ws.onopen = function (event) {
      console.log("Socket connection established: ", event);
      setIsConnected(true);

      ws.send("Hello World");
    };

    ws.onmessage = function (event) {
      console.log("Message: ", event.data);
    };

    ws.onclose = function (event) {
      console.log("Disconnected", event);
      setIsConnected(false);
    };
  });
}
