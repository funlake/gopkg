package utils

import (
	"github.com/funlake/gopkg/utils/log"
	"unsafe"
)

func RoutineRecover(){
	if err := recover() ; err != nil{
		log.Error("%s",err)
	}
}
func Wrap(fun func()){
	go func() {
		defer RoutineRecover()
		fun()
	}()
}
func StrToByte(s string) []byte{
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0],x[1],x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func ByteToStr(s []byte) string{
	return *(*string)(unsafe.Pointer(&s))
}