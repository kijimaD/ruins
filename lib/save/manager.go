package save

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"time"

	gc "github.com/kijimaD/ruins/lib/components"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// コンポーネント名の定数
const (
	ComponentAIVision          = "AIVision"
	ComponentOperator          = "Operator"
	ComponentBlockView         = "BlockView"
	ComponentBlockPass         = "BlockPass"
	ComponentFactionAllyData   = "FactionAllyData"
	ComponentFactionEnemyData  = "FactionEnemyData"
	ComponentInParty           = "InParty"
	ComponentItem              = "Item"
	ComponentLocationInBackpack = "LocationInBackpack"
	ComponentLocationEquipped  = "LocationEquipped"
	ComponentLocationOnField   = "LocationOnField"
	ComponentLocationNone      = "LocationNone"
	ComponentEquipmentChanged  = "EquipmentChanged"
)

// Data はセーブデータの最上位構造
type Data struct {
	Version   string        `json:"version"`
	Timestamp time.Time     `json:"timestamp"`
	World     WorldSaveData `json:"world"`
}

// WorldSaveData はワールド全体のセーブデータ
type WorldSaveData struct {
	Entities []EntitySaveData `json:"entities"`
	// TODO: リソース情報も追加予定
	// Resources ResourcesSaveData `json:"resources"`
}

// EntitySaveData は単一エンティティのセーブデータ
type EntitySaveData struct {
	StableID   StableID                 `json:"stable_id"`
	Components map[string]ComponentData `json:"components"`
}

// ComponentData はコンポーネントデータ
type ComponentData struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// EntityReference はエンティティ参照のセーブデータ
type EntityReference struct {
	TargetStableID *StableID `json:"target_stable_id,omitempty"`
}

// SerializationManager は安定ID + リフレクションベースのシリアライゼーションを管理
type SerializationManager struct {
	saveDirectory     string
	stableIDManager   *StableIDManager
	componentRegistry *ComponentRegistry
}

// NewSerializationManager は新しいSerializationManagerを作成
func NewSerializationManager(saveDir string) *SerializationManager {
	return &SerializationManager{
		saveDirectory:     saveDir,
		stableIDManager:   NewStableIDManager(),
		componentRegistry: NewComponentRegistry(),
	}
}

// SaveWorld はワールド全体をファイルに保存
func (sm *SerializationManager) SaveWorld(world w.World, slotName string) error {
	// コンポーネントレジストリを初期化
	err := sm.componentRegistry.InitializeFromWorld(world)
	if err != nil {
		return fmt.Errorf("failed to initialize component registry: %w", err)
	}

	// 保存ディレクトリを作成
	err = os.MkdirAll(sm.saveDirectory, 0755)
	if err != nil {
		return fmt.Errorf("failed to create save directory: %w", err)
	}

	// ワールドデータを抽出
	worldData := sm.extractWorldData(world)

	// セーブデータを作成
	saveData := Data{
		Version:   "1.0.0",
		Timestamp: time.Now(),
		World:     worldData,
	}

	// JSONにシリアライズ
	data, err := json.MarshalIndent(saveData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal save data: %w", err)
	}

	// ファイルに書き込み
	fileName := filepath.Join(sm.saveDirectory, slotName+".json")
	err = os.WriteFile(fileName, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write save file: %w", err)
	}

	return nil
}

// LoadWorld はファイルからワールドを復元
func (sm *SerializationManager) LoadWorld(world w.World, slotName string) error {
	// コンポーネントレジストリを初期化
	err := sm.componentRegistry.InitializeFromWorld(world)
	if err != nil {
		return fmt.Errorf("failed to initialize component registry: %w", err)
	}

	// ファイルを読み込み
	fileName := filepath.Join(sm.saveDirectory, slotName+".json")
	data, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("failed to read save file: %w", err)
	}

	// JSONをパース
	var saveData Data
	err = json.Unmarshal(data, &saveData)
	if err != nil {
		return fmt.Errorf("failed to unmarshal save data: %w", err)
	}

	// バージョンチェック
	if saveData.Version != "1.0.0" {
		return fmt.Errorf("unsupported save data version: %s", saveData.Version)
	}

	// ワールドをクリア
	sm.clearWorld(world)
	sm.stableIDManager.Clear()

	// ワールドデータを復元
	err = sm.restoreWorldData(world, saveData.World)
	if err != nil {
		return fmt.Errorf("failed to restore world data: %w", err)
	}

	return nil
}

