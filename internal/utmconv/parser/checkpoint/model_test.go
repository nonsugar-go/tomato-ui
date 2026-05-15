package checkpoint

import (
	"encoding/json"
	"testing"
)

func TestCheckPointObjectUnmarshal(t *testing.T) {
	rawJSON := `[
  {
    "nat-settings": {
      "ipv4-address": "192.0.2.5",
      "method": "static",
      "install-on": "CPGW1",
      "auto-rule": true,
      "ipv6-address": ""
    },
    "interfaces": [],
    "comments": "",
    "color": "orange",
    "icon": "Objects/host",
    "meta-info": {
      "creator": "secadmin",
      "validation-state": "ok",
      "last-modify-time": {
        "iso-8601": "2026-05-02T15:37+0900",
        "posix": 1777703842647
      },
      "creation-time": {
        "iso-8601": "2026-05-02T10:42+0900",
        "posix": 1777686147131
      },
      "lock": "unlocked",
      "last-modifier": "secadmin"
    },
    "type": "host",
    "tags": [],
    "uid": "6928aae5-878f-425c-942c-cc2ca36da718",
    "ipv4-address": "192.168.111.5",
    "domain": {
      "uid": "41e821a0-3720-11e3-aa6e-0800200c9fde",
      "domain-type": "domain",
      "name": "SMC User"
    },
    "name": "DMZSRV",
    "read-only": true,
    "available-actions": {
      "edit": "false",
      "clone": "false",
      "delete": "false"
    }
  },
  {
    "nat-settings": {
      "ipv4-address": "192.0.2.20",
      "method": "hide",
      "hide-behind": "ip-address",
      "install-on": "CPGW1",
      "auto-rule": true,
      "ipv6-address": ""
    },
    "interfaces": [],
    "comments": "",
    "color": "crete blue",
    "icon": "Objects/host",
    "meta-info": {
      "creator": "secadmin",
      "validation-state": "ok",
      "last-modify-time": {
        "iso-8601": "2026-05-02T21:45+0900",
        "posix": 1777725948705
      },
      "creation-time": {
        "iso-8601": "2026-05-02T10:11+0900",
        "posix": 1777684319815
      },
      "lock": "unlocked",
      "last-modifier": "secadmin"
    },
    "type": "host",
    "tags": [],
    "uid": "7efff40a-d27d-4789-bd7c-6088511f620c",
    "ipv4-address": "192.168.1.69",
    "domain": {
      "uid": "41e821a0-3720-11e3-aa6e-0800200c9fde",
      "domain-type": "domain",
      "name": "SMC User"
    },
    "name": "SmartConsole_VM",
    "read-only": true,
    "available-actions": {
      "edit": "false",
      "clone": "false",
      "delete": "false"
    }
  }
]
`

	var hosts []CheckPointHost
	err := json.Unmarshal([]byte(rawJSON), &hosts)
	if err != nil {
		t.Fatalf("Unmarshalに失敗しました: %v", err)
	}
	if len(hosts) != 2 {
		t.Errorf("期待するオブジェクト数は 2 ですが、%d でした", len(hosts))
	}

	target := hosts[0]
	if target.Name != "DMZSRV" {
		t.Errorf("名前が一致しません: 期待=CPGW1, 結果=%s", target.Name)
	}

	if target.Ipv4Address != "192.168.111.5" {
		t.Errorf("IPが一致しません: 期待=192.168.111.5, 結果=%s", target.Ipv4Address)
	}

	if target.NatSettings.Method != "static" {
		t.Errorf("NATメソッドが一致しません: 期待=static, 結果=%s", target.NatSettings.Method)
	}
}
