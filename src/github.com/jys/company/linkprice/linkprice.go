package linkprice

import (
	"bytes"
	"encoding/json"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"io"
	"net/http"
	"strconv"
	"time"
	"ui"
)

type response struct {
	Result    string     `json:"result"`
	ListCount int        `json:"list_count"`
	OrderList []LinkData `json:"order_list"`
}

type LinkData struct {
	TrlogId         string `json:"trlog_id"`
	MId             string `json:"m_id"`
	OCd             string `json:"o_cd"`
	PCd             string `json:"p_cd"`
	PNm             string `json:"p_nm"`
	Ccd             string `json:"c_cd"`
	ItCnt           string `json:"it_cnt"`
	UserId          string `json:"user_id"`
	Status          string `json:"status"`
	Yyyymmdd        string `json:"yyyymmdd"`
	Hhmiss          string `json:"hhmiss"`
	Sales           int    `json:"sales"`
	Commission      int    `json:"commission"`
	PurRate         string `json:"pur_rate"`
	TransComment    string `json:"trans_comment"`
	MemberShipId    string `json:"membership_id"`
	CreateTimeStamp string `json:"create_time_stamp"`
	AppliedPgmId    string `json:"applied_pgm_id"`
	PgmName         string `json:"pgm_name"`
}

type finish struct {
	name   string
	orders map[string][]LinkData
}

type value struct {
	index int
	date  string
	size  int
}

type Request struct {
	PageSize int
	DateList map[string]int
}

type view struct {
	model   ui.Model
	program *tea.Program
}

type LinkPrice interface {
	GetRequest(req *Request) *map[string][]LinkData
}

func GetRequest(req *Request) *map[string][]LinkData {
	sizeRequest(req)

	model := ui.NewModel(getPrettyName(req))
	program := tea.NewProgram(model)

	go func(p *tea.Program) {
		p.Start()
	}(program)

	v := &view{model, program}

	return mainRequest(v, req)
}

func sizeRequest(req *Request) {
	ch := make(chan *value)
	fin := make(chan *value)

	for k := range req.DateList {
		go (func(c <-chan *value, f chan<- *value) {
			v := <-c
			data := getData(v)

			defer (func() {
				fin <- &value{0, v.date, (data.ListCount / req.PageSize) + 1}
			})()
		})(ch, fin)

		ch <- &value{1, k, 1}
	}

	for i := 0; i < len(req.DateList); i++ {
		v := <-fin
		req.DateList[v.date] = v.size
	}
}

func mainRequest(v *view, req *Request) *map[string][]LinkData {
	ch := make(chan *value)
	fin := make(chan finish)

	result := make(map[string][]LinkData)

	for k := range req.DateList {
		for index := 1; index <= getSumPage(req, k); index++ {
			go request(ch, fin)
			ch <- &value{index, k, req.PageSize}
		}
	}

	for k := range req.DateList {
		for index := 1; index <= getSumPage(req, k); index++ {
			element := <-fin
			for _, orders := range element.orders {
				for _, order := range orders {
					result[order.OCd] = append(result[order.OCd], order)
				}
			}
			v.program.Send(v.model.NextTick())
		}
	}

	time.Sleep(1 * time.Millisecond)

	v.program.Send(v.model.NextTick())

	return &result
}

func getPrettyName(req *Request) []string {
	requests := []string{"페이지 계산"}

	for k, v := range req.DateList {
		for i := 1; i <= v; i++ {
			requests = append(requests, fmt.Sprintf("%s > %d 페이지", k, i))
		}
	}

	requests = append(requests, "fin")

	return requests
}

func getSumPage(req *Request, key string) int {
	sumPage := 0
	for k, v := range req.DateList {
		if k == key {
			sumPage += v
		}
	}
	return sumPage
}

func request(ch <-chan *value, fin chan<- finish) {
	data := <-ch

	res := getData(data)

	orders := make(map[string][]LinkData)

	for _, order := range res.OrderList {
		orders[order.OCd] = append(orders[order.OCd], order)
	}

	defer (func() {
		fin <- finish{name: fmt.Sprintf("%s-%d", data.date, data.index), orders: orders}
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
	b := bytes.Buffer{}
	tmp.Execute(&b, map[string]string{"Page": strconv.Itoa(data.index), "Date": data.date, "PerPage": strconv.Itoa(data.size)})

	return b.String()
}
