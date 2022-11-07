package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
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
	group := sync.WaitGroup{}

	for i := 1; i <= 4; i++ {
		group.Add(1)
		go (func(index int) {
			tmp, _ := template.New("url").Parse("http://api.linkprice.com/affiliate/translist.php?a_id=&auth_key=&yyyymmdd={{.Date}}&page={{.Page}}&per_page=1000")

			b := bytes.Buffer{}

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
			data3 := map[string]int{}

			for _, order := range res.OrderList {
				data[order.UserId] = append(data[order.UserId], order)
				if data3[order.OCd] == 0 {
					data3[order.OCd] = 1
				} else {
					data3[order.OCd] = data3[order.OCd] + 1
				}
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

			for k, v := range data3 {
				if v > 1 {
					fmt.Println("==========================================")
					fmt.Println("중복된 order_code : ", k, " count : ", v)
				}
			}

			defer group.Done()
			defer resp.Body.Close()
		})(i)
	}

	group.Wait()
}
