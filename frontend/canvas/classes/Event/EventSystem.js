class EventSystem {
  events;

  constructor() {
    if (EventSystem.instance == null) {
      this.events = {};
      EventSystem.instance = this;
    }

    return EventSystem.instance;
  }

  subscribe(event, callback) {
    if (!this.events[event]) {
      this.events[event] = [];
    }

    this.events[event].push(callback);
  }

  unsubscribe(event, callback) {
    if (this.events[event]) {
      this.events[event] = this.events[event].filter((cb) => cb !== callback);
    }
  }

  dispatch(event, data) {
    if (this.events[event]) {
      this.events[event].forEach((callback) =>
        callback({
          data,
          timestamp: Date.now(),
          event: event,
        })
      );
    }
  }
}

const eventSystem = new EventSystem();
Object.freeze(eventSystem);

export default eventSystem;
