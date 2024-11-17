

https://pkg.go.dev/crypto









## 哈希函数

哈希函数是一种将任意长度的数据映射为固定长度哈希值的函数。它的主要特点是：

1. 固定输出长度：哈希函数生成的哈希值长度是固定的，无论输入数据的长度如何，输出的哈希值长度保持相同。
2. 不可逆性：哈希函数是单向函数，即从哈希值无法还原出原始数据。无法通过哈希值逆向推导出输入数据。
3. 雪崩效应：输入数据的微小变化会导致生成的哈希值产生巨大的差异。这种特性称为雪崩效应，保证了输入数据的细微变化会导致完全不同的哈希值。
4. 唯一性：理想情况下，每个不同的输入数据应该生成唯一的哈希值。但由于哈希值的固定长度限制，不同的输入数据可能会产生相同的哈希值，这种情况称为哈希碰撞。好的哈希函数应该尽量降低碰撞的概率。

* func RegisterHash(h Hash, f func() hash.Hash)  全局注册哈希函数
* type Hash
  * func (h Hash) Available() bool
  * func (h Hash) HashFunc() Hash
  * func (h Hash) New() hash.Hash
  * func (h Hash) Size() int
  * func (h Hash) String() string

内部支持算法

https://pkg.go.dev/crypto/md5

强哈希https://pkg.go.dev/crypto/sha256

```go
const (
	MD4         Hash = 1 + iota // import golang.org/x/crypto/md4
	MD5                         // import crypto/md5
	SHA1                        // import crypto/sha1
	SHA224                      // import crypto/sha256
	SHA256                      // import crypto/sha256
	SHA384                      // import crypto/sha512
	SHA512                      // import crypto/sha512
	MD5SHA1                     // no implementation; MD5+SHA1 used for TLS RSA
	RIPEMD160                   // import golang.org/x/crypto/ripemd160
	SHA3_224                    // import golang.org/x/crypto/sha3
	SHA3_256                    // import golang.org/x/crypto/sha3
	SHA3_384                    // import golang.org/x/crypto/sha3
	SHA3_512                    // import golang.org/x/crypto/sha3
	SHA512_224                  // import crypto/sha512
	SHA512_256                  // import crypto/sha512
	BLAKE2s_256                 // import golang.org/x/crypto/blake2s
	BLAKE2b_256                 // import golang.org/x/crypto/blake2b
	BLAKE2b_384                 // import golang.org/x/crypto/blake2b
	BLAKE2b_512                 // import golang.org/x/crypto/blake2b

)
```

示例

```go
package main

import (
	"crypto/md5"
	"fmt"
)

func main() {
	h := md5.New()
	h.Write([]byte("hello world"))
	fmt.Println("%x", h.Sum(nil))
}
```





## 非对称加密

https://pkg.go.dev/crypto/rsa

非对称加密是一种密码学技术，使用一对密钥（公钥和私钥）来进行加密和解密操作。与对称加密算法不同，非对称加密算法使用两个相关联的密钥，其中一个用于加密数据，另一个用于解密数据。

在非对称加密中，公钥是公开的，可以向任何人公开，用于加密数据。而私钥是保密的，仅由数据的接收方持有，用于解密数据。

非对称加密算法具有以下特点：

1. 加密：使用公钥对数据进行加密，生成密文。只有持有私钥的接收方可以解密该密文。

2. 解密：使用私钥对密文进行解密，还原为原始数据。私钥是唯一能够解密的密钥。

3. 数字签名：非对称加密算法还可以用于生成和验证数字签名。发送方使用私钥对数据进行签名，接收方使用公钥验证签名的有效性。

4. 密钥交换：非对称加密算法也可以用于安全地交换对称加密算法所需的密钥。发送方使用接收方的公钥加密对称密钥，接收方使用私钥解密获取对称密钥。

常见的非对称加密算法包括 RSA（Rivest-Shamir-Adleman）、DSA（Digital Signature Algorithm）、ECC（Elliptic Curve Cryptography）等。这些算法基于不同的数学原理，提供了不同的安全性和性能特征。

非对称加密算法在安全通信、数字签名、密钥交换等场景中得到广泛应用。同时，由于非对称加密算法的计算复杂性较高，通常用于加密较短的数据或加密对称密钥等关键信息，而不适用于大量数据的加密。因此，非对称加密通常与对称加密结合使用，以实现更高效的加密通信方案。

