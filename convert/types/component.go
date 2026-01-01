package types

import (
	"fmt"
	"sync"

	"github.com/neura-flow/common/log"
	"github.com/neura-flow/common/named"
)

type Component interface {
	named.Named
	Builder
}

type component struct {
	named.Named
	Builder
}

func NewComponent(named named.Named, b Builder) Component {
	return &component{
		Named:   named,
		Builder: b,
	}
}

type Instance struct {
	named.Named
}

func Register(comp named.Named) {
	cm.Register(comp)
}

var cm = &Manager{}

func GlobalManager() *Manager {
	return cm
}

type Manager struct {
	named.Named
	m sync.Map
}

func (cm *Manager) Load(name named.Named) (Component, bool) {
	bld, ok := cm.load(name)
	if !ok {
		return nil, false
	}
	c, ok := bld.(Component)
	if ok {
		return c, true
	}
	return nil, false
}

func (cm *Manager) NewInstance(name named.Name, bundle *Bundle) (named.Named, error) {
	c, ok := cm.Load(name)
	if !ok {
		return nil, fmt.Errorf("builder: %s not found", name)
	}
	return c.Build(name, bundle)
}

func (cm *Manager) Has(name named.Name) bool {
	_, ok := cm.m.Load(name.Name())
	return ok
}

func (cm *Manager) load(name named.Named) (bld named.Named, ok bool) {
	log.DefaultLogger().Debugf("Load " + name.Name())
	val, ok := cm.m.Load(name.Name())
	if !ok {
		return
	}
	bld, ok = val.(named.Named)
	return
}

func (cm *Manager) Register(builders ...named.Named) {
	for _, bld := range builders {
		if bld == nil {
			continue
		}
		log.DefaultLogger().Infof("Register builder: %s", bld.Name())
		if _, ok := cm.m.LoadOrStore(bld.Name(), bld); ok {
			log.DefaultLogger().Panicf("Register failed, duplicate name: %s", bld.Name())
		}
	}
}

func (cm *Manager) Keys() []string {
	keys := make([]string, 0)
	cm.m.Range(func(key, value interface{}) bool {
		keys = append(keys, key.(string))
		return true
	})
	return keys
}

func (cm *Manager) Range(f func(bld named.Named)) {
	cm.m.Range(func(key, value interface{}) bool {
		f(value.(named.Named))
		return true
	})
}
