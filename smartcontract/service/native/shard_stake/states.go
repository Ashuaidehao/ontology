package shard_stake

import (
	"encoding/hex"
	"fmt"
	"github.com/ontio/ontology-crypto/keypair"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/common/serialization"
	"github.com/ontio/ontology/smartcontract/service/native/utils"
	"io"
	"sort"
)

type View uint64 // shard consensus epoch index

type PeerViewInfo struct {
	PeerPubKey          string
	Owner               common.Address
	WholeFee            uint64 // each epoch handling fee
	FeeBalance          uint64 // each epoch handling fee not be withdrawn
	WholeStakeAmount    uint64 // node + user stake amount
	WholeUnfreezeAmount uint64 // all user can withdraw amount
	UserStakeAmount     uint64 // user stake amount
	MaxAuthorization    uint64 // max user stake amount
	Proportion          uint64 // proportion to user
}

func (this *PeerViewInfo) Serialize(w io.Writer) error {
	if err := serialization.WriteString(w, this.PeerPubKey); err != nil {
		return fmt.Errorf("serialize peer public key failed, err: %s", err)
	}
	if err := utils.WriteAddress(w, this.Owner); err != nil {
		return fmt.Errorf("serialize owner failed, err: %s", err)
	}
	if err := utils.WriteVarUint(w, this.WholeFee); err != nil {
		return fmt.Errorf("serialize whole fee failed, err: %s", err)
	}
	if err := utils.WriteVarUint(w, this.FeeBalance); err != nil {
		return fmt.Errorf("serialize fee balance failed, err: %s", err)
	}
	if err := utils.WriteVarUint(w, this.WholeStakeAmount); err != nil {
		return fmt.Errorf("serialize whole stake amount failed, err: %s", err)
	}
	if err := utils.WriteVarUint(w, this.WholeUnfreezeAmount); err != nil {
		return fmt.Errorf("serialize whole unfreeze amount failed, err: %s", err)
	}
	if err := utils.WriteVarUint(w, this.UserStakeAmount); err != nil {
		return fmt.Errorf("serialize user stake amount failed, err: %s", err)
	}
	if err := utils.WriteVarUint(w, this.MaxAuthorization); err != nil {
		return fmt.Errorf("serialize max authorization failed, err: %s", err)
	}
	if err := utils.WriteVarUint(w, this.Proportion); err != nil {
		return fmt.Errorf("serialize propotion failed, err: %s", err)
	}
	return nil
}
func (this *PeerViewInfo) Deserialize(r io.Reader) error {
	var err error = nil
	if this.PeerPubKey, err = serialization.ReadString(r); err != nil {
		return fmt.Errorf("deserialize: read peer pub key failed, err: %s", err)
	}
	if this.Owner, err = utils.ReadAddress(r); err != nil {
		return fmt.Errorf("deserialize: read owner failed, err: %s", err)
	}
	if this.WholeFee, err = utils.ReadVarUint(r); err != nil {
		return fmt.Errorf("deserialize: read whole fee failed, err: %s", err)
	}
	if this.FeeBalance, err = utils.ReadVarUint(r); err != nil {
		return fmt.Errorf("deserialize: read fee balance failed, err: %s", err)
	}
	if this.WholeStakeAmount, err = utils.ReadVarUint(r); err != nil {
		return fmt.Errorf("deserialize: read whole stake amount failed, err: %s", err)
	}
	if this.WholeUnfreezeAmount, err = utils.ReadVarUint(r); err != nil {
		return fmt.Errorf("deserialize: read whole unfreeze amount failed, err: %s", err)
	}
	if this.UserStakeAmount, err = utils.ReadVarUint(r); err != nil {
		return fmt.Errorf("deserialize: read user stake amount failed, err: %s", err)
	}
	if this.MaxAuthorization, err = utils.ReadVarUint(r); err != nil {
		return fmt.Errorf("deserialize: read max authorization failed, err: %s", err)
	}
	if this.Proportion, err = utils.ReadVarUint(r); err != nil {
		return fmt.Errorf("deserialize: read proportion failed, err: %s", err)
	}
	return nil
}

type ViewInfo struct {
	Peers map[keypair.PublicKey]*PeerViewInfo
}

func (this *ViewInfo) GetPeer(pubKey string) (*PeerViewInfo, keypair.PublicKey, error) {
	if this.Peers == nil {
		return nil, nil, fmt.Errorf("GetPeer: peers is nil")
	}
	pubKeyData, err := hex.DecodeString(pubKey)
	if err != nil {
		return nil, nil, fmt.Errorf("GetPeer: decode param pub key failed, err: %s", err)
	}
	paramPubkey, err := keypair.DeserializePublicKey(pubKeyData)
	if err != nil {
		return nil, nil, fmt.Errorf("GetPeer: deserialize param pub key failed, err: %s", err)
	}
	shardPeerStakeInfo, ok := this.Peers[paramPubkey]
	if !ok {
		return nil, nil, fmt.Errorf("GetPeer: peer %s not exist", pubKey)
	}
	return shardPeerStakeInfo, paramPubkey, nil
}

