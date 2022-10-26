package handler

import (
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"math/rand"
	"net/http"
	"reflect"
	"time"
)

//CreateShortUrlHandler Эндпоинт POST / принимает в теле запроса строку URL для сокращения
//и возвращает ответ с кодом 201 и сокращённым URL в виде текстовой строки в теле.
func CreateShortUrlHandler(w http.ResponseWriter, r *http.Request) {
	/*
		забираем url адресс из тела
		генерим число и добавляем в мапу id[url]
		возвращаем сгенерированное id
	*/

	m := map[int]string{}

	rand.Seed(time.Now().UnixNano())
	min := 1
	max := 1000
	count := rand.Intn(max-min+1) + min

	defer r.Body.Close()
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("error: %s", err)
	} else {
		m[count] = string(payload)
	}

	fmt.Println(m)
	fmt.Fprintf(w, "payload: %s\n", payload)

}

//GetShortUrlByIdHandler Эндпоинт GET /{id} принимает в качестве URL-параметра идентификатор сокращённого URL
//и возвращает ответ с кодом 307 и оригинальным URL в HTTP-заголовке Location.
func GetShortUrlByIdHandler(w http.ResponseWriter, r *http.Request) {
	/*
		принимаем id из url-параметра
		ищем по ключу id значение url

	*/
	vars := mux.Vars(r)
	fmt.Println(reflect.TypeOf(vars["id"]))

	fmt.Fprintf(w, "id: %s\n", vars["id"])
}
