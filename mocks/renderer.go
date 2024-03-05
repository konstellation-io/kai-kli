// Code generated by MockGen. DO NOT EDIT.
// Source: interface.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	kai "github.com/konstellation-io/kli/api/kai"
	entity "github.com/konstellation-io/kli/internal/entity"
	configuration "github.com/konstellation-io/kli/internal/services/configuration"
	krt "github.com/konstellation-io/krt/pkg/krt"
)

// MockRenderer is a mock of Renderer interface.
type MockRenderer struct {
	ctrl     *gomock.Controller
	recorder *MockRendererMockRecorder
}

// MockRendererMockRecorder is the mock recorder for MockRenderer.
type MockRendererMockRecorder struct {
	mock *MockRenderer
}

// NewMockRenderer creates a new mock instance.
func NewMockRenderer(ctrl *gomock.Controller) *MockRenderer {
	mock := &MockRenderer{ctrl: ctrl}
	mock.recorder = &MockRendererMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRenderer) EXPECT() *MockRendererMockRecorder {
	return m.recorder
}

// RenderAddMaintainerToProduct mocks base method.
func (m *MockRenderer) RenderAddMaintainerToProduct(product, userEmail string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RenderAddMaintainerToProduct", product, userEmail)
}

// RenderAddMaintainerToProduct indicates an expected call of RenderAddMaintainerToProduct.
func (mr *MockRendererMockRecorder) RenderAddMaintainerToProduct(product, userEmail interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RenderAddMaintainerToProduct", reflect.TypeOf((*MockRenderer)(nil).RenderAddMaintainerToProduct), product, userEmail)
}

// RenderAddUserToProduct mocks base method.
func (m *MockRenderer) RenderAddUserToProduct(product, userEmail string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RenderAddUserToProduct", product, userEmail)
}

// RenderAddUserToProduct indicates an expected call of RenderAddUserToProduct.
func (mr *MockRendererMockRecorder) RenderAddUserToProduct(product, userEmail interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RenderAddUserToProduct", reflect.TypeOf((*MockRenderer)(nil).RenderAddUserToProduct), product, userEmail)
}

// RenderCallout mocks base method.
func (m *MockRenderer) RenderCallout(v *entity.Version) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RenderCallout", v)
}

// RenderCallout indicates an expected call of RenderCallout.
func (mr *MockRendererMockRecorder) RenderCallout(v interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RenderCallout", reflect.TypeOf((*MockRenderer)(nil).RenderCallout), v)
}

// RenderConfiguration mocks base method.
func (m *MockRenderer) RenderConfiguration(scope string, config map[string]string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RenderConfiguration", scope, config)
}

// RenderConfiguration indicates an expected call of RenderConfiguration.
func (mr *MockRendererMockRecorder) RenderConfiguration(scope, config interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RenderConfiguration", reflect.TypeOf((*MockRenderer)(nil).RenderConfiguration), scope, config)
}

// RenderKliVersion mocks base method.
func (m *MockRenderer) RenderKliVersion(version, buildDate string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RenderKliVersion", version, buildDate)
}

// RenderKliVersion indicates an expected call of RenderKliVersion.
func (mr *MockRendererMockRecorder) RenderKliVersion(version, buildDate interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RenderKliVersion", reflect.TypeOf((*MockRenderer)(nil).RenderKliVersion), version, buildDate)
}

// RenderLogin mocks base method.
func (m *MockRenderer) RenderLogin(serverName string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RenderLogin", serverName)
}

// RenderLogin indicates an expected call of RenderLogin.
func (mr *MockRendererMockRecorder) RenderLogin(serverName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RenderLogin", reflect.TypeOf((*MockRenderer)(nil).RenderLogin), serverName)
}

// RenderLogout mocks base method.
func (m *MockRenderer) RenderLogout(serverName string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RenderLogout", serverName)
}

// RenderLogout indicates an expected call of RenderLogout.
func (mr *MockRendererMockRecorder) RenderLogout(serverName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RenderLogout", reflect.TypeOf((*MockRenderer)(nil).RenderLogout), serverName)
}

// RenderLogs mocks base method.
func (m *MockRenderer) RenderLogs(productID string, logs []entity.Log, outFormat entity.LogOutFormat, showAllLabels bool) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RenderLogs", productID, logs, outFormat, showAllLabels)
}

