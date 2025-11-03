package cmd

import (
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/kijimaD/ruins/lib/raw"
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
	writeString(file, "各ステージ・階層ごとのアイテム出現確率を示します。\n\n")

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

	// 全アイテム名を収集（ヘッダー用）
	itemNames := make(map[string]bool)
	for _, entry := range table.Entries {
		itemNames[entry.ItemName] = true
	}

	// ソートされたアイテム名リスト
	sortedItems := make([]string, 0, len(itemNames))
	for name := range itemNames {
		sortedItems = append(sortedItems, name)
	}
	sort.Strings(sortedItems)

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

	// テーブルヘッダー
	header := "| 深度 |"
	for _, item := range sortedItems {
		header += fmt.Sprintf(" %s |", item)
	}
	writeString(file, header+"\n")

	// セパレータ
	separator := "|------|"
	for range sortedItems {
		separator += "--------|"
	}
	writeString(file, separator+"\n")

	// 各深度の行
	for depth := 1; depth <= maxDepth; depth++ {
		row := fmt.Sprintf("| %d |", depth)

		// 各アイテムの出現確率を計算
		probs := calculateProbabilities(table, depth)

		for _, item := range sortedItems {
			if prob, ok := probs[item]; ok {
				row += fmt.Sprintf(" %.1f%% |", prob*100)
			} else {
				row += " - |"
			}
		}
		writeString(file, row+"\n")
	}

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
