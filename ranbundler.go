package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/pelletier/go-toml"
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

type Config struct {
	InputPath    string   `toml:"input_path" json:"input_path"`
	OutputDir    string   `toml:"output_dir" json:"output_dir"`
	IgnoredPaths []string `toml:"ignored_paths" json:"ignored_paths"`
	DeviceTypes  []string `toml:"device_types" json:"device_types"`
}

func parseTomlConfigurationFile() (Config, error) {
	file, err := os.Open("config.toml")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var config Config

	b, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	err = toml.Unmarshal(b, &config)
	if err != nil {
		panic(err)
	}
	return config, nil
}

func parseJsonConfigurationFile(jsonFilePath string) (Config, error) {
	jsonFile, err := os.Open(jsonFilePath)
	if err != nil {
		fmt.Printf("Error opening JSON file: %v\n", err)
		return Config{}, err
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Printf("Error reading JSON file: %v\n", err)
		return Config{}, err
	}
	var config Config
	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		fmt.Printf("Error unmarshalling JSON file: %v\n", err)
		return Config{}, err
	}
	return config, nil
}

/**
 * Get the Configuration from the configuration file
 * @param {string} configurationFile - The configuration file
 * @return {Config} - The configuration
 * @return {error} - The error
 */
func getConfigurationFile(configurationFile string) (Config, error) {
	if strings.HasSuffix(configurationFile, ".toml") {
		return parseTomlConfigurationFile()
	} else if strings.HasSuffix(configurationFile, ".json") {
		return parseJsonConfigurationFile(configurationFile)
	}
	return Config{}, fmt.Errorf("invalid configuration file")
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide the configuration file path")
		return
	}
	getConfigurationFile(os.Args[1])
	config, err := parseJsonConfigurationFile("config.json")
	if err != nil {
		fmt.Println(err)
	}
	path := config.InputPath
	outputDir := config.OutputDir
	var wg sync.WaitGroup
	for _, deviceType := range config.DeviceTypes {
		wg.Add(1)
		go func(deviceType string) {
			defer wg.Done()
			var deviceBuildPath = outputDir + "/" + deviceType + "/"
			walkFilePath(path, deviceBuildPath, deviceType, config.IgnoredPaths)
			modifyPlatformConfigurationFile("platform.ts", deviceType)
		}(deviceType)
		command := exec.Command("pnpm", "install")
		command.Dir = currentWorkingDirectory + "/" + outputDir + "/" + deviceType
		err := command.Run()
		if err != nil {
			fmt.Println(err)
		}
	}
	wg.Wait()
}
