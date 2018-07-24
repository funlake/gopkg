package timer

import (
	"testing"
)
func BenchmarkCron_Exec(b *testing.B) {
	cron := NewTimer()
	cron.Ready()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next(){
			//timer.secondWheel.GetSlot(3)
			cron.SetInterval(3, func() {

			})
			cron.SetInterval(63, func() {

			})
			cron.SetInterval(3663, func() {

			})
		}
	})
}

func TestCron_Exec(t *testing.T) {
	cron := NewTimer()
	cron.Ready()
	r := cron.SetInterval(61, func() {

	})
	if r.pos != 1{
		t.Error("Wrong slot position",r.pos)
	}
	if cron.minuteWheel.slot[r.pos] == nil{
		t.Error("Not put into minute slot correctly")
	}
	//if timer.secondWheel.slot[1] == nil{
	//	t.Error("Not put into second slot")
	//}
	r2 := cron.SetInterval(33, func() {

	})

	if r2.pos != 33 {
		t.Error("Wrong slot position",r2.pos)
	}

	if cron.secondWheel.slot[r2.pos] == nil{
		t.Error("Not put into second slot correctly")
	}
}

//func TestCron_Ready(t *testing.T) {
//	m := `{"code":-2,"data":[],"msg":"参数缺失"}`
//	type r struct{
//		Code int `json:"code"`
//		Data interface{} `json:"data"`
//		Msg string `json:"msg"`
//	}
//	var R r
//	json.Unmarshal([]byte(m),&R)
//	t.Log(R)
//}
