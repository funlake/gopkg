package timewheel

import (
)

type Cron struct{
	SecondWheel *TimeWheel
	MinuteWheel *TimeWheel
	HourWheel   *TimeWheel
}
func (schedular *Cron) initialize() {
	schedular.SecondWheel = &TimeWheel{
		interval:1,
		maxpos:60,
		curpos:0,
		slot:make(map[int] chan *SlotItem),
	}
	schedular.MinuteWheel = &TimeWheel{
		interval:60,
		maxpos:60,
		curpos:0,
		slot:make(map[int] chan *SlotItem),
	}
	schedular.HourWheel = &TimeWheel{
		interval:3600,
		maxpos:24,
		curpos:0,
		slot:make(map[int] chan *SlotItem),
	}
}
func (schedular *Cron) Exec(timeout int,fun func()) *SlotItem{
	if timeout <= schedular.MinuteWheel.interval{
		return schedular.SecondWheel.SetInterval(timeout, fun)
	}
	if timeout > schedular.MinuteWheel.interval && timeout <= schedular.HourWheel.interval {
		//用指针方能正确赋予sm.left值
		var sm *SlotItem
		sm = schedular.MinuteWheel.SetInterval(timeout,func(){
				//累计剩余left秒数，判断是否需要cross本次定时tick
				if sm.left > 0 &&  (sm.left + schedular.MinuteWheel.interval <= timeout){
					sm.left = sm.left + schedular.MinuteWheel.interval
					//logs.Info(sm.left)
					//go to next round
					return
				}
				if timeout%schedular.MinuteWheel.interval > 0 {
					var ss *SlotItem
					//sm.left意思是定时任务离走完还剩余的秒数
					//mention : what happen if timeout - sm.left == 0 ?
					//需要修改Timewheel,先执行回调，再更新步数
					ss = schedular.Exec((timeout - sm.left) % schedular.MinuteWheel.interval, func() {
						if sm != nil {
							//记录离此次分钟定时器结束还有多少秒
							//后续的定时间隔计算需考虑此sm.left
							if schedular.SecondWheel.curpos == 0 {
								//重置left,回到宇宙形成之初
								sm.left = 0
							} else {
								//sm.left = (schedular.SecondWheel.maxpos - schedular.SecondWheel.curpos + 1) * schedular.SecondWheel.interval
								sm.left = schedular.MinuteWheel.interval - ((timeout - sm.left) % schedular.MinuteWheel.interval)
							}
							schedular.SecondWheel.StopInterval(ss)
						}
						fun()
					})
				} else {
					fun()
				}

		})
		return sm
	}
	if timeout > schedular.HourWheel.interval{
		var sh *SlotItem
		sh = schedular.HourWheel.SetInterval(timeout, func() {
			if sh.left > 0 &&  (sh.left  + schedular.HourWheel.interval  <= timeout){
				sh.left = sh.left + schedular.HourWheel.interval
				//logs.Info(sh.left)
				//go to next round
				return
			}
			if timeout%schedular.HourWheel.interval > 0 {
				var sm *SlotItem
				sm = schedular.Exec((timeout - sh.left) % schedular.HourWheel.interval,func() {
					if sm != nil {
						if schedular.MinuteWheel.curpos == 0 {
							//重置left,回到宇宙形成之初
							sm.left = 0
						} else {
							//sh.left = (schedular.MinuteWheel.maxpos - schedular.MinuteWheel.curpos + 1) * schedular.MinuteWheel.interval
							sh.left = schedular.HourWheel.interval - ((timeout - sh.left) % schedular.HourWheel.interval)
						}
						schedular.MinuteWheel.StopInterval(sm)
					}
					fun()
				})
			} else {
				fun()
			}
		})
		return sh
	}
	return nil
}
func(schedular *Cron)Ready(){
	schedular.initialize()
	go schedular.SecondWheel.Start(func() {
		//every minutes pass
		schedular.MinuteWheel.UpdatePos(func() {
			//every hours pass
			schedular.HourWheel.UpdatePos(func() {
				//nothing to do
			})
			schedular.HourWheel.Invoke()
		})
		schedular.MinuteWheel.Invoke()
	})
}
