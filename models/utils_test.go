package models_test

import (
	"testing"

	"github.com/ONSdigital/dp-feedback-api/config"
	"github.com/ONSdigital/dp-feedback-api/models"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	unsafeStr         = `<script>document.getElementById("$demo").innerHTML = "Hello JavaScript!";</script>`
	sanitizedStr      = `&lt;script&gt;document.getElementById(\&#34;\\$demo\&#34;).innerHTML = \&#34;Hello JavaScript!\&#34;;&lt;/script&gt;`
	htmlSanitizedStr  = `&lt;script&gt;document.getElementById(&#34;$demo&#34;).innerHTML = &#34;Hello JavaScript!&#34;;&lt;/script&gt;`
	sqlSanitizedStr   = `<script>document.getElementById(\"$demo\").innerHTML = \"Hello JavaScript!\";</script>`
	nosqlSanitizedStr = `<script>document.getElementById("\$demo").innerHTML = "Hello JavaScript!";</script>`
)

func TestSanitize(t *testing.T) {
	Convey("given a fully disabled sanitize config", t, func() {
		cfg := &config.Sanitize{
			HTML:  false,
			SQL:   false,
			NoSQL: false,
		}

		Convey("Then no sanitization is performed", func() {
			s := models.Sanitize(cfg, unsafeStr)
			So(s, ShouldEqual, unsafeStr)
		})
	})

	Convey("given a fully enabled sanitize config", t, func() {
		cfg := &config.Sanitize{
			HTML:  true,
			SQL:   true,
			NoSQL: true,
		}

		Convey("Then a full sanitization is performed", func() {
			s := models.Sanitize(cfg, unsafeStr)
			So(s, ShouldEqual, sanitizedStr)
		})
	})

	Convey("given an html-only sanitize config", t, func() {
		cfg := &config.Sanitize{
			HTML:  true,
			SQL:   false,
			NoSQL: false,
		}

		Convey("Then an html-only sanitization is performed", func() {
			s := models.Sanitize(cfg, unsafeStr)
			So(s, ShouldEqual, htmlSanitizedStr)
		})
	})

	Convey("given a sql-only sanitize config", t, func() {
		cfg := &config.Sanitize{
			HTML:  false,
			SQL:   true,
			NoSQL: false,
		}

		Convey("Then a sql-only sanitization is performed", func() {
			s := models.Sanitize(cfg, unsafeStr)
			So(s, ShouldEqual, sqlSanitizedStr)
		})
	})

	Convey("given a nosql-only sanitize config", t, func() {
		cfg := &config.Sanitize{
			HTML:  false,
			SQL:   false,
			NoSQL: true,
		}

		Convey("Then a nosql-only sanitization is performed", func() {
			s := models.Sanitize(cfg, unsafeStr)
			So(s, ShouldEqual, nosqlSanitizedStr)
		})
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
