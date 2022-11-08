package linkprice

import (
	"bytes"
	"encoding/json"
	"fetch"
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"io"
	"lnlist"
	"net/http"
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

type finish struct {
	name   string
	orders []list.Item
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
	lnModel *lnlist.Model
	model   fetch.Model
	program *tea.Program
}

type LinkPrice interface {
	GetRequest(model lnlist.Model, req *Request)
}

func GetRequest(lnModel *lnlist.Model, req *Request) {
	sizeRequest(req)

	model := fetch.NewModel(getPrettyName(req))
	program := tea.NewProgram(model)

	go func(p *tea.Program) {
		p.Start()
	}(program)

	v := &view{lnModel, model, program}

	mainRequest(v, req)
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
		value := <-fin
		req.DateList[value.date] = value.size
	}
}

func mainRequest(v *view, req *Request) {
	ch := make(chan *value)
	fin := make(chan finish)

	for k := range req.DateList {
		for index := 1; index <= getSumPage(req, k); index++ {
			go request(ch, fin)
			ch <- &value{index, k, req.PageSize}
		}
	}

	for k := range req.DateList {
		for index := 1; index <= getSumPage(req, k); index++ {
			element := <-fin
			for _, order := range element.orders {
				v.lnModel.List.InsertItem(0, order)
			}
			v.program.Send(v.model.NextTick())
		}
	}

	time.Sleep(1 * time.Millisecond)

	v.program.Send(v.model.NextTick())
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

	var orders []list.Item

	for _, order := range res.OrderList {
		orders = append(orders, lnlist.Item{Subject: order.OCd, Desc: fmt.Sprintf("user_id : %s, m_id : %s, trlog_id : %s, p_code : %s, p_name : %s, count : %s, sales : %d, commission : %d, status : %s, date: %s comment : %s",
			order.UserId, order.MId, order.TrlogId, order.PCd, order.PNm, order.ItCnt, order.Sales, order.Commission, order.Status, order.Yyyymmdd+order.Hhmiss, order.TransComment)})
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
	tmp, _ := template.New("url").Parse("http://api.linkprice.com/affiliate/translist.php?a_id=&auth_key=&yyyymmdd={{.Date}}&page={{.Page}}&per_page={{.PerPage}}")

	b := bytes.Buffer{}
	tmp.Execute(&b, map[string]string{"Page": strconv.Itoa(data.index), "Date": data.date, "PerPage": strconv.Itoa(data.size)})

	return b.String()
}
