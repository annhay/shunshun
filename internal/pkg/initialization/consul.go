package initialization

import (
	"github.com/google/uuid"
	capi "github.com/hashicorp/consul/api"
)

// ConsulClient consul客户端实例
type ConsulClient struct {
	Client *capi.Client
}

func NewConsul(addr string) *ConsulClient {
	cfg := capi.DefaultConfig()
	cfg.Address = addr //consul 连接地址匹配
	client, err := capi.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	return &ConsulClient{Client: client}
}

// ConsulKV 服务注册信息配置
type ConsulKV struct {
	Name    string   //服务名
	Tags    []string //服务标签
	Address string   //服务 ip
	Port    int      //服务端口号
}

// RegisterServer consul注册方法
func (c *ConsulClient) RegisterServer(cfg ConsulKV) (string, error) {
	registration := capi.AgentServiceRegistration{
		ID:      uuid.NewString(),
		Name:    cfg.Name,
		Tags:    cfg.Tags,
		Port:    cfg.Port,
		Address: cfg.Address,
	}
	err := c.Client.Agent().ServiceRegister(&registration)
	return registration.ID, err
}

// DeregisterServer consul注销方法
func (c *ConsulClient) DeregisterServer(serviceID string) error {
	return c.Client.Agent().ServiceDeregister(serviceID)
}
