package types

import (
	"github.com/celestiaorg/nmt/namespace"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

type Message struct {
	// NamespaceID defines the namespace of this message, i.e. the
	// namespace it will use in the namespaced Merkle tree.
	//
	// TODO: spec out constrains and
	// introduce dedicated type instead of just []byte
	NamespaceID namespace.ID

	// Data is the actual data contained in the message
	// (e.g. a block of a virtual sidechain).
	Data []byte
}

type Messages struct {
	MessagesList []Message `json:"msgs"`
}

var (
	MessageEmpty  = Message{}
	MessagesEmpty = Messages{}
)

func MessageFromProto(p *tmproto.Message) Message {
	if p == nil {
		return MessageEmpty
	}
	return Message{
		NamespaceID: p.NamespaceId,
		Data:        p.Data,
	}
}

func MessagesFromProto(p *tmproto.Messages) Messages {
	if p == nil {
		return MessagesEmpty
	}

	msgs := make([]Message, 0, len(p.MessagesList))

	for i := 0; i < len(p.MessagesList); i++ {
		msgs = append(msgs, MessageFromProto(p.MessagesList[i]))
	}
	return Messages{MessagesList: msgs}
}

func (msg Message) ToProto() *tmproto.Message {
	return &tmproto.Message{
		NamespaceId: msg.NamespaceID,
		Data:        msg.Data,
	}
}

func (msgs Messages) ToProto() tmproto.Messages {
	pmsgs := make([]*tmproto.Message, len(msgs.MessagesList))
	for i, m := range msgs.MessagesList {
		pmsgs[i] = m.ToProto()
	}
	return tmproto.Messages{MessagesList: pmsgs}
}

type IntermediateStateRoots struct {
	RawRootsList []tmbytes.HexBytes `json:"intermediate_roots"`
}

func IntermediateStateRootsFromProto(isrs tmproto.IntermediateStateRoots) IntermediateStateRoots {
	roots := make([]tmbytes.HexBytes, len(isrs.RawRootsList))
	for i, r := range isrs.RawRootsList {
		roots[i] = tmbytes.HexBytes(r)
	}
	return IntermediateStateRoots{RawRootsList: roots}
}

func (isrs IntermediateStateRoots) ToProto() tmproto.IntermediateStateRoots {
	roots := make([][]byte, len(isrs.RawRootsList))
	for i, r := range isrs.RawRootsList {
		roots[i] = r
	}
	return tmproto.IntermediateStateRoots{RawRootsList: roots}
}
