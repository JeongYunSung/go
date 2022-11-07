package main

import (
	"fmt"
	"linkprice"
	"os"
	"runtime/trace"
	"time"
)

var (
	pageSize = 1000
	dateList = map[string]int{"202208": 0, "202209": 0}
	//dateList = map[string]int{"202207": 0, "202208": 0, "202209": 0, "202210": 0, "202211": 0}
)

func main() {
	f, _ := os.Create("trace.out")
	trace.Start(f)

	now := time.Now()

	for k, v := range *linkprice.GetRequest(&linkprice.Request{PageSize: pageSize, DateList: dateList}) {
		if k == "2305559427" {
			for _, order := range v {
				fmt.Printf("user_id : %s, m_id : %s, trlog_id : %s, order_code : %s, p_code : %s, p_name : %s\ncount : %s, sales : %d, commission : %d, status : %s, date: %s comment : %s \n",
					k, order.MId, order.TrlogId, order.OCd, order.PCd, order.PNm, order.ItCnt, order.Sales, order.Commission, order.Status, order.Yyyymmdd+order.Hhmiss, order.TransComment)
			}
		}
	}

	fmt.Println("time : ", time.Since(now))

	defer f.Close()
	defer trace.Stop()
}
