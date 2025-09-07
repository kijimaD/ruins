# GameLog Package

ゲームログシステムのパッケージです。メソッドチェーンによる色付きログを提供します。

## 主な機能

- メソッドチェーンによる直感的なログ作成
- 色付きテキストフラグメント
- プリセット関数による統一的な色付け
- スレッドセーフなログストレージ

## 基本的な使い方

```go
// シンプルなログ
gamelog.New().
    Append("プレイヤーがアイテムを入手した").
    Log(gamelog.LogKindField)

// 色付きログ
gamelog.New().
    PlayerName("Hero").
    Append("が").
    ItemName("Iron Sword").
    Append("を入手した。").
    Log(gamelog.LogKindField)
```

## プリセット関数

### 基本プリセット
| 関数 | 色 | 用途 |
|------|----|------|
| `Success(text)` | 緑色 | 成功メッセージ |
| `Warning(text)` | 黄色 | 警告メッセージ |
| `Error(text)` | 赤色 | エラーメッセージ |
| `System(text)` | 水色 | システムメッセージ |

### ゲーム要素プリセット
| 関数 | 色 | 用途 |
|------|----|------|
| `PlayerName(name)` | 緑色 | プレイヤー名 |
| `NPCName(name)` | 黄色 | NPC名 |
| `ItemName(item)` | シアン色 | アイテム名 |
| `Location(place)` | オレンジ色 | 場所名 |
| `Action(action)` | 紫色 | アクション名 |
| `Money(amount)` | 黄色 | 金額 |
| `Damage(num)` | 赤色 | ダメージ数値 |

### 戦闘専用プリセット
| 関数 | 色 | 用途 |
|------|----|------|
| `Encounter(text)` | 赤色 | 敵との遭遇 |
| `Victory(text)` | 緑色 | 勝利メッセージ |
| `Defeat(text)` | 赤色 | 敗北メッセージ |
| `Magic(text)` | 紫色 | 魔法関連 |

## ログ種別

| 種別 | 用途 |
|------|------|
| `LogKindField` | フィールド探索ログ |
| `LogKindBattle` | 戦闘ログ |
| `LogKindScene` | シーンログ |

## ログストレージ

```go
var (
    FieldLog  *SafeSlice  // フィールドログ
    BattleLog *SafeSlice  // 戦闘ログ
    SceneLog  *SafeSlice  // シーンログ
)
```

色付きエントリの取得：
```go
entries := gamelog.FieldLog.GetRecentEntries(5)
for _, entry := range entries {
    for _, fragment := range entry.Fragments {
        // fragment.Text と fragment.Color を使用
    }
}
```

## カスタム色

```go
import "github.com/kijimaD/ruins/lib/colors"

// 定義済み色を使用
gamelog.New().
    ColorRGBA(colors.ColorBlue).
    Append("青色のテキスト").
    Log(gamelog.LogKindField)

// カスタム色を作成
gamelog.New().
    ColorRGBA(colors.NamedColor(255, 0, 0)). // 赤色
    Append("カスタム色のテキスト").
    Log(gamelog.LogKindField)
```
