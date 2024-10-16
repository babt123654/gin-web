package stacks

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"gin-web/api/v1/wallet"
	"github.com/stretchr/testify/assert"
	"math/big"
	"reflect"
	"testing"
)

func TestSerializeBoolCV(t *testing.T) {
	test1 := &wallet.BooleanCV{wallet.BoolTrue}
	test2 := &wallet.BooleanCV{wallet.BoolFalse}
	assert.Equal(t, "03", hex.EncodeToString(wallet.serializeCV(test1)))
	assert.Equal(t, "04", hex.EncodeToString(wallet.serializeCV(test2)))
}

func TestSerializeResponseCV(t *testing.T) {
	test1 := &wallet.ResponseCV{wallet.ResponseOk, wallet.NewUintCV(big.NewInt(10))}
	test2 := &wallet.ResponseCV{wallet.ResponseErr, wallet.NewUintCV(big.NewInt(10))}
	assert.Equal(t, "07010000000000000000000000000000000a", hex.EncodeToString(wallet.serializeCV(test1)))
	assert.Equal(t, "08010000000000000000000000000000000a", hex.EncodeToString(wallet.serializeCV(test2)))
}

func TestSerializeListCV(t *testing.T) {
	u1 := wallet.NewUintCV(big.NewInt(1))
	u2 := &wallet.SomeCV{wallet.OptionalSome, wallet.NewUintCV(big.NewInt(2))}
	test := &wallet.ListCV{wallet.List, []wallet.ClarityValue{&wallet.BooleanCV{wallet.BoolTrue}, &wallet.BooleanCV{wallet.BoolFalse}, u1, u2}}
	s := wallet.serializeCV(test)
	assert.Equal(t, "0b00000004030401000000000000000000000000000000010a0100000000000000000000000000000002", hex.EncodeToString(s))
}

func TestSerializeTupleCV(t *testing.T) {
	test := make(map[string]wallet.ClarityValue)
	test["ac"] = &wallet.BooleanCV{wallet.BoolTrue}
	test["ab"] = &wallet.BooleanCV{wallet.BoolFalse}
	test["cd"] = wallet.NewUintCV(big.NewInt(1))
	test["ba"] = &wallet.SomeCV{wallet.OptionalSome, wallet.NewUintCV(big.NewInt(2))}

	tuple := &wallet.TupleCV{wallet.Tuple, test}
	s := wallet.serializeCV(tuple)
	assert.Equal(t, "0c0000000402616204026163030262610a01000000000000000000000000000000020263640100000000000000000000000000000001", hex.EncodeToString(s))
}

func TestSerializeStringCV(t *testing.T) {
	s1 := &wallet.StringCV{wallet.IntASCII, "test"}
	s2 := &wallet.StringCV{wallet.IntUTF8, "test"}
	str1 := wallet.serializeCV(s1)
	str2 := wallet.serializeCV(s2)
	assert.Equal(t, "0d0000000474657374", hex.EncodeToString(str1))
	assert.Equal(t, "0e0000000474657374", hex.EncodeToString(str2))
}

func TestDeserializeCV1(t *testing.T) {
	test := wallet.NewUintCV(big.NewInt(100))
	serialized := hex.EncodeToString(wallet.serializeCV(test))
	res := wallet.DeserializeCV(serialized).(wallet.UintCV)
	assert.Equal(t, true, reflect.DeepEqual(*test, res))
}

func TestDeserializeCV2(t *testing.T) {
	test := &wallet.BufferCV{wallet.Buffer, []byte("memo")}
	fmt.Printf("pre:   %+v\n", test)
	serialized := hex.EncodeToString(wallet.serializeCV(test))
	fmt.Printf("serialized: %+v\n", serialized)
	res := wallet.DeserializeCV(serialized).(wallet.BufferCV)
	fmt.Printf("after: %+v\n", res)

	assert.Equal(t, true, reflect.DeepEqual(*test, res))
}

