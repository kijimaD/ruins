// Package components はゲームコンポーネントの定義と実装を提供する。
//
// 責務:
// - ECS（Entity Component System）のコンポーネント定義
// - キャラクター、アイテム、フィールドオブジェクトの属性管理
// - 戦闘、移動、AI、描画などの機能別コンポーネント提供
//
// 使い分け:
// - GameComponentList: エンティティ作成時のコンポーネント情報格納
// - Components: ECSクエリで使用するコンポーネント実体
// - 各構造体: 個別のコンポーネントデータ
//
// 仕様:
// - NullComponent: 状態マーカー（Player, Dead, InParty等）
// - SliceComponent: データ保持（Pools, Attributes, Attack等）
// - 死亡状態はDeadコンポーネントで明示的に管理
// - HP.Current == 0 での死亡判定とDeadコンポーネント付与を併用
package components
