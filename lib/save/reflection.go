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
	r.registerNullComponent(reflect.TypeOf(&gc.Player{}), components.Player)
	r.registerNullComponent(reflect.TypeOf(&gc.FactionAllyData{}), components.FactionAlly)
	r.registerNullComponent(reflect.TypeOf(&gc.FactionEnemyData{}), components.FactionEnemy)
	r.registerNullComponent(reflect.TypeOf(&gc.InParty{}), components.InParty)
	r.registerNullComponent(reflect.TypeOf(&gc.Item{}), components.Item)

	// アイテム位置情報コンポーネント
	r.registerNullComponent(reflect.TypeOf(&gc.LocationInBackpack{}), components.ItemLocationInBackpack)
	r.registerNullComponent(reflect.TypeOf(&gc.LocationOnField{}), components.ItemLocationOnField)
	r.registerNullComponent(reflect.TypeOf(&gc.LocationNone{}), components.ItemLocationNone)
	r.registerComponent(reflect.TypeOf(&gc.LocationEquipped{}), components.ItemLocationEquipped, r.extractItemLocationEquipped, r.restoreItemLocationEquipped, r.resolveLocationEquippedRefs)

	// 装備変更フラグ
	r.registerNullComponent(reflect.TypeOf(&gc.EquipmentChanged{}), components.EquipmentChanged)

	// データコンポーネント
	r.registerComponent(reflect.TypeOf(&gc.Name{}), components.Name, r.extractName, r.restoreName, nil)
	r.registerComponent(reflect.TypeOf(&gc.Pools{}), components.Pools, r.extractPools, r.restorePools, nil)
	r.registerComponent(reflect.TypeOf(&gc.Attributes{}), components.Attributes, r.extractAttributes, r.restoreAttributes, nil)
	r.registerComponent(reflect.TypeOf(&gc.Description{}), components.Description, r.extractDescription, r.restoreDescription, nil)
	r.registerComponent(reflect.TypeOf(&gc.Job{}), components.Job, r.extractJob, r.restoreJob, nil)

	// アイテム関連コンポーネント
	r.registerComponent(reflect.TypeOf(&gc.Wearable{}), components.Wearable, r.extractWearable, r.restoreWearable, nil)
	r.registerComponent(reflect.TypeOf(&gc.Card{}), components.Card, r.extractCard, r.restoreCard, nil)
	r.registerComponent(reflect.TypeOf(&gc.Material{}), components.Material, r.extractMaterial, r.restoreMaterial, nil)
	r.registerComponent(reflect.TypeOf(&gc.Consumable{}), components.Consumable, r.extractConsumable, r.restoreConsumable, nil)
	r.registerComponent(reflect.TypeOf(&gc.Attack{}), components.Attack, r.extractAttack, r.restoreAttack, nil)
	r.registerComponent(reflect.TypeOf(&gc.Recipe{}), components.Recipe, r.extractRecipe, r.restoreRecipe, nil)
	r.registerComponent(reflect.TypeOf(&gc.ProvidesHealing{}), components.ProvidesHealing, r.extractProvidesHealing, r.restoreProvidesHealing, nil)
	r.registerComponent(reflect.TypeOf(&gc.InflictsDamage{}), components.InflictsDamage, r.extractInflictsDamage, r.restoreInflictsDamage, nil)

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
			case "Player":
				return struct{}{}, entity.HasComponent(world.Components.Player)
			case "FactionAllyData":
				return struct{}{}, entity.HasComponent(world.Components.FactionAlly)
			case "FactionEnemyData":
				return struct{}{}, entity.HasComponent(world.Components.FactionEnemy)
			case "InParty":
				return struct{}{}, entity.HasComponent(world.Components.InParty)
			case "Item":
				return struct{}{}, entity.HasComponent(world.Components.Item)
			case "LocationInBackpack":
				return struct{}{}, entity.HasComponent(world.Components.ItemLocationInBackpack)
			case "LocationOnField":
				return struct{}{}, entity.HasComponent(world.Components.ItemLocationOnField)
			case "LocationNone":
				return struct{}{}, entity.HasComponent(world.Components.ItemLocationNone)
			case "EquipmentChanged":
				return struct{}{}, entity.HasComponent(world.Components.EquipmentChanged)
			}
			return nil, false
		},
		RestoreFunc: func(world w.World, entity ecs.Entity, _ interface{}) error {
			// NullComponentを追加
			switch elemType.Name() {
			case "Operator":
				entity.AddComponent(world.Components.Operator, &gc.Operator{})
			case "BlockView":
				entity.AddComponent(world.Components.BlockView, &gc.BlockView{})
			case "BlockPass":
				entity.AddComponent(world.Components.BlockPass, &gc.BlockPass{})
			case "Player":
				entity.AddComponent(world.Components.Player, &gc.Player{})
			case "FactionAllyData":
				entity.AddComponent(world.Components.FactionAlly, &gc.FactionAllyData{})
			case "FactionEnemyData":
				entity.AddComponent(world.Components.FactionEnemy, &gc.FactionEnemyData{})
			case "InParty":
				entity.AddComponent(world.Components.InParty, &gc.InParty{})
			case "Item":
				entity.AddComponent(world.Components.Item, &gc.Item{})
			case "LocationInBackpack":
				entity.AddComponent(world.Components.ItemLocationInBackpack, &gc.LocationInBackpack{})
			case "LocationOnField":
				entity.AddComponent(world.Components.ItemLocationOnField, &gc.LocationOnField{})
			case "LocationNone":
				entity.AddComponent(world.Components.ItemLocationNone, &gc.LocationNone{})
			case "EquipmentChanged":
				entity.AddComponent(world.Components.EquipmentChanged, &gc.EquipmentChanged{})
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

func (r *ComponentRegistry) resolveAIVisionRefs(_ w.World, _ ecs.Entity, _ interface{}, _ *StableIDManager) error {
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

// Name コンポーネントの処理
func (r *ComponentRegistry) extractName(world w.World, entity ecs.Entity) (interface{}, bool) {
	if !entity.HasComponent(world.Components.Name) {
		return nil, false
	}
	name := world.Components.Name.Get(entity).(*gc.Name)
	return *name, true
}

func (r *ComponentRegistry) restoreName(world w.World, entity ecs.Entity, data interface{}) error {
	name, ok := data.(gc.Name)
	if !ok {
		return fmt.Errorf("invalid Name data type: %T", data)
	}
	entity.AddComponent(world.Components.Name, &name)
	return nil
}

// Pools コンポーネントの処理
func (r *ComponentRegistry) extractPools(world w.World, entity ecs.Entity) (interface{}, bool) {
	if !entity.HasComponent(world.Components.Pools) {
		return nil, false
	}
	pools := world.Components.Pools.Get(entity).(*gc.Pools)
	return *pools, true
}

func (r *ComponentRegistry) restorePools(world w.World, entity ecs.Entity, data interface{}) error {
	pools, ok := data.(gc.Pools)
	if !ok {
		return fmt.Errorf("invalid Pools data type: %T", data)
	}
	entity.AddComponent(world.Components.Pools, &pools)
	return nil
}

// Attributes コンポーネントの処理
func (r *ComponentRegistry) extractAttributes(world w.World, entity ecs.Entity) (interface{}, bool) {
	if !entity.HasComponent(world.Components.Attributes) {
		return nil, false
	}
	attributes := world.Components.Attributes.Get(entity).(*gc.Attributes)
	return *attributes, true
}

func (r *ComponentRegistry) restoreAttributes(world w.World, entity ecs.Entity, data interface{}) error {
	attributes, ok := data.(gc.Attributes)
	if !ok {
		return fmt.Errorf("invalid Attributes data type: %T", data)
	}
	entity.AddComponent(world.Components.Attributes, &attributes)
	return nil
}

// ItemLocationEquipped コンポーネントの処理
func (r *ComponentRegistry) extractItemLocationEquipped(world w.World, entity ecs.Entity) (interface{}, bool) {
	if !entity.HasComponent(world.Components.ItemLocationEquipped) {
		return nil, false
	}
	equipped := world.Components.ItemLocationEquipped.Get(entity).(*gc.LocationEquipped)
	return *equipped, true
}

func (r *ComponentRegistry) restoreItemLocationEquipped(world w.World, entity ecs.Entity, data interface{}) error {
	equipped, ok := data.(gc.LocationEquipped)
	if !ok {
		return fmt.Errorf("invalid LocationEquipped data type: %T", data)
	}
	entity.AddComponent(world.Components.ItemLocationEquipped, &equipped)
	return nil
}

func (r *ComponentRegistry) resolveLocationEquippedRefs(_ w.World, _ ecs.Entity, _ interface{}, _ *StableIDManager) error {
	// エンティティ参照の解決はSerializationManagerで実装
	return nil
}

// Description コンポーネントの処理
func (r *ComponentRegistry) extractDescription(world w.World, entity ecs.Entity) (interface{}, bool) {
	if !entity.HasComponent(world.Components.Description) {
		return nil, false
	}
	desc := world.Components.Description.Get(entity).(*gc.Description)
	return *desc, true
}

func (r *ComponentRegistry) restoreDescription(world w.World, entity ecs.Entity, data interface{}) error {
	desc, ok := data.(gc.Description)
	if !ok {
		return fmt.Errorf("invalid Description data type: %T", data)
	}
	entity.AddComponent(world.Components.Description, &desc)
	return nil
}

// Job コンポーネントの処理
func (r *ComponentRegistry) extractJob(world w.World, entity ecs.Entity) (interface{}, bool) {
	if !entity.HasComponent(world.Components.Job) {
		return nil, false
	}
	job := world.Components.Job.Get(entity).(*gc.Job)
	return *job, true
}

func (r *ComponentRegistry) restoreJob(world w.World, entity ecs.Entity, data interface{}) error {
	job, ok := data.(gc.Job)
	if !ok {
		return fmt.Errorf("invalid Job data type: %T", data)
	}
	entity.AddComponent(world.Components.Job, &job)
	return nil
}

// Wearable コンポーネントの処理
func (r *ComponentRegistry) extractWearable(world w.World, entity ecs.Entity) (interface{}, bool) {
	if !entity.HasComponent(world.Components.Wearable) {
		return nil, false
	}
	wearable := world.Components.Wearable.Get(entity).(*gc.Wearable)
	return *wearable, true
}

func (r *ComponentRegistry) restoreWearable(world w.World, entity ecs.Entity, data interface{}) error {
	wearable, ok := data.(gc.Wearable)
	if !ok {
		return fmt.Errorf("invalid Wearable data type: %T", data)
	}
	entity.AddComponent(world.Components.Wearable, &wearable)
	return nil
}

// Card コンポーネントの処理
func (r *ComponentRegistry) extractCard(world w.World, entity ecs.Entity) (interface{}, bool) {
	if !entity.HasComponent(world.Components.Card) {
		return nil, false
	}
	card := world.Components.Card.Get(entity).(*gc.Card)
	return *card, true
}

func (r *ComponentRegistry) restoreCard(world w.World, entity ecs.Entity, data interface{}) error {
	card, ok := data.(gc.Card)
	if !ok {
		return fmt.Errorf("invalid Card data type: %T", data)
	}
	entity.AddComponent(world.Components.Card, &card)
	return nil
}

// Material コンポーネントの処理
func (r *ComponentRegistry) extractMaterial(world w.World, entity ecs.Entity) (interface{}, bool) {
	if !entity.HasComponent(world.Components.Material) {
		return nil, false
	}
	material := world.Components.Material.Get(entity).(*gc.Material)
	return *material, true
}

func (r *ComponentRegistry) restoreMaterial(world w.World, entity ecs.Entity, data interface{}) error {
	material, ok := data.(gc.Material)
	if !ok {
		return fmt.Errorf("invalid Material data type: %T", data)
	}
	entity.AddComponent(world.Components.Material, &material)
	return nil
}

// Consumable コンポーネントの処理
func (r *ComponentRegistry) extractConsumable(world w.World, entity ecs.Entity) (interface{}, bool) {
	if !entity.HasComponent(world.Components.Consumable) {
		return nil, false
	}
	consumable := world.Components.Consumable.Get(entity).(*gc.Consumable)
	return *consumable, true
}

func (r *ComponentRegistry) restoreConsumable(world w.World, entity ecs.Entity, data interface{}) error {
	consumable, ok := data.(gc.Consumable)
	if !ok {
		return fmt.Errorf("invalid Consumable data type: %T", data)
	}
	entity.AddComponent(world.Components.Consumable, &consumable)
	return nil
}

// Attack コンポーネントの処理
func (r *ComponentRegistry) extractAttack(world w.World, entity ecs.Entity) (interface{}, bool) {
	if !entity.HasComponent(world.Components.Attack) {
		return nil, false
	}
	attack := world.Components.Attack.Get(entity).(*gc.Attack)
	return *attack, true
}

func (r *ComponentRegistry) restoreAttack(world w.World, entity ecs.Entity, data interface{}) error {
	attack, ok := data.(gc.Attack)
	if !ok {
		return fmt.Errorf("invalid Attack data type: %T", data)
	}
	entity.AddComponent(world.Components.Attack, &attack)
	return nil
}

// Recipe コンポーネントの処理
func (r *ComponentRegistry) extractRecipe(world w.World, entity ecs.Entity) (interface{}, bool) {
	if !entity.HasComponent(world.Components.Recipe) {
		return nil, false
	}
	recipe := world.Components.Recipe.Get(entity).(*gc.Recipe)
	return *recipe, true
}

func (r *ComponentRegistry) restoreRecipe(world w.World, entity ecs.Entity, data interface{}) error {
	recipe, ok := data.(gc.Recipe)
	if !ok {
		return fmt.Errorf("invalid Recipe data type: %T", data)
	}
	entity.AddComponent(world.Components.Recipe, &recipe)
	return nil
}

// ProvidesHealing コンポーネントの処理
func (r *ComponentRegistry) extractProvidesHealing(world w.World, entity ecs.Entity) (interface{}, bool) {
	if !entity.HasComponent(world.Components.ProvidesHealing) {
		return nil, false
	}
	healing := world.Components.ProvidesHealing.Get(entity).(*gc.ProvidesHealing)

	// Amounterインターフェースを具体的な型に変換してシリアライズ
	var amountData map[string]interface{}
	switch a := healing.Amount.(type) {
	case gc.RatioAmount:
		amountData = map[string]interface{}{
			"type":  "ratio",
			"ratio": a.Ratio,
		}
	case gc.NumeralAmount:
		amountData = map[string]interface{}{
			"type":    "numeral",
			"numeral": a.Numeral,
		}
	default:
		// デフォルトでRatio 0.5を使用
		amountData = map[string]interface{}{
			"type":  "ratio",
			"ratio": 0.5,
		}
	}

	return map[string]interface{}{
		"amount": amountData,
	}, true
}

func (r *ComponentRegistry) restoreProvidesHealing(world w.World, entity ecs.Entity, data interface{}) error {
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid ProvidesHealing data type: %T", data)
	}

	amountData, ok := dataMap["amount"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid ProvidesHealing amount data")
	}

	var amount gc.Amounter
	amountType, _ := amountData["type"].(string)
	switch amountType {
	case "ratio":
		if ratio, ok := amountData["ratio"].(float64); ok {
			amount = gc.RatioAmount{Ratio: ratio}
		}
	case "numeral":
		if numeral, ok := amountData["numeral"].(float64); ok {
			amount = gc.NumeralAmount{Numeral: int(numeral)}
		}
	default:
		return fmt.Errorf("unknown amount type: %s", amountType)
	}

	healing := &gc.ProvidesHealing{
		Amount: amount,
	}
	entity.AddComponent(world.Components.ProvidesHealing, healing)
	return nil
}

// InflictsDamage コンポーネントの処理
func (r *ComponentRegistry) extractInflictsDamage(world w.World, entity ecs.Entity) (interface{}, bool) {
	if !entity.HasComponent(world.Components.InflictsDamage) {
		return nil, false
	}
	damage := world.Components.InflictsDamage.Get(entity).(*gc.InflictsDamage)
	return *damage, true
}

func (r *ComponentRegistry) restoreInflictsDamage(world w.World, entity ecs.Entity, data interface{}) error {
	damage, ok := data.(gc.InflictsDamage)
	if !ok {
		return fmt.Errorf("invalid InflictsDamage data type: %T", data)
	}
	entity.AddComponent(world.Components.InflictsDamage, &damage)
	return nil
}
