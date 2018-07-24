package timer

func NewTimer() *timer {
	return &timer{}
}
type timer struct{
	secondWheel *timeWheel
	minuteWheel *timeWheel
	hourWheel   *timeWheel
}
func (schedular *timer) initialize() {
	schedular.secondWheel = &timeWheel{
		interval:1,
		maxpos:60,
		curpos:0,
		slot:make(map[int] chan *SlotItem),
	}
	schedular.minuteWheel = &timeWheel{
		interval:60,
		maxpos:60,
		curpos:0,
		slot:make(map[int] chan *SlotItem),
	}
	schedular.hourWheel = &timeWheel{
		interval:3600,
		maxpos:24,
		curpos:0,
		slot:make(map[int] chan *SlotItem),
	}
}
func (schedular *timer) SetInterval(timeout int,fun func()) *SlotItem{
	if timeout <= schedular.minuteWheel.interval{
		return schedular.secondWheel.SetInterval(timeout, fun)
	}
	if timeout > schedular.minuteWheel.interval && timeout <= schedular.hourWheel.interval {
		//用指针方能正确赋予sm.left值
		var sm *SlotItem
		sm = schedular.minuteWheel.SetInterval(timeout,func(){
				//累计剩余left秒数，判断是否需要cross本次定时tick
				if sm.left > 0 &&  (sm.left + schedular.minuteWheel.interval <= timeout){
					sm.left = sm.left + schedular.minuteWheel.interval
					//logs.Info(sm.left)
					//go to next round
					return
				}
				if timeout%schedular.minuteWheel.interval > 0 {
					var ss *SlotItem
					//sm.left意思是定时任务离走完还剩余的秒数
					//mention : what happen if timeout - sm.left == 0 ?
					//需要修改Timewheel,先执行回调，再更新步数
					ss = schedular.SetInterval((timeout - sm.left) % schedular.minuteWheel.interval, func() {
						if sm != nil {
							//记录离此次分钟定时器结束还有多少秒
							//后续的定时间隔计算需考虑此sm.left
							if schedular.secondWheel.curpos == 0 {
								//重置left,回到宇宙形成之初
								sm.left = 0
							} else {
								//sm.left = (schedular.secondWheel.maxpos - schedular.secondWheel.curpos + 1) * schedular.secondWheel.interval
								sm.left = schedular.minuteWheel.interval - ((timeout - sm.left) % schedular.minuteWheel.interval)
							}
							schedular.secondWheel.StopInterval(ss)
						}
						fun()
					})
				} else {
					fun()
				}

		})
		return sm
	}
	if timeout > schedular.hourWheel.interval{
		var sh *SlotItem
		sh = schedular.hourWheel.SetInterval(timeout, func() {
			if sh.left > 0 &&  (sh.left  + schedular.hourWheel.interval  <= timeout){
				sh.left = sh.left + schedular.hourWheel.interval
				//logs.Info(sh.left)
				//go to next round
				return
			}
			if timeout%schedular.hourWheel.interval > 0 {
				var sm *SlotItem
				sm = schedular.SetInterval((timeout - sh.left) % schedular.hourWheel.interval,func() {
					if sm != nil {
						if schedular.minuteWheel.curpos == 0 {
							//重置left,回到宇宙形成之初
							sm.left = 0
						} else {
							//sh.left = (schedular.minuteWheel.maxpos - schedular.minuteWheel.curpos + 1) * schedular.minuteWheel.interval
							sh.left = schedular.hourWheel.interval - ((timeout - sh.left) % schedular.hourWheel.interval)
						}
						schedular.minuteWheel.StopInterval(sm)
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
func (schedular *timer) StopInterval(item * SlotItem)  {
	if item.timeout <= schedular.minuteWheel.interval {
		schedular.secondWheel.StopInterval(item)
	}
	if item.timeout > schedular.minuteWheel.interval && item.timeout <= schedular.hourWheel.interval {
		schedular.minuteWheel.StopInterval(item)
	}
	if item.timeout > schedular.hourWheel.interval{
		schedular.hourWheel.StopInterval(item)
	}
}
func(schedular *timer)Ready(){
	schedular.initialize()
	go schedular.secondWheel.Start(func() {
		//every minutes pass
		schedular.minuteWheel.UpdatePos(func() {
			//every hours pass
			schedular.hourWheel.UpdatePos(func() {
				//nothing to do
			})
			schedular.hourWheel.Invoke()
		})
		schedular.minuteWheel.Invoke()
	})
}
