package docker

import (
	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/go-connections/nat"
)

func containerJSON(ops ...func(*dockertypes.ContainerJSON)) dockertypes.ContainerJSON {
	c := &dockertypes.ContainerJSON{
		ContainerJSONBase: &dockertypes.ContainerJSONBase{
			Name:       "fake",
			HostConfig: &container.HostConfig{},
			State:      &dockertypes.ContainerState{},
		},
		Config: &container.Config{},
		NetworkSettings: &dockertypes.NetworkSettings{
			NetworkSettingsBase: dockertypes.NetworkSettingsBase{},
		},
	}

	for _, op := range ops {
		op(c)
	}

	return *c
}

func status(status string) func(*dockertypes.ContainerJSON) {
	return func(c *dockertypes.ContainerJSON) {
		c.ContainerJSONBase.State.Status = status
	}
}

func name(name string) func(*dockertypes.ContainerJSON) {
	return func(c *dockertypes.ContainerJSON) {
		c.ContainerJSONBase.Name = name
	}
}

func networkMode(mode string) func(*dockertypes.ContainerJSON) {
	return func(c *dockertypes.ContainerJSON) {
		c.ContainerJSONBase.HostConfig.NetworkMode = container.NetworkMode(mode)
	}
}

func nodeIP(ip string) func(*dockertypes.ContainerJSON) {
	return func(c *dockertypes.ContainerJSON) {
		c.ContainerJSONBase.Node = &dockertypes.ContainerNode{
			IPAddress: ip,
		}
	}
}

func ports(portMap nat.PortMap) func(*dockertypes.ContainerJSON) {
	return func(c *dockertypes.ContainerJSON) {
		c.NetworkSettings.NetworkSettingsBase.Ports = portMap
	}
}

func withNetwork(name string, ops ...func(*network.EndpointSettings)) func(*dockertypes.ContainerJSON) {
	return func(c *dockertypes.ContainerJSON) {
		if c.NetworkSettings.Networks == nil {
			c.NetworkSettings.Networks = map[string]*network.EndpointSettings{}
		}
		c.NetworkSettings.Networks[name] = &network.EndpointSettings{}
		for _, op := range ops {
			op(c.NetworkSettings.Networks[name])
		}
	}
}

func ipv4(ip string) func(*network.EndpointSettings) {
	return func(s *network.EndpointSettings) {
		s.IPAddress = ip
	}
}

func ipv6(ip string) func(*network.EndpointSettings) {
	return func(s *network.EndpointSettings) {
		s.GlobalIPv6Address = ip
	}
}

func swarmTask(id string, ops ...func(*swarm.Task)) swarm.Task {
	task := &swarm.Task{
		ID: id,
	}

	for _, op := range ops {
		op(task)
	}

	return *task
}

func taskSlot(slot int) func(*swarm.Task) {
	return func(task *swarm.Task) {
		task.Slot = slot
	}
}

func taskNetworkAttachment(id, name, driver string, addresses []string) func(*swarm.Task) {
	return func(task *swarm.Task) {
		task.NetworksAttachments = append(task.NetworksAttachments, swarm.NetworkAttachment{
			Network: swarm.Network{
				ID: id,
				Spec: swarm.NetworkSpec{
					Annotations: swarm.Annotations{
						Name: name,
					},
					DriverConfiguration: &swarm.Driver{
						Name: driver,
					},
				},
			},
			Addresses: addresses,
		})
	}
}

func taskStatus(ops ...func(*swarm.TaskStatus)) func(*swarm.Task) {
	return func(task *swarm.Task) {
		status := &swarm.TaskStatus{}

		for _, op := range ops {
			op(status)
		}

		task.Status = *status
	}
}

func taskState(state swarm.TaskState) func(*swarm.TaskStatus) {
	return func(status *swarm.TaskStatus) {
		status.State = state
	}
}

func taskContainerStatus(id string) func(*swarm.TaskStatus) {
	return func(status *swarm.TaskStatus) {
		status.ContainerStatus = &swarm.ContainerStatus{
			ContainerID: id,
		}
	}
}

func swarmService(ops ...func(*swarm.Service)) swarm.Service {
	service := &swarm.Service{
		ID: "serviceID",
		ServiceStatus: &swarm.ServiceStatus{
			RunningTasks: 1,
			DesiredTasks: 1,
		},
		Spec: swarm.ServiceSpec{
			Annotations: swarm.Annotations{
				Name: "defaultServiceName",
			},
		},
	}

	for _, op := range ops {
		op(service)
	}

	return *service
}

func serviceStatus(status *swarm.ServiceStatus) func(service *swarm.Service) {
	return func(service *swarm.Service) {
		service.ServiceStatus = status
	}
}

func serviceName(name string) func(service *swarm.Service) {
	return func(service *swarm.Service) {
		service.Spec.Annotations.Name = name
	}
}

func serviceLabels(labels map[string]string) func(service *swarm.Service) {
	return func(service *swarm.Service) {
		service.Spec.Annotations.Labels = labels
	}
}

func withEndpoint(ops ...func(*swarm.Endpoint)) func(*swarm.Service) {
	return func(service *swarm.Service) {
		endpoint := &swarm.Endpoint{}

		for _, op := range ops {
			op(endpoint)
		}

		service.Endpoint = *endpoint
	}
}

func virtualIP(networkID, addr string) func(*swarm.Endpoint) {
	return func(endpoint *swarm.Endpoint) {
		if endpoint.VirtualIPs == nil {
			endpoint.VirtualIPs = []swarm.EndpointVirtualIP{}
		}
		endpoint.VirtualIPs = append(endpoint.VirtualIPs, swarm.EndpointVirtualIP{
			NetworkID: networkID,
			Addr:      addr,
		})
	}
}

func withEndpointSpec(ops ...func(*swarm.EndpointSpec)) func(*swarm.Service) {
	return func(service *swarm.Service) {
		endpointSpec := &swarm.EndpointSpec{}

		for _, op := range ops {
			op(endpointSpec)
		}

		service.Spec.EndpointSpec = endpointSpec
	}
}

func modeDNSRR(spec *swarm.EndpointSpec) {
	spec.Mode = swarm.ResolutionModeDNSRR
}

func modeVIP(spec *swarm.EndpointSpec) {
	spec.Mode = swarm.ResolutionModeVIP
}
