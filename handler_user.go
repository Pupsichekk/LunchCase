package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
)

func postUser(apiCfg *apiConfig, user User) error {
	_, err := apiCfg.DB.Query("INSERT INTO users(id, email, password, name, admin, checked) values ($1, $2, $3, $4, $5, $6) RETURNING id, name, password, name, admin, checked", user.ID, user.Email, user.Password, user.Name, user.Admin, user.Checked)
	if err != nil {
		return err
	}
	return nil
}

func (apiCfg *apiConfig) handlerUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
		Admin    bool   `json:"admin"`
		Checked  bool   `json:"checked"`
	}
	params := parameters{}

	if r.Method == "POST" {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Error parsing json: %v", err))
		}
		err = json.Unmarshal(body, &params)
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Eror unmarshaling json: %v", err))
		}
		err = postUser(apiCfg, User{
			ID:       uuid.New(),
			Email:    params.Email,
			Password: params.Password,
			Name:     params.Name,
			Admin:    params.Admin,
			Checked:  params.Checked,
		})
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Eror inserting user: %v", err))
		}
		respondWithJson(w, 200, struct{}{})
		return
	}
	respondWithError(w, 501, fmt.Sprintf("Eror handling request: %v", fmt.Errorf("unsupported method: %v", r.Method)))
}
