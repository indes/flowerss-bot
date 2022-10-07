package model

import (
	"encoding/hex"
	"hash/fnv"
)

func GenHashID(sLink string, id string) string {
	idString := string(sLink) + "||" + id
	f := fnv.New32()
	f.Write([]byte(idString))

	encoded := hex.EncodeToString(f.Sum(nil))
	return encoded
}
