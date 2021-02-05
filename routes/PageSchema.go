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

func GetAllPages(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	collection := db.ConnectPageSchemas()

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

func CreatePageSchema(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	creater, err := jwt.ExtractUsernameFromToken(r)

	if err != nil {
		res.JSON(w, 500, "Internal Server Error")
		return
	}

	page := r.PostFormValue("page")
	key := r.PostFormValue("key")
	label := r.PostFormValue("label")

	if govalidator.IsNull(page) || govalidator.IsNull(key) || govalidator.IsNull(label) {
		res.JSON(w, 400, "Data can not empty")
		return
	}

	page = models.Santize(page)
	key = models.Santize(key)
	label = models.Santize(label)

	collection := db.ConnectPageSchemas()
	var result bson.M
	errFindUsername := collection.FindOne(context.TODO(), bson.M{"key": key}).Decode(&result)

	if errFindUsername == nil {
		res.JSON(w, 409, "Page does exists")
		return
	}

	newPage := bson.M{"page": page, "key": key, "label": label, "sortable": false, "selected": false, "created_at": time.Now(), "updated_at": time.Now(), "creater": creater}

	_, errs := collection.InsertOne(context.TODO(), newPage)

	if errs != nil {
		res.JSON(w, 500, "Create page schema has failed")
		return
	}

	res.JSON(w, 201, "Create Succesfully")

}

func EditPageSchema(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	page := r.PostFormValue("page")
	key := r.PostFormValue("key")
	label := r.PostFormValue("label")
	sortable := r.PostFormValue("sortable")
	selected := r.PostFormValue("selected")
	username, err := jwt.ExtractUsernameFromToken(r)

	if err != nil {
		res.JSON(w, 500, "Internal Server Error")
		return
	}

	collection := db.ConnectPageSchemas()

	var result bson.M
	errFind := collection.FindOne(context.TODO(), bson.M{"key": id}).Decode(&result)
	if errFind != nil {
		res.JSON(w, 404, "page schema Not Found")
		return
	}

	creater := fmt.Sprintf("%v", result["creater"])
	oldpage := fmt.Sprintf("%v", result["page"])
	oldkey := fmt.Sprintf("%v", result["key"])
	oldlabel := fmt.Sprintf("%v", result["label"])
	oldsortable := fmt.Sprintf("%v", result["sortable"])
	oldselected := fmt.Sprintf("%v", result["selected"])

	if username != creater {
		res.JSON(w, 403, "Permission Denied")
		return
	}
	if page == "" {
		page = oldpage
	}
	if key == "" {
		key = oldkey
	}
	if label == "" {
		label = oldlabel
	}
	if sortable == "" {
		sortable = oldsortable
	}
	if selected == "" {
		selected = oldselected
	}

	filter := bson.M{"key": id}
	fmt.Println(result)
	update := bson.M{"$set": bson.M{"page": page, "key": key, "label": label, "sortable": sortable, "selected": selected}}

	_, errUpdate := collection.UpdateOne(context.TODO(), filter, update)

	if errUpdate != nil {
		res.JSON(w, 500, "Edit has failed")
		return
	}

	res.JSON(w, 200, "Edit Successfully")

}

func DeletePageSchema(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	username, err := jwt.ExtractUsernameFromToken(r)
	collection := db.ConnectPageSchemas()

	if err != nil {
		res.JSON(w, 500, "Internal Server Error")
		return
	}

	var result bson.M
	errFind := collection.FindOne(context.TODO(), bson.M{"key": id}).Decode(&result)

	if errFind != nil {
		res.JSON(w, 404, "Page Schema Not Found")
		return
	}

	creater := fmt.Sprintf("%v", result["creater"])

	if username != creater {
		res.JSON(w, 403, "Permission Denied")
		return
	}

	errDelete := collection.FindOneAndDelete(context.TODO(), bson.M{"key": id}).Decode(&result)

	if errDelete != nil {
		res.JSON(w, 500, "Delete has failed")
		return
	}

	res.JSON(w, 200, "Delete Successfully")

}
