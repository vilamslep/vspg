package config

const (
	DAILY = iota
	WEEKLY
	MONTHLY
)

func getKindPrewiew(kind int) string {
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
	Daily  ScheduleItem `yaml:"daily"`
	Weekly ScheduleItem `yaml:"weekly"`
	Month  ScheduleItem `yaml:"monthly"`
}

type ScheduleItem struct {
	Dbs       []string `yaml:"dbs"`
	KeepCount int `yaml:"keep_count"`
	Repeat    []int `yaml:"repeat"`
	Kind      int
}

//TODO
func NewScheduleItem(conf map[string]string, kind int) *ScheduleItem {
	si := ScheduleItem{}
	return &si
}

func (si ScheduleItem) GetKindPrewiew() string {
	return getKindPrewiew(si.Kind)
}

func (si ScheduleItem) NeedToRun() bool {
	return true
}

func (si ScheduleItem) checkDailySchedules() bool {
	return true
}

func (si ScheduleItem) checkWeeklySchedules() bool {
	return true
}

func (si ScheduleItem) checkMonthlySchedules() bool {
	return true
}
