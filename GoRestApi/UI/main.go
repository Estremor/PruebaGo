package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/gorilla/mux"
)

type Customer struct {
	Id        string `json:"Id,omitempty"`
	FirstName string `json:"FirstName,omitempty"`
	LastName  string `json:"LastName,omitempty"`
	Score     int64  `json:"Score,omitempty"`
	Min       int32  `json:"Min,omitempty"`
}

var CustonArra []Customer
var ListWhite []Customer

func GetData() {
	file, error := ioutil.ReadFile("C:\\Users\\Alexander\\Desktop\\GoRestApi\\UI\\CustomerDa.json")

	if error != nil {
		log.Fatal("ocurrio un error", error)
	}

	str := string(file)
	//fmt.Print(str)
	json.Unmarshal([]byte(str), &CustonArra)
}

func GetCustomer(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)
	GetData()
	json.NewEncoder(w).Encode(CustonArra)
}

func GetOneCustomer(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)
	params := mux.Vars(req)
	for _, item := range CustonArra {
		if item.Id == params["Id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Customer{})
}

func CreateCustomer(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)

	params := mux.Vars(req)
	var client Customer
	GetData()
	_ = json.NewDecoder(req.Body).Decode(&client)
	client.Id = params["Id"]
	CustonArra = append(CustonArra, client)

	str, _ := json.Marshal(CustonArra)
	dataBytes := []byte(string(str))
	ioutil.WriteFile("C:\\Users\\Alexander\\Desktop\\GoRestApi\\UI\\CustomerDa.json", []byte(""), 0)
	ioutil.WriteFile("C:\\Users\\Alexander\\Desktop\\GoRestApi\\UI\\CustomerDa.json", dataBytes, 0)

	json.NewEncoder(w).Encode(CustonArra)

}

type byScoreF []Customer

func (a byScoreF) Len() int           { return len(a) }
func (a byScoreF) Less(i, j int) bool { return a[i].Score < a[j].Score && a[i].Min == 60 }
func (a byScoreF) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func ResetData() {
	date := time.Now()
	lentAr := len(CustonArra)
	actual := int32(date.Minute())

	for index := 0; index < lentAr; index++ {
		if ((actual - CustonArra[index].Min) > 4) || CustonArra[index].Min > actual {
			CustonArra[index].Min = 60
		}
	}
}

func ByScore(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)
	GetData()
	ResetData()

	sort.Sort(byScoreF(CustonArra))

	var score int64
	for index := (len(CustonArra) - 1); index >= 0; index-- {
		if CustonArra[index].Min == 60 {
			score = CustonArra[index].Score
			break
		}
	}

	resultcus := Filter(CustonArra, func(v Customer) bool {
		return v.Score == score //&& v.Min == 60
	})

	i := 0
	for index := 0; index < len(CustonArra); index++ {
		if CustonArra[index].Id == resultcus[i].Id {
			CustonArra[index].Min = int32(time.Now().Minute())
			i++
			if i == len(resultcus) {
				break
			}
		}
	}

	str, _ := json.Marshal(CustonArra)
	dataBytes := []byte(string(str))
	ioutil.WriteFile("C:\\Users\\Alexander\\Desktop\\GoRestApi\\UI\\CustomerDa.json", []byte(""), 0)
	ioutil.WriteFile("C:\\Users\\Alexander\\Desktop\\GoRestApi\\UI\\CustomerDa.json", dataBytes, 0)

	json.NewEncoder(w).Encode(resultcus)
}

func EnabledIDs(w http.ResponseWriter, req *http.Request) {
	GetData()
	ResetData()
	resultcus := Filter(CustonArra, func(v Customer) bool {
		return v.Min == 60
	})
	json.NewEncoder(w).Encode(resultcus)
}

func CanBeServed(w http.ResponseWriter, req *http.Request) {
	enableCors(&w)

	params := mux.Vars(req)

	response, err := http.Get("http://localhost:8000/EnabledIDs/")
	if err != nil {
		json.NewEncoder(w).Encode(err)
	}

	var clients []Customer
	_ = json.NewDecoder(response.Body).Decode(&clients)

	id := params["CustomerId"]
	res := Any(clients, func(v Customer) bool {
		return v.Id == id && v.Min == 60
	})

	json.NewEncoder(w).Encode(res)

}

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/Customer", GetCustomer).Methods("GET")
	router.HandleFunc("/Customer/{Id}", GetOneCustomer).Methods("GET")
	router.HandleFunc("/Customer/{Id}", CreateCustomer).Methods("POST")
	router.HandleFunc("/ByScore/", ByScore).Methods("GET")
	router.HandleFunc("/CanBeServed/{CustomerId}", CanBeServed).Methods("GET")
	router.HandleFunc("/EnabledIDs/", EnabledIDs).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", router))
	println("Hola mundo")
}

func Filter(vs []Customer, f func(Customer) bool) []Customer {
	vsf := make([]Customer, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}
func Any(vs []Customer, f func(Customer) bool) bool {
	for _, v := range vs {
		if f(v) {
			return true
		}
	}
	return false
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