// extractWorldData はワールドからセーブデータを抽出
//
//nolint:gocyclo // コンポーネント種別ごとの処理が必要なため複雑度が高い
func (sm *SerializationManager) extractWorldData(world w.World) WorldSaveData {
	entities := []EntitySaveData{}
	processedEntities := make(map[ecs.Entity]bool) // 重複処理防止

	// 各コンポーネント型を持つエンティティを検索
	for _, typeInfo := range sm.componentRegistry.GetAllTypes() {
		entityCount := 0

		// この型のコンポーネントを持つ全エンティティを取得
		switch typeInfo.Name {
		case "Position":
			world.Manager.Join(world.Components.Position).Visit(ecs.Visit(func(entity ecs.Entity) {
				entityCount++
				sm.processEntityForSave(entity, world, &entities, processedEntities)
			}))
		case "Velocity":
			world.Manager.Join(world.Components.Velocity).Visit(ecs.Visit(func(entity ecs.Entity) {
				entityCount++
				sm.processEntityForSave(entity, world, &entities, processedEntities)
			}))
		case ComponentAIVision:
			world.Manager.Join(world.Components.AIVision).Visit(ecs.Visit(func(entity ecs.Entity) {
				entityCount++
				sm.processEntityForSave(entity, world, &entities, processedEntities)
			}))
		case "AIRoaming":
			world.Manager.Join(world.Components.AIRoaming).Visit(ecs.Visit(func(entity ecs.Entity) {
				entityCount++
				sm.processEntityForSave(entity, world, &entities, processedEntities)
			}))
		case "AIChasing":
			world.Manager.Join(world.Components.AIChasing).Visit(ecs.Visit(func(entity ecs.Entity) {
				entityCount++
				sm.processEntityForSave(entity, world, &entities, processedEntities)
			}))
		case "Camera":
			world.Manager.Join(world.Components.Camera).Visit(ecs.Visit(func(entity ecs.Entity) {
				entityCount++
				sm.processEntityForSave(entity, world, &entities, processedEntities)
			}))
		case "SpriteRender":
			world.Manager.Join(world.Components.SpriteRender).Visit(ecs.Visit(func(entity ecs.Entity) {
				entityCount++
				sm.processEntityForSave(entity, world, &entities, processedEntities)
			}))
		case "GridElement":
			world.Manager.Join(world.Components.GridElement).Visit(ecs.Visit(func(entity ecs.Entity) {
				entityCount++
				sm.processEntityForSave(entity, world, &entities, processedEntities)
			}))
		case ComponentOperator:
			world.Manager.Join(world.Components.Operator).Visit(ecs.Visit(func(entity ecs.Entity) {
				entityCount++
				sm.processEntityForSave(entity, world, &entities, processedEntities)
			}))
		case ComponentBlockView:
			world.Manager.Join(world.Components.BlockView).Visit(ecs.Visit(func(entity ecs.Entity) {
				entityCount++
				sm.processEntityForSave(entity, world, &entities, processedEntities)
			}))
		case ComponentBlockPass:
			world.Manager.Join(world.Components.BlockPass).Visit(ecs.Visit(func(entity ecs.Entity) {
				entityCount++
				sm.processEntityForSave(entity, world, &entities, processedEntities)
			}))
		case ComponentFactionAllyData:
			world.Manager.Join(world.Components.FactionAlly).Visit(ecs.Visit(func(entity ecs.Entity) {
				entityCount++
				sm.processEntityForSave(entity, world, &entities, processedEntities)
			}))
		case ComponentFactionEnemyData:
			world.Manager.Join(world.Components.FactionEnemy).Visit(ecs.Visit(func(entity ecs.Entity) {
				entityCount++
				sm.processEntityForSave(entity, world, &entities, processedEntities)
			}))
		case ComponentInParty:
			world.Manager.Join(world.Components.InParty).Visit(ecs.Visit(func(entity ecs.Entity) {
				entityCount++
				sm.processEntityForSave(entity, world, &entities, processedEntities)
			}))
		case "Name":
			world.Manager.Join(world.Components.Name).Visit(ecs.Visit(func(entity ecs.Entity) {
				entityCount++
				sm.processEntityForSave(entity, world, &entities, processedEntities)
			}))
		case "Pools":
			world.Manager.Join(world.Components.Pools).Visit(ecs.Visit(func(entity ecs.Entity) {
				entityCount++
				sm.processEntityForSave(entity, world, &entities, processedEntities)
			}))
		case "Attributes":
			world.Manager.Join(world.Components.Attributes).Visit(ecs.Visit(func(entity ecs.Entity) {
				entityCount++
				sm.processEntityForSave(entity, world, &entities, processedEntities)
			}))
		case ComponentItem:
			world.Manager.Join(world.Components.Item).Visit(ecs.Visit(func(entity ecs.Entity) {
				entityCount++
				sm.processEntityForSave(entity, world, &entities, processedEntities)
			}))
		case ComponentLocationInBackpack:
			world.Manager.Join(world.Components.ItemLocationInBackpack).Visit(ecs.Visit(func(entity ecs.Entity) {
				entityCount++
				sm.processEntityForSave(entity, world, &entities, processedEntities)
			}))
		case ComponentLocationEquipped:
			world.Manager.Join(world.Components.ItemLocationEquipped).Visit(ecs.Visit(func(entity ecs.Entity) {
				entityCount++
				sm.processEntityForSave(entity, world, &entities, processedEntities)
			}))
		case ComponentLocationOnField:
			world.Manager.Join(world.Components.ItemLocationOnField).Visit(ecs.Visit(func(entity ecs.Entity) {
				entityCount++
				sm.processEntityForSave(entity, world, &entities, processedEntities)
			}))
		case ComponentLocationNone:
			world.Manager.Join(world.Components.ItemLocationNone).Visit(ecs.Visit(func(entity ecs.Entity) {
				entityCount++
				sm.processEntityForSave(entity, world, &entities, processedEntities)
			}))
		case "Description":
			world.Manager.Join(world.Components.Description).Visit(ecs.Visit(func(entity ecs.Entity) {
				entityCount++
				sm.processEntityForSave(entity, world, &entities, processedEntities)
			}))
		case "Wearable":
			world.Manager.Join(world.Components.Wearable).Visit(ecs.Visit(func(entity ecs.Entity) {
				entityCount++
				sm.processEntityForSave(entity, world, &entities, processedEntities)
			}))
		case "Card":
			world.Manager.Join(world.Components.Card).Visit(ecs.Visit(func(entity ecs.Entity) {
				entityCount++
				sm.processEntityForSave(entity, world, &entities, processedEntities)
			}))
		case "Material":
			world.Manager.Join(world.Components.Material).Visit(ecs.Visit(func(entity ecs.Entity) {
				entityCount++
				sm.processEntityForSave(entity, world, &entities, processedEntities)
			}))
		case "Consumable":
			world.Manager.Join(world.Components.Consumable).Visit(ecs.Visit(func(entity ecs.Entity) {
				entityCount++
				sm.processEntityForSave(entity, world, &entities, processedEntities)
			}))
		case "Attack":
			world.Manager.Join(world.Components.Attack).Visit(ecs.Visit(func(entity ecs.Entity) {
				entityCount++
				sm.processEntityForSave(entity, world, &entities, processedEntities)
			}))
		case "Recipe":
			world.Manager.Join(world.Components.Recipe).Visit(ecs.Visit(func(entity ecs.Entity) {
				entityCount++
				sm.processEntityForSave(entity, world, &entities, processedEntities)
			}))
		case ComponentEquipmentChanged:
			world.Manager.Join(world.Components.EquipmentChanged).Visit(ecs.Visit(func(entity ecs.Entity) {
				entityCount++
				sm.processEntityForSave(entity, world, &entities, processedEntities)
			}))
		case "ProvidesHealing":
			world.Manager.Join(world.Components.ProvidesHealing).Visit(ecs.Visit(func(entity ecs.Entity) {
				entityCount++
				sm.processEntityForSave(entity, world, &entities, processedEntities)
			}))
		case "InflictsDamage":
			world.Manager.Join(world.Components.InflictsDamage).Visit(ecs.Visit(func(entity ecs.Entity) {
				entityCount++
				sm.processEntityForSave(entity, world, &entities, processedEntities)
			}))
		}
	}

	return WorldSaveData{
		Entities: entities,
	}
}

