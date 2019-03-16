
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:45</date>
//</624450121391673344>

package storage

//我们要使用的DB模式。实际/当前数据库架构可能不同
//直到运行迁移。
const CurrentDbSchema = DbSchemaHalloween

//曾经有一段时间我们根本没有模式。
const DbSchemaNone = ""

//“纯度”是我们与Swarm 0.3.5一起发布的第一个级别数据库的正式模式。
const DbSchemaPurity = "purity"

//“万圣节”在这里是因为我们有一个螺丝钉在垃圾回收索引。
//因此，我们必须重建gc索引以消除错误的
//这需要很长的时间。此模式用于记账，
//所以重建索引只运行一次。
const DbSchemaHalloween = "halloween"

