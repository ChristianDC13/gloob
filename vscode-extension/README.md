# Gloob Language Support for VS Code

Syntax highlighting support for the Gloob programming language.

## Features

- Syntax highlighting for Gloob files (`.gloob`)
- Support for all Gloob keywords and operators
- String and number literals
- Comments
- Object syntax
- Automatic bracket matching

## Installation

1. Open VS Code
2. Go to Extensions (Cmd+Shift+X)
3. Click the "..." menu and select "Install from VSIX..."
4. Select this extension's `.vsix` file

Or install for development:
1. Clone this repository
2. Open the `vscode-extension` folder in VS Code
3. Press F5 to launch a new VS Code window with the extension loaded

## Building

```bash
# Install vsce (VS Code Extension Manager)
npm install -g vsce

# Build the extension
cd vscode-extension
vsce package
```

This will create a `.vsix` file that can be installed in VS Code.

