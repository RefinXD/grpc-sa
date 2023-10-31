package places

import (
	"context"
	"fmt"
)


type PlacesService interface {
	UploadPlaceInfo(place Place) error
	UpdatePlace(updatePlace UpdatePlace) error
	SearchPlaces(name PlaceName) error 
	FilterPlaces(filter Filter) error
	RemovePlaces(name PlaceName) error
}


type placesService struct {
	placesClient PlacesClient
}

func NewPlaceService(placesClient PlacesClient) PlacesService {

	return placesService{placesClient}
}

func (base placesService) UploadPlaceInfo(place Place) error {

	res, err := base.placesClient.UploadPlaceInfo(context.Background(),&place)
	if err != nil{
		return err;
	}
	fmt.Println(res)
	fmt.Println("Service:uploadPLaceINfo")
	return nil
}


func (base placesService) UpdatePlace(updatePlace UpdatePlace) error {

	res, err := base.placesClient.UpdatePlace(context.Background(),&updatePlace)
	if err != nil{
		return err;
	}
	fmt.Println(res)
	fmt.Println("Service:update")
	return nil
}


func (base placesService) SearchPlaces(name PlaceName) error {

	res, err := base.placesClient.SearchPlaces(context.Background(),&name)
	if err != nil{
		return err;
	}
	fmt.Println(res)
	fmt.Println("Service:SearchPlaces")
	return nil
}



func (base placesService) FilterPlaces(filter Filter) error {

	res, err := base.placesClient.FilterPlaces(context.Background(),&filter)
	if err != nil{
		return err;
	}
	fmt.Println(res)
	fmt.Println("Service:FilterPlaces")
	return nil
}


func (base placesService) RemovePlaces(name PlaceName) error {

	res, err := base.placesClient.RemovePlaces(context.Background(),&name)
	if err != nil{
		return err;
	}
	fmt.Println(res)
	fmt.Println("Service:RemovePlaces")
	return nil
}
