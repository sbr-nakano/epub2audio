# TTSエンジン抽象化設計

## 概要
Voicevox、Voicepeakなど異なるTTSエンジンを統一されたインターフェースで扱い、読み上げ時の構造保持（ポーズ/遅延）を各エンジンに適した形式で実現する。

## 1. TTSエンジンの特性比較

| 特性 | Voicevox | Voicepeak |
|------|----------|-----------|
| API形式 | REST API | CLI |
| ポーズ表現 | 音声合成クエリ内で指定 | 全角スペース |
| 話者指定 | speaker_id | ナレーター名 |
| 出力形式 | WAV | WAV |
| バッチ処理 | 非対応（1リクエスト1音声） | 対応（複数行） |
| 料金 | 無料 | 有料ライセンス |

## 2. 統一インターフェース

```go
// internal/tts/interface.go
package tts

type Engine interface {
    // エンジンの初期化
    Initialize(config Config) error
    
    // エンジン名
    Name() string
    
    // 音声合成
    Synthesize(request SynthesisRequest) (*SynthesisResult, error)
    
    // バッチ処理のサポート
    SupportsBatch() bool
    
    // バッチ音声合成
    SynthesizeBatch(requests []SynthesisRequest) ([]*SynthesisResult, error)
    
    // サポートする機能の確認
    Capabilities() Capabilities
    
    // 話者リストの取得
    GetSpeakers() ([]Speaker, error)
    
    // ポーズ変換
    ConvertPause(pause Pause) string
    
    // 終了処理
    Close() error
}

type Config struct {
    // 共通設定
    OutputDir    string
    TempDir      string
    
    // エンジン固有の設定
    EngineConfig map[string]interface{}
}

type SynthesisRequest struct {
    Text         string
    OutputPath   string
    Speaker      Speaker
    AudioConfig  AudioConfig
    Pauses       []PausePosition  // テキスト内のポーズ位置
}

type SynthesisResult struct {
    FilePath     string
    Duration     float64  // 秒
    Success      bool
    Error        error
}

type AudioConfig struct {
    Speed        float64  // 話速（0.5-2.0）
    Pitch        float64  // 音高（-0.15-0.15）
    Intonation   float64  // 抑揚（0.0-2.0）
    Volume       float64  // 音量（0.0-1.0）
    PrePause     float64  // 前の無音（秒）
    PostPause    float64  // 後の無音（秒）
}

type Speaker struct {
    ID           string
    Name         string
    Language     string
    Gender       string
    Description  string
    Styles       []string  // 感情表現など
}

type Capabilities struct {
    MaxTextLength      int
    SupportsPause      bool
    SupportsSSML       bool
    SupportsPitch      bool
    SupportsSpeed      bool
    SupportsVolume     bool
    SupportsEmotions   bool
    SupportsBatch      bool
}

type PausePosition struct {
    Position int    // テキスト内の位置
    Pause    Pause
}

type Pause struct {
    Type     PauseType
    Duration int      // ミリ秒
}

type PauseType int

const (
    PauseTypeShort   PauseType = iota  // 短いポーズ（読点相当）
    PauseTypeMedium                     // 中程度のポーズ（句点相当）
    PauseTypeLong                       // 長いポーズ（段落終了）
    PauseTypeChapter                    // 章の区切り
)
```

## 3. Voicevox実装

```go
// internal/tts/voicevox.go
package tts

import (
    "github.com/go-resty/resty/v2"
)

type VoicevoxEngine struct {
    client      *resty.Client
    apiEndpoint string
    speakers    []Speaker
}

func NewVoicevoxEngine() *VoicevoxEngine {
    return &VoicevoxEngine{
        client:      resty.New(),
        apiEndpoint: "http://localhost:50021",
    }
}

func (v *VoicevoxEngine) Initialize(config Config) error {
    // APIエンドポイントの設定
    if endpoint, ok := config.EngineConfig["endpoint"].(string); ok {
        v.apiEndpoint = endpoint
    }
    
    // 話者リストの取得
    return v.loadSpeakers()
}

func (v *VoicevoxEngine) Synthesize(request SynthesisRequest) (*SynthesisResult, error) {
    // 1. テキストをポーズ付きに変換
    textWithPauses := v.insertPauses(request.Text, request.Pauses)
    
    // 2. 音声合成用クエリを生成
    query, err := v.createAudioQuery(textWithPauses, request.Speaker.ID)
    if err != nil {
        return nil, err
    }
    
    // 3. クエリを調整（速度、音高など）
    v.adjustQuery(query, request.AudioConfig)
    
    // 4. 音声合成
    audio, err := v.synthesis(query, request.Speaker.ID)
    if err != nil {
        return nil, err
    }
    
    // 5. ファイルに保存
    err = v.saveAudio(audio, request.OutputPath)
    
    return &SynthesisResult{
        FilePath: request.OutputPath,
        Success:  err == nil,
        Error:    err,
    }, nil
}

// Voicevox用のポーズ挿入
func (v *VoicevoxEngine) insertPauses(text string, pauses []PausePosition) string {
    // Voicevoxでは音声合成クエリ内でポーズを指定
    // ここでは特殊マーカーを挿入し、後でクエリ生成時に処理
    result := text
    for _, p := range pauses {
        marker := v.getPauseMarker(p.Pause)
        result = insertAt(result, p.Position, marker)
    }
    return result
}

func (v *VoicevoxEngine) getPauseMarker(pause Pause) string {
    // Voicevox用の内部マーカー
    switch pause.Type {
    case PauseTypeShort:
        return "{{PAUSE_S}}"
    case PauseTypeMedium:
        return "{{PAUSE_M}}"
    case PauseTypeLong:
        return "{{PAUSE_L}}"
    case PauseTypeChapter:
        return "{{PAUSE_C}}"
    default:
        return ""
    }
}

// 音声合成クエリの調整
func (v *VoicevoxEngine) adjustQuery(query *AudioQuery, config AudioConfig) {
    query.SpeedScale = config.Speed
    query.PitchScale = config.Pitch
    query.IntonationScale = config.Intonation
    query.VolumeScale = config.Volume
    query.PrePhonemeLength = config.PrePause
    query.PostPhonemeLength = config.PostPause
}
```