func (this *ViewInfo) Serialize(w io.Writer) error {
	err := utils.WriteVarUint(w, uint64(len(this.Peers)))
	if err != nil {
		return fmt.Errorf("serialize: wirte peers len faield, err: %s", err)
	}
	peerInfoList := make([]*PeerViewInfo, 0)
	for _, info := range this.Peers {
		peerInfoList = append(peerInfoList, info)
	}
	sort.SliceStable(peerInfoList, func(i, j int) bool {
		return peerInfoList[i].PeerPubKey > peerInfoList[j].PeerPubKey
	})
	for index, info := range peerInfoList {
		err = info.Serialize(w)
		if err != nil {
			return fmt.Errorf("serialize: index %d, err: %s", index, err)
		}
	}
	return nil
}

func (this *ViewInfo) Deserialize(r io.Reader) error {
	num, err := utils.ReadVarUint(r)
	if err != nil {
		return fmt.Errorf("deserialze: read peers num failed, err: %s", err)
	}
	this.Peers = make(map[keypair.PublicKey]*PeerViewInfo)
	for i := uint64(0); i < num; i++ {
		info := &PeerViewInfo{}
		err = info.Deserialize(r)
		if err != nil {
			return fmt.Errorf("deserialize: index %d, err: %s", i, err)
		}
		pubKeyData, err := hex.DecodeString(info.PeerPubKey)
		if err != nil {
			return fmt.Errorf("deserialze: decode pub key failed, err: %s", err)
		}
		pubKey, err := keypair.DeserializePublicKey(pubKeyData)
		if err != nil {
			return fmt.Errorf("deserialze: deserialize pub key failed, err: %s", err)
		}
		this.Peers[pubKey] = info
	}
	return nil
}

type UserPeerStakeInfo struct {
	PeerPubKey     string
	StakeAmount    uint64
	UnfreezeAmount uint64
}

type UserStakeInfo struct {
	Peers map[keypair.PublicKey]*UserPeerStakeInfo
}

func (this *UserStakeInfo) Serialize(w io.Writer) error {
	err := utils.WriteVarUint(w, uint64(len(this.Peers)))
	if err != nil {
		return fmt.Errorf("serialize: wirte peers len faield, err: %s", err)
	}
	userPeerInfoList := make([]*UserPeerStakeInfo, 0)
	for _, info := range this.Peers {
		userPeerInfoList = append(userPeerInfoList, info)
	}
	sort.SliceStable(userPeerInfoList, func(i, j int) bool {
		return userPeerInfoList[i].PeerPubKey > userPeerInfoList[j].PeerPubKey
	})
	for index, info := range userPeerInfoList {
		err = serialization.WriteString(w, info.PeerPubKey)
		if err != nil {
			return fmt.Errorf("serialize peer public key failed, index %d, err: %s", index, err)
		}
		err = utils.WriteVarUint(w, info.StakeAmount)
		if err != nil {
			return fmt.Errorf("serialize stake amount failed, index %d, err: %s", index, err)
		}
		err = utils.WriteVarUint(w, info.UnfreezeAmount)
		if err != nil {
			return fmt.Errorf("serialize unfreeze amount failed, index %d, err: %s", index, err)
		}
	}
	return nil
}

func (this *UserStakeInfo) Deserialize(r io.Reader) error {
	num, err := utils.ReadVarUint(r)
	if err != nil {
		return fmt.Errorf("deserialze: read peers num failed, err: %s", err)
	}
	this.Peers = make(map[keypair.PublicKey]*UserPeerStakeInfo)
	for i := uint64(0); i < num; i++ {
		info := &UserPeerStakeInfo{}
		peerPubKey, err := serialization.ReadString(r)
		if err != nil {
			return fmt.Errorf("deserialze: read peer pub key failed, index %d, err: %s", i, err)
		}
		info.PeerPubKey = peerPubKey
		pubKeyData, err := hex.DecodeString(peerPubKey)
		if err != nil {
			return fmt.Errorf("deserialze: decode param pub key failed, err: %s", err)
		}
		pubKey, err := keypair.DeserializePublicKey(pubKeyData)
		if err != nil {
			return fmt.Errorf("deserialze: deserialize param pub key failed, err: %s", err)
		}
		stakeAmount, err := utils.ReadVarUint(r)
		if err != nil {
			return fmt.Errorf("deserialze: deserialize whole fee failed, err: %s", err)
		}
		info.StakeAmount = stakeAmount
		unfreezeAmount, err := utils.ReadVarUint(r)
		if err != nil {
			return fmt.Errorf("deserialze: deserialize whole fee failed, err: %s", err)
		}
		info.UnfreezeAmount = unfreezeAmount
		this.Peers[pubKey] = info
	}
	return nil
}
