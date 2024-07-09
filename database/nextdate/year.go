package nextdate

// CalculateNextDateAfterYear возвращает следующую дату, исходя из правила repeat "Y"
func CalculateNextDateAfterYear() string {
	nextDateDT := startDate.AddDate(1, 0, 0)
	for nextDateDT.Before(now) || nextDateDT.Equal(now) {
		nextDateDT = nextDateDT.AddDate(1, 0, 0)
	}
	nextDate = nextDateDT.Format(dateFormat)
	return nextDate
}
