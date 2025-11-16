package metadata

import (
	"context"
	"fmt"
	"sync"
)

const (
	KeyZone       = "zone"
	KeyBiz        = "biz"
	KeyService    = "service"
	KeyNode       = "node"
	KeyComponent  = "component"
	KeyCode       = "code"
	KeyTopic      = "topic"
	KeyQueue      = "queue"
	KeyGroup      = "group"
	KeyPartition  = "partition"
	KeyMessageKey = "messageKey"
	KeyRouteKey   = "routeKey"
	KeyTimestamp  = "timestamp"
	KeyId         = "id"
	KeyOffset     = "offset"
	KeyLatency    = "latency"
	KeySize       = "size"
	KeyTraceId    = "traceId"
)

var Global Metadata = &metadata{}

type KV interface {
	Key() string
	Value() interface{}
	String() string
}

type Metadata interface {
	Value(key string) interface{}
	Set(key string, val interface{})
	Range(f func(kv KV))
	List() []KV
}

type kv struct {
	key   string
	value interface{}
}

func NewKV(key string, value interface{}) KV {
	return kv{
		key:   key,
		value: value,
	}
}

func (v kv) Key() string {
	return v.key
}

func (v kv) Value() interface{} {
	return v.value
}

func (v kv) String() string {
	return fmt.Sprintf("%s:%v", v.key, v.value)
}

type metadata struct {
	m sync.Map
}

// MergeMetadata 按顺序合并多个元数据对象，后面的覆盖前面的
func MergeMetadata(md ...Metadata) Metadata {
	var list []KV
	for _, m := range md {
		list = append(list, m.List()...)
	}
	return FromKVList(list...)
}

func New() Metadata {
	return &metadata{}
}

func FromMap(md map[string]interface{}) Metadata {
	m := &metadata{}
	for key, val := range md {
		m.m.Store(key, NewKV(key, val))
	}
	return m
}

func FromKVList(list ...KV) Metadata {
	md := &metadata{}
	for _, kv := range list {
		md.m.Store(kv.Key(), kv)
	}
	return md
}

type Key struct{}

func FromContext(ctx context.Context) Metadata {
	if md, ok := ctx.Value(Key{}).(Metadata); ok && md != nil {
		return md
	}
	return FromMap(map[string]interface{}{})
}

func ToContext(ctx context.Context, md Metadata) context.Context {
	return context.WithValue(ctx, Key{}, md)
}

func (md *metadata) Value(key string) interface{} {
	v, _ := md.m.Load(key)
	if kv, ok := v.(KV); ok {
		return kv.Value()
	}
	return nil
}

func (md *metadata) List() []KV {
	var list []KV
	md.Range(func(kv KV) {
		list = append(list, kv)
	})
	return list
}

func (md *metadata) Range(f func(KV)) {
	md.m.Range(func(key, value interface{}) bool {
		if v, ok := value.(KV); ok {
			f(v)
			return true
		}
		return false
	})
}

func (md *metadata) Set(key string, val interface{}) {
	md.m.Store(key, NewKV(key, val))
}

func Clone(md Metadata) Metadata {
	md1 := &metadata{}
	md.Range(func(k KV) {
		md1.Set(k.Key(), k.Value())
	})
	return md1
}
