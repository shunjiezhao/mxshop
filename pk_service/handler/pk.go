package handler

import (
	"fmt"
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

type DivideQuestions interface {
	Divide(cnt int) ([]string, []string, error)
}

var QuestionCnt int = 10

type PKService struct {
	db      *gorm.DB
	logger  *zap.Logger
	watcher *RedisWatcher
	Random  Random
	Choose  Choose
	DivideQuestions
	proto.UnimplementedPKServer
}

type divideQuestion struct {
}

//TODO:完成题库的题目提取
func (d *divideQuestion) Divide(cnt int) ([]string, []string, error) {
	var ques, ans []string
	for i := 0; i < 10; i++ {
		ques = append(ques, fmt.Sprintf("题目%d", i))
		ans = append(ans, string(i%4+'a')) // a b c d
	}
	return ques, ans, nil
}

type Config struct {
	DB     *gorm.DB
	Logger *zap.Logger
	Random Random
	Choose Choose
	DivideQuestions
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
	if config.DivideQuestions == nil {
		s.DivideQuestions = &divideQuestion{}
	} else {
		s.DivideQuestions = config.DivideQuestions
	}

	go s.watcher.Watch()
	return s
}
