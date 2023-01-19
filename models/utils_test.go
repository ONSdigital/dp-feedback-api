package models_test

import (
	"testing"

	"github.com/ONSdigital/dp-feedback-api/models"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	unsafeStr    = `<script>document.getElementById("$demo").innerHTML = "Hello JavaScript!";</script>`
	sanitizedStr = `&lt;script&gt;document.getElementById(\&#34;\\$demo\&#34;).innerHTML = \&#34;Hello JavaScript!\&#34;;&lt;/script&gt;`
)

func TestSanitize(t *testing.T) {
	Convey("strings are correctly sanitized", t, func() {
		s := models.Sanitize(unsafeStr)
		So(s, ShouldEqual, sanitizedStr)
	})
}

func TestMysqlRealEscapeString(t *testing.T) {
	Convey(`Given a sting containing the character '\\' then it is correctly escaped`, t, func() {
		s := models.MysqlRealEscapeString(`test \\ test`)
		So(s, ShouldEqual, `test \\\\ test`)
	})

	Convey(`Given a sting containing the character '\n' then it is correctly escaped`, t, func() {
		s := models.MysqlRealEscapeString(`test \n test`)
		So(s, ShouldEqual, `test \\n test`)
	})

	Convey(`Given a sting containing the character '\r' then it is correctly escaped`, t, func() {
		s := models.MysqlRealEscapeString(`test \r test`)
		So(s, ShouldEqual, `test \\r test`)
	})

	Convey(`Given a sting containing the character '\'' then it is correctly escaped`, t, func() {
		s := models.MysqlRealEscapeString(`test \' test`)
		So(s, ShouldEqual, `test \\\' test`)
	})

	Convey(`Given a sting containing the character '"' then it is correctly escaped`, t, func() {
		s := models.MysqlRealEscapeString(`test " test`)
		So(s, ShouldEqual, `test \" test`)
	})

	Convey("Given a sting containing the character ascii = '\032' then it is correctly escaped ", t, func() {
		before := []byte{'a', 'b', '\032', 'c', 'd'}
		s := models.MysqlRealEscapeString(string(before))
		after := []byte(s)
		So(after, ShouldResemble, []byte{'a', 'b', '\\', 'Z', 'c', 'd'})
	})

	Convey("Given a sting containing the character ascii = 0 then it is correctly escaped ", t, func() {
		before := []byte{'a', 'b', 0x00, 'c', 'd'}
		s := models.MysqlRealEscapeString(string(before))
		after := []byte(s)
		So(after, ShouldResemble, []byte{'a', 'b', '\\', 0, 'c', 'd'})
	})
}
