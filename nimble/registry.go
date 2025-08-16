package nimble

import (
	"context"
	"errors"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"log"
)

type Registry struct {
	ComponentMap map[string]Component
	InstancesMap map[string]ComponentImpl
	RoutesMap    map[string]*Route
	logger       *zap.Logger
}

func NewRegistry(components []Component, routes []*Route, logger *zap.Logger) (*Registry, error) {
	compMap := make(map[string]Component)
	for _, component := range components {
		compMap[component.Key()] = component
	}

	reg := &Registry{
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

func (m *Registry) Instantiate(uri string) (ComponentImpl, error) {
	instance, ok := m.InstancesMap[uri]
	if ok {
		return instance, nil
	}

	params, err := ParseParams(uri)
	if err != nil {
		return nil, errors.New("Invalid URI: " + uri)
	}

	comp, ok := m.ComponentMap[params.Key]
	if !ok {
		return nil, errors.New("Cannot find Component: " + params.Key)
	}

	inst, err := comp.Instantiate(params)
	if err != nil {
		return nil, errors.New("Cannot Instantiate " + params.Key + ": " + err.Error())
	}

	m.InstancesMap[uri] = inst

	return inst, nil

}

func (m *Registry) RegisterRoute(route *Route) error {
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

func (m *Registry) Start() error {
	m.logger.Info("Register Starting")

	for _, route := range m.RoutesMap {
		route.Start()
		m.logger.Info("Started Route", zap.String("id", route.id), zap.String("name", route.Name))
	}

	for uri, component := range m.InstancesMap {
		go component.Start()
		m.logger.Info("Started Component", zap.String("id", uri))

	}

	m.logger.Info("Register Started")
	return nil
}

func (m *Registry) Stop() error {
	log.Println("Register Stop")
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
	fx.Invoke(func(lc fx.Lifecycle, reg *Registry) error {
		lc.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				return reg.Stop()
			},
		})
		return reg.Start()
	}),
)
