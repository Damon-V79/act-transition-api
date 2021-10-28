package consumer

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/Damon-V79/act-transition-api/internal/mocks"
	"github.com/Damon-V79/act-transition-api/internal/model"
	"github.com/golang/mock/gomock"
)

func TestStart(t *testing.T) {

	testEvents := []model.TransitionEvent{
		{
			ID:     1,
			Type:   model.Created,
			Status: model.Processed,
			Entity: &model.Transition{
				ID:   911,
				Name: "UniqueName",
				From: "Moscow",
				To:   "Sochi",
			},
		},
		{
			ID:     11,
			Type:   model.Created,
			Status: model.Processed,
			Entity: &model.Transition{
				ID:   1024,
				Name: "OrdinaryName",
				From: "Moscow",
				To:   "Novosibirsk",
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	mockRepo := mocks.NewMockEventRepo(mockCtrl)

	ctx := context.Background()

	channelSize := 512
	consumerCount := uint64(1)
	consumeSize := uint64(5)
	consumeTimeout := 10 * time.Millisecond
	events := make(chan []model.TransitionEvent, channelSize)

	mockRepo.EXPECT().Lock(gomock.Any(), gomock.Any()).AnyTimes().Return(testEvents, nil)

	consumer := NewDbConsumer(
		consumerCount,
		consumeSize,
		consumeTimeout,
		mockRepo,
		events,
	)

	consumer.Start(ctx)

	for recvEvents := range events {
		if !reflect.DeepEqual(recvEvents, testEvents) {
			t.Fatal("The test result did not match with expected result")
		}
		break
	}

	consumer.Close()

}
