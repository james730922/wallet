package workerpool

import (
	"github.com/james730922/wallet/service/internal/utils/errs"
)

func newDispatcher(buf int) *dispatcher {
	return &dispatcher{
		taskStorage: make(chan Task, buf),
	}
}

type dispatcher struct {
	taskStorage chan Task // 工作儲存庫
}

func (d *dispatcher) receiveTask(task Task) error {
	select {
	case d.taskStorage <- task:
		return nil
	default:
		return errs.BufSpaceFilled
	}
}
