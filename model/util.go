package model

import "encoding/base64"

func genHashID(s *Source, id string) string {
	idString := s.Link + "||" + id
	encoded := base64.StdEncoding.EncodeToString([]byte(idString))
	return encoded
}
