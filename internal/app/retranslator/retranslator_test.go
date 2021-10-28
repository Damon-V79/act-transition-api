package retranslator

import (
	"context"
	"testing"
	"time"

	"github.com/Damon-V79/act-transition-api/internal/mocks"
	"github.com/golang/mock/gomock"
)

func TestStart(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	mockRepo := mocks.NewMockEventRepo(mockCtrl)
	mockSender := mocks.NewMockEventSender(mockCtrl)

	mockRepo.EXPECT().Lock(gomock.Any(), gomock.Any()).AnyTimes()
	mockRepo.EXPECT().Unlock(gomock.Any(), gomock.Any()).AnyTimes()
	mockRepo.EXPECT().Remove(gomock.Any(), gomock.Any()).AnyTimes()

	mockSender.EXPECT().Send(gomock.Any(), gomock.Any()).AnyTimes()

	cfg := Config{
		ChannelSize:    512,
		ConsumerCount:  2,
		ConsumeSize:    10,
		ConsumeTimeout: 10 * time.Second,
		ProducerCount:  2,
		WorkerCount:    2,
		Repo:           mockRepo,
		Sender:         mockSender,
	}

	ctx := context.Background()
	retranslator := NewRetranslator(cfg)
	retranslator.Start(ctx)
	retranslator.Close()
}