func TestDeserializeCV5(t *testing.T) {
	test := *wallet.NewStandardPrincipalCV("ST2MCYPWTFMD2MGR5YY695EJG0G1R4J2BTJPRGM7H")
	js, _ := json.Marshal(test)
	fmt.Printf("pre:   %s\n", string(js))
	serialized := hex.EncodeToString(wallet.serializeCV(test))
	fmt.Printf("serialized: %+v\n", serialized)
	res := wallet.DeserializeCV(serialized).(wallet.StandardPrincipalCV)
	jsAfter, _ := json.Marshal(res)
	fmt.Printf("after: %+v\n", string(jsAfter))
	assert.Equal(t, true, reflect.DeepEqual(test, res))
}

func TestDeserializeCV6(t *testing.T) {
	test, _ := wallet.NewContractPrincipalCV("SP001SFSMC2ZY76PD4M68P3WGX154XCH7NE3TYMX.pox-pools-1-cycle")
	serialized := hex.EncodeToString(wallet.serializeCV(test))
	res := wallet.DeserializeCV(serialized).(wallet.ContractPrincipalCV)
	assert.Equal(t, true, reflect.DeepEqual(*test, res))
}

func TestDeserializeCV9(t *testing.T) {
	test := wallet.NoneCV{wallet.OptionalNone}
	serialized := hex.EncodeToString(wallet.serializeCV(test))
	res := wallet.DeserializeCV(serialized).(wallet.NoneCV)
	assert.Equal(t, true, reflect.DeepEqual(test, res))
}

func TestDeserializeCV10(t *testing.T) {
	test := &wallet.SomeCV{wallet.OptionalSome, wallet.NewUintCV(big.NewInt(1))}
	serialized := hex.EncodeToString(wallet.serializeCV(test))
	res := wallet.DeserializeCV(serialized).(wallet.SomeCV)
	reflect.DeepEqual(*test, res)
}