## 4. Voicepeak実装

```go
// internal/tts/voicepeak.go
package tts

import (
    "os/exec"
    "strings"
)

type VoicepeakEngine struct {
    exePath    string
    narrator   string
    workDir    string
}

func NewVoicepeakEngine() *VoicepeakEngine {
    return &VoicepeakEngine{
        exePath: "voicepeak.exe",  // or Mac/Linux path
    }
}

func (v *VoicepeakEngine) Initialize(config Config) error {
    // 実行ファイルパスの設定
    if path, ok := config.EngineConfig["exe_path"].(string); ok {
        v.exePath = path
    }
    
    // デフォルトナレーターの設定
    if narrator, ok := config.EngineConfig["narrator"].(string); ok {
        v.narrator = narrator
    }
    
    v.workDir = config.TempDir
    return v.checkExecutable()
}

func (v *VoicepeakEngine) Synthesize(request SynthesisRequest) (*SynthesisResult, error) {
    // 1. テキストにポーズを挿入
    textWithPauses := v.insertPauses(request.Text, request.Pauses)
    
    // 2. Voicepeak用のコマンドを構築
    args := v.buildCommand(textWithPauses, request)
    
    // 3. 実行
    cmd := exec.Command(v.exePath, args...)
    err := cmd.Run()
    
    return &SynthesisResult{
        FilePath: request.OutputPath,
        Success:  err == nil,
        Error:    err,
    }, nil
}

// Voicepeak用のポーズ挿入（全角スペース）
func (v *VoicepeakEngine) insertPauses(text string, pauses []PausePosition) string {
    result := text
    offset := 0
    
    for _, p := range pauses {
        pauseStr := v.getPauseString(p.Pause)
        position := p.Position + offset
        result = result[:position] + pauseStr + result[position:]
        offset += len(pauseStr)
    }
    
    return result
}

func (v *VoicepeakEngine) getPauseString(pause Pause) string {
    // 全角スペースの数でポーズの長さを制御
    spaceCount := 0
    switch pause.Type {
    case PauseTypeShort:
        spaceCount = 2
    case PauseTypeMedium:
        spaceCount = 4
    case PauseTypeLong:
        spaceCount = 8
    case PauseTypeChapter:
        spaceCount = 16
    }
    
    return strings.Repeat("　", spaceCount)
}

func (v *VoicepeakEngine) buildCommand(text string, request SynthesisRequest) []string {
    args := []string{
        "--narrator", request.Speaker.Name,
        "--output", request.OutputPath,
        "--speed", fmt.Sprintf("%.2f", request.AudioConfig.Speed),
        "--pitch", fmt.Sprintf("%.2f", request.AudioConfig.Pitch),
    }
    
    // テキストは標準入力またはファイル経由で渡す
    args = append(args, "--text", text)
    
    return args
}

// バッチ処理のサポート
func (v *VoicepeakEngine) SynthesizeBatch(requests []SynthesisRequest) ([]*SynthesisResult, error) {
    // Voicepeakはテキストファイルによるバッチ処理をサポート
    batchFile := filepath.Join(v.workDir, "batch.txt")
    
    // バッチファイルの作成
    var lines []string
    for _, req := range requests {
        textWithPauses := v.insertPauses(req.Text, req.Pauses)
        lines = append(lines, fmt.Sprintf("%s\t%s", req.OutputPath, textWithPauses))
    }
    
    err := os.WriteFile(batchFile, []byte(strings.Join(lines, "\n")), 0644)
    if err != nil {
        return nil, err
    }
    
    // バッチ実行
    cmd := exec.Command(v.exePath, "--batch", batchFile)
    err = cmd.Run()
    
    // 結果の作成
    results := make([]*SynthesisResult, len(requests))
    for i, req := range requests {
        results[i] = &SynthesisResult{
            FilePath: req.OutputPath,
            Success:  err == nil,
            Error:    err,
        }
    }
    
    return results, nil
}
```

