#!/usr/bin/env python3
import argparse
import pexpect
import sys

def run_ssh(host: str, user: str, password: str, commands: list, prompt: str, debug: bool = False):
    p = pexpect.spawn(f"ssh {user}@{host}", encoding="utf-8", timeout=10, searchwindowsize=65536)
    p.setwinsize(100, 1000)
    p.logfile_read = sys.stdout
    p.expect(r"[Pp]assword:[ ]?$")
    p.sendline(password)
    for cmd in commands:
        p.expect(prompt, timeout=60)
        if debug:
            print(f"\n[*] Executing: {cmd}")
        p.sendline(cmd)
    p.expect(prompt, timeout=60)
    p.sendline("exit")
    p.expect(pexpect.EOF)

def main():
    parser = argparse.ArgumentParser(description="SSH 経由でコマンドを実行する CLI ツール")

    parser.add_argument("-host", required=True, help="接続先ホストの IP アドレス")
    parser.add_argument("-user", required=True, help="SSH ユーザー名")
    # Pythonの予約語 'pass' を避けるため dest で変数名を指定
    parser.add_argument("-pass", dest="password", required=True, help="SSH パスワード")
    parser.add_argument("-f", dest="file", required=True, help="コマンドファイルパス")
    parser.add_argument("-prompt", default=r"[$#>][ ]?$", help="プロンプトの正規表現 (デフォルト: '[$#>][ ]?$')")
    parser.add_argument("-debug", action="store_true", help="デバッグ表示")

    args = parser.parse_args()

    if args.debug:
        print(f"[*] Target Host: {args.host}")
        print(f"[*] SSH User: {args.user}")
        print(f"[*] SSH Password: {'*' * len(args.password)}")
        print(f"[*] Command File: {args.file}")
        print(f"[*] Prompt Regex: {args.prompt}")
        print(f"[*] Debug Mode: {'ON' if args.debug else 'OFF'}")

    try:
        with open(args.file, "r", encoding="utf-8") as f:
            commands = [line.strip() for line in f if line.strip()]
            print(f"[*] Loaded {len(commands)} commands.")
            
            run_ssh(host=args.host, user=args.user, password=args.password, commands=commands, prompt=args.prompt,
                    debug=args.debug)

    except FileNotFoundError:
        print(f"Error: {args.file} が見つかりません。")
        sys.exit(1)

if __name__ == "__main__":
    main()