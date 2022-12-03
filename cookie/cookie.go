package cookie

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"log"
	"net/http"
)

func generateRandom(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func Cook() string {
	src := []byte("Хотели кушать - и съели Кука.") // данные, которые хотим зашифровать
	//fmt.Printf("original: %s\n", src)

	// будем использовать AES256, создав ключ длиной 32 байта
	key, err := generateRandom(2 * aes.BlockSize)
	if err != nil {
		log.Printf("error: %v\n", err)

	}

	aesblock, err := aes.NewCipher(key)
	if err != nil {
		log.Printf("error: %v\n", err)

	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		log.Printf("error: %v\n", err)

	}

	// создаём вектор инициализации
	nonce, err := generateRandom(aesgcm.NonceSize())
	if err != nil {
		log.Println(err)
	}

	dst := aesgcm.Seal(nil, nonce, src, nil)
	return string(dst)
}

func CreateAndCheckCookieInGET(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("r.Cookies(): ", r.Cookies())
	if len(r.Cookies()) == 0 {
		//fmt.Println("hasn`t cookie")
		dst := Cook()
		dst2 := http.Cookie{Name: "leo", Value: dst}
		http.SetCookie(w, &dst2)
	} else {
		//fmt.Println("has cookie")
		_, err := r.Cookie("leo")
		if err != nil {
			dst := Cook()
			dst2 := http.Cookie{Name: "leo", Value: dst}
			http.SetCookie(w, &dst2)
		}
	}
}

func CheckInPOST(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("leo")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			http.Error(w, "cookie not found", http.StatusBadRequest)
		default:
			log.Println(err)
			http.Error(w, "server error", http.StatusInternalServerError)
		}
		return
	}
}
