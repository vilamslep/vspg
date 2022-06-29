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
	Monthly  ScheduleItem `yaml:"monthly"`
}

type ScheduleItem struct {
	Dbs       []string `yaml:"dbs"`
	KeepCount int `yaml:"keep_count"`
	Repeat    []int `yaml:"repeat"`
	Kind      int
}

func (si ScheduleItem) Empty() bool {
	return false
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
	// if self.repeat == None:
    //         return False
    //     if self.kind == Periodicity.DAILY:
    //         return self.__check_daily_schedules()
    //     elif self.kind == Periodicity.WEEKLY:
    //         return self.__check_weekly_schedules()
    //     elif self.kind == Periodicity.MONTHLY:
    //         return self.__check_monthly_schedules()
	return true
}

func (si ScheduleItem) checkDailySchedules() bool {
	// today = datetime.now()
        
    //     weekday = today.weekday()+1
            
    //     for day in self.repeat:
    //         if (day == 0 or day == weekday):
    //             return True
	return true
}

func (si ScheduleItem) checkWeeklySchedules() bool {
	// today = datetime.now()
        
	// week_number = today.isocalendar().week

	// if (today.weekday()+1) != 7:
	// 	return False

	// for week in self.repeat:
	// 	if (week == 0 or week == week_number):
	// 		return True
	return true
}

func (si ScheduleItem) checkMonthlySchedules() bool {
	// today = datetime.now()
        
    //     if today.month == 12:
    //         year = today.year + 1
    //         next_month = 1
    //     else:
    //         year = today.year
    //         next_month = today.month + 1 

    //     finish_day = (datetime(year, next_month, 1) - timedelta(days=1)).day

    //     if today.day != finish_day:
    //         return False

    //     for month in self.repeat:
    //         if (month == 0 or month == today.month):
    //             return True
	return true
}
