package config

import "time"

const (
	DAILY = iota
	WEEKLY
	MONTHLY
)

func GetKindPrewiew(kind int) string {
	switch kind {
	case 0:
		return "Daily"
	case 1:
		return "Weekly"
	case 2:
		return "Monthly"
	default:
		return ""
	}
}

type Schedule struct {
	Daily   ScheduleItem `yaml:"daily"`
	Weekly  ScheduleItem `yaml:"weekly"`
	Monthly ScheduleItem `yaml:"monthly"`
	Weekend []int        `yaml:"weekend"`
}

type ScheduleItem struct {
	Dbs       []string `yaml:"dbs"`
	KeepCount int      `yaml:"keep_count"`
	Repeat    []int    `yaml:"repeat"`
	Kind      int
	Today     time.Time
	Weekend   []int
}

func (si ScheduleItem) Empty() bool {
	return false
}

//TODO
func NewScheduleItem(conf map[string]string, kind int, weekends []int) *ScheduleItem {
	si := ScheduleItem{Today: time.Now(), Weekend: weekends}
	return &si
}

func (si ScheduleItem) NeedToRun() bool {
	switch si.Kind {
	case DAILY:
		return si.checkDailySchedules()
	case WEEKLY:
		return si.checkWeeklySchedules()
	case MONTHLY:
		return si.checkMonthlySchedules()
	default:
		return false
	}
}

func (si ScheduleItem) checkDailySchedules() bool {

	if si.stopOnWeekend() {
		return false
	}

	nwk := weekday(si.Today)
	for i := range si.Repeat {
		if si.Repeat[i] == 0 {
			return true
		}
		if nwk == si.Repeat[i] {
			return true
		}
	}
	return false
}

func (si ScheduleItem) checkWeeklySchedules() bool {

	nwk := weekday(si.Today)
	if len(si.Weekend) == 0 && nwk != 7 {
		return false
	}

	r := false
	_, w := si.Today.ISOWeek()
	for i := range si.Repeat {
		if si.Repeat[i] == 0 {
			r = true
			break
		}
		if w == si.Repeat[i] {
			r = true
			break
		}
	}

	if len(si.Weekend) > 0 {
		alwd := make([]int, 0, 7-len(si.Weekend))
		for i := 1; i <= 7; i++ {
			f := false
			for j := range si.Weekend {
				if si.Weekend[j] == i {
					f = true
					break
				}
			}
			if !f {
				alwd = append(alwd, i)
			}
		}
		lastD := alwd[len(alwd)-1]
		r = lastD == nwk
	}

	return r
}

func (si ScheduleItem) checkMonthlySchedules() bool {

	startDay := time.Date(si.Today.Year(), si.Today.Month(), si.Today.Day(), 0, 0, 0, 0, si.Today.Location())
	endMonth := time.Date(si.Today.Year(), si.Today.Month()+1, 1, 0, 0, 0, 0, si.Today.Location()).Add(-24 * time.Hour)

	if len(si.Weekend) == 0 {
		return startDay == endMonth
	}

	emwd := weekday(endMonth)

	endOnWeekends := false
	for i := range si.Weekend {
		if si.Weekend[i] == emwd {
			endOnWeekends = true
			break
		}
	}

	if !endOnWeekends {
		return startDay == endMonth
	}

	alwd := make([]int, 0, 7-len(si.Weekend))
	for i := 1; i <= 7; i++ {
		f := false
		for j := range si.Weekend {
			if si.Weekend[j] == i {
				f = true
				break
			}
		}
		if !f {
			alwd = append(alwd, i)
		}
	}
	lastD := alwd[len(alwd)-1]
	diff := emwd - lastD
	newEndMonth := time.Date(endMonth.Year(), endMonth.Month(), endMonth.Day()-diff, 0, 0, 0, 0, endMonth.Location())
	return startDay == newEndMonth
}

func (si ScheduleItem) stopOnWeekend() bool {
	if len(si.Weekend) == 0 {
		return false
	}
	nwk := weekday(si.Today)
	for i := range si.Weekend {
		if si.Weekend[i] == nwk {
			return true
		}
	}

	return false
}

func weekday(date time.Time) int {
	wd := int(date.Weekday())
	if wd == 0 {
		wd = 7
	}
	return wd
}
