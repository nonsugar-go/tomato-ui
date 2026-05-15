package checkpoint

type ShowPackage struct {
	Hosts []CheckPointHost
}

type CheckPointHost struct {
	NatSettings      NatSettings      `json:"nat-settings,omitempty"`
	Interfaces       []any            `json:"interfaces"`
	Comments         string           `json:"comments"`
	Color            string           `json:"color"`
	Icon             string           `json:"icon"`
	MetaInfo         MetaInfo         `json:"meta-info"`
	Type             string           `json:"type"`
	Tags             []any            `json:"tags"`
	UID              string           `json:"uid"`
	Ipv4Address      string           `json:"ipv4-address"`
	Domain           Domain           `json:"domain"`
	Name             string           `json:"name"`
	ReadOnly         bool             `json:"read-only"`
	AvailableActions AvailableActions `json:"available-actions"`
}

type NatSettings struct {
	Ipv4Address string `json:"ipv4-address"`
	Method      string `json:"method"`
	HideBehind  string `json:"hide-behind"`
	InstallOn   string `json:"install-on"`
	AutoRule    bool   `json:"auto-rule"`
	Ipv6Address string `json:"ipv6-address"`
}

type LastModifyTime struct {
	Iso8601 string `json:"iso-8601"`
	Posix   int64  `json:"posix"`
}

type CreationTime struct {
	Iso8601 string `json:"iso-8601"`
	Posix   int64  `json:"posix"`
}

type MetaInfo struct {
	Creator         string         `json:"creator"`
	ValidationState string         `json:"validation-state"`
	LastModifyTime  LastModifyTime `json:"last-modify-time"`
	CreationTime    CreationTime   `json:"creation-time"`
	Lock            string         `json:"lock"`
	LastModifier    string         `json:"last-modifier"`
}

type Domain struct {
	UID        string `json:"uid"`
	DomainType string `json:"domain-type"`
	Name       string `json:"name"`
}

type AvailableActions struct {
	Edit   string `json:"edit"`
	Clone  string `json:"clone"`
	Delete string `json:"delete"`
}
