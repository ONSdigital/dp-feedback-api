package models_test

import (
	"testing"

	"github.com/ONSdigital/dp-feedback-api/models"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSanitize(t *testing.T) {
	Convey("HTML strings are correctly sanitized", t, func() {
		s := models.Sanitize(unsafeHTML)
		So(s, ShouldEqual, sanitizedHTML)
	})

	Convey("SQL strings are correctly sanitized", t, func() {
		s := models.Sanitize(unsafeSQL)
		So(s, ShouldEqual, sanitizedSQL)
	})

	Convey("NoSQL strings are correctly sanitized", t, func() {
		s := models.Sanitize(unsafeNoSQL)
		So(s, ShouldEqual, sanitizedNoSQL)
	})
}

// func TestMysqlRealEscapeString(t *testing.T) {
// 	Convey("The character '\\' is correctly escaped", t, func() {
// 		s := models.MysqlRealEscapeString("test \\ test")
// 		So(s, ShouldEqual, "test \\\\ test")
// 	})

// 	Convey("Given a sting containing the character ascii = 0 ", t, func() {
// 		before := []byte{'a', 'b', 0x00, 'c', 'd'}
// 		s := models.MysqlRealEscapeString(string(before))
// 		after := []byte(s)
// 		So(after, ShouldResemble, []byte{'a', 'b', '\\', 0, 'c', 'd'})
// 	})
// }
