# Frontend Directory Structure

This frontend follows modern JavaScript development best practices with clear separation of concerns.

## Directory Structure

```
frontend/
├── public/               # Static assets and entry point
│   ├── index.html       # Main HTML file (entry point)
│   └── style.css        # Global styles
├── src/                 # Source code (ES6 modules)
│   ├── main.js          # Application entry point
│   ├── core/            # Core system files
│   │   ├── SceneRunner.js     # p5.js scene management
│   │   └── EventSystem.js     # Event bus system
│   ├── game/            # Game-specific code
│   │   ├── entities/    # Game entities (Player, etc.)
│   │   └── scenes/      # Game scenes
│   │       ├── BaseScene.js   # Base scene class
│   │       └── walker/        # Walker scene implementation
│   ├── network/         # Network/WebSocket code
│   │   ├── SocketManager.js   # WebSocket management
│   │   └── models.js          # Network message models
│   ├── ui/              # User interface components
│   │   └── DOMController.js   # DOM manipulation
│   └── utils/           # Utility functions
│       ├── uuid.js      # UUID generation
│       └── objectMerge.js     # Object utilities
└── assets/              # Static assets (images, libraries)
    └── lib/             # Third-party libraries
        └── p5/          # p5.js library files
```

## Key Benefits

1. **Clear Separation**: Source code vs static assets
2. **Module Organization**: Logical grouping by functionality
3. **Maintainability**: Easy to find and modify specific features
4. **Scalability**: Easy to add new scenes, entities, or UI components
5. **Best Practices**: Follows modern frontend architecture patterns

## Usage

- Serve the `public/` directory from your web server
- The entry point is `public/index.html`
- All JavaScript modules are loaded as ES6 modules
- p5.js library is loaded globally before modules

## Development

- Add new game entities to `src/game/entities/`
- Add new scenes to `src/game/scenes/`
- Add UI components to `src/ui/`
- Add utilities to `src/utils/`
- Network-related code goes in `src/network/`