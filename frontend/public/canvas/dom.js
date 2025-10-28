// DOM ELEMENTS

export function initDom() {
  getDomReferences();
  registerUiEventListeners();
}

function getDomReferences() {}

function registerUiEventListeners() {}

// Util Functions
export const el = (id) => {
  const el = document.getElementById(id);
  if (!el) {
    console.warn("No element found for: ", id);
    return document.createElement("div"); // Safely handle missing element
  }

  return el;
};
