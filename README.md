# epub2audio

様々なテキスト形式（EPUB、Markdown、HTML）から高品質な音声ファイルを生成するツール。VoicevoxやVoicepeakなど複数のTTSエンジンに対応。

## 特徴

- 📚 **多様な入力形式**: EPUB、Markdown、HTMLファイルをサポート
- 🎯 **読み上げ最適化**: 読み上げに適した部分を自動抽出
- 🎵 **構造保持**: 見出しや段落構造に応じた適切なポーズを挿入
- 🔄 **複数エンジン対応**: VoicevoxとVoicepeakを切り替え可能
- ⚡ **高速処理**: 並列処理によるバッチ音声生成

## インストール

### 前提条件

- Go 1.21以上
- ffmpeg（音声ファイルの処理用）
- Voicevox または Voicepeak（TTSエンジン）

### ビルド

```bash
git clone https://github.com/sbr-nakano/epub2audio.git
cd epub2audio
go build -o epub2audio cmd/epub2audio/main.go
```

## 使い方

### 基本的な使用方法

```bash
# EPUBファイルを音声に変換
epub2audio convert book.epub -o audiobook.mp3

# Markdownファイルを変換（Voicevox使用）
epub2audio convert README.md --engine voicevox -o readme.wav

# HTMLファイルを変換（カスタム設定）
epub2audio convert article.html --config myconfig.yml
```

### 複数ファイルの処理

```bash
# ディレクトリ内の全Markdownファイルを変換
epub2audio convert *.md --output-dir ./audio/

# バッチファイルによる一括処理
epub2audio batch filelist.txt --parallel 4
```

### その他のコマンド

```bash
# 利用可能な話者の一覧表示
epub2audio list speakers --engine voicevox

# テキスト抽出のみ（ドライラン）
epub2audio convert book.epub --dry-run --extract-text extracted.txt

# ヘルプの表示
epub2audio --help
```

## 設定

### 設定ファイル例（config.yml）

```yaml
# 入力設定
input:
  parser: auto  # auto, markdown, epub, html
  encoding: utf-8

# 処理設定
processing:
  max_chunk_size: 140
  preserve_structure: true

# ポーズ設定
pause:
  strategy: standard
  standard:
    heading_1: 2000  # ミリ秒
    heading_2: 1500
    paragraph: 800

# TTSエンジン設定
tts:
  engine: voicepeak  # voicepeak, voicevox
  
  voicepeak:
    narrator: "Japanese Female"
    speed: 1.0
    
  voicevox:
    speaker_id: 1  # ずんだもん（ノーマル）
    speed_scale: 1.0

# 出力設定
output:
  format: mp3
  bitrate: 128k
  merge: true
```

## 対応フォーマット

### 入力形式
- EPUB (.epub)
- Markdown (.md)
- HTML (.html, .htm)

### 出力形式
- MP3（推奨）
- WAV
- M4A

## TTSエンジン

### Voicepeak
- 商用利用可能なTTSエンジン
- 高品質な日本語音声
- ポーズは全角スペースで制御

### Voicevox
- 無料のオープンソースTTSエンジン
- 多様なキャラクターボイス
- REST APIによる制御

## アーキテクチャ

```
┌─────────────────────────────────────────────────┐
│                   CLI Layer                      │
├─────────────────────────────────────────────────┤
│              Orchestration Layer                 │
├──────────────────┬──────────────────┬───────────┤
│  Input Parser    │ Text Processor   │    TTS    │
├──────────────────┼──────────────────┼───────────┤
│ • Markdown       │ • Pause Injector │ • Voicevox│
│ • EPUB          │ • Text Splitter  │ • Voicepeak│
│ • HTML          │ • Normalizer     │           │
└──────────────────┴──────────────────┴───────────┘
```

## 開発

### ディレクトリ構成

```
epub2audio/
├── cmd/           # CLIエントリーポイント
├── internal/      # 内部パッケージ
│   ├── parser/    # 入力パーサー
│   ├── tts/       # TTSエンジン
│   └── processor/ # テキスト処理
├── pkg/           # 公開パッケージ
├── configs/       # 設定ファイル例
└── docs/          # ドキュメント
```

### テスト実行

```bash
# 全テストを実行
go test ./...

# カバレッジ付きでテスト
go test -cover ./...

# 特定のパッケージのテスト
go test ./internal/parser/epub
```

### ビルド

```bash
# リリースビルド
make build

# クロスコンパイル
make build-all
```

## トラブルシューティング

### Voicepeakが動作しない
- Voicepeakがインストールされているか確認
- 実行パスが正しく設定されているか確認
- ライセンスが有効か確認

### Voicevoxに接続できない
- Voicevoxが起動しているか確認（デフォルト: http://localhost:50021）
- ファイアウォールの設定を確認

### 音声が途切れる
- max_chunk_sizeを調整（デフォルト: 140文字）
- ポーズ設定を調整

## ライセンス

MIT License

## 貢献

プルリクエストを歓迎します。大きな変更の場合は、まずissueを作成して変更内容を議論してください。

## 作者

nakano shota (@sbr-nakano)
