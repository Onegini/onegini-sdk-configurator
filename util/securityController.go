package util

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
)

func ReadCordovaSecurityPreferences(config *Config) (rootDetection bool, debugDetection bool) {
	rootDetectionSet := false
	debugDetectionSet := false

	for _, pref := range config.Cordova.Preferences {
		if pref.Name == "OneginiRootDetectionEnabled" {
			rootDetectionSet = true
			var err error
			rootDetection, err = strconv.ParseBool(pref.Value)
			if err != nil {
				os.Stderr.WriteString(fmt.Sprintf("ERROR: could not parse 'OneginiRootDetectionEnabled' preference: %v\n", err.Error()))
				os.Exit(1)
			}
		}
		if pref.Name == "OneginiDebugDetectionEnabled" {
			debugDetectionSet = true
			var err error
			debugDetection, err = strconv.ParseBool(pref.Value)
			if err != nil {
				os.Stderr.WriteString(fmt.Sprintf("ERROR: could not parse 'OneginiDebugDetectionEnabled' preference: %v\n", err.Error()))
				os.Exit(1)
			}
		}
	}

	if !rootDetectionSet {
		rootDetection = true
	}
	if !debugDetectionSet {
		debugDetection = true
	}
	return
}

func WriteAndroidSecurityController(appDir string, config *Config, debugDetection bool, rootDetection bool) {
	fileContents := `package %s;

@SuppressWarnings({ "unused", "WeakerAccess" })
public final class SecurityController {
  public static final boolean debugDetection = %s;
  public static final boolean rootDetection = %s;
}`
	packageId := getPackageIdentifierFromConfig(config)
	fileContents = fmt.Sprintf(fileContents, packageId, strconv.FormatBool(debugDetection), strconv.FormatBool(rootDetection))
	storePath := getAndroidSecurityControllerPath(appDir, config)

	if rootDetection && debugDetection {
		os.Remove(storePath)
	} else {
		if err := ioutil.WriteFile(storePath, []byte(fileContents), os.ModePerm); err != nil {
			log.Fatal("WARNING! Could not update security controller. This might be dangerous!")
		}
	}
}

func WriteIOSSecurityController(appDir string, appName string, config *Config, debugDetection bool, rootDetection bool) {
	group := "Configuration"
	headerContents := `#import <Foundation/Foundation.h>

@interface SecurityController : NSObject
+ (bool)rootDetection;
+ (bool)debugDetection;
@end
`

	modelContents := `#import "SecurityController.h"

@implementation SecurityController
+(bool)rootDetection{
    return %s;
}
+(bool)debugDetection{
    return %s;
}
@end
`
	var (
		sDebugDetection string
		sRootDetection  string
	)

	if debugDetection {
		sDebugDetection = "YES"
	} else {
		sDebugDetection = "NO"
	}

	if rootDetection {
		sRootDetection = "YES"
	} else {
		sRootDetection = "NO"
	}

	modelContents = fmt.Sprintf(modelContents, sRootDetection, sDebugDetection)
	xcodeProjPath := getIosXcodeProjPath(appDir, appName, config)
	configModelPath := getIosConfigModelPath(appDir, appName, config)

	headerStorePath := path.Join(configModelPath, "SecurityController.h")
	modelStorePath := path.Join(configModelPath, "SecurityController.m")

	if rootDetection && debugDetection {
		removeFileFromXcodeProj(headerStorePath, xcodeProjPath, group)
		removeFileFromXcodeProj(modelStorePath, xcodeProjPath, group)
		os.Remove(headerStorePath)
		os.Remove(modelStorePath)
	} else {
		ioutil.WriteFile(headerStorePath, []byte(headerContents), os.ModePerm)
		ioutil.WriteFile(modelStorePath, []byte(modelContents), os.ModePerm)
		addFileToXcodeProj(headerStorePath, xcodeProjPath, appName, group)
		addFileToXcodeProj(modelStorePath, xcodeProjPath, appName, group)
	}
}
