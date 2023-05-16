/*
## Задача № 1
Написать API для указанных маршрутов(endpoints)
"/first"  // Случайное число
"/second" // Случайное число
"/summa"  // Сумма двух случайных чисел

результат вернуть в виде JSON

"math/rand"
number := rand.Intn(100)

## Задача № 2

сайт: postman.com
Установить программу и создать учетную запись.

*/

package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
)

type TNumber struct {
	Number int `json:"Number"`
}

type TNumberData struct {
	First  int
	Second int
	Summa  int
}

type TRouter struct {
	Data TNumberData
}

var Router TRouter

func (self *TRouter) Init() {
	self.Data.First = rand.Intn(100)
	self.Data.Second = rand.Intn(100)
	self.Data.Summa = self.Data.First + self.Data.Second
}
func (self *TRouter) ReturnFirst(w http.ResponseWriter, r *http.Request) {
	var Number TNumber
	Number.Number = self.Data.First
	json.NewEncoder(w).Encode(Number)
}
func (self *TRouter) ReturnSecond(w http.ResponseWriter, r *http.Request) {
	var Number TNumber
	Number.Number = self.Data.Second
	json.NewEncoder(w).Encode(Number)
}
func (self *TRouter) ReturnTotal(w http.ResponseWriter, r *http.Request) {
	var Number TNumber
	Number.Number = self.Data.First + self.Data.Second
	json.NewEncoder(w).Encode(Number)
}
func (self *TRouter) Index(w http.ResponseWriter, r *http.Request) {
	type TIndexDoc struct {
		Route       string `json:"Route"`
		Description string `json:"Description"`
	}
	type TIndex struct {
		Index   []TIndexDoc `json:"IndexDoc"`
		Numbers TNumberData `json:"Numbers"`
	}
	var Answer TIndex
	Answer.Index = []TIndexDoc{TIndexDoc{"/First", "Первое число"}, TIndexDoc{"/Second", "Второе число"}, TIndexDoc{"/Summa", "Сумма двух чисел"}}
	Answer.Numbers = self.Data
	json.NewEncoder(w).Encode(Answer)
}

func main() {
	Router.Init()

	http.HandleFunc("/", Router.Index)
	http.HandleFunc("/First", Router.ReturnFirst)
	http.HandleFunc("/Second", Router.ReturnSecond)
	http.HandleFunc("/Summa", Router.ReturnTotal)

	http.ListenAndServe(":1234", nil)

}
