// Package effects はゲームエフェクトシステムの実装を提供する。ステート遷移は対象としない。
//
// 主要な設計原則:
//
// 1. 型安全性: Effectインターフェースによる型安全なエフェクト管理
// 2. 責務分離: エフェクト定義、ターゲット選択、実行処理の明確な分離
// 3. テスタビリティ: 各エフェクトの独立したテストが可能
// 4. 拡張性: 新しいエフェクトタイプの簡単な追加
//
// 基本的な使用方法:
//
//	processor := effects.NewProcessor()
//
//	// ダメージエフェクト
//	damage := effects.Damage{Amount: 50, Source: effects.DamageSourceWeapon}
//	processor.AddEffect(damage, &attacker, target)
//
//	// ターゲットセレクタを使用（戦闘時：ゲームログ出力あり）
//	healing := effects.Healing{Amount: gc.NumeralAmount{Numeral: 30}}
//	processor.AddTargetedEffectWithLogger(healing, &healer, effects.TargetParty{}, gamelog.BattleLog, world)
//
//	// 非戦闘時の回復（ゲームログ出力なし）
//	processor.AddTargetedEffect(healing, &healer, effects.TargetParty{}, world)
//
//	// エフェクト実行
//	if err := processor.Execute(world); err != nil {
//	    log.Printf("エフェクト実行エラー: %v", err)
//	}
//
// エフェクトタイプ:
//
// combat.go    - 戦闘関連エフェクト（ダメージ、Healing統合型、スタミナ）
// recovery.go  - 非戦闘時の全回復エフェクト（FullRecoveryHP/SP）
// item.go      - アイテム関連エフェクト（使用、消費）
//
// パッケージ構成:
//
// processor.go - エフェクト実行管理システム
// target.go    - ターゲット選択システム
package effects
