package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

type value struct {
	index int
	date  string
	size  int
}

var (
	pageSize = 1000
	dateList = map[string]int{"202207": 0, "202208": 0, "202209": 0, "202210": 0, "202211": 0}
)

func main() {
	f, _ := os.Create("trace.out")
	trace.Start(f)

	now := time.Now()

	sizeRequest()
	mainRequest()

	fmt.Println("time : ", time.Since(now))

	defer f.Close()
	defer trace.Stop()
}

func sizeRequest() {
	ch := make(chan *value)
	fin := make(chan *value)

	for k := range dateList {
		go (func(c <-chan *value, f chan<- *value) {
			v := <-c
			data := getData(v)

			defer (func() {
				fin <- &value{0, v.date, (data.ListCount / pageSize) + 1}
			})()
		})(ch, fin)

		ch <- &value{1, k, 1}
	}

	for i := 0; i < len(dateList); i++ {
		value := <-fin
		dateList[value.date] = value.size
	}
}

func mainRequest() {
	ch := make(chan *value)
	fin := make(chan int)

	for k := range dateList {
		for index := 1; index <= getSumPage(k); index++ {
			go request(ch, fin)
			ch <- &value{index, k, pageSize}
		}
	}

	for k := range dateList {
		for index := 1; index <= getSumPage(k); index++ {
			<-fin
		}
	}
}

func getSumPage(key string) int {
	sumPage := 0
	for k, v := range dateList {
		if k == key {
			sumPage += v
		}
	}
	return sumPage
}

func request(ch <-chan *value, fin chan<- int) {
	data := <-ch

	res := getData(data)

	orders := make(map[string][]linkData)

	for _, order := range res.OrderList {
		orders[order.OCd] = append(orders[order.OCd], order)
	}

	for k, v := range orders {
		for _, order := range v {
			if order.OCd == "2305559427" {
				fmt.Printf("user_id : %s, m_id : %s, trlog_id : %s, order_code : %s, p_code : %s, p_name : %s\ncount : %s, sales : %d, commission : %d, status : %s, date: %s comment : %s \n",
					k, order.MId, order.TrlogId, order.OCd, order.PCd, order.PNm, order.ItCnt, order.Sales, order.Commission, order.Status, order.Yyyymmdd+order.Hhmiss, order.TransComment)
			}
		}
	}

	defer (func() {
		fin <- 1
	})()
}

func getData(data *value) response {
	resp, _ := http.Get(getURL(data))
	body, _ := io.ReadAll(resp.Body)

	str := string(body)
	res := response{}

	json.Unmarshal([]byte(str), &res)

	return res
}

func getURL(data *value) string {
	tmp, _ := template.New("url").Parse("http://api.linkprice.com/affiliate/translist.php?a_id=&auth_key=&yyyymmdd={{.Date}}&page={{.Page}}&per_page={{.PerPage}}")

	b := bytes.Buffer{}
	tmp.Execute(&b, map[string]string{"Page": strconv.Itoa(data.index), "Date": data.date, "PerPage": strconv.Itoa(data.size)})

	return b.String()
}
