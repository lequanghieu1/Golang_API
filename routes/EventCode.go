package routes

import (
	"context"
	"fmt"
	"net/http"
	_ "reflect"

	"github.com/asaskevich/govalidator"
	jwt "github.com/conglt10/web-golang/auth"
	db "github.com/conglt10/web-golang/database"
	"github.com/conglt10/web-golang/models"
	res "github.com/conglt10/web-golang/utils"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
)

func GetAllEvents(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	collection := db.ConnectEventCodes()

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

func CreateEventCode(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	creater, err := jwt.ExtractUsernameFromToken(r)

	if err != nil {
		res.JSON(w, 500, "Internal Server Error")
		return
	}

	name := r.PostFormValue("name")
	code := r.PostFormValue("code")

	if govalidator.IsNull(name) || govalidator.IsNull(code) {
		res.JSON(w, 400, "Data can not empty")
		return
	}

	name = models.Santize(name)
	code = models.Santize(code)

	collection := db.ConnectEventCodes()
	var result bson.M
	errFindUsername := collection.FindOne(context.TODO(), bson.M{"code": code}).Decode(&result)

	if errFindUsername == nil {
		res.JSON(w, 409, "code does exists")
		return
	}

	newEvC := bson.M{"name": name, "code": code, "creater": creater}

	_, errs := collection.InsertOne(context.TODO(), newEvC)

	if errs != nil {
		res.JSON(w, 500, "Create event code has failed")
		return
	}

	res.JSON(w, 201, "Create Succesfully")

}

func EditEventCode(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	username, err := jwt.ExtractUsernameFromToken(r)
	id := ps.ByName("id")
	name := r.PostFormValue("name")
	code := r.PostFormValue("code")

	if err != nil {
		res.JSON(w, 500, "Internal Server Error")
		return
	}

	collection := db.ConnectEventCodes()

	var result bson.M
	errFind := collection.FindOne(context.TODO(), bson.M{"name": id}).Decode(&result)
	fmt.Println(errFind)
	if errFind != nil {
		res.JSON(w, 404, "Event code Not Found")
		return
	}
	creater := fmt.Sprintf("%v", result["creater"])
	oldname := fmt.Sprintf("%v", result["name"])
	oldcode := fmt.Sprintf("%v", result["code"])
	if username != creater {
		res.JSON(w, 403, "Permission Denied")
		return
	}
	if name == "" {
		name = oldname
	}
	if code == "" {
		code = oldcode
	}
	filter := bson.M{"name": id}
	update := bson.M{"$set": bson.M{"name": name, "code": code}}

	_, errUpdate := collection.UpdateOne(context.TODO(), filter, update)

	if errUpdate != nil {
		res.JSON(w, 500, "Edit has failed")
		return
	}

	res.JSON(w, 200, "Edit Successfully")

}

func DeleteEventCode(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	username, err := jwt.ExtractUsernameFromToken(r)
	collection := db.ConnectEventCodes()

	if err != nil {
		res.JSON(w, 500, "Internal Server Error")
		return
	}

	var result bson.M
	errFind := collection.FindOne(context.TODO(), bson.M{"name": id}).Decode(&result)

	if errFind != nil {
		res.JSON(w, 404, "User Not Found")
		return
	}

	creater := fmt.Sprintf("%v", result["creater"])

	if username != creater {
		res.JSON(w, 403, "Permission Denied")
		return
	}

	errDelete := collection.FindOneAndDelete(context.TODO(), bson.M{"name": id}).Decode(&result)

	if errDelete != nil {
		res.JSON(w, 500, "Delete has failed")
		return
	}

	res.JSON(w, 200, "Delete Successfully")

}
