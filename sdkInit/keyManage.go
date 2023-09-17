package sdkInit

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/crypto/ecies"
)

func Test() {
	// // 获取当前工作目录
	// wd, err := os.Getwd()
	// if err != nil {
	// 	panic(err)
	// }
	keyPath := "/home/go/src/fabric-go-sdk/fixtures/crypto-config/peerOrganizations/org1.example.com/peers/peer1.org1.example.com/msp/keystore/priv_sk"

	keyBytes, err := ioutil.ReadFile(keyPath)
	keyblock, _ := pem.Decode(keyBytes)
	if keyblock == nil {
		fmt.Errorf("failed to decode PEM block")
	}

	key, err := x509.ParsePKCS8PrivateKey(keyblock.Bytes)
	if err != nil {
		fmt.Println(err)
	}

	ecdsaPrivateKey, _ := key.(*ecdsa.PrivateKey)
	priKey := ecies.ImportECDSA(ecdsaPrivateKey)

	publicKey := &ecdsaPrivateKey.PublicKey
	eciesPublicKey := ecies.ImportECDSAPublic(publicKey)

	message := []byte("hello, world")

	encrypted, err := ecies.Encrypt(rand.Reader, eciesPublicKey, message, nil, nil)
	if err != nil {
		panic(err)
	}
	// 2. 对密文进行哈希
	hash := sha256.Sum256(encrypted)

	// 3. 对摘要进行签名
	signature, err := ecdsa.SignASN1(rand.Reader, ecdsaPrivateKey, hash[:])
	if err != nil {
		// 处理签名错误
	}
	total := append(signature, encrypted...)
	//验签

	sign := total[:len(signature)]
	var esig struct {
		R, S *big.Int
	}
	if _, err := asn1.Unmarshal(sign, &esig); err != nil {
		fmt.Println(err)
	}
	valid := ecdsa.Verify(publicKey, hash[:], esig.R, esig.S)
	if !valid {
		fmt.Println("valid failed")
	} else {
		fmt.Println("valid success")
	}

	data, err := ECCDecrypt(encrypted, *priKey)
	if err != nil {
		panic(err)
	}
	fmt.Println("解密后", string(data))

}

//生成一个p256ECC密钥，返回密钥的位置
func GenerateKey(UserID string) string {
	// 生成密钥对
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	// 将私钥编码为 DER 格式
	privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		panic(err)
	}
	// 将 DER 格式的私钥编码为 PEM 格式
	privateKeyPem := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	filename := "./kafka_crypto/" + UserID + ".private.pem"
	// 将 PEM 格式的私钥写入文件
	privateKeyFile, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer privateKeyFile.Close()
	if err := pem.Encode(privateKeyFile, privateKeyPem); err != nil {
		panic(err)
	}
	return filename
}

//获取client的prikey
func GetClientPrivateKey(keypath string) *ecdsa.PrivateKey {
	keyBytes, err := ioutil.ReadFile(keypath)
	keyblock, _ := pem.Decode(keyBytes)
	if keyblock == nil {
		fmt.Errorf("failed to decode PEM block")
	}
	key, err := x509.ParseECPrivateKey(keyblock.Bytes)
	if err != nil {
		panic(err)
	}
	return key
}

//获取节点的prikey
func GetNodePrivateKey(keypath string) *ecdsa.PrivateKey {
	keyBytes, err := ioutil.ReadFile(keypath)
	keyblock, _ := pem.Decode(keyBytes)
	if keyblock == nil {
		fmt.Errorf("failed to decode PEM block")
	}
	key, err := x509.ParsePKCS8PrivateKey(keyblock.Bytes)
	if err != nil {
		panic(err)
	}
	ecdsaPrivateKey, _ := key.(*ecdsa.PrivateKey)

	return ecdsaPrivateKey
}

