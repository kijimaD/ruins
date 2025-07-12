デッキ構築型RPG(予定)。

- [play](https://kijimad.github.io/ruins/)
- [godoc](https://kijimad.github.io/ruins/godoc/pkg/github.com/kijimaD/ruins/lib/)

## Play Images

<img src="./vrtimages/Battle.png" width="50%" /><img src="./vrtimages/Dungeon.png" width="50%" />

<img src="./vrtimages/HomeMenu.png" width="50%" /><img src="./vrtimages/MainMenu.png" width="50%" />

<img src="./vrtimages/Intro.png" width="50%" /><img src="./vrtimages/InventoryMenu.png" width="50%" />

<img src="./vrtimages/CraftMenu.png" width="50%" /><img src="./vrtimages/EquipMenu.png" width="50%" />

<img src="./vrtimages/DungeonMenu.png" width="50%" /><img src="./vrtimages/DungeonSelect.png" width="50%" />

<img src="./vrtimages/DebugMenu.png" width="50%" /><img src="./vrtimages/GameOver.png" width="50%" />

各画像は全自動でキャプチャされる。最新なことをCIで保証している。

## Reference

ゲーム作成で参考にしたコード等。

- ゲームシステム面で、KONAMIのビデオゲーム『パワプロクンポケット』シリーズ10・11・12を参考にした
- https://github.com/x-hgg-x/sokoban-go
  - 最初にコピペして作成をはじめ、改変していった
  - ECSの使い方まわりで大きく参考にした
- https://github.com/x-hgg-x/goecsengine
  - sokoban-goの依存ライブラリ。改変しないと不便なことがけっこうあったので、コピペしてコードに含めた
- https://bfnihtly.bracketproductions.com/
  - 設定ファイルによるファクトリまわりを参考にした
- https://krkrz.github.io/krkr2doc/kag3doc/contents/
  - サウンドノベルに必要な記法を参考にした
- https://ebitengine.org/en/examples/raycasting.html
  - レイキャストの実装の参考にした

使用した素材類。

- https://hpgpixer.jp/image_icons/vehicle/icon_vehicle.html
- http://jikasei.me/font/jf-dotfont/
- https://github.com/googlefonts/morisawa-biz-ud-gothic
- https://www.pixilart.com

## TODO

- [ ] フォント管理、リソース管理をworldから分離する
- [ ] UIを共通化、わかりやすくする
- [ ] stateにUIとロジックすべて入っているのを直す
- [ ] UI resouceを分離する
- [ ] テーブル表示UIを作る(性能表示用)
- [ ] キーボードベースの操作にする
  - [ ] キーボード操作のエミュレートで自動テスト可能にする

## Memo

- UI
  - 実装と外観を分離して理解しやすい
  - キーボードベースの操作である
  - テスト可能である
  - メニュー操作が多い
  - 共通のメニューコンポーネントを使用する
    - コンポーネントは状態を持つ。コンポーネントの変化した状態によって各ステートでのアクションを発火したり変更を反映できる
    - テスト可能である
    - メニューではカーソルキーによる上下の選択と、タブキーによるタブの選択ができる
