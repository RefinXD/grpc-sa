package places

import (
	context "context"
	"fmt"

	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
	// "go.mongodb.org/mongo-driver/mongo/readpref"
)


type Connection struct{
	PlacesCollection *mongo.Collection
}

type placesServer struct {
	con Connection
}

func NewPlacesServer(con Connection) PlaceServiceServer{
	return placesServer{con:con}
}

func (p placesServer) mustEmbedUnimplementedPlaceServiceServer() {}

func (p placesServer) UploadPlaceInfo(ctx context.Context,req *Place) (*Place, error) {
	filter := bson.D{{"name",req.Name}}
	rep := p.con.PlacesCollection.FindOne(ctx,filter)
	var repRes Place
	err := rep.Decode(&repRes)
	if err == nil{
		println("Duplicate name")
		return nil,err
	}

	placesResult, err := p.con.PlacesCollection.InsertOne(ctx,req)
	if err != nil{
		log.Fatal(err)
	}
	//result := req.Name
	id := placesResult.InsertedID.(primitive.ObjectID)
	res := Place {
		Id: id.String(),
		Name: req.Name,
		Owner: req.Owner,
		Capacity: req.Capacity,
		AvailableSeat: req.Capacity,
		Facilities: req.Facilities,	
	}
	return &res,nil
}

func (p placesServer) UpdatePlace(ctx context.Context,req *UpdatePlace) (*Place, error) {
	filter := bson.D{{"name",req.TargetName}}
	object := bson.M{
        "$set": req.NewInfo,
    }
	rep := p.con.PlacesCollection.FindOne(ctx,filter)
	var repRes Place
	err := rep.Decode(&repRes)
	log.Println(repRes)
	if err != nil{
		println("Place Not Found")
		return nil,err
	}
	placesResult, err := p.con.PlacesCollection.UpdateOne(ctx,filter,object)
	if err != nil{
		log.Fatal(err)
	}
	print(placesResult.UpsertedID)
	res := Place {
		Id: placesResult.UpsertedID.(primitive.ObjectID).String(),
		Name: req.NewInfo.Name,
		Owner: req.NewInfo.Owner,
		AvailableSeat: req.NewInfo.AvailableSeat,
		Capacity: req.NewInfo.Capacity,
		Facilities: req.NewInfo.Facilities,
	}
	fmt.Println(res)
	return &res,nil
}

func (p placesServer) GetPlaceInfo(ctx context.Context,req *PlaceId) (*Place, error) {
	var place Place
	id,err := primitive.ObjectIDFromHex(req.Id)
	if err != nil{
		log.Fatal(err)
	}
	filter := bson.M{"_id":id}
	placesResult := p.con.PlacesCollection.FindOne(ctx,filter)
	log.Println(placesResult.Decode(&place))
	
	return &place,nil
}

func (p placesServer) SearchPlaces(ctx context.Context,req *PlaceName) (*PlaceList, error) {
	filter := bson.M{"name": primitive.Regex{Pattern: req.Name, Options: ""}}
	placesResult,err := p.con.PlacesCollection.Find(ctx,filter)
	if err != nil{
		log.Fatal(err)
	}
	var res PlaceList
	for placesResult.Next(context.TODO()) {
		var result Place
		if err := placesResult.Decode(&result); err != nil {
			log.Fatal(err)
		}
		res.Place = append(res.Place, &result)
	}
	if err := placesResult.Err(); err != nil {
		log.Fatal(err)
	}
	return &res,nil
}

func (p placesServer) FilterPlaces(ctx context.Context,req *Filter) (*PlaceList, error) {
	filter := bson.M{}
	if len(req.Facilities) == 0{
		filter = bson.M{"availableseat":bson.M{"$gte":req.MinCapacity}}
	} else{
		filter = bson.M{"facilities":bson.M{"$in" :req.Facilities},"availableseat":bson.M{"$gt":req.MinCapacity}}
	}
	placesResult,err := p.con.PlacesCollection.Find(ctx,filter)
	var test PlaceList
	placesResult.All(ctx,test.Place)
	if err != nil{
		log.Fatal(err)
	}
	var res PlaceList
	for placesResult.Next(context.TODO()) {
		var result Place
		if err := placesResult.Decode(&result); err != nil {
			log.Fatal(err)
		}
		res.Place = append(res.Place, &result)
	}
	if err := placesResult.Err(); err != nil {
		log.Fatal(err)
	}
	
	return &res,nil
}

func (p placesServer) RemovePlaces(ctx context.Context,req *PlaceName) (*Empty, error) {
	filter := bson.D{{"name",req.Name}}
	rep := p.con.PlacesCollection.FindOne(ctx,filter)
	var repRes Place
	err := rep.Decode(&repRes)
	if err != nil{
		println("No place found with given name")
		return nil,err
	}
	placesResult, err := p.con.PlacesCollection.DeleteOne(ctx,filter)
	if err != nil{
		log.Fatal(err)
	}
	var res Empty
	print(placesResult)
	return &res,nil
}