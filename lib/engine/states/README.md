# BaseState 使用方法

## 概要
BaseStateは、全てのstateで共通するtransition管理機能を提供します。

## 使用方法

### 1. 構造体での埋め込み
```go
type YourState struct {
    states.BaseState  // 埋め込み
    ui *ebitenui.UI
    // その他のフィールド
}
```

### 2. Updateメソッドでの使用
```go
func (st *YourState) Update(world w.World) states.Transition {
    // 通常の処理...
    
    // 即座に遷移したい場合
    if someCondition {
        return states.Transition{Type: states.TransPop}
    }
    
    // 他の処理で遷移を設定したい場合
    if anotherCondition {
        st.SetTransition(states.Transition{Type: states.TransSwitch, NewStates: []states.State{&NextState{}}})
    }
    
    // 最後に必ずConsumeTransitionを呼ぶ
    return st.ConsumeTransition()
}
```

### 3. 自動化される処理

- **OnResume時**: `ClearTransition()`が自動実行され、前のstate遷移がクリアされる
- **Enterキーリセット**: 重複実行防止のため、グローバルキー状態が自動リセット

## メリット

- ✅ 全てのstateで統一されたtransition処理
- ✅ OnResume時の自動クリア（`st.trans = nil`が不要）
- ✅ 共通のEnterキー重複実行防止
- ✅ 既存コードの`if st.trans != nil { ... }`パターンを`return st.ConsumeTransition()`に置き換えるだけ

## 移行方法

1. `states.BaseState`を埋め込む
2. `trans *states.Transition`フィールドを削除
3. Updateメソッドの最後を`return st.ConsumeTransition()`に変更
4. OnResumeでの`st.trans = nil`を削除（自動実行される）