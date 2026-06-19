// Package mock_domain is a generated GoMock package.
package mock_domain

import (
	context "context"
	reflect "reflect"
	domain "warehouse-service/internal/domain"

	gomock "github.com/golang/mock/gomock"
)

// MockWarehouseRepository is a mock of WarehouseRepository interface.
type MockWarehouseRepository struct {
	ctrl     *gomock.Controller
	recorder *MockWarehouseRepositoryMockRecorder
}

// MockWarehouseRepositoryMockRecorder is the mock recorder for MockWarehouseRepository.
type MockWarehouseRepositoryMockRecorder struct {
	mock *MockWarehouseRepository
}

// NewMockWarehouseRepository creates a new mock instance.
func NewMockWarehouseRepository(ctrl *gomock.Controller) *MockWarehouseRepository {
	mock := &MockWarehouseRepository{ctrl: ctrl}
	mock.recorder = &MockWarehouseRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWarehouseRepository) EXPECT() *MockWarehouseRepositoryMockRecorder {
	return m.recorder
}

// AddPackageToManifest mocks base method.
func (m *MockWarehouseRepository) AddPackageToManifest(ctx context.Context, awb, manifestID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddPackageToManifest", ctx, awb, manifestID)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddPackageToManifest indicates an expected call of AddPackageToManifest.
func (mr *MockWarehouseRepositoryMockRecorder) AddPackageToManifest(ctx, awb, manifestID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddPackageToManifest", reflect.TypeOf((*MockWarehouseRepository)(nil).AddPackageToManifest), ctx, awb, manifestID)
}

// CreateManifest mocks base method.
func (m *MockWarehouseRepository) CreateManifest(ctx context.Context, manifest *domain.Manifest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateManifest", ctx, manifest)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateManifest indicates an expected call of CreateManifest.
func (mr *MockWarehouseRepositoryMockRecorder) CreateManifest(ctx, manifest interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateManifest", reflect.TypeOf((*MockWarehouseRepository)(nil).CreateManifest), ctx, manifest)
}

// DispatchManifest mocks base method.
func (m *MockWarehouseRepository) DispatchManifest(ctx context.Context, manifestID string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DispatchManifest", ctx, manifestID)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DispatchManifest indicates an expected call of DispatchManifest.
func (mr *MockWarehouseRepositoryMockRecorder) DispatchManifest(ctx, manifestID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DispatchManifest", reflect.TypeOf((*MockWarehouseRepository)(nil).DispatchManifest), ctx, manifestID)
}

// GetManifestByID mocks base method.
func (m *MockWarehouseRepository) GetManifestByID(ctx context.Context, id string) (*domain.Manifest, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetManifestByID", ctx, id)
	ret0, _ := ret[0].(*domain.Manifest)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetManifestByID indicates an expected call of GetManifestByID.
func (mr *MockWarehouseRepositoryMockRecorder) GetManifestByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetManifestByID", reflect.TypeOf((*MockWarehouseRepository)(nil).GetManifestByID), ctx, id)
}

// GetPackageByAWB mocks base method.
func (m *MockWarehouseRepository) GetPackageByAWB(ctx context.Context, awb string) (*domain.Package, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPackageByAWB", ctx, awb)
	ret0, _ := ret[0].(*domain.Package)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPackageByAWB indicates an expected call of GetPackageByAWB.
func (mr *MockWarehouseRepositoryMockRecorder) GetPackageByAWB(ctx, awb interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPackageByAWB", reflect.TypeOf((*MockWarehouseRepository)(nil).GetPackageByAWB), ctx, awb)
}

// SavePackage mocks base method.
func (m *MockWarehouseRepository) SavePackage(ctx context.Context, pkg *domain.Package) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SavePackage", ctx, pkg)
	ret0, _ := ret[0].(error)
	return ret0
}

// SavePackage indicates an expected call of SavePackage.
func (mr *MockWarehouseRepositoryMockRecorder) SavePackage(ctx, pkg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SavePackage", reflect.TypeOf((*MockWarehouseRepository)(nil).SavePackage), ctx, pkg)
}

// UpdatePackageStatus mocks base method.
func (m *MockWarehouseRepository) UpdatePackageStatus(ctx context.Context, awb, status string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePackageStatus", ctx, awb, status)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdatePackageStatus indicates an expected call of UpdatePackageStatus.
func (mr *MockWarehouseRepositoryMockRecorder) UpdatePackageStatus(ctx, awb, status interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePackageStatus", reflect.TypeOf((*MockWarehouseRepository)(nil).UpdatePackageStatus), ctx, awb, status)
}

// MockKafkaProducer is a mock of KafkaProducer interface.
type MockKafkaProducer struct {
	ctrl     *gomock.Controller
	recorder *MockKafkaProducerMockRecorder
}

// MockKafkaProducerMockRecorder is the mock recorder for MockKafkaProducer.
type MockKafkaProducerMockRecorder struct {
	mock *MockKafkaProducer
}

// NewMockKafkaProducer creates a new mock instance.
func NewMockKafkaProducer(ctrl *gomock.Controller) *MockKafkaProducer {
	mock := &MockKafkaProducer{ctrl: ctrl}
	mock.recorder = &MockKafkaProducerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockKafkaProducer) EXPECT() *MockKafkaProducerMockRecorder {
	return m.recorder
}

// PublishEvent mocks base method.
func (m *MockKafkaProducer) PublishEvent(ctx context.Context, topic, key string, value []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PublishEvent", ctx, topic, key, value)
	ret0, _ := ret[0].(error)
	return ret0
}

// PublishEvent indicates an expected call of PublishEvent.
func (mr *MockKafkaProducerMockRecorder) PublishEvent(ctx, topic, key, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PublishEvent", reflect.TypeOf((*MockKafkaProducer)(nil).PublishEvent), ctx, topic, key, value)
}
