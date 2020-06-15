package sr25519

import (
	"encoding/hex"
	"errors"
	"fmt"
	sr25519 "github.com/ChainSafe/go-schnorrkel"
	"github.com/JFJun/substrate-go/ss58"
	r255 "github.com/gtank/ristretto255"
)


type KeyPair struct {
	Wif string
	Address string
}

func GenerateKey()([]byte,[]byte, error){
	secret,err:=sr25519.GenerateMiniSecretKey()
	if err != nil {
		return nil, nil, err
	}
	if len(secret.Encode())!=32  {
		return nil, nil, errors.New("private key or public key length i not equal 32")
	}
	priv:=secret.Encode()
	pub:=secret.Public().Encode()
	return priv[:],pub[:],nil
}

func CreateAddress(pubKey,prefix []byte)(string,error){
	return ss58.Encode(pubKey,prefix)
}

func PrivateKeyToAddress(privateKey,prefix []byte)(string,error){
	var p [32]byte
	copy(p[:],privateKey[:])
	secret,err:=sr25519.NewMiniSecretKeyFromRaw(p)
	if err != nil {
		panic(err)
	}

	pub:=secret.Public().Encode()
	return ss58.Encode(pub[:],prefix)
}

func PrivateKeyToHex(privateKey []byte)(string,error){
	if len(privateKey)!=32 {
		return "",errors.New("private key length is not equal 32")
	}
	privHex:=hex.EncodeToString(privateKey)
	return "0x"+privHex,nil
}
// todo
func PrivateKeyToWif(privateKey []byte)(string,error){
	if len(privateKey)!=32 {
		return "",errors.New("private key length is not equal 32")
	}
	return "",nil
}
func Sign(privateKey,message []byte)([]byte,error){
	var sigBytes []byte
	var key ,nonce [32]byte
	copy(key[:],privateKey[:32])
	signContext:=sr25519.NewSigningContext([]byte("substrate"),message)
	if	len(privateKey)==32{	// Is seed

		sk,err:=sr25519.NewMiniSecretKeyFromRaw(key)
		if err != nil {
			return nil, err
		}

		signContext.AppendMessage([]byte("proto-name"), []byte("Schnorr-sig"))
		pub := sk.Public()
		pubc := pub.Compress()
		signContext.AppendMessage([]byte("sign:pk"), pubc[:])

		r, err := sr25519.NewRandomScalar()
		if err != nil {

			return nil, err
		}
		R := r255.NewElement().ScalarBaseMult(r)
		signContext.AppendMessage([]byte("sign:R"), R.Encode([]byte{}))

		// form k
		kb := signContext.ExtractBytes([]byte("sign:c"), 64)
		k := r255.NewScalar()
		k.FromUniformBytes(kb)

		// form scalar from secret key x
		x, err := sr25519.ScalarFromBytes(sk.ExpandEd25519().Encode())
		if err != nil {
			return nil, err
		}
		// s = kx + r
		s := x.Multiply(x, k).Add(x, r)
		sig:=sr25519.Signature{R: R, S: s}
		sbs:=sig.Encode()
		sigBytes = sbs[:]
		varifySigContent:=sr25519.NewSigningContext([]byte("substrate"),message)
		if !sk.Public().Verify(&sig,varifySigContent){
			return nil,errors.New("verify sign error")
		}
	}else if len(privateKey)==64 {		//Is private key
		copy(nonce[:],privateKey[32:])
		sk:=sr25519.NewSecretKey(key,nonce)
		sig,err:=sk.Sign(signContext)
		if err != nil {
			return nil, fmt.Errorf("sr25519 sign error,err=%v",err)
		}
		sbs:=sig.Encode()
		sigBytes = sbs[:]
		pub,_:=sk.Public()
		if !pub.Verify(sig,sr25519.NewSigningContext([]byte("substrate"),message)){
			return nil,errors.New("verify sign error")
		}
	}
	return sigBytes[:],nil

}
