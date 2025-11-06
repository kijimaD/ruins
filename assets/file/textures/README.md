aseprite でパッキングしている。

ファイル名が数字で終わると aseprite が JSON の filename が正しく出力されないので、末尾にアンダースコアを追加している。

```
aseprite \
    --batch ./tiles/*.png \
    --sheet ./dist/tiles.png \
    --data ./dist/tiles.json \
    --sheet-type packed \
    --format json-array \
    --filename-format "{title}"

aseprite \
    --batch ./single/*.png \
    --sheet ./dist/single.png \
    --data ./dist/single.json \
    --sheet-type packed \
    --format json-array \
    --filename-format "{title}"
```

```
tree -d
.
├── dist    # asepriteによる生成物。手動ではいじらない
├── single  # 1つで機能するスプライト
└── tiles   # オートタイルのスプライト
```
