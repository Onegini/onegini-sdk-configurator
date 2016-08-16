//Copyright 2016 Onegini B.V.
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.

package cmd

import (
	"os"
	"path"

	"fmt"

	"github.com/spf13/cobra"
	"github.com/Onegini/onegini-sdk-configurator/util"
)

var iosCmd = &cobra.Command{
	Use:   "ios",
	Short: "Configure an iOS project",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		config := util.ParseConfig(tsConfigLocation, cmd)
		var appTarget string

		if isCordova {
			util.ParseCordovaConfig(appDir, config)
			rootDetection, debugDetection = util.ReadCordovaSecurityPreferences(config)
			appTarget = config.Cordova.AppName

			verifyCordovaIosPlatformInstalled()

		} else {
			appTarget = targetName
		}
		verifyAppTarget(appTarget, cmd)

		util.PrepareIosPaths(appDir, appTarget, config)
		util.WriteIOSConfigModel(appDir, appTarget, config)
		util.WriteIOSSecurityController(appDir, appTarget, config, debugDetection, rootDetection)
		util.ConfigureIOSCertificates(config, appTarget, appDir)

		util.PrintSuccessMessage(config, debugDetection, rootDetection)
	},
}

func verifyAppTarget(appTarget string, cmd *cobra.Command) {
	if (len(appTarget) == 0) {
		if isCordova {
			os.Stderr.WriteString(fmt.Sprintln("ERROR: No application identifier found in your 'config.xml'. Please make sure that you have set one."))
			os.Exit(1)
		} else {
			fmt.Print("ERROR: No target name provided. Provide one using 'onegini-sdk-configurator ios -t <target-name>'\n\n")
			cmd.Help()
			os.Exit(1)
		}
	}
}

func verifyCordovaIosPlatformInstalled() {
	_, err := os.Stat(path.Join(appDir, "platforms", "ios"))
	if os.IsNotExist(err) {
		os.Stderr.WriteString(fmt.Sprintln("ERROR: Your project does not seem to have the iOS platform added. Please try `cordova platform add ios`"))
		os.Exit(1)
	}
}