//用户发给节点，用自己私钥签名，对方公钥加密
func ClientEncryptionByPubECC(client_keypath string, node_pubkey *ecdsa.PublicKey, message []byte) ([]byte, int) {
	//get client key
	key := GetClientPrivateKey(client_keypath)
	//encrypt
	eciesPublicKey := ecies.ImportECDSAPublic(node_pubkey)
	encrypted, err := ecies.Encrypt(rand.Reader, eciesPublicKey, message, nil, nil)
	if err != nil {
		panic(err)
	}
	//hash and sign
	hash := sha256.Sum256(encrypted)
	signature, err := ecdsa.SignASN1(rand.Reader, key, hash[:])
	if err != nil {
		panic(err)
	}

	return append(signature, encrypted...), len(signature)
}

//用户公钥加密，节点私钥签名
func NodeEncryptionByPubECC(node_keypath string, client_pubkey *ecdsa.PublicKey, message []byte) ([]byte, int) {
	key := GetNodePrivateKey(node_keypath)
	//encrypt
	params := ecies.ECIES_AES128_SHA256
	eciesPublicKey := ecies.ImportECDSAPublic(client_pubkey)
	eciesPublicKey.Params = params
	encrypted, err := ecies.Encrypt(rand.Reader, eciesPublicKey, message, nil, nil)
	if err != nil {
		panic(err)
	}
	//hash and sign
	hash := sha256.Sum256(encrypted)
	signature, err := ecdsa.SignASN1(rand.Reader, key, hash[:])
	if err != nil {
		panic(err)
	}
	return append(signature, encrypted...), len(signature)
}

func DecryptionByPri(keypath string, validkey *ecdsa.PublicKey, message []byte, len_of_sign int) []byte {
	//varify
	signature := message[:len_of_sign]
	ciphertext := message[len_of_sign:]
	hash := sha256.Sum256(ciphertext)
	var esig struct {
		R, S *big.Int
	}
	if _, err := asn1.Unmarshal(signature, &esig); err != nil {
		fmt.Println(err)
	}
	valid := ecdsa.Verify(validkey, hash[:], esig.R, esig.S)
	if !valid {
		panic("valid failed")
	} else {
		ecdsaPrivateKey := GetNodePrivateKey(keypath)
		priKey := ecies.ImportECDSA(ecdsaPrivateKey)
		data, err := ECCDecrypt(ciphertext, *priKey)
		if err != nil {
			panic(err)
		}
		return data
	}
}

func PrivateKeyToECICS(privateKey []byte) (*ecies.PrivateKey, error) {
	// 解析DER编码的私钥数据
	key, err := x509.ParseECPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	priKey := ecies.ImportECDSA(key)
	return priKey, nil
}

func ECCDecrypt(ct []byte, prk ecies.PrivateKey) ([]byte, error) {
	pt, err := prk.Decrypt(ct, nil, nil)
	return pt, err
}

//UnmarshalECCPublicKey extract ECC public key from marshaled objects
// func UnmarshalECCPublicKey(object []byte) (usinfo UserInfo) {
// 	var public ecdsa.PublicKey
// 	var pub ecdsa.PublicKey
// 	var pubparams elliptic.CurveParams
// 	rt := new(UserInfo)

// 	errmarsh := json.Unmarshal(object, &rt)
// 	if errmarsh != nil {
// 		fmt.Println("err at UnmarshalECCPublicKey()")
// 		panic(errmarsh)
// 	}
// 	fmt.Println("ahahah", rt.PublicKey.MyCurve.B)

// 	pubparams.B = rt.PublicKey.MyCurve.B
// 	pubparams.BitSize = rt.PublicKey.MyCurve.BitSize
// 	pubparams.Gx = rt.PublicKey.MyCurve.Gx
// 	pubparams.Gy = rt.PublicKey.MyCurve.Gy
// 	pubparams.P = rt.PublicKey.MyCurve.P
// 	pubparams.N = rt.PublicKey.MyCurve.N
// 	pubparams.Name = rt.PublicKey.MyCurve.Name
// 	public.Curve = &pubparams
// 	public.X = rt.PublicKey.X
// 	public.Y = rt.PublicKey.Y
// 	mapstructure.Decode(public, &pub)
// 	usinfo.PublicKey = &pub
// 	usinfo.Attribute = rt.Attribute
// 	usinfo.UserId = rt.UserId
// 	usinfo.Username = rt.Username

// 	return
// }