* type Decrypter 非对称加密接口
* type DecrypterOpts any 
* type PrivateKey any 使用某种算法生成的私钥
* type PublicKey any 使用某种算法生成的公钥
* type Signer interface 用于签名操作
* type SignerOpts interface

```go
type Decrypter interface {
	// Public returns the public key corresponding to the opaque,
	// private key.
	Public() PublicKey

	// Decrypt decrypts msg. The opts argument should be appropriate for
	// the primitive used. See the documentation in each implementation for
	// details.
	Decrypt(rand io.Reader, msg []byte, opts DecrypterOpts) (plaintext []byte, err error)
}
type DecrypterOpts any

type Signer interface {
	Public() PublicKey
	Sign(rand io.Reader, digest []byte, opts SignerOpts) (signature []byte, err error)
}
type SignerOpts interface {
	HashFunc() Hash
}
```

rsa加密示例

* func GenerateKey(random io.Reader, bits int) (*PrivateKey, error)  生成指定长度随机私钥
* 

```go
package main
import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
)

func main() {
	// 生成密钥对
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println("Error generating RSA key:", err)
		return
	}

	// 获取公钥
	publicKey := &privateKey.PublicKey

	// 原始数据
	message := []byte("Hello, world!")

	// 加密
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, message)
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return
	}

	fmt.Println("Encrypted data:", ciphertext)

	// 解密
	plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, ciphertext)
	if err != nil {
		fmt.Println("Error decrypting data:", err)
		return
	}

	fmt.Println("Decrypted data:", string(plaintext))
}

```







## 加密和签名

加密（Encrypt）和签名（Sign）是密码学中常用的两种操作，用于不同的安全目的。

1. 加密（Encrypt）：
加密是将明文（原始数据）通过使用密钥（对称密钥或公钥）进行转换，生成密文（加密后的数据）。加密的目的是保护数据的机密性，确保只有授权的人能够解密并访问原始数据。

- 对称加密：使用相同的密钥进行加密和解密操作。发送方和接收方共享相同的密钥，因此要确保密钥的安全性。
- 非对称加密：使用一对相关联的密钥（公钥和私钥）进行加密和解密操作。公钥用于加密数据，私钥用于解密数据。公钥可以公开，而私钥保密。

2. 签名（Sign）：
签名是用于验证数据完整性和身份认证的技术。它使用私钥对数据进行加密，生成数字签名。接收方可以使用公钥验证签名的有效性，确保数据的完整性和未被篡改。签名的目的是防止数据被篡改并验证数据的来源。

- 数字签名：
   - 私钥用于生成数字签名。
   - 公钥用于验证签名的有效性。

需要注意的是，加密和签名是不同的操作，其使用的密钥和目的也不同。加密主要关注数据的保密性，而签名则关注数据的完整性和认证。

在实际应用中，加密和签名通常结合使用，以确保数据的安全性和完整性。例如，可以使用非对称加密算法对对称密钥进行加密和传输，然后使用对称加密算法对数据进行加密。同时，可以使用私钥对数据进行签名，然后使用公钥验证签名的有效性。这样，数据在传输过程中既能保密又能验证完整性。







https加密过程

1. 客户端在浏览器中输入一个https网址，然后连接到server的443端口 采用https协议的server必须有一套数字证书（一套公钥和密钥） 首先server将证书（公钥）传送到客户端 客户端解析证书，验证成功，则生成一个随机数（私钥），并用证书将该随机数加密后传回server server用密钥解密后，获得这个随机值，然后将要传输的信息和私钥通过某种算法混合在一起（加密）传到客户端 客户端用之前的生成的随机数（私钥）解密服务器端传来的信息

2. 首先浏览器会从内置的证书列表中索引，找到服务器下发证书对应的机构，如果没有找到，此时就会提示用户该证书是不是由权威机构颁发，是不可信任的。如果查到了对应的机构，则取出该机构颁发的公钥。

   用机构的证书公钥解密得到证书的内容和证书签名，内容包括网站的网址、网站的公钥、证书的有效期等。浏览器会先验证证书签名的合法性。签名通过后，浏览器验证证书记录的网址是否和当前网址是一致的，不一致会提示用户。如果网址一致会检查证书有效期，证书过期了也会提示用户。这些都通过认证时，浏览器就可以安全使用证书中的网站公钥了。



















