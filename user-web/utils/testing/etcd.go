package testing

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"testing"
)

//TODO:等v5.0镜像出来后 配置测试环境

const (
	image         = "bitnami/etcd:latest"
	containerPort = "2379/tcp"
)

var Dsn string

const defaultEtcdURI = "localhost:2379"

func RunWithMongoInDocker(m *testing.M) int {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	cli.ContainerRemove(ctx, "Test_Etcd", types.ContainerRemoveOptions{Force: true})
	//, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, platform *specs.Platform, containerName string) (container.ContainerCreateCreatedBody, error) {
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: image,
		Tty:   false,
		ExposedPorts: map[nat.Port]struct{}{
			containerPort: {},
		},
		Env: []string{""},
	}, &container.HostConfig{
		PortBindings: nat.PortMap{
			containerPort: []nat.PortBinding{
				// port = 0 会分配空闲端口
				{HostIP: "127.0.0.1", HostPort: "0"},
			},
		},
	}, nil, nil, "Test_Postgres")
	if err != nil {
		log.Fatalln("check docker server is start?")
	}
	containerID := resp.ID

	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("panic: %s", err)
		}
		fmt.Println("remove container")
		cli.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{Force: true})
	}()

	err = cli.ContainerStart(ctx, containerID, types.ContainerStartOptions{})
	// 获取容器信息
	inspect, err := cli.ContainerInspect(ctx, containerID)
	if err != nil {
		panic(err)
	}
	// 获取分配端口号
	Portbindings := inspect.NetworkSettings.Ports[containerPort][0]
	fmt.Printf("listening address :%v\n", Portbindings)
	fmt.Printf("Postgres Port : %s", Portbindings.HostPort)
	// postgres 连接url
	Dsn = fmt.Sprintf("host=localhost user=postgres dbname=postgres password=az123. port=%s sslmode=disable TimeZone=Asia/Shanghai", Portbindings.HostPort)
	return m.Run()
}

// 在这 这前调用 RunWithMongoInDocker 来启动 etcd 服务
func NewClient(c context.Context) (*clientv3.Client, error) {
	return nil, nil
}
func NewDefaultClient(c context.Context) {
}
