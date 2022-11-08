package main

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"linkprice"
	"lnlist"
	"os"
	"runtime/trace"
	"time"
)

var (
	pageSize = 1000
	//dateList = map[string]int{"202210": 0, "202211": 0}
	dateList = map[string]int{"202207": 0, "202208": 0, "202209": 0, "202210": 0, "202211": 0}
)

func main() {
	f, _ := os.Create("trace.out")
	trace.Start(f)

	model := lnlist.Model{List: list.New(nil, list.NewDefaultDelegate(), 0, 0)}
	model.List.Title = "링크프라이스"

	linkprice.GetRequest(&model, &linkprice.Request{PageSize: pageSize, DateList: dateList})

	time.Sleep(1000 * time.Millisecond)

	tea.NewProgram(model, tea.WithAltScreen()).Start()

	defer f.Close()
	defer trace.Stop()
}
