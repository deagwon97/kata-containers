// Copyright (c) 2017-2018 Intel Corporation
// Copyright (c) 2018 Huawei Corporation
//
// SPDX-License-Identifier: Apache-2.0
//

package drivers

import (
	"testing"

	"github.com/kata-containers/kata-containers/src/runtime/pkg/device/config"
	"github.com/stretchr/testify/assert"
)

func TestGetVFIODetails(t *testing.T) {
	type testData struct {
		deviceStr   string
		expectedStr string
	}

	data := []testData{
		{"0000:02:10.0", "0000:02:10.0"},
		{"0000:0210.0", ""},
		{"f79944e4-5a3d-11e8-99ce-", ""},
		{"f79944e4-5a3d-11e8-99ce", ""},
		{"test", ""},
		{"", ""},
	}

	for _, d := range data {
		deviceBDF, deviceSysfsDev, vfioDeviceType, err := GetVFIODetails(d.deviceStr, "")

		switch vfioDeviceType {
		case config.VFIOPCIDeviceNormalType:
			assert.Equal(t, d.expectedStr, deviceBDF)
		case config.VFIOPCIDeviceMediatedType, config.VFIOAPDeviceMediatedType:
			assert.Equal(t, d.expectedStr, deviceSysfsDev)
		default:
			assert.NotNil(t, err)
		}

		if d.expectedStr == "" {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
		}
	}

}

func TestVFIOGroupNameFromDevPath(t *testing.T) {
	testCases := []struct {
		name        string
		path        string
		expected    string
		expectError bool
	}{
		{
			name:     "regular iommu group",
			path:     "/dev/vfio/17",
			expected: "17",
		},
		{
			name:     "unsafe no-iommu group",
			path:     "/dev/vfio/noiommu-17",
			expected: "17",
		},
		{
			name:        "vfio control device",
			path:        "/dev/vfio/vfio",
			expectError: true,
		},
		{
			name:        "iommufd cdev",
			path:        "/dev/vfio/devices/vfio17",
			expectError: true,
		},
		{
			name:        "non-numeric no-iommu group",
			path:        "/dev/vfio/noiommu-gpu",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			group, err := VFIOGroupNameFromDevPath(tc.path)
			if tc.expectError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, group)
		})
	}
}
