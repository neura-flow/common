package config

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	vaultApi "github.com/hashicorp/vault/api"
	jsoniter "github.com/json-iterator/go"
	kjson "github.com/knadh/koanf/parsers/json"
	kyaml "github.com/knadh/koanf/parsers/yaml"
	kenv "github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/knadh/koanf/v2"
	"github.com/neura-flow/common/util"
)

const (
	DefaultAppName       = "datasix"
	DefaultCfgFile       = "./configs/config.yaml"
	DefaultAppPort       = 10001
	DefaultConsulAddress = "http://127.0.0.1:8500"
	DefaultVaultAddress  = "http://127.0.0.1:8200"
	DefaultConsulPrefix  = "config"
	DefaultVaultPrefix   = "secret"
	AppPort              = "app.port"
	AppFile              = "app.file"
	ConsulAddress        = "consul.address"
	ConsulPrefix         = "consul.prefix"
	ConsulEnabled        = "consul.enabled"
	ConsulToken          = "consul_token"
	VaultAddress         = "vault.address"
	VaultToken           = "vault.token"
	VaultEnabled         = "vault.enabled"
	VaultPrefix          = "vault.prefix"
	GinMode              = "gin.mode"
)

type Option = func(*Config)

func WithCfgVar(n string) Option {
	return func(c *Config) {
		c.cfgVar = n
	}
}

type Config struct {
	k      *koanf.Koanf
	parser *kjson.JSON

	cfgVar string
}

