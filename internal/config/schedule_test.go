package config

import (
	"testing"
	"time"
)

func TestCheckDailySchedules(t *testing.T) {
	t.Log("Testning daily schedules")
	{
		testID := 1
		t.Logf("Testing on monday with weekends. Id %d", testID)
		{
			s := getTestScheduleEverydayWithWeekendsMonday()
			if ok := s.Daily.checkDailySchedules(); !ok {
				t.Error("Schedules error. Should be true ")
			}
		}

		testID++
		t.Logf("Testing on Friday with weekends. Id %d", testID)
		{
			s := getTestScheduleEverydayWithWeekendsFriday()
			if ok := s.Daily.checkDailySchedules(); !ok {
				t.Error("Schedules error. Should be true ")
			}
		}

		testID++
		t.Logf("Testing on Sunday with weekends. Id %d", testID)
		{
			s := getTestScheduleEverydayWithWeekendsSunday()
			if ok := s.Daily.checkDailySchedules(); ok {
				t.Error("Schedules error. Should be false ")
			}
		}

		t.Logf("Testing on monday with weekends. Id %d", testID)
		{
			s := getTestScheduleEverydayWithoutWeekendsMonday()
			if ok := s.Daily.checkDailySchedules(); !ok {
				t.Error("Schedules error. Should be true ")
			}
		}

		testID++
		t.Logf("Testing on Friday without weekends. Id %d", testID)
		{
			s := getTestScheduleEverydayWithoutWeekendsFriday()
			if ok := s.Daily.checkDailySchedules(); !ok {
				t.Error("Schedules error. Should be true ")
			}
		}

		testID++
		t.Logf("Testing on Sunday without weekends. Id %d", testID)
		{
			s := getTestScheduleEverydayWithoutWeekendsSunday()
			if ok := s.Daily.checkDailySchedules(); !ok {
				t.Error("Schedules error. Should be true ")
			}
		}
	}
}

func TestCheckWeeklySchedules(t *testing.T) {
	t.Log("Testning weekly schedules")
	{
		testID := 1
		t.Logf("Testing on monday with weekends. Id %d", testID)
		{
			s := getTestScheduleEverydayWithWeekendsMonday()
			if ok := s.Weekly.checkWeeklySchedules(); ok {
				t.Error("Schedules error. Should be false ")
			}
		}

		testID++
		t.Logf("Testing on Friday with weekends. Id %d", testID)
		{
			s := getTestScheduleEverydayWithWeekendsFriday()
			if ok := s.Weekly.checkWeeklySchedules(); !ok {
				t.Error("Schedules error. Should be true ")
			}
		}

		testID++
		t.Logf("Testing on Sunday with weekends. Id %d", testID)
		{
			s := getTestScheduleEverydayWithWeekendsSunday()
			if ok := s.Weekly.checkWeeklySchedules(); ok {
				t.Error("Schedules error. Should be false ")
			}
		}

		t.Logf("Testing on monday with weekends. Id %d", testID)
		{
			s := getTestScheduleEverydayWithoutWeekendsMonday()
			if ok := s.Weekly.checkWeeklySchedules(); ok {
				t.Error("Schedules error. Should be false ")
			}
		}

		testID++
		t.Logf("Testing on Friday without weekends. Id %d", testID)
		{
			s := getTestScheduleEverydayWithoutWeekendsFriday()
			if ok := s.Weekly.checkWeeklySchedules(); ok {
				t.Error("Schedules error. Should be false ")
			}
		}

		testID++
		t.Logf("Testing on Sunday without weekends. Id %d", testID)
		{
			s := getTestScheduleEverydayWithoutWeekendsSunday()
			if ok := s.Weekly.checkWeeklySchedules(); !ok {
				t.Error("Schedules error. Should be true ")
			}
		}
	}
}

func TestCheckMonthlySchedules(t *testing.T) {
	t.Log("Testning monthly schedules")
	{
		testID := 1
		t.Logf("Testing on Sunday with weekends. Id %d", testID)
		{
			s := getTestScheduleMonthlyWithWeekendsSunday()
			if ok := s.Monthly.checkMonthlySchedules(); ok {
				t.Error("Schedules error. Should be false ")
			}
		}

		testID++
		t.Logf("Testing on Friday with weekends. Id %d", testID)
		{
			s := getTestScheduleMonthlyWithWeekendsFriday()
			if ok := s.Monthly.checkMonthlySchedules(); !ok {
				t.Error("Schedules error. Should be true ")
			}
		}

		testID++
		t.Logf("Testing on Sunday without weekends. Id %d", testID)
		{
			s := getTestScheduleMonthlyWithoutWeekendsSunday()
			if ok := s.Monthly.checkMonthlySchedules(); !ok {
				t.Error("Schedules error. Should be true ")
			}
		}

		testID++
		t.Logf("Testing on Friday without weekends. Id %d", testID)
		{
			s := getTestScheduleMonthlyWithoutWeekendsFriday()
			if ok := s.Monthly.checkMonthlySchedules(); ok {
				t.Error("Schedules error. Should be false ")
			}
		}
	}
}

