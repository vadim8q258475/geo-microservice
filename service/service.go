package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	pb "github.com/vadim8q258475/geo-microservice/pb"

	"github.com/ekomobile/dadata/v2/api/suggest"
	"github.com/ekomobile/dadata/v2/client"
	"go.uber.org/zap"
)

type GrpcService struct {
	pb.UnimplementedGeoServiceServer
	api       *suggest.Api
	apiKey    string
	secretKey string
	logger    *zap.Logger
}

func NewGeoGrpcService(apiKey, secretKey string, logger *zap.Logger) *GrpcService {
	var err error
	endpointUrl, err := url.Parse("https://suggestions.dadata.ru/suggestions/api/4_1/rs/")
	if err != nil {
		return nil
	}

	creds := client.Credentials{
		ApiKeyValue:    apiKey,
		SecretKeyValue: secretKey,
	}

	api := suggest.Api{
		Client: client.NewClient(endpointUrl, client.WithCredentialProvider(&creds)),
	}

	return &GrpcService{
		api:       &api,
		apiKey:    apiKey,
		secretKey: secretKey,
		logger:    logger,
	}
}

func (s *GrpcService) AddressSearch(ctx context.Context, in *pb.AddressSearchRequest) (*pb.AddressSearchResponse, error) {
	res := &pb.AddressSearchResponse{Addresses: make([]*pb.Address, 0)}
	rawRes, err := s.api.Address(context.Background(), &suggest.RequestParams{Query: in.Query})
	if err != nil {
		s.logger.Error(fmt.Sprintf("api request error: %w", err))
		return nil, err
	}

	for _, r := range rawRes {
		if r.Data.City == "" || r.Data.Street == "" {
			continue
		}
		res.Addresses = append(res.Addresses, &pb.Address{City: r.Data.City, Street: r.Data.Street, House: r.Data.House, Lat: r.Data.GeoLat, Lon: r.Data.GeoLon})
	}
	s.logger.Info("success address search")
	return res, nil
}

func (s *GrpcService) GeoCode(ctx context.Context, in *pb.GeoCodeRequest) (*pb.AddressSearchResponse, error) {
	httpClient := &http.Client{}
	var data = strings.NewReader(fmt.Sprintf(`{"lat": %s, "lon": %s}`, in.Lat, in.Lng))
	req, err := http.NewRequest("POST", "https://suggestions.dadata.ru/suggestions/api/4_1/rs/geolocate/address", data)
	if err != nil {
		s.logger.Error(fmt.Sprintf("create request error: %w", err))
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", s.apiKey))
	resp, err := httpClient.Do(req)
	if err != nil {
		s.logger.Error(fmt.Sprintf("api request error: %w", err))
		return nil, err
	}
	var geoCode GeoCode

	err = json.NewDecoder(resp.Body).Decode(&geoCode)
	if err != nil {
		s.logger.Error(fmt.Sprintf("decode error: %w", err))
		return nil, err
	}
	result := &pb.AddressSearchResponse{Addresses: make([]*pb.Address, 0)}

	for _, r := range geoCode.Suggestions {
		var address pb.Address
		address.City = string(r.Data.City)
		address.Street = string(r.Data.Street)
		address.House = r.Data.House
		address.Lat = r.Data.GeoLat
		address.Lon = r.Data.GeoLon
		result.Addresses = append(result.Addresses, &address)
	}
	s.logger.Info("success geocode search")
	return result, nil
}
