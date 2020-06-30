package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/badoux/checkmail"
	"github.com/gorilla/mux"
)

type Email struct {
	Address string `json:"address,omitempty"`
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
}

func main()  {
	router := mux.NewRouter()
	router.HandleFunc("/verify", verifyEmail).Methods("POST")


	fmt.Println("Server Running on 127.0.0.1:8080")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func verifyEmail(w http.ResponseWriter, r *http.Request) {
	var input Email

	body, _ := ioutil.ReadAll(r.Body)

	_ = json.Unmarshal(body, &input)

	w.Header().Add("Content-Type", "application/json")

	err := checkmail.ValidateHost(input.Address)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		responseErr := &Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		_ = json.NewEncoder(w).Encode(responseErr)
		fmt.Println(err)
		return
	}

	err = checkmail.ValidateHost(input.Address)
	if smtpErr, ok := err.(checkmail.SmtpError); ok && err != nil {
		w.WriteHeader(http.StatusBadRequest)
		responseErr := &Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		_ = json.NewEncoder(w).Encode(responseErr)
		fmt.Printf("Code: %s, Msg: %s", smtpErr.Code(), smtpErr)
		return
	}

	response := &Response{
		Code:    http.StatusOK,
		Message: "Success verify email!",
	}

	_ = json.NewEncoder(w).Encode(response)
}