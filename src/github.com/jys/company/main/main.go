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
	//dateList = map[string]int{"202201": 0, "202202": 0, "202203": 0, "202204": 0, "202205": 0, "202206": 0}
	dateList = map[string]int{"202211": 0}
	//dateList = map[string]int{"202207": 0, "202208": 0, "202209": 0, "202210": 0, "202211": 0}
	//dateList = map[string]int{"202201": 0, "202202": 0, "202203": 0, "202204": 0, "202205": 0, "202206": 0, "202207": 0, "202208": 0, "202209": 0, "202210": 0, "202211": 0}
	//dateList = map[string]int{"20220901": 0, "20220909": 0, "20221001": 0, "20221009": 0, "20221101": 0, "20221109": 0}
)

type data struct {
	s   int
	f   int
	fs  int
	n   int
	nn  int
	ns  int
	nf  int
	sf  int
	sfs int
	nfs int
	fsf int
	ff  int
	ss  int
	ot  int
}

func (d *data) add(link []linkprice.LinkData) {
	if len(link) == 1 {
		for _, order := range link {
			switch order.Status {
			case "210":
				d.s++
			case "310":
				d.f++
			case "300":
				d.fs++
			case "100":
				d.n++
			}
		}
		return
	} else if len(link) == 2 {
		if contains([]string{link[0].Status, link[1].Status}, "100") && contains([]string{link[0].Status, link[1].Status}, "210") {
			d.ns++
		} else if contains([]string{link[0].Status, link[1].Status}, "100") && contains([]string{link[0].Status, link[1].Status}, "300") {
			d.nfs++
		} else if contains([]string{link[0].Status, link[1].Status}, "100") && contains([]string{link[0].Status, link[1].Status}, "310") {
			d.nf++
		} else if contains([]string{link[0].Status, link[1].Status}, "210") && contains([]string{link[0].Status, link[1].Status}, "300") {
			d.sfs++
		} else if contains([]string{link[0].Status, link[1].Status}, "210") && contains([]string{link[0].Status, link[1].Status}, "310") {
			d.sf++
		} else if contains([]string{link[0].Status, link[1].Status}, "300") && contains([]string{link[0].Status, link[1].Status}, "310") {
			d.fsf++
		} else if link[0].Status == "310" && link[1].Status == "310" {
			d.ff++
		} else if link[0].Status == "100" && link[1].Status == "100" {
			d.nn++
		} else if link[0].Status == "210" && link[1].Status == "210" {
			d.ss++
		}
		return
	}
	for _, order := range link {
		fmt.Printf("[1] - m_id : %s, u_id : %s, o_cd : %s, p_cd : %s, p_nm : %s, cnt : %s\n[2] - amt : %d, commission : %d, status : %s, date : %s, mebership : %s, credate : %s, c_cd : %s, pgm_id : %s, pgm_name : %s\n",
			order.MId, order.UserId, order.OCd, order.PCd, order.PNm, order.ItCnt, order.Sales, order.Commission, order.Status, order.Yyyymmdd+order.Hhmiss, order.MemberShipId, order.CreateTimeStamp, order.Ccd, order.AppliedPgmId, order.PgmName)
	}
	d.ot += len(link)
}

func contains(elems []string, v string) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}

func main() {
	f, _ := os.Create("trace.out")
	trace.Start(f)

	time.Sleep(10 * time.Millisecond)

	total := 0

	//market := []string{}

	for _, v := range *linkprice.GetRequest(&linkprice.Request{PageSize: pageSize, DateList: dateList}) {
		//if !contains(market, v[0].MId) {
		//	market = append(market, v[0].MId)
		//}
		d[v[0].MId].add(v)
		total += len(v)
		for _, order := range v {
			//if len(v) > 1 && contains([]string{"2315151421", "12327689169", "8154429078792692", "203604415", "12222281597", "2006399345"}, order.OCd) {
			if contains([]string{"2036044415", "8154429078792692"}, order.OCd) {
				//if len(v) == 1 {
				fmt.Printf("[1] - m_id : %s, u_id : %s, o_cd : %s, p_cd : %s, p_nm : %s, cnt : %s\n[2] - amt : %d, commission : %d, status : %s, date : %s, mebership : %s, credate : %s, c_cd : %s, pgm_id : %s, pgm_name : %s\n",
					order.MId, order.UserId, order.OCd, order.PCd, order.PNm, order.ItCnt, order.Sales, order.Commission, order.Status, order.Yyyymmdd+order.Hhmiss, order.MemberShipId, order.CreateTimeStamp, order.Ccd, order.AppliedPgmId, order.PgmName)
			}
		}
		//}
		//}
		//}
		//fmt.Printf("user_id : %s, m_id : %s, trlog_id : %s, order_code : %s, p_code : %s, p_name : %s\ncount : %s, sales : %d, commission : %d, status : %s, date: %s comment : %s \n",
		//	order.UserId, order.MId, order.TrlogId, order.OCd, order.PCd, order.PNm, order.ItCnt, order.Sales, order.Commission, order.Status, order.Yyyymmdd+order.Hhmiss, order.TransComment)
		//	order.UserId, order.MId, order.TrlogId, order.OCd, order.PCd, order.PNm, order.ItCnt, order.Sales, order.Commission, order.Status, order.Yyyymmdd+order.Hhmiss, order.TransComment)
		//}
		//if len(v) == 1 && order.MId == "gmarket" && order.Status == "310" {
		//	if !contains([]string{"auction", "gmarket", "aliexpress"}, order.MId) {
		//}
	}

	time.Sleep(10 * time.Millisecond)

	for k, v := range d {
		fmt.Println(k)
		fmt.Println("========================================")
		fmt.Printf("구매확정 : %d\n취소확정 : %d\n취소대기 : %d\n결제 : %d\n결제 + 구매확정 : %d\n결제 + 취소확정 : %d\n구매확정 + 취소확정 : %d\n구매확정 + 취소대기 : %d\n일반 + 취소대기 : %d\n취소확정 + 취소대기 : %d\n결제 + 결제 : %d\n취소확정 + 취소확정 : %d\n구매확정 + 구매확정 : %d\n예외건 : %d\n",
			v.s, v.f, v.fs, v.n, v.ns, v.nf, v.sf, v.sfs, v.nfs, v.fsf, v.nn, v.ff, v.ss*2, v.ot)
		fmt.Println("========================================")
	}

	//fmt.Println(market)

	fmt.Println("total : ", total)

	myTotal := 0
	for _, v := range d {
		myTotal += v.s + v.f + v.fs + v.n + v.ns*2 + v.nf*2 + v.sf*2 + v.sfs*2 + v.nfs*2 + v.fsf*2 + v.nn*2 + v.ff*2 + v.ss*2 + v.ot
	}
	fmt.Println("myTotal : ", myTotal)

	defer f.Close()
	defer trace.Stop()
}
