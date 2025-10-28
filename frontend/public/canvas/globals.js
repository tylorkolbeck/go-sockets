// App Variables
export const globals = {
  vec: null,
};

export function initGlobals() {
  // p5.js aliases
  globals.vec = createVector;

  // Make vec available globally for easier access
  window.vec = globals.vec;
}

// Export vec as a function that returns the current value
export const vec = (...args) => globals.vec(...args);