// processEntityForSave はエンティティを処理してセーブデータに追加
func (sm *SerializationManager) processEntityForSave(entity ecs.Entity, world w.World, entities *[]EntitySaveData, processed map[ecs.Entity]bool) {
	// 既に処理済みかチェック
	if processed[entity] {
		return
	}
	processed[entity] = true

	// 安定IDを取得
	stableID := sm.stableIDManager.GetStableID(entity)

	// エンティティデータを作成
	entityData := EntitySaveData{
		StableID:   stableID,
		Components: make(map[string]ComponentData),
	}

	// 各コンポーネント型をチェック
	for _, typeInfo := range sm.componentRegistry.GetAllTypes() {
		if data, hasComponent := typeInfo.ExtractFunc(world, entity); hasComponent {
			// エンティティ参照の処理
			processedData := sm.processEntityReferences(data, typeInfo)

			entityData.Components[typeInfo.Name] = ComponentData{
				Type: typeInfo.Name,
				Data: processedData,
			}
		}
	}

	// コンポーネントがあるエンティティのみ保存
	if len(entityData.Components) > 0 {
		*entities = append(*entities, entityData)
	}
}

// processEntityReferences はエンティティ参照を安定IDに変換
func (sm *SerializationManager) processEntityReferences(data interface{}, typeInfo *ComponentTypeInfo) interface{} {
	// AIVisionのTargetEntityを特別処理
	if typeInfo.Name == ComponentAIVision {
		if vision, ok := data.(gc.AIVision); ok {
			visionRef := struct {
				ViewDistance gc.Pixel  `json:"view_distance"`
				TargetRef    *StableID `json:"target_ref,omitempty"`
			}{
				ViewDistance: vision.ViewDistance,
			}

			if vision.TargetEntity != nil {
				targetStableID := sm.stableIDManager.GetStableID(*vision.TargetEntity)
				visionRef.TargetRef = &targetStableID
			}

			return visionRef
		}
	}

	// LocationEquippedのOwnerを特別処理
	if typeInfo.Name == ComponentLocationEquipped {
		if equipped, ok := data.(gc.LocationEquipped); ok {
			ownerStableID := sm.stableIDManager.GetStableID(equipped.Owner)
			equippedRef := struct {
				OwnerRef      StableID                 `json:"owner_ref"`
				EquipmentSlot gc.EquipmentSlotNumber `json:"equipment_slot"`
			}{
				OwnerRef:      ownerStableID,
				EquipmentSlot: equipped.EquipmentSlot,
			}
			return equippedRef
		}
	}

	// 他の型はそのまま返す
	return data
}

