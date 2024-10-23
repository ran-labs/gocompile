package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

/**
 * Check if the directive type matches
 * @param {string} directiveType - The directive type
 * @param {string} line - The line of the file
 * @return {bool} - If the directive type matches
 */
func isMatchingDirectiveType(
	directiveType string,
	component string,
) bool {
	if strings.Contains(component, "OF") {
		return strings.Contains(component, directiveType) || strings.Contains(component, "*")
	}
	return true
}

/**
 * Extract the exclusive component, and parent component string
 * @param {string} fileContentString - The file content string
 * @return {string} - The exclusive component string
 * @return {string} - The parent component string prefix
 * @return {string} - The parent component string suffix
 */
func removeInvalidExclusiveComponent(
	fileContentString string,
	directiveType string,
) string {
	newFileContentString := fileContentString
	var exclusiveComponentRegex = regexp.MustCompile(`(?s)<EXCLUSIVE(.*?)<\/EXCLUSIVE>`)
	var exclusiveComponentStrings = exclusiveComponentRegex.FindAllStringIndex(fileContentString, -1)
	for _, exclusiveComponent := range exclusiveComponentStrings {
		exclusiveComponent := fileContentString[exclusiveComponent[0]:exclusiveComponent[1]]
		if !isMatchingDirectiveType(directiveType, exclusiveComponent) {
			newFileContentString = strings.ReplaceAll(newFileContentString, exclusiveComponent, "")
		}
	}
	return newFileContentString
}

func parseJavascriptFile(path string, directiveType string) string {
	fileContent, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}
	fileContentString := string(fileContent)
	if !strings.Contains(fileContentString, "<EXCLUSIVE") {
		return string(fileContent)
	}
	fileContent = []byte(removeInvalidExclusiveComponent(fileContentString, directiveType))
	fmt.Println("File content:", string(fileContent))
	return string(fileContent)
}

/**
 * Check if the file is a javascript file
 * @param {string} path - The path of the file
 * @return {bool} - If the file is a javascript file
 */
func fileContainsUiComponents(path string) bool {
	fileExtensions := []string{".jsx", ".tsx", ".astro", ".svelte", ".vue"}
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
			if fileContainsUiComponents(path) {
				var fileContent = parseJavascriptFile(path, directiveType)
				if fileContent == "" {
					return nil
				}
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
	ignored_paths := []string{"node_modules", "build-target", "src-tauri", ".git", "gocomploy"}
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
