package timewheel

import (
	"testing"
	"encoding/json"
)
func BenchmarkCron_Exec(b *testing.B) {
	cron := & Cron{}
	cron.Ready()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next(){
			//cron.secondWheel.GetSlot(3)
			cron.Exec(3, func() {

			})
			cron.Exec(63, func() {

			})
			cron.Exec(3663, func() {

			})
		}
	})
}

func TestCron_Exec(t *testing.T) {
	cron := & Cron{}
	cron.Ready()
	r := cron.Exec(61, func() {

	})
	if r.pos != 1{
		t.Error("Wrong slot position",r.pos)
	}
	if cron.MinuteWheel.slot[r.pos] == nil{
		t.Error("Not put into minute slot correctly")
	}
	//if cron.secondWheel.slot[1] == nil{
	//	t.Error("Not put into second slot")
	//}
	r2 := cron.Exec(33, func() {

	})

	if r2.pos != 33 {
		t.Error("Wrong slot position",r2.pos)
	}

	if cron.SecondWheel.slot[r2.pos] == nil{
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
