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
	NATRules      []NATRule
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

	NatRulePackage struct {
		Description string `json:"_description"`
		Value       string `json:"value"`
	} `json:"nat_rule_package"`

	PredefinedServices struct {
		Description string   `json:"_description"`
		Value       []string `json:"value"`
	} `json:"predefined_services"`

	ZoneReplacementMap struct {
		Description string    `json:"_description"`
		Value       []ZoneMap `json:"value"`
	} `json:"zone_replacement_map"`

	AddressReplacementMap struct {
		Description string       `json:"_description"`
		Value       []ServiceMap `json:"value"`
	} `json:"address_replacement_map"`

	ServiceReplacementMap struct {
		Description string       `json:"_description"`
		Value       []ServiceMap `json:"value"`
	} `json:"service_replacement_map"`
}

type ZoneMap struct {
	Before string `json:"before"`
	After  string `json:"after"`
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
	cpCli.NatRulePackage.Description = "cli 出力時の nat-rule の package"
	cpCli.NatRulePackage.Value = "standard"
	cpCli.PredefinedServices.Description = "事前定義サービス名・サービス グループ名の配列"
	cpCli.PredefinedServices.Value = []string{
		"ICMP Protocol", "Instagram",
		// Check Point の事前定義済みサービス (tcp)
		"AOL",
		"AP-Defender",
		"AT-Defender",
		"Backage",
		"BGP",
		"Bionet-Setup",
		"CheckPointExchangeAgent",
		"Citrix_ICA",
		"ConnectedOnLine",
		"CP_Exnet_PK",
		"CP_Exnet_resolve",
		"CP_redundant",
		"CP_reporting",
		"CP_rtm",
		"CP_seam",
		"CP_SmartPortal",
		"CP_SSL_Network_Extender",
		"CPD",
		"CPD_amon",
		"CPM",
		"CPMI",
		"CrackDown",
		"CreativePartnerClnt",
		"CreativePartnerSrvr",
		"DaCryptic",
		"DameWare",
		"daytime-tcp",
		"DerSphere",
		"DerSphere_II",
		"Direct_Connect_TCP",
		"discard-tcp",
		"DNP3",
		"domain-tcp",
		"DoT",
		"echo-tcp",
		"EDGE",
		"eDonkey_4661",
		"eDonkey_4662",
		"Entrust-Admin",
		"Entrust-KeyMgmt",
		"exec",
		"FIBMGR",
		"finger",
		"Freak2k",
		"ftp",
		"ftp-bidir",
		"ftp-pasv",
		"ftp-port",
		"FW1",
		"FW1_amon",
		"FW1_clntauth_http",
		"FW1_clntauth_telnet",
		"FW1_CPRID",
		"FW1_cvp",
		"FW1_ela",
		"FW1_ica_mgmt_tools",
		"FW1_ica_pull",
		"FW1_ica_push",
		"FW1_ica_services",
		"FW1_key",
		"FW1_lea",
		"FW1_log",
		"FW1_mgmt",
		"FW1_netso",
		"FW1_omi",
		"FW1_omi-sic",
		"FW1_pslogon",
		"FW1_pslogon_NG",
		"FW1_sam",
		"FW1_sds_logon",
		"FW1_sds_logon_NG",
		"FW1_snauth",
		"FW1_topo",
		"FW1_uaa",
		"FW1_ufp",
		"GateCrasher",
		"GNUtella_rtr_TCP",
		"GNUtella_TCP",
		"gopher",
		"GoToMyPC",
		"H323",
		"H323_any",
		"HackaTack_31785",
		"HackaTack_31787",
		"HackaTack_31788",
		"HackaTack_31790",
		"HackaTack_31792",
		"Hotline_client",
		"http",
		"HTTP_and_HTTPS_proxy",
		"HTTP_proxy",
		"https",
		"HTTPS_proxy",
		"icap",
		"ICCP",
		"ICKiller",
		"ident",
		"IKE_NAT_TRAVERSAL_TCP",
		"IKE_tcp",
		"imap",
		"IMAP-SSL",
		"iMesh",
		"InCommand",
		"indeni",
		"IPSO_Clustering_Mgmt_Protocol",
		"irc1",
		"irc2",
		"jabber",
		"Jade",
		"Kaos",
		"KaZaA",
		"Kerberos_v5_TCP",
		"Kuang2",
		"ldap",
		"ldap-ssl",
		"login",
		"lotus",
		"lpdw0rm",
		"Madster",
		"microsoft-ds",
		"Mneah",
		"Modbus",
		"MS-SQL-Monitor",
		"MS-SQL-Server",
		"MSN_Messenger_File_Transfer",
		"MSNMS",
		"MSNP",
		"Multidropper",
		"MySQL",
		"Napster_Client_6600-6699",
		"Napster_directory_4444",
		"Napster_directory_5555",
		"Napster_directory_6666",
		"Napster_directory_7777",
		"Napster_directory_8888_primary",
		"Napster_redirector",
		"nbsession",
		"NCP",
		"netbios-session",
		"netshow",
		"netstat",
		"nfsd-tcp",
		"nntp",
		"ntp-tcp",
		"OAS-NameServer",
		"OAS-ORB",
		"OPC",
		"OpenWindows",
		"Orbix-1570",
		"Orbix-1571",
		"pcANYWHERE-data",
		"pcTELECOMMUTE-FileSync",
		"pop-2",
		"pop-3",
		"POP3S",
		"Port_6667_trojans",
		"PostgreSQL",
		"pptp-tcp",
		"RainWall_Command",
		"RAT",
		"Real-Audio",
		"RealSecure",
		"Remote_Desktop_Protocol",
		"Remote_Storm",
		"rtsp",
		"SCCP",
		"securidprop",
		"Shadyshell",
		"shell",
		"SIC-TCP",
		"sip-tcp",
		"sip_any-tcp",
		"sip_tls_authentication",
		"sip_tls_not_inspected",
		"SkyDance-T",
		"smb",
		"smtp",
		"SMTPS",
		"snmp-tcp",
		"SocketsdesTroie",
		"sqlnet1",
		"sqlnet2-1521",
		"sqlnet2-1525",
		"sqlnet2-1526",
		"Squid_NTLM",
		"ssh",
		"ssh_version_2",
		"ssl_v2",
		"ssl_v3",
		"StoneBeat-Control",
		"StoneBeat-Daemon",
		"SubSeven",
		"SubSeven-G",
		"T.120",
		"TACACSplus",
		"tcp-high-ports",
		"telnet",
		"Terrortrojan",
		"TheFlu",
		"time-tcp",
		"tls1.0",
		"tls1.1",
		"tls1.2",
		"tns",
		"TransScout",
		"Trinoo",
		"UltorsTrojan",
		"unknown_protocol_tcp",
		"UserCheck",
		"uucp",
		"wais",
		"winframe",
		"WinHole",
		"X11",
		"Xanadu",
		"Yahoo_Messenger_messages",
		"Yahoo_Messenger_Voice_Chat_TCP",
		"Yahoo_Messenger_Webcams",
		// Check Point の事前定義済みサービス (udp)
		"archie",
		"BFD-Multihop",
		"BFD-Single_hop",
		"biff",
		"Blubster",
		"bootp",
		"Citrix_ICA_Browsing",
		"CP_SecureAgent-udp",
		"CU-SeeMe",
		"daytime-udp",
		"dhcp",
		"dhcp-relay",
		"dhcp-rep-localmodule",
		"dhcp-req-localmodule",
		"Direct_Connect_UDP",
		"discard-udp",
		"domain-udp",
		"E2ECP",
		"echo-udp",
		"eDonkey_4665",
		"FreeTel-outgoing-server",
		"FW1_load_agent",
		"FW1_scv_keep_alive",
		"FW1_snmp",
		"GNUtella_rtr_UDP",
		"GNUtella_UDP",
		"H323_ras",
		"H323_ras_only",
		"HackaTack_31789",
		"HackaTack_31791",
		"Hotline_tracker",
		"ICQ_locator",
		"IKE",
		"IKE_NAT_TRAVERSAL",
		"interphone",
		"kerberos-udp",
		"Kerberos_v5_UDP",
		"L2TP",
		"ldap_udp",
		"MetaIP-UAT",
		"mgcp_CA",
		"mgcp_MG",
		"microsoft-ds-udp",
		"MS-SQL-Monitor_UDP",
		"MS-SQL-Server_UDP",
		"MSN_Messenger_1863_UDP",
		"MSN_Messenger_5190",
		"MSN_Messenger_Voice",
		"MSSQL_resolver",
		"name",
		"nbdatagram",
		"nbname",
		"NEW-RADIUS",
		"NEW-RADIUS-ACCOUNTING",
		"nfsd",
		"NoBackO",
		"ntp-udp",
		"pcANYWHERE-stat",
		"quic",
		"RADIUS",
		"RADIUS-ACCOUNTING",
		"RainWall_Daemon",
		"RainWall_Status",
		"RainWall_Stop",
		"RDP",
		"Remote_Desktop_Protocol_UDP",
		"RexxRave",
		"rip",
		"RIPng",
		"rtp",
		"securid-udp",
		"sip",
		"sip_any",
		"smb-udp",
		"snmp",
		"snmp-read",
		"snmp-trap",
		"SWTP_Gateway",
		"SWTP_SMS",
		"syslog",
		"TACACS",
		"tftp",
		"time-udp",
		"tunnel_test",
		"UA_CS",
		"UA_PHONE",
		"udp-high-ports",
		"unknown_protocol_udp",
		"VPN1_IPSEC_encapsulation",
		"wap_wdp",
		"wap_wdp_enc",
		"wap_wtp",
		"wap_wtp_enc",
		"who",
		"WinMX",
		"Yahoo_Messenger_Voice_Chat_UDP",
		// Check Point の事前定義済みサービス (icmp)
		"dest-unreach",
		"echo-reply",
		"echo-request",
		"info-reply",
		"info-req",
		"mask-reply",
		"mask-request",
		"param-prblm",
		"redirect",
		"source-quench",
		"time-exceeded",
		"timestamp",
		"timestamp-reply",
		// Check Point の事前定義済みサービス (other)
		"AH",
		"backweb",
		"dhcp-reply",
		"dhcp-request",
		"dhcpv6-relay",
		"dhcpv6-reply",
		"dhcpv6-request",
		"egp",
		"ESP",
		"FreeTel-incoming",
		"FreeTel-outgoing-client",
		"ftp_mapped",
		"FW1_Encapsulation",
		"ggp",
		"gre",
		"gtp_v0_path_mgmt",
		"gtp_v1_path_mgmt",
		"gtp_v2_path_mgmt",
		"high_udp_for_secure_SCCP",
		"http_mapped",
		"HTTP_wo_SCV",
		"icmp-proto",
		"igmp",
		"igrp",
		"MGCP_dynamic_ports",
		"Mobility_Header",
		"ospf",
		"Partial-Destination-v6",
		"Partial-Source-v6",
		"pim",
		"rip-response",
		"sip_dynamic_ports",
		"SIT",
		"SIT_with_Intra_Tunnel_Inspection",
		"Sitara",
		"SKIP",
		"smtp_mapped",
		"Snmp-Read-Only",
		"traceroute",
		"tunnel_test_mapped",
		"vrrp",
		"X11-verify",
		"ZSP",
		// Check Point の事前定義済みサービス グループ
		"AD_Dcerpc_services",
		"AOL_Messenger",
		"Authenticated",
		"CIFS",
		"Citrix_metaFrame",
		"DAIP_Control_services",
		"daytime",
		"Direct_Connect",
		"discard",
		"dns",
		"echo",
		"eDonkey",
		"Entrust-CA",
		"FreeTel-outgoing",
		"FW1_clntauth",
		"GNUtella",
		"Hotline",
		"HTTPS default services",
		"icmp-requests",
		"Integrity_Server",
		"IPSEC",
		"IPv6_group",
		"irc",
		"kerberos",
		"Messenger_Applications",
		"MS-SQL",
		"MSExchange",
		"MSExchange-2000",
		"MSExchange-RemoteAdmin",
		"MSExchange-SiteConnector",
		"MSExchange_2007",
		"MSN_Messenger",
		"Napster",
		"NBT",
		"NetMeeting",
		"NFS",
		"NIS",
		"ntp",
		"OAS",
		"Orbix",
		"P2P_File_Sharing_Applications",
		"pcANYWHERE",
		"pcTELECOMMUTE",
		"PPTP",
		"RainWall-Control",
		"RealPlayer",
		"securid",
		"sqlnet2",
		"StoneBeat",
		"time",
		"Trojan_Services",
		"Web",
		"Web_Proxy",
		"Yahoo_Messenger",
	}
	// TODO:
	cpCli.ZoneReplacementMap.Description = "ゾーン名置換用マップ"
	cpCli.ZoneReplacementMap.Value = []ZoneMap{
		{Before: "LAN", After: "LANZone"},
		{Before: "WAN", After: "WANZone"},
		{Before: "DMZ", After: "DMZZone"},
	}
	cpCli.AddressReplacementMap.Description = "アドレス名置換用マップ"
	cpCli.AddressReplacementMap.Value = []ServiceMap{}
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
