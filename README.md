# gowallet

A bitcoin wallet application written in golang. 
Supports random wallet and brain wallet.

<b>The brain wallet uses a secret phrase and a salt phrase</b> to generate the private key.<br/>

<b>Secret phrase at least 16 characters</b>, containing uppercase letters, lowercase letters, numbers, and special characters.<br/>
<b>Salt phrase at least 6 characters.</b><br/>

The secret phrase and the salt phrase support a hex notation similar to '\xFF' or '\xff' to represent a character.</br>

It is advisable to use more complex secret phrases and to write secret phrases on paper.<br/>
It is also recommended that salt phrases be memorized in the brain.<br/>


Usage of address:<br/>
　-b　　Brain wallet mode.<br/>
　-brain<br/>
　　　　Brain wallet mode.<br/>
　-o string<br/>
  　　　　Output file name.<br/>
　-output string<br/>
　　　　Output file name.<br/>


Donations are welcome at <code>1Brn37oiWcDoTVqeP1EzbVtCz3dJ7W1Z57</code>


# go钱包
go钱包是用GO语言编写的比特币钱包软件。支持随机钱包和脑钱包。</br>

**脑钱包使用一个秘密短语和一个盐短语生成私钥。**</br>
秘密短语至少16个字符，包含大写字母，小写字母，数字和特殊字符。</br>
盐短语至少6个字符。</br>
秘密短语和盐短语允许使用类似于'\xFF'这样的十六进制表示法表示一个字符</br>

建议使用较为复杂的秘密短语并将秘密短语记在纸上。</br>
同时建议将盐短语记在脑中。

