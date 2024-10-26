package raw

type DropTable struct {
	Name    string
	XpBase  float64
	Entries []DropTableEntry `toml:"entries"`
}

type DropTableEntry struct {
	Material string
	Weight   float64
}