func TestDeserializeCVWithJson(t *testing.T) {
	jsonData := `
{
  "functionArgs": [
    {
      "type": 0,
      "value": 100000000000
    },
    {
      "type": 1,
      "value": 100000000000
    },
    {
      "buffer": "dGVzdA==",
      "type": 2
    },
    {
      "type": 3
    },
    {
      "type": 4
    },
    {
      "type": 5,
      "address": {
        "hash160": "a8cf5b9a7d1a2a4305f78c92ba50040382484bd4",
        "type": 0,
        "version": 26
      }
    },
    {
      "type": 6,
      "address": {
        "hash160": "0000e5f9a305ff1cd6692864587c87425275913d",
        "type": 0,
        "version": 22
      },
      "contractName": {
        "content": "pox-pools-1-cycle",
        "lengthPrefixBytes": 1,
        "maxLengthBytes": 128,
        "type": 2
      }
    },
    {
      "type": 7,
      "value": {
        "type": 14,
        "data": "ok ok ok"
      }
    },
    {
      "type": 8,
      "value": {
        "type": 14,
        "data": "error error error"
      }
    },
    {
      "type": 9
    },
    {
      "type": 10,
      "value": {
        "type": 10,
        "value": {
          "type": 2,
          "buffer": "dGVzdA=="
        }
      }
    },
    {
      "type": 11,
      "list": [
        {
          "type": 3
        },
        {
          "type": 4
        }
      ]
    },
    {
      "type": 12,
      "data": {
        "hashbytes": {
          "buffer": "Bc9SpEvz5oKbT4wiHMZ1NVv4O30=",
          "type": 2
        },
        "version": {
          "buffer": "AA==",
          "type": 2
        }
      }
    },
    {
      "data": "testStringAsciiCV",
      "type": 13
    },
    {
      "data": "testStringUtf8CV",
      "type": 14
    }
  ]
}
`
	args := wallet.getFunctionArgs(jsonData)

	res := wallet.DeserializeJson(args)
	test0 := wallet.NewIntCV(big.NewInt(100000000000))
	test1 := wallet.NewUintCV(big.NewInt(100000000000))
	test2 := &wallet.BufferCV{wallet.Buffer, []byte("test")}
	test3 := &wallet.BooleanCV{wallet.BoolTrue}
	test4 := &wallet.BooleanCV{wallet.BoolFalse}
	test5 := *wallet.NewStandardPrincipalCV("ST2MCYPWTFMD2MGR5YY695EJG0G1R4J2BTJPRGM7H")
	test6, _ := wallet.NewContractPrincipalCV("SP001SFSMC2ZY76PD4M68P3WGX154XCH7NE3TYMX.pox-pools-1-cycle")
	test7 := &wallet.ResponseCV{wallet.ResponseOk, &wallet.StringCV{wallet.IntUTF8, "ok ok ok"}}
	test8 := &wallet.ResponseCV{wallet.ResponseErr, &wallet.StringCV{wallet.IntUTF8, "error error error"}}
	test9 := &wallet.NoneCV{wallet.OptionalNone}
	test10 := &wallet.SomeCV{wallet.OptionalSome, &wallet.SomeCV{wallet.OptionalSome, &wallet.BufferCV{wallet.Buffer, []byte("test")}}}
	test11 := &wallet.ListCV{wallet.List, []wallet.ClarityValue{test3, test4}}
	test12, _ := wallet.GetPoxAddress("1Xik14zRm29UsyS6DjhYg4iZeZqsDa8D3")
	test13 := &wallet.StringCV{wallet.IntASCII, "testStringAsciiCV"}
	test14 := &wallet.StringCV{wallet.IntUTF8, "testStringUtf8CV"}
	for i, j := range res {
		wallet.serializeCV(j)
		if i == 0 {
			assert.Equal(t, true, reflect.DeepEqual(j, test0))
		} else if i == 1 {
			assert.Equal(t, true, reflect.DeepEqual(j, test1))
		} else if i == 2 {
			assert.Equal(t, true, reflect.DeepEqual(j, test2))
		} else if i == 3 {
			assert.Equal(t, true, reflect.DeepEqual(j, test3))
		} else if i == 4 {
			assert.Equal(t, true, reflect.DeepEqual(j, test4))
		} else if i == 5 {
			assert.Equal(t, true, reflect.DeepEqual(j, test5))
		} else if i == 6 {
			assert.Equal(t, true, reflect.DeepEqual(j, test6))
		} else if i == 7 {
			assert.Equal(t, true, reflect.DeepEqual(j, test7))
		} else if i == 8 {
			assert.Equal(t, true, reflect.DeepEqual(j, test8))
		} else if i == 9 {
			assert.Equal(t, true, reflect.DeepEqual(j, test9))
		} else if i == 10 {
			assert.Equal(t, true, reflect.DeepEqual(j, test10))
		} else if i == 11 {
			assert.Equal(t, true, reflect.DeepEqual(j, test11))
		} else if i == 12 {
			assert.Equal(t, true, reflect.DeepEqual(j, test12))
		} else if i == 13 {
			assert.Equal(t, true, reflect.DeepEqual(j, test13))
		} else if i == 14 {
			assert.Equal(t, true, reflect.DeepEqual(j, test14))
		}
	}
}

func TestDeserializePostCondition(t *testing.T) {
	str := []string{"000216c03b5520cf3a0bd270d8e41e5e19a464aef6294c010000000000002710", "010316e685b016b3b6cd9ebf35f38e5ae29392e2acd51d0f616c65782d7661756c742d76312d3116e685b016b3b6cd9ebf35f38e5ae29392e2acd51d176167653030302d676f7665726e616e63652d746f6b656e04616c657803000000000078b854"}
	var res []wallet.PostConditionInterface
	for _, s := range str {
		v := wallet.DeserializePostCondition(s)
		res = append(res, v)
	}
	fmt.Printf("%+v", res)
}
