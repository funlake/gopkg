package jobworker

import (
	"testing"
)

func TestDispatcher_Put(t *testing.T) {
	dispatcher := NewDispather(2,10)
	for i:=0;i<10;i++{
		dispatcher.Put(&simpleJob{})
	}
	if dispatcher.Put(&simpleJob{}){
		t.Error("Should be stucked")
	}else{
		t.Log("Ok")
	}
}

func BenchmarkDispatcher_Put(b *testing.B) {
	dispatcher := NewDispather(20000,500000)
	b.SetParallelism(100)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next(){
			dispatcher.Put(&simpleJob{})
		}
	})
}