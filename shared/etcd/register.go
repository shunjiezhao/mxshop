package etcd

import (
	"context"
	c3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"log"
	"time"
)

// 服务注册及健康检查

const schema = "zsj"

type RegisterClient interface {
	Register(serName, addr string, lease int64) error // 注册
	Watch()                                           // 监听
	Close() error                                     // 关闭
}

// 创建注册服务
type ServiceRegister struct {
	cli     *c3.Client
	leaseID c3.LeaseID //租约id
	// 租约
	keepAliveChan <-chan *c3.LeaseKeepAliveResponse
	key           string
	val           string
	Logger        *zap.Logger
}

func (s *ServiceRegister) Register(serName, addr string, lease int64) error {
	s.key = "/" + schema + "/" + serName + "/" + addr // 这里需要与前端进行协商
	s.val = addr
	return s.putKeyWithLease(lease)
}

//设置租约+注册
func (s *ServiceRegister) putKeyWithLease(lease int64) error {
	//设置租约
	resp, err := s.cli.Grant(context.Background(), lease)
	if err != nil {
		return err
	}
	//注册服务并绑定租约
	_, err = s.cli.Put(context.Background(), s.key, s.val, c3.WithLease(resp.ID))
	if err != nil {

		return err
	}
	//设置续租 定期发送需求请求
	leaseRespChan, err := s.cli.KeepAlive(context.Background(), resp.ID)
	if err != nil {
		return err
	}
	s.leaseID = resp.ID
	s.keepAliveChan = leaseRespChan
	s.Logger.Info("【Etcd】: put key", zap.Int("leaseID", int(s.leaseID)), zap.String("key", s.key),
		zap.String("val", s.val))
	return nil

}

//Watch 监听 续租情况
//TODO:优雅的退出
func (s *ServiceRegister) Watch() {
	for {
		select {
		case <-s.keepAliveChan:
		}
	}
	s.Logger.Info("【Etcd】 关闭续租")
}

// Close 注销服务
func (s *ServiceRegister) Close() error {
	//撤销租约
	if _, err := s.cli.Revoke(context.Background(), s.leaseID); err != nil {
		return err
	}
	s.Logger.Info("【Etcd】 撤销租约")
	return s.cli.Close()
}

//NewServiceRegister 新建注册服务
func NewServiceRegister(endpoints []string, logger *zap.Logger) (RegisterClient, error) {
	cli, err := c3.New(c3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}

	ser := &ServiceRegister{
		cli:    cli,
		Logger: logger,
	}
	return ser, nil
}
