
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:42</date>
//</624450107735019520>


package params

//Gastable组织不同以太坊阶段的天然气价格。
type GasTable struct {
	ExtcodeSize uint64
	ExtcodeCopy uint64
	ExtcodeHash uint64
	Balance     uint64
	SLoad       uint64
	Calls       uint64
	Suicide     uint64

	ExpByte uint64

//当
//退款账户是这样的
//不存在。这种逻辑是相似的
//打电话。可以留空。零均值
//不收费。
	CreateBySuicide uint64
}

//包含不同以太坊阶段天然气价格的变量。
var (
//GastableHomestead包含
//宅基地阶段。
	GasTableHomestead = GasTable{
		ExtcodeSize: 20,
		ExtcodeCopy: 20,
		Balance:     20,
		SLoad:       50,
		Calls:       40,
		Suicide:     0,
		ExpByte:     10,
	}

//GastableEIP150包含天然气再定价
//EIP150相位。
	GasTableEIP150 = GasTable{
		ExtcodeSize: 700,
		ExtcodeCopy: 700,
		Balance:     400,
		SLoad:       200,
		Calls:       700,
		Suicide:     5000,
		ExpByte:     10,

		CreateBySuicide: 25000,
	}
//GastableEIP158包含
//EIP155/EIP158阶段。
	GasTableEIP158 = GasTable{
		ExtcodeSize: 700,
		ExtcodeCopy: 700,
		Balance:     400,
		SLoad:       200,
		Calls:       700,
		Suicide:     5000,
		ExpByte:     50,

		CreateBySuicide: 25000,
	}
//天然气君士坦丁堡包含天然气再定价
//君士坦丁堡时期。
	GasTableConstantinople = GasTable{
		ExtcodeSize: 700,
		ExtcodeCopy: 700,
		ExtcodeHash: 400,
		Balance:     400,
		SLoad:       200,
		Calls:       700,
		Suicide:     5000,
		ExpByte:     50,

		CreateBySuicide: 25000,
	}
)

