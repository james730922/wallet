package workerpool

import "github.com/james730922/wallet/service/internal/thirdparty/logger"

func newWorker(id int, dispatcher *dispatcher) *worker {
	return &worker{
		id:             id,
		quit:           make(chan bool),
		taskDepositary: dispatcher,
	}
}

type worker struct {
	id             int         // 工人編號
	isLog          bool        // 是否開啟紀錄
	quit           chan bool   // 關閉執行緒
	taskDepositary *dispatcher // 工作倉庫
}

func (w *worker) start() {
	go func() {
		for {
			select {
			case task := <-w.taskDepositary.taskStorage:
				err := w.processTask(task)

				if w.isLog {
					if err != nil {
						logger.ApLog().Debugf("worker #%d is done with error:%s", w.id, err)
					} else {
						logger.ApLog().Debugf("worker #%d is done with no error", w.id)
					}
				}

			case <-w.quit:
				return
			}
		}
	}()
}

func (w *worker) close() {
	w.quit <- true
}

func (w *worker) processTask(task Task) error {
	err := task.Exec()
	return err
}

func (w *worker) log(ok bool) {
	w.isLog = ok
}
