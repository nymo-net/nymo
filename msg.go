package nymo

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"math/big"
	"time"

	"github.com/nymo-net/nymo/pb"
	"google.golang.org/protobuf/proto"
)

type message struct {
	Sender   *address
	SendTime time.Time
	Content  string
}

func (u *user) DecryptMessage(msg *pb.Message) *message {
	if msg.TargetCohort != u.cohort {
		return nil
	}

	eKeyX, eKeyY := elliptic.UnmarshalCompressed(curve, msg.EphemeralPub)
	if eKeyX == nil {
		return nil
	}

	secret, iv := curve.ScalarMult(eKeyX, eKeyY, u.key.D.Bytes())
	cp, err := aes.NewCipher(secret.Bytes())
	if err != nil {
		return nil
	}

	encMsg := msg.EncMessage
	cipher.NewCBCDecrypter(cp, iv.Bytes()[:16]).CryptBlocks(encMsg, encMsg)
	block := trimBlock(encMsg)
	if block == nil {
		return nil
	}

	enc := new(pb.EncryptedMessage)
	err = proto.Unmarshal(block, enc)
	if err != nil {
		return nil
	}

	ret := new(pb.RealMessage)
	err = proto.Unmarshal(enc.Msg, ret)
	if err != nil {
		return nil
	}

	x, y := elliptic.UnmarshalCompressed(curve, ret.SenderID)
	if x == nil || !ecdsa.Verify(
		&ecdsa.PublicKey{Curve: curve, X: x, Y: y}, enc.Msg,
		new(big.Int).SetBytes(enc.Signature[:curveByteLen]),
		new(big.Int).SetBytes(enc.Signature[curveByteLen:]),
	) {
		return nil
	}

	return &message{
		Sender:   newAddress(x, y),
		SendTime: time.UnixMilli(msg.Generated),
		Content:  ret.Message,
	}
}

func (u *user) NewMessage(recipient *address, msg string) (*pb.Message, error) {
	ephemeralKey, err := ecdsa.GenerateKey(curve, cReader)
	if err != nil {
		return nil, err
	}

	secret, iv := curve.ScalarMult(recipient.x, recipient.y, ephemeralKey.D.Bytes())
	cp, err := aes.NewCipher(secret.Bytes())
	if err != nil {
		return nil, err
	}

	rMsg := &pb.RealMessage{
		Message:  msg,
		SenderID: elliptic.MarshalCompressed(curve, u.key.X, u.key.Y),
	}
	rMsgBuf, err := proto.Marshal(rMsg)
	if err != nil {
		return nil, err
	}

	sigR, sigS, err := ecdsa.Sign(cReader, u.key, rMsgBuf)
	if err != nil {
		return nil, err
	}

	enc := pb.EncryptedMessage{
		Msg:       rMsgBuf,
		Signature: make([]byte, curveByteLen*2),
	}
	sigR.FillBytes(enc.Signature[:curveByteLen])
	sigS.FillBytes(enc.Signature[curveByteLen:])

	marshal, err := proto.Marshal(&enc)
	if err != nil {
		return nil, err
	}

	marshal = padBlock(marshal)
	cipher.NewCBCEncrypter(cp, iv.Bytes()[:16]).CryptBlocks(marshal, marshal)

	return &pb.Message{
		TargetCohort: recipient.cohort,
		Generated:    time.Now().UnixMilli(),
		EphemeralPub: elliptic.MarshalCompressed(curve, ephemeralKey.X, ephemeralKey.Y),
		EncMessage:   marshal,
	}, nil
}
