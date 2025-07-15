package entities

import (
	"fmt"
	"log"
	"reflect"

	w "github.com/kijimaD/ruins/lib/world"

	ecs "github.com/x-hgg-x/goecs/v2"
)

// ComponentList is a list of preloaded entities with components
type ComponentList struct {
	Game []interface{}
}

// AddEntities adds entities with engine and game components
func AddEntities(world w.World, entityComponentList ComponentList) []ecs.Entity {
	// Create new entities and add engine components
	entities := make([]ecs.Entity, len(entityComponentList.Game))
	for iEntity := range entityComponentList.Game {
		entities[iEntity] = world.Manager.NewEntity()
		AddEntityComponents(entities[iEntity], world.Components.Game, entityComponentList.Game[iEntity])
	}

	// Add game components
	if entityComponentList.Game != nil {
		if len(entityComponentList.Game) != len(entities) {
			log.Fatal("incorrect size for game component list")
		}
		for iEntity := range entities {
			AddEntityComponents(entities[iEntity], world.Components.Game, entityComponentList.Game[iEntity])
		}
	}
	return entities
}

// AddEntityComponents adds loaded components to an entity
func AddEntityComponents(entity ecs.Entity, ecsComponentList interface{}, components interface{}) ecs.Entity {
	// 追加先のコンポーネントリスト。コンポーネントのスライス群
	ecv := reflect.ValueOf(ecsComponentList).Elem()
	// 追加するコンポーネント
	cv := reflect.ValueOf(components)
	for iField := 0; iField < cv.NumField(); iField++ {
		if !cv.Field(iField).IsNil() {
			component := cv.Field(iField).Elem()
			value := reflect.New(reflect.TypeOf(component.Interface()))

			switch component.Kind() {
			case reflect.Struct:
				// 追加対象コンポーネントの型名を使って、追加先コンポーネントのフィールドを対応付けて値を設定する
				value.Elem().Set(component)
				ecsComponent := ecv.FieldByName(component.Type().Name()).Interface().(ecs.DataComponent)
				entity.AddComponent(ecsComponent, value.Interface())
			case reflect.Interface:
				// Stringer インターフェースだけ対応している。Componentsに対応するフィールド名が必須なため
				if component.Type().Implements(reflect.TypeOf((*fmt.Stringer)(nil)).Elem()) {
					method := component.MethodByName("String")
					if !method.IsValid() {
						log.Fatal("String() に失敗した")
					}
					results := method.Call(nil)
					if len(results) != 1 {
						log.Fatal("String() の返り値の取得に失敗した")
					}
					v := component.Elem().Interface()
					value.Elem().Set(reflect.ValueOf(v))

					result := results[0].Interface().(string)
					ecsComponent := ecv.FieldByName(result).Interface().(ecs.DataComponent)
					entity.AddComponent(ecsComponent, value.Interface())
				}
			default:
				log.Fatalf("GameComponentListフィールドに指定された型の処理は定義されていない: %s", component.Kind())
			}
		}
	}
	return entity
}
