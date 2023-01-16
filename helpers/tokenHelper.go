package helpers

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt"
	"github.com/snehadeep-wagh/go-backend/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetails struct {
	Email    string
	Fname    string
	Lname    string
	Uid      string
	UserType string
	jwt.StandardClaims
}

var secrete_key = os.Getenv("SECRETE_KEY")
var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

func CreateAllTokens(email string, fName string, lName string, usrType string, uid string) (signedToken string, signedRefreshToken string, err error) {
	claim := &SignedDetails{
		Email:    email,
		Fname:    fName,
		Lname:    lName,
		Uid:      uid,
		UserType: usrType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaim := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	signedToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claim).SignedString([]byte(secrete_key))
	signedRefreshToken, err = jwt.NewWithClaims(jwt.SigningMethodES256, refreshClaim).SignedString([]byte(secrete_key))

	if err != nil {
		log.Panic(err)
		return
	}

	return signedToken, signedRefreshToken, err
}

func UpdateAllTokens(token string, refrehToken string, userId string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()

	// create the primitive update object for the mongo
	var updateObj primitive.D
	// add the required fields to update the data
	updateObj = append(updateObj, bson.E{"token", token})
	updateObj = append(updateObj, bson.E{"refresh_token", refrehToken})
	updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{"updated_at", updated_at})

	// If true, a new document will be inserted if the filter does not match any documents in the collection
	upsert := true
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}
	_, err := userCollection.UpdateOne(ctx,
		bson.M{"user_id": userId},
		bson.D{{"$set", updateObj}},
		&opt)

	if err != nil {
		log.Panic(err)
		return
	}

	return
}

func ValidateToken(signedToken string) (claim *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(t *jwt.Token) (interface{}, error) { return []byte(secrete_key), nil })

	if err != nil {
		msg = err.Error()
		return
	}

	claim, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = fmt.Sprintf("the token is invalid")
		msg = err.Error()
		return
	}

	if claim.ExpiresAt < time.Now().Local().Unix() {
		msg = fmt.Sprintf("token is expired")
		msg = err.Error()
		return
	}

	return claim, msg
}
