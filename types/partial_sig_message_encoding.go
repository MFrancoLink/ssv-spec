// Code generated by fastssz. DO NOT EDIT.
// Hash: ae91558f5b93329d35593cd02e0186d6707d167466df4a3c2d68f8324748ce62
// Version: 0.1.3
package types

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"
)

// MarshalSSZ ssz marshals the PartialSignatureMessages object
func (p *PartialSignatureMessages) MarshalSSZ() ([]byte, error) {
	return ssz.MarshalSSZ(p)
}

// MarshalSSZTo ssz marshals the PartialSignatureMessages object to a target array
func (p *PartialSignatureMessages) MarshalSSZTo(buf []byte) (dst []byte, err error) {
	dst = buf
	offset := int(20)

	// Field (0) 'Type'
	dst = ssz.MarshalUint64(dst, uint64(p.Type))

	// Field (1) 'Slot'
	dst = ssz.MarshalUint64(dst, uint64(p.Slot))

	// Offset (2) 'Messages'
	dst = ssz.WriteOffset(dst, offset)
	offset += len(p.Messages) * 144

	// Field (2) 'Messages'
	if size := len(p.Messages); size > 13 {
		err = ssz.ErrListTooBigFn("PartialSignatureMessages.Messages", size, 13)
		return
	}
	for ii := 0; ii < len(p.Messages); ii++ {
		if dst, err = p.Messages[ii].MarshalSSZTo(dst); err != nil {
			return
		}
	}

	return
}

// UnmarshalSSZ ssz unmarshals the PartialSignatureMessages object
func (p *PartialSignatureMessages) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size < 20 {
		return ssz.ErrSize
	}

	tail := buf
	var o2 uint64

	// Field (0) 'Type'
	p.Type = PartialSigMsgType(ssz.UnmarshallUint64(buf[0:8]))

	// Field (1) 'Slot'
	p.Slot = phase0.Slot(ssz.UnmarshallUint64(buf[8:16]))

	// Offset (2) 'Messages'
	if o2 = ssz.ReadOffset(buf[16:20]); o2 > size {
		return ssz.ErrOffset
	}

	if o2 < 20 {
		return ssz.ErrInvalidVariableOffset
	}

	// Field (2) 'Messages'
	{
		buf = tail[o2:]
		num, err := ssz.DivideInt2(len(buf), 144, 13)
		if err != nil {
			return err
		}
		p.Messages = make([]*PartialSignatureMessage, num)
		for ii := 0; ii < num; ii++ {
			if p.Messages[ii] == nil {
				p.Messages[ii] = new(PartialSignatureMessage)
			}
			if err = p.Messages[ii].UnmarshalSSZ(buf[ii*144 : (ii+1)*144]); err != nil {
				return err
			}
		}
	}
	return err
}

// SizeSSZ returns the ssz encoded size in bytes for the PartialSignatureMessages object
func (p *PartialSignatureMessages) SizeSSZ() (size int) {
	size = 20

	// Field (2) 'Messages'
	size += len(p.Messages) * 144

	return
}

// HashTreeRoot ssz hashes the PartialSignatureMessages object
func (p *PartialSignatureMessages) HashTreeRoot() ([32]byte, error) {
	return ssz.HashWithDefaultHasher(p)
}

// HashTreeRootWith ssz hashes the PartialSignatureMessages object with a hasher
func (p *PartialSignatureMessages) HashTreeRootWith(hh ssz.HashWalker) (err error) {
	indx := hh.Index()

	// Field (0) 'Type'
	hh.PutUint64(uint64(p.Type))

	// Field (1) 'Slot'
	hh.PutUint64(uint64(p.Slot))

	// Field (2) 'Messages'
	{
		subIndx := hh.Index()
		num := uint64(len(p.Messages))
		if num > 13 {
			err = ssz.ErrIncorrectListSize
			return
		}
		for _, elem := range p.Messages {
			if err = elem.HashTreeRootWith(hh); err != nil {
				return
			}
		}
		hh.MerkleizeWithMixin(subIndx, num, 13)
	}

	hh.Merkleize(indx)
	return
}

// GetTree ssz hashes the PartialSignatureMessages object
func (p *PartialSignatureMessages) GetTree() (*ssz.Node, error) {
	return ssz.ProofTree(p)
}

// MarshalSSZ ssz marshals the PartialSignatureMessage object
func (p *PartialSignatureMessage) MarshalSSZ() ([]byte, error) {
	return ssz.MarshalSSZ(p)
}

// MarshalSSZTo ssz marshals the PartialSignatureMessage object to a target array
func (p *PartialSignatureMessage) MarshalSSZTo(buf []byte) (dst []byte, err error) {
	dst = buf

	// Field (0) 'PartialSignature'
	if size := len(p.PartialSignature); size != 96 {
		err = ssz.ErrBytesLengthFn("PartialSignatureMessage.PartialSignature", size, 96)
		return
	}
	dst = append(dst, p.PartialSignature...)

	// Field (1) 'SigningRoot'
	dst = append(dst, p.SigningRoot[:]...)

	// Field (2) 'Signer'
	dst = ssz.MarshalUint64(dst, uint64(p.Signer))

	// Field (3) 'ValidatorIndex'
	dst = ssz.MarshalUint64(dst, uint64(p.ValidatorIndex))

	return
}

// UnmarshalSSZ ssz unmarshals the PartialSignatureMessage object
func (p *PartialSignatureMessage) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size != 144 {
		return ssz.ErrSize
	}

	// Field (0) 'PartialSignature'
	if cap(p.PartialSignature) == 0 {
		p.PartialSignature = make([]byte, 0, len(buf[0:96]))
	}
	p.PartialSignature = append(p.PartialSignature, buf[0:96]...)

	// Field (1) 'SigningRoot'
	copy(p.SigningRoot[:], buf[96:128])

	// Field (2) 'Signer'
	p.Signer = OperatorID(ssz.UnmarshallUint64(buf[128:136]))

	// Field (3) 'ValidatorIndex'
	p.ValidatorIndex = phase0.ValidatorIndex(ssz.UnmarshallUint64(buf[136:144]))

	return err
}

// SizeSSZ returns the ssz encoded size in bytes for the PartialSignatureMessage object
func (p *PartialSignatureMessage) SizeSSZ() (size int) {
	size = 144
	return
}

// HashTreeRoot ssz hashes the PartialSignatureMessage object
func (p *PartialSignatureMessage) HashTreeRoot() ([32]byte, error) {
	return ssz.HashWithDefaultHasher(p)
}

// HashTreeRootWith ssz hashes the PartialSignatureMessage object with a hasher
func (p *PartialSignatureMessage) HashTreeRootWith(hh ssz.HashWalker) (err error) {
	indx := hh.Index()

	// Field (0) 'PartialSignature'
	if size := len(p.PartialSignature); size != 96 {
		err = ssz.ErrBytesLengthFn("PartialSignatureMessage.PartialSignature", size, 96)
		return
	}
	hh.PutBytes(p.PartialSignature)

	// Field (1) 'SigningRoot'
	hh.PutBytes(p.SigningRoot[:])

	// Field (2) 'Signer'
	hh.PutUint64(uint64(p.Signer))

	// Field (3) 'ValidatorIndex'
	hh.PutUint64(uint64(p.ValidatorIndex))

	hh.Merkleize(indx)
	return
}

// GetTree ssz hashes the PartialSignatureMessage object
func (p *PartialSignatureMessage) GetTree() (*ssz.Node, error) {
	return ssz.ProofTree(p)
}
