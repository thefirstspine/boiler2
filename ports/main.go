package ports

import (
	"fmt"
	"net"
	"github.com/thefirstspine/boiler2/utils"
)

func GetFirstFreePort(start int, end int, exclude []int) int {
	for i := start; i <= end; i++ {
		if !utils.IntInSlice(i, exclude) && TestPortIsFree(i) {
			return i
		}
	}
	return -1
}

func TestPortIsFree(port int) bool {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
  defer l.Close()

  return err == nil
}
