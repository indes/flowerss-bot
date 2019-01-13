package model

import "encoding/base64"

func genHashID(sid uint, id string) string {
	idString := string(sid) + "||" + id
	encoded := base64.StdEncoding.EncodeToString([]byte(idString))
	return encoded
}
