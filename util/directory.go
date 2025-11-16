package util

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/samuel/go-zookeeper/zk"
)

func CreateZkNode(path string, conn *zk.Conn) error {
	return CreateZkNodeWithData(path, nil, conn)
}

func CreateZkNodeWithData(path string, data []byte, conn *zk.Conn) error {
	exist, _, err := conn.Exists(path)
	if err != nil {
		return err
	}
	if !exist {
		for _, dir := range GetParentPathsWithRoot(path, false) {
			if exist, _, err = conn.Exists(dir); err != nil {
				return err
			} else if !exist {
				if _, err = conn.Create(dir, nil, 0, zk.WorldACL(zk.PermAll)); err != nil {
					return err
				}
			}
		}
		_, err = conn.Create(path, data, 0, zk.WorldACL(zk.PermAll))
		return err
	}
	return nil
}

// GetParentPaths 返回给定文件路径的所有父级路径列表
// 例如：输入 "/a/b/c/d.txt" 返回 ["/", "/a", "/a/b", "/a/b/c"]
func GetParentPaths(filePath string) []string {
	// 清理路径，移除多余的路径分隔符和相对路径符号
	cleanPath := filepath.Clean(filePath)

	// 如果是空路径或根路径，直接返回空切片
	if cleanPath == "" || cleanPath == "/" || cleanPath == "." || cleanPath == ".." {
		return []string{}
	}

	var parentPaths []string

	// 处理绝对路径的情况
	if filepath.IsAbs(cleanPath) {
		// 分割路径组件
		components := strings.Split(cleanPath, string(filepath.Separator))

		// 重建父级路径
		currentPath := ""
		for i, comp := range components {
			if comp == "" {
				// 根目录情况
				if i == 0 {
					currentPath = string(filepath.Separator)
					parentPaths = append(parentPaths, currentPath)
				}
				continue
			}

			if currentPath == string(filepath.Separator) {
				currentPath = currentPath + comp
			} else if currentPath == "" {
				currentPath = comp
			} else {
				currentPath = filepath.Join(currentPath, comp)
			}

			// 不包含最后一个组件（文件或最终目录）
			if i < len(components)-1 {
				parentPaths = append(parentPaths, currentPath)
			}
		}
	} else {
		// 处理相对路径
		dir := filepath.Dir(cleanPath)

		// 如果已经是当前目录，没有父级路径
		if dir == "." {
			return []string{}
		}

		// 递归构建父级路径列表
		for dir != "." && dir != "" {
			parentPaths = append([]string{dir}, parentPaths...) // 插入到开头
			dir = filepath.Dir(dir)
		}
	}

	return parentPaths
}

// GetParentPathsWithRoot 返回父级路径列表，可选择是否包含根路径
func GetParentPathsWithRoot(filePath string, includeRoot bool) []string {
	parents := GetParentPaths(filePath)

	if !includeRoot && len(parents) > 0 {
		// 移除根路径（第一个元素）
		if parents[0] == string(filepath.Separator) || parents[0] == "." {
			if len(parents) > 1 {
				return parents[1:]
			}
			return []string{}
		}
	}

	return parents
}

func JoinDir(args ...interface{}) string {
	path := ""
	for _, arg := range args {
		if arg != "" {
			p := ToString(arg)
			if strings.HasPrefix(p, "/") {
				path = path + p
			} else {
				path = path + "/" + p
			}
		}
	}
	return path
}

// IsExist 判断给定的文件路径是否存在
func IsExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// IsDir 判断给定的路径是否是文件夹
func IsDir(path string) bool {
	if stat, err := os.Stat(path); err == nil {
		return stat.IsDir()
	}
	return false
}

// CreateFileIfNecessary 给定的文件不存在则创建
func CreateFileIfNecessary(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		if file, err := os.Create(path); err == nil {
			file.Close()
		}
	}
	exist := IsExist(path)
	return exist
}

// MkdirIfNecessary 给定的目录不存在则创建
func MkdirIfNecessary(path string) error {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			// os.Chmod(path, 0777)
			return err
		}
	}
	return nil
}
