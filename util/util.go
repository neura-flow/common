package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	uuid "github.com/satori/go.uuid"
)

// String returns a pointer to the string value passed in.
func String(v string) *string {
	return &v
}

func Int(v int) *int {
	return &v
}

func Int32(v int32) *int32 {
	return &v
}

func Int64(v int64) *int64 {
	return &v
}

func Bool(b bool) *bool {
	return &b
}

func Float64(v float64) *float64 {
	return &v
}

func Duration(v time.Duration) *time.Duration {
	return &v
}

func Time(v time.Time) *time.Time {
	return &v
}

func Map(v map[string]string) *map[string]string {
	return &v
}

func Abs(n int64) int64 {
	y := n >> 63
	return (n ^ y) - y
}

func ParseInt(s string) *int {
	v, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return nil
	}
	return Int(int(v))
}

func ParseInt64(s string) *int64 {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return nil
	}
	return &v
}

func Array(v []string) []*string {
	if v == nil {
		return nil
	}
	if len(v) == 0 {
		return nil
	}
	nv := make([]*string, 0)
	for _, item := range v {
		nv = append(nv, String(item))
	}
	return nv
}

func IsBlank(value *string) bool {
	if value == nil || strings.TrimSpace(*value) == "" {
		return true
	}
	return false
}

func IsNotBlank(value *string) bool {
	return !IsBlank(value)
}

func AnyBlank(values ...*string) bool {
	if len(values) == 0 {
		return false
	}

	isBlank := false
	for _, value := range values {
		if value == nil || strings.TrimSpace(*value) == "" {
			isBlank = true
			break
		}
	}
	return isBlank
}

// IsNotBlankArray 其中任何一个元素不为空，识别为不空
func IsNotBlankArray(value []string) bool {
	var notBlank bool
	for _, item := range value {
		if IsNotBlank(&item) {
			notBlank = true
			break
		}
	}
	return notBlank
}

func IsBlankArray(value []string) bool {
	if len(value) == 0 {
		return true
	}
	return !IsNotBlankArray(value)
}

func OneOf(ele string, targets []string) bool {
	for _, target := range targets {
		if strings.EqualFold(ele, target) {
			return true
		}
	}
	return false
}

func StringVal(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func BoolToStr(b *bool) *string {
	v := "false"
	if b != nil && *b {
		v = "true"
	}
	return &v
}

func StrToBool(s *string) *bool {
	if s == nil {
		return Bool(false)
	}
	b, err := strconv.ParseBool(*s)
	if err != nil {
		return Bool(false)
	}
	return Bool(b)
}

func IntToStr(i *int) *string {
	if i == nil {
		return nil
	}
	v := fmt.Sprint(*i)
	return &v
}

func IfInt32(condition bool, a, b int32) int32 {
	v := If(condition, Int32(a), Int32(b))
	c := v.(*int32)
	return *c
}

func IfInt(condition bool, a, b int) int {
	v := If(condition, a, b)
	c := v.(int)
	return c
}

func IfStr(condition bool, a, b string) string {
	v := If(condition, String(a), String(b))
	c := v.(*string)
	return *c
}

func If(condition bool, a, b interface{}) interface{} {
	if condition {
		return a
	}
	return b
}

func ExternalIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			ip := getIpFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip, nil
		}
	}
	return nil, fmt.Errorf("error connected to the network")
}

func getIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil // not an ipv4 address
	}

	return ip
}

func FormatFloat(f float64, precision int) string {
	return strconv.FormatFloat(f, 'f', precision, 64)
}

func Snake2Camel(name string) string {
	// _id 不转化
	if strings.EqualFold(name, "_id") {
		return name
	}
	name = strings.Replace(name, "_", " ", -1)
	name = strings.Title(name)
	return Lcfirst(strings.Replace(name, " ", "", -1))
}

func Camel2Snake(name string) string {
	buffer := NewBuffer()
	for i, r := range name {
		if unicode.IsUpper(r) {
			if i != 0 {
				buffer.Append('_')
			}
			buffer.Append(unicode.ToLower(r))
		} else {
			buffer.Append(r)
		}
	}
	return buffer.String()
}

func Ucfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

func Lcfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

type Buffer struct {
	*bytes.Buffer
}

func NewBuffer() *Buffer {
	return &Buffer{Buffer: new(bytes.Buffer)}
}

func (b *Buffer) Append(i interface{}) *Buffer {
	switch val := i.(type) {
	case int:
		b.append(strconv.Itoa(val))
	case int64:
		b.append(strconv.FormatInt(val, 10))
	case uint:
		b.append(strconv.FormatUint(uint64(val), 10))
	case uint64:
		b.append(strconv.FormatUint(val, 10))
	case string:
		b.append(val)
	case []byte:
		b.Write(val)
	case rune:
		b.WriteRune(val)
	}
	return b
}

func (b *Buffer) append(s string) *Buffer {
	defer func() {
		if err := recover(); err != nil {
			log.Println("*****内存不够了！******")
		}
	}()
	b.WriteString(s)
	return b
}

func Mill(t time.Time) int64 {
	return t.Unix() * 1000
}

const (
	defaultPageSize int = 200
)

func PageSize(pageSize int) int {
	if pageSize <= 0 || pageSize > defaultPageSize {
		return defaultPageSize
	}
	return pageSize
}

func IntArrToStr(v []int) (ret []string) {
	if len(v) == 0 {
		ret = make([]string, 0)
		return
	}
	ret = make([]string, 0)
	for _, item := range v {
		ret = append(ret, fmt.Sprint(item))
	}
	return ret
}

