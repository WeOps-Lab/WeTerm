package index

import (
	"weterm/model"
	"weterm/pages/example"
)

var installMenu = []MenuItem{
	{
		Name: "单机版-未实现",
		Action: func(bs *model.AppModel) {
			example.SetUpFormSamplePage(bs)
			bs.CorePages.SwitchToPage("form_sample")
		},
	},
	{
		Name: "标准版(3节点)-未实现",
		Action: func(bs *model.AppModel) {
			example.SetUpFormSamplePage(bs)
			bs.CorePages.SwitchToPage("form_sample")
		},
	},
	{
		Name: "高可用版(7节点)-未实现",
		Action: func(bs *model.AppModel) {
			example.SetUpFormSamplePage(bs)
			bs.CorePages.SwitchToPage("form_sample")
		},
	},
	{
		Name: "自定义安装-未实现",
		Action: func(bs *model.AppModel) {
			example.SetUpFormSamplePage(bs)
			bs.CorePages.SwitchToPage("form_sample")
		},
	},
}