## 5. エンジンファクトリー

```go
// internal/tts/factory.go
package tts

type EngineFactory struct {
    engines map[string]func() Engine
}

func NewEngineFactory() *EngineFactory {
    return &EngineFactory{
        engines: map[string]func() Engine{
            "voicevox":  func() Engine { return NewVoicevoxEngine() },
            "voicepeak": func() Engine { return NewVoicepeakEngine() },
        },
    }
}

func (f *EngineFactory) CreateEngine(name string, config Config) (Engine, error) {
    constructor, ok := f.engines[name]
    if !ok {
        return nil, fmt.Errorf("unknown TTS engine: %s", name)
    }
    
    engine := constructor()
    if err := engine.Initialize(config); err != nil {
        return nil, fmt.Errorf("failed to initialize %s: %w", name, err)
    }
    
    return engine, nil
}

// 利用可能なエンジンのリスト
func (f *EngineFactory) AvailableEngines() []string {
    engines := make([]string, 0, len(f.engines))
    for name := range f.engines {
        engines = append(engines, name)
    }
    return engines
}
```

## 6. ポーズ戦略

```go
// internal/tts/pause_strategy.go
package tts

type PauseStrategy interface {
    // ブロックタイプからポーズを決定
    GetPause(blockType parser.BlockType, level int) Pause
    
    // テキスト内の自然なポーズ位置を検出
    DetectNaturalPauses(text string) []PausePosition
}

type StandardPauseStrategy struct {
    engine Engine
}

func (s *StandardPauseStrategy) GetPause(blockType parser.BlockType, level int) Pause {
    switch blockType {
    case parser.BlockTypeHeading:
        if level == 1 {
            return Pause{Type: PauseTypeChapter, Duration: 2000}
        } else if level == 2 {
            return Pause{Type: PauseTypeLong, Duration: 1500}
        } else {
            return Pause{Type: PauseTypeMedium, Duration: 1000}
        }
    case parser.BlockTypeParagraph:
        return Pause{Type: PauseTypeMedium, Duration: 800}
    case parser.BlockTypeListItem:
        return Pause{Type: PauseTypeShort, Duration: 400}
    default:
        return Pause{Type: PauseTypeShort, Duration: 300}
    }
}

// 句読点による自然なポーズの検出
func (s *StandardPauseStrategy) DetectNaturalPauses(text string) []PausePosition {
    var pauses []PausePosition
    
    for i, r := range text {
        switch r {
        case '。', '.':
            pauses = append(pauses, PausePosition{
                Position: i + 1,
                Pause:    Pause{Type: PauseTypeMedium, Duration: 600},
            })
        case '、', ',':
            pauses = append(pauses, PausePosition{
                Position: i + 1,
                Pause:    Pause{Type: PauseTypeShort, Duration: 300},
            })
        case '？', '！':
            pauses = append(pauses, PausePosition{
                Position: i + 1,
                Pause:    Pause{Type: PauseTypeLong, Duration: 800},
            })
        }
    }
    
    return pauses
}
```

## 7. 設定例

```yaml
tts:
  engine: voicevox  # or voicepeak
  
  # 共通設定
  common:
    output_format: wav
    temp_dir: ./tmp/tts
    
  # Voicevox設定
  voicevox:
    endpoint: http://localhost:50021
    speaker_id: 1  # ずんだもん（ノーマル）
    speed: 1.0
    pitch: 0.0
    intonation: 1.0
    
  # Voicepeak設定
  voicepeak:
    exe_path: /usr/local/bin/voicepeak
    narrator: "Japanese Female"
    speed: 1.0
    pitch: 0.0
    
  # ポーズ設定
  pause:
    strategy: standard
    custom_rules:
      heading_1: 2000  # ms
      heading_2: 1500
      paragraph: 800
      list_item: 400
```

## 8. 拡張性

1. **新しいTTSエンジンの追加**
   - Engineインターフェースを実装
   - EngineFactoryに登録
   - Google Cloud TTS、Amazon Polly、Azure Speech等

2. **カスタムポーズ戦略**
   - PauseStrategyインターフェースを実装
   - 言語別、用途別の戦略を追加可能

3. **後処理機能**
   - 音量正規化
   - ノイズ除去
   - フォーマット変換（WAV→MP3）