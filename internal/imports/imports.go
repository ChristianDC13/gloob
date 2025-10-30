package imports

import (
	"fmt"
	"gloob-interpreter/internal/parser"
	"os"
	"path/filepath"
	"strings"
)

// ProcessImports processes all import statements recursively and returns
// a flattened list of statements with imports resolved.
// It handles circular import detection and relative path resolution.
func ProcessImports(program *parser.Program, basePath string) (*parser.Program, error) {
	visited := make(map[string]bool)

	// Get the directory of the base file
	baseDir := filepath.Dir(basePath)
	if baseDir == "" {
		baseDir = "."
	}

	// Process the program recursively
	statements, err := processStatementsWithImports(program.Statements, baseDir, visited)
	if err != nil {
		return nil, err
	}

	return &parser.Program{Statements: statements}, nil
}

// processStatementsWithImports recursively processes a list of statements,
// expanding any import statements into the imported file's statements.
func processStatementsWithImports(statements []parser.Statement, baseDir string, visited map[string]bool) ([]parser.Statement, error) {
	result := make([]parser.Statement, 0)

	for _, stmt := range statements {
		// Check if this is an import statement
		if importStmt, ok := stmt.(*parser.ImportStatement); ok {
			// Resolve the import path
			importPath := resolveImportPath(importStmt.Path, baseDir)

			// Check for circular imports
			absPath, err := filepath.Abs(importPath)
			if err != nil {
				return nil, fmt.Errorf("failed to resolve import path %s: %v", importPath, err)
			}

			if visited[absPath] {
				return nil, fmt.Errorf("circular import detected: %s", importPath)
			}

			// Mark this file as visited
			visited[absPath] = true

			// Read and parse the imported file
			importedStatements, err := loadAndParseFile(importPath, visited)
			if err != nil {
				return nil, fmt.Errorf("failed to import %s: %v", importStmt.Path, err)
			}

			// Add the imported statements to the result
			result = append(result, importedStatements...)

			// Unmark after processing (to allow the same file to be imported in different branches)
			// Comment this out if you want to prevent importing the same file multiple times
			// delete(visited, absPath)
		} else {
			// Not an import statement, add it as-is
			result = append(result, stmt)
		}
	}

	return result, nil
}

// resolveImportPath resolves an import path relative to the base directory.
// If the path doesn't have a .gloob extension, it adds one.
func resolveImportPath(importPath, baseDir string) string {
	// Add .gloob extension if not present
	if !strings.HasSuffix(importPath, ".gloob") && !strings.HasSuffix(importPath, ".gb") {
		importPath += ".gloob"
	}

	// If the path is absolute, use it as-is
	if filepath.IsAbs(importPath) {
		return importPath
	}

	// Otherwise, resolve it relative to the base directory
	return filepath.Join(baseDir, importPath)
}

// loadAndParseFile loads a file, parses it, and recursively processes its imports.
func loadAndParseFile(filePath string, visited map[string]bool) ([]parser.Statement, error) {
	// Read the file
	sourceCode, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	// Parse the file
	p := parser.NewParser(nil)
	program := p.ProduceAST(string(sourceCode))

	// Get the directory of this file for resolving its imports
	fileDir := filepath.Dir(filePath)

	// Process imports recursively
	statements, err := processStatementsWithImports(program.Statements, fileDir, visited)
	if err != nil {
		return nil, err
	}

	return statements, nil
}
