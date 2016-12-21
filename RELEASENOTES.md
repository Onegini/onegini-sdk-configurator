# Release notes

## 3.1.1

### Bug fixes

* Fixed path resolving issues when using a relative path as `app-dir`
* Auto locate the Xcode project file
* Prevent the configurator from creating multiple Xcode references wen run multiple times

## 3.1.0

### Features

* The Max PIN failures property in the ConfigModel for Android is no longer a required property.

## 3.0.0

Please note that this release is only compatible with the following SDK versions:
* Android SDK 6.00.00 and higher
* iOS SDK 5.00.00 and higher

### Features

* Add support for Android SDK versions 6.00.00 and higher
* Add support for iOS SDK versions 5.00.00 and higher
* Add a version flag
 
### Bug fixes

* Fixed a bug that forced a specific Gradle project layout for Android

## 2.0.0

### Features

* Complete rebuild of the SDK configurator in go
* CLI api using flags

## 1.0.0

### Features

* Initial release