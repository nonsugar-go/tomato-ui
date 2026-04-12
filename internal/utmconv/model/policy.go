package model

type Policy struct {
	Name        string
	Description string

	// 有効/無効
	Enabled bool

	// マッチ条件
	Match PolicyMatch

	// 動作
	Action PolicyAction

	// 追加オプション
	Logging  Logging
	Schedule string
	Tags     []string
	Group    string

	// 拡張（ベンダー依存吸収用）
	Extensions map[string]any
}

type PolicyMatch struct {
	FromZones []string
	ToZones   []string

	Sources      []AddressRef
	Destinations []AddressRef

	Applications []string
	Services     []ServiceRef

	Users         []string
	HIPs          []string
	URLCategories []string

	NegateSource      bool
	NegateDestination bool
}

type PolicyAction struct {
	Type ActionType

	// セキュリティプロファイル（UTM系）
	Profiles []string
}

type ActionType string

const (
	ActionAllow ActionType = "allow"
	ActionDeny  ActionType = "deny"
	ActionDrop  ActionType = "drop"
	ActionReset ActionType = "reset"
)

type Logging struct {
	LogAtStart bool
	LogAtEnd   bool
	LogProfile string
}

type AddressRef struct {
	Name string
	// 将来的に type / subnet / fqdn なども持てる
}

type ServiceRef struct {
	Name string
}
