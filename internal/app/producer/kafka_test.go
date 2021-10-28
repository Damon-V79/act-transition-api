package producer

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Damon-V79/act-transition-api/internal/app/repo"
	"github.com/Damon-V79/act-transition-api/internal/mocks"
	"github.com/Damon-V79/act-transition-api/internal/model"
	"github.com/gammazero/workerpool"
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
	mockSender := mocks.NewMockEventSender(mockCtrl)

	channelSize := 512
	events := make(chan []model.TransitionEvent, channelSize)

	ctx := context.Background()

	producerCount := uint64(1)
	workerCount := 2
	workerPool := workerpool.New(workerCount)
	producer := NewKafkaProducer(
		producerCount,
		mockSender,
		events,
		repo.NewDbCleaner(mockRepo),
		repo.NewDbUpdater(mockRepo),
		workerPool,
	)

	events <- testEvents

	mockSender.EXPECT().Send(gomock.Any(), &testEvents[0]).AnyTimes().Return(nil)
	mockSender.EXPECT().Send(gomock.Any(), &testEvents[1]).AnyTimes().Return(errors.New("Test error"))

	mockRepo.EXPECT().Remove(gomock.Any(), []uint64{11}).AnyTimes().Return(nil)
	mockRepo.EXPECT().Unlock(gomock.Any(), []uint64{1}).AnyTimes().Return(nil)

	producer.Start(ctx)

	time.Sleep(100 * time.Millisecond)

	producer.Close()

}
