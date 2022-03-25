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

// Message contains a decrypted message (sent to the user).
type Message struct {
	// Sender is the sender of the message.
	Sender *Address
	// SendTime is the sender-specified time.
	SendTime time.Time
	// Content contains the protocol-specified message data.
	// It is a UTF-8 encoded string in vanilla implementation.
	Content []byte
}

func (u *User) decryptMessage(msg *pb.Message) *Message {
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
	cipher.NewCBCDecrypter(cp, iv.Bytes()[8:][:aes.BlockSize]).CryptBlocks(encMsg, encMsg)
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
	h := hasher(enc.Msg)
	if x == nil || !ecdsa.Verify(
		&ecdsa.PublicKey{Curve: curve, X: x, Y: y}, h[:],
		new(big.Int).SetBytes(enc.Signature[:curveByteLen]),
		new(big.Int).SetBytes(enc.Signature[curveByteLen:]),
	) {
		return nil
	}

	return &Message{
		Sender:   newAddress(x, y),
		SendTime: time.UnixMilli(ret.SendTime),
		Content:  ret.Message,
	}
}

// NewMessage send a new message with content msg to the recipient.
// Database.StoreMessage might be called synchronously.
func (u *User) NewMessage(recipient *Address, msg []byte) error {
	ephemeralKey, err := ecdsa.GenerateKey(curve, cReader)
	if err != nil {
		return err
	}

	secret, iv := curve.ScalarMult(recipient.x, recipient.y, ephemeralKey.D.Bytes())
	cp, err := aes.NewCipher(secret.Bytes())
	if err != nil {
		return err
	}

	rMsg := &pb.RealMessage{
		Message:  msg,
		SendTime: time.Now().UnixMilli(),
		SenderID: elliptic.MarshalCompressed(curve, u.key.X, u.key.Y),
	}
	rMsgBuf, err := proto.Marshal(rMsg)
	if err != nil {
		return err
	}

	h := hasher(rMsgBuf)
	sigR, sigS, err := ecdsa.Sign(cReader, u.key, h[:])
	if err != nil {
		return err
	}

	enc := pb.EncryptedMessage{
		Msg:       rMsgBuf,
		Signature: make([]byte, curveByteLen*2),
	}
	sigR.FillBytes(enc.Signature[:curveByteLen])
	sigS.FillBytes(enc.Signature[curveByteLen:])

	marshal, err := proto.Marshal(&enc)
	if err != nil {
		return err
	}

	marshal = padBlock(marshal)
	cipher.NewCBCEncrypter(cp, iv.Bytes()[8:][:aes.BlockSize]).CryptBlocks(marshal, marshal)

	mMsg, err := proto.Marshal(&pb.Message{
		TargetCohort: recipient.cohort,
		EphemeralPub: elliptic.MarshalCompressed(curve, ephemeralKey.X, ephemeralKey.Y),
		EncMessage:   marshal,
	})
	if err != nil {
		return err
	}

	msgHash := hasher(mMsg)
	return u.db.StoreMessage(msgHash, &pb.MsgContainer{
		Msg: mMsg,
		Pow: calcPoW(msgHash[:]),
	}, func() (uint32, error) { return recipient.cohort, nil })
}
