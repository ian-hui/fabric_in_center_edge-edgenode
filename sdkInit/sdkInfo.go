package sdkInit

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	contextAPI "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
)

type OrgInfo struct {
	OrgAdminUser          string // like "Admin"
	OrgName               string // like "Org1"
	OrgMspId              string // like "Org1MSP"
	OrgUser               string // like "User1"
	OrgMspClient          *mspclient.Client
	OrgAdminClientContext *contextAPI.ClientProvider
	OrgResMgmt            *resmgmt.Client
	OrgPeerNum            int
	//Peers                 []*fab.Peer
	OrgAnchorFile string // like ./channel-artifacts/Org2MSPanchors.tx
}

type SdkEnvInfo struct {
	// 通道信息
	ChannelID     string // like "simplecc"
	ChannelConfig string // like os.Getenv("GOPATH") + "/src/github.com/hyperledger/fabric-samples/test-network/channel-artifacts/testchannel.tx"

	// 组织信息
	Orgs []*OrgInfo
	// 排序服务节点信息
	OrdererAdminUser     string // like "Admin"
	OrdererOrgName       string // like "OrdererOrg"
	OrdererEndpoint      string
	OrdererClientContext *contextAPI.ClientProvider
	// 链码信息
	ChaincodeID      string
	ChaincodeGoPath  string
	ChaincodePath    string
	ChaincodeVersion string
	ChClient         *channel.Client
	EvClient         *event.Client
}

type Application struct {
	SdkEnvInfo *SdkEnvInfo
}

type UserInfo struct {
	UserId    string
	Username  string
	Attribute []string
	PublicKey *MyPublicKey `json:"public_key"`
}

type MyPublicKey struct {
	*MyCurve
	X, Y *big.Int
}



type retrieve struct {
	CurveParams *elliptic.CurveParams `json:"Curve"`
	MyX         *big.Int              `json:"X"`
	MyY         *big.Int              `json:"Y"`
}

type UserInfoGetter struct {
	UserId    string
	Username  string
	Attribute []string
	PublicKey *retrieve
}

type MyCurve struct {
	P       *big.Int // the order of the underlying field
	N       *big.Int // the order of the base point
	B       *big.Int // the constant of the curve equation
	Gx, Gy  *big.Int // (x,y) of the base point
	BitSize int      // the size of the underlying field
	Name    string   // the canonical name of the curve
}

//自定义mypublic的序列化和反序列化
func (key *MyPublicKey) MarshalJSON() ([]byte, error) {
	// 将 X 和 Y 转换成字符串
	type Alias MyPublicKey
	return json.Marshal(&struct {
		Curve *MyCurve `json:"Curve"`
		X     string   `json:"x"`
		Y     string   `json:"y"`
	}{
		Curve: key.MyCurve,
		X:     key.X.String(),
		Y:     key.Y.String(),
	})
}

func (key *MyPublicKey) UnmarshalJSON(data []byte) error {
	// 将 X 和 Y 解析为字符串，并转换成 big.Int
	type Alias MyPublicKey
	aux := &struct {
		Curve *MyCurve `json:"Curve"`
		X     string   `json:"x"`
		Y     string   `json:"y"`
	}{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	x, ok := new(big.Int).SetString(aux.X, 10)
	if !ok {
		return fmt.Errorf("failed to parse x as big.Int")
	}
	y, ok := new(big.Int).SetString(aux.Y, 10)
	if !ok {
		return fmt.Errorf("failed to parse y as big.Int")
	}
	key.X = x
	key.Y = y
	key.MyCurve = aux.Curve
	return nil
}

//因为ecdsa在标准库中，所以只能是这样的写法
func ConversionEcdsaPub2MyPub(pubkey *ecdsa.PublicKey) *MyPublicKey {
	pubcur := MyCurve{
		B:       pubkey.Curve.Params().B,
		N:       pubkey.Curve.Params().N,
		P:       pubkey.Curve.Params().P,
		Gx:      pubkey.Curve.Params().Gx,
		Gy:      pubkey.Curve.Params().Gy,
		BitSize: pubkey.Curve.Params().BitSize,
		Name:    pubkey.Curve.Params().Name,
	}
	pubk := MyPublicKey{
		MyCurve: &pubcur,
		X:       pubkey.X,
		Y:       pubkey.Y,
	}
	return &pubk
}

func (pubkey *MyPublicKey) ConversionMyPub2EcdsaPub() *ecdsa.PublicKey {
	pubparams := elliptic.CurveParams{
		B:       pubkey.MyCurve.B,
		BitSize: pubkey.MyCurve.BitSize,
		Gx:      pubkey.MyCurve.Gx,
		Gy:      pubkey.MyCurve.Gy,
		P:       pubkey.MyCurve.P,
		N:       pubkey.MyCurve.N,
		Name:    pubkey.MyCurve.Name,
	}
	ecdsaPub := ecdsa.PublicKey{
		Curve: &pubparams,
		X:     pubkey.X,
		Y:     pubkey.Y,
	}
	return &ecdsaPub
}

//自定义mycurve的序列化和反序列化
func (curve *MyCurve) MarshalJSON() ([]byte, error) {
	type Alias MyCurve
	return json.Marshal(&struct {
		P       string `json:"P"`
		N       string `json:"N"`
		B       string `json:"B"`
		Gx      string `json:"Gx"`
		Gy      string `json:"Gy"`
		BitSize int    `json:"BitSize"`
		Name    string `json:"Name"`
	}{
		P:       curve.P.String(),
		N:       curve.N.String(),
		B:       curve.B.String(),
		Gx:      curve.Gx.String(),
		Gy:      curve.Gy.String(),
		BitSize: curve.BitSize,
		Name:    curve.Name,
	})
}

func (curve *MyCurve) UnmarshalJSON(data []byte) error {
	type Alias MyCurve
	aux := &struct {
		P       string `json:"P"`
		N       string `json:"N"`
		B       string `json:"B"`
		Gx      string `json:"Gx"`
		Gy      string `json:"Gy"`
		BitSize int    // the size of the underlying field
		Name    string // the canonical name of the curve
	}{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	P, ok := new(big.Int).SetString(aux.P, 10)
	if !ok {
		return fmt.Errorf("failed to parse P as big.Int")
	}
	N, ok := new(big.Int).SetString(aux.N, 10)
	if !ok {
		return fmt.Errorf("failed to parse N as big.Int")
	}
	B, ok := new(big.Int).SetString(aux.B, 10)
	if !ok {
		return fmt.Errorf("failed to parse B as big.Int")
	}
	Gx, ok := new(big.Int).SetString(aux.Gx, 10)
	if !ok {
		return fmt.Errorf("failed to parse Gx as big.Int")
	}
	Gy, ok := new(big.Int).SetString(aux.Gy, 10)
	if !ok {
		return fmt.Errorf("failed to parse Gy as big.Int")
	}
	curve.P = P
	curve.N = N
	curve.B = B
	curve.Gx = Gx
	curve.Gy = Gy
	curve.BitSize = aux.BitSize
	curve.Name = aux.Name
	return nil

}
