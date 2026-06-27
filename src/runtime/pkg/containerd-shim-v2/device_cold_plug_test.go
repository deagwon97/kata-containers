// Copyright (c) 2026 NVIDIA CORPORATION.
//
// SPDX-License-Identifier: Apache-2.0
//

package containerdshim

import (
	"testing"

	"github.com/stretchr/testify/assert"
	podresourcesv1 "k8s.io/kubelet/pkg/apis/podresources/v1"
)

func TestPodResourcesCDIDevices(t *testing.T) {
	podRes := &podresourcesv1.PodResources{
		Containers: []*podresourcesv1.ContainerResources{
			{
				Name: "legacy-device-plugin",
				Devices: []*podresourcesv1.ContainerDevices{
					{
						ResourceName: "vendor.example/gpu",
						DeviceIds:    []string{"gpu0", "gpu1"},
					},
				},
			},
			{
				Name: "dra",
				DynamicResources: []*podresourcesv1.DynamicResource{
					{
						ClaimName:      "claim",
						ClaimNamespace: "default",
						ClaimResources: []*podresourcesv1.ClaimResource{
							{
								DriverName: "vfio-gpu.kata-containers.io",
								PoolName:   "node-0",
								DeviceName: "vfio0",
								CDIDevices: []*podresourcesv1.CDIDevice{
									{Name: "nvidia.com/pgpu=vfio0"},
									{Name: ""},
								},
							},
						},
					},
				},
			},
		},
	}

	assert.Equal(t, []string{
		"vendor.example/gpu=gpu0",
		"vendor.example/gpu=gpu1",
		"nvidia.com/pgpu=vfio0",
	}, podResourcesCDIDevices(podRes))
}
