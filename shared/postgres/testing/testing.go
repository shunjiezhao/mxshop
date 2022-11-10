package testing

import (
	"context"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

const (
	image         = "postgres:latest"
	containerPort = "5432/tcp"
)

var Dsn string

const defaultMongoURI = "host=localhost user=postgres dbname=mxshop password=az123. port=5432 sslmode=disable TimeZone=Asia/Shanghai"

func RunWithMongoInDocker(m *testing.M) int {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	cli.ContainerRemove(ctx, "Test_Postgres", types.ContainerRemoveOptions{Force: true})
	//, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, platform *specs.Platform, containerName string) (container.ContainerCreateCreatedBody, error) {
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: image,
		Tty:   false,
		ExposedPorts: map[nat.Port]struct{}{
			containerPort: {},
		},
		Env: []string{"POSTGRES_PASSWORD=az123."},
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

// 在这 这前调用 RunWithMongoInDocker 来启动 mongo 服务
func NewClient(c context.Context) (*gorm.DB, error) {
	if Dsn == "" {
		return nil, fmt.Errorf("postgres uri not set please run RunWithMongoInDocker function")
	}

	// 等待一下
	time.Sleep(2 * time.Second)
	DB, err := gorm.Open(postgres.Open(Dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	sqlDB, err := DB.DB()

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(10)
	return DB, nil
}
func NewDefaultClient(c context.Context) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(Dsn), &gorm.Config{})
}
