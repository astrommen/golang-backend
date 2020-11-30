package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	// "io"
	"log"
	"net/http"

	"golang-backend/helpers"
	"golang-backend/interfaces"
	"golang-backend/users"

	"github.com/gorilla/mux"
)

type Login struct {
	Username string
	Password string
}

type Register struct {
	Username string
	Email    string
	Password string
}

func readBody(r *http.Request) []byte {
	
	body, err := ioutil.ReadAll(r.Body)
	helpers.HandleErr(err)

	return body
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func apiResponse(call map[string]interface{}, w http.ResponseWriter) {
	if call["message"] == "all is fine" {
		resp := call
		json.NewEncoder(w).Encode(resp)
	} else {
		resp := interfaces.ErrResponse{Message: "Wrong username or password"}
		json.NewEncoder(w).Encode(resp)
	}
}

func register(w http.ResponseWriter, r *http.Request) {
	//Allow CORSE here by * or specific origin
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.Header().Set("Access-Control-Allow-Headers", "*")

	// w.Header().Set("Content-Type", "application/javascript")

	// Read body
	body := readBody(r)
	fmt.Println(body)

	// Handle registration
	var formattedBody Register
	err := json.Unmarshal(body, &formattedBody)
	helpers.HandleErr(err)

	register := users.Register(formattedBody.Username, formattedBody.Email, formattedBody.Password)
	

	// Refactor register to use apiRespons
	apiResponse(register, w)

	// Prepare response
	
}

func login(w http.ResponseWriter, r *http.Request) {
	//Allow CORSE here by * or specific origin
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	w.Header().Set("Content-Type", "application/javascript")

	// Refactor login to use readBody
	body := readBody(r)

	var formattedBody Login
	err := json.Unmarshal(body, &formattedBody)
	helpers.HandleErr(err)
	
	login := users.Login(formattedBody.Username, formattedBody.Password)

	// Refactor login to use apiResponse
	apiResponse(login, w)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["id"]
	auth := r.Header.Get("Authorization")

	user := users.GetUser(userId, auth)

	apiResponse(user, w)
}

// func handler(w http.ResponseWriter, r *http.Request) {
// 	resp, err := http.Get("http://jsonplaceholder.typicode.com/posts") // your request to the api

// 	w.Header().Set("Content-Type", "application/javascript")

// 	if err == nil && resp.StatusCode == http.StatusOK {
// 			io.Copy(w, resp.Body)
// 	} else {
// 			json.NewEncoder(w).Encode(err)
// 	}
// }



func StartApi() {
	router := mux.NewRouter()
	// Panic handler middleware
	router.Use(helpers.PanicHandler)
	router.HandleFunc("/login", login).Methods("POST", "OPTIONS")
	router.HandleFunc("/register", register).Methods("POST", "OPTIONS")
	router.HandleFunc("/user/{id}", getUser).Methods("GET")
	// router.HandleFunc("/handler", handler).Methods("GET")
	fmt.Println("App is working on port :8888")
	log.Fatal(http.ListenAndServe(":8888", router))
}

