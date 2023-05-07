package main

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

// XHTMLに読み上げ用のポーズ文字列を追加する
//
// - `<span>内容</span>`タグを`<span>内容　</span>`のように</span>タグ直前に全角スペースを2個追加すること
// - `<h1>内容</h1>`タグを `<h1>内容　　　　　　　　　　　　　　　　　　</h1>`のように</h1>タグの直前に全角スペース18個追加すること
// - <h2><h3>とレベルを低くするごとに全角スペースを2つ減らすこと
// - <p><h1〜6>タグの終わりで改行をすること
// - <br>タグで改行をすること
func addPause(input string) string {
	for i := 1; i <= 6; i++ {
		re := regexp.MustCompile(fmt.Sprintf(`<h%d>(.*?)<\/h%d>`, i, i))
		input = re.ReplaceAllStringFunc(input, func(s string) string {
			tag := s[0:3]
			content := s[3 : len(s)-5]
			spaces := strings.Repeat("　", 18-2*(i-1))
			return fmt.Sprintf("%s%s%s", tag, content+spaces, s[len(s)-5:])
		})
	}

	input = regexp.MustCompile(`<span>(.*?)<\/span>`).ReplaceAllStringFunc(input, func(s string) string {
		content := s[6 : len(s)-7]
		return fmt.Sprintf("<span>%s　　</span>", content)
	})

	input = strings.ReplaceAll(input, "</p>", "</p>\n")
	input = regexp.MustCompile(`<\/h[1-8]>`).ReplaceAllStringFunc(input, func(s string) string {
		return s + "\n"
	})
	input = strings.ReplaceAll(input, "<br>", "<br>\n")

	return input
}

// XHTMLに読み上げ用のポーズ文字列を追加しxhtmlタグを削除する
func xhtml2txt(input string) string {
	input = addPause(input)
	input = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(input, "")
	input = regexp.MustCompile(`^\s*$\n`).ReplaceAllString(input, "")
	input = strings.TrimRight(input, "\n")
	return input
}

// 文字列を指定された最大文字数で分割する
//
// - 指定された最大文字数を必ず守ること
// - , . 、 。 または改行の文章の区切りでのみ分割し、それ以外では分割しないこと
func splitString(input string, maxChar int) []string {
	inputLen := utf8.RuneCountInString(input)
	if inputLen <= maxChar {
		return []string{input}
	}

	// 処理
	var result []string
	start := 0
	lineRunes := []rune(input)
	// 行の文字数が残っている限り処理を繰り返す
	for start < len(lineRunes) {
		// 最大文字数以上にならないように終了位置を調整
		end := start + maxChar
		if end > inputLen {
			end = inputLen
		}

		if end >= len(lineRunes) {
			result = append(result, string(lineRunes[start:]))
			break
		}

		splitIndex := getIndexByDelimiter(string(lineRunes[start:end]))

		// 区切り文字が見つかった場合
		if splitIndex != -1 {
			end = start + splitIndex
		} else {
			// 区切り文字が見つからない場合、最大文字数以内で最も長い区切り文字を探す
			end = start + maxChar
		}

		result = append(result, string(lineRunes[start:end]))
		start = end
	}

	return result
}

func getIndexByDelimiter(src string) int {
	// 区切り文字のパターンを定義
	splitPattern := regexp.MustCompile("[,.、。，．\n]")

	substr := string(src)
	splitIndex := -1
	// 区切り文字の位置を検索
	for _, loc := range splitPattern.FindAllStringIndex(substr, -1) {
		splitIndexInRunes := utf8.RuneCountInString(substr[:loc[1]])
		if splitIndexInRunes > splitIndex {
			splitIndex = splitIndexInRunes
		}
	}

	return splitIndex
}
