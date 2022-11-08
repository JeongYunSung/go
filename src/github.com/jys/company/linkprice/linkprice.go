package linkprice

import (
	"bytes"
	"encoding/json"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"io"
	"net/http"
	"strconv"
	"text/template"
	"time"
	"ui"
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
	orders map[string][]linkData
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
	GetRequest(req *Request) *map[string][]linkData
}

func GetRequest(req *Request) *map[string][]linkData {
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
		value := <-fin
		req.DateList[value.date] = value.size
	}
}

func mainRequest(v *view, req *Request) *map[string][]linkData {
	ch := make(chan *value)
	fin := make(chan finish)

	result := make(map[string][]linkData)

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
					result[order.OCd] = append(result[order.OCd], orders...)
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

	orders := make(map[string][]linkData)

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
	tmp, _ := template.New("url").Parse("http://api.linkprice.com/affiliate/translist.php?a_id=A100675064&auth_key=3491b643ffeb0d893e34b5dc7f714964&yyyymmdd={{.Date}}&page={{.Page}}&per_page={{.PerPage}}")

	b := bytes.Buffer{}
	tmp.Execute(&b, map[string]string{"Page": strconv.Itoa(data.index), "Date": data.date, "PerPage": strconv.Itoa(data.size)})

	return b.String()
}
