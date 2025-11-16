package metadata

import (
	"context"
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetadata(t *testing.T) {
	md := FromMap(map[string]interface{}{"foo": "bar"})
	assert.Equal(t, md.Value("foo"), "bar")
	ctx := ToContext(context.Background(), md)
	md1 := FromContext(ctx)
	assert.Equal(t, md1.Value("foo"), "bar")
	list := md.List()
	assert.Equal(t, list[0].Key(), "foo")
	assert.Equal(t, list[0].Value(), "bar")
	md2 := FromKVList(list...)
	assert.Equal(t, md2.Value("foo"), "bar")
	md2 = FromMap(map[string]interface{}{"x": "y"})
	md3 := MergeMetadata(md, md2)
	assert.Equal(t, md3.Value("foo"), "bar")
	assert.Equal(t, md3.Value("x"), "y")
}

func ExampleMetadata() {
	//从 map 创建 Metadata
	md1 := FromMap(map[string]interface{}{"x": "1"})
	//从 KV 列表创建 Metadata
	md2 := FromKVList(NewKV("y", "2"))
	//合并多个 Metadata
	md3 := MergeMetadata(md1, md2)
	//将 Metadata 放到 context
	ctx := ToContext(context.Background(), md3)
	//从 context 获取 Metadata
	md := FromContext(ctx)
	//获取 KV 列表
	list := md.List()
	sort.Slice(list, func(i, j int) bool {
		return list[j].Key() > list[i].Key()
	})
	fmt.Println(list)
	//遍历，转成 map
	var m = make(map[string]interface{})
	md.Range(func(kv KV) {
		m[kv.Key()] = kv.Value()
	})
	fmt.Printf("map[%s:%v %s:%v]\n", "x", m["x"], "y", m["y"])
	//获取单个值
	x := md.Value("x")
	y := md.Value("y")
	fmt.Println(x, y)

	// Output:
	//[x:1 y:2]
	//map[x:1 y:2]
	//1 2
}

func TestKV(t *testing.T) {
	kv := NewKV("foo", "bar")
	assert.NotNil(t, kv)
	assert.Equal(t, kv.Key(), "foo")
	assert.Equal(t, kv.Value(), "bar")
}
