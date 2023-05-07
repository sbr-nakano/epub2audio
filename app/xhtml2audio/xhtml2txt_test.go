package main

import (
	"testing"
)

func TestAddPause(t *testing.T) {
	input := "<h1>見出し</h1><p>本文と<br>改行</p><p>本文と<span>目立つ文字</span></p>"
	expected := "<h1>見出し　　　　　　　　　　　　　　　　　　</h1>\n<p>本文と<br>\n改行</p>\n<p>本文と<span>目立つ文字　　</span></p>\n"

	output := addPause(input)
	if output != expected {
		t.Errorf("期待される出力: %s\n実際の出力: %s\n", expected, output)
	}
}

func TestXhtml2txt(t *testing.T) {
	input := "<h1>見出し</h1><p>本文と<br>改行</p><p>本文と<span>目立つ文字</span></p>"
	expected := "見出し　　　　　　　　　　　　　　　　　　\n本文と\n改行\n本文と目立つ文字　　"

	output := xhtml2txt(input)
	if output != expected {
		t.Errorf("期待される出力: %s\n実際の出力: %s\n", expected, output)
	}
}

func TestSplitString_Long(t *testing.T) {
	input := "あいう、えおかき。くけ,こさし.すせそ\nたちつてとな"
	maxChar := 5
	expected := []string{
		"あいう、",
		"えおかき。",
		"くけ,",
		"こさし.",
		"すせそ\n",
		"たちつてと",
		"な",
	}

	output := splitString(input, maxChar)
	if len(output) != len(expected) {
		t.Errorf("期待される配列の長さ: %d\n実際の配列の長さ: %d\n", len(expected), len(output))
	}

	for i := range output {
		if output[i] != expected[i] {
			t.Errorf("期待される出力: %s\n実際の出力: %s\n", expected[i], output[i])
		}
	}
}

func TestSplitString_NoDelimita(t *testing.T) {
	input := "あいうえお"
	maxChar := 5
	expected := []string{
		"あいうえお",
	}

	output := splitString(input, maxChar)
	if len(output) != len(expected) {
		t.Errorf("期待される配列の長さ: %d\n実際の配列の長さ: %d\n", len(expected), len(output))
	}

	for i := range output {
		if output[i] != expected[i] {
			t.Errorf("期待される出力: %s\n実際の出力: %s\n", expected[i], output[i])
		}
	}

	maxChar = 4
	expected = []string{
		"あいうえ",
		"お",
	}

	output = splitString(input, maxChar)
	if len(output) != len(expected) {
		t.Errorf("期待される配列の長さ: %d\n実際の配列の長さ: %d\n", len(expected), len(output))
	}

	for i := range output {
		if output[i] != expected[i] {
			t.Errorf("期待される出力: %s\n実際の出力: %s\n", expected[i], output[i])
		}
	}
}

func TestGetIndexByDelimiter(t *testing.T) {
	src := "あ、いうえお。かき"
	expected := 7

	result := getIndexByDelimiter(src)
	if result != expected {
		t.Errorf("期待される出力: %d, 実際の出力: %d", expected, result)
	}
}
