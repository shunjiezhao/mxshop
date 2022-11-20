package handler

import (
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"gorm.io/gorm"
	proto "server/pk_service/proto/gen/v1/pk"
	"server/shared/queue"
)

type Random interface {
	RandomUserID() queue.UserId
}
type Choose interface {
	ChooseOne([]queue.UserId) queue.UserId
}

type PKService struct {
	db      *gorm.DB
	logger  *zap.Logger
	watcher *RedisWatcher
	Random  Random
	Choose  Choose
	proto.UnimplementedPKServer
}

type Config struct {
	DB     *gorm.DB
	Logger *zap.Logger
	Random Random
	Choose Choose
	queue.UserPublisher
}

func NewService(config *Config) *PKService {
	s := &PKService{
		db:     config.DB,
		logger: config.Logger,
		Random: config.Random,
		Choose: config.Choose,
	}
	s.watcher = &RedisWatcher{
		Add:           make(chan queue.UserId, 0),
		Ctx:           context.Background(),
		UserPublisher: config.UserPublisher,
	}
	go s.watcher.Watch()
	return s
}
