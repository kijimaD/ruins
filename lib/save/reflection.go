package save

import (
	"fmt"
	"reflect"
	"sync"

	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// ComponentTypeInfo はコンポーネント型の情報
type ComponentTypeInfo struct {
	Name           string                                                         // 型名
	Type           reflect.Type                                                   // 型情報
	ComponentRef   interface{}                                                    // ECSコンポーネントへの参照
	ExtractFunc    func(w.World, ecs.Entity) (interface{}, bool)                  // エンティティからコンポーネントを抽出
	RestoreFunc    func(w.World, ecs.Entity, interface{}) error                   // エンティティにコンポーネントを復元
	ResolveRefFunc func(w.World, ecs.Entity, interface{}, *StableIDManager) error // エンティティ参照を解決
}

// ComponentRegistry はコンポーネント型の自動検出と管理を行う
type ComponentRegistry struct {
	types       map[reflect.Type]*ComponentTypeInfo
	nameToType  map[string]reflect.Type
	mutex       sync.RWMutex
	initialized bool
}

// NewComponentRegistry は新しいComponentRegistryを作成
func NewComponentRegistry() *ComponentRegistry {
	return &ComponentRegistry{
		types:      make(map[reflect.Type]*ComponentTypeInfo),
		nameToType: make(map[string]reflect.Type),
	}
}

// InitializeFromWorld はワールドから自動的にコンポーネント型を検出・登録
func (r *ComponentRegistry) InitializeFromWorld(world w.World) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.initialized {
		return nil // 既に初期化済み
	}

	components := world.Components

	// リフレクションを使って全コンポーネント型を自動登録
	r.registerComponent(reflect.TypeOf(&gc.Position{}), components.Position, r.extractPosition, r.restorePosition, nil)
	r.registerComponent(reflect.TypeOf(&gc.Velocity{}), components.Velocity, r.extractVelocity, r.restoreVelocity, nil)
	r.registerComponent(reflect.TypeOf(&gc.AIVision{}), components.AIVision, r.extractAIVision, r.restoreAIVision, r.resolveAIVisionRefs)
	r.registerComponent(reflect.TypeOf(&gc.AIRoaming{}), components.AIRoaming, r.extractAIRoaming, r.restoreAIRoaming, nil)
	r.registerComponent(reflect.TypeOf(&gc.AIChasing{}), components.AIChasing, r.extractAIChasing, r.restoreAIChasing, nil)
	r.registerComponent(reflect.TypeOf(&gc.Camera{}), components.Camera, r.extractCamera, r.restoreCamera, nil)
	r.registerComponent(reflect.TypeOf(&gc.SpriteRender{}), components.SpriteRender, r.extractSpriteRender, r.restoreSpriteRender, nil)
	r.registerComponent(reflect.TypeOf(&gc.GridElement{}), components.GridElement, r.extractGridElement, r.restoreGridElement, nil)

	// NullComponentは特別扱い
	r.registerNullComponent(reflect.TypeOf(&gc.Operator{}), components.Operator)
	r.registerNullComponent(reflect.TypeOf(&gc.BlockView{}), components.BlockView)
	r.registerNullComponent(reflect.TypeOf(&gc.BlockPass{}), components.BlockPass)

	r.initialized = true
	return nil
}

// registerComponent は単一コンポーネント型を登録
func (r *ComponentRegistry) registerComponent(
	typ reflect.Type,
	componentRef interface{},
	extractFunc func(w.World, ecs.Entity) (interface{}, bool),
	restoreFunc func(w.World, ecs.Entity, interface{}) error,
	resolveRefFunc func(w.World, ecs.Entity, interface{}, *StableIDManager) error,
) {
	// ポインタ型から要素型を取得
	elemType := typ.Elem()

	info := &ComponentTypeInfo{
		Name:           elemType.Name(),
		Type:           elemType,
		ComponentRef:   componentRef,
		ExtractFunc:    extractFunc,
		RestoreFunc:    restoreFunc,
		ResolveRefFunc: resolveRefFunc,
	}

	r.types[elemType] = info
	r.nameToType[elemType.Name()] = elemType
}

