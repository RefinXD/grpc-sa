package main

import (
	"client/places"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)



func main(){

	creds := insecure.NewCredentials()
	cc,err := grpc.Dial("localhost:9000",grpc.WithTransportCredentials(creds))
	if err != nil{
		log.Fatal(err)
	}

	placesClient := places.NewPlacesClient(cc)
	placesService := places.NewPlaceService(placesClient)

	// placesService.UploadPlaceInfo(places.Place{
	// 	Name: "toilet",
	// 	Capacity: 1,
	// 	Facilities: []string{
	// 		"toilet",
	// 	},
	// })

	// placesService.UploadPlaceInfo(places.Place{
	// 	Name: "wifi",
	// 	Capacity: 1,
	// 	Facilities: []string{
	// 		"wifi",
	// 	},
	// })

	// placesService.UploadPlaceInfo(places.Place{
	// 	Name: "both",
	// 	Capacity: 1,
	// 	Facilities: []string{
	// 		"toilet",
	// 		"wifi",
	// 	},
	// })


	// placesService.UpdatePlace(places.UpdatePlace{
	// 	Target: &places.Place{
	// 	Name: "test",
	// 	Capacity: 1,
	// 	Facilities: []string{
	// 		"yes",
	// 	},
	// },
	// 	NewInfo : &places.Place{
	// 	Name: "test2",
	// 	Capacity: 1,
	// 	Facilities: []string{
	// 		"yes",
	// 	},
	// },
	// })


	//placesService.SearchPlaces(places.PlaceName{Name: "ac"})
	placesService.FilterPlaces(places.Filter{Facilities:[]string{"toilet"}})
}