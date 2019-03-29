package main

import (
	"bytes"
	"encoding/binary"
)

// ForeignStatement describes the parameters we need to know to follow an arbitrary
// outside statement that's being witnessed on our connected server.
type ForeignStatement struct {
	// Was this the initial statement of this log? If so, the hash calculation
	// differs.
	InitialStatement bool

	// The Log's ID - needed to fetch the proof from the server
	LogID [32]byte

	// This is the actual statement we're trying to prove validity for
	StatementPreimage string

	// This is the signature on the statement
	Signature [64]byte

	// This is the public key that signed the statement
	PubKey [33]byte

	// Index is the sequential index of the statement in the log
	Index uint64

	// Proof is optional, used for historic proofs (we can fetch the current
	// proof from the server if it's meant to keep live).
	Proof []byte
}

// Bytes serializes a ForeignStatement object into a byte slice
func (f *ForeignStatement) Bytes() []byte {
	var b bytes.Buffer

	if f.InitialStatement {
		b.Write([]byte{0x01})
	} else {
		b.Write([]byte{0x00})
	}

	b.Write(f.LogID[:])
	b.Write(f.Signature[:])
	b.Write(f.PubKey[:])
	binary.Write(&b, binary.BigEndian, f.Index)
	binary.Write(&b, binary.BigEndian, uint32(len(f.StatementPreimage)))
	b.Write([]byte(f.StatementPreimage))
	if f.Proof == nil {
		binary.Write(&b, binary.BigEndian, uint32(0))
	} else {
		binary.Write(&b, binary.BigEndian, uint32(len(f.Proof)))
		b.Write(f.Proof)
	}

	return b.Bytes()
}

// ForeignStatementFromBytes deserializes a byte slice into a commitment object
func ForeignStatementFromBytes(b []byte) *ForeignStatement {
	f := ForeignStatement{}
	buf := bytes.NewBuffer(b)

	f.InitialStatement = bytes.Equal(buf.Next(1), []byte{0x01})
	copy(f.LogID[:], buf.Next(32))
	copy(f.Signature[:], buf.Next(64))
	copy(f.PubKey[:], buf.Next(33))

	binary.Read(buf, binary.BigEndian, &f.Index)

	iLen := uint32(0)
	binary.Read(buf, binary.BigEndian, &iLen)
	f.StatementPreimage = string(buf.Next(int(iLen)))

	binary.Read(buf, binary.BigEndian, &iLen)
	if iLen > 0 {
		f.Proof = buf.Next(int(iLen))
	}

	return &f
}
