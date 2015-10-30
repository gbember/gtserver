package word

import (
	"sync"

	"github.com/gbember/gt/logger"
	"github.com/gbember/gt/module"
)

type Worder interface {
	Init() interface{}
	Run(interface{})
	Close(interface{})
}

var (
	wds []*word = make([]*word, 0, 20)
	mut sync.Mutex
)

type words []*word

type word struct {
	restartNum    int
	maxRestartNum int
	worder        Worder
	v             interface{}
}

func RegisterModule() {
	module.Register(words(wds))
}

//注册一个worder
func RegisterWord(worder Worder, maxRestartNum int) {
	mut.Lock()
	defer mut.Unlock()
	w := new(word)
	if maxRestartNum < 1 {
		maxRestartNum = 1
	}
	w.maxRestartNum = maxRestartNum
	w.worder = worder
	wds = append(wds, w)
}

func (wds words) OnInit() {
	for i := 0; i < len(wds); i++ {
		wds[i].init()
	}
	logger.Info("word start...")
}

func (wds words) Run(chan bool) {
	for i := 0; i < len(wds); i++ {
		go wds[i].run()
	}
}

func (wds words) OnDestroy() {
	for i := 0; i < len(wds); i++ {
		wds[i].close()
	}
}

func (w *word) init() {
	w.v = w.worder.Init()
}

func (w *word) run() {
	//是否重启
	defer func() {
		if x := recover(); x != nil {
			w.restartNum++
			if w.restartNum <= w.maxRestartNum {
				go w.run()
				logger.Error("word restart: %v", x)
			} else {
				logger.Error("word exit: %v", x)
			}
		}
	}()
	w.worder.Run(w.v)
}

func (w *word) close() {
	w.worder.Close(w.v)
}
