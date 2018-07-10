package cmd

import (
	"fmt"
        "github.com/viphxin/xingo/utils"
)

type TimerCommand struct {
}

func NewTimerCommand() *TimerCommand{
	return &TimerCommand{}
}
func (this *TimerCommand)Name()string{
	return "timer"
}

func (this *TimerCommand)Help()string{
	return fmt.Sprintf("timer:\r\n" +
		"----------- count: 剩余定时任务数\r\n")
}

func (this *TimerCommand)Run(args []string) string{
	if len(args) == 0{
		return this.Help()
	}else{
		switch args[0] {
		case "count":
			return fmt.Sprintf("%d", utils.GlobalObject.GetSafeTimer().TotalCnt())
		default:
			return "未实现"
		}
	}
	return "OK"
}
