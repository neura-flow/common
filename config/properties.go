package config

import (
	"encoding/json"
	"strings"

	"github.com/neura-flow/common/util"
)

// LoadProperties 加载 kv 格式的配置, 该配置是前端页面保存的组件属性配置格式，属性名格式：下划线
// 解析时，把下划线转化为 '.', 比如: a_b_c 转化为 a.b.c, 然后保存到 viper, 再 unmarshal
func LoadProperties(kvm map[string]interface{}, target any) error {
	c := newPropertiesToMap("_")
	bytes := []byte(util.ToJson(c.Do(kvm)))
	if err := json.Unmarshal(bytes, target); err != nil {
		return err
	}
	return nil
}

type propertiesToMap struct {
	keyDelimiter string
}

func newPropertiesToMap(keyDelimiter string) *propertiesToMap {
	if keyDelimiter == "" {
		keyDelimiter = "."
	}
	return &propertiesToMap{
		keyDelimiter: keyDelimiter,
	}
}

func (p *propertiesToMap) Do(props map[string]any) map[string]any {
	output := make(map[string]any)
	for k, v := range props {
		path := strings.Split(k, p.keyDelimiter)
		lastKey := path[len(path)-1]
		deepestMap := p.deepSearch(output, path[0:len(path)-1])
		deepestMap[lastKey] = v
	}
	return output
}

func (p *propertiesToMap) deepSearch(m map[string]any, path []string) map[string]any {
	for _, k := range path {
		m2, ok := m[k]
		if !ok {
			// intermediate key does not exist => create it and continue from there
			m3 := make(map[string]any)
			m[k] = m3
			m = m3
			continue
		}
		m3, ok := m2.(map[string]any)
		if !ok {
			// intermediate key is a value => replace with a new map
			m3 = make(map[string]any)
			m[k] = m3
		}
		// continue search from here
		m = m3
	}
	return m
}
