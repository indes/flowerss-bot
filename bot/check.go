package bot

import (
	valid "github.com/asaskevich/govalidator"
)

func CheckUrl(s string) bool {
	return valid.IsURL(s)
}
