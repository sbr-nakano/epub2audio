# 入力パーサー抽象化設計

## 概要
様々なテキスト形式（Markdown、EPUB、HTML）から読み上げに適した部分を抽出し、統一されたデータ構造に変換する。

## 1. 共通インターフェース

```go
// internal/parser/interface.go
package parser

type Parser interface {
    // ファイルを解析してDocumentを返す
    Parse(filePath string) (*Document, error)
    
    // このパーサーがファイルを処理できるか判定
    CanParse(filePath string) bool
    
    // サポートする拡張子のリスト
    SupportedExtensions() []string
}

type Document struct {
    Title    string
    Author   string
    Language string
    Chapters []Chapter
    Metadata map[string]string
}

type Chapter struct {
    ID       string
    Title    string
    Level    int  // 1=章, 2=節, 3=項
    Order    int
    Content  []ContentBlock
}

type ContentBlock struct {
    Type     BlockType
    Level    int        // 見出しレベル（h1=1, h2=2, ...）
    Content  string     // テキストコンテンツ
    Attributes map[string]string
    Children []ContentBlock
}

type BlockType int

const (
    BlockTypeHeading BlockType = iota
    BlockTypeParagraph
    BlockTypeList
    BlockTypeListItem
    BlockTypeBlockquote
    BlockTypeCode      // 読み上げから除外
    BlockTypeImage     // 代替テキストを使用
    BlockTypeTable     // 特別な処理が必要
    BlockTypeFootnote  // オプションで読み上げ
)
```

## 2. 各パーサーの実装

### 2.1 Markdownパーサー

```go
// internal/parser/markdown.go
package parser

import (
    "github.com/yuin/goldmark"
    "github.com/yuin/goldmark/ast"
)

type MarkdownParser struct {
    // 読み上げから除外する要素
    excludeCodeBlocks bool
    excludeFootnotes  bool
    excludeTables     bool
}

func (p *MarkdownParser) Parse(filePath string) (*Document, error) {
    // 1. ファイルを読み込み
    // 2. Goldmarkでパース
    // 3. ASTを走査して読み上げ可能な要素を抽出
    // 4. Document構造に変換
}

// 読み上げに適した部分の判定
func (p *MarkdownParser) shouldInclude(node ast.Node) bool {
    switch node.Kind() {
    case ast.KindCodeBlock, ast.KindFencedCodeBlock:
        return !p.excludeCodeBlocks
    case ast.KindHTMLBlock:
        return false  // HTMLブロックは除外
    case ast.KindImage:
        return true   // alt textを抽出
    default:
        return true
    }
}
```

### 2.2 EPUBパーサー

```go
// internal/parser/epub.go
package parser

type EPUBParser struct {
    // 目次の解析方法
    tocStrategy TOCStrategy
    // 本文抽出のフィルター
    contentFilter ContentFilter
}

func (p *EPUBParser) Parse(filePath string) (*Document, error) {
    // 1. EPUBファイルを解凍
    // 2. content.opfから読み順を取得
    // 3. toc.ncx または nav.xhtmlから目次構造を取得
    // 4. 各XHTMLファイルから本文を抽出
    // 5. 読み上げ不要な要素を除外
}

// 読み上げから除外すべき要素
func (p *EPUBParser) filterContent(elem *html.Node) bool {
    // 除外: <aside>, <nav>, <script>, <style>
    // 除外: class="pagenum", class="footnote"
    // 含む: <p>, <h1>-<h6>, <blockquote>, <li>
}
```

### 2.3 HTMLパーサー

```go
// internal/parser/html.go
package parser

type HTMLParser struct {
    // 本文抽出の戦略
    extractionStrategy ExtractionStrategy
    // カスタムセレクター
    contentSelectors []string
}

// 読み上げに適した部分の自動検出
func (p *HTMLParser) detectMainContent(doc *html.Node) *html.Node {
    // 1. <main>要素を探す
    // 2. <article>要素を探す
    // 3. role="main"属性を持つ要素を探す
    // 4. id="content"やclass="content"を探す
    // 5. 最も多くのテキストを含む要素を選択
}
```

## 3. パーサーファクトリー

```go
// internal/parser/factory.go
package parser

type ParserFactory struct {
    parsers map[string]Parser
}

func NewParserFactory() *ParserFactory {
    return &ParserFactory{
        parsers: map[string]Parser{
            ".md":   NewMarkdownParser(),
            ".epub": NewEPUBParser(),
            ".html": NewHTMLParser(),
            ".htm":  NewHTMLParser(),
        },
    }
}

func (f *ParserFactory) GetParser(filePath string) (Parser, error) {
    ext := strings.ToLower(filepath.Ext(filePath))
    parser, ok := f.parsers[ext]
    if !ok {
        return nil, fmt.Errorf("unsupported file type: %s", ext)
    }
    return parser, nil
}

// 自動検出モード
func (f *ParserFactory) DetectAndParse(filePath string) (*Document, error) {
    parser, err := f.GetParser(filePath)
    if err != nil {
        return nil, err
    }
    return parser.Parse(filePath)
}
```

## 4. 読み上げ適性の判定基準

### 4.1 含めるべき要素
- 見出し（h1-h6）
- 段落（p）
- リスト項目（li）
- 引用（blockquote）
- 画像の代替テキスト（alt属性）
- 表のヘッダーとセル（簡略化して）

### 4.2 除外すべき要素
- スクリプト（script）
- スタイル（style）
- ナビゲーション（nav）
- 広告・バナー
- ページ番号
- 索引
- 参考文献（オプション）

### 4.3 特別な処理が必要な要素
- ルビ（ruby）→ 読み仮名として処理
- 略語（abbr）→ 展開形を読む
- 数式（math）→ 説明テキストに変換
- 表（table）→ 行ごとに読み上げ

## 5. 設定例

```yaml
parser:
  markdown:
    exclude_code_blocks: true
    exclude_footnotes: false
    heading_detection: "atx"  # or "setext"
    
  epub:
    toc_strategy: "ncx"  # or "nav"
    extract_metadata: true
    ignore_page_numbers: true
    
  html:
    content_selectors:
      - "main"
      - "article"
      - "#content"
      - ".post-content"
    exclude_selectors:
      - ".advertisement"
      - ".sidebar"
      - "nav"
      - "footer"
```

## 6. エラーハンドリング

```go
type ParserError struct {
    Type    ErrorType
    Message string
    File    string
    Line    int
}

type ErrorType int

const (
    ErrorTypeFileNotFound ErrorType = iota
    ErrorTypeInvalidFormat
    ErrorTypeEncodingError
    ErrorTypeNoContent
    ErrorTypePartialContent
)
```

## 7. 拡張ポイント

1. **新しいフォーマットの追加**
   - Parserインターフェースを実装
   - ParserFactoryに登録

2. **カスタム抽出ルール**
   - ContentFilterインターフェースで拡張
   - サイト固有のルールを追加可能

3. **前処理・後処理**
   - テキストの正規化
   - 文字コード変換
   - 記号の読み替え