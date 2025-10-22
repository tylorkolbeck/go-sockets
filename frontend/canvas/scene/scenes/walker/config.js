export const config = {
  canvas: {
    domId: "canvas",
    height: 800,
    width: 800,
    background: {
      r: 1,
      g: 1,
      b: 1,
    },
  },
  ui: {
    x: 700,
    y: 20,
    controls: {
      timerConfig: {
        color: "red",
        x: 0,
        y: 0,
        fontSize: 24,
        label: "Run Time:",
      },

      frameRateConfig: {
        color: "red",
        x: 0,
        y: 15,
        fontSize: 24,
        label: "Frame Rate:",
      },
    },
  },
};
