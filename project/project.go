package project

import (
	"os"
	"path/filepath"
	"strings"
)

// env-tag
const (
	DevEnv        = "Dev"
	TestEnv       = "Test"
	ProdEnv       = "Product"
	PrereleaseEnv = "Prerelease"
)

// IsProduct 当前环境是不是生产环境
func IsProduct() bool {
	env := os.Getenv("RunTimeEnv")
	if strings.EqualFold(ProdEnv, env) {
		return true
	}
	return false
}

// IsDev 当前环境是不是开发环境
func IsDev() bool {
	env := os.Getenv("RunTimeEnv")
	if strings.EqualFold(DevEnv, env) {
		return true
	}
	return false
}

// IsTest 当前环境是不是测试环境
func IsTest() bool {
	env := os.Getenv("RunTimeEnv")
	if strings.EqualFold(TestEnv, env) {
		return true
	}
	return false
}

// IsPrerelease 当前环境是不是预发布环境
func IsPrerelease() bool {
	env := os.Getenv("RunTimeEnv")
	if strings.EqualFold(PrereleaseEnv, env) {
		return true
	}
	return false
}

// GetRunTime 获取当前系统环境
func GetRunTime() string {
	env := os.Getenv("RunTimeEnv")
	if len(env) != 0 {
		return env
	}
	return DevEnv
}

// GetProjectName 获取系统名称
func GetProjectName() string {
	name := os.Getenv("AppName")
	if len(name) != 0 {
		return name
	} else {
		executable, _ := os.Executable()
		dir := filepath.Dir(executable)
		return filepath.Base(dir)
	}
	return "UnKnowSystem"
}

// GetRootDir
// 获取上层业务项目的根路径
// 获取项目根路径的过程中，必须至少满足以下其中一个条件
// 1. 编译后的二制程序位于 项目根目录下
// 2. 编译后的二制程序位于 项目根目录下的bin目录下
// 3. 执行当前测试时，当前路径必须在项目根目录下或者位于 src 以内的子目录内(包括src目录下)
func GetRootDir() string {
	//定位当前运行中的可执行文件的绝对路径，这里包括正式的可执行文件可用于单元测试的可执行文件
	filename, _ := filepath.Abs(os.Args[0])
	binfile := filepath.Base(filename)

	if strings.Contains(binfile, ".test") || strings.HasPrefix(binfile, "___") {
		//认为当前的执行环境为单元测试，此时考虑使用 当前路径 来定位根目录
		//这里的curPath为当前cd进入的路径
		curPath, err := os.Getwd()
		if err != nil {
			return ""
		}
		return _findRootDir(curPath)
	} else if strings.Contains(binfile, "debug") {
		// delve debug 文件, filename所在的文件夹为模块目录，在获取本地private.env和configs下的schema时会出错，需要向上获取根目录
		return _findRootDir(filepath.Dir(filename))
	} else {
		//认为当前的执行环境正常编译后的二进制程序
		rootDir := filepath.Dir(filename)
		if strings.Contains(filename, "bin") {
			rootDir = filepath.Dir(filepath.Dir(filename))
		}
		return rootDir
	}
}

func _findRootDir(curPath string) string {
	//考虑使用go.mod的位置来确定项目根目录
	//以当前目录为基准，向上遍历，找到存在go.mod的目录，即认为是项目根目录
	assertRoot := curPath
	for {
		if FileExists(assertRoot + string(filepath.Separator) + "go.mod") {
			//表示命中根目录(仅适用于go modules项目)
			return assertRoot
		}
		parentPath := filepath.Dir(assertRoot)
		if assertRoot == parentPath {
			//已经到根目录下，仍没有定位到项目根目录，返回最初的当前路径
			return curPath
		} else {
			assertRoot = parentPath
		}
	}
}

// FileExists 判断一个文件或目录是否存在
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil && !os.IsExist(err) {
		return false
	}
	return true
}
