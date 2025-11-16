package util

import (
	"fmt"
	"testing"
	"time"
)

func TestParseTime(t *testing.T) {
	datetime, err := Parse(String("2020-06-01 10:10:46"), String(YYYY_MM_DD_HH_mm_ss), String(EST))
	if err != nil {
		t.Fatal(err)
		return
	}
	fmt.Print(datetime.String())
}

func TestAddMin(t *testing.T) {
	now := time.Now()
	datetime, err := AddMin(&now, 10, String(EST))
	if err != nil {
		t.Fatal(err)
		return
	}
	fmt.Print(datetime.String())
}

func TestFormat(t *testing.T) {
	now := time.Now()
	str, _ := Format(&now, String(YYYY_MM_DD), String(EST))
	fmt.Printf("date: %s", *str)
	fmt.Print("\n")
	str, _ = Format(&now, String(YYYY_MM_DD_HH_mm_ss), String(EST))
	fmt.Printf("datetime: %s", *str)
	fmt.Print("\n")
}

func TestBetween(t *testing.T) {
	now, err := Parse(String("2021-07-30 01:04:05"), String(YYYY_MM_DD_HH_mm_ss), String(UTC))
	if err != nil {
		return
	}
	s, _ := Format(now, String(YYYY_MM_DD_HH_mm_ss), String(EST))
	fmt.Printf("dateStart: %s", *s)
	fmt.Println("")
	dateStart, _ := Trunc(now, String(EST))
	fmt.Printf("dateStart: %s", dateStart)
	fmt.Println("")
	scaleStart, _ := AddMin(dateStart, 10, String(EST))
	scaleEnd, _ := AddMin(dateStart, 1440, String(EST))

	start, _ := Format(scaleStart, String(YYYY_MM_DD_HH_mm_ss), String(EST))
	end, _ := Format(scaleEnd, String(YYYY_MM_DD_HH_mm_ss), String(EST))
	fmt.Printf("start: %s end:%s\n", *start, *end)

	if now.After(*scaleStart) && now.Before(*scaleEnd) {
		fmt.Print("in scale window")
	} else {
		fmt.Print("not in scale window")
	}
}

func TestParse2(t *testing.T) {
	now, err := Parse(String("2021-07-30 10:04:05"), String(YYYY_MM_DD_HH_mm_ss), String(EST))
	if err != nil {
		return
	}
	s, _ := Format(now, String(YYYY_MM_DD_HH_mm_ss), String(EST))
	fmt.Printf("dateStart: %s", *s)
	fmt.Println("")
}

func TestHour(t *testing.T) {
	now, err := Parse(String("2021-07-30 00:04:05"), String(YYYY_MM_DD_HH_mm_ss), String(EST))
	if err != nil {
		t.Fatal(err)
	}
	//now = Now()
	hour, err := Hour(now, String(EST))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("hour: %d \n", hour)
}

func TestParseUtcTimestamp(t *testing.T) {
	tm := ParseUtcTimestamp("1636084800000")
	if tm != nil {
		fmt.Printf("%v", tm.String())
	} else {
		fmt.Printf("%v", "nil")
	}

}

func TestTimeLocal(t *testing.T) {
	s, _ := Format(Now(), String(YYYY_MM_DD_HH_mm_ss), String(time.Local.String()))
	fmt.Println(*s)
}

func TestFormatISO8601(t *testing.T) {
	s, _ := Format(Now(), String(ISO8601oOmitted), String(EST))
	fmt.Println(*s)
}

func TestFormatISO(t *testing.T) {
	pattern := "Mon Jan 2 15:04:05 MST 2006"
	s, _ := Format(Now(), String(pattern), String(EST))
	fmt.Println(*s)
}

func TestGetWeekDay(t *testing.T) {
	start, _ := Parse(String("2022-09-05"), String(YYYY_MM_DD), String(EST))
	for i := 0; i < 20; i++ {
		day, _ := AddDay(start, i, String(EST))
		weekDay, err := GetWeekDay(day, EST)
		if err != nil {
			t.Fatal(err)
		}
		dayStr := FormatDateTime(day, EST)
		fmt.Println(fmt.Sprintf("day: %s weekDay: %d", dayStr, weekDay))
	}

}
