//Copyright 2017 Onegini B.V.
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

package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func WriteAndroidAppScheme(config *Config) {
	if config.ConfigureForNativeScript {
		return
	}

	manifestPath := config.getAndroidManifestPath()
	manifest, err := ioutil.ReadFile(manifestPath)
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("ERROR: Cannot read the Android Manifest: %v.\n", err))
		os.Exit(1)
	}

	scheme := strings.Split(config.Options.RedirectUrl, "://")[0]

	if config.ConfigureForCordova {
		newRegexp := regexp.MustCompile(`(?s)<intent-filter android:label="OneginiRedirectionIntent" android:name="OneginiRedirectionIntent">.*android:scheme="([^"]*)".*</intent-filter>`)
		oldRegexp := regexp.MustCompile(`(?s)<activity\s+.*android:name="MainActivity".*>.*<intent-filter>.*android:scheme="([^"]*)".*</intent-filter>.*</activity>`)

		schemeRegexp := regexp.MustCompile(`android:scheme="[^"]*"`)

		if newRegexp.Match(manifest) {
			os.Stderr.WriteString(fmt.Sprintf("NEW MATCH FOUND\n"))
			manifest = newRegexp.ReplaceAllFunc(manifest, func(input []byte) (output []byte) {
				output = schemeRegexp.ReplaceAll(input, []byte("android:scheme=\""+scheme+"\""))
				return
			})
		} else {
			os.Stderr.WriteString(fmt.Sprintf("OLD MATCH FOUND\n"))
			manifest = oldRegexp.ReplaceAllFunc(manifest, func(input []byte) (output []byte) {
				output = schemeRegexp.ReplaceAll(input, []byte("android:scheme=\""+scheme+"\""))
				return
			})
		}
		ioutil.WriteFile(manifestPath, manifest, os.ModePerm)
	}
}
