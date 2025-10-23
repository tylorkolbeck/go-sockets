export class EventEmitter {
  events;

  constructor() {
    this.events = {};
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
