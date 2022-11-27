package handler

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"hash/fnv"
	"server/pk_service/global"
	proto "server/pk_service/proto/gen/v1/pk"
	"server/shared/queue"
)

// 参加活动
func (P *PKService) Join(ctx context.Context, req *proto.JoinRequest) (*proto.JoinResponse, error) {
	var (
		divide []string // 题目列表
		answer []string // 答案列表
	)
	// check user id
	Uid := req.Id
	c, cancel := context.WithTimeout(context.Background(), global.MaxWaitRedisTime)
	defer cancel()
	if err := P.UserIsExists(c, Uid); err != nil {
		P.logger.Info("user is not exist", zap.Error(err))
		return nil, status.Errorf(codes.FailedPrecondition, "")
	}
	// 现在获取成功
	switch req.FindType {
	case proto.FindType_Random:
		// 随机选在线的用户
		// 推入等待列表
		P.watcher.Add <- queue.UserId(Uid)
	case proto.FindType_Choose:
		// 先检查用户是否参加活动
		err := P.UserIsExists(c, req.OtherId)
		if err != nil {
			P.logger.Info("user is not exist", zap.Error(err))
			return nil, status.Errorf(codes.FailedPrecondition, "")
		}
		// 发起挑战
		// 挑战 id

	}
	// 返回建立成功开始建立
	// 查询题目
	divide, answer, _ = P.Divide(QuestionCnt)
	resp := &proto.JoinResponse{
		Question: divide,
		Answer:   answer,
	}
	return resp, nil
}

func (P *PKService) Create(context.Context, *proto.CreateRequest) (*proto.CreateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}

//用户是否参加活动
func (P *PKService) UserIsExists(c context.Context, Uid int32) error {
	//TODO: 这里需要 活动号
	if ok, err := global.RedisClient.GetBit(c, global.RedisPartyPrefix+":1", int64(hashUid(Uid))).Result(); err != nil || ok == 0 {
		if err == redis.Nil {
			P.logger.Info("无法获取 user bitmap")
			return status.Error(codes.Internal, "无法获取")
		}
		if ok == 0 {
			return status.Error(codes.Internal, "没有该用户")
		}
	}
	return nil
}

func hashUid(Uid int32) int64 {
	hash := fnv.New32()
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, uint32(Uid))
	hash.Write(bs)
	return int64(hash.Sum32()) // uint32 -> int64 is safe
}

// 加入活动 设置某位 bit 位1
func (P *PKService) TakePartIn(c context.Context, req *proto.TakePartInRequest) (*emptypb.Empty, error) {
	partyKey := fmt.Sprintf("%s:%d", global.RedisPartyPrefix, req.Id)
	global.RedisClient.SetBit(c, partyKey, hashUid(req.Uid), 1)
	return &emptypb.Empty{}, nil
}