func IntArrToMap(v []int) map[int]bool {
	m := make(map[int]bool)
	for _, item := range v {
		m[item] = true
	}
	return m
}

func StrArrToMap(v []string) map[string]bool {
	m := make(map[string]bool)
	for _, item := range v {
		m[item] = true
	}
	return m
}

func Now() *time.Time {
	v := time.Now()
	return &v
}

func ToJson(v interface{}) string {
	if v == nil {
		return ""
	}
	bytes, _ := json.Marshal(v)
	return string(bytes)
}

func ToJsonBytes(v interface{}) []byte {
	return []byte(ToJson(v))
}

func GUID() (id string) {
	v := uuid.NewV4().String()
	v = strings.Replace(v, "-", "", -1)
	return v
}

func Str2Bool(str string) bool {
	tbool, _ := strconv.ParseBool(str)
	return tbool
}
func Str2Int(str string) int {
	tint, _ := strconv.Atoi(str)
	return tint
}

func NowMilliSecond() int64 {
	return time.Now().UnixNano() / 1e6
}

// Distinct 字段列表去重
func Distinct(fields []string) []string {
	if len(fields) <= 1 {
		return fields
	}
	fieldMap := make(map[string]bool)
	for _, field := range fields {
		fieldMap[field] = true
	}
	distincts := make([]string, 0)
	for k := range fieldMap {
		distincts = append(distincts, k)
	}
	return distincts
}

func GetRealValue(val reflect.Value) reflect.Value {
	kind := val.Kind()
	if kind == reflect.Ptr {
		return GetRealValue(val.Elem())
	} else if kind == reflect.Interface {
		if val.CanInterface() {
			return GetRealValue(reflect.ValueOf(val.Interface()))
		} else {
			return GetRealValue(val.Elem())
		}
	} else {
		return val
	}
}

func SortString(list []string) []string {
	if len(list) > 0 {
		sort.SliceStable(list, func(i, j int) bool {
			return strings.Compare(list[j], list[i]) > 0
		})
	}
	return list
}

func SortInt(list []int) []int {
	if len(list) > 0 {
		sort.SliceStable(list, func(i, j int) bool {
			return list[j] > list[i]
		})
	}
	return list
}

func FmtWeek(weekDays []string) string {
	text := make([]string, 0)
	for _, item := range SortString(weekDays) {
		v := ""
		switch item {
		case "*":
			v = "无限制"
			return v
		case "0":
			v = "周日"
		case "1":
			v = "周一"
		case "2":
			v = "周二"
		case "3":
			v = "周三"
		case "4":
			v = "周四"
		case "5":
			v = "周五"
		case "6":
			v = "周六"
		default:

		}
		if IsNotBlank(&v) {
			text = append(text, v)
		}
	}
	return strings.Join(text, ",")
}

func FmtHours(hours []string) string {
	hourMap := make(map[string]string)
	for i := 0; i < 24; i++ {
		hourMap[fmt.Sprint(i)] = fmt.Sprintf("[%d-%d)", i, i+1)
	}
	text := make([]string, 0)
	for _, item := range SortString(hours) {
		if strings.EqualFold(item, "*") {
			text = []string{"无限制"}
			break
		}
		if v, found := hourMap[item]; found {
			text = append(text, v)
		}
	}
	return strings.Join(text, ",")
}

func Div(a, b float64) float64 {
	if b == 0 {
		return 0
	}
	return Decimal(a / b)
}

func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}

func ShortId(id string) string {
	if len(id) > 8 {
		return id[:8]
	} else {
		return id
	}
}

func JoinInts(eles []int) string {
	s := make([]string, len(eles))
	for i, ele := range eles {
		s[i] = fmt.Sprintf("%d", ele)
	}
	return strings.Join(s, ",")
}

func GetEnvAny(names ...string) string {
	for _, n := range names {
		if val := os.Getenv(n); val != "" {
			return val
		}
	}
	return ""
}

func ParseUrls(urls []string) ([]*url.URL, error) {
	uList := make([]*url.URL, 0, len(urls))
	for i, ev := range urls {
		urlList, err := urlsFromStr(ev)
		if err != nil {
			return nil, fmt.Errorf("failed to parse url field '%v': %v", strings.Join(urls, ".")+"."+strconv.Itoa(i), err)
		}
		uList = append(uList, urlList...)
	}

	return uList, nil
}

func ParseUrl(url string) (*url.URL, error) {
	urlList, err := urlsFromStr(url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse url field '%v': %v", url, err)
	}
	if len(urlList) > 0 {
		u := urlList[0]
		return u, nil
	}
	return nil, nil
}

func urlsFromStr(str string) (urls []*url.URL, err error) {
	for _, s := range strings.Split(str, ",") {
		if s = strings.TrimSpace(s); s == "" {
			continue
		}
		var u *url.URL
		if u, err = url.Parse(s); err != nil {
			return
		}
		urls = append(urls, u)
	}
	return
}

// ExpandAddress 展开地址, 支持单个地址元素用逗号分隔
func ExpandAddress(addresses []string) []string {
	list := make([]string, 0)
	for _, addr := range addresses {
		for _, splitAddr := range strings.Split(addr, ",") {
			if trimmed := strings.TrimSpace(splitAddr); len(trimmed) > 0 {
				list = append(list, trimmed)
			}
		}
	}
	return list
}

func ToMarkDown(slice []string, header bool) string {
	var buf strings.Builder
	buf.WriteString("| ")
	buf.WriteString(strings.Join(slice, " | "))
	buf.WriteString(" |\n")
	if header {
		buf.WriteString("| ")
		buf.WriteString(strings.Repeat("--- |", len(slice)))
		buf.WriteString("\n")
	}
	return buf.String()
}
