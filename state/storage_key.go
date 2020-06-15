package state

import (
	"encoding/hex"
	"errors"
	"github.com/JFJun/substrate-go/util"
	"github.com/JFJun/substrate-go/xxhash"
	codes "github.com/freehere107/go-scale-codec"

	"strings"
)

func CreateStorageKey(meta *codes.MetadataDecoder, module, fn string, arg []byte, arg2 []byte)(string,error){
	var (
		data []byte
		err error
	)
	if meta==nil {
		return "", errors.New("metadata is null")
	}
	for _,mod:=range meta.Metadata.Metadata.Modules {
		if mod.Name== module{
			for _,method:=range mod.Storage{
				if method.Name==fn {
					if method.Type.MapType!=nil {
						hashMethod:=method.Type.MapType.Hasher
						data,err = createKey(module,fn,hashMethod,"",arg,nil,IsMap)
					}
					if method.Type.DoubleMapType!=nil {
						hasherMethod:=method.Type.DoubleMapType.Hasher
						keyHasherMethod:=method.Type.DoubleMapType.Key2Hasher
						data,err = createKey(module,fn,hasherMethod,keyHasherMethod,arg,nil,IsDoubleMap)
					}
					if method.Type.PlainType!=nil{
						HashMethod:="Twox64Concat"
						data,err = createKey(module,fn,HashMethod,"",arg,nil,IsPlainType)
					}
				}
			}
		}
	}
	if err != nil {
		return "", err
	}
	if data==nil {
		return "",errors.New("create storage key data is null")
	}
	return "0x"+hex.EncodeToString(data),nil
}

const  (
	IsMap = iota
	IsDoubleMap
	IsPlainType
)

func createKey(module,fn ,hasherMethod,keyHasherMethod string,arg1,arg2 []byte,mapType int)([]byte,error){
	var key []byte
	if mapType==IsMap {
		hasher,err:=util.SelectHash(hasherMethod)
		if err != nil {
			return nil, err
		}
		_,err=hasher.Write(arg1)
		if err != nil {
			return nil,err
		}
		key=hasher.Sum(nil)
		if strings.Contains(hasherMethod,"Concat") {
			key = append(key,arg1...)
		}

	}
	if mapType==IsDoubleMap {
		hasher,err:=util.SelectHash(hasherMethod)
		if err != nil {
			return nil, err
		}
		keyHasher,err:=util.SelectHash(hasherMethod)
		if err != nil {
			return nil, err
		}
		_,err = hasher.Write(arg1)
		if err != nil {
			return nil,err
		}
		data1:=hasher.Sum(nil)
		if strings.Contains(hasherMethod,"Concat") {

			key = append(data1,arg1...)
		}
		_,err = keyHasher.Write(arg2)
		if err != nil {
			return nil,err
		}
		data2:=hasher.Sum(nil)
		if strings.Contains(keyHasherMethod,"Concat") {
			key = append(data2,arg2...)
		}
		key = append(key,data1...)
		key = append(key,data2...)
	}
	if mapType==IsPlainType {
		//hasher,err:=util.SelectHash(hasherMethod)
		//if err != nil {
		//	return nil, err
		//}
		//_,err=hasher.Write(arg1)
		//if err != nil {
		//	return nil,err
		//}
		//key=hasher.Sum(nil)
		//if strings.Contains(hasherMethod,"Concat") {
		//	key = append(key,arg1...)
		//}
	}
	return append(createPrefixedKey(module,fn),key...),nil
}



func createPrefixedKey(module, fn string) []byte {
	return append(xxhash.New128([]byte(module)).Sum(nil), xxhash.New128([]byte(fn)).Sum(nil)...)
}

type AccountInfo struct {
	Nonce    U32			`json:"nonce"`
	Refcount U8				`json:"ref_count"`
	Data     struct {
		Free       U128		`json:"free"`
		Reserved   U128		`json:"reserved"`
		MiscFrozen U128		`json:"misc_frozen"`
		FreeFrozen U128		`json:"free_frozen"`
	}						`json:"data"`
}