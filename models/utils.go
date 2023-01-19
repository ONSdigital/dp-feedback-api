package models

import (
	"html"
	"strings"
)

// Sanitize sanitizes the input string to prevent html, mysql and nosql (mongodb only) injection attacks
func Sanitize(toSanitize string) string {
	return html.EscapeString(
		MysqlRealEscapeString(
			MongodbEscapeString(toSanitize),
		),
	)
}

// MysqlRealEscapeString escapes the control characters used by sql commands
func MysqlRealEscapeString(value string) string {
	var sb strings.Builder
	for i := 0; i < len(value); i++ {
		c := value[i]
		switch c {
		case '\\', 0, '\n', '\r', '\'', '"':
			sb.WriteByte('\\')
			sb.WriteByte(c)
		case '\032':
			sb.WriteByte('\\')
			sb.WriteByte('Z')
		default:
			sb.WriteByte(c)
		}
	}
	return sb.String()
}

// MongodbEscapeString escapes the control characters used by mongodb commands
// TODO at the moment we only check for '$' char. We may need a deeper discussion about this
func MongodbEscapeString(value string) string {
	var sb strings.Builder
	for i := 0; i < len(value); i++ {
		c := value[i]
		switch c {
		case '$':
			sb.WriteByte('\\')
			sb.WriteByte(c)
		default:
			sb.WriteByte(c)
		}
	}
	return sb.String()
}
