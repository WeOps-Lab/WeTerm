package healthcheck

import (
	"fmt"
	"sort"
	"strings"

	capi "github.com/hashicorp/consul/api"
)

type ConsulHealth struct {
	BaseHealthChecker
	c *capi.Client
}

func NewConsulHealth() ConsulHealth {
	config := capi.DefaultConfig()
	client, err := capi.NewClient(config)
	if err != nil {
		panic(err)
	}
	return ConsulHealth{
		c: client,
	}
}

func (h ConsulHealth) Check() []HealthResult {
	result := []HealthResult{}
	// 获取consul的选举情况
	result = append(
		result,
		HealthResult{
			status:  Common,
			message: "Consul 选举状态:",
		})
	leader, err := h.c.Status().Leader()
	if err != nil {
		result = append(result, HealthResult{status: Error, message: err.Error()})
	} else {
		result = append(result, HealthResult{status: Healthy, message: fmt.Sprintf("[aqua]Leader: %s", leader)})
	}
	peers, err := h.c.Status().Peers()
	if err != nil {
		result = append(result, HealthResult{status: Error, message: err.Error()})
	} else {
		result = append(result, HealthResult{status: Healthy, message: fmt.Sprintf("Peers: %s", strings.Join(peers, ","))})
	}
	// 获取consul集群状态
	members, err := h.c.Agent().Members(false)
	//    Agent状态描述
	// 	  AgentMemberNone    = 0
	//	  AgentMemberAlive   = 1
	//	  AgentMemberLeaving = 2
	//	  AgentMemberLeft    = 3
	//	  AgentMemberFailed  = 4
	//    这里简单处理，不正常的都返回Error
	result = append(
		result,
		HealthResult{
			status:  Common,
			message: "Consul 集群状态:",
		})
	if err != nil {
		result = append(result, HealthResult{status: Error, message: err.Error()})
	}
	for i := range members {
		m := members[i]
		if m.Status == 1 {
			result = append(result, HealthResult{status: Healthy, message: h.buildMemberMessage(m)})
		} else {
			result = append(result, HealthResult{status: Error, message: h.buildMemberMessage(m)})
		}
	}

	// 获取consul注册的服务状态
	result = append(
		result,
		HealthResult{
			status:  Common,
			message: "Consul 服务状态:",
		})
	services, err := h.c.Agent().Checks()
	if err != nil {
		result = append(result, HealthResult{status: Error, message: err.Error()})
	}
	// 排序
	servicesKeys := make([]string, 0, len(services))
	for servicesKey := range services {
		servicesKeys = append(servicesKeys, servicesKey)
	}
	sort.SliceStable(servicesKeys, func(i, j int) bool {
		return services[servicesKeys[i]].ServiceName < services[servicesKeys[j]].ServiceName
	})
	for _, k := range servicesKeys {
		s := services[k]
		if s.Status == "passing" {
			result = append(result, HealthResult{status: Healthy, message: h.buildServiceMessage(s)})
		} else {
			result = append(result, HealthResult{status: Error, message: h.buildServiceMessage(s)})
		}
	}
	return result
}

func (h ConsulHealth) buildMemberMessage(member *capi.AgentMember) string {
	return fmt.Sprintf("[aqua]ID: [white]%s [aqua]Name: [white]%s [aqua]Addr: [white]%s, [aqua]Build: [white]%s, [aqua]Role: [white]%s", member.Tags["id"], member.Name, member.Addr, member.Tags["build"], member.Tags["role"])
}

func (h ConsulHealth) buildServiceMessage(service *capi.AgentCheck) string {
	return fmt.Sprintf("Service: [yellow]%s [aqua]ID: [yellow]%s \n              [aqua]Output: [white]%s", service.ServiceName, service.ServiceID, service.Output)
}
