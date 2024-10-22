package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
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
func parseJavascriptFile(path string, directiveType string) string {
	// TODO: Deal with comments
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return err.Error()
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
	return strings.Join(fileContentLines, "")
}

/**
 * Check if the file is a javascript file
 * @param {string} path - The path of the file
 * @return {bool} - If the file is a javascript file
 */
func fileContainsUiComponents(path string) bool {
	fileExtensions := []string{".jsx", ".tsx", ".astro", ".svelte"}
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
func walkFilePath(srcDir string, directiveBuildPath string, directiveType string, ignoredPaths []string) {
	filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		newPath := srcDir + strings.ReplaceAll(path, srcDir, directiveBuildPath)
		if pathIsIgnored(path, ignoredPaths) {
			// fmt.Println("Ignored path:", path)
			return filepath.SkipDir
		}
		if info.IsDir() {
			os.MkdirAll(newPath, 0777)
		} else {
			// TODO: consider returning conditional stating that the file has not Exclusive component
			if fileContainsUiComponents(path) {
				var fileContent = parseJavascriptFile(path, directiveType)
				file, err := os.Create(newPath)
				if err != nil {
					fmt.Println(err)
					return err
				}
				defer file.Close()
				_, err = file.WriteString(fileContent)
				if err != nil {
					fmt.Println(err)
					return err
				}
			} else {
				srcFile, err := os.Open(path)
				if err != nil {
					fmt.Println(err)
					return err
				}
				defer srcFile.Close()

				dstFile, err := os.Create(newPath)
				if err != nil {
					fmt.Println(err)

					return err
				}
				defer dstFile.Close()

				_, err = io.Copy(dstFile, srcFile)
				if err != nil {
					fmt.Println(err)

					return err
				}
			}
		}
		if err != nil {
			fmt.Println(err)

			return nil
		}

		return nil
	})
}

func pathIsIgnored(path string, ignoredPaths []string) bool {
	var currentDirArray []string = strings.Split(path, "/")
	if len(currentDirArray) == 0 {
		return false
	}
	var currentDir string = currentDirArray[len(currentDirArray)-1]
	for _, ignoredPath := range ignoredPaths {
		if currentDir == ignoredPath {
			return true
		}
	}
	return false
}

func modifyPlatformConfigurationFile(path string, deviceType string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Read the entire file
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	var fileString string = string(data)
	platformStringIndex := strings.Index(fileString, "PLATFORM")
	platformStringTerminationIndex := strings.Index(fileString[platformStringIndex:], "}")
	if platformStringIndex == -1 &&
		platformStringTerminationIndex == -1 &&
		platformStringTerminationIndex < platformStringIndex {
		return "", fmt.Errorf("PLATFORM not found in the file")
	}
	newPlatformString := fmt.Sprintf(
		`PLATFORM: { MODE: "%s", NAME: "%s", ID: "%s" }`, deviceType, deviceType, deviceType,
	)
	newData := strings.ReplaceAll(fileString, fileString[platformStringIndex:platformStringIndex+platformStringTerminationIndex+1], newPlatformString)
	return newData, nil
}

func main() {
	var path string = "/home/sanner/Coding/RAN/ran-app-native/"
	fmt.Println("Path:", path)
	deviceTypes := []string{"mobile", "web"}
	ignored_paths := []string{"node_modules", "build-target", ".git", "gocomploy"}
	var wg sync.WaitGroup
	for _, deviceType := range deviceTypes {
		wg.Add(1)
		go func(deviceType string) {
			defer wg.Done()
			var deviceBuildPath = "build-target/" + deviceType + "/"
			walkFilePath(path, deviceBuildPath, deviceType, ignored_paths)
		}(deviceType)
	}
	wg.Wait()
	modifyPlatformConfigurationFile("platform.ts", "mobile")
}
