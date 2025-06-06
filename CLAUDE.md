# CLAUDE.md

このファイルは、Claude Code (claude.ai/code) がこのリポジトリで作業する際のガイダンスを提供します。

## プロジェクト概要

epub2audioは、複数のTTSエンジン（Voicevox、Voicepeak）を使用して、様々なテキスト形式（EPUB、Markdown、HTML）を音声に変換するGoベースのアプリケーションです。拡張可能なアーキテクチャを採用しています：

```
入力ファイル → パーサー → テキスト処理 → TTSエンジン → ffmpeg → 音声ファイル
(.epub/.md/.html)         (構造保持・ポーズ挿入)  (Voicevox/Voicepeak)
```

## アーキテクチャ

### 現在の実装状況
2つの主要モジュールが存在していますが、統合が必要です：

1. **epub2xhtml**: EPUBファイルからXHTMLコンテンツを抽出
2. **xhtml2audio**: XHTMLをテキストに変換しTTS用に準備

### 計画中の新アーキテクチャ
```
├── internal/
│   ├── parser/         # 入力パーサー（EPUB、Markdown、HTML）
│   ├── tts/           # TTSエンジン（Voicevox、Voicepeak）
│   ├── processor/     # テキスト処理（ポーズ挿入、分割）
│   └── orchestrator/  # パイプライン制御
```

主要インターフェース：
- **Parser**: 各入力形式を統一されたDocument構造に変換
- **TTSEngine**: 各TTSエンジンを統一APIで制御
- **PauseStrategy**: エンジンごとに最適なポーズ表現を生成

## 開発コマンド

```bash
# ビルド（appディレクトリから）
cd app
go build -o epub2audio

# テスト実行
cd app/xhtml2audio
go test -v

# 特定のテストを実行
go test -v -run TestAddPause
```

## 設定

プロジェクトは `config/config.yml` を使用：
- `epub2audio.output`: 音声ファイルの出力ディレクトリ
- `epub2audio.tmp`: 処理用の一時ディレクトリ
- `ffmpeg.path`: ffmpeg実行ファイルのパス

## テストアプローチ

テストはGoの標準テストパッケージを使用。テストファイルは `*_test.go` パターンに従います。現在のテストはxhtml2audioモジュールのテキスト処理関数をカバーしています。

## 現在の状態と今後の実装

### 実装済み
- EPUBからXHTML抽出の基本ロジック
- テキスト処理（ポーズ挿入、分割）の基本ロジック
- Voicepeak用の全角スペースによるポーズ制御

### 実装予定
- 入力パーサーの抽象化（Markdown、HTML対応）
- TTSエンジンの抽象化（Voicevox対応）
- CLIインターフェース
- 設定ファイルシステム
- バッチ処理と並列化

## 設計ドキュメント

詳細な設計は以下のドキュメントを参照：
- `doc/architecture_proposal.md` - 全体アーキテクチャ
- `doc/input_parser_design.md` - 入力パーサー設計
- `doc/tts_engine_design.md` - TTSエンジン設計
- `doc/implementation_plan.md` - 実装計画