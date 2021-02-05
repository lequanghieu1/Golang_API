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

func GetAllRoles(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	collection := db.ConnectRolesAccess()

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

func CreateRolesAccess(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	creater, err := jwt.ExtractUsernameFromToken(r)

	if err != nil {
		res.JSON(w, 500, "Internal Server Error")
		return
	}

	id_extra := r.PostFormValue("id_extra")
	name_model := r.PostFormValue("name_model")
	add := r.PostFormValue("add")
	read := r.PostFormValue("read")
	update := r.PostFormValue("update")
	delete := r.PostFormValue("delete")

	if govalidator.IsNull(id_extra) || govalidator.IsNull(name_model) || govalidator.IsNull(add) || govalidator.IsNull(read) || govalidator.IsNull(update) || govalidator.IsNull(delete) {
		res.JSON(w, 400, "Data can not empty")
		return
	}
	var result bson.M
	collections := db.ConnectUsers()
	errFindUsername := collections.FindOne(context.TODO(), bson.M{"username": id_extra}).Decode(&result)

	if errFindUsername != nil {
		res.JSON(w, 409, "id_extra don't match")
		return
	}
	id_extra = models.Santize(id_extra)
	name_model = models.Santize(name_model)
	add = models.Santize(add)
	read = models.Santize(read)
	update = models.Santize(update)
	delete = models.Santize(delete)

	collection := db.ConnectRolesAccess()
	errFindId := collection.FindOne(context.TODO(), bson.M{"id_extra": id_extra}).Decode(&result)

	if errFindId == nil {
		res.JSON(w, 409, "id_extra is exists")
		return
	}
	newRoles := bson.M{"id_extra": id_extra, "name_model": name_model, "add": add, "read": read, "update": update, "delete": delete, "created_at": time.Now(), "updated_at": time.Now(), "creater": creater}

	_, errs := collection.InsertOne(context.TODO(), newRoles)

	if errs != nil {
		res.JSON(w, 500, "Create roles access has failed")
		return
	}

	res.JSON(w, 201, "Create Succesfully")

}

func EditRolesAccess(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	name_model := r.PostFormValue("name_model")
	add := r.PostFormValue("add")
	read := r.PostFormValue("read")
	update := r.PostFormValue("update")
	delete := r.PostFormValue("delete")
	username, err := jwt.ExtractUsernameFromToken(r)

	if err != nil {
		res.JSON(w, 500, "Internal Server Error")
		return
	}

	collection := db.ConnectRolesAccess()

	var result bson.M
	errFind := collection.FindOne(context.TODO(), bson.M{"id_extra": id}).Decode(&result)
	fmt.Println(errFind)
	if errFind != nil {
		res.JSON(w, 404, "Roles access Not Found")
		return
	}

	creater := fmt.Sprintf("%v", result["creater"])

	if username != creater {
		res.JSON(w, 403, "Permission Denied")
		return
	}

	filter := bson.M{"id_extra": id}
	updates := bson.M{"$set": bson.M{"name_model": name_model, "add": add, "read": read, "update": update, "delete": delete}}

	_, errUpdate := collection.UpdateOne(context.TODO(), filter, updates)

	if errUpdate != nil {
		res.JSON(w, 500, "Edit has failed")
		return
	}

	res.JSON(w, 200, "Edit Successfully")

}

func DeleteRolesAccess(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	username, err := jwt.ExtractUsernameFromToken(r)
	collection := db.ConnectRolesAccess()

	if err != nil {
		res.JSON(w, 500, "Internal Server Error")
		return
	}

	var result bson.M
	errFind := collection.FindOne(context.TODO(), bson.M{"id_extra": id}).Decode(&result)

	if errFind != nil {
		res.JSON(w, 404, "Roles access Not Found")
		return
	}

	creater := fmt.Sprintf("%v", result["creater"])

	if username != creater {
		res.JSON(w, 403, "Permission Denied")
		return
	}

	errDelete := collection.FindOneAndDelete(context.TODO(), bson.M{"id_extra": id}).Decode(&result)

	if errDelete != nil {
		res.JSON(w, 500, "Delete has failed")
		return
	}

	res.JSON(w, 200, "Delete Successfully")

}
