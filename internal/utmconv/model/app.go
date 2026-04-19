package model

type App struct {
	Filename       string
	Vendor         string
	To             string
	IgnoreWarnings bool // 変換時の警告を無視するかどうか
	Tag            []Tag
	Addresses      []Address
	AddressGroups  []AddressGroup
	Services       []Service
	ServiceGroups  []ServiceGroup
	Policies       []Policy
}
