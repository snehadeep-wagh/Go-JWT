package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	db "github.com/snehadeep-wagh/go-backend/database"
	"github.com/snehadeep-wagh/go-backend/helpers"
	"github.com/snehadeep-wagh/go-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var dbClient = db.Client
var userCollection = db.OpenCollection(dbClient, "user")
var validate = validator.New()

func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	// check if the user accessing it is admin
	err := helpers.CheckUserType(r, "ADMIN")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()
	cur, err := userCollection.Find(ctx, bson.D{{}})
	if err != nil {
		http.Error(w, "Problem while fetching all the documents occured!", http.StatusInternalServerError)
		return
	}

	var users_list []primitive.M
	for cur.Next(ctx) {
		var user_details bson.M
		if err := cur.Decode(&user_details); err != nil {
			http.Error(w, "Problem while decoding the result!", http.StatusInternalServerError)
			return
		}

		users_list = append(users_list, user_details)
	}
	defer cur.Close(ctx)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users_list)
}

func GetUserById(w http.ResponseWriter, r *http.Request) {
	fmt.Print("this is getuserbyid!")
	params := mux.Vars(r)
	userId := params["userId"]

	w.Header().Set("content-type", "application/json")

	// here we check if the user is having same id or not
	if err := helpers.CheckUserTypeWithUserId(r, userId); err != nil {
		// return the error
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, "error: "+err.Error(), 400)
		return
	}

	// create the context with time 100sec
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var user models.User

	// get the data from user data
	err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
	defer cancel()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, "Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func Signup(w http.ResponseWriter, r *http.Request) {
	// fmt.Print("this is signup!")
	w.Header().Set("content-type", "application/json")
	var user models.User
	// unmarshal
	json.NewDecoder(r.Body).Decode(&user)
	fmt.Print(user)

	// validate the details entered by the user
	validationErr := validate.Struct(user)
	if validationErr != nil {
		http.Error(w, "Validation Error: Check the details entered", http.StatusBadRequest)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	emailCount, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		// log.Panic(err)
		http.Error(w, "Error occured while checking the email.", http.StatusInternalServerError)
		return
	}

	// Hash the password and store it into the user struct object
	password := HashPassword(*user.Password)
	user.Password = &password

	phoneCount, err := userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		// log.Panic(err)
		http.Error(w, "Error occured while checking the phone number.", http.StatusInternalServerError)
		return
	}

	if phoneCount > 0 || emailCount > 0 {
		fmt.Fprint(w, "Email or phone already exists!")
		return
	}

	// now add the other details to the user
	create_time, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	update_time, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.Updated_at = update_time
	user.Created_at = create_time
	user.Id = primitive.NewObjectID()
	user.User_id = user.Id.Hex()

	token, refreshToken, err := helpers.CreateAllTokens(*user.Email, *user.First_name, *user.Last_name, *user.User_type, *&user.User_id)
	user.Token = &token
	user.Refresh_token = &refreshToken

	insertCount, insertErr := userCollection.InsertOne(ctx, user)
	if insertErr != nil {
		fmt.Print(w, "Unable to insert the item!")
		http.Error(w, "Unable to insert the item!", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(insertCount)
	json.NewEncoder(w).Encode(user)
}

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	// user detail details entered by the user
	var user models.User

	// user details fetched from the database
	var dbUser models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Problem in getting the user entered details!", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()
	dbErr := userCollection.FindOne(ctx, bson.M{"email": *user.Email}).Decode(&dbUser)
	if dbErr != nil {
		// fmt.Print("Error: " + dbErr.Error())
		http.Error(w, "Email or password is incorrect!", http.StatusInternalServerError)
		return
	}

	isValidPass, msg := VerifyPassword(*dbUser.Password, *user.Password)

	if isValidPass != true {
		http.Error(w, "Error: "+msg, http.StatusBadRequest)
		return
	}

	token, refreshToken, _ := helpers.CreateAllTokens(*dbUser.Email, *dbUser.First_name, *dbUser.Last_name, *dbUser.User_type, dbUser.User_id)

	helpers.UpdateAllTokens(token, refreshToken, dbUser.User_id)

	dbErr = userCollection.FindOne(ctx, bson.M{"user_id": dbUser.User_id}).Decode(&dbUser)
	if dbErr != nil {
		http.Error(w, "Error: "+dbErr.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dbUser)
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	// compare the password provided by the user and the stored password at
	// the time of signup
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(providedPassword))
	check := true
	msg := ""

	if err != nil {
		check = false
		msg = fmt.Sprintf(userPassword + " " + providedPassword + " " + "password is incorrect!")
	}

	return check, msg
}

func HashPassword(pass string) string {
	encrPass, err := bcrypt.GenerateFromPassword([]byte(pass), 2)
	if err != nil {
		log.Panic(err)
	}

	return string(encrPass)
}
