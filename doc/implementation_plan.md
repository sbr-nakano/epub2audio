# 実装計画書

## 1. プロジェクト概要

### 1.1 目的
様々なテキスト形式（Markdown、EPUB、HTML）から、複数のTTSエンジン（Voicevox、Voicepeak）を使用して高品質な音声ファイルを生成する。

### 1.2 主要機能
- 多様な入力形式のサポート
- 文章構造を保持した読み上げ（適切なポーズ挿入）
- 複数TTSエンジンの切り替え
- バッチ処理による効率的な音声生成

## 2. 実装フェーズ

### Phase 1: 基盤構築（2週間）
1. **プロジェクト構造の整備**
   - ディレクトリ構成の変更
   - 依存関係の管理（go.mod）
   - 基本的なCLI構造

2. **共通インターフェースの実装**
   - Parser interface
   - TTS Engine interface
   - Document/Chapter/ContentBlock構造体

3. **設定管理**
   - YAMLベースの設定ファイル
   - 環境変数サポート
   - デフォルト値の定義

### Phase 2: パーサー実装（2週間）
1. **EPUBパーサー**
   - 既存コードのリファクタリング
   - インターフェースへの適合
   - テストケースの作成

2. **Markdownパーサー**
   - Goldmarkライブラリの統合
   - 読み上げ適性フィルター
   - 構造解析

3. **HTMLパーサー**
   - メインコンテンツの自動検出
   - 不要要素の除外
   - セレクターベースの抽出

### Phase 3: TTSエンジン実装（2週間）
1. **Voicepeakエンジン**
   - 既存コードのリファクタリング
   - バッチ処理の実装
   - エラーハンドリング

2. **Voicevoxエンジン**
   - REST APIクライアント
   - 音声合成クエリの生成
   - ストリーミング対応

3. **ポーズ戦略**
   - 標準ポーズ戦略
   - カスタムルールサポート
   - エンジン別の最適化

### Phase 4: 統合とオーケストレーション（1週間）
1. **パイプライン実装**
   - ファイル検出と振り分け
   - 並列処理の実装
   - 進捗レポート

2. **音声ファイル処理**
   - ffmpeg統合
   - 音声ファイルのマージ
   - フォーマット変換

3. **CLIの完成**
   - コマンドライン引数
   - ヘルプメッセージ
   - エラーメッセージ

### Phase 5: テストと最適化（1週間）
1. **単体テスト**
   - 各モジュールのテスト
   - モックの作成
   - カバレッジ測定

2. **統合テスト**
   - E2Eテストシナリオ
   - 実際のファイルでのテスト
   - パフォーマンス測定

3. **ドキュメント**
   - README.mdの更新
   - APIドキュメント
   - 使用例の追加

## 3. 技術スタック

### 3.1 言語とフレームワーク
- Go 1.21+
- cobra (CLI framework)
- viper (configuration)
- zap (logging)

### 3.2 主要ライブラリ
```go
// go.mod
module github.com/sbr-nakano/epub2audio

go 1.21

require (
    github.com/spf13/cobra v1.8.0
    github.com/spf13/viper v1.18.0
    github.com/yuin/goldmark v1.6.0
    golang.org/x/net v0.19.0
    github.com/go-resty/resty/v2 v2.11.0
    go.uber.org/zap v1.26.0
    gopkg.in/yaml.v3 v3.0.1
)
```

### 3.3 外部依存
- ffmpeg（音声処理）
- Voicevox（TTS）
- Voicepeak（TTS）

## 4. ディレクトリ構成（最終形）

```
epub2audio/
├── cmd/
│   └── epub2audio/
│       ├── main.go
│       ├── root.go
│       ├── convert.go
│       └── list.go
├── internal/
│   ├── parser/
│   │   ├── interface.go
│   │   ├── factory.go
│   │   ├── markdown/
│   │   │   ├── parser.go
│   │   │   └── parser_test.go
│   │   ├── epub/
│   │   │   ├── parser.go
│   │   │   ├── parser_test.go
│   │   │   └── toc.go
│   │   └── html/
│   │       ├── parser.go
│   │       └── parser_test.go
│   ├── tts/
│   │   ├── interface.go
│   │   ├── factory.go
│   │   ├── pause_strategy.go
│   │   ├── voicevox/
│   │   │   ├── engine.go
│   │   │   ├── engine_test.go
│   │   │   └── client.go
│   │   └── voicepeak/
│   │       ├── engine.go
│   │       └── engine_test.go
│   ├── processor/
│   │   ├── text_processor.go
│   │   ├── splitter.go
│   │   └── normalizer.go
│   ├── orchestrator/
│   │   ├── pipeline.go
│   │   ├── worker_pool.go
│   │   └── progress.go
│   └── config/
│       ├── config.go
│       └── validator.go
├── pkg/
│   ├── audio/
│   │   ├── merger.go
│   │   └── converter.go
│   └── utils/
│       ├── file.go
│       └── string.go
├── configs/
│   ├── default.yml
│   ├── voicevox.yml
│   └── voicepeak.yml
├── test/
│   ├── fixtures/
│   ├── e2e/
│   └── benchmark/
├── docs/
│   ├── README.md
│   ├── architecture.md
│   ├── api/
│   └── examples/
├── Makefile
├── go.mod
├── go.sum
└── .github/
    └── workflows/
        └── ci.yml
```

