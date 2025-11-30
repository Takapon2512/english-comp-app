package utils

import (
	"encoding/json"
	"fmt"
	"strings"
)

// extractFirstJSONObject 文字列から最初のJSONオブジェクトを抽出する
func ExtractFirstJSONObject(input string) (string, error) {
	input = strings.TrimSpace(input)

	// JSONオブジェクトの開始位置を探す
	start := strings.Index(input, "{")
	if start == -1 {
		return "", fmt.Errorf("JSONオブジェクトが見つかりません")
	}

	// 対応する閉じ括弧を探す
	braceCount := 0
	end := -1

	for i := start; i < len(input); i++ {
		switch input[i] {
		case '{':
			braceCount++
		case '}':
			braceCount--
			if braceCount == 0 {
				end = i + 1
				break
			}
		}
	}

	if end == -1 {
		return "", fmt.Errorf("JSONオブジェクトの終了が見つかりません")
	}

	jsonStr := input[start:end]

	// JSONの妥当性を確認
	var temp interface{}
	if err := json.Unmarshal([]byte(jsonStr), &temp); err != nil {
		return "", fmt.Errorf("無効なJSON形式です: %w", err)
	}

	return jsonStr, nil
}

// ParseJSONStringArray JSON文字列を文字列配列にパースする
func ParseJSONStringArray(jsonStr string) ([]string, error) {
	if jsonStr == "" {
		return []string{}, nil
	}

	var result []string
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil, fmt.Errorf("JSON文字列の配列への変換に失敗しました: %w", err)
	}

	return result, nil
}

// StringArrayToJSON 文字列配列をJSON文字列に変換する
func StringArrayToJSON(arr []string) (string, error) {
	if len(arr) == 0 {
		return "[]", nil
	}

	jsonBytes, err := json.Marshal(arr)
	if err != nil {
		return "", fmt.Errorf("文字列配列のJSON変換に失敗しました: %w", err)
	}

	return string(jsonBytes), nil
}
