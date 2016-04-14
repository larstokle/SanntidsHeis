package orderLog

import(
	."globals"
	//"io/ioutil"
	"os"
)	


 	

func check(e error) {
    if e != nil {
        panic(e)
    }
}


func WriteInternalOrders(orders []byte) {

	f, err := os.Create("InternalLog.txt")
	check(err)
	defer f.Close()
	ret,err := f.Write(orders)
	_ = ret
	check(err)
}


func FindInternalOrders() []byte{
	f, err := os.Open("InternalLog.txt")
	if(os.IsNotExist(err)){
		return make([]byte,N_FLOORS)
	}
	check(err)
	orderByte := make([]byte,N_FLOORS)
	_,err = f.Read(orderByte)
	f.Close()
	return orderByte
}