package main

import (
	"client/places"
	"encoding/json"
	//"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)
var creds = insecure.NewCredentials()
var	cc,err = grpc.Dial("localhost:9000",grpc.WithTransportCredentials(creds))


var placesClient = places.NewPlaceServiceClient(cc)
var	placesService = places.NewPlaceService(placesClient)


func main(){

	//placesService.FilterPlaces(places.Filter{Facilities:[]string{"toilet"}})

	r := http.NewServeMux()
	r.HandleFunc("/update", updateHandler)
	r.HandleFunc("/upload", uploadHandler)
	r.HandleFunc("/filter", filterHandler)
	r.HandleFunc("/search", searchHandler)
	r.HandleFunc("/delete", deleteHandler)
	server := &http.Server{
		Addr: ":8080",
		Handler:r,
	}
	if err:= server.ListenAndServe();err != nil{
		panic(err)
	}
}
func updateHandler(w http.ResponseWriter, r *http.Request) {
	if (r.Method != http.MethodPatch){
		http.Error(w,"Method not allowed",http.StatusMethodNotAllowed)
	}
	var place places.UpdatePlace

	respBody, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(respBody, &place)
	if err != nil{
		log.Fatal(err)
	}
	log.Println(place)
	res,err:= placesService.UpdatePlace(place);
	jsonBytes,err := json.Marshal(res);
	
	w.Header().Set("Content-Type", "application/json")

	w.Write(jsonBytes)

}

func filterHandler(w http.ResponseWriter, r *http.Request) {

	if (r.Method != http.MethodGet){
		http.Error(w,"Method not allowed",http.StatusMethodNotAllowed)
	}
	//log.Println(r.Body)
	var filter places.Filter
	err := json.NewDecoder(r.Body).Decode(&filter)
	log.Println(filter)
	if err != nil{
		log.Fatal(err)
	}
	res,err := placesService.FilterPlaces(filter);
	jsonBytes,err := json.Marshal(res);
	
	w.Header().Set("Content-Type", "application/json")

	w.Write(jsonBytes)

}

func deleteHandler(w http.ResponseWriter, r *http.Request) {

	if (r.Method != http.MethodDelete){
		http.Error(w,"Method not allowed",http.StatusMethodNotAllowed)
	}
	var name places.PlaceName
	respBody, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(respBody, &name)
	if err != nil{
		log.Fatal(err)
	}
	res,err := placesService.RemovePlaces(name);
	jsonBytes,err := json.Marshal(res);
	
	w.Header().Set("Content-Type", "application/json")

	w.Write(jsonBytes)

}

func searchHandler(w http.ResponseWriter, r *http.Request) {

	if (r.Method != http.MethodGet){
		http.Error(w,"Method not allowed",http.StatusMethodNotAllowed)
	}
	var name places.PlaceName
	respBody, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(respBody, &name)
	if err != nil{
		log.Fatal(err)
	}
	res,err:= placesService.SearchPlaces(name);
	jsonBytes,err := json.Marshal(res);
	
	w.Header().Set("Content-Type", "application/json")

	w.Write(jsonBytes)

}
func uploadHandler(w http.ResponseWriter, r *http.Request) {

	if (r.Method != http.MethodPost){
		http.Error(w,"Method not allowed",http.StatusMethodNotAllowed)
	}
	var place places.Place
	respBody, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(respBody, &place)
	if err != nil{
		log.Fatal(err)
	}
	res,err := placesService.UploadPlaceInfo(place);
	jsonBytes,err := json.Marshal(res);
	
	w.Header().Set("Content-Type", "application/json")

	w.Write(jsonBytes)

}
