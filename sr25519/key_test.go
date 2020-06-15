package sr25519

import (
	"encoding/hex"
	"fmt"
	sr25519 "github.com/ChainSafe/go-schnorrkel"
	"github.com/JFJun/substrate-go/config"
	"github.com/JFJun/substrate-go/ss58"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/blake2b"
	"testing"
)

func TestCreateAddress(t *testing.T) {
	secret:="718b88b8a4df1334265cc726d8cdaedde106fe7eaddea82a1bb457e0011c38d4"
	priv,err:=hex.DecodeString(secret)
	if err != nil {
		panic(err)
	}
	//fmt.Println(len(priv))
	if len(priv)!=32 {
		return
	}
	var s [32]byte
	copy(s[:],priv)
	fmt.Println(s)
	key,err:=sr25519.NewMiniSecretKeyFromRaw(s)
	if err != nil {
		panic(err)
	}
	p:=key.ExpandEd25519().Encode()
	fmt.Println("=============")
	fmt.Println(hex.EncodeToString(p[:]))
	_,pubK:=key.ExpandEd25519(),key.Public()
	pp:=pubK.Encode()
	fmt.Println(hex.EncodeToString(pp[:]))
	pub:=pubK.Encode()
	fmt.Println(ss58.Encode(pub[:],config.SubstratePrefix))

}
var(
	ssPrefix = []byte{0x53, 0x53, 0x35, 0x38, 0x50, 0x52, 0x45}
)
func TestCreateAddress2(t *testing.T) {
	address:="Dj1nP5Ebtjo5GZhksnh1zyDWgAcmiuqmXC6PEe6YzD7JNT6"
	data:=base58.Decode(address)
	fmt.Println(data)
	fmt.Println(data[:33])
	fmt.Println(hex.EncodeToString(data[1:33]))
	var d []byte
	d = append(d,ssPrefix...)
	d = append(d,data[:33]...)
	s:=blake2b.Sum512(d)
	fmt.Println(s)
}

func TestCreateAddress3(t *testing.T) {
	priv,pub,err:=GenerateKey()
	fmt.Println(len(priv),len(pub))
	if err != nil {
		fmt.Println(11)
		panic(err)
	}
	private,err:=PrivateKeyToHex(priv)
	if err != nil {
		panic(err)
	}
	address,err:=CreateAddress(pub,config.SubstratePrefix)
	if err != nil {
		panic(err)
	}
	fmt.Println(private)
	fmt.Println(address)
	/*
	0x3cc8f0463493f30438c7735f9208351f35b43140c7d1544fe8be9b8591ab53aa
	F2WdpYJw37tjBWLAgUZQRN1hxhwUF285L1L3NaqhWKw3tUU
	*/
}

func TestSign(t *testing.T) {
	/*

	0x4719bf9d6a2cfb930c29e048e25a269ab7619087c51c4e09f1c2b6d6885e100ae1a69e063de542665247c137dc1de785b8505cc3c97e40ec6d87710edf7d6551
	0x62959952742388d1a4a20d1a39814833acd2113b5865cc7cf99ff3f48ce041795b17d120aff74d8d0d6e4d5fa1a4f6bac475a886ffea6a4bbd1366e6b329c18a
	*/
	secret:="41c3d017d032c499e59d1ba4ba14e4405a7f95f83fa850a4d323ab52b7d8bd35"
	priv,err:=hex.DecodeString(secret)
	if err != nil {
		panic(err)
	}
	//fmt.Println(len(priv))
	if len(priv)!=32 {
		return
	}
	var s [32]byte
	copy(s[:],priv)
	fmt.Println(s)
	key,err:=sr25519.NewMiniSecretKeyFromRaw(s)
	if err != nil {
		panic(err)
	}

	pub,_:=key.ExpandEd25519().Public()
	p:=pub.Encode()
	fmt.Println(hex.EncodeToString(p[:]))
	//priv,err:=hex.DecodeString(secret)
	//if err != nil {
	//	panic(err)
	//}
	////fmt.Println(len(priv))
	//if len(priv)!=32 {
	//	return
	//}
	//
	//message:=[]byte("Test123")
	//sign,err:=Sign(priv,message)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(hex.EncodeToString(sign))

}
func TestGenerateKey(t *testing.T) {

}

func TestPrivateKeyToHex(t *testing.T) {
/*
	0xe547526543c32f8822448ff3d3232a04e5452889f550193def1657f3019fc30f21418da1dad997e136f6b3240ac261d5a71c7d952aeb5c87d2c2279dc5b90801
	0xf67567cbd08a43da678c9534d607034a5b5fc8c7b89104b8e7da448a3294375e542ce2e595bfa2cce3a4c180d3a85aaab297dadd36b60318b7ffe9f34b660d8b
	*/
	privateKey,_:=hex.DecodeString("e547526543c32f8822448ff3d3232a04e5452889f550193def1657f3019fc30f21418da1dad997e136f6b3240ac261d5a71c7d952aeb5c87d2c2279dc5b90801")
	sig,_:=hex.DecodeString("f67567cbd08a43da678c9534d607034a5b5fc8c7b89104b8e7da448a3294375e542ce2e595bfa2cce3a4c180d3a85aaab297dadd36b60318b7ffe9f34b660d8b")
	var key,nonce [32]byte
	copy(key[:],privateKey[:32])
	copy(nonce[:],privateKey[32:])
	fmt.Println(len(privateKey))
	fmt.Println(key)
	fmt.Println(nonce)
	secret:=sr25519.NewSecretKey(key,nonce)
	var sigBytes [64]byte
	copy(sigBytes[:],sig)
	fmt.Println(len(sig))
	signs:=sr25519.Signature{}
	signs.Decode(sigBytes)
	pub,err:=secret.Public()
	if err != nil {
		fmt.Println(err)
		return
	}
	isOK:=pub.Verify(&signs,sr25519.NewSigningContext([]byte("substrate"),[]byte("Test123")))
	fmt.Println(isOK)
}