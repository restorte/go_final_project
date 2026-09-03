package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

func NextDate(now, date, repeat string) (string, error) {
	if repeat == "" {
		return "", nil
	}

	nowTime, err := time.Parse(dateFormat, now)
	if err != nil {
		return "", err
	}
	taskTime, err := time.Parse(dateFormat, date)
	if err != nil {
		return "", err
	}

	if strings.HasPrefix(repeat, "d ") {
		parts := strings.Fields(repeat)
		if len(parts) != 2 {
			return "", errors.New("invalid d format")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil || days <= 0 || days > 400 {
			return "", errors.New("invalid d interval")
		}
		next := taskTime
		for {
			next = next.AddDate(0, 0, days)
			if next.After(nowTime) {
				return next.Format(dateFormat), nil
			}
		}
	}

	if repeat == "y" {
		next := taskTime
		for {
			next = next.AddDate(1, 0, 0)
			if next.After(nowTime) {
				return next.Format(dateFormat), nil
			}
		}
	}

	if strings.HasPrefix(repeat, "w ") {
		parts := strings.Fields(repeat)
		if len(parts) != 2 {
			return "", errors.New("invalid w format")
		}
		weekdaysStr := strings.Split(parts[1], ",")
		var weekdays []time.Weekday
		for _, s := range weekdaysStr {
			d, err := strconv.Atoi(s)
			if err != nil || d < 1 || d > 7 {
				return "", errors.New("invalid weekday")
			}
			var wd time.Weekday
			if d == 7 {
				wd = time.Sunday
			} else {
				wd = time.Weekday(d)
			}
			weekdays = append(weekdays, wd)
		}
		next := taskTime
		for {
			next = next.AddDate(0, 0, 1)
			if next.After(nowTime) {
				for _, wd := range weekdays {
					if next.Weekday() == wd {
						return next.Format(dateFormat), nil
					}
				}
			}
			if next.After(nowTime.AddDate(1, 0, 0)) {
				return "", errors.New("no suitable date found")
			}
		}
	}

	if strings.HasPrefix(repeat, "m ") {
		parts := strings.Fields(repeat)
		if len(parts) < 2 {
			return "", errors.New("invalid m format")
		}
		daysStr := strings.Split(parts[1], ",")
		days := []int{}
		for _, s := range daysStr {
			d, err := strconv.Atoi(s)
			if err != nil {
				return "", errors.New("invalid day")
			}
			if d >= -2 && d <= 31 && d != 0 {
				days = append(days, d)
			} else {
				return "", errors.New("invalid day value")
			}
		}
		if len(days) == 0 {
			return "", errors.New("no days specified")
		}
		months := []int{}
		if len(parts) > 2 {
			monthsStr := strings.Split(parts[2], ",")
			for _, s := range monthsStr {
				m, err := strconv.Atoi(s)
				if err != nil || m < 1 || m > 12 {
					return "", errors.New("invalid month")
				}
				months = append(months, m)
			}
		}
		next := taskTime
		for {
			next = next.AddDate(0, 0, 1)
			if next.After(nowTime) {
				if len(months) > 0 {
					found := false
					for _, m := range months {
						if int(next.Month()) == m {
							found = true
							break
						}
					}
					if !found {
						continue
					}
				}
				day := next.Day()
				lastDay := time.Date(next.Year(), next.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
				secondLast := lastDay - 1
				for _, d := range days {
					if d == -1 && day == lastDay {
						return next.Format(dateFormat), nil
					}
					if d == -2 && day == secondLast {
						return next.Format(dateFormat), nil
					}
					if d == day {
						return next.Format(dateFormat), nil
					}
				}
			}
			if next.After(nowTime.AddDate(2, 0, 0)) {
				return "", errors.New("no suitable date found")
			}
		}
	}

	return "", errors.New("unsupported repeat rule")
}

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	now := r.FormValue("now")
	if now == "" {
		now = time.Now().Format(dateFormat)
	}
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")
	if date == "" || repeat == "" {
		http.Error(w, "Missing date or repeat", http.StatusBadRequest)
		return
	}
	next, err := NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write([]byte(next))
}
