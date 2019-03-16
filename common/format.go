
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:34</date>
//</624450073501110272>


package common

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

//prettyDuration是时间的一个漂亮打印版本。Duration值会减少
//格式文本表示中不必要的精度。
type PrettyDuration time.Duration

var prettyDurationRe = regexp.MustCompile(`\.[0-9]+`)

//string实现了stringer接口，允许漂亮地打印持续时间
//数值四舍五入为三位小数。
func (d PrettyDuration) String() string {
	label := fmt.Sprintf("%v", time.Duration(d))
	if match := prettyDurationRe.FindString(label); len(match) > 4 {
		label = strings.Replace(label, match, match[:4], 1)
	}
	return label
}

//prettyage是时间的一个漂亮打印版本。持续时间值
//最大值为一个最重要的单位，包括天/周/年。
type PrettyAge time.Time

//AgeUnits是一个非常适合打印使用的单位列表。
var ageUnits = []struct {
	Size   time.Duration
	Symbol string
}{
	{12 * 30 * 24 * time.Hour, "y"},
	{30 * 24 * time.Hour, "mo"},
	{7 * 24 * time.Hour, "w"},
	{24 * time.Hour, "d"},
	{time.Hour, "h"},
	{time.Minute, "m"},
	{time.Second, "s"},
}

//string实现了stringer接口，允许漂亮地打印持续时间
//四舍五入到最重要的时间单位。
func (t PrettyAge) String() string {
//计算时差并处理0角箱
	diff := time.Since(time.Time(t))
	if diff < time.Second {
		return "0"
	}
//返回前累计3个分量的精度
	result, prec := "", 0

	for _, unit := range ageUnits {
		if diff > unit.Size {
			result = fmt.Sprintf("%s%d%s", result, diff/unit.Size, unit.Symbol)
			diff %= unit.Size

			if prec += 1; prec >= 3 {
				break
			}
		}
	}
	return result
}