// registerNullComponent はNullComponent型を登録
func (r *ComponentRegistry) registerNullComponent(typ reflect.Type, componentRef interface{}) {
	elemType := typ.Elem()

	info := &ComponentTypeInfo{
		Name:         elemType.Name(),
		Type:         elemType,
		ComponentRef: componentRef,
		ExtractFunc: func(world w.World, entity ecs.Entity) (interface{}, bool) {
			// NullComponentの存在チェック
			switch elemType.Name() {
			case "Operator":
				return struct{}{}, entity.HasComponent(world.Components.Operator)
			case "BlockView":
				return struct{}{}, entity.HasComponent(world.Components.BlockView)
			case "BlockPass":
				return struct{}{}, entity.HasComponent(world.Components.BlockPass)
			}
			return nil, false
		},
		RestoreFunc: func(world w.World, entity ecs.Entity, data interface{}) error {
			// NullComponentを追加
			switch elemType.Name() {
			case "Operator":
				entity.AddComponent(world.Components.Operator, &gc.Operator{})
			case "BlockView":
				entity.AddComponent(world.Components.BlockView, &gc.BlockView{})
			case "BlockPass":
				entity.AddComponent(world.Components.BlockPass, &gc.BlockPass{})
			}
			return nil
		},
		ResolveRefFunc: nil, // NullComponentには参照解決は不要
	}

	r.types[elemType] = info
	r.nameToType[elemType.Name()] = elemType
}

// GetTypeInfo は型情報を取得
func (r *ComponentRegistry) GetTypeInfo(typ reflect.Type) (*ComponentTypeInfo, bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	info, exists := r.types[typ]
	return info, exists
}

// GetTypeInfoByName は名前から型情報を取得
func (r *ComponentRegistry) GetTypeInfoByName(name string) (*ComponentTypeInfo, bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	typ, exists := r.nameToType[name]
	if !exists {
		return nil, false
	}

	info, exists := r.types[typ]
	return info, exists
}

// GetAllTypes は登録されている全ての型を取得
func (r *ComponentRegistry) GetAllTypes() []*ComponentTypeInfo {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	types := make([]*ComponentTypeInfo, 0, len(r.types))
	for _, info := range r.types {
		types = append(types, info)
	}

	return types
}

// 各コンポーネント型の抽出・復元関数
func (r *ComponentRegistry) extractPosition(world w.World, entity ecs.Entity) (interface{}, bool) {
	if !entity.HasComponent(world.Components.Position) {
		return nil, false
	}
	pos := world.Components.Position.Get(entity).(*gc.Position)
	return *pos, true
}

func (r *ComponentRegistry) restorePosition(world w.World, entity ecs.Entity, data interface{}) error {
	pos, ok := data.(gc.Position)
	if !ok {
		return fmt.Errorf("invalid Position data type: %T", data)
	}
	entity.AddComponent(world.Components.Position, &pos)
	return nil
}

func (r *ComponentRegistry) extractVelocity(world w.World, entity ecs.Entity) (interface{}, bool) {
	if !entity.HasComponent(world.Components.Velocity) {
		return nil, false
	}
	vel := world.Components.Velocity.Get(entity).(*gc.Velocity)
	return *vel, true
}

func (r *ComponentRegistry) restoreVelocity(world w.World, entity ecs.Entity, data interface{}) error {
	vel, ok := data.(gc.Velocity)
	if !ok {
		return fmt.Errorf("invalid Velocity data type: %T", data)
	}
	entity.AddComponent(world.Components.Velocity, &vel)
	return nil
}

func (r *ComponentRegistry) extractAIVision(world w.World, entity ecs.Entity) (interface{}, bool) {
	if !entity.HasComponent(world.Components.AIVision) {
		return nil, false
	}
	vision := world.Components.AIVision.Get(entity).(*gc.AIVision)

	// 注意: エンティティ参照のStableID変換は上位レベルで行う
	return *vision, true
}

