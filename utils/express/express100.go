package express

import (
	"fmt"
	"github.com/bysir-zl/bygo/util"
	"github.com/bysir-zl/bygo/util/http_util"
)

const (
	KuaiDi100Api = "http://www.kuaidi100.com/applyurl?key=%s&com=%s&nu=%s"
	KuaiDi100Key = "ad7eab3fb63b2589"
)

func GetExpressInfoBy100(come,number string) (rsp string, err error) {
	url := fmt.Sprintf(KuaiDi100Api, KuaiDi100Key, come, number)

	_, rsp, err = http_util.Get(url, util.OrderKV{}, nil)
	if err != nil {
		return
	}

	return
}
