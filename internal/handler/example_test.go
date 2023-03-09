// example testing
package handler

import (
	"encoding/hex"
	"log"
)

func ExampleGenUlid() {
	short := GenUlid()
	shortURL := hex.EncodeToString([]byte(short))
	log.Printf("Short URL is %s \n", shortURL)
}
