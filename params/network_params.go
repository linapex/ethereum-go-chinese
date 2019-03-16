
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:42</date>
//</624450107772768256>


package params

//这些是需要在客户端之间保持不变的网络参数，但是
//不一定与共识有关。

const (
//BloomBitsBlocks是单个BloomBit部分向量的块数。
//包含在服务器端。
	BloomBitsBlocks uint64 = 4096

//BloomBitsBlocksClient是单个BloomBit部分向量的块数。
//在轻型客户端包含
	BloomBitsBlocksClient uint64 = 32768

//BloomConfirms是在Bloom部分
//考虑可能是最终的，并计算其旋转位。
	BloomConfirms = 256

//chtfrequenceclient是在客户端创建cht的块频率。
	CHTFrequencyClient = 32768

//chtfrequencyserver是在服务器端创建cht的块频率。
//最终，这可以与客户端版本合并，但这需要
//完整的数据库升级，所以应该留一段合适的时间。
	CHTFrequencyServer = 4096

//BloomTrieFrequency是在两个对象上创建BloomTrie的块频率。
//服务器/客户端。
	BloomTrieFrequency = 32768

//HelperTrieConfirmations是预期客户端之前的确认数
//提供所给的帮助者。
	HelperTrieConfirmations = 2048

//HelperTrieProcessConfirmations是HelperTrie之前的确认数
//生成
	HelperTrieProcessConfirmations = 256
)

