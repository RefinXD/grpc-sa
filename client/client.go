package main

import (
	"client/places"
	"strings"

	//"context"
	"encoding/json"
	"fmt"

	//"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/rs/cors"
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
	r.HandleFunc("/searchbyowner", searchByOwnerHandler)
	r.HandleFunc("/delete", deleteHandler)
	r.HandleFunc("/info", getHandler)
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},     // Allow any origin
		AllowedHeaders: []string{"*"},     // Allow any headers
		AllowedMethods: []string{"GET", "POST","PATCH", "PUT", "DELETE"},
		Debug:          true,              // Enable debugging (optional)
	}).Handler(r)
	

	// Start the server on port 8080
	fmt.Println("Server listening on :8080")
	http.ListenAndServe(":8080", corsHandler)
}
func updateHandler(w http.ResponseWriter, r *http.Request) {

	if (r.Method != http.MethodPatch){
		http.Error(w,"Method not allowed",http.StatusMethodNotAllowed)
	}
	reqToken := strings.Split(r.Header.Get("Authorization"), " ")[1]
	var place places.UpdatePlace

	respBody, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(respBody, &place)
	if err != nil{
		log.Fatal(err)
	}
	log.Println(place)
	res,err:= placesService.UpdatePlace(place,reqToken);
	if err != nil{
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	jsonBytes,err := json.Marshal(res);
	
	w.Header().Set("Content-Type", "application/json")

	w.Write(jsonBytes)

}

func filterHandler(w http.ResponseWriter, r *http.Request) {

	if (r.Method != http.MethodPost){
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
	if err != nil{
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
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
	reqToken := strings.Split(r.Header.Get("Authorization"), " ")[1]
	res,err := placesService.RemovePlaces(name,reqToken);
	if err != nil{
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	fmt.Println(res)
	fmt.Println("deleted")
	
	w.Header().Set("Content-Type", "application/json")

	w.Write([]byte("Successfully Deleted"))

}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("tests",r.URL.Query())
	if (r.Method != http.MethodGet){
		http.Error(w,"Method not allowed",http.StatusMethodNotAllowed)
	}
	var name places.PlaceName
	respBody, _ := ioutil.ReadAll(r.Body)
	fmt.Println(respBody)
	name.Name= r.URL.Query().Get("name")
	if (len(respBody) == 0){
		respBody= []byte{123, 10, 9, 34, 110, 97, 109, 101, 34, 58, 34, 34, 10, 125, 10}
	}
	err := json.Unmarshal(respBody, &name)
	if err != nil{
		log.Fatal(err)
	}
	res,err:= placesService.SearchPlaces(name);
	fmt.Println("success")
	if err != nil{
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	jsonBytes,err := json.Marshal(res);
	
	w.Header().Set("Content-Type", "application/json")

	w.Write(jsonBytes)

}

func searchByOwnerHandler(w http.ResponseWriter, r *http.Request) {
	if (r.Method != http.MethodGet){
		http.Error(w,"Method not allowed",http.StatusMethodNotAllowed)
	}
	var name places.OwnerName
	respBody, _ := ioutil.ReadAll(r.Body)
	fmt.Println(respBody)
	if (len(respBody) == 0){
		respBody= []byte{123, 10, 9, 34, 110, 97, 109, 101, 34, 58, 34, 34, 10, 125, 10}
	}
	err := json.Unmarshal(respBody, &name)
	if err != nil{
		log.Fatal(err)
	}
	name.OwnerName = r.URL.Query().Get("ownerName")
	res,err:= placesService.SearchPlacesByOwner(name);
	fmt.Println("success")
	if err != nil{
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	jsonBytes,err := json.Marshal(res);
	
	w.Header().Set("Content-Type", "application/json")

	w.Write(jsonBytes)

}





func uploadHandler(w http.ResponseWriter, r *http.Request) {
	
	if (r.Method != http.MethodPost){
		http.Error(w,"Method not allowed",http.StatusMethodNotAllowed)
	}
	req := r.Header

	fmt.Println("Keys:")
    for key := range req {
        fmt.Println(key)
    }
	reqToken := strings.Split(r.Header.Get("Authorization"), " ")[1]
	fmt.Println(reqToken)
	var place places.Place
	fmt.Println(1)
	respBody, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(respBody, &place)
	if err != nil{
		log.Fatal(err)
	}
	fmt.Println(3)
	res,err := placesService.UploadPlaceInfo(place,reqToken);
	fmt.Println(res,err)
	if err != nil || res == nil{

		w.WriteHeader(400)
		if res == nil{
			w.Write([]byte(err.Error()))
		}else{
		w.Write([]byte(err.Error()))
		}
		return
	}
	jsonBytes,err := json.Marshal(res);

	
	w.Header().Set("Content-Type", "application/json")

	w.Write(jsonBytes)

}

func getHandler(w http.ResponseWriter, r *http.Request) {

	if (r.Method != http.MethodGet){
		http.Error(w,"Method not allowed",http.StatusMethodNotAllowed)
	}
	var name places.PlaceId
	respBody, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(respBody, &name)
	if err != nil{
		log.Fatal(err)
	}
	res,err:= placesService.GetPlaceInfo(name);
	if err != nil{
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	jsonBytes,err := json.Marshal(res);
	
	w.Header().Set("Content-Type", "application/json")

	w.Write(jsonBytes)

}