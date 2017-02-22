package logs

import (
	"testing"

	"github.com/zhangweilun/logs"
)

/**
* 
* @author willian
* @created 2017-02-19 09:38
* @email 18702515157@163.com  
**/

func TestLogs(t *testing.T) {
	logs.Init("/logs")
	logs.Debug("s","s")
}