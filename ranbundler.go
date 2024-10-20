package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

/**
 * Find the end line index of the exclusive component
 * @param {string} line - The line of the file
 * @param {int} lineIndex - The index of the line
 * @param {bufio.Reader} reader - The reader of the file
 * @param {[]string} fileContentLines - The content of the file
 * @return {int} - The end line index of the exclusive component
 */
func findExclusiveComponentEndLineIndex(
	line string,
	lineIndex int,
	reader *bufio.Reader,
	fileContentLines []string,
) int {
	if strings.Contains(line, "</EXCLUSIVE") {
		return lineIndex
	} else {
		var componentEndLineIndex int = lineIndex
		for {
			componentEndLineIndex++
			nextLine, err := reader.ReadString('\n')
			if err != nil {
				break
			}
			var currentLine string = string(nextLine)
			fileContentLines = append(fileContentLines, currentLine)
			if strings.Contains(currentLine, "</EXCLUSIVE") {
				return componentEndLineIndex
			}
		}
	}
	return -1
}

/**
 * Parse the javascript file
 * @param {string} path - The path of the file
 */
func parseJavascriptFile(path string, directiveType string) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var lineIndex int = 0
	// Read and print lines
	var fileContentLines = []string{}
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		fileContentLines = append(fileContentLines, line)
		// Exclusive component can be on the same line or on following line
		var lineString string = string(line)
		if strings.Contains(lineString, "<EXCLUSIVE") && !strings.Contains(lineString, directiveType) {
			var componentEndLineIndex int = findExclusiveComponentEndLineIndex(
				lineString,
				lineIndex,
				reader,
				fileContentLines,
			)
			if lineIndex > componentEndLineIndex {
			}
			var exclusiveComponent []string = fileContentLines[lineIndex : componentEndLineIndex+1]
			fmt.Println("Content of Exclusive Component:", exclusiveComponent)
		}
		lineIndex++
	}
}

/**
 * Check if the file is a javascript file
 * @param {string} path - The path of the file
 * @return {bool} - If the file is a javascript file
 */
func isJavascriptFile(path string) bool {
	fileExtensions := []string{".js", ".jsx", ".ts", ".tsx", ".vue", ".svelte"}
	ext := filepath.Ext(path)
	for _, fileExt := range fileExtensions {
		if ext == fileExt {
			return true
		}
	}
	return false
}

/**
 * Walk through the file path
 * @param {string} path - The path of the file
 * @param {string} directiveType - The directive type
 */
func walkFilePath(path string, directiveType string) {
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return nil
		}
		if isJavascriptFile(path) {
			parseJavascriptFile(path, directiveType)
		}
		return nil
	})
}

func main() {
	var path string = "/home/sanner/Coding/RAN/ran-app-native/src"
	var directiveType string = "web"
	fmt.Println("Path:", path)
	walkFilePath(path, directiveType)
}
