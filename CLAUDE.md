# CLAUDE.md

このファイルは、Claude Code (claude.ai/code) がこのリポジトリで作業する際のガイダンスを提供します。

## プロジェクト概要

epub2audioは、Voicepeak（日本語音声合成）を使用してEPUBファイルを音声に変換するGoベースのアプリケーションです。パイプラインアーキテクチャを採用しています：

```
EPUB → XHTML抽出 → テキスト処理 → Voicepeak (TTS) → ffmpeg → MP3
```

## アーキテクチャ

コードベースは2つの主要モジュールで構成されています：

1. **epub2xhtml**: EPUBファイルからXHTMLコンテンツを抽出
   - EPUBアーカイブを解凍
   - .xhtml/.htmlファイルを抽出
   - 目次（_toc.xhtml）を解析

2. **xhtml2audio**: XHTMLをテキストに変換しTTS用に準備
   - 音声制御用のポーズ文字（全角スペース　）を追加
   - HTMLタグを除去
   - テキストを140文字のチャンクに分割

主な実装詳細：
- ポーズ制御は日本語全角スペース（　）を使用: h1=18個、h2=16個、レベルごとに2個ずつ減少
- テキストは句読点（。、・？！\r\n）を考慮して140文字境界で分割
- 現在、main.goファイルは空 - コアロジックは存在するが統合が必要

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

## 現在の状態

プロジェクト構造は完成していますが、実装は未完成です：
- コアのテキスト処理ロジックは実装済み
- メインエントリーポイントの実装が必要
- Voicepeakとffmpegの統合が保留中
- コマンドラインインターフェースが未実装