// RenderLogs indicates an expected call of RenderLogs.
func (mr *MockRendererMockRecorder) RenderLogs(productID, logs, outFormat, showAllLabels interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RenderLogs", reflect.TypeOf((*MockRenderer)(nil).RenderLogs), productID, logs, outFormat, showAllLabels)
}

// RenderProcessDeleted mocks base method.
func (m *MockRenderer) RenderProcessDeleted(process string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RenderProcessDeleted", process)
}

// RenderProcessDeleted indicates an expected call of RenderProcessDeleted.
func (mr *MockRendererMockRecorder) RenderProcessDeleted(process interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RenderProcessDeleted", reflect.TypeOf((*MockRenderer)(nil).RenderProcessDeleted), process)
}

// RenderProcessRegistered mocks base method.
func (m *MockRenderer) RenderProcessRegistered(process *entity.RegisteredProcess) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RenderProcessRegistered", process)
}

// RenderProcessRegistered indicates an expected call of RenderProcessRegistered.
func (mr *MockRendererMockRecorder) RenderProcessRegistered(process interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RenderProcessRegistered", reflect.TypeOf((*MockRenderer)(nil).RenderProcessRegistered), process)
}

// RenderProcesses mocks base method.
func (m *MockRenderer) RenderProcesses(processes []krt.Process) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RenderProcesses", processes)
}

// RenderProcesses indicates an expected call of RenderProcesses.
func (mr *MockRendererMockRecorder) RenderProcesses(processes interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RenderProcesses", reflect.TypeOf((*MockRenderer)(nil).RenderProcesses), processes)
}

// RenderProductBinded mocks base method.
func (m *MockRenderer) RenderProductBinded(productID string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RenderProductBinded", productID)
}

// RenderProductBinded indicates an expected call of RenderProductBinded.
func (mr *MockRendererMockRecorder) RenderProductBinded(productID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RenderProductBinded", reflect.TypeOf((*MockRenderer)(nil).RenderProductBinded), productID)
}

// RenderProductCreated mocks base method.
func (m *MockRenderer) RenderProductCreated(product string, server *configuration.Server, initLocal bool) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RenderProductCreated", product, server, initLocal)
}

// RenderProductCreated indicates an expected call of RenderProductCreated.
func (mr *MockRendererMockRecorder) RenderProductCreated(product, server, initLocal interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RenderProductCreated", reflect.TypeOf((*MockRenderer)(nil).RenderProductCreated), product, server, initLocal)
}

// RenderProducts mocks base method.
func (m *MockRenderer) RenderProducts(products []kai.Product) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RenderProducts", products)
}

// RenderProducts indicates an expected call of RenderProducts.
func (mr *MockRendererMockRecorder) RenderProducts(products interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RenderProducts", reflect.TypeOf((*MockRenderer)(nil).RenderProducts), products)
}

// RenderPublishVersion mocks base method.
func (m *MockRenderer) RenderPublishVersion(product, versionTag string, triggers []entity.TriggerEndpoint) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RenderPublishVersion", product, versionTag, triggers)
}

// RenderPublishVersion indicates an expected call of RenderPublishVersion.
func (mr *MockRendererMockRecorder) RenderPublishVersion(product, versionTag, triggers interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RenderPublishVersion", reflect.TypeOf((*MockRenderer)(nil).RenderPublishVersion), product, versionTag, triggers)
}

// RenderPushVersion mocks base method.
func (m *MockRenderer) RenderPushVersion(product, versionTag string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RenderPushVersion", product, versionTag)
}

// RenderPushVersion indicates an expected call of RenderPushVersion.
func (mr *MockRendererMockRecorder) RenderPushVersion(product, versionTag interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RenderPushVersion", reflect.TypeOf((*MockRenderer)(nil).RenderPushVersion), product, versionTag)
}

// RenderRegisteredProcesses mocks base method.
func (m *MockRenderer) RenderRegisteredProcesses(registries []*entity.RegisteredProcess) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RenderRegisteredProcesses", registries)
}

// RenderRegisteredProcesses indicates an expected call of RenderRegisteredProcesses.
func (mr *MockRendererMockRecorder) RenderRegisteredProcesses(registries interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RenderRegisteredProcesses", reflect.TypeOf((*MockRenderer)(nil).RenderRegisteredProcesses), registries)
}

