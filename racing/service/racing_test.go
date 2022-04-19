package service

import (
	"context"
	"testing"
	"time"

	"github.com/elliottasmith/entain/racing/db/mocks"
	"github.com/elliottasmith/entain/racing/proto/racing"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
	"syreclabs.com/go/faker"
)

func TestListRaces(t *testing.T) {
	// Create a new mock RacesRepo
	mockRacingRepo := new(mocks.RacesRepo)

	// Create a new racing service using the mock RacesRepo
	rs := NewRacingService(mockRacingRepo)

	// Create a new context with the option to cancel once processing is complete and a default timeout of 1 second.
	// As this is an inexpensive service 1 second should be suitable.
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)

	// Add the cancel function to the defer call.
	defer cancel()

	// Create a new list races request, filtering visible only
	req := &racing.ListRacesRequest{
		Filter: &racing.ListRacesRequestFilter{
		},
		Order: &racing.ListRacesRequestOrder{
		},
	}
	
	// Create a mock race to be returned from the list function
	mockRace := []*racing.Race{
		{
			Id: faker.RandomInt64(1, 10),
			MeetingId: faker.RandomInt64(1, 10),
			Name: faker.Team().Name(),
			Number: faker.RandomInt64(1, 10),
			Visible: faker.RandomInt64(0, 1) == 0,
			AdvertisedStartTime: timestamppb.New(faker.Time().Between(time.Now().AddDate(0, 0, -1), time.Now().AddDate(0, 0, 2))),
		},
	}

	// Create a mock result for the list function to return 1 race.
	mockRacingRepo.On("List", req.Filter, req.Order).Return(mockRace, nil).Times(1)

	// Call the list races service with the context and request
	res, err := rs.ListRaces(ctx, req)

	// Assert there are no errors and that the race response matches the expected race
	assert.NoError(t, err)
	assert.Equal(t, &racing.ListRacesResponse{
		Races: mockRace,
	}, res)
}

func TestGetRace(t *testing.T) {
	// Create a new mock RacesRepo
	mockRacingRepo := new(mocks.RacesRepo)

	// Create a new racing service using the mock RacesRepo
	rs := NewRacingService(mockRacingRepo)

	// Create a new context with the option to cancel once processing is complete and a default timeout of 1 second.
	// As this is an inexpensive service 1 second should be suitable.
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)

	// Add the cancel function to the defer call.
	defer cancel()

	// Create a new get race request
	req := &racing.GetRaceRequest{
		Id: 4,
	}
	
	// Create a mock race to be returned from the get function
	mockRace := &racing.Race{
		Id: 4,
		MeetingId: faker.RandomInt64(1, 10),
		Name: faker.Team().Name(),
		Number: faker.RandomInt64(1, 10),
		Visible: faker.RandomInt64(0, 1) == 0,
		AdvertisedStartTime: timestamppb.New(faker.Time().Between(time.Now().AddDate(0, 0, -1), time.Now().AddDate(0, 0, 2))),
	}

	// Create a mock result for the get function to return 1 race.
	mockRacingRepo.On("Get", &req.Id).Return(mockRace, nil).Times(1)

	// Call the get race service with the context and request
	res, err := rs.GetRace(ctx, req)

	// Assert there are no errors and that the race response matches the expected race
	assert.NoError(t, err)
	assert.Equal(t, &racing.GetRaceResponse{
		Race: mockRace,
	}, res)
}