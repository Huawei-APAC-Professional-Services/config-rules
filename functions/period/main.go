package main

import (
	"fmt"

	"github.com/Huawei-APAC-Professional-Services/config-rules/event"
	huaweicontext "huaweicloud.com/go-runtime/go-api/context"
)

func handler(event event.ConfigEvent, ctx huaweicontext.RuntimeContext) (interface{}, error) {
	fmt.Println(event.DomainId)
	return "test", nil
}

func main() {
	huaweicontext.Register(handler)
}