// restoreWorldData はセーブデータからワールドを復元
func (sm *SerializationManager) restoreWorldData(world w.World, worldData WorldSaveData) error {
	// 第1段階: 全エンティティを作成して安定IDマッピング
	entityMap := make(map[StableID]ecs.Entity)
	entityDataMap := make(map[StableID]EntitySaveData)

	for _, entityData := range worldData.Entities {
		entity := world.Manager.NewEntity()

		// 安定IDマッピングを登録
		err := sm.stableIDManager.RegisterEntity(entity, entityData.StableID)
		if err != nil {
			return fmt.Errorf("failed to register entity mapping: %w", err)
		}

		entityMap[entityData.StableID] = entity
		entityDataMap[entityData.StableID] = entityData
	}

	// 第2段階: コンポーネントを復元（エンティティ参照なし）
	for stableID, entityData := range entityDataMap {
		entity := entityMap[stableID]

		for componentName, componentData := range entityData.Components {
			typeInfo, exists := sm.componentRegistry.GetTypeInfoByName(componentName)
			if !exists {
				fmt.Printf("Warning: unknown component type: %s\n", componentName)
				continue
			}

			// JSONからコンポーネントデータを復元
			restoredData, err := sm.restoreComponentData(componentData.Data, typeInfo)
			if err != nil {
				return fmt.Errorf("failed to restore component %s: %w", componentName, err)
			}

			// コンポーネントをエンティティに追加
			err = typeInfo.RestoreFunc(world, entity, restoredData)
			if err != nil {
				return fmt.Errorf("failed to add component %s to entity: %w", componentName, err)
			}
		}
	}

	// 第3段階: エンティティ参照を解決
	for stableID, entityData := range entityDataMap {
		entity := entityMap[stableID]

		for componentName, componentData := range entityData.Components {
			typeInfo, exists := sm.componentRegistry.GetTypeInfoByName(componentName)
			if !exists || typeInfo.ResolveRefFunc == nil {
				continue
			}

			err := sm.resolveEntityReferences(world, entity, componentData.Data, typeInfo)
			if err != nil {
				return fmt.Errorf("failed to resolve references for %s: %w", componentName, err)
			}
		}
	}

	return nil
}

