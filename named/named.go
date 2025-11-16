package named

import "strings"

const DefaultName = Name("default")

type Named interface {
	//Name 获取完整名称
	Name() string
	//ShortName 获取短名称，名称的最后部分
	ShortName() Name
	//Namespace 获取命名空间，不包含短名称的部分
	Namespace() Namespace
}

type Name string

func (n Name) Name() string {
	return string(n)
}

// Namespace 获取命名空间，不包含短名称的部分
func (n Name) Namespace() Namespace {
	idx := strings.LastIndexByte(n.Name(), '.')
	if idx <= 0 {
		return ""
	}
	return Namespace(n.Name()[:idx])
}

// ShortName 获取短名称，名称的最后部分
func (n Name) ShortName() Name {
	names := strings.Split(n.Name(), ".")
	return Name(names[len(names)-1])
}

// IsDefaultName 判断名称是不是默认名称
func IsDefaultName(name Named) bool {
	lastName := Name(name.Name()).ShortName()
	return lastName.Name() == DefaultName.Name() || name.Name() == ""
}

// IsShortName 判断名称是不是短名称
func IsShortName(name Named) bool {
	return Name(name.Name()).ShortName().Name() == name.Name()
}

// Namespace 命名空间
type Namespace string

// Name 获取名称
func (n Namespace) Name() string {
	return string(n)
}

// Namespace 获取命名空间，不包含短名称的部分
func (n Namespace) Namespace() Namespace {
	return Name(n).Namespace()
}

// ShortName 获取短名称，名称的最后部分
func (n Namespace) ShortName() Name {
	return Name(n).ShortName()
}

// Join 将一个名称连接到命名空间
func (n Namespace) Join(name Name) Name {
	if name == "" {
		name = DefaultName
	}
	if n == "" {
		return name
	}
	return Name(n.Name() + "." + name.Name())
}

func (n Namespace) JoinNS(ns Namespace) Namespace {
	return Namespace(n.Name() + "." + ns.Name())
}
