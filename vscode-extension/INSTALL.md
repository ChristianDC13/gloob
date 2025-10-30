# Installing the Gloob VS Code Extension

## Option 1: Development Mode (Recommended for testing)

1. Open the `vscode-extension` folder in VS Code
2. Press `F5` to launch a new Extension Development Host window
3. In the new window, open any `.gloob` file
4. Syntax highlighting should be working!

## Option 2: Install from VSIX

### Build the VSIX file

```bash
# First, install vsce globally (if not already installed)
npm install -g vsce

# Navigate to the extension folder
cd vscode-extension

# Package the extension
vsce package
```

This creates a `.vsix` file that you can install.

### Install the VSIX in VS Code

1. Open VS Code
2. Go to Extensions (Cmd+Shift+X on Mac, Ctrl+Shift+X on Windows/Linux)
3. Click the `...` menu at the top right
4. Select "Install from VSIX..."
5. Choose the generated `.vsix` file

## Features Included

- ✅ Syntax highlighting for all Gloob keywords
- ✅ String literal highlighting (single and double quotes)
- ✅ Number highlighting
- ✅ Comment support (//)
- ✅ Operator highlighting
- ✅ Bracket matching
- ✅ Auto-closing brackets and quotes

## Testing

Open the `sample.gloob` file in the extension folder to see examples of all syntax highlighting features.