//Test sets
func getTestScheduleEverydayWithWeekendsMonday() Schedule {
	today := time.Date(2022, 6, 27, 0, 0, 0, 0, time.UTC)
	return Schedule{
		Daily:   ScheduleItem{Repeat: []int{0}, Today: today, Weekend: []int{6, 7}},
		Weekly:  ScheduleItem{Repeat: []int{0}, Today: today, Weekend: []int{6, 7}},
		Monthly: ScheduleItem{Repeat: []int{0}, Today: today, Weekend: []int{6, 7}},
		Weekend: []int{6, 7},
	}
}

func getTestScheduleEverydayWithWeekendsFriday() Schedule {
	today := time.Date(2022, 7, 1, 0, 0, 0, 0, time.UTC)
	return Schedule{
		Daily:   ScheduleItem{Repeat: []int{0}, Today: today, Weekend: []int{6, 7}},
		Weekly:  ScheduleItem{Repeat: []int{0}, Today: today, Weekend: []int{6, 7}},
		Monthly: ScheduleItem{Repeat: []int{0}, Today: today, Weekend: []int{6, 7}},
		Weekend: []int{6, 7},
	}
}

func getTestScheduleEverydayWithWeekendsSunday() Schedule {
	today := time.Date(2022, 07, 3, 0, 0, 0, 0, time.Local)
	return Schedule{
		Daily:   ScheduleItem{Repeat: []int{0}, Today: today, Weekend: []int{6, 7}},
		Weekly:  ScheduleItem{Repeat: []int{0}, Today: today, Weekend: []int{6, 7}},
		Monthly: ScheduleItem{Repeat: []int{0}, Today: today, Weekend: []int{6, 7}},
		Weekend: []int{6, 7},
	}
}

func getTestScheduleEverydayWithoutWeekendsMonday() Schedule {
	today := time.Date(2022, 6, 27, 0, 0, 0, 0, time.UTC)
	return Schedule{
		Daily:   ScheduleItem{Repeat: []int{0}, Today: today},
		Weekly:  ScheduleItem{Repeat: []int{0}, Today: today},
		Monthly: ScheduleItem{Repeat: []int{0}, Today: today},
	}
}

func getTestScheduleEverydayWithoutWeekendsFriday() Schedule {
	today := time.Date(2022, 7, 1, 0, 0, 0, 0, time.UTC)
	return Schedule{
		Daily:   ScheduleItem{Repeat: []int{0}, Today: today},
		Weekly:  ScheduleItem{Repeat: []int{0}, Today: today},
		Monthly: ScheduleItem{Repeat: []int{0}, Today: today},
	}
}

func getTestScheduleEverydayWithoutWeekendsSunday() Schedule {
	today := time.Date(2022, 7, 3, 0, 0, 0, 0, time.UTC)
	return Schedule{
		Daily:   ScheduleItem{Repeat: []int{0}, Today: today},
		Weekly:  ScheduleItem{Repeat: []int{0}, Today: today},
		Monthly: ScheduleItem{Repeat: []int{0}, Today: today},
	}
}

func getTestScheduleMonthlyWithWeekendsFriday() Schedule {
	today := time.Date(2022, 7, 29, 0, 0, 0, 0, time.UTC)
	return Schedule{
		Daily:   ScheduleItem{Repeat: []int{0}, Today: today, Weekend: []int{6,7}},
		Weekly:  ScheduleItem{Repeat: []int{0}, Today: today, Weekend: []int{6,7}},
		Monthly: ScheduleItem{Repeat: []int{0}, Today: today, Weekend: []int{6,7}},
	}
}

func getTestScheduleMonthlyWithWeekendsSunday() Schedule {
	today := time.Date(2022, 7, 31, 0, 0, 0, 0, time.UTC)
	return Schedule{
		Daily:   ScheduleItem{Repeat: []int{0}, Today: today, Weekend: []int{6,7}},
		Weekly:  ScheduleItem{Repeat: []int{0}, Today: today, Weekend: []int{6,7}},
		Monthly: ScheduleItem{Repeat: []int{0}, Today: today, Weekend: []int{6,7}},
	}
}

func getTestScheduleMonthlyWithoutWeekendsFriday() Schedule {
	today := time.Date(2022, 7, 29, 0, 0, 0, 0, time.UTC)
	return Schedule{
		Daily:   ScheduleItem{Repeat: []int{0}, Today: today},
		Weekly:  ScheduleItem{Repeat: []int{0}, Today: today},
		Monthly: ScheduleItem{Repeat: []int{0}, Today: today},
	}
}

func getTestScheduleMonthlyWithoutWeekendsSunday() Schedule {
	today := time.Date(2022, 7, 31, 0, 0, 0, 0, time.UTC)
	return Schedule{
		Daily:   ScheduleItem{Repeat: []int{0}, Today: today},
		Weekly:  ScheduleItem{Repeat: []int{0}, Today: today},
		Monthly: ScheduleItem{Repeat: []int{0}, Today: today},
	}
}


