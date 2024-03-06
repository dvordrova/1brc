package lib

import "fmt"

func FormattedNumber(num int64) string {
	if num >= 0 {
		return fmt.Sprintf("%d.%d", num/10, num%10)
	} else {
		return fmt.Sprintf("-%d.%d", -num/10, -num%10)
	}
}

func FormattedAvg(sum, count int64) string {
	avg := sum / count
	if avg >= 0 {
		avg = (sum + count - 1) / count
	}
	return FormattedNumber(avg)
}