## 5. 主要な実装タスク

### 5.1 リファクタリングタスク
- [ ] 既存のepub2xhtml.goをparser/epub/に移動
- [ ] 既存のxhtml2txt.goをprocessor/に統合
- [ ] main.goをCLI構造に変更

### 5.2 新規実装タスク
- [ ] Markdownパーサー
- [ ] HTMLパーサー
- [ ] Voicevoxエンジン
- [ ] パイプラインオーケストレーター
- [ ] 設定管理システム
- [ ] プログレスレポート
- [ ] バッチ処理最適化

### 5.3 テストタスク
- [ ] 単体テストの作成（目標カバレッジ80%）
- [ ] 統合テストシナリオ
- [ ] ベンチマークテスト
- [ ] 実ファイルでの動作確認

## 6. コマンドライン仕様

```bash
# 基本的な使用方法
epub2audio convert input.epub -o output.mp3

# 複数ファイルの変換
epub2audio convert *.md --engine voicevox --output-dir ./audio/

# 設定ファイル指定
epub2audio convert input.html --config custom.yml

# ドライラン（テキスト抽出のみ）
epub2audio convert input.epub --dry-run --extract-text output.txt

# 話者リスト表示
epub2audio list speakers --engine voicevox

# バッチ処理
epub2audio batch batch-list.txt --parallel 4

# ヘルプ
epub2audio --help
epub2audio convert --help
```

## 7. 設定ファイル仕様

```yaml
# default.yml
version: "1.0"

input:
  parser: auto  # auto, markdown, epub, html
  encoding: utf-8
  
processing:
  max_chunk_size: 140
  preserve_structure: true
  normalize_text: true
  
pause:
  strategy: standard  # standard, custom
  standard:
    heading_1: 2000
    heading_2: 1500
    heading_3: 1000
    paragraph: 800
    list_item: 400
    
tts:
  engine: voicepeak  # voicepeak, voicevox
  parallel_workers: 2
  retry_count: 3
  retry_delay: 1000
  
output:
  format: mp3  # mp3, wav, m4a
  bitrate: 128k
  merge: true
  keep_temp_files: false
  
logging:
  level: info  # debug, info, warn, error
  file: ./logs/epub2audio.log
```

## 8. 成功基準

### 8.1 機能要件
- [x] .md, .epub, .html ファイルの読み込み
- [x] 読み上げに適した部分の自動抽出
- [x] 文章構造を保持したポーズの挿入
- [x] Voicevox/Voicepeakの切り替え
- [x] 高品質な音声ファイルの生成

### 8.2 非機能要件
- [ ] 1GBのEPUBファイルを10分以内に処理
- [ ] メモリ使用量500MB以下
- [ ] 並列処理による高速化
- [ ] エラー時の適切なリカバリー
- [ ] 詳細なログ出力

### 8.3 品質基準
- [ ] テストカバレッジ80%以上
- [ ] golintエラー0
- [ ] ドキュメント完備
- [ ] CI/CDパイプライン構築

## 9. リスクと対策

### 9.1 技術的リスク
| リスク | 影響度 | 対策 |
|--------|--------|------|
| TTSエンジンのAPI変更 | 高 | 抽象化レイヤーで吸収 |
| 大容量ファイルでのメモリ不足 | 中 | ストリーミング処理 |
| 文字エンコーディング問題 | 中 | 自動検出と変換 |

### 9.2 スケジュールリスク
- Voicevox APIの学習曲線
- ffmpeg統合の複雑さ
- テストデータの準備

## 10. 今後の拡張予定

### 10.1 短期（3ヶ月）
- [ ] PDF入力サポート
- [ ] Google Cloud TTS対応
- [ ] Web UI

### 10.2 中期（6ヶ月）
- [ ] リアルタイム処理
- [ ] 音声編集機能
- [ ] クラウド版

### 10.3 長期（1年）
- [ ] 多言語対応
- [ ] AI音声認識との統合
- [ ] 音声ブック管理システム