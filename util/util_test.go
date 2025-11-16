package util

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

func Test_If(t *testing.T) {
	v := IfStr(true, "a", "b")
	fmt.Print(v)
}

func TestGetIp(t *testing.T) {
	ip, err := ExternalIP()
	if err != nil {
		t.Fatal(err)
		return
	}
	fmt.Printf("ip: %s", ip.String())
}

func TestFormatFloat(t *testing.T) {
	fmt.Print(FormatFloat(0.111111, 2))
}

func TestSnake2Camel(t *testing.T) {
	fmt.Println(Snake2Camel("__user_a_a_id"))
}

func TestCamel2Snake(t *testing.T) {
	fmt.Println(Camel2Snake("userAAId"))
}

func TestGetUUID(t *testing.T) {
	fmt.Println(GUID())
}

func TestGetenvs(t *testing.T) {
	for _, k := range os.Environ() {
		fmt.Println(k)
	}
}

func TestReflectTime(t *testing.T) {
	tradeDate := "2022-01-01"
	val := reflect.ValueOf(tradeDate)
	fmt.Println(val.Type().String())
}

func TestAbs(t *testing.T) {
	fmt.Println(Abs(-64))
}

func TestFmtWeek(t *testing.T) {
	s := []string{"0", "1", "2", "*"}
	fmt.Println(FmtWeek(s))
}

func TestFmtHours(t *testing.T) {
	s := []string{"3", "0", "1", "2"}
	fmt.Println(FmtHours(s))
}
