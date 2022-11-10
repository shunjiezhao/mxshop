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

//ListenLeaseRespChan 监听 续租情况
//TODO:优雅的退出
func (s *ServiceRegister) ListenLeaseRespChan() {
	i := 0
	for {
		select {
		case leaseKeepResp := <-s.keepAliveChan:
			i++
			if i == 100 {
				s.Logger.Info("【Etcd】 续约成功", zap.Any("resp", leaseKeepResp))
				i = 0
			}
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
func NewServiceRegister(endpoints []string, serName, addr string, lease int64, logger *zap.Logger) (*ServiceRegister, error) {
	cli, err := c3.New(c3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}

	ser := &ServiceRegister{
		cli:    cli,
		key:    "/" + schema + "/" + serName + "/" + addr, // 这里需要与前端进行协商
		val:    addr,
		Logger: logger,
	}

	//申请租约设置时间keepalive
	if err := ser.putKeyWithLease(lease); err != nil {
		return nil, err
	}
	return ser, nil
}
