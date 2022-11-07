package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime/trace"
	"strconv"
	"text/template"
	"time"
)

type response struct {
	Result    string     `json:"result"`
	ListCount int        `json:"list_count"`
	OrderList []linkData `json:"order_list"`
}

type linkData struct {
	TrlogId      string `json:"trlog_id"`
	MId          string `json:"m_id"`
	OCd          string `json:"o_cd"`
	PCd          string `json:"p_cd"`
	PNm          string `json:"p_nm"`
	ItCnt        string `json:"it_cnt"`
	UserId       string `json:"user_id"`
	Status       string `json:"status"`
	Yyyymmdd     string `json:"yyyymmdd"`
	Hhmiss       string `json:"hhmiss"`
	Sales        int    `json:"sales"`
	Commission   int    `json:"commission"`
	PurRate      string `json:"pur_rate"`
	TransComment string `json:"trans_comment"`
}

func main() {
	f, _ := os.Create("trace.out")
	trace.Start(f)

	defer f.Close()
	defer trace.Stop()

	ch := make(chan int)
	fin := make(chan int)

	now := time.Now()

	for i := 1; i <= 4; i++ {
		go request(ch, fin)
		ch <- i
	}

	for i := 0; i < 4; i++ {
		<-fin
	}

	fmt.Println("time : ", time.Since(now))
}

func request(ch <-chan int, fin chan<- int) {
	tmp, _ := template.New("url").Parse("http://api.linkprice.com/affiliate/translist.php?a_id=&auth_key=&yyyymmdd={{.Date}}&page={{.Page}}&per_page=1000")

	b := bytes.Buffer{}

	index := <-ch

	defer (func() {
		fin <- index
	})()

	err3 := tmp.Execute(&b, map[string]string{"Page": strconv.Itoa(index), "Date": "202209"})
	if err3 != nil {
		log.Fatal(err3)
	}

	resp, err := http.Get(b.String())

	if err != nil {
		log.Fatal(err)
	}

	r, _ := io.ReadAll(resp.Body)

	str := string(r)
	res := response{}

	err2 := json.Unmarshal([]byte(str), &res)

	if err2 != nil {
		log.Fatal(err2)
	}

	data := make(map[string][]linkData)

	for _, order := range res.OrderList {
		data[order.UserId] = append(data[order.UserId], order)
	}

	fmt.Printf("%d 페이지의 %d 건의 데이터 조회 완료 \n", index, len(res.OrderList))

	for k, v := range data {
		if len(v) > 1 {
			fmt.Println("==========================================")
			for _, order := range v {
				fmt.Printf("user_id : %s, m_id : %s, trlog_id : %s, order_code : %s, p_code : %s, p_name : %s\ncount : %s, sales : %d, commission : %d, status : %s, date: %s comment : %s \n",
					k, order.MId, order.TrlogId, order.OCd, order.PCd, order.PNm, order.ItCnt, order.Sales, order.Commission, order.Status, order.Yyyymmdd+order.Hhmiss, order.TransComment)
			}
		}
	}

	fmt.Println("[JYS]중복 체크 로직")
}
