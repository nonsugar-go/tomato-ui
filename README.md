# Tomato User Interface (TUI)

🍅 UI

## TODO

- [X] utmconv: UTM コンフィグ解析とベンダー間変換ツール
- [ ] vm-tui: VMware ESXi の管理ツール

### utmconv

#### 対象

- [X] PaloAlto (Panorama 含む)
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

- [X] panorama.xlsx: Panorama 解析結果の Excel 出力
- [X] checkpoint_address.conf: Check Point 用のコンフィグ (host / network)
- [X] checkpoint_address_group.conf: Check Point 用のコンフィグ (address-group)
- [ ] checkpoint_service.conf: Check Point 用のコンフィグ (service)
- [ ] checkpoint_policy.conf: Check Point 用のコンフィグ (policy)
- [ ] checkpoint_nat.conf: Check Point 用のコンフィグ (nat)

### vm-tui

- [ ] 仮想マシンの一覧を表示する (名前, ステータス, IP アドレス)
- [ ] 仮想マシンを起動・シャットダウン・停止
