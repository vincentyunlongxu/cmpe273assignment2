package main

import "gopkg.in/mgo.v2/bson"

type Request struct{
	Name	string	`json:"name"`
	Address	string	`json:"address"`
	City	string	`json:"city"`
	State	string	`json:"state"`
	Zip		string	`json:"zip"`
}

type Response struct{
	Id		bson.ObjectId	`json:"id" bson:"_id"`
	Name 	string			`json:"name" bson:"name"`
	Address string			`json:"address" bson:"address"`
	City 	string			`json:"city" bson:"city"`
	State 	string			`json:"state" bson:"state"`
	Zip 	string			`json:"zip" bson:"zip"`
	Coordinate struct{
		Lat float64			`json:"lat" bson:"lat"`
		Lng	float64			`json:"lng" bson:"lng"`
	}						`json:"coordinate" bson:"coordinate"`
}