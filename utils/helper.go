package utils

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"

	"finance-bot/model"
)


func Slugify(s string) string {
    s = strings.ToLower(s)
    var b strings.Builder
    for _, r := range s {
        if unicode.IsLetter(r) || unicode.IsNumber(r) {
            b.WriteRune(r)
        } else if unicode.IsSpace(r) || r == '-' {
            b.WriteRune('-')
        }
        // abaikan karakter lain
    }
    return strings.Trim(b.String(), "-")
}

func FormatDate(date time.Time) string {
	day := date.Format("02")
	month := date.Format("01")
	year := date.Format("06")

	return fmt.Sprintf("%s/%s/%s", day, month, year)
}


func FormatRupiah(amount float64) string {
	number := fmt.Sprintf("%.0f", amount)

	var result []string
	for i, j := len(number), 0; i > 0; i -= 3 {
		start := i - 3
		if start < 0 {
			start = 0
		}
		result = append([]string{number[start:i]}, result...)
		j++
	}

	return "Rp" + strings.Join(result, ".")
}

func FormatTransactionType(tipe string) (string, error) {
	if tipe == string(model.EXPENSE) {
		return "Keluar", nil
	} else if tipe == string(model.INCOME) {
        return "Masuk", nil
    }

    return "", errors.New("invalid transaction type")
}

func ParseTransactionType(t string) (model.TransactionType, error) {
	switch t {
	case string(model.INCOME), string(model.EXPENSE):
		return model.TransactionType(t), nil
	default:
		return "", errors.New("invalid transaction type")
	}
}

func EscapeMarkdown(text string) string {
	replacer := strings.NewReplacer(
		"_", "\\_",
		"*", "\\*",
		"[", "\\[",
		"]", "\\]",
		"(", "\\(",
		")", "\\)",
		"~", "\\~",
		"`", "\\`",
		">", "\\>",
		"#", "\\#",
		"+", "\\+",
		"-", "\\-",
		"=", "\\=",
		"|", "\\|",
		"{", "\\{",
		"}", "\\}",
		".", "\\.",
		"!", "\\!",
	)
	return replacer.Replace(text)
}

func EncodeChatID(chatID int64) string {
	chatIDStr := strconv.FormatInt(chatID, 10)
	return base64.URLEncoding.EncodeToString([]byte(chatIDStr))
}