func (r *ComponentRegistry) restoreAIVision(world w.World, entity ecs.Entity, data interface{}) error {
	vision, ok := data.(gc.AIVision)
	if !ok {
		return fmt.Errorf("invalid AIVision data type: %T", data)
	}
	// エンティティ参照の復元は後で行うため、一旦nilで設定
	vision.TargetEntity = nil
	entity.AddComponent(world.Components.AIVision, &vision)
	return nil
}

func (r *ComponentRegistry) resolveAIVisionRefs(world w.World, entity ecs.Entity, data interface{}, idManager *StableIDManager) error {
	// エンティティ参照の解決はSerializationManagerで実装
	return nil
}

func (r *ComponentRegistry) extractAIRoaming(world w.World, entity ecs.Entity) (interface{}, bool) {
	if !entity.HasComponent(world.Components.AIRoaming) {
		return nil, false
	}
	roaming := world.Components.AIRoaming.Get(entity).(*gc.AIRoaming)
	return *roaming, true
}

func (r *ComponentRegistry) restoreAIRoaming(world w.World, entity ecs.Entity, data interface{}) error {
	roaming, ok := data.(gc.AIRoaming)
	if !ok {
		return fmt.Errorf("invalid AIRoaming data type: %T", data)
	}
	entity.AddComponent(world.Components.AIRoaming, &roaming)
	return nil
}

func (r *ComponentRegistry) extractAIChasing(world w.World, entity ecs.Entity) (interface{}, bool) {
	if !entity.HasComponent(world.Components.AIChasing) {
		return nil, false
	}
	chasing := world.Components.AIChasing.Get(entity).(*gc.AIChasing)
	return *chasing, true
}

func (r *ComponentRegistry) restoreAIChasing(world w.World, entity ecs.Entity, data interface{}) error {
	chasing, ok := data.(gc.AIChasing)
	if !ok {
		return fmt.Errorf("invalid AIChasing data type: %T", data)
	}
	entity.AddComponent(world.Components.AIChasing, &chasing)
	return nil
}

func (r *ComponentRegistry) extractCamera(world w.World, entity ecs.Entity) (interface{}, bool) {
	if !entity.HasComponent(world.Components.Camera) {
		return nil, false
	}
	camera := world.Components.Camera.Get(entity).(*gc.Camera)
	return *camera, true
}

func (r *ComponentRegistry) restoreCamera(world w.World, entity ecs.Entity, data interface{}) error {
	camera, ok := data.(gc.Camera)
	if !ok {
		return fmt.Errorf("invalid Camera data type: %T", data)
	}
	entity.AddComponent(world.Components.Camera, &camera)
	return nil
}

func (r *ComponentRegistry) extractSpriteRender(world w.World, entity ecs.Entity) (interface{}, bool) {
	if !entity.HasComponent(world.Components.SpriteRender) {
		return nil, false
	}
	sprite := world.Components.SpriteRender.Get(entity).(*gc.SpriteRender)
	return *sprite, true
}

func (r *ComponentRegistry) restoreSpriteRender(world w.World, entity ecs.Entity, data interface{}) error {
	sprite, ok := data.(gc.SpriteRender)
	if !ok {
		return fmt.Errorf("invalid SpriteRender data type: %T", data)
	}
	entity.AddComponent(world.Components.SpriteRender, &sprite)
	return nil
}

func (r *ComponentRegistry) extractGridElement(world w.World, entity ecs.Entity) (interface{}, bool) {
	if !entity.HasComponent(world.Components.GridElement) {
		return nil, false
	}
	grid := world.Components.GridElement.Get(entity).(*gc.GridElement)
	return *grid, true
}

func (r *ComponentRegistry) restoreGridElement(world w.World, entity ecs.Entity, data interface{}) error {
	grid, ok := data.(gc.GridElement)
	if !ok {
		return fmt.Errorf("invalid GridElement data type: %T", data)
	}
	entity.AddComponent(world.Components.GridElement, &grid)
	return nil
}
