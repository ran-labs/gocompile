package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

/**
 * Removes the exclusive component from the file content
 * @param {string} line - The line of the file
 * @param {int} lineIndex - The index of the line
 * @param {bufio.Reader} reader - The reader of the file
 * @param {[]string} fileContentLines - The content of the file
 * @return {int} - The end line index of the exclusive component
 */
func removeExclusiveComponent(
	line string,
	lineIndex int,
	reader *bufio.Reader,
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
			// fileContentLines = append(fileContentLines, currentLine)
			if strings.Contains(currentLine, "</EXCLUSIVE") {
				return componentEndLineIndex
			}
		}
	}
	return -1
}

/**
 * Check if the directive type matches
 * @param {string} directiveType - The directive type
 * @param {string} line - The line of the file
 * @return {bool} - If the directive type matches
 */
func isMatchingDirectiveType(
	directiveType string,
	line string,
) bool {
	return strings.Contains(line, directiveType) || strings.Contains(line, "*")
}

func findExclusiveComponentParams(
	directiveType string,
	lineIndex int,
	reader *bufio.Reader,
	fileContentLines []string,
) ([]string, bool) {
	tempFileContentLines := fileContentLines
	for {
		lineIndex++
		nextLine, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		currentLine := string(nextLine)
		tempFileContentLines = append(tempFileContentLines, currentLine)
		if strings.Contains(currentLine, "OF") && !isMatchingDirectiveType(directiveType, currentLine) {
			return fileContentLines, false
		} // check if the directive type matches
		if strings.Contains(currentLine, ">") {
			fileContentLines = tempFileContentLines
			return fileContentLines, true
		}
	}
	// Case occurs when exclusive component is not closed
	// TODO: Add error handling (Throw an error and stop the program)
	return nil, false
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
		var useExclusiveComponent bool
		fileContentLines, useExclusiveComponent = findExclusiveComponentParams(directiveType, lineIndex, reader, fileContentLines)
		if !useExclusiveComponent {
			removeExclusiveComponent(
				lineString,
				lineIndex,
				reader,
			)
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

func parsePathDifference(srcPath string, dstDir string) string {
	rel, err := filepath.Rel(srcPath, dstDir)
	if err != nil {
		fmt.Println(err)
	}
	return rel
}

/**
 * Walk through the file path
 * @param {string} path - The path of the file
 * @param {string} directiveType - The directive type
 */
func walkFilePath(srcDir string, directiveType string, dstDir string) {
	filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			var pathDifference string = parsePathDifference(srcDir, path)
			newDir := dstDir + pathDifference
			fmt.Println("Directory:", newDir)
			os.MkdirAll(dstDir+pathDifference, 0777)
		} else {
			// fmt.Println("File:", path)
		}
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
	fmt.Println("Path:", path)
	deviceTypes := []string{"mobile"}

	var buildPath string = path + "/build-target"
	os.Mkdir(buildPath, 0777)
	for _, deviceType := range deviceTypes {
		var deviceBuildPath string = buildPath + "/" + deviceType
		os.Mkdir(deviceBuildPath, 0777)
		walkFilePath(path, deviceType, deviceBuildPath)
	}
}