// restoreComponentData はJSONデータからコンポーネントデータを復元
func (sm *SerializationManager) restoreComponentData(jsonData interface{}, typeInfo *ComponentTypeInfo) (interface{}, error) {
	// AIVisionを特別処理（カスタムシリアライズ形式のため）
	if typeInfo.Name == "AIVision" {
		dataMap, ok := jsonData.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid AIVision JSON data type: %T", jsonData)
		}

		// ViewDistanceを取得
		viewDistanceVal, exists := dataMap["view_distance"]
		if !exists {
			return nil, fmt.Errorf("view_distance not found in AIVision data")
		}
		viewDistance, ok := viewDistanceVal.(float64)
		if !ok {
			return nil, fmt.Errorf("invalid view_distance type: %T", viewDistanceVal)
		}

		// AIVision構造体を作成（TargetEntityは後で解決）
		vision := gc.AIVision{
			ViewDistance: gc.Pixel(viewDistance),
			TargetEntity: nil,
		}
		return vision, nil
	}

	// ProvidesHealingを特別処理（Amounterインターフェースのため）
	if typeInfo.Name == "ProvidesHealing" {
		// この型はreflection.goのrestoreProvidesHealingで処理されるため、
		// ここではデータをそのまま返す
		return jsonData, nil
	}

	// LocationEquippedを特別処理（カスタムシリアライズ形式のため）
	if typeInfo.Name == "LocationEquipped" {
		dataMap, ok := jsonData.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid LocationEquipped JSON data type: %T", jsonData)
		}

		// EquipmentSlotを取得
		equipmentSlotVal, exists := dataMap["equipment_slot"]
		if !exists {
			return nil, fmt.Errorf("equipment_slot not found in LocationEquipped data")
		}
		equipmentSlot, ok := equipmentSlotVal.(float64)
		if !ok {
			return nil, fmt.Errorf("invalid equipment_slot type: %T", equipmentSlotVal)
		}

		// LocationEquipped構造体を作成（Ownerは後で解決）
		equipped := gc.LocationEquipped{
			Owner:         0, // 一時的に無効なエンティティID
			EquipmentSlot: gc.EquipmentSlotNumber(equipmentSlot),
		}
		return equipped, nil
	}

	// 通常のコンポーネント処理
	// JSONデータをバイトに変換
	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		return nil, err
	}

	// 型に応じて適切な構造体を作成
	targetValue := reflect.New(typeInfo.Type).Interface()

	// JSONからデコード
	err = json.Unmarshal(jsonBytes, targetValue)
	if err != nil {
		return nil, err
	}

	// ポインタから値を取得
	return reflect.ValueOf(targetValue).Elem().Interface(), nil
}

// resolveEntityReferences はエンティティ参照を解決
func (sm *SerializationManager) resolveEntityReferences(world w.World, entity ecs.Entity, jsonData interface{}, typeInfo *ComponentTypeInfo) error {
	if typeInfo.Name == "AIVision" {
		// JSONデータを変換
		jsonBytes, err := json.Marshal(jsonData)
		if err != nil {
			return err
		}

		var visionRef struct {
			ViewDistance gc.Pixel  `json:"view_distance"`
			TargetRef    *StableID `json:"target_ref,omitempty"`
		}

		err = json.Unmarshal(jsonBytes, &visionRef)
		if err != nil {
			return err
		}

		// AIVisionコンポーネントを取得
		if entity.HasComponent(world.Components.AIVision) {
			vision := world.Components.AIVision.Get(entity).(*gc.AIVision)

			// エンティティ参照を解決
			if visionRef.TargetRef != nil {
				if targetEntity, exists := sm.stableIDManager.GetEntity(*visionRef.TargetRef); exists {
					vision.TargetEntity = &targetEntity
				} else {
					fmt.Printf("Warning: target entity not found for stable ID: %v\n", *visionRef.TargetRef)
				}
			}
		}
	}

	if typeInfo.Name == "LocationEquipped" {
		// JSONデータを変換
		jsonBytes, err := json.Marshal(jsonData)
		if err != nil {
			return err
		}

		var equippedRef struct {
			OwnerRef      StableID                 `json:"owner_ref"`
			EquipmentSlot gc.EquipmentSlotNumber `json:"equipment_slot"`
		}

		err = json.Unmarshal(jsonBytes, &equippedRef)
		if err != nil {
			return err
		}

		// LocationEquippedコンポーネントを取得
		if entity.HasComponent(world.Components.ItemLocationEquipped) {
			equipped := world.Components.ItemLocationEquipped.Get(entity).(*gc.LocationEquipped)

			// エンティティ参照を解決
			if ownerEntity, exists := sm.stableIDManager.GetEntity(equippedRef.OwnerRef); exists {
				equipped.Owner = ownerEntity
			} else {
				fmt.Printf("Warning: owner entity not found for stable ID: %v\n", equippedRef.OwnerRef)
			}
		}
	}

	return nil
}

// clearWorld はワールドの全エンティティをクリア
func (sm *SerializationManager) clearWorld(world w.World) {
	// 全エンティティを削除
	entitiesToDelete := make([]ecs.Entity, 0)

	world.Manager.Join().Visit(ecs.Visit(func(entity ecs.Entity) {
		entitiesToDelete = append(entitiesToDelete, entity)
	}))

	for _, entity := range entitiesToDelete {
		world.Manager.DeleteEntity(entity)
	}
}

// GetStableIDManager は安定IDマネージャーを取得
func (sm *SerializationManager) GetStableIDManager() *StableIDManager {
	return sm.stableIDManager
}

// GetComponentRegistry はコンポーネントレジストリを取得
func (sm *SerializationManager) GetComponentRegistry() *ComponentRegistry {
	return sm.componentRegistry
}