func New(opts ...Option) *Config {
	c := &Config{
		k:      koanf.New("."),
		parser: kjson.Parser(),
		cfgVar: AppFile,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func NewFromKoanf(k *koanf.Koanf) *Config {
	return &Config{
		k:      k,
		parser: kjson.Parser(),
	}
}

func NewFromMapObject(v interface{}) (*Config, error) {
	return NewFromJson([]byte(util.ToJson(v)))
}

func NewFromYamlFile(file string) (*Config, error) {
	if strings.TrimSpace(file) == "" {
		return nil, fmt.Errorf("file is required")
	}
	buf, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return NewFromYaml(buf)
}

func NewFromYaml(data []byte) (*Config, error) {
	c := New()
	if err := c.k.Load(rawbytes.Provider(data), kyaml.Parser()); err != nil {
		return nil, err
	}
	return c, nil
}

func NewFromJsonFile(file string) (*Config, error) {
	if strings.TrimSpace(file) == "" {
		return nil, fmt.Errorf("file is required")
	}
	buf, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return NewFromJson(buf)
}

func NewFromJson(data []byte) (*Config, error) {
	c := New()
	if err := c.k.Load(rawbytes.Provider(data), kjson.Parser()); err != nil {
		return nil, err
	}
	return c, nil
}

func NewFromMap(kvm map[string]interface{}) (*Config, error) {
	c := New()
	buf, err := jsoniter.Marshal(kvm)
	if err != nil {
		return nil, err
	}
	if err = c.k.Load(rawbytes.Provider(buf), kjson.Parser()); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Config) BindAndReadConfig(def interface{}) error {
	if err := c.BindFlags(def); err != nil {
		return err
	}
	return c.ReadConfig()
}

// BindFlags 根据结构体def中的定义的 json 属性绑定到 Flags, 然后解析 flags 并读取
func (c *Config) BindFlags(def interface{}) error {
	InitFlags(def)
	return NewFlagSource(c.k).Load()
}

func (c *Config) ReadConfig() error {
	s := []Source{
		NewEnvSource(c.k),
		NewFileSource(c.k, c.string(c.configFileVar(), DefaultCfgFile)),
	}
	if err := c.load(s); err != nil {
		return err
	}

	s = []Source{}
	if c.k.Bool(ConsulEnabled) {
		s = append(s, c.newConsulSource())
	}
	if c.k.Bool(VaultEnabled) {
		s = append(s, c.newVaultSource())
	}
	return c.load(s)
}

func (c *Config) load(s []Source) error {
	for _, item := range s {
		if err := item.Load(); err != nil {
			return err
		}
	}
	return nil
}

func (c *Config) newConsulSource() Source {
	return NewConsulSource(
		c.k,
		c.string(ConsulAddress, DefaultConsulAddress),
		c.k.String(ConsulToken),
		c.string(ConsulPrefix, DefaultConsulPrefix),
		c.AppName(),
		c.k.Bool(ConsulEnabled),
	)
}

func (c *Config) newVaultSource() Source {
	return NewVaultSource(
		c.k,
		c.string(VaultAddress, DefaultVaultAddress),
		c.k.String(strings.Replace(VaultToken, ".", "_", -1)),
		c.string(VaultPrefix, DefaultVaultPrefix),
		c.AppName(),
		c.k.Bool(VaultEnabled),
	)
}

func (c *Config) string(key string, def string) string {
	v := strings.Trim(c.k.String(key), " ")
	if v != "" {
		return v
	}
	return def
}

func (c *Config) configFileVar() string {
	if c.cfgVar != "" {
		return c.cfgVar
	}
	return AppFile
}

func (c *Config) AppName() string {
	return c.string(c.configFileVar(), DefaultAppName)
}

func (c *Config) AppPort() int {
	if v := c.k.Int(AppPort); v > 0 {
		return v
	}
	return DefaultAppPort
}

func (c *Config) GinMode() string {
	return c.string(GinMode, gin.DebugMode)
}

func (c *Config) ConsulAddress() string {
	return c.string(ConsulAddress, DefaultConsulAddress)
}

func (c *Config) VaultAddress() string {
	return c.string(VaultAddress, DefaultVaultAddress)
}

func (c *Config) AppFile() string {
	return c.string(AppFile, DefaultCfgFile)
}

// Dump find key and unmarshal to target
func (c *Config) Dump(key string, target interface{}) error {
	var b []byte
	var err error
	if key == "" {
		b, err = c.k.Marshal(c.parser)
	} else {
		b, err = jsoniter.Marshal(c.k.Get(key))
	}
	if err != nil {
		return err
	}
	return jsoniter.Unmarshal(b, target)
}

func (c *Config) Koanf() *koanf.Koanf {
	return c.k
}

func Dump(source interface{}, target interface{}) error {
	v, err := jsoniter.Marshal(source)
	if err != nil {
		return err
	}
	return jsoniter.Unmarshal(v, target)
}

type Source interface {
	Load() error
}

type KeyValue struct {
	Key    string
	Value  []byte
	Format string
}

type flagSource struct {
	k *koanf.Koanf
}

func NewFlagSource(k *koanf.Koanf) Source {
	return &flagSource{
		k: k,
	}
}

func (s *flagSource) Load() error {
	kvm := make(map[string]string, 10)
	flag.VisitAll(func(f *flag.Flag) {
		kvm[f.Name] = f.Value.String()
	})
	buf, _ := jsoniter.Marshal(kvm)
	if err := s.k.Load(rawbytes.Provider(buf), kjson.Parser()); err != nil {
		return err
	}
	return nil
}

func InitFlags(v interface{}) {
	initFlags(reflect.ValueOf(v))
	flag.Parse()
}

func initFlags(v reflect.Value) {
	v = reflect.Indirect(v)
	rt := v.Type()
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		ft := rt.Field(i)
		f = reflect.Indirect(f)
		tag := ft.Tag.Get("json")
		tag = strings.Split(tag, ",")[0]
		var name = tag
		if f.Kind() == reflect.Struct {
			initFlags(f)
		} else {
			desc := ft.Tag.Get("desc")
			switch f.Kind() {
			case reflect.String:
				flag.String(name, f.String(), desc)
			case reflect.Bool:
				flag.Bool(name, f.Bool(), desc)
			case reflect.Int, reflect.Int64, reflect.Int16, reflect.Int8, reflect.Int32:
				flag.Int64(name, f.Int(), desc)
			case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
				flag.Uint64(name, f.Uint(), desc)
			case reflect.Float32, reflect.Float64:
				flag.Float64(name, f.Float(), desc)
			default:

			}
		}
	}
}

type fileSource struct {
	config string
	k      *koanf.Koanf
}

func NewFileSource(k *koanf.Koanf, config string) Source {
	return &fileSource{
		k:      k,
		config: config,
	}
}

func (f *fileSource) Load() error {
	if f.config != "" {
		buf, err := os.ReadFile(f.config)
		if err != nil {
			return err
		}
		if err = f.k.Load(rawbytes.Provider(buf), kyaml.Parser()); err != nil {
			return err
		}
	}
	return nil
}

func LoadFile(path string, dst interface{}) error {
	k := koanf.New(".")
	l := NewFileSource(k, path)
	if err := l.Load(); err != nil {
		return err
	}
	return NewFromKoanf(k).Dump("", &dst)
}

type envSource struct {
	k *koanf.Koanf
}

func NewEnvSource(k *koanf.Koanf) Source {
	return &envSource{k: k}
}

func (e *envSource) Load() error {
	return e.k.Load(
		kenv.Provider(
			"",
			".",
			func(s string) string {
				return strings.Replace(strings.ToLower(s), "_", ".", -1)
			},
		),
		nil,
	)
}

type consulSource struct {
	k       *koanf.Koanf
	address string
	token   string
	prefix  string
	appName string
	enabled bool
}

func NewConsulSource(k *koanf.Koanf, address, token, prefix, appName string, enabled bool) Source {
	return &consulSource{
		k:       k,
		address: address,
		token:   token,
		prefix:  prefix,
		appName: appName,
		enabled: enabled,
	}
}

func (s *consulSource) Load() error {
	if !s.enabled {
		return nil
	}
	config := api.DefaultConfig()
	config.Address = s.address
	if strings.Trim(s.token, " ") != "" {
		config.Token = strings.Trim(s.token, " ")
	}
	client, err := api.NewClient(config)
	if err != nil {
		return err
	}

	prefixes := []string{s.prefix + "/application/", s.prefix + "/" + s.appName + "/"}
	return s.readConsul(client, prefixes)
}

func (s *consulSource) readConsul(client *api.Client, prefixes []string) error {
	data := make(map[string]string, 150)
	for _, key := range prefixes {
		pairs, _, err := client.KV().List(key, &api.QueryOptions{})
		if err != nil {
			return err
		}
		for _, pair := range pairs {
			data[pair.Key[len(key):]] = string(pair.Value)
		}
	}
	buf, _ := jsoniter.Marshal(data)
	if err := s.k.Load(rawbytes.Provider(buf), kjson.Parser()); err != nil {
		return err
	}
	return nil
}

type vaultSource struct {
	k       *koanf.Koanf
	address string
	token   string
	prefix  string
	appName string
	enabled bool
}

func NewVaultSource(k *koanf.Koanf, address, token, prefix, appName string, enabled bool) Source {
	return &vaultSource{
		k:       k,
		address: address,
		token:   token,
		prefix:  prefix,
		appName: appName,
		enabled: enabled,
	}
}

func (s *vaultSource) Load() error {
	if !s.enabled {
		return nil
	}
	client, err := vaultApi.NewClient(&vaultApi.Config{
		Address: s.address,
	})
	if err != nil {
		return err
	}
	client.SetToken(s.token)
	prefixes := []string{s.prefix + "/application", "secret/" + s.appName}
	return s.readVault(client, prefixes)
}

func (s *vaultSource) readVault(client *vaultApi.Client, prefixes []string) error {
	data := make(map[string]string, 100)
	for _, key := range prefixes {
		secret, err := client.Logical().Read(key)
		if err != nil {
			return err
		}
		if secret == nil || secret.Data == nil {
			continue
		}
		for k, v := range secret.Data {
			switch v.(type) {
			case string:
				data[k] = v.(string)
			default:
				if s, err := jsoniter.MarshalToString(v); err == nil {
					data[k] = s
				}
			}
		}
	}
	buf, _ := jsoniter.Marshal(data)
	if err := s.k.Load(rawbytes.Provider(buf), kjson.Parser()); err != nil {
		return err
	}
	return nil
}

// ArgsDef 应用启动所需的 args
//
//	-app.name datasix  # 服务名称
//	-app.port 8882     # http 端口设置
//	-config ./resource/config.yaml  # 配置文件地址
//	-consul.address http://127.0.0.1:8500 # 默认
//	-consul.enabled true # 可选值(true,false), 默认: true
//	-consul.prefix config # consul 配置前缀, 默认: config
//	-vault.address http://127.0.0.1:8200  # 默认
//	-vault.enabled false # 可选值(true,false), 默认: false
//	-vault.prefix secret # vault 配置前缀, 默认: secret
type ArgsDef struct {
	AppName       string `json:"app.name" desc:"app name registered in consul"`
	AppPort       int    `json:"app.port" desc:"app listen port(default: 8882)"`
	AppFile       string `json:"app.file" desc:"app file(default: ./configs/config.yaml)"`
	ConsulAddress string `json:"consul.address" desc:"consul address, default: http://127.0.0.1:8500"`
	ConsulEnabled string `json:"consul.enabled" desc:"consul enabled, default: true"`
	ConsulPrefix  string `json:"consul.prefix" desc:"consul config prefix, default: config"`
	VaultAddress  string `json:"vault.address" desc:"vault address, default: http://127.0.0.1:8200"`
	VaultEnabled  string `json:"vault.enabled" desc:"vault enabled, default: false"`
	VaultPrefix   string `json:"vault.prefix" desc:"vault config prefix, default: secret"`
}
