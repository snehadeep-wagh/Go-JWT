package helpers

import (
	"errors"
	"net/http"
	"github.com/gorilla/mux"
)

func CheckUserType(r *http.Request, role string) (err error) {
	err = nil
	params := mux.Vars(r)
	user_type := params["user_type"]

	if user_type != role {
		err = errors.New("Unauthorized access to the data")
		return err
	}
	return err
}

func CheckUserTypeWithUserId(r *http.Request, id string) (err error) {
	// check if the user is not admin and have the same id
	params := mux.Vars(r)
	user_type := params["user_type"]
	uid := params["user_id"]
	err = nil

	if user_type == "USER" && uid != id {
		err = errors.New("Unauthorized access to the data")
		return err
	}

	// check if the usertype
	// in this case this is redundant
	err = CheckUserType(r, user_type)
	return err
}
