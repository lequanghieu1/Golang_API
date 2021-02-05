package routes

import (
	"context"
	"fmt"
	"net/http"
	_ "reflect"
	"time"

	"github.com/asaskevich/govalidator"
	jwt "github.com/conglt10/web-golang/auth"
	db "github.com/conglt10/web-golang/database"
	"github.com/conglt10/web-golang/models"
	res "github.com/conglt10/web-golang/utils"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
)

func GetAllUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	collection := db.ConnectUsers()

	var result []bson.M
	data, err := collection.Find(context.Background(), bson.M{})

	if err != nil {
		res.JSON(w, 500, "Internal Server Error")
		return
	}

	defer data.Close(context.Background())
	for data.Next(context.Background()) {
		var elem bson.M
		err := data.Decode(&elem)

		if err != nil {
			res.JSON(w, 500, "Internal Server Error")
			return
		}

		result = append(result, elem)
	}

	res.JSON(w, 200, result)
}

func CreateManager(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	creater, err := jwt.ExtractUsernameFromToken(r)

	if err != nil {
		res.JSON(w, 500, "Internal Server Error")
		return
	}

	username := r.PostFormValue("username")
	fullname := r.PostFormValue("fullname")
	password := r.PostFormValue("password")

	if govalidator.IsNull(username) || govalidator.IsNull(fullname) || govalidator.IsNull(password) {
		res.JSON(w, 400, "Data can not empty")
		return
	}

	username = models.Santize(username)
	fullname = models.Santize(fullname)
	password = models.Santize(password)

	collection := db.ConnectUsers()
	var result bson.M
	errFindUsername := collection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&result)

	if errFindUsername == nil {
		res.JSON(w, 409, "User does exists")
		return
	}
	password, err = models.Hash(password)

	if err != nil {
		res.JSON(w, 500, "Register has failed")
		return
	}
	newUser := bson.M{"username": username, "fullname": fullname, "password": password, "is_login": false, "created_at": time.Now(), "updated_at": time.Now(), "creater": creater}

	_, errs := collection.InsertOne(context.TODO(), newUser)

	if errs != nil {
		res.JSON(w, 500, "Create manager has failed")
		return
	}

	res.JSON(w, 201, "Create Succesfully")

}

func EditManager(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	fullname := r.PostFormValue("fullname")
	password := r.PostFormValue("password")
	username, err := jwt.ExtractUsernameFromToken(r)

	if err != nil {
		res.JSON(w, 500, "Internal Server Error")
		return
	}

	collection := db.ConnectUsers()

	var result bson.M
	errFind := collection.FindOne(context.TODO(), bson.M{"username": id}).Decode(&result)
	fmt.Println(errFind)
	if errFind != nil {
		res.JSON(w, 404, "User Not Found")
		return
	}

	creater := fmt.Sprintf("%v", result["creater"])
	oldfullname := fmt.Sprintf("%v", result["fullname"])
	oldpassword := fmt.Sprintf("%v", result["password"])

	if username != creater {
		res.JSON(w, 403, "Permission Denied")
		return
	}
	if fullname == "" {
		fullname = oldfullname
	}
	if password == "" {
		password = oldpassword
	}
	password, err = models.Hash(password)

	if err != nil {
		res.JSON(w, 500, "Register has failed")
		return
	}
	filter := bson.M{"username": id}
	update := bson.M{"$set": bson.M{"fullname": fullname, "password": password}}

	_, errUpdate := collection.UpdateOne(context.TODO(), filter, update)

	if errUpdate != nil {
		res.JSON(w, 500, "Edit has failed")
		return
	}

	res.JSON(w, 200, "Edit Successfully")

}

func DeleteManager(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	username, err := jwt.ExtractUsernameFromToken(r)
	collection := db.ConnectUsers()

	if err != nil {
		res.JSON(w, 500, "Internal Server Error")
		return
	}

	var result bson.M
	errFind := collection.FindOne(context.TODO(), bson.M{"username": id}).Decode(&result)

	if errFind != nil {
		res.JSON(w, 404, "User Not Found")
		return
	}

	creater := fmt.Sprintf("%v", result["creater"])

	if username != creater {
		res.JSON(w, 403, "Permission Denied")
		return
	}

	errDelete := collection.FindOneAndDelete(context.TODO(), bson.M{"username": id}).Decode(&result)

	if errDelete != nil {
		res.JSON(w, 500, "Delete has failed")
		return
	}

	res.JSON(w, 200, "Delete Successfully")

}
