package main

import (
	"fmt"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"io/ioutil"
	"gopkg.in/mgo.v2"
	// "os"
)

var session *mgo.Session

func Create(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
	var userInfo Request
	json.NewDecoder(req.Body).Decode(&userInfo)
	response := Response{}
	response.Id = bson.NewObjectId()
	response.Name = userInfo.Name
	response.Address = userInfo.Address
	response.City = userInfo.City
	response.State = userInfo.State
	response.Zip = userInfo.Zip
	requestURL := createURL(response.Address, response.City, response.State)
	getLocation(&response, requestURL)
	session.DB("cmpe273").C("assignment2").Insert(response)

	location, _ := json.Marshal(response)
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(201)
	fmt.Fprintf(rw, "%s", location)
}

func Get(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(404)
		return
	}
	userId := bson.ObjectIdHex(id)
	response := Response{}
	if err := session.DB("cmpe273").C("assignment2").FindId(userId).One(&response); err != nil{
		w.WriteHeader(404)
		return
	}
	location, _ := json.Marshal(response)
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", location)
}

func Update(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(404)
		return
	}
	var userInfo Request
	json.NewDecoder(r.Body).Decode(&userInfo)
	response := Response{}
	response.Id = bson.ObjectIdHex(id)
	response.Name = userInfo.Name
	response.Address = userInfo.Address
	response.City = userInfo.City
	response.State = userInfo.State
	response.Zip = userInfo.Zip
	requestURL := createURL(response.Address, response.City, response.State)
	getLocation(&response, requestURL)

	if err := session.DB("cmpe273").C("assignment2").Update(bson.M{"_id": response.Id}, bson.M{"$set":bson.M{"address" : response.Address, "city" : response.City, "state" : response.State, "zip" : response.Zip, "coordinate.lat" : response.Coordinate.Lat, "coordinate.lng" : response.Coordinate.Lng}}); err != nil {
		w.WriteHeader(404)
		return
	}

	if err := session.DB("cmpe273").C("assignment2").FindId(response.Id).One(&response); err != nil {
		w.WriteHeader(404)
		return
	}

	location, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s", location)
}

func Delete(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(404)
		return
	}
	userId := bson.ObjectIdHex(id)

	if err := session.DB("cmpe273").C("assignment2").RemoveId(userId); err != nil {
		w.WriteHeader(404)
		return
	}
	w.WriteHeader(200)
}

func main() {

	session, _ = mgo.Dial("mongodb://admin:12345@ds051838.mongolab.com:51838/cmpe273")

	defer session.Close()

	mux := httprouter.New()
	//connectDB()
	mux.GET("/locations/:id", Get)
    mux.POST("/locations/", Create)
    mux.DELETE("/locations/:id", Delete)
    mux.PUT("/locations/:id", Update)
    server := http.Server{
            Addr:        "0.0.0.0:8080",
            Handler: mux,
    }
    server.ListenAndServe()	
}


func createURL(address string, city string, state string) string{
	var getURL string
	spStr := strings.Split(address, " ")
	for i := 0; i < len(spStr); i++ {
		if i == 0 {
			getURL = spStr[i] + "+"
		}else if i == len(spStr) - 1 {
			getURL = getURL + spStr[i] + ","
		}else{
			getURL = getURL + spStr[i] + "+"
		}
	}
	spStr = strings.Split(city, " ")
	for i := 0; i < len(spStr); i++ {
		if i == 0{
			getURL = getURL + "+" + spStr[i]
		}else if i == len(spStr) - 1 {
			getURL = getURL + "+" + spStr[i] + ","
		}else{
			getURL = getURL + "+" + spStr[i]
		}
	}
	getURL = getURL + "+" + state
	return getURL
}

func getLocation(response *Response, formatURL string) {
	urlLeft := "http://maps.google.com/maps/api/geocode/json?address="
	urlRight := "&sensor=false"
	urlFormat := urlLeft + formatURL + urlRight

	getLocation, err := http.Get(urlFormat)
	if err != nil{
		fmt.Println("Get Location Error", err)
		panic(err)
	}
	
	body, err := ioutil.ReadAll(getLocation.Body)
	if err != nil{
		fmt.Println("Get Location Error", err)
		panic(err)
	}

	var data MyJsonName
	byt := []byte(body)
	if err := json.Unmarshal(byt, &data); err != nil{
		panic(err)
	}
	response.Coordinate.Lat = data.Results[0].Geometry.Location.Lat
	response.Coordinate.Lng = data.Results[0].Geometry.Location.Lng
}