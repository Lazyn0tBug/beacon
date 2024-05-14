package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// FindProjectRoot 向上查找并返回包含go.mod文件的目录路径
func FindProjectRoot() (string, error) {
	cur_dir, err := os.Getwd()
	if err != nil {
		return "", nil
	}

	for {
		// 构建go.mod文件的完整路径
		goModPath := filepath.Join(cur_dir, "go.mod")
		// 检查文件是否存在
		if _, err := os.Stat(goModPath); err == nil {
			// 找到go.mod文件，返回当前目录
			return cur_dir, nil
		}
		// 如果已经到达根目录，则停止查找
		if cur_dir == filepath.Dir(cur_dir) {
			break
		}
		// 向上移动到父目录
		cur_dir = filepath.Dir(cur_dir)
	}
	// 未找到go.mod文件，返回错误
	return "", fmt.Errorf("could not find project root: go.mod not found")
}

// checkLoggingDir 函数用于检查日志目录是否存在，若不存在则创建该目录
// 返回日志目录的路径和错误信息
func checkConfigDir() (string, error) {
	root, err := FindProjectRoot()
	if err != nil {
		return ".", err
	}

	config_dir := filepath.Join(root, "config")

	if _, err := os.Stat(config_dir); os.IsNotExist(err) {
		return "", err // 返回空字符串和创建目录时遇到的错误
	}
	return config_dir, nil
}

// checkLoggingDir 函数用于检查日志目录是否存在，若不存在则创建该目录
// 返回日志目录的路径和错误信息
func checkLoggingDir() (string, error) {
	root, err := FindProjectRoot()
	if err != nil {
		return ".", err
	}

	logging_dir := filepath.Join(root, "logs")

	if _, err := os.Stat(logging_dir); os.IsNotExist(err) {
		// 如果目录不存在，则创建它
		if err := os.MkdirAll(logging_dir, 0755); err != nil {
			return "", err // 返回空字符串和创建目录时遇到的错误
		}
	}
	return logging_dir, nil
}

// getConfigFile 构建并返回特定配置文件的路径
// configFileName: 配置文件名，如 "logging.toml" 或 "database.toml"
// 返回配置文件的完整路径和可能发生的错误
func getConfigFile(configFileName string) (string, error) {
	configDir, err := checkConfigDir()
	if err != nil {
		// 如果检查配置目录出错，返回默认的配置文件名和错误
		return configFileName, err
	}

	// 构建完整的配置文件路径
	configPath := filepath.Join(configDir, configFileName)
	return configPath, nil
}

// 原来的 getLoggingConfigFile 和 getDBconfigFile 可以调用新的 getConfigFile 函数

func getLoggingConfigFile() (string, error) {
	return getConfigFile("logging.toml")
}

func getDBconfigFile() (string, error) {
	return getConfigFile("database.toml")
}
