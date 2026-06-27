package model

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/pelletier/go-toml/v2"
)

func TestNewDefaultAppConfig(t *testing.T) {
	appConfig := NewDefaultAppConfig()

	jsonData := []byte(`{
  "_description": "utmconv の設定ファイル",
  "checkpoint": {
    "cli": {
      "mgmt_cli_user": {
        "_description": "mgmt_cli tool のユーザー名",
        "value": "secadmin"
      },
      "mgmt_cli_password": {
        "_description": "mgmt_cli tool のパスワード",
        "value": "Lab@12345"
      },
      "ignore-warnings": {
        "_description": "mgmt_cli tool に ignore-warnings true を付加するかどうか？",
        "value": true
      },
      "access_rule_layer": {
        "_description": "cli 出力時の access-rule の layer",
        "value": "Network"
      },
      "access_rule_section": {
        "_description": "cli 出力時の access-rule を追加するセクション タイトル",
        "value": "New rules"
      },
      "nat_rule_package": {
        "_description": "cli 出力時の nat-rule の package",
        "value": "standard"
      },
      "threat_rule_layer": {
        "_description": "cli 出力時の threat-rule の layer",
        "value": "Threat Prevention"
      },
      "predefined_services": {
        "_description": "事前定義サービス名・サービス グループ名の配列",
        "value": [
          "ICMP Protocol",
          "Instagram",
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
          "Yahoo_Messenger"
        ]
      },
      "zone_replacement_map": {
        "_description": "ゾーン名置換用マップ",
        "value": [
          {
            "before": "LAN",
            "after": "LANZone"
          },
          {
            "before": "WAN",
            "after": "WANZone"
          },
          {
            "before": "DMZ",
            "after": "DMZZone"
          }
        ]
      },
      "address_replacement_map": {
        "_description": "アドレス名置換用マップ",
        "value": []
      },
      "service_replacement_map": {
        "_description": "サービス名置換用マップ",
        "value": [
          {
            "before": "service-http",
            "after": "http"
          },
          {
            "before": "service-https",
            "after": "https"
          }
        ]
      }
    }
  },
  "paloalto": {
    "conf": {
      "application_default_replacement_map": {
        "_description": "サービスが application-default または any のときのサービス名置換用マップ",
        "value": [
          {
            "application": "icmp",
            "services": [
              "ICMP Protocol"
            ]
          },
          {
            "application": "ping",
            "services": [
              "echo-request"
            ]
          },
          {
            "application": "traceroute",
            "services": [
              "traceroute"
            ]
          },
          {
            "application": "ssh",
            "services": [
              "ssh"
            ]
          },
          {
            "application": "ssh-tunnel",
            "services": [
              "ssh"
            ]
          },
          {
            "application": "syslog",
            "services": [
              "syslog"
            ]
          },
          {
            "application": "ipsec-esp",
            "services": [
              "ESP",
              "IKE",
              "IKE_NAT_TRAVERSAL"
            ]
          },
          {
            "application": "ipsec-esp-udp",
            "services": [
              "ESP",
              "IKE",
              "IKE_NAT_TRAVERSAL"
            ]
          },
          {
            "application": "instagram-base",
            "services": [
              "Instagram"
            ]
          }
        ]
      }
    }
  }
}`)

	tomlData := []byte(`# Check Point
[checkpoint]
# Check Point の CLI
[checkpoint.cli]
# mgmt_cli tool のユーザー名
[checkpoint.cli.mgmt_cli_user]
value = 'secadmin'

# mgmt_cli tool のパスワード
[checkpoint.cli.mgmt_cli_password]
value = 'Lab@12345'

# mgmt_cli tool に ignore-warnings true を付加するかどうか？
[checkpoint.cli.ignore-warnings]
value = true

# cli 出力時の access-rule の layer
[checkpoint.cli.access_rule_layer]
value = 'Network'

# cli 出力時の access-rule を追加するセクション タイトル
[checkpoint.cli.access_rule_section]
value = 'New rules'

# cli 出力時の nat-rule の package
[checkpoint.cli.nat_rule_package]
value = 'standard'

# cli 出力時の threat-rule の layer
[checkpoint.cli.threat_rule_layer]
value = 'Threat Prevention'

# 事前定義サービス名・サービス グループ名の配列
[checkpoint.cli.predefined_services]
value = [
  'ICMP Protocol',
  'Instagram',
  'AOL',
  'AP-Defender',
  'AT-Defender',
  'Backage',
  'BGP',
  'Bionet-Setup',
  'CheckPointExchangeAgent',
  'Citrix_ICA',
  'ConnectedOnLine',
  'CP_Exnet_PK',
  'CP_Exnet_resolve',
  'CP_redundant',
  'CP_reporting',
  'CP_rtm',
  'CP_seam',
  'CP_SmartPortal',
  'CP_SSL_Network_Extender',
  'CPD',
  'CPD_amon',
  'CPM',
  'CPMI',
  'CrackDown',
  'CreativePartnerClnt',
  'CreativePartnerSrvr',
  'DaCryptic',
  'DameWare',
  'daytime-tcp',
  'DerSphere',
  'DerSphere_II',
  'Direct_Connect_TCP',
  'discard-tcp',
  'DNP3',
  'domain-tcp',
  'DoT',
  'echo-tcp',
  'EDGE',
  'eDonkey_4661',
  'eDonkey_4662',
  'Entrust-Admin',
  'Entrust-KeyMgmt',
  'exec',
  'FIBMGR',
  'finger',
  'Freak2k',
  'ftp',
  'ftp-bidir',
  'ftp-pasv',
  'ftp-port',
  'FW1',
  'FW1_amon',
  'FW1_clntauth_http',
  'FW1_clntauth_telnet',
  'FW1_CPRID',
  'FW1_cvp',
  'FW1_ela',
  'FW1_ica_mgmt_tools',
  'FW1_ica_pull',
  'FW1_ica_push',
  'FW1_ica_services',
  'FW1_key',
  'FW1_lea',
  'FW1_log',
  'FW1_mgmt',
  'FW1_netso',
  'FW1_omi',
  'FW1_omi-sic',
  'FW1_pslogon',
  'FW1_pslogon_NG',
  'FW1_sam',
  'FW1_sds_logon',
  'FW1_sds_logon_NG',
  'FW1_snauth',
  'FW1_topo',
  'FW1_uaa',
  'FW1_ufp',
  'GateCrasher',
  'GNUtella_rtr_TCP',
  'GNUtella_TCP',
  'gopher',
  'GoToMyPC',
  'H323',
  'H323_any',
  'HackaTack_31785',
  'HackaTack_31787',
  'HackaTack_31788',
  'HackaTack_31790',
  'HackaTack_31792',
  'Hotline_client',
  'http',
  'HTTP_and_HTTPS_proxy',
  'HTTP_proxy',
  'https',
  'HTTPS_proxy',
  'icap',
  'ICCP',
  'ICKiller',
  'ident',
  'IKE_NAT_TRAVERSAL_TCP',
  'IKE_tcp',
  'imap',
  'IMAP-SSL',
  'iMesh',
  'InCommand',
  'indeni',
  'IPSO_Clustering_Mgmt_Protocol',
  'irc1',
  'irc2',
  'jabber',
  'Jade',
  'Kaos',
  'KaZaA',
  'Kerberos_v5_TCP',
  'Kuang2',
  'ldap',
  'ldap-ssl',
  'login',
  'lotus',
  'lpdw0rm',
  'Madster',
  'microsoft-ds',
  'Mneah',
  'Modbus',
  'MS-SQL-Monitor',
  'MS-SQL-Server',
  'MSN_Messenger_File_Transfer',
  'MSNMS',
  'MSNP',
  'Multidropper',
  'MySQL',
  'Napster_Client_6600-6699',
  'Napster_directory_4444',
  'Napster_directory_5555',
  'Napster_directory_6666',
  'Napster_directory_7777',
  'Napster_directory_8888_primary',
  'Napster_redirector',
  'nbsession',
  'NCP',
  'netbios-session',
  'netshow',
  'netstat',
  'nfsd-tcp',
  'nntp',
  'ntp-tcp',
  'OAS-NameServer',
  'OAS-ORB',
  'OPC',
  'OpenWindows',
  'Orbix-1570',
  'Orbix-1571',
  'pcANYWHERE-data',
  'pcTELECOMMUTE-FileSync',
  'pop-2',
  'pop-3',
  'POP3S',
  'Port_6667_trojans',
  'PostgreSQL',
  'pptp-tcp',
  'RainWall_Command',
  'RAT',
  'Real-Audio',
  'RealSecure',
  'Remote_Desktop_Protocol',
  'Remote_Storm',
  'rtsp',
  'SCCP',
  'securidprop',
  'Shadyshell',
  'shell',
  'SIC-TCP',
  'sip-tcp',
  'sip_any-tcp',
  'sip_tls_authentication',
  'sip_tls_not_inspected',
  'SkyDance-T',
  'smb',
  'smtp',
  'SMTPS',
  'snmp-tcp',
  'SocketsdesTroie',
  'sqlnet1',
  'sqlnet2-1521',
  'sqlnet2-1525',
  'sqlnet2-1526',
  'Squid_NTLM',
  'ssh',
  'ssh_version_2',
  'ssl_v2',
  'ssl_v3',
  'StoneBeat-Control',
  'StoneBeat-Daemon',
  'SubSeven',
  'SubSeven-G',
  'T.120',
  'TACACSplus',
  'tcp-high-ports',
  'telnet',
  'Terrortrojan',
  'TheFlu',
  'time-tcp',
  'tls1.0',
  'tls1.1',
  'tls1.2',
  'tns',
  'TransScout',
  'Trinoo',
  'UltorsTrojan',
  'unknown_protocol_tcp',
  'UserCheck',
  'uucp',
  'wais',
  'winframe',
  'WinHole',
  'X11',
  'Xanadu',
  'Yahoo_Messenger_messages',
  'Yahoo_Messenger_Voice_Chat_TCP',
  'Yahoo_Messenger_Webcams',
  'archie',
  'BFD-Multihop',
  'BFD-Single_hop',
  'biff',
  'Blubster',
  'bootp',
  'Citrix_ICA_Browsing',
  'CP_SecureAgent-udp',
  'CU-SeeMe',
  'daytime-udp',
  'dhcp',
  'dhcp-relay',
  'dhcp-rep-localmodule',
  'dhcp-req-localmodule',
  'Direct_Connect_UDP',
  'discard-udp',
  'domain-udp',
  'E2ECP',
  'echo-udp',
  'eDonkey_4665',
  'FreeTel-outgoing-server',
  'FW1_load_agent',
  'FW1_scv_keep_alive',
  'FW1_snmp',
  'GNUtella_rtr_UDP',
  'GNUtella_UDP',
  'H323_ras',
  'H323_ras_only',
  'HackaTack_31789',
  'HackaTack_31791',
  'Hotline_tracker',
  'ICQ_locator',
  'IKE',
  'IKE_NAT_TRAVERSAL',
  'interphone',
  'kerberos-udp',
  'Kerberos_v5_UDP',
  'L2TP',
  'ldap_udp',
  'MetaIP-UAT',
  'mgcp_CA',
  'mgcp_MG',
  'microsoft-ds-udp',
  'MS-SQL-Monitor_UDP',
  'MS-SQL-Server_UDP',
  'MSN_Messenger_1863_UDP',
  'MSN_Messenger_5190',
  'MSN_Messenger_Voice',
  'MSSQL_resolver',
  'name',
  'nbdatagram',
  'nbname',
  'NEW-RADIUS',
  'NEW-RADIUS-ACCOUNTING',
  'nfsd',
  'NoBackO',
  'ntp-udp',
  'pcANYWHERE-stat',
  'quic',
  'RADIUS',
  'RADIUS-ACCOUNTING',
  'RainWall_Daemon',
  'RainWall_Status',
  'RainWall_Stop',
  'RDP',
  'Remote_Desktop_Protocol_UDP',
  'RexxRave',
  'rip',
  'RIPng',
  'rtp',
  'securid-udp',
  'sip',
  'sip_any',
  'smb-udp',
  'snmp',
  'snmp-read',
  'snmp-trap',
  'SWTP_Gateway',
  'SWTP_SMS',
  'syslog',
  'TACACS',
  'tftp',
  'time-udp',
  'tunnel_test',
  'UA_CS',
  'UA_PHONE',
  'udp-high-ports',
  'unknown_protocol_udp',
  'VPN1_IPSEC_encapsulation',
  'wap_wdp',
  'wap_wdp_enc',
  'wap_wtp',
  'wap_wtp_enc',
  'who',
  'WinMX',
  'Yahoo_Messenger_Voice_Chat_UDP',
  'dest-unreach',
  'echo-reply',
  'echo-request',
  'info-reply',
  'info-req',
  'mask-reply',
  'mask-request',
  'param-prblm',
  'redirect',
  'source-quench',
  'time-exceeded',
  'timestamp',
  'timestamp-reply',
  'AH',
  'backweb',
  'dhcp-reply',
  'dhcp-request',
  'dhcpv6-relay',
  'dhcpv6-reply',
  'dhcpv6-request',
  'egp',
  'ESP',
  'FreeTel-incoming',
  'FreeTel-outgoing-client',
  'ftp_mapped',
  'FW1_Encapsulation',
  'ggp',
  'gre',
  'gtp_v0_path_mgmt',
  'gtp_v1_path_mgmt',
  'gtp_v2_path_mgmt',
  'high_udp_for_secure_SCCP',
  'http_mapped',
  'HTTP_wo_SCV',
  'icmp-proto',
  'igmp',
  'igrp',
  'MGCP_dynamic_ports',
  'Mobility_Header',
  'ospf',
  'Partial-Destination-v6',
  'Partial-Source-v6',
  'pim',
  'rip-response',
  'sip_dynamic_ports',
  'SIT',
  'SIT_with_Intra_Tunnel_Inspection',
  'Sitara',
  'SKIP',
  'smtp_mapped',
  'Snmp-Read-Only',
  'traceroute',
  'tunnel_test_mapped',
  'vrrp',
  'X11-verify',
  'ZSP',
  'AD_Dcerpc_services',
  'AOL_Messenger',
  'Authenticated',
  'CIFS',
  'Citrix_metaFrame',
  'DAIP_Control_services',
  'daytime',
  'Direct_Connect',
  'discard',
  'dns',
  'echo',
  'eDonkey',
  'Entrust-CA',
  'FreeTel-outgoing',
  'FW1_clntauth',
  'GNUtella',
  'Hotline',
  'HTTPS default services',
  'icmp-requests',
  'Integrity_Server',
  'IPSEC',
  'IPv6_group',
  'irc',
  'kerberos',
  'Messenger_Applications',
  'MS-SQL',
  'MSExchange',
  'MSExchange-2000',
  'MSExchange-RemoteAdmin',
  'MSExchange-SiteConnector',
  'MSExchange_2007',
  'MSN_Messenger',
  'Napster',
  'NBT',
  'NetMeeting',
  'NFS',
  'NIS',
  'ntp',
  'OAS',
  'Orbix',
  'P2P_File_Sharing_Applications',
  'pcANYWHERE',
  'pcTELECOMMUTE',
  'PPTP',
  'RainWall-Control',
  'RealPlayer',
  'securid',
  'sqlnet2',
  'StoneBeat',
  'time',
  'Trojan_Services',
  'Web',
  'Web_Proxy',
  'Yahoo_Messenger'
]

# ゾーン名置換用マップ
[checkpoint.cli.zone_replacement_map]
[[checkpoint.cli.zone_replacement_map.value]]
before = 'LAN'
after = 'LANZone'

[[checkpoint.cli.zone_replacement_map.value]]
before = 'WAN'
after = 'WANZone'

[[checkpoint.cli.zone_replacement_map.value]]
before = 'DMZ'
after = 'DMZZone'

# アドレス名置換用マップ
[checkpoint.cli.address_replacement_map]
value = []

# サービス名置換用マップ
[checkpoint.cli.service_replacement_map]
[[checkpoint.cli.service_replacement_map.value]]
before = 'service-http'
after = 'http'

[[checkpoint.cli.service_replacement_map.value]]
before = 'service-https'
after = 'https'

# PaloAlto
[paloalto]
# PaloAlto の Conf
[paloalto.conf]
# サービスが application-default または any のときのサービス名置換用マップ
[paloalto.conf.application_default_replacement_map]
[[paloalto.conf.application_default_replacement_map.value]]
application = 'icmp'
services = ['ICMP Protocol']

[[paloalto.conf.application_default_replacement_map.value]]
application = 'ping'
services = ['echo-request']

[[paloalto.conf.application_default_replacement_map.value]]
application = 'traceroute'
services = ['traceroute']

[[paloalto.conf.application_default_replacement_map.value]]
application = 'ssh'
services = ['ssh']

[[paloalto.conf.application_default_replacement_map.value]]
application = 'ssh-tunnel'
services = ['ssh']

[[paloalto.conf.application_default_replacement_map.value]]
application = 'syslog'
services = ['syslog']

[[paloalto.conf.application_default_replacement_map.value]]
application = 'ipsec-esp'
services = ['ESP', 'IKE', 'IKE_NAT_TRAVERSAL']

[[paloalto.conf.application_default_replacement_map.value]]
application = 'ipsec-esp-udp'
services = ['ESP', 'IKE', 'IKE_NAT_TRAVERSAL']

[[paloalto.conf.application_default_replacement_map.value]]
application = 'instagram-base'
services = ['Instagram']`)

	appConfigJSON := &AppConfig{}
	if err := json.Unmarshal(jsonData, appConfigJSON); err != nil {
		t.Fatalf("failed to unmarshal json: %v", err)
	}
	diff := cmp.Diff(appConfig, appConfigJSON)
	if diff != "" {
		t.Errorf("appConfig diff(json): %s", diff)
	}

	appConfigTOML := &AppConfig{}
	if err := toml.Unmarshal(tomlData, appConfigTOML); err != nil {
		t.Fatalf("failed to unmarshal toml: %v", err)
	}
	diff = cmp.Diff(appConfig, appConfigTOML)
	if diff != "" {
		t.Errorf("appConfig diff(toml): %s", diff)
	}
}
