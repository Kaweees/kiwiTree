package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	// "path/filepath"
	// "strings"
)

// Represents the possible node types in the file tree.
type Node int

const (
	Root Node = iota
	Parent
	Branch
	Leaf
	FinalLeaf
	Symlink
)

// String method to convert the current node to a string.
func (n Node) String() string {
	return [...]string{".", "..", "├──", "│  ", "└──", "->"}[n]
}

func main() {
	// Get root directory from command line args or use current directory
	root := "."
	if len(os.Args) > 1 {
		root = os.Args[1]
	}

	err := printTree(root, "", true)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func getNodeType(fileInfo os.FileInfo) Node {
	if fileInfo.Name() == "." {
		return Root
	} else if fileInfo.Name() == ".." {
		return Parent
	} else if fileInfo.IsDir() {
		return Branch
	} else if fileInfo.Mode()&os.ModeSymlink != 0 {
		return Symlink
	} else if fileInfo.Mode()&os.ModeDir == 0 {
		return Leaf
	} else {
		return FinalLeaf
	}
}

func printTree(root string, prefix string, isLast bool) error {
	// Get file info for the current path
	fileInfo, err := os.Stat(root)

	if err != nil {
		return err
	}

	// Get the node type
	node := getNodeType(fileInfo)

	// Print the current node with appropriate prefix
	switch node {
	case Branch:
		if prefix == "" {
			fmt.Printf("%s├── %s/\n", prefix, fileInfo.Name())
			prefix += "│   "
		} else {
			fmt.Printf("%s├── %s/\n", prefix, fileInfo.Name())
			prefix += "│   "
		}
	case Leaf:
		fmt.Printf("%s├── %s\n", prefix, fileInfo.Name())
		prefix += "    "
	case FinalLeaf:
		fmt.Printf("%s└── %s\n", prefix, fileInfo.Name())
		prefix += "    "
	default:
		// Root node is the only node that should not have a prefix
		if node != Root {
			return fmt.Errorf("unknown node type: %v", node)
		} else {
			fmt.Printf("%s%s\n", prefix, fileInfo.Name())
		}
	}

	// If it's not a directory, return
	if !fileInfo.IsDir() {
		return nil
	}

	// Read directory contents
	entries, err := os.ReadDir(root)
	if err != nil {
		return err
	}

	// Filter out hidden files and sort
	var visibleEntries []os.DirEntry
	for _, entry := range entries {
		if !strings.HasPrefix(entry.Name(), ".") {
			visibleEntries = append(visibleEntries, entry)
		}
	}

	// Recursively print each entry
	for i, entry := range visibleEntries {
		isLastEntry := i == len(visibleEntries)-1
		err := printTree(filepath.Join(root, entry.Name()), prefix, isLastEntry)
		if err != nil {
			return err
		}
	}

	return nil
}
