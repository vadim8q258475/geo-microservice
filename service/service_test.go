package service

import (
	"context"
	"testing"

	"github.com/ekomobile/dadata/v2/api/suggest"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/vadim8q258475/geo-microservice/mock"
	pb "github.com/vadim8q258475/geo-microservice/pb"
	"go.uber.org/zap"
)

func TestGrpcService_AddressSearch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockApi := mock.NewMockApi(ctrl)
	logger := zap.NewNop()

	service := GrpcService{api: mockApi, logger: logger}
	// api err
	mockApi.EXPECT().Address(gomock.Any(), gomock.Any()).Return(nil, assert.AnError)
	req := pb.AddressSearchRequest{Query: "query"}
	_, err := service.AddressSearch(context.Background(), &req)
	assert.Error(t, err)

	rawRes := []*suggest.AddressSuggestion{}
	mockApi.EXPECT().Address(gomock.Any(), gomock.Any()).Return(rawRes, nil)
	resp, err := service.AddressSearch(context.Background(), &req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestNewGrpcService(t *testing.T) {
	apiKey := "apiKey"
	secretKey := "secretKey"
	logger := zap.NewNop()

	service := NewGeoGrpcService(apiKey, secretKey, logger)

	assert.NotNil(t, service.api)
	assert.Equal(t, apiKey, "apiKey")
	assert.Equal(t, secretKey, "secretKey")
}
