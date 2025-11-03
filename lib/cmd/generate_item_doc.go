package cmd

import (
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/kijimaD/ruins/lib/raw"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v2"
)

// CmdGenerateItemDoc はアイテム出現確率のドキュメントを生成するコマンド
var CmdGenerateItemDoc = &cli.Command{
	Name:        "generate-item-doc",
	Usage:       "generate-item-doc",
	Description: "Generate item spawn probability documentation",
	Action:      runGenerateItemDoc,
	Flags:       []cli.Flag{},
}

func runGenerateItemDoc(_ *cli.Context) error {
	// raw.tomlを読み込む
	master, err := raw.LoadFromFile("metadata/entities/raw/raw.toml")
	if err != nil {
		return fmt.Errorf("error loading raw.toml: %w", err)
	}

	// Markdownファイルを開く
	file, err := os.Create("docs/item_tables.md")
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Warning: failed to close file: %v\n", err)
		}
	}()

	// ヘッダー
	writeString(file, "# アイテム出現テーブル\n\n")
	writeString(file, "各ステージ・階層ごとのアイテム出現確率を示す。\n\n")

	// 各ItemTableを処理
	for _, table := range master.Raws.ItemTables {
		generateTableDoc(file, table)
	}

	log.Println("Generated docs/item_tables.md")
	return nil
}

// writeString はファイルへの書き込みとエラーハンドリングを行う
func writeString(file *os.File, s string) {
	if _, err := file.WriteString(s); err != nil {
		log.Fatalf("Error writing to file: %v\n", err)
	}
}

func generateTableDoc(file *os.File, table raw.ItemTable) {
	writeString(file, fmt.Sprintf("## %s\n\n", table.Name))

	// 最大深度を決定
	maxDepth := 0
	for _, entry := range table.Entries {
		if entry.MaxDepth > maxDepth {
			maxDepth = entry.MaxDepth
		}
	}
	if maxDepth == 0 {
		maxDepth = 50 // デフォルト最大深度
	}

	// 各深度での最大アイテム数を計算
	maxItems := 0
	for depth := 1; depth <= maxDepth; depth++ {
		probs := calculateProbabilities(table, depth)
		if len(probs) > maxItems {
			maxItems = len(probs)
		}
	}

	// tablewriterを初期化
	tw := tablewriter.NewWriter(file)

	// ヘッダーを設定（深度 + 最大アイテム数分のダミーカラム）
	header := []string{"深度"}
	for i := 0; i < maxItems; i++ {
		header = append(header, "-")
	}
	tw.SetHeader(header)

	// Markdown形式の設定
	tw.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	tw.SetCenterSeparator("|")
	tw.SetAutoFormatHeaders(false)

	// 各深度の行を追加
	for depth := 1; depth <= maxDepth; depth++ {
		// 各アイテムの出現確率を計算
		probs := calculateProbabilities(table, depth)

		// アイテム名をソート
		items := make([]string, 0, len(probs))
		for item := range probs {
			items = append(items, item)
		}
		sort.Strings(items)

		// 行データを作成
		row := []string{fmt.Sprintf("%d", depth)}

		// アイテムと確率を追加
		for _, item := range items {
			prob := probs[item]
			row = append(row, fmt.Sprintf("%s %.1f%%", item, prob*100))
		}

		// 不足分は空欄で埋める
		for i := len(items); i < maxItems; i++ {
			row = append(row, "-")
		}

		tw.Append(row)
	}

	// テーブルをレンダリング
	tw.Render()

	writeString(file, "\n")
}

func calculateProbabilities(table raw.ItemTable, depth int) map[string]float64 {
	// 深度範囲内のエントリをフィルタリング
	validEntries := make([]raw.ItemTableEntry, 0, len(table.Entries))
	for _, entry := range table.Entries {
		// MinDepthチェック（0は制限なし）
		if entry.MinDepth > 0 && depth < entry.MinDepth {
			continue
		}
		// MaxDepthチェック（0は制限なし）
		if entry.MaxDepth > 0 && depth > entry.MaxDepth {
			continue
		}
		validEntries = append(validEntries, entry)
	}

	// 総重みを計算
	var totalWeight float64
	for _, entry := range validEntries {
		totalWeight += entry.Weight
	}

	// 確率を計算
	probs := make(map[string]float64)
	if totalWeight > 0 {
		for _, entry := range validEntries {
			probs[entry.ItemName] = entry.Weight / totalWeight
		}
	}

	return probs
}
