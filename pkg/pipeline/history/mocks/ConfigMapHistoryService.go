// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	chartConfig "github.com/devtron-labs/devtron/internal/sql/repository/chartConfig"
	"github.com/devtron-labs/devtron/pkg/bean"
	history "github.com/devtron-labs/devtron/pkg/pipeline/history/bean"

	mock "github.com/stretchr/testify/mock"

	pipelineConfig "github.com/devtron-labs/devtron/internal/sql/repository/pipelineConfig"

	repository "github.com/devtron-labs/devtron/pkg/pipeline/history/repository"

	time "time"
)

// ConfigMapHistoryService is an autogenerated mock type for the ConfigMapHistoryService type
type ConfigMapHistoryService struct {
	mock.Mock
}

// ConvertConfigDataToComponentLevelDto provides a mock function with given fields: config, configType, userHasAdminAccess
func (_m *ConfigMapHistoryService) ConvertConfigDataToComponentLevelDto(config *bean.ConfigData, configType repository.ConfigType, userHasAdminAccess bool) (*history.ComponentLevelHistoryDetailDto, error) {
	ret := _m.Called(config, configType, userHasAdminAccess)

	var r0 *history.ComponentLevelHistoryDetailDto
	var r1 error
	if rf, ok := ret.Get(0).(func(*bean.ConfigData, repository.ConfigType, bool) (*history.ComponentLevelHistoryDetailDto, error)); ok {
		return rf(config, configType, userHasAdminAccess)
	}
	if rf, ok := ret.Get(0).(func(*bean.ConfigData, repository.ConfigType, bool) *history.ComponentLevelHistoryDetailDto); ok {
		r0 = rf(config, configType, userHasAdminAccess)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*history.ComponentLevelHistoryDetailDto)
		}
	}

	if rf, ok := ret.Get(1).(func(*bean.ConfigData, repository.ConfigType, bool) error); ok {
		r1 = rf(config, configType, userHasAdminAccess)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateCMCSHistoryForDeploymentTrigger provides a mock function with given fields: pipeline, deployedOn, deployedBy
func (_m *ConfigMapHistoryService) CreateCMCSHistoryForDeploymentTrigger(pipeline *pipelineConfig.Pipeline, deployedOn time.Time, deployedBy int32) error {
	ret := _m.Called(pipeline, deployedOn, deployedBy)

	var r0 error
	if rf, ok := ret.Get(0).(func(*pipelineConfig.Pipeline, time.Time, int32) error); ok {
		r0 = rf(pipeline, deployedOn, deployedBy)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateHistoryFromAppLevelConfig provides a mock function with given fields: appLevelConfig, configType
func (_m *ConfigMapHistoryService) CreateHistoryFromAppLevelConfig(appLevelConfig *chartConfig.ConfigMapAppModel, configType repository.ConfigType) error {
	ret := _m.Called(appLevelConfig, configType)

	var r0 error
	if rf, ok := ret.Get(0).(func(*chartConfig.ConfigMapAppModel, repository.ConfigType) error); ok {
		r0 = rf(appLevelConfig, configType)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateHistoryFromEnvLevelConfig provides a mock function with given fields: envLevelConfig, configType
func (_m *ConfigMapHistoryService) CreateHistoryFromEnvLevelConfig(envLevelConfig *chartConfig.ConfigMapEnvModel, configType repository.ConfigType) error {
	ret := _m.Called(envLevelConfig, configType)

	var r0 error
	if rf, ok := ret.Get(0).(func(*chartConfig.ConfigMapEnvModel, repository.ConfigType) error); ok {
		r0 = rf(envLevelConfig, configType)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetDeployedHistoryByPipelineIdAndWfrId provides a mock function with given fields: pipelineId, wfrId, configType
func (_m *ConfigMapHistoryService) GetDeployedHistoryByPipelineIdAndWfrId(pipelineId int, wfrId int, configType repository.ConfigType) (*repository.ConfigmapAndSecretHistory, bool, []string, error) {
	ret := _m.Called(pipelineId, wfrId, configType)

	var r0 *repository.ConfigmapAndSecretHistory
	var r1 bool
	var r2 []string
	var r3 error
	if rf, ok := ret.Get(0).(func(int, int, repository.ConfigType) (*repository.ConfigmapAndSecretHistory, bool, []string, error)); ok {
		return rf(pipelineId, wfrId, configType)
	}
	if rf, ok := ret.Get(0).(func(int, int, repository.ConfigType) *repository.ConfigmapAndSecretHistory); ok {
		r0 = rf(pipelineId, wfrId, configType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*repository.ConfigmapAndSecretHistory)
		}
	}

	if rf, ok := ret.Get(1).(func(int, int, repository.ConfigType) bool); ok {
		r1 = rf(pipelineId, wfrId, configType)
	} else {
		r1 = ret.Get(1).(bool)
	}

	if rf, ok := ret.Get(2).(func(int, int, repository.ConfigType) []string); ok {
		r2 = rf(pipelineId, wfrId, configType)
	} else {
		if ret.Get(2) != nil {
			r2 = ret.Get(2).([]string)
		}
	}

	if rf, ok := ret.Get(3).(func(int, int, repository.ConfigType) error); ok {
		r3 = rf(pipelineId, wfrId, configType)
	} else {
		r3 = ret.Error(3)
	}

	return r0, r1, r2, r3
}

// GetDeployedHistoryDetailForCMCSByPipelineIdAndWfrId provides a mock function with given fields: pipelineId, wfrId, configType, userHasAdminAccess
func (_m *ConfigMapHistoryService) GetDeployedHistoryDetailForCMCSByPipelineIdAndWfrId(pipelineId int, wfrId int, configType repository.ConfigType, userHasAdminAccess bool) ([]*history.ComponentLevelHistoryDetailDto, error) {
	ret := _m.Called(pipelineId, wfrId, configType, userHasAdminAccess)

	var r0 []*history.ComponentLevelHistoryDetailDto
	var r1 error
	if rf, ok := ret.Get(0).(func(int, int, repository.ConfigType, bool) ([]*history.ComponentLevelHistoryDetailDto, error)); ok {
		return rf(pipelineId, wfrId, configType, userHasAdminAccess)
	}
	if rf, ok := ret.Get(0).(func(int, int, repository.ConfigType, bool) []*history.ComponentLevelHistoryDetailDto); ok {
		r0 = rf(pipelineId, wfrId, configType, userHasAdminAccess)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*history.ComponentLevelHistoryDetailDto)
		}
	}

	if rf, ok := ret.Get(1).(func(int, int, repository.ConfigType, bool) error); ok {
		r1 = rf(pipelineId, wfrId, configType, userHasAdminAccess)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDeployedHistoryList provides a mock function with given fields: pipelineId, baseConfigId, configType, componentName
func (_m *ConfigMapHistoryService) GetDeployedHistoryList(pipelineId int, baseConfigId int, configType repository.ConfigType, componentName string) ([]*history.DeployedHistoryComponentMetadataDto, error) {
	ret := _m.Called(pipelineId, baseConfigId, configType, componentName)

	var r0 []*history.DeployedHistoryComponentMetadataDto
	var r1 error
	if rf, ok := ret.Get(0).(func(int, int, repository.ConfigType, string) ([]*history.DeployedHistoryComponentMetadataDto, error)); ok {
		return rf(pipelineId, baseConfigId, configType, componentName)
	}
	if rf, ok := ret.Get(0).(func(int, int, repository.ConfigType, string) []*history.DeployedHistoryComponentMetadataDto); ok {
		r0 = rf(pipelineId, baseConfigId, configType, componentName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*history.DeployedHistoryComponentMetadataDto)
		}
	}

	if rf, ok := ret.Get(1).(func(int, int, repository.ConfigType, string) error); ok {
		r1 = rf(pipelineId, baseConfigId, configType, componentName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDeploymentDetailsForDeployedCMCSHistory provides a mock function with given fields: pipelineId, configType
func (_m *ConfigMapHistoryService) GetDeploymentDetailsForDeployedCMCSHistory(pipelineId int, configType repository.ConfigType) ([]*history.ConfigMapAndSecretHistoryDto, error) {
	ret := _m.Called(pipelineId, configType)

	var r0 []*history.ConfigMapAndSecretHistoryDto
	var r1 error
	if rf, ok := ret.Get(0).(func(int, repository.ConfigType) ([]*history.ConfigMapAndSecretHistoryDto, error)); ok {
		return rf(pipelineId, configType)
	}
	if rf, ok := ret.Get(0).(func(int, repository.ConfigType) []*history.ConfigMapAndSecretHistoryDto); ok {
		r0 = rf(pipelineId, configType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*history.ConfigMapAndSecretHistoryDto)
		}
	}

	if rf, ok := ret.Get(1).(func(int, repository.ConfigType) error); ok {
		r1 = rf(pipelineId, configType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetHistoryForDeployedCMCSById provides a mock function with given fields: id, pipelineId, configType, componentName, userHasAdminAccess
func (_m *ConfigMapHistoryService) GetHistoryForDeployedCMCSById(id int, pipelineId int, configType repository.ConfigType, componentName string, userHasAdminAccess bool) (*history.HistoryDetailDto, error) {
	ret := _m.Called(id, pipelineId, configType, componentName, userHasAdminAccess)

	var r0 *history.HistoryDetailDto
	var r1 error
	if rf, ok := ret.Get(0).(func(int, int, repository.ConfigType, string, bool) (*history.HistoryDetailDto, error)); ok {
		return rf(id, pipelineId, configType, componentName, userHasAdminAccess)
	}
	if rf, ok := ret.Get(0).(func(int, int, repository.ConfigType, string, bool) *history.HistoryDetailDto); ok {
		r0 = rf(id, pipelineId, configType, componentName, userHasAdminAccess)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*history.HistoryDetailDto)
		}
	}

	if rf, ok := ret.Get(1).(func(int, int, repository.ConfigType, string, bool) error); ok {
		r1 = rf(id, pipelineId, configType, componentName, userHasAdminAccess)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MergeAppLevelAndEnvLevelConfigs provides a mock function with given fields: appLevelConfig, envLevelConfig, configType, configMapSecretNames
func (_m *ConfigMapHistoryService) MergeAppLevelAndEnvLevelConfigs(appLevelConfig *chartConfig.ConfigMapAppModel, envLevelConfig *chartConfig.ConfigMapEnvModel, configType repository.ConfigType, configMapSecretNames []string) (string, error) {
	ret := _m.Called(appLevelConfig, envLevelConfig, configType, configMapSecretNames)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(*chartConfig.ConfigMapAppModel, *chartConfig.ConfigMapEnvModel, repository.ConfigType, []string) (string, error)); ok {
		return rf(appLevelConfig, envLevelConfig, configType, configMapSecretNames)
	}
	if rf, ok := ret.Get(0).(func(*chartConfig.ConfigMapAppModel, *chartConfig.ConfigMapEnvModel, repository.ConfigType, []string) string); ok {
		r0 = rf(appLevelConfig, envLevelConfig, configType, configMapSecretNames)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(*chartConfig.ConfigMapAppModel, *chartConfig.ConfigMapEnvModel, repository.ConfigType, []string) error); ok {
		r1 = rf(appLevelConfig, envLevelConfig, configType, configMapSecretNames)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewConfigMapHistoryService interface {
	mock.TestingT
	Cleanup(func())
}

// NewConfigMapHistoryService creates a new instance of ConfigMapHistoryService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewConfigMapHistoryService(t mockConstructorTestingTNewConfigMapHistoryService) *ConfigMapHistoryService {
	mock := &ConfigMapHistoryService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}