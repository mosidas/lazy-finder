# lazy-finder

macOS の Finder のカラムビューを、[lazygit](https://github.com/jesseduffield/lazygit) 風の枠線の画面で操作するターミナル用のファイラー。

```
┌─親──────┐┌─Documents──────────┐┌─プレビュー──────────┐
│ Desktop ││ ● report.md        ││ # 2024 Report       │
│ Documen…││   notes.txt        ││                     │
│ Downloa…││   img/             ││ ## Summary          │
└─────────┘└────────────────────┘└─────────────────────┘
 lazy-finder  ~/Documents
 ✓ 1 件をコピーしました                    [copy 1 件]
 enter/l:開く  h:親へ  space:選択  y:コピー  p:貼り付け  d:ゴミ箱へ …
```

## 動作環境

- macOS(Apple Silicon と Intel)、Linux(amd64 と arm64)
- `code` コマンド。フォルダーを VS Code で開く操作に使う

## インストール

### Releases から

1. [Releases](https://github.com/mosidas/lazy-finder/releases) から `lazy-finder-<バージョン>-<os>-<arch>.tar.gz` をダウンロードして展開する。Apple Silicon の Mac は `darwin-arm64`、Intel の Mac は `darwin-amd64`、Linux は `linux-amd64` か `linux-arm64` である。
2. `lazy-finder` を PATH の通ったディレクトリ(例: `/usr/local/bin`)へ移す。
3. macOS では、配布物を公証していないため初回の実行が Gatekeeper に止められる。ターミナルで `xattr -d com.apple.quarantine /usr/local/bin/lazy-finder` を実行し、隔離属性を外す。
4. `lazy-finder --help` を実行する。

### ソースからビルドする

1. Go 1.25 以上をインストールする。
2. `make build` を実行する。`./bin/lazy-finder` ができる。
3. `./bin/lazy-finder --help` を実行する。

## 使い方

1. `lazy-finder <ディレクトリ>` を実行する。ディレクトリを省くと、カレントディレクトリで開く。
2. `j` と `k` でカーソルを動かす。
3. `l` でフォルダーに入る。`h` で親フォルダーへ戻る。
4. `space` でファイルを選択し、`y` でコピーする。
5. コピー先のフォルダーへ移動し、`p` で貼り付ける。
6. `q` で終了する。

キーの一覧は `?` で表示する。`d` の削除は、ファイルをゴミ箱へ移動する。

テーマを変えるには、[lazy-finder のオプションとキー](docs/reference.md) の `--theme` を指定する。

## 開発

道具とビルドは「ソースからビルドする」の手順と同じである。テストは次のコマンドで実行する。

```sh
make test
```

## ライセンス

[MIT License](LICENSE)
