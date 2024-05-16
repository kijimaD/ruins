package components

import (
	"reflect"

	ecs "github.com/x-hgg-x/goecs/v2"
)

// Components contains engine and game components
type Components struct {
	Game interface{}
}

// InitComponents initializes components
func InitComponents(manager *ecs.Manager, gameComponents interface{}) *Components {
	components := &Components{Game: gameComponents}
	initFields(manager, components.Game)
	return components
}

func initFields(manager *ecs.Manager, components interface{}) {
	if components != nil {
		v := reflect.ValueOf(components).Elem()
		for iField := 0; iField < v.NumField(); iField++ {
			component := v.Field(iField)
			switch component.Interface().(type) {
			case *ecs.NullComponent:
				component.Set(reflect.ValueOf(manager.NewNullComponent()))
			case *ecs.SliceComponent:
				component.Set(reflect.ValueOf(manager.NewSliceComponent()))
			case *ecs.DenseSliceComponent:
				component.Set(reflect.ValueOf(manager.NewDenseSliceComponent()))
			case *ecs.MapComponent:
				component.Set(reflect.ValueOf(manager.NewMapComponent()))
			}
		}
	}
}
