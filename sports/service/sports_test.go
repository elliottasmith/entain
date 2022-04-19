package service

import (
	"context"
	"testing"
	"time"

	"github.com/elliottasmith/entain/sports/db/mocks"
	"github.com/elliottasmith/entain/sports/proto/sports"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
	"syreclabs.com/go/faker"
)

func TestListRaces(t *testing.T) {
	// Create a new mock EventsRepo
	mockEventsRepo := new(mocks.EventsRepo)

	// Create a new sports service using the mock EventsRepo
	es := NewSportsService(mockEventsRepo)

	// Create a new context with the option to cancel once processing is complete and a default timeout of 1 second.
	// As this is an inexpensive service 1 second should be suitable.
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)

	// Add the cancel function to the defer call.
	defer cancel()

	// Create a new list events request
	req := &sports.ListEventsRequest{
		Filter: &sports.ListEventsRequestFilter{
		},
		Order: &sports.ListEventsRequestOrder{
		},
	}
	
	// Create a mock event to be returned from the list function
	mockEvent := []*sports.Event{
		{
			Id: faker.RandomInt64(1, 10),
			SportType: faker.Team().Creature(),
			League: faker.Company().Name(),
			Country: faker.Team().State(),
			LocationId: faker.RandomInt64(1, 10),
			Name: faker.Team().Name(),
			Round: faker.RandomInt64(1, 24),
			Game: faker.RandomInt64(1, 8),
			Visible: faker.RandomInt64(0, 1) == 0,
			AdvertisedStartTime: timestamppb.New(faker.Time().Between(time.Now().AddDate(0, 0, -1), time.Now().AddDate(0, 0, 2))),
		},
	}

	// Create a mock result for the list function to return 1 event.
	mockEventsRepo.On("List", req.Filter, req.Order).Return(mockEvent, nil).Times(1)

	// Call the list events service with the context and request
	res, err := es.ListEvents(ctx, req)

	// Assert there are no errors and that the event response matches the expected event
	assert.NoError(t, err)
	assert.Equal(t, &sports.ListEventsResponse{
		Events: mockEvent,
	}, res)
}