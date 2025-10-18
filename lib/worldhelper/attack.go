package worldhelper

import (
	"fmt"

	gc "github.com/kijimaD/ruins/lib/components"
	"github.com/kijimaD/ruins/lib/raw"
	w "github.com/kijimaD/ruins/lib/world"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// GetAttackFromCommandTable は敵のCommandTableからランダムに攻撃を選択する
// digger_rsのNaturalAttackDefenseと同様に、カードエンティティを生成せずに攻撃パラメータを取得する
func GetAttackFromCommandTable(world w.World, enemyEntity ecs.Entity) (*gc.Attack, string, error) {
	// CommandTableコンポーネントを取得
	commandTableComp := world.Components.CommandTable.Get(enemyEntity)
	if commandTableComp == nil {
		return nil, "", fmt.Errorf("enemy has no CommandTable component")
	}

	commandTableName := commandTableComp.(*gc.CommandTable).Name
	rawMaster := world.Resources.RawMaster.(*raw.Master)

	// CommandTableを取得
	commandTable, err := rawMaster.GetCommandTable(commandTableName)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get command table: %w", err)
	}

	// 重み付きランダムでカード名を選択
	cardName := commandTable.SelectByWeight()
	if cardName == "" {
		return nil, "", fmt.Errorf("no card selected from command table")
	}

	// カード名からEntitySpecを取得（エンティティは生成しない）
	cardSpec, err := rawMaster.NewCardSpec(cardName)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get card spec: %w", err)
	}

	// Attackコンポーネントがない場合はエラー
	if cardSpec.Attack == nil {
		return nil, "", fmt.Errorf("card %s has no Attack component", cardName)
	}

	return cardSpec.Attack, cardName, nil
}

// GetAttackFromCard はカードエンティティから攻撃パラメータを取得する
// プレイヤーが装備しているカードエンティティから攻撃情報を取得する
func GetAttackFromCard(world w.World, cardEntity ecs.Entity) (*gc.Attack, string, error) {
	// Attackコンポーネントを取得
	attackComp := world.Components.Attack.Get(cardEntity)
	if attackComp == nil {
		return nil, "", fmt.Errorf("card has no Attack component")
	}

	// 名前を取得
	nameComp := world.Components.Name.Get(cardEntity)
	if nameComp == nil {
		return nil, "", fmt.Errorf("card has no Name component")
	}

	attack := attackComp.(*gc.Attack)
	cardName := nameComp.(*gc.Name).Name

	return attack, cardName, nil
}
