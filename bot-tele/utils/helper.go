package utils

import (
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"bot-tele/model"
)

var validCommands = map[string]bool{
	"/start":     true,
	"/close":     true,
	"/bantuan":   true,
	"🆘Bantuan":   true,
	"/harian":    true,
	"📆Harian":    true,
	"/bulanan":   true,
	"📅Bulanan":   true,
	"/hapus":     true,
	"🔥Hapus":     true,
	"/daftar":    true,
	"📝Daftar":    true,
	"/dashboard": true,
	"📊Dashboard": true,
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

func HasTransactionAmount(input string) bool {
	input = strings.ToLower(input)
	input = strings.TrimSpace(input)

	// Regex untuk mencari angka dengan satuan lokal (rb, ribu, jt, juta, k, dll)
	amountRegex := regexp.MustCompile(`(?i)\b\d{1,3}(\.\d{3})*(rb|k|ribu|jt|juta)?\b|\b\d+\b`)

	return amountRegex.MatchString(input)
}

func ParseTransactionType(t string) (model.TransactionType, error) {
	switch t {
	case string(model.INCOME), string(model.EXPENSE):
		return model.TransactionType(t), nil
	default:
		return "", errors.New("invalid transaction type")
	}
}

func FormatDailyReport(transactions []model.Transaction) string {
	var (
		report       strings.Builder
		totalOut     float64
		totalIn      float64
	)

	report.WriteString("📊 *Laporan Hari Ini*\n")

	for _, tx := range transactions {
		var transactionType string
		if tx.TransactionType == "EXPENSE" {
			transactionType = "🔴"
		} else if tx.TransactionType == "INCOME" {
			transactionType = "🟢"
		}

		formatAmount := FormatRupiah(tx.Amount)

		line := fmt.Sprintf("#%d %s %s %s\n", tx.ID, transactionType, formatAmount, tx.OriginalText)
		report.WriteString(line)

		if tx.TransactionType == "EXPENSE" {
			totalOut += tx.Amount
		} else if tx.TransactionType == "INCOME" {
			totalIn += tx.Amount
		}
	}

	report.WriteString("\n")
	report.WriteString(fmt.Sprintf("🟢 Total Pemasukan: %s\n", FormatRupiah(totalIn)))
	report.WriteString(fmt.Sprintf("🔴 Total Pengeluaran: %s\n", FormatRupiah(totalOut)))

	return report.String()
}

func FormatMonthlyReport(transactions []model.Transaction) string {
	var (
		report       strings.Builder
		totalOut     float64
		totalIn      float64
	)

	report.WriteString("📆 *Laporan Bulan Ini*\n")

	for _, tx := range transactions {
		var transactionType string
		if tx.TransactionType == "EXPENSE" {
			transactionType = "🔴"
		} else if tx.TransactionType == "INCOME" {
			transactionType = "🟢"
		}

		formatAmount := FormatRupiah(tx.Amount)
		formatDate := FormatDate(tx.TransactionDate)

		line := fmt.Sprintf("#%d %s %s  %s %s\n", tx.ID, transactionType, formatDate, formatAmount, tx.OriginalText)
		report.WriteString(line)

		if tx.TransactionType == "EXPENSE" {
			totalOut += tx.Amount
		} else if tx.TransactionType == "INCOME" {
			totalIn += tx.Amount
		}
	}

	report.WriteString("\n")
	report.WriteString(fmt.Sprintf("🟢 Total Pemasukan: %s\n", FormatRupiah(totalIn)))
	report.WriteString(fmt.Sprintf("🔴 Total Pengeluaran: %s\n", FormatRupiah(totalOut)))

	return report.String()
}

func IsCommand(input string) bool {
	input = strings.TrimSpace(input)
	if input == "" {
		return false
	}

	// Ambil kata pertama saja
	parts := strings.Fields(input)
	firstWord := parts[0]

	_, isValid := validCommands[firstWord]
	return isValid
}

// Format date ke bahasa Indonesia
func FormatDateIndonesian(t time.Time) string {
	monthNames := map[time.Month]string{
		time.January:   "Januari",
		time.February:  "Februari",
		time.March:     "Maret",
		time.April:     "April",
		time.May:       "Mei",
		time.June:      "Juni",
		time.July:      "Juli",
		time.August:    "Agustus",
		time.September: "September",
		time.October:   "Oktober",
		time.November:  "November",
		time.December:  "Desember",
	}
	
	return fmt.Sprintf("%d %s %d", t.Day(), monthNames[t.Month()], t.Year())
}

// Format date ke bahasa Indonesia pendek
func FormatDateIndonesianShort(t time.Time) string {
	monthNames := map[time.Month]string{
		time.January:   "Jan",
		time.February:  "Feb",
		time.March:     "Mar",
		time.April:     "Apr",
		time.May:       "Mei",
		time.June:      "Jun",
		time.July:      "Jul",
		time.August:    "Ags",
		time.September: "Sep",
		time.October:   "Okt",
		time.November:  "Nov",
		time.December:  "Des",
	}
	
	return fmt.Sprintf("%d %s %d", t.Day(), monthNames[t.Month()], t.Year())
}

// Format dengan hari dalam bahasa Indonesia
func FormatDateTimeIndonesian(t time.Time) string {
	dayNames := map[time.Weekday]string{
		time.Sunday:    "Minggu",
		time.Monday:    "Senin",
		time.Tuesday:   "Selasa",
		time.Wednesday: "Rabu",
		time.Thursday:  "Kamis",
		time.Friday:    "Jumat",
		time.Saturday:  "Sabtu",
	}
	
	monthNames := map[time.Month]string{
		time.January:   "Januari",
		time.February:  "Februari",
		time.March:     "Maret",
		time.April:     "April",
		time.May:       "Mei",
		time.June:      "Juni",
		time.July:      "Juli",
		time.August:    "Agustus",
		time.September: "September",
		time.October:   "Oktober",
		time.November:  "November",
		time.December:  "Desember",
	}
	
	return fmt.Sprintf("%s, %d %s %d", 
		dayNames[t.Weekday()], 
		t.Day(), 
		monthNames[t.Month()], 
		t.Year())
}

// Format berbagai pilihan
func FormatDateCustom(t time.Time, format string) string {
	switch format {
	case "iso":
		return t.Format("2006-01-02")
	case "iso-time":
		return t.Format("2006-01-02 15:04:05")
	case "iso-full":
		return t.Format(time.RFC3339)
	case "id":
		return FormatDateIndonesian(t)
	case "id-short":
		return FormatDateIndonesianShort(t)
	case "id-full":
		return FormatDateTimeIndonesian(t)
	case "us":
		return t.Format("January 02, 2006")
	case "us-short":
		return t.Format("Jan 02, 2006")
	default:
		return t.Format("2006-01-02")
	}
}