// RenderRemoveMaintainerFromProduct mocks base method.
func (m *MockRenderer) RenderRemoveMaintainerFromProduct(product, userEmail string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RenderRemoveMaintainerFromProduct", product, userEmail)
}

// RenderRemoveMaintainerFromProduct indicates an expected call of RenderRemoveMaintainerFromProduct.
func (mr *MockRendererMockRecorder) RenderRemoveMaintainerFromProduct(product, userEmail interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RenderRemoveMaintainerFromProduct", reflect.TypeOf((*MockRenderer)(nil).RenderRemoveMaintainerFromProduct), product, userEmail)
}

// RenderRemoveUserFromProduct mocks base method.
func (m *MockRenderer) RenderRemoveUserFromProduct(product, userEmail string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RenderRemoveUserFromProduct", product, userEmail)
}

// RenderRemoveUserFromProduct indicates an expected call of RenderRemoveUserFromProduct.
func (mr *MockRendererMockRecorder) RenderRemoveUserFromProduct(product, userEmail interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RenderRemoveUserFromProduct", reflect.TypeOf((*MockRenderer)(nil).RenderRemoveUserFromProduct), product, userEmail)
}

// RenderServers mocks base method.
func (m *MockRenderer) RenderServers(servers []*configuration.Server) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RenderServers", servers)
}

// RenderServers indicates an expected call of RenderServers.
func (mr *MockRendererMockRecorder) RenderServers(servers interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RenderServers", reflect.TypeOf((*MockRenderer)(nil).RenderServers), servers)
}

// RenderStartVersion mocks base method.
func (m *MockRenderer) RenderStartVersion(product, versionTag string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RenderStartVersion", product, versionTag)
}

// RenderStartVersion indicates an expected call of RenderStartVersion.
func (mr *MockRendererMockRecorder) RenderStartVersion(product, versionTag interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RenderStartVersion", reflect.TypeOf((*MockRenderer)(nil).RenderStartVersion), product, versionTag)
}

// RenderStopVersion mocks base method.
func (m *MockRenderer) RenderStopVersion(product, versionTag string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RenderStopVersion", product, versionTag)
}

// RenderStopVersion indicates an expected call of RenderStopVersion.
func (mr *MockRendererMockRecorder) RenderStopVersion(product, versionTag interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RenderStopVersion", reflect.TypeOf((*MockRenderer)(nil).RenderStopVersion), product, versionTag)
}

// RenderUnpublishVersion mocks base method.
func (m *MockRenderer) RenderUnpublishVersion(product, versionTag string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RenderUnpublishVersion", product, versionTag)
}

// RenderUnpublishVersion indicates an expected call of RenderUnpublishVersion.
func (mr *MockRendererMockRecorder) RenderUnpublishVersion(product, versionTag interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RenderUnpublishVersion", reflect.TypeOf((*MockRenderer)(nil).RenderUnpublishVersion), product, versionTag)
}

// RenderVersion mocks base method.
func (m *MockRenderer) RenderVersion(productID string, v *entity.Version) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RenderVersion", productID, v)
}

// RenderVersion indicates an expected call of RenderVersion.
func (mr *MockRendererMockRecorder) RenderVersion(productID, v interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RenderVersion", reflect.TypeOf((*MockRenderer)(nil).RenderVersion), productID, v)
}

// RenderVersions mocks base method.
func (m *MockRenderer) RenderVersions(productID string, versions []*entity.Version) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RenderVersions", productID, versions)
}

// RenderVersions indicates an expected call of RenderVersions.
func (mr *MockRendererMockRecorder) RenderVersions(productID, versions interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RenderVersions", reflect.TypeOf((*MockRenderer)(nil).RenderVersions), productID, versions)
}

// RenderWorkflows mocks base method.
func (m *MockRenderer) RenderWorkflows(workflows []krt.Workflow) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RenderWorkflows", workflows)
}

// RenderWorkflows indicates an expected call of RenderWorkflows.
func (mr *MockRendererMockRecorder) RenderWorkflows(workflows interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RenderWorkflows", reflect.TypeOf((*MockRenderer)(nil).RenderWorkflows), workflows)
}
