// Package effects はゲームエフェクトシステムの実装を提供する。
//
// 新しいアーキテクチャ:
//
// このパッケージは型安全で拡張性の高いエフェクトシステムを提供します。
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
//	damage := effects.CombatDamage{Amount: 50, Source: effects.DamageSourceWeapon}
//	processor.AddEffect(damage, &attacker, target)
//
//	// ターゲットセレクタを使用
//	healing := effects.CombatHealing{Amount: gc.NumeralAmount{Numeral: 30}}
//	processor.AddTargetedEffect(healing, &healer, effects.PartyTargets{}, world)
//
//	// エフェクト実行
//	if err := processor.Execute(world); err != nil {
//	    log.Printf("エフェクト実行エラー: %v", err)
//	}
//
// エフェクトタイプ:
//
// combat.go    - 戦闘関連エフェクト（ダメージ、回復、スタミナ）
// recovery.go  - 非戦闘時の回復エフェクト（ログ出力なし）
// movement.go  - 移動関連エフェクト（ワープ、脱出）
// item.go      - アイテム関連エフェクト（使用、消費）
//
// 下位互換性:
//
// global.go では既存コードとの互換性を保つためのアダプター層を提供しています。
// 既存のAddEffect(), RunEffectQueue(), ItemTrigger()関数は引き続き使用できます。
package effects
