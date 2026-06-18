// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"
	domain "settlement-service/internal/domain"

	gomock "github.com/golang/mock/gomock"
)

// MockSettlementRepository is a mock of SettlementRepository interface.
type MockSettlementRepository struct {
	ctrl     *gomock.Controller
	recorder *MockSettlementRepositoryMockRecorder
}

// MockSettlementRepositoryMockRecorder is the mock recorder for MockSettlementRepository.
type MockSettlementRepositoryMockRecorder struct {
	mock *MockSettlementRepository
}

// NewMockSettlementRepository creates a new mock instance.
func NewMockSettlementRepository(ctrl *gomock.Controller) *MockSettlementRepository {
	mock := &MockSettlementRepository{ctrl: ctrl}
	mock.recorder = &MockSettlementRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSettlementRepository) EXPECT() *MockSettlementRepositoryMockRecorder {
	return m.recorder
}

// CreateCommissionLog mocks base method.
func (m *MockSettlementRepository) CreateCommissionLog(ctx context.Context, log *domain.CommissionLog) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCommissionLog", ctx, log)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateCommissionLog indicates an expected call of CreateCommissionLog.
func (mr *MockSettlementRepositoryMockRecorder) CreateCommissionLog(ctx, log interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCommissionLog", reflect.TypeOf((*MockSettlementRepository)(nil).CreateCommissionLog), ctx, log)
}

// GetCommissionByAWB mocks base method.
func (m *MockSettlementRepository) GetCommissionByAWB(ctx context.Context, awb string) (*domain.CommissionLog, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCommissionByAWB", ctx, awb)
	ret0, _ := ret[0].(*domain.CommissionLog)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCommissionByAWB indicates an expected call of GetCommissionByAWB.
func (mr *MockSettlementRepositoryMockRecorder) GetCommissionByAWB(ctx, awb interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCommissionByAWB", reflect.TypeOf((*MockSettlementRepository)(nil).GetCommissionByAWB), ctx, awb)
}

// GetCommissionsByCourier mocks base method.
func (m *MockSettlementRepository) GetCommissionsByCourier(ctx context.Context, courierID string) ([]domain.CommissionLog, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCommissionsByCourier", ctx, courierID)
	ret0, _ := ret[0].([]domain.CommissionLog)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCommissionsByCourier indicates an expected call of GetCommissionsByCourier.
func (mr *MockSettlementRepositoryMockRecorder) GetCommissionsByCourier(ctx, courierID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCommissionsByCourier", reflect.TypeOf((*MockSettlementRepository)(nil).GetCommissionsByCourier), ctx, courierID)
}

// GetCourierSummary mocks base method.
func (m *MockSettlementRepository) GetCourierSummary(ctx context.Context, courierID string) (*domain.CourierSummary, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCourierSummary", ctx, courierID)
	ret0, _ := ret[0].(*domain.CourierSummary)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCourierSummary indicates an expected call of GetCourierSummary.
func (mr *MockSettlementRepositoryMockRecorder) GetCourierSummary(ctx, courierID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCourierSummary", reflect.TypeOf((*MockSettlementRepository)(nil).GetCourierSummary), ctx, courierID)
}

// MarkAsPaid mocks base method.
func (m *MockSettlementRepository) MarkAsPaid(ctx context.Context, courierID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MarkAsPaid", ctx, courierID)
	ret0, _ := ret[0].(error)
	return ret0
}

// MarkAsPaid indicates an expected call of MarkAsPaid.
func (mr *MockSettlementRepositoryMockRecorder) MarkAsPaid(ctx, courierID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MarkAsPaid", reflect.TypeOf((*MockSettlementRepository)(nil).MarkAsPaid), ctx, courierID)
}

// MockPricingServiceClient is a mock of PricingServiceClient interface.
type MockPricingServiceClient struct {
	ctrl     *gomock.Controller
	recorder *MockPricingServiceClientMockRecorder
}

// MockPricingServiceClientMockRecorder is the mock recorder for MockPricingServiceClient.
type MockPricingServiceClientMockRecorder struct {
	mock *MockPricingServiceClient
}

// NewMockPricingServiceClient creates a new mock instance.
func NewMockPricingServiceClient(ctrl *gomock.Controller) *MockPricingServiceClient {
	mock := &MockPricingServiceClient{ctrl: ctrl}
	mock.recorder = &MockPricingServiceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPricingServiceClient) EXPECT() *MockPricingServiceClientMockRecorder {
	return m.recorder
}

// GetCommissionRate mocks base method.
func (m *MockPricingServiceClient) GetCommissionRate(ctx context.Context, serviceType string) (float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCommissionRate", ctx, serviceType)
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCommissionRate indicates an expected call of GetCommissionRate.
func (mr *MockPricingServiceClientMockRecorder) GetCommissionRate(ctx, serviceType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCommissionRate", reflect.TypeOf((*MockPricingServiceClient)(nil).GetCommissionRate), ctx, serviceType)
}
