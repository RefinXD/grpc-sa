package places

import (
	context "context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/metadata"
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

type HexId struct {
    ID primitive.ObjectID `bson:"_id"`
}

type UserRoleJson struct {
	Role string `json:"role"`;
}

func (p placesServer) mustEmbedUnimplementedPlaceServiceServer() {}

func (p placesServer) UploadPlaceInfo(ctx context.Context,req *Place) (*Place, error) {
	//Verify user role
	md,_ := metadata.FromIncomingContext(ctx)
	token := "Bearer " + strings.Split(md.Get("authorization")[0], " ")[1]
	fmt.Println(token)

	newReq, err := http.NewRequest("GET", "http://user-service:8081/verifyuserdetail", nil)
	newReq.Header.Add("Authorization", token)
	client := &http.Client{}
    resp, err := client.Do(newReq)
    if err != nil {
		
        log.Println("Error on response.\n[ERROR] -", err)
    }
    defer resp.Body.Close()
	var target UserRoleJson;
	bodyBytes, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Fatal(err)
    }
    json.Unmarshal(bodyBytes,&target)
    if err != nil {
        log.Println("Error while reading the response bytes:", err)
    }
    log.Println(target.Role)
	if (target.Role != "OWNER"){
		err = errors.New("User not authorized to do this action")
		return nil,err
	}
	
	filter := bson.D{{"name",req.Name}}	
	fmt.Println(filter)
	rep := p.con.PlacesCollection.FindOne(ctx,filter)
	var repRes Place
	decodeErr := rep.Decode(&repRes)
	if decodeErr == nil{
		println("Duplicate name")
		return nil,err
	}

	placesResult, err := p.con.PlacesCollection.InsertOne(ctx,req)
	if err != nil{
		log.Fatal(err)
	}
	fmt.Println(req)
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
	//var wg sync.WaitGroup
	//defer wg.Done()

	//Verify user role
	md,_ := metadata.FromIncomingContext(ctx)
	fmt.Println(md)
	token := "Bearer " + strings.Split(md.Get("authorization")[0], " ")[1]
	fmt.Println(token)
	newReq, err := http.NewRequest("GET", "http://user-service:8081/verifyuserdetail", nil)
	newReq.Header.Add("Authorization", token)
	client := &http.Client{}
	fmt.Println(1)
    resp, err := client.Do(newReq)
    if err != nil {
		
         log.Println("Error on response.\n[ERROR] -", err)
		 }
    defer resp.Body.Close()
	var target UserRoleJson;
	bodyBytes, err := io.ReadAll(resp.Body)
    if err != nil {
         log.Fatal(err)
     }
    json.Unmarshal(bodyBytes,&target)
    if err != nil {
         log.Println("Error while reading the response bytes:", err)
    }
    log.Println(target.Role)
	if (target.Role != "OWNER"){
		if req.NewInfo.Capacity != 0 ||  req.NewInfo.Name != "" || req.NewInfo.Facilities != nil{
		err = errors.New("User not authorized to do this action")
		return nil,err
		}
	}
	






	filter := bson.D{{"name",req.TargetName}}
	fmt.Println(filter)
	rep := p.con.PlacesCollection.FindOne(ctx,filter)
	var repRes Place
	err = rep.Decode(&repRes)
	// log.Println(repRes)
	if err != nil{
		println("Place Not Found")
		return nil,err
	}

	fields := bson.M{}
	if req.NewInfo.Facilities == nil{
		fields["facilities"] = repRes.Facilities;
	}else{
		fields["facilities"] = req.NewInfo.Facilities
	}

	if &req.NewInfo.AvailableSeat == nil{
		fields["availableseat"] = repRes.AvailableSeat;
	}else{
		fields["availableseat"] = req.NewInfo.AvailableSeat
	}
	if req.NewInfo.Capacity == 0{
		fields["capacity"] = repRes.Capacity;
	}else{
		fields["capacity"] = req.NewInfo.Capacity
	}

	// if req.NewInfo.Name == nil{
	// 	fields["name"] = repRes.Name;
	// }else{
	// 	fields["name"] = req.NewInfo.Name
	// }
	
	object := bson.M{
		"$set":fields,
	}
	placesResult, err := p.con.PlacesCollection.UpdateOne(ctx,filter,object)
	if err != nil{
		log.Println("what")
		log.Fatal(err)
	}
	log.Println("2")
	log.Println(placesResult)
	log.Println("3")
	res := Place {
		Id: repRes.Id,
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
	fmt.Println(req)
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
	fmt.Println("Name",req.Name)
	filter := bson.M{"name": primitive.Regex{Pattern: req.Name, Options: ""}}
	placesResult,err := p.con.PlacesCollection.Find(ctx,filter);
	if err != nil{
		log.Fatal(err)
	}
	var res PlaceList
	for placesResult.Next(context.TODO()) {
		fmt.Println(placesResult)
		var temp Place;
		var id HexId;
		if err := placesResult.Decode(&temp); err != nil {
			log.Fatal(err)
		}
		if err := placesResult.Decode(&id); err != nil {
			log.Fatal(err)
		}
		temp.Id = id.ID.Hex()
		res.Place = append(res.Place, &temp)
	}
	fmt.Println("search done")
	if err := placesResult.Err(); err != nil {
		log.Fatal(err)
	}
	return &res,nil
}

func (p placesServer) SearchPlacesByOwner(ctx context.Context, req *OwnerName) (*PlaceList, error) {
	fmt.Println("Name",req.OwnerName)
	filter := bson.M{"owner":req.OwnerName}
	placesResult,err := p.con.PlacesCollection.Find(ctx,filter);
	if err != nil{
		log.Fatal(err)
	}
	var res PlaceList
	for placesResult.Next(context.TODO()) {
		fmt.Println(placesResult)
		var temp Place;
		var id HexId;
		if err := placesResult.Decode(&temp); err != nil {
			log.Fatal(err)
		}
		if err := placesResult.Decode(&id); err != nil {
			log.Fatal(err)
		}
		temp.Id = id.ID.Hex()
		res.Place = append(res.Place, &temp)
	}
	fmt.Println("search done")
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
		filter = bson.M{"facilities":bson.M{"$in" :req.Facilities},"availableseat":bson.M{"$gte":req.MinCapacity}}
	}
	fmt.Println(filter)
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


	//Verify user role
	md,_ := metadata.FromIncomingContext(ctx)
	token := "Bearer " + strings.Split(md.Get("authorization")[0], " ")[1]
	//fmt.Println(token)
	newReq, err := http.NewRequest("GET", "http://user-service:8081/verifyuserdetail", nil)
	newReq.Header.Add("Authorization", token)
	client := &http.Client{}
	fmt.Println(1)
    resp, err := client.Do(newReq)
    if err != nil {
		
        log.Println("Error on response.\n[ERROR] -", err)
    }
    defer resp.Body.Close()
	var target UserRoleJson;
	bodyBytes, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Fatal(err)
    }
    json.Unmarshal(bodyBytes,&target)
    if err != nil {
        log.Println("Error while reading the response bytes:", err)
    }
    log.Println(target.Role)
	if (target.Role != "OWNER"){
		err = errors.New("User not authorized to do this action")
		return nil,err
	}
	

	filter := bson.D{{"name",req.Name}}
	rep := p.con.PlacesCollection.FindOne(ctx,filter)
	var repRes Place
	err = rep.Decode(&repRes)
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