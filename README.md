# Tomato User Interface (TUI)

🍅 UI

## TODO

- [ ] utmconv: UTM コンフィグ解析とベンダー間変換ツール
- [ ] vm-tui: VMware ESXi の管理ツール

### utmconv

#### 対象

- [ ] PaloAlto (Panorama 含む)
- [ ] FortiGate
- [ ] Check Point

#### 使用方法

```bash
## Panorama のコンフィグから解析結果を Excel へ出力
$ utmconv -vendor paloalto -in panorama.xml

## Panorama のコンフィグから Check Point 用のコンフィグへ変換
$ utmconv -vendor paloalto -to checkpoint
```

#### 出力ファイル (例)

- panorama.xlsx: Panorama 解析結果の Excel 出力
- checkpoint_host_DG.conf: Check Point 用のコンフィグ (host)
- checkpoint_network_DG.conf: Check Point 用のコンフィグ (network)
- checkpoint_service_DG.conf: Check Point 用のコンフィグ (service)
- checkpoint_nat_DG.conf: Check Point 用のコンフィグ (nat)
- checkpoint_policy_DG.conf: Check Point 用のコンフィグ (policy)

### vm-tui

- 仮想マシンの一覧を表示する (名前, ステータス, IP アドレス)
- 仮想マシンを起動・シャットダウン・停止