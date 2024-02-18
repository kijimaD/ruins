package raw

type Weapon struct {
	Accuracy          int // 命中率。0~100%
	BaseDamage        int // ベース攻撃力
	AttackCount       int // 攻撃回数
	EnergyConsumption int // 攻撃で消費するエネルギー
}
