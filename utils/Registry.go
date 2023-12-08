package utils

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"sync"
)

type PluginType string

const (
	WorkflowPlugin PluginType = "workflow"
)

type Plugin interface {
	Configure() error
}

var registry map[PluginType]map[string]Plugin
var mutex sync.Mutex

func init() {
	registry = make(map[PluginType]map[string]Plugin)
}

func RegisterPlugin(pluginType PluginType, name string, v Plugin) {
	mutex.Lock()
	defer mutex.Unlock()

	log.Debugf("[RegisterWorkflowPlugin] type: %v, name: %v", pluginType, name)

	_, ok := registry[pluginType]
	if !ok {
		registry[pluginType] = make(map[string]Plugin)
	}

	_, ok = registry[pluginType][name]
	if ok {
		panic(fmt.Sprintf("plugin already exists, type: %v, name: %v", pluginType, name))
	}
	registry[pluginType][name] = v
}

func GetPlugin(pluginType PluginType, name string) (Plugin, error) {
	mutex.Lock()
	defer mutex.Unlock()

	if registry == nil {
		return nil, errors.New("empty workflow registry")
	}

	plugins, ok := registry[pluginType]
	if !ok {
		return nil, errors.New(fmt.Sprintf("empty plugin type: %v, name: %v", pluginType, name))
	}
	p, ok := plugins[name]
	if !ok {
		return nil, errors.New(fmt.Sprintf("empty plugin, type: %v, name: %v", pluginType, name))
	}
	return p, nil
}
