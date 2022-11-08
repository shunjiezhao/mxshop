package discovery

import (
	"context"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"sync"
	"time"
)

type ServiceDiscovery struct {
	cli *clientv3.Client
	//TODO: 是否需要将 prefix 相同的归为一类
	serviceList map[string]string
	lock        sync.Mutex
}

//NewServiceDiscovery  新建发现服务

func NewServiceDiscovery(endpoints []string) *ServiceDiscovery {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		zap.L().Fatal("can not connect the etcd server", zap.Error(err))
	}
	zap.L().Info("【Etcd】connect success", zap.Any("addr", endpoints))
	return &ServiceDiscovery{
		cli:         cli,
		serviceList: make(map[string]string),
	}
}

//WatchService 初始化服务列表和监视
func (s *ServiceDiscovery) WatchService(prefix string) error {
	resp, err := s.cli.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		zap.L().Info("【Etcd】can not get server list from etcd", zap.Error(err))
		return err
	}
	for _, kv := range resp.Kvs {
		s.SetServiceList(string(kv.Key), string(kv.Value))
	}
	// 监听这个前缀（健康检查）
	go s.watch(prefix)
	return nil
}

func (s *ServiceDiscovery) SetServiceList(key string, val string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.serviceList[key] = val
	zap.L().Info("【Etcd】:put key", zap.String("key", key), zap.String("val", val))
}

func (s *ServiceDiscovery) watch(prefix string) {
	rch := s.cli.Watch(context.Background(), prefix, clientv3.WithPrefix())
	zap.L().Info("【Etcd】:watch key", zap.String("prefix", prefix))
	// 发生了变化
	for change := range rch {
		for _, ev := range change.Events {
			switch ev.Type {
			case mvccpb.PUT:
				s.SetServiceList(string(ev.Kv.Key), string(ev.Kv.Value))
			case mvccpb.DELETE:
				s.DelServiceList(string(ev.Kv.Key))
			}
		}
	}
}

//GetServices 获取服务地址
func (s *ServiceDiscovery) GetServices() []string {
	s.lock.Lock()
	defer s.lock.Unlock()
	addrs := make([]string, 0)
	for _, v := range s.serviceList {
		addrs = append(addrs, v)
	}
	return addrs
}

func (s *ServiceDiscovery) DelServiceList(key string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.serviceList, key)
	zap.L().Info("【Etcd】:del key", zap.String("key", key))
}
