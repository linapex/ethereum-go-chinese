
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:36</date>
//</624450083189952512>


//+构建GouuZZ

package runtime

//引信是Go-Fuzz工具的基本切入点
//
//对于有效的可分析/不可运行代码，返回1，0
//对于无效的操作码。
func Fuzz(input []byte) int {
	_, _, err := Execute(input, input, &Config{
		GasLimit: 3000000,
	})

//无效操作码
	if err != nil && len(err.Error()) > 6 && string(err.Error()[:7]) == "invalid" {
		return 0
	}

	return 1
}

