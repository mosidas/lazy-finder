# lazy-finder

macOS の Finder（カラムビュー）をベースに、[lazygit](https://github.com/jesseduffield/lazygit) 風の枠線UIでファイル操作ができる TUI ファイラです。Go + [Bubble Tea](https://github.com/charmbracelet/bubbletea) 製。

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

## 特長

- **3カラムのカラムビュー**: 親ディレクトリ / カレント / プレビュー（フォルダは中身、ファイルはテキストプレビュー）
- **lazygit スタイル + ayu light テーマ**: 角丸枠線・フォーカスパネルのハイライト・下部ステータス/ヘルプバー。明るい背景に青（フォルダ）/オレンジ（アクセント）の ayu light 配色（端末の背景色に依存せず明色で描画）
- **基本操作**: コピー / 切り取り / 貼り付け / 名前変更 / 新規フォルダ / 削除（ゴミ箱へ移動）
- **複数選択**: space でマークして一括操作
- **VS Code 連携**: `e` でカレント（またはカーソル上のフォルダ）を `code` で開く
- **安全な削除**: 完全削除ではなく macOS の `~/.Trash` へ移動（Linux では freedesktop のゴミ箱）

## インストール / 実行

```sh
make build          # ./bin/lazy-finder を生成
./bin/lazy-finder            # カレントディレクトリで起動
./bin/lazy-finder ~/Projects # 任意のディレクトリで起動

# もしくは直接
go run . [起動ディレクトリ]
go install .        # $GOBIN にインストール
```

依存: Go 1.25+。VS Code 連携には `code` コマンドが PATH に必要です（VS Code の「シェルコマンド: PATH内に'code'コマンドをインストール」）。

## キーバインド

| キー | 動作 |
| --- | --- |
| `↑`/`k`, `↓`/`j` | カーソル移動 |
| `g` / `G` | 先頭 / 末尾へ |
| `enter` / `l` / `→` | フォルダに入る / ファイルを標準アプリで開く |
| `h` / `←` / `backspace` | 親フォルダへ |
| `~` | ホームディレクトリへ |
| `space` | 選択（マーク）のトグル |
| `y` / `x` | コピー / 切り取り |
| `p` | 貼り付け（カレントへ） |
| `d` | ゴミ箱へ移動（確認あり） |
| `r` | 名前変更 |
| `n` | 新規フォルダ |
| `o` | 標準アプリで開く |
| `e` | VS Code で開く |
| `.` | 隠しファイルの表示切替 |
| `R` / `Ctrl+R` | 再読込 |
| `?` | ヘルプ |
| `q` / `Ctrl+C` | 終了 |

選択（マーク）がある場合は `y`/`x`/`d` はマークした全項目が対象、無ければカーソル上の項目が対象になります。

## 構成

```
main.go                  エントリポイント
internal/fsops/          ファイルシステム操作（一覧 / コピー / 移動 / ゴミ箱 / プレビュー）
internal/ui/             Bubble Tea の Model / Update / View・キーバインド・スタイル
```

## 開発

```sh
make test    # ユニットテスト
make vet     # go vet
make fmt     # gofmt
```
