package date

import (
	"fmt"
	"time"
)

const (
	TimeZoneOffset     int    = 7
	TimeZoneName       string = "WIB"
	FileDateTimeLayout string = "20060102150405"
	FileDateLayout     string = "20060102"
)

var (
	serverLoc = time.FixedZone("SERVER_TZ", TimeZoneOffset*3600)
)

func TimeToDateStringFileFormat(timeObject time.Time) string {
	loc := time.FixedZone("", TimeZoneOffset*3600)
	return timeObject.In(loc).Format(FileDateLayout)
}
func Now() time.Time {
	return time.Now().In(serverLoc)
}
func DateStringToTime(dateStr string) time.Time {
	dateTimeStr := fmt.Sprintf("%s 00:00:00 +0700", dateStr)
	timeObjc, _ := time.Parse("02/01/2006 15:04:05 -0700", dateTimeStr)
	return timeObjc
}
func ChangeTo(src time.Time, h, m, s int) time.Time {
	src = ToServerTZ(src)
	return time.Date(
		src.Year(),
		src.Month(),
		src.Day(),
		h,
		m,
		s,
		0,
		src.Location(),
	)
}
func ToServerTZ(t time.Time) time.Time {
	return t.In(serverLoc)
}
