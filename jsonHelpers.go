package main

import (
	"encoding/json"
	"net/http"
)

type errResBody struct {
	Error string `json:"error"`
}

func respondWithError(w http.ResponseWriter, code int, msg string) error {
	eRes := &errResBody{
		Error: msg,
	}

	data, err := json.Marshal(eRes)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	w.Write(data)

	return nil
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) error {
	res, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	w.Write(res)
	return nil
}
