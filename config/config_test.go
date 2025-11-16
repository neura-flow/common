package config

import (
	"fmt"
	"testing"

	"github.com/neura-flow/common/util"
	"github.com/spf13/viper"
)

func TestConfig(t *testing.T) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("json")
	v.Set("database.host", "localhost")
	v.Set("database.port", 3306)
	v.Set("database.userName", "admin")
	v.Set("database.password", "123456")
	v.Set("database.extra.k1", "1")
	v.Set("database.extra.k2", "2")
	v.Set("database.list", "2,2,3,4")

	var cfg *config
	if err := v.Unmarshal(&cfg); err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%v\n", util.ToJson(cfg))
}

func TestLoadProperties(t *testing.T) {
	kvm := make(map[string]interface{})
	kvm["database_host"] = "localhost"
	kvm["database_port"] = 3306
	kvm["database_userName"] = "admin"
	kvm["database_password"] = "123456"
	kvm["database_extra.k1"] = "1"
	kvm["database_extra.k2"] = "2"
	kvm["database_list"] = []string{"2", "2", "3", "4"}
	var cfg = &config{}
	if err := LoadProperties(kvm, cfg); err != nil {
		t.Fatal(err)
	} else {
		fmt.Printf("%v\n", util.ToJson(cfg))
	}
}

func TestPropertiesToMap(t *testing.T) {
	kvm := make(map[string]interface{})
	kvm["database.host"] = "localhost"
	kvm["database.port"] = 3306
	kvm["database.userName"] = "admin"
	kvm["database.password"] = "123456"
	kvm["database.extra.k1"] = "1"
	kvm["database.extra.k2"] = "2"
	kvm["database.list"] = []string{"2", "2", "3", "4"}

	converter := newPropertiesToMap(".")
	result := converter.Do(kvm)
	fmt.Printf("%v\n", util.ToJson(result))
}

type config struct {
	Database *Database `json:"database"`
}
type Database struct {
	Host     string                 `json:"host"`
	Port     int                    `json:"port"`
	UserName string                 `json:"userName"`
	Password string                 `json:"password"`
	Extra    map[string]interface{} `json:"extra"`
	List     []string               `json:"list"`
}
