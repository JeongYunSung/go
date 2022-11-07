package linkprice

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"text/template"
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

type Request struct {
	PageSize int
	DateList map[string]int
}

type LinkPrice interface {
	GetRequest(req *Request) *map[string][]linkData
}

func GetRequest(req *Request) *map[string][]linkData {
	sizeRequest(req)
	return mainRequest(req)
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

func mainRequest(req *Request) *map[string][]linkData {
	ch := make(chan *value)
	fin := make(chan map[string][]linkData)

	result := make(map[string][]linkData)

	for k := range req.DateList {
		for index := 1; index <= getSumPage(req, k); index++ {
			go request(ch, fin)
			ch <- &value{index, k, req.PageSize}
		}
	}

	for k := range req.DateList {
		for index := 1; index <= getSumPage(req, k); index++ {
			for k, v := range <-fin {
				result[k] = append(result[k], v...)
			}
		}
	}

	return &result
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

func request(ch <-chan *value, fin chan<- map[string][]linkData) {
	data := <-ch

	res := getData(data)

	orders := make(map[string][]linkData)

	for _, order := range res.OrderList {
		orders[order.OCd] = append(orders[order.OCd], order)
	}

	defer (func() {
		fin <- orders
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
