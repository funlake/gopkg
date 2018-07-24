##### Intro
Simple timewheel pkg implement by golang
###### Install
go get github.com/funlake/timewheel
###### Import
import github.com/funlake/timewheel
```
cron := &timewheel.Cron{}
cron.Exec(3,func(){
    //do anything in each 3 seconds
})
```


