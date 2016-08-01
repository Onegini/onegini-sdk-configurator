package util

import (
	"os"
	"path"

	"fmt"

	"os/exec"

	"gitlab.onegini.com/mobile-platform/onegini-sdk-configurator/data"
	"strings"
)

var (
	removeFileScriptPath string
	addFileScriptPath string
)

func init() {
	tempPath := path.Join(os.TempDir(), "onegini-sdk-configurator")
	removeFileScriptPath = path.Join(tempPath, "lib", "removeFileFromXcodeProject.rb")
	addFileScriptPath = path.Join(tempPath, "lib", "addFileToXcodeProject.rb")

	if err := data.RestoreAsset(tempPath, "lib/removeFileFromXcodeProject.rb"); err != nil {
		os.Stderr.WriteString(fmt.Sprintf("ERROR: Could not restore required asset: %v\n", err))
		os.Exit(1)
	}

	if err := data.RestoreAsset(tempPath, "lib/addFileToXcodeProject.rb"); err != nil {
		os.Stderr.WriteString(fmt.Sprintf("ERROR: Could not restore required asset: %v\n", err))
		os.Exit(1)
	}
}

func removeFileFromXcodeProj(filepath string, xcodeProjPath string, group string) {
	ruby := checkForRuby()
	checkForXcodeprojGem()

	cmd := exec.Command(
		ruby,
		removeFileScriptPath,
		xcodeProjPath,
		filepath,
		group,
	)

	startCmd(cmd)
}

func addFileToXcodeProj(filePath string, xcodeProjPath string, appName string, group string) {
	ruby := checkForRuby();
	checkForXcodeprojGem();

	cmd := exec.Command(
		ruby,
		addFileScriptPath,
		xcodeProjPath,
		filePath,
		appName,
		group,
	)

	startCmd(cmd)
}

func checkForRuby() (ruby string) {
	ruby, lookErr := exec.LookPath("ruby")
	if lookErr != nil {
		os.Stderr.WriteString(fmt.Sprintln("ERROR: Could not find ruby executable in $PATH"))
		os.Exit(1)
	}

	return ruby
}

func checkForXcodeprojGem() {
	cmd := exec.Command("gem", "list")

	result, err := cmd.CombinedOutput()
	if (err != nil) {
		os.Stderr.WriteString(fmt.Sprintf("ERROR: Cannot find the xcodeproj gem: %v\n", err.Error()))
		os.Exit(1)
	}

	if !strings.Contains(string(result), "xcodeproj") {
		os.Stderr.WriteString(fmt.Sprintln("ERROR: The required gem 'xcodeproj' is not installed. Install it using: 'gem install xcodeproj'\n In case you are using the OS X System ruby use: 'sudo gem install xcodeproj'"))
		os.Exit(1)
	}
}

func startCmd(cmd *exec.Cmd) {
	outByte, err := cmd.CombinedOutput()

	if len(outByte) > 0 {
		fmt.Printf("%v\n", string(outByte))
	}
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("ERROR: Could not modify Xcode project: %v\n", err.Error()))
		os.Exit(1)
	}
}
