package nimble

import (
	"context"
	"errors"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Registry interface {
	Logger() *zap.Logger
	Instantiate(uri string) (ComponentImpl, error)
	RegisterRoute(route *Route) error
	Start() error
	Stop() error
}

type RegistryImpl struct {
	ComponentMap map[string]Component
	InstancesMap map[string]ComponentImpl
	RoutesMap    map[string]*Route
	logger       *zap.Logger
}

func NewRegistry(components []Component, routes []*Route, logger *zap.Logger) (Registry, error) {
	compMap := make(map[string]Component)
	for _, component := range components {
		compMap[component.Key()] = component
	}

	reg := &RegistryImpl{
		ComponentMap: compMap,
		InstancesMap: make(map[string]ComponentImpl),
		RoutesMap:    make(map[string]*Route),
		logger:       logger,
	}

	for _, route := range routes {
		err := reg.RegisterRoute(route)
		if err != nil {
			return nil, err
		}
	}

	return reg, nil
}

func (m *RegistryImpl) Logger() *zap.Logger {
	return m.logger
}

func (m *RegistryImpl) Instantiate(uri string) (ComponentImpl, error) {
	params, err := ParseParams(uri)
	if err != nil {
		return nil, errors.New("Invalid URI: " + uri)
	}

	id := params.GetId()
	m.logger.Info("Instantiate", zap.String("id", id), zap.String("uri", uri))

	instance, ok := m.InstancesMap[id]
	if ok {
		return instance, nil
	}

	comp, ok := m.ComponentMap[params.Key]
	if !ok {
		return nil, errors.New("Cannot find Component: " + params.Key)
	}

	inst, err := comp.Instantiate(m, params)
	if err != nil {
		return nil, errors.New("Cannot Instantiate " + params.Key + ": " + err.Error())
	}

	m.InstancesMap[id] = inst

	return inst, nil

}

func (m *RegistryImpl) RegisterRoute(route *Route) error {
	route.Registry = m
	m.RoutesMap[route.id] = route

	for _, uri := range route.dependencies {
		component, err := m.Instantiate(uri)
		if err != nil {
			return err
		}
		route.components[uri] = component
	}

	return nil
}

func (m *RegistryImpl) Start() error {
	m.logger.Info("Register Starting...")

	for _, route := range m.RoutesMap {
		m.logger.Info("Starting Route", zap.String("id", route.id), zap.String("name", route.Name))
		route.Start()
	}

	for uri, component := range m.InstancesMap {
		m.logger.Info("Starting Component...", zap.String("id", uri))
		component.Start()
	}

	m.logger.Info("Register Started")
	return nil
}

func (m *RegistryImpl) Stop() error {
	m.logger.Info("Register Stopping...")

	for uri, component := range m.InstancesMap {
		m.logger.Info("Stopping Component...", zap.String("id", uri))
		component.Stop()
	}

	for _, route := range m.RoutesMap {
		m.logger.Info("Stopping Route...", zap.String("id", route.id), zap.String("name", route.Name))
		route.Stop()
	}

	m.logger.Info("Register Stopped")

	return nil
}

var Module = fx.Module(
	"nimble",
	fx.Provide(
		fx.Annotate(
			NewRegistry,
			fx.ParamTags(`group:"nimble-components"`, `group:"nimble-routes"`),
		),
	),
	fx.Invoke(func(lc fx.Lifecycle, reg Registry) error {
		lc.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				return reg.Stop()
			},
		})
		return reg.Start()
	}),
)
