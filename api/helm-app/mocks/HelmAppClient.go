// Code generated by mockery v2.31.4. DO NOT EDIT.

package mocks

import (
	context "context"
	"github.com/devtron-labs/devtron/api/helm-app/gRPC"
	mock "github.com/stretchr/testify/mock"
)

// HelmAppClient is an autogenerated mock type for the HelmAppClient type
type HelmAppClient struct {
	mock.Mock
}

// DeleteApplication provides a mock function with given fields: ctx, in
func (_m *HelmAppClient) DeleteApplication(ctx context.Context, in *gRPC.ReleaseIdentifier) (*gRPC.UninstallReleaseResponse, error) {
	ret := _m.Called(ctx, in)

	var r0 *gRPC.UninstallReleaseResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *gRPC.ReleaseIdentifier) (*gRPC.UninstallReleaseResponse, error)); ok {
		return rf(ctx, in)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *gRPC.ReleaseIdentifier) *gRPC.UninstallReleaseResponse); ok {
		r0 = rf(ctx, in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gRPC.UninstallReleaseResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *gRPC.ReleaseIdentifier) error); ok {
		r1 = rf(ctx, in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAppDetail provides a mock function with given fields: ctx, in
func (_m *HelmAppClient) GetAppDetail(ctx context.Context, in *gRPC.AppDetailRequest) (*gRPC.AppDetail, error) {
	ret := _m.Called(ctx, in)

	var r0 *gRPC.AppDetail
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *gRPC.AppDetailRequest) (*gRPC.AppDetail, error)); ok {
		return rf(ctx, in)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *gRPC.AppDetailRequest) *gRPC.AppDetail); ok {
		r0 = rf(ctx, in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gRPC.AppDetail)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *gRPC.AppDetailRequest) error); ok {
		r1 = rf(ctx, in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAppStatus provides a mock function with given fields: ctx, in
func (_m *HelmAppClient) GetAppStatus(ctx context.Context, in *gRPC.AppDetailRequest) (*gRPC.AppStatus, error) {
	ret := _m.Called(ctx, in)

	var r0 *gRPC.AppStatus
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *gRPC.AppDetailRequest) (*gRPC.AppStatus, error)); ok {
		return rf(ctx, in)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *gRPC.AppDetailRequest) *gRPC.AppStatus); ok {
		r0 = rf(ctx, in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gRPC.AppStatus)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *gRPC.AppDetailRequest) error); ok {
		r1 = rf(ctx, in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDeploymentDetail provides a mock function with given fields: ctx, in
func (_m *HelmAppClient) GetDeploymentDetail(ctx context.Context, in *gRPC.DeploymentDetailRequest) (*gRPC.DeploymentDetailResponse, error) {
	ret := _m.Called(ctx, in)

	var r0 *gRPC.DeploymentDetailResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *gRPC.DeploymentDetailRequest) (*gRPC.DeploymentDetailResponse, error)); ok {
		return rf(ctx, in)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *gRPC.DeploymentDetailRequest) *gRPC.DeploymentDetailResponse); ok {
		r0 = rf(ctx, in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gRPC.DeploymentDetailResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *gRPC.DeploymentDetailRequest) error); ok {
		r1 = rf(ctx, in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDeploymentHistory provides a mock function with given fields: ctx, in
func (_m *HelmAppClient) GetDeploymentHistory(ctx context.Context, in *gRPC.AppDetailRequest) (*gRPC.HelmAppDeploymentHistory, error) {
	ret := _m.Called(ctx, in)

	var r0 *gRPC.HelmAppDeploymentHistory
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *gRPC.AppDetailRequest) (*gRPC.HelmAppDeploymentHistory, error)); ok {
		return rf(ctx, in)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *gRPC.AppDetailRequest) *gRPC.HelmAppDeploymentHistory); ok {
		r0 = rf(ctx, in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gRPC.HelmAppDeploymentHistory)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *gRPC.AppDetailRequest) error); ok {
		r1 = rf(ctx, in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDesiredManifest provides a mock function with given fields: ctx, in
func (_m *HelmAppClient) GetDesiredManifest(ctx context.Context, in *gRPC.ObjectRequest) (*gRPC.DesiredManifestResponse, error) {
	ret := _m.Called(ctx, in)

	var r0 *gRPC.DesiredManifestResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *gRPC.ObjectRequest) (*gRPC.DesiredManifestResponse, error)); ok {
		return rf(ctx, in)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *gRPC.ObjectRequest) *gRPC.DesiredManifestResponse); ok {
		r0 = rf(ctx, in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gRPC.DesiredManifestResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *gRPC.ObjectRequest) error); ok {
		r1 = rf(ctx, in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetNotes provides a mock function with given fields: ctx, request
func (_m *HelmAppClient) GetNotes(ctx context.Context, request *gRPC.InstallReleaseRequest) (*gRPC.ChartNotesResponse, error) {
	ret := _m.Called(ctx, request)

	var r0 *gRPC.ChartNotesResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *gRPC.InstallReleaseRequest) (*gRPC.ChartNotesResponse, error)); ok {
		return rf(ctx, request)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *gRPC.InstallReleaseRequest) *gRPC.ChartNotesResponse); ok {
		r0 = rf(ctx, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gRPC.ChartNotesResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *gRPC.InstallReleaseRequest) error); ok {
		r1 = rf(ctx, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetValuesYaml provides a mock function with given fields: ctx, in
func (_m *HelmAppClient) GetValuesYaml(ctx context.Context, in *gRPC.AppDetailRequest) (*gRPC.ReleaseInfo, error) {
	ret := _m.Called(ctx, in)

	var r0 *gRPC.ReleaseInfo
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *gRPC.AppDetailRequest) (*gRPC.ReleaseInfo, error)); ok {
		return rf(ctx, in)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *gRPC.AppDetailRequest) *gRPC.ReleaseInfo); ok {
		r0 = rf(ctx, in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gRPC.ReleaseInfo)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *gRPC.AppDetailRequest) error); ok {
		r1 = rf(ctx, in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Hibernate provides a mock function with given fields: ctx, in
func (_m *HelmAppClient) Hibernate(ctx context.Context, in *gRPC.HibernateRequest) (*gRPC.HibernateResponse, error) {
	ret := _m.Called(ctx, in)

	var r0 *gRPC.HibernateResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *gRPC.HibernateRequest) (*gRPC.HibernateResponse, error)); ok {
		return rf(ctx, in)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *gRPC.HibernateRequest) *gRPC.HibernateResponse); ok {
		r0 = rf(ctx, in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gRPC.HibernateResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *gRPC.HibernateRequest) error); ok {
		r1 = rf(ctx, in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InstallRelease provides a mock function with given fields: ctx, in
func (_m *HelmAppClient) InstallRelease(ctx context.Context, in *gRPC.InstallReleaseRequest) (*gRPC.InstallReleaseResponse, error) {
	ret := _m.Called(ctx, in)

	var r0 *gRPC.InstallReleaseResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *gRPC.InstallReleaseRequest) (*gRPC.InstallReleaseResponse, error)); ok {
		return rf(ctx, in)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *gRPC.InstallReleaseRequest) *gRPC.InstallReleaseResponse); ok {
		r0 = rf(ctx, in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gRPC.InstallReleaseResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *gRPC.InstallReleaseRequest) error); ok {
		r1 = rf(ctx, in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InstallReleaseWithCustomChart provides a mock function with given fields: ctx, in
func (_m *HelmAppClient) InstallReleaseWithCustomChart(ctx context.Context, in *gRPC.HelmInstallCustomRequest) (*gRPC.HelmInstallCustomResponse, error) {
	ret := _m.Called(ctx, in)

	var r0 *gRPC.HelmInstallCustomResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *gRPC.HelmInstallCustomRequest) (*gRPC.HelmInstallCustomResponse, error)); ok {
		return rf(ctx, in)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *gRPC.HelmInstallCustomRequest) *gRPC.HelmInstallCustomResponse); ok {
		r0 = rf(ctx, in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gRPC.HelmInstallCustomResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *gRPC.HelmInstallCustomRequest) error); ok {
		r1 = rf(ctx, in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsReleaseInstalled provides a mock function with given fields: ctx, in
func (_m *HelmAppClient) IsReleaseInstalled(ctx context.Context, in *gRPC.ReleaseIdentifier) (*gRPC.BooleanResponse, error) {
	ret := _m.Called(ctx, in)

	var r0 *gRPC.BooleanResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *gRPC.ReleaseIdentifier) (*gRPC.BooleanResponse, error)); ok {
		return rf(ctx, in)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *gRPC.ReleaseIdentifier) *gRPC.BooleanResponse); ok {
		r0 = rf(ctx, in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gRPC.BooleanResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *gRPC.ReleaseIdentifier) error); ok {
		r1 = rf(ctx, in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListApplication provides a mock function with given fields: ctx, req
func (_m *HelmAppClient) ListApplication(ctx context.Context, req *gRPC.AppListRequest) (gRPC.ApplicationService_ListApplicationsClient, error) {
	ret := _m.Called(ctx, req)

	var r0 gRPC.ApplicationService_ListApplicationsClient
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *gRPC.AppListRequest) (gRPC.ApplicationService_ListApplicationsClient, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *gRPC.AppListRequest) gRPC.ApplicationService_ListApplicationsClient); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(gRPC.ApplicationService_ListApplicationsClient)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *gRPC.AppListRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RollbackRelease provides a mock function with given fields: ctx, in
func (_m *HelmAppClient) RollbackRelease(ctx context.Context, in *gRPC.RollbackReleaseRequest) (*gRPC.BooleanResponse, error) {
	ret := _m.Called(ctx, in)

	var r0 *gRPC.BooleanResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *gRPC.RollbackReleaseRequest) (*gRPC.BooleanResponse, error)); ok {
		return rf(ctx, in)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *gRPC.RollbackReleaseRequest) *gRPC.BooleanResponse); ok {
		r0 = rf(ctx, in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gRPC.BooleanResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *gRPC.RollbackReleaseRequest) error); ok {
		r1 = rf(ctx, in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TemplateChart provides a mock function with given fields: ctx, in
func (_m *HelmAppClient) TemplateChart(ctx context.Context, in *gRPC.InstallReleaseRequest) (*gRPC.TemplateChartResponse, error) {
	ret := _m.Called(ctx, in)

	var r0 *gRPC.TemplateChartResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *gRPC.InstallReleaseRequest) (*gRPC.TemplateChartResponse, error)); ok {
		return rf(ctx, in)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *gRPC.InstallReleaseRequest) *gRPC.TemplateChartResponse); ok {
		r0 = rf(ctx, in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gRPC.TemplateChartResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *gRPC.InstallReleaseRequest) error); ok {
		r1 = rf(ctx, in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UnHibernate provides a mock function with given fields: ctx, in
func (_m *HelmAppClient) UnHibernate(ctx context.Context, in *gRPC.HibernateRequest) (*gRPC.HibernateResponse, error) {
	ret := _m.Called(ctx, in)

	var r0 *gRPC.HibernateResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *gRPC.HibernateRequest) (*gRPC.HibernateResponse, error)); ok {
		return rf(ctx, in)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *gRPC.HibernateRequest) *gRPC.HibernateResponse); ok {
		r0 = rf(ctx, in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gRPC.HibernateResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *gRPC.HibernateRequest) error); ok {
		r1 = rf(ctx, in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateApplication provides a mock function with given fields: ctx, in
func (_m *HelmAppClient) UpdateApplication(ctx context.Context, in *gRPC.UpgradeReleaseRequest) (*gRPC.UpgradeReleaseResponse, error) {
	ret := _m.Called(ctx, in)

	var r0 *gRPC.UpgradeReleaseResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *gRPC.UpgradeReleaseRequest) (*gRPC.UpgradeReleaseResponse, error)); ok {
		return rf(ctx, in)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *gRPC.UpgradeReleaseRequest) *gRPC.UpgradeReleaseResponse); ok {
		r0 = rf(ctx, in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gRPC.UpgradeReleaseResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *gRPC.UpgradeReleaseRequest) error); ok {
		r1 = rf(ctx, in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateApplicationWithChartInfo provides a mock function with given fields: ctx, in
func (_m *HelmAppClient) UpdateApplicationWithChartInfo(ctx context.Context, in *gRPC.InstallReleaseRequest) (*gRPC.UpgradeReleaseResponse, error) {
	ret := _m.Called(ctx, in)

	var r0 *gRPC.UpgradeReleaseResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *gRPC.InstallReleaseRequest) (*gRPC.UpgradeReleaseResponse, error)); ok {
		return rf(ctx, in)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *gRPC.InstallReleaseRequest) *gRPC.UpgradeReleaseResponse); ok {
		r0 = rf(ctx, in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gRPC.UpgradeReleaseResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *gRPC.InstallReleaseRequest) error); ok {
		r1 = rf(ctx, in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewHelmAppClient creates a new instance of HelmAppClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewHelmAppClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *HelmAppClient {
	mock := &HelmAppClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}