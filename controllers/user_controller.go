package controllers

import (
	"github.com/nikola43/ecodadys_api/models"
	"github.com/nikola43/ecodadys_api/utils"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}

	// decode body
	decodeError := utils.DecodeHttpRequestPayload(w, r, user)
	if decodeError != nil {
		utils.RespondHttpError(w, http.StatusUnprocessableEntity, "Invalid resquest payload")
		return
	}
	err := user.Login(con)
	utils.RespondHttpRequest(w, err, user)
}

func SingUp(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}

	// Decode body
	decodeError := utils.DecodeHttpRequestPayload(w, r, user)
	if decodeError != nil {
		utils.RespondHttpError(w, http.StatusUnprocessableEntity, "Invalid resquest payload")
		return
	}

	// Insert user
	err := user.SingUp()
	utils.RespondHttpRequest(w, err, nil)
}

func GetUserByID(w http.ResponseWriter, r *http.Request) {
	user := &models.User{ID: utils.ReadHttpRequestIntegerParam(w, r, "id")}
	err := user.GetUserByID(con)
	utils.RespondHttpRequest(w, err, user)
}

func InviteUser(w http.ResponseWriter, r *http.Request) {
	user := &models.User{ID: utils.ReadHttpRequestIntegerParam(w, r, "invite_id")}
	err := user.InviteUser(con, utils.ReadHttpRequestIntegerParam(w, r, "invited_id"))
	utils.RespondHttpRequest(w, err, nil)
}

func GetNumberOfInvitedUsers(w http.ResponseWriter, r *http.Request) {
	user := &models.User{ID: utils.ReadHttpRequestIntegerParam(w, r, "id")}
	resp, err := user.GetNumberOfInvitedUsers(con)
	utils.RespondHttpRequest(w, err, resp)
}

func GetNumberOfUsers(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	resp, err := user.GetNumberOfUsers(con)
	utils.RespondHttpRequest(w, err, resp)
}

func GetAll(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	users, err := user.GetAll(con)
	utils.RespondHttpRequest(w, err, users)
}
