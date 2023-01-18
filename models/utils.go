package models

import (
	"html"
)

func Sanitize(toSanitize string) string {
	return html.EscapeString(toSanitize)
}

// // MysqlRe// as per
// func MysqlRealEscapeString(value string) string {
// 	var sb strings.Builder
// 	for i := 0; i < len(value); i++ {
// 		c := value[i]
// 		switch c {
// 		case '\\', 0, '\n', '\r', '\'', '"':
// 			sb.WriteByte('\\')
// 			sb.WriteByte(c)
// 		case '\032':
// 			sb.WriteByte('\\')
// 			sb.WriteByte('Z')
// 		default:
// 			sb.WriteByte(c)
// 		}
// 	}
// 	return sb.String()
// }

// func ParseMySQL(value string) string {
// 	// db := &sql.DB{}
// 	db := sql.OpenDB(nil)
// 	_, err := db.Prepare(value)
// 	if err != nil {
// 		return value
// 		// log.Fatal(err)
// 	} else {
// 		return ""
// 	}
// }
