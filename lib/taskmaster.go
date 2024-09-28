package lib

import (
	"sync"

	webview "github.com/webview/webview_go"
)

type TaskMaster struct {
	wg *sync.WaitGroup
	wv *webview.WebView
}

func NewTaskmaster(wv *webview.WebView) *TaskMaster {
	var wg sync.WaitGroup
	return &TaskMaster{
		wg: &wg,
		wv: wv,
	}
}

func (c *TaskMaster) Task(name string, handle func(s string)) {
	c.wg.Add(1)
	(*c.wv).Bind(name, func(s string) {
		defer c.wg.Done()
		handle(s)
	})
}

func (c *TaskMaster) WaitForAllTasks() {
	c.wg.Wait()
}
