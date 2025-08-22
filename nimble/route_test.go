package nimble

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRoute_NewRoute_CreatesRouteWithDefaults(t *testing.T) {
	// given
	name := "route1"
	// when
	route := NewRoute(name)
	// then
	assert.NotNil(t, route)
	assert.Equal(t, name, route.Name)
	assert.Equal(t, 1, route.Concurrency)
	assert.NotNil(t, route.mainBlock)
	assert.NotNil(t, route.currentBlock)
	assert.Empty(t, route.dependencies)
	assert.NotNil(t, route.components)
}

func TestRoute_From_SetsFromAndDependencies(t *testing.T) {
	// given
	uri := "component:foo"
	// when
	route := From(uri)
	// then
	assert.Equal(t, uri, route.From)
	assert.Contains(t, route.dependencies, uri)
}

func TestRoute_To_AddsDependencyAndBlock(t *testing.T) {
	// given
	route := NewRoute("r")
	route.Registry = &mockRegistry{}
	uri := "component:bar"
	// when
	route.To(uri)
	// then
	assert.Contains(t, route.dependencies, uri)
}

func TestRoute_End_SetsCurrentBlockToParent(t *testing.T) {
	// given
	parent := &LinearBlock{}
	child := &LinearBlock{Parent: parent}
	route := NewRoute("r")
	route.currentBlock = child
	// when
	route.End()
	// then
	assert.Equal(t, parent, route.currentBlock)
}

func TestRoute_Start_SetsContextAndStartsGoroutines(t *testing.T) {
	// given
	route := NewRoute("r")
	route.Concurrency = 2
	route.components["r"] = &mockComponentImpl{}
	// when
	route.Start()
	// then
	assert.NotNil(t, route.Ctx)
	assert.NotNil(t, route.cancel)
	route.Stop()
}

func TestRoute_Stop_CancelsContext(t *testing.T) {
	// given
	route := NewRoute("r")
	route.Ctx, route.cancel = context.WithCancel(context.Background())
	// when
	route.Stop()
	// then
	assert.NotNil(t, route.Ctx.Err())
}
