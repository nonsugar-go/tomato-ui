# Tomato User Interface (TUI)

🍅 UI

```plaintext
_____                          _             _   _ ___ 
|_   _|__  _ __ ___   __ _ _ __| |_ ___      | | | |_ _|
| |/ _ \| '_ ' _ \ / _' | '__| __/ _ \_____| | | || | 
| | (_) | | | | | | (_| | |  | || (_) |_____| |_| || | 
|_|\___/|_| |_| |_|\__,_|_|   \__\___/      \___/|___|
```


## TODO

- [X] utmconv: UTM コンフィグ解析とベンダー間変換ツール
- [ ] vm-tui: VMware ESXi の管理ツール
- [X] push-cli.py: SSH 経由でコマンドを実行する CLI ツール
- [ ] cp-dump: Check Point Management API を使用して、オブジェクトやポリシーを JSON 形式でエクスポートする CLI ツール。

### utmconv

#### 対象

- [X] PaloAlto (Panorama 含む)
- [ ] FortiGate
- [ ] Check Point

#### 使用方法

```bash
utmconv -in panorama.xml -to cp -ignore-warnings
```

#### 出力ファイル (例)

- [X] panorama.xlsx: Panorama 解析結果の Excel 出力
- [X] checkpoint_address.conf: Check Point 用のコンフィグ (host / network)
- [X] checkpoint_address_group.conf: Check Point 用のコンフィグ (address-group)
- [X] checkpoint_service.conf: Check Point 用のコンフィグ (service)
- [X] checkpoint_service_group.conf: Check Point 用のコンフィグ (service-group)
- [ ] checkpoint_policy.conf: Check Point 用のコンフィグ (policy)
- [ ] checkpoint_nat.conf: Check Point 用のコンフィグ (nat)

### vm-tui

- [ ] 仮想マシンの一覧を表示する (名前, ステータス, IP アドレス)
- [ ] 仮想マシンを起動・シャットダウン・停止

### push-cli.py

#### 使用例

```zsh
push-cli.py -host 192.168.1.41 -user admin -pass Lab@12345 -f checkpoint_service.conf
```

### cp-dump

Check Point Management API を使用して、オブジェクトやポリシーを JSON 形式でエクスポートする CLI ツール。

#### Features

- ホスト / ネットワーク / グループなどのオブジェクト取得
- アクセスルール（Access Control Policy）の取得
- JSON形式で標準出力またはファイル出力
- API のページネーションを自動処理（optional）
- シンプルな CLI インターフェース

#### Usage

```bash
cp-dump <resource> [options]

cp-dump hosts --server 1.2.3.4 --user admin --password xxx
cp-dump rules --layer "Network"
```

| Resource | Description           |
| -------- | --------------------- |
| hosts    | ホストオブジェクト      |
| networks | ネットワークオブジェクト |
| groups   | グループ                |
| services | サービス（tcp/udpなど） |
| rules    | アクセスルール          |
| nat      | NATルール              |
| all      | すべて取得             |

