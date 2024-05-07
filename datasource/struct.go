package datasource

import "sort"

//Boolean 布尔
type Boolean struct {
}

//Decimal 小数
type Decimal struct {
	Max float64 `json:"max,omitempty"`
	Min float64 `json:"min,omitempty"`
}

//Integer 数字
type Integer struct {
	Max int64 `json:"max,omitempty"`
	Min int64 `json:"min,omitempty"`
}

//String 字符串
type String struct {
	MaxLen int32 `json:"maxLen"`
}

//Event 事件
type Event struct {
	Key        string      `json:"key"`             //标识
	Name       int32       `json:"name"`            //事件名称
	Type       int64       `json:"type"`            //事件类型
	Parameters []Parameter `json:"items,omitempty"` //参数
}

//Command 事件
type Command struct {
	Key    string      `json:"key"`           //标识
	Name   string      `json:"name"`          //事件名称
	Type   int64       `json:"type"`          //事件类型
	Inputs []Parameter `json:"req,omitempty"` //输入参数
	Outs   []Parameter `json:"res,omitempty"` //输出参数
}

//Parameter 参数
type Parameter struct {
	//type	integer($int32)
	Key     string   `json:"key,omitempty"`     //标识
	Type    int32    `json:"type,omitempty"`    //数据类型; 1 整数, 2 浮点,3 字符, 4 布尔
	Boolean *Boolean `json:"boolean,omitempty"` //布尔
	Decimal *Decimal `json:"double,omitempty"`  //浮点
	Integer *Integer `json:"integer,omitempty"` //数字
	String  *String  `json:"string,omitempty"`  //字符串
}

//ParameterSort 排序方法
type ParameterSort func(i, j *Parameter) bool

//Sort 排序
func (p ParameterSort) Sort(items []Parameter) {
	ps := &ParameterSorter{
		items: items,
		by:    p,
	}
	sort.Sort(ps)
}

//ParameterSorter 参数排序
type ParameterSorter struct {
	items []Parameter
	by    func(i, j *Parameter) bool
}

//Len 长度
func (p *ParameterSorter) Len() int {
	return len(p.items)
}

//Swap 位置交换
func (p *ParameterSorter) Swap(i, j int) {
	p.items[i], p.items[j] = p.items[j], p.items[i]
}

//Less less
func (p *ParameterSorter) Less(i, j int) bool {
	return p.by(&p.items[i], &p.items[j])
}

//Extension 扩展信息
type Extension struct {
	Type   int        `json:type`               //1：jar	2:js
	Script JavaScript `json:"script,omitempty"` //js 脚本
	Java   Java       `json:"jar,omitempty"`    //java
}

//JavaScript js
type JavaScript struct {
	Data string `json:data`
}

//Java java
type Java struct {
	URL       string `json:url`
	ClassName string `json:className`
}
