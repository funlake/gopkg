package jobworker

import "testing"

func TestDispatcher_Put(t *testing.T) {
	dispatcher := NewRequestDispather(10,10)
	dispatcher.Put(&simpleJob{})
}

func BenchmarkDispatcher_Put(b *testing.B) {
	dispatcher := NewRequestDispather(100,50)
	b.SetParallelism(10)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next(){
			dispatcher.Put(&simpleJob{})
		}
	})
}