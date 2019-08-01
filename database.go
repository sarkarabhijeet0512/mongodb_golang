package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type user struct {
	Name string
	Age  int
	City string
	Dob  string
}

func createnew() interface{} {
	var name, dob, home string
	var age int

	fmt.Println("Enter name: ")
	fmt.Scan(&name)
	fmt.Println("Enter age: ")
	fmt.Scan(&age)
	fmt.Println("Enter date of birth (dd-mm-yyyy format): ")
	fmt.Scan(&dob)
	fmt.Println("Enter your hometown: ")
	fmt.Scan(&home)

	people := user{Name: name, Age: age, Dob: dob, City: home}
	return people
}

func insert(client *mongo.Client, collection *mongo.Collection) {

	persons := createnew()

	insertResult, err := collection.InsertOne(context.TODO(), persons)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted a Single Document: ", insertResult.InsertedID)
	db(client)
}

func db(client *mongo.Client) {
	collection := client.Database("mydb").Collection("persons")
	var choice string
	fmt.Println("what operation do you want to perform (insert/delete/update/get/exit)")
	_, err := fmt.Scan(&choice)
	if err != nil {
		fmt.Println("Invalid format;", err)
		os.Exit(0)
	}

	switch choice {
	case "insert":
		insert(client, collection)

	case "delete":
		delete(client, collection)

	case "update":
		update(client, collection)

	case "get":
		search(client, collection)
	case "exit":
		os.Exit(0)

	default:
		fmt.Println("invalid input\nuse appropriate values")
	}
}

func main() {
	clientOptions := options.Client().ApplyURI("mongodb+srv://Abhi_0512:12345ok@db-xhiqs.mongodb.net/test")
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	db(client)

}

func delete(client *mongo.Client, collection *mongo.Collection) {
	var objID string
	var person user
	fmt.Println("enter id")
	fmt.Scan(&objID)
	id, _ := primitive.ObjectIDFromHex(objID)
	result := collection.FindOneAndDelete(context.TODO(), bson.M{"_id": id})
	result.Decode(&person)
	if person.Age != 0 {
		fmt.Println(person, " deleted successfully")
	} else {
		fmt.Println("delete failed wrong id")

	}
	db(client)
}

func search(client *mongo.Client, collection *mongo.Collection) {
	var objID string
	var person user
	fmt.Print("Enter the id of the person: ")
	fmt.Scan(&objID)
	id, _ := primitive.ObjectIDFromHex(objID)
	result := collection.FindOne(context.TODO(), bson.M{"_id": id})
	result.Decode(&person)
	if person.Age != 0 {
		fmt.Println(person, " found successfully")
	} else {
		fmt.Println("Search failed")
	}
	db(client)
}

func update(client *mongo.Client, collection *mongo.Collection) {
	var objID string
	var people user
	fmt.Println("Enter the id of the person: ")
	fmt.Scan(&objID)
	id, _ := primitive.ObjectIDFromHex(objID)
	result := collection.FindOne(context.TODO(), bson.M{"_id": id})
	result.Decode(&people)
	if people.Age != 0 {
		fmt.Println(people, " found successfully")
		people = updaterequire(client, collection, people)
		collection.UpdateOne(context.TODO(), bson.M{"_id": id}, people)
		fmt.Println(people, "updated successfully")
	} else {
		fmt.Println("Enter valid id")
	}
	db(client)
}

func updaterequire(client *mongo.Client, collection *mongo.Collection, person user) user {
	var name, dob, home string
	var age int
	var field string

	fmt.Println("Enter the field you want to update: \n Name \n Age \n dob \n home")
	fmt.Println("Enter one of the above fields: ")
	fmt.Scan(&field)

	switch field {
	case "name":
		fmt.Print("Enter new name: ")
		fmt.Scan(&name)
		person.Name = name
	case "age":
		fmt.Print("Enter new age: ")
		fmt.Scan(&age)
		person.Age = age
	case "dob":
		fmt.Print("Enter new dob (in dd-mm-yyyy format): ")
		fmt.Scan(&dob)
		person.Dob = dob
	case "home":
		fmt.Print("Enter new home: ")
		fmt.Scan(&home)
		person.City = home

	default:
		fmt.Println("enter valid values")
	}
	return person
}
