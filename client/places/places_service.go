package places

import (
	"context"
	"fmt"
)


type PlacesService interface {
	UploadPlaceInfo(place Place) (*Place, error)
	UpdatePlace(updatePlace UpdatePlace) (*Place, error)
	SearchPlaces(name PlaceName) (*PlaceList,error) 
	FilterPlaces(filter Filter) (*PlaceList,error)
	RemovePlaces(name PlaceName) (*Empty,error)
	GetPlaceInfo(placeId PlaceId) (*Place, error)
}


type placesService struct {
	placesClient PlaceServiceClient
}

func NewPlaceService(placesClient PlaceServiceClient) PlacesService {

	return placesService{placesClient}
}

func (base placesService) UploadPlaceInfo(place Place) (*Place, error) {

	res, err := base.placesClient.UploadPlaceInfo(context.Background(),&place)
	if err != nil{
		return nil,err;
	}
	fmt.Println(res)
	fmt.Println("Service:uploadPLaceINfo")
	return res,nil
}


func (base placesService) UpdatePlace(updatePlace UpdatePlace) (*Place,error) {

	res, err := base.placesClient.UpdatePlace(context.Background(),&updatePlace)
	if err != nil{
		fmt.Println(err)
		return nil,err;
	}
	fmt.Println(res)
	fmt.Println("Service:update")
	return res,nil
}


func (base placesService) SearchPlaces(name PlaceName) (*PlaceList,error) {

	res, err := base.placesClient.SearchPlaces(context.Background(),&name)
	if err != nil{
		return nil,err;
	}
	fmt.Println(res)
	fmt.Println("Service:SearchPlaces")
	return res,nil
}



func (base placesService) FilterPlaces(filter Filter) (*PlaceList,error) {
	res, err := base.placesClient.FilterPlaces(context.Background(),&filter)
	if err != nil{
		return nil,err;
	}
	fmt.Println(res)
	fmt.Println("Service:FilterPlaces")
	return res,nil
}


func (base placesService) RemovePlaces(name PlaceName) (*Empty,error) {

	res, err := base.placesClient.RemovePlaces(context.Background(),&name)
	if err != nil{
		return nil,err;
	}
	fmt.Println(res)
	fmt.Println("Service:RemovePlaces")
	return nil,nil
}

func (base placesService) GetPlaceInfo(id PlaceId) (*Place,error) {

	res, err := base.placesClient.GetPlaceInfo(context.Background(),&id)
	if err != nil{
		return nil,err;
	}
	fmt.Println(res)
	fmt.Println("Service:GetPlaceInfo")
	return res,nil
}
