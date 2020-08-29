package bot

import (
	valid "github.com/asaskevich/govalidator"
)

// CheckURL check if the string is a URL
func CheckURL(s string) bool {
	return valid.IsURL(s)
}
