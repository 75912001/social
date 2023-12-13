package util

import "runtime"

// IsWindows win
func IsWindows() bool {
	return `windows` == runtime.GOOS
}

// IsLinux linux
func IsLinux() bool {
	return `linux` == runtime.GOOS
}

// RunMode 运行模式
type RunMode int

const (
	RunModeRelease RunMode = 0 //release 模式
	RunModeDebug   RunMode = 1 //debug 模式
)

// GRunMode 运行模式
var GRunMode = RunModeRelease

// IsDebug 是否为调试模式
func IsDebug() bool {
	return GRunMode == RunModeDebug
}

// IsRelease 是否为发行模式
func IsRelease() bool {
	return GRunMode == RunModeRelease
}
