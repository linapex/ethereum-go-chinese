
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:42</date>
//</624450108024426496>


/*
包rlp实现rlp序列化格式。

RLP（递归线性前缀）的目的是任意编码
嵌套的二进制数据数组，rlp是使用的主要编码方法
在以太坊中序列化对象。RLP的唯一目的是编码
结构；编码特定的原子数据类型（例如字符串、整数、
浮点数）保留到高阶协议；在以太坊整数中
必须用不带前导零的大尾数二进制形式表示
（因此使整数值为零等于空字节
数组）。

RLP值由类型标记区分。类型标记位于
输入流中的值，并定义字节的大小和类型
接下来就是这样。
**/

package rlp

