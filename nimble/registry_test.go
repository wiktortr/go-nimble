package nimble

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"testing"
)

type mockRegistry struct {
	mock.Mock
}

func (m *mockRegistry) Logger() *zap.Logger {
	return m.Called().Get(0).(*zap.Logger)
}

func (m *mockRegistry) Instantiate(uri string) (ComponentImpl, error) {
	args := m.Called(uri)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(ComponentImpl), args.Error(1)
}

func (m *mockRegistry) RegisterRoute(_ *Route) error {
	return m.Called().Error(0)
}

func (m *mockRegistry) Start() error {
	return m.Called().Error(0)
}

func (m *mockRegistry) Stop() error {
	return m.Called().Error(0)
}

func TestRegistryImpl_Logger_ReturnsLogger(t *testing.T) {
	// given
	logger := zap.NewNop()
	reg := &RegistryImpl{logger: logger}
	// when
	result := reg.Logger()
	// then
	assert.Equal(t, logger, result)
}

func TestRegistryImpl_Instantiate_InvalidURI_ReturnsError(t *testing.T) {
	// given
	logger := zap.NewNop()
	reg := &RegistryImpl{logger: logger}
	// when
	inst, err := reg.Instantiate("invalid_uri")
	// then
	assert.Nil(t, inst)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid URI")
}

func TestRegistryImpl_Instantiate_ComponentNotFound_ReturnsError(t *testing.T) {
	// given
	logger := zap.NewNop()
	reg := &RegistryImpl{
		ComponentMap: map[string]Component{},
		InstancesMap: map[string]ComponentImpl{},
		logger:       logger,
	}
	uri := "comp://test"

	// when
	inst, err := reg.Instantiate(uri)

	// then
	assert.Nil(t, inst)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Cannot find Component")
}

func TestRegistryImpl_Instantiate_AlreadyInstantiated_ReturnsInstance(t *testing.T) {
	// given
	logger := zap.NewNop()
	comp := &mockComponent{}
	inst := &mockComponentImpl{}
	reg := &RegistryImpl{
		ComponentMap: map[string]Component{"test": comp},
		InstancesMap: map[string]ComponentImpl{"comp://test?id=1": inst},
		logger:       logger,
	}
	uri := "comp://test"

	// when
	result, err := reg.Instantiate(uri)

	// then
	assert.NoError(t, err)
	assert.Equal(t, inst, result)
}

func TestRegistryImpl_Instantiate_InstantiateError_ReturnsError(t *testing.T) {
	// given
	logger := zap.NewNop()
	comp := &mockComponent{}
	reg := &RegistryImpl{
		ComponentMap: map[string]Component{"comp": comp},
		InstancesMap: map[string]ComponentImpl{},
		logger:       logger,
	}
	uri := "comp://test"
	comp.On("Instantiate", reg, mock.Anything).Return(&mockComponentImpl{}, errors.New("fail"))

	// when
	inst, err := reg.Instantiate(uri)

	// then
	assert.Nil(t, inst)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Cannot Instantiate")
	comp.AssertExpectations(t)
}

func TestRegistryImpl_Instantiate_SuccessfulInstantiation(t *testing.T) {
	// given
	logger := zap.NewNop()
	comp := &mockComponent{}
	inst := &mockComponentImpl{}
	reg := &RegistryImpl{
		ComponentMap: map[string]Component{"test": comp},
		InstancesMap: map[string]ComponentImpl{},
		logger:       logger,
	}
	comp.On("Instantiate", reg, mock.Anything).Return(inst, nil)
	uri := "test:key=id"

	// when
	result, err := reg.Instantiate(uri)

	// then
	assert.NoError(t, err)
	assert.Equal(t, inst, result)
	comp.AssertExpectations(t)
}

func TestRegistryImpl_RegisterRoute_SuccessfulRegistration(t *testing.T) {
	// given
	logger := zap.NewNop()
	comp := &mockComponent{}
	inst := &mockComponentImpl{}
	reg := &RegistryImpl{
		ComponentMap: map[string]Component{"test": comp},
		InstancesMap: map[string]ComponentImpl{},
		RoutesMap:    map[string]*Route{},
		logger:       logger,
	}
	route := &Route{
		id:           "route1",
		dependencies: []string{"test://abc"},
		components:   make(map[string]ComponentImpl),
	}
	comp.On("Instantiate", reg, mock.Anything).Return(inst, nil)

	// when
	err := reg.RegisterRoute(route)

	// then
	assert.NoError(t, err)
	assert.Equal(t, reg, route.Registry)
	assert.Equal(t, route, reg.RoutesMap["route1"])
	assert.Equal(t, inst, route.components["test://abc"])
	comp.AssertExpectations(t)
}

func TestRegistryImpl_RegisterRoute_InstantiateError_ReturnsError(t *testing.T) {
	// given
	logger := zap.NewNop()
	comp := &mockComponent{}
	reg := &RegistryImpl{
		ComponentMap: map[string]Component{"test": comp},
		InstancesMap: map[string]ComponentImpl{},
		RoutesMap:    map[string]*Route{},
		logger:       logger,
	}
	route := &Route{
		id:           "route1",
		dependencies: []string{"test:key=id"},
		components:   make(map[string]ComponentImpl),
	}
	comp.On("Instantiate", reg, mock.Anything).Return(&mockComponentImpl{}, errors.New("fail"))

	// when
	err := reg.RegisterRoute(route)

	// then
	assert.Error(t, err)
	comp.AssertExpectations(t)
}

//func TestRegistryImpl_Start_StartsRoutesAndComponents(t *testing.T) {
//	// given
//	logger := zap.NewNop()
//	route := &mockRoute{id: "r1", Name: "route1"}
//	route.On("Start").Return()
//	comp := &mockComponentImpl{}
//	reg := &RegistryImpl{
//		RoutesMap:    map[string]*Route{"r1": (*Route)(route)},
//		InstancesMap: map[string]ComponentImpl{"c1": comp},
//		logger:       logger,
//	}
//
//	// when
//	err := reg.Start()
//
//	// then
//	assert.NoError(t, err)
//	route.AssertExpectations(t)
//}
//
//func TestRegistryImpl_Stop_StopsComponentsAndRoutes(t *testing.T) {
//	// given
//	logger := zap.NewNop()
//	route := &mockRoute{id: "r1", Name: "route1"}
//	route.On("Stop").Return()
//	comp := &mockComponentImplForRegistry{}
//	reg := &RegistryImpl{
//		RoutesMap:    map[string]*Route{"r1": (*Route)(route)},
//		InstancesMap: map[string]ComponentImpl{"c1": comp},
//		logger:       logger,
//	}
//	// when
//	err := reg.Stop()
//	// then
//	assert.NoError(t, err)
//	route.AssertExpectations(t)
//}
