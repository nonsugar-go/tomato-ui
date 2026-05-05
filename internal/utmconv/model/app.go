package model

type App struct {
	Filename      string
	Vendor        string
	To            string
	AppConfig     AppConfig
	Tag           []Tag
	Addresses     []Address
	AddressGroups []AddressGroup
	Services      []Service
	ServiceGroups []ServiceGroup
	Policies      []Policy
}

type AppConfig struct {
	Description string           `json:"_description"`
	CheckPoint  CheckPointConfig `json:"checkpoint"`
	PaloAlto    PaloAltoConfig   `json:"paloalto"`
}

type CheckPointConfig struct {
	Cli CheckPointCli `json:"cli"`
}

type PaloAltoConfig struct {
	Conf PaloAltoConf `json:"conf"`
}

type CheckPointCli struct {
	MgmtCliUser struct {
		Description string `json:"_description"`
		Value       string `json:"value"`
	} `json:"mgmt_cli_user"`

	MgmtCliPassword struct {
		Description string `json:"_description"`
		Value       string `json:"value"`
	} `json:"mgmt_cli_password"`

	IgnoreWarnings struct {
		Description string `json:"_description"`
		Value       bool   `json:"value"`
	} `json:"ignore-warnings"`

	AccessRuleLayer struct {
		Description string `json:"_description"`
		Value       string `json:"value"`
	} `json:"access_rule_layer"`

	AccessRuleSection struct {
		Description string `json:"_description"`
		Value       string `json:"value"`
	} `json:"access_rule_section"`

	PredefinedServices struct {
		Description string   `json:"_description"`
		Value       []string `json:"value"`
	} `json:"predefined_services"`

	ServiceReplacementMap struct {
		Description string       `json:"_description"`
		Value       []ServiceMap `json:"value"`
	} `json:"service_replacement_map"`
}

type ServiceMap struct {
	Before string `json:"before"`
	After  string `json:"after"`
}

type PaloAltoConf struct {
	ApplicationDefaultReplacementMap struct {
		Description string      `json:"_description"`
		Value       []AppSvcMap `json:"value"`
	} `json:"application_default_replacement_map"`
}

type AppSvcMap struct {
	Application string   `json:"application"`
	Services    []string `json:"services"`
}

func NewDefaultAppConfig() *AppConfig {
	cfg := AppConfig{Description: "utmconv の設定ファイル"}

	cpCli := CheckPointCli{}
	cpCli.MgmtCliUser.Description = "mgmt_cli tool のユーザー名"
	cpCli.MgmtCliUser.Value = "secadmin"
	cpCli.MgmtCliPassword.Description = "mgmt_cli tool のパスワード"
	cpCli.MgmtCliPassword.Value = "Lab@12345"
	cpCli.IgnoreWarnings.Description = "mgmt_cli tool に ignore-warnings true を付加するかどうか？"
	cpCli.IgnoreWarnings.Value = true
	cpCli.AccessRuleLayer.Description = "cli 出力時の access-rule の layer"
	cpCli.AccessRuleLayer.Value = "Network"
	cpCli.AccessRuleSection.Description = "cli 出力時の access-rule を追加するセクション タイトル"
	cpCli.AccessRuleSection.Value = "New rules"
	cpCli.PredefinedServices.Description = "事前定義サービス名の配列"
	cpCli.PredefinedServices.Value = []string{
		"ICMP Protocol", "echo-request", "traceroute", "ssh", "syslog", "ESP", "IKE", "IKE_NAT_TRAVERSAL", "Instagram"}
	cpCli.ServiceReplacementMap.Description = "サービス名置換用マップ"
	cpCli.ServiceReplacementMap.Value = []ServiceMap{
		{Before: "service-http", After: "http"},
		{Before: "service-https", After: "https"},
	}
	cfg.CheckPoint.Cli = cpCli

	paConf := PaloAltoConf{}
	paConf.ApplicationDefaultReplacementMap.Description = "サービスが application-default または any のときのサービス名置換用マップ"
	paConf.ApplicationDefaultReplacementMap.Value = []AppSvcMap{
		{Application: "icmp", Services: []string{"ICMP Protocol"}}, // check point のサービス icmp-proto は返答を許可が無効、Accounting もできない
		{Application: "ping", Services: []string{"echo-request"}},
		{Application: "traceroute", Services: []string{"traceroute"}},
		{Application: "ssh", Services: []string{"ssh"}},
		{Application: "ssh-tunnel", Services: []string{"ssh"}},
		{Application: "syslog", Services: []string{"syslog"}},
		{Application: "ipsec-esp", Services: []string{"ESP", "IKE", "IKE_NAT_TRAVERSAL"}},
		{Application: "ipsec-esp-udp", Services: []string{"ESP", "IKE", "IKE_NAT_TRAVERSAL"}},
		{Application: "instagram-base", Services: []string{"Instagram"}},
	}
	cfg.PaloAlto.Conf = paConf

	return &cfg
}
