package solana

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/lunixbochs/struc"
	"github.com/mr-tron/base58"
)

type Padding []byte

type Hash PublicKey

///

type Signature [32]byte

func SignatureFromBase58(in string) (out Signature, err error) {
	val, err := base58.Decode(in)
	if err != nil {
		return
	}

	if len(val) != 64 {
		err = fmt.Errorf("invalid length, expected 64, got %d", len(val))
		return
	}
	copy(out[:], val)
	return
}

func (p Signature) MarshalJSON() ([]byte, error) {
	return json.Marshal(base58.Encode(p[:]))
}
func (p *Signature) UnmarshalJSON(data []byte) (err error) {
	var s string
	err = json.Unmarshal(data, &s)
	if err != nil {
		return
	}

	dat, err := base58.Decode(s)
	if err != nil {
		return err
	}

	if len(dat) != 64 {
		return errors.New("invalid data length for public key")
	}

	target := Signature{}
	copy(target[:], dat)
	*p = target
	return
}

func (p Signature) String() string {
	return base58.Encode(p[:])
}

///

type PublicKey [32]byte

func PublicKeyFromBase58(in string) (out PublicKey, err error) {
	val, err := base58.Decode(in)
	if err != nil {
		return
	}

	if len(val) != 32 {
		err = fmt.Errorf("invalid length, expected 32, got %d", len(val))
		return
	}
	copy(out[:], val)
	return
}

func (p PublicKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(base58.Encode(p[:]))
}
func (p *PublicKey) UnmarshalJSON(data []byte) (err error) {
	var s string
	err = json.Unmarshal(data, &s)
	if err != nil {
		return
	}

	dat, err := base58.Decode(s)
	if err != nil {
		return err
	}

	if len(dat) != 32 {
		return errors.New("invalid data length for public key")
	}

	target := PublicKey{}
	copy(target[:], dat)
	*p = target
	return
}

func (p PublicKey) String() string {
	return base58.Encode(p[:])
}

///

type Base58 []byte

func (t Base58) MarshalJSON() ([]byte, error) {
	return json.Marshal(base58.Encode(t))
}

func (t *Base58) UnmarshalJSON(data []byte) (err error) {
	var s string
	err = json.Unmarshal(data, &s)
	if err != nil {
		return
	}

	*t, err = base58.Decode(s)
	return
}

func (t Base58) String() string {
	return base58.Encode(t)
}

///

type U64 uint64

func (i U64) MarshalJSON() (data []byte, err error) {
	if i > 0xffffffff {
		encodedInt, err := json.Marshal(uint64(i))
		if err != nil {
			return nil, err
		}
		data = append([]byte{'"'}, encodedInt...)
		data = append(data, '"')
		return data, nil
	}
	return json.Marshal(uint64(i))
}

func (i *U64) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return errors.New("empty value")
	}

	if data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}

		val, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return err
		}

		*i = U64(val)

		return nil
	}

	var v uint64
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*i = U64(v)

	return nil
}

///

type ByteWrapper struct {
	io.Reader
}

func (w *ByteWrapper) ReadByte() (byte, error) {
	var b [1]byte
	_, err := w.Read(b[:])
	return b[0], err
}

/// ShortVec
type ShortVec uint16

func (v ShortVec) Pack(p []byte, opt *struc.Options) (int, error) {
	// JAVASCRIPT
	// let rem_len = len;
	// for (;;) {
	//   let elem = rem_len & 0x7f;
	//   rem_len >>= 7;
	//   if (rem_len == 0) {
	//     bytes.push(elem);
	//     break;
	//   } else {
	//     elem |= 0x80;
	//     bytes.push(elem);
	//   }
	// }

	// RUST
	//     // Pass a non-zero value to serialize_tuple() so that serde_json will
	// // generate an open bracket.
	// let mut seq = serializer.serialize_tuple(1)?;

	// let mut rem_len = self.0;
	// loop {
	//     let mut elem = (rem_len & 0x7f) as u8;
	//     rem_len >>= 7;
	//     if rem_len == 0 {
	//         seq.serialize_element(&elem)?;
	//         break;
	//     } else {
	//         elem |= 0x80;
	//         seq.serialize_element(&elem)?;
	//     }
	// }
	// seq.end()

	return 0, nil
}
func (v *ShortVec) Unpack(r io.Reader, length int, opt *struc.Options) error {
	var l, s int
	for {

		// JAVASCRIPT
		//   let elem = bytes.shift();
		//   len |= (elem & 0x7f) << (size * 7);
		//   size += 1;
		//   if ((elem & 0x80) === 0) {
		//     break;
		//   }
	}

	// RUST
	// let mut len: usize = 0;
	// let mut size: usize = 0;
	// loop {
	//     let elem: u8 = seq
	//         .next_element()?
	//         .ok_or_else(|| de::Error::invalid_length(size, &self))?;

	//     len |= (elem as usize & 0x7f) << (size * 7);
	//     size += 1;

	//     if elem as usize & 0x80 == 0 {
	//         break;
	//     }

	//     if size > size_of::<u16>() + 1 {
	//         return Err(de::Error::invalid_length(size, &self));
	//     }
	// }

	// Ok(ShortU16(len as u16))

	// TODO: have a func that would return `size` also, separately? called twice?
	return nil
}
func (v *ShortVec) Size(opt *struc.Options) int {
	var buf [8]byte
	return binary.PutUvarint(buf[:], uint64(*v))
}
func (v *ShortVec) String() string {
	return strconv.FormatUint(uint64(*v), 10)
}

var shortVecOverflow = errors.New("short_vec: varint overflows a 16-bit integer")

func readShortVec(r io.ByteReader) (uint64, error) {
	var x uint64
	var s uint
	for i := 0; ; i++ {
		b, err := r.ReadByte()
		if err != nil {
			return x, err
		}
		if b < 0x80 {
			if i > 4 || i == 4 && b > 1 {
				return x, shortVecOverflow
			}
			return x | uint64(b)<<s, nil
		}
		x |= uint64(b&0x7f) << s
		s += 7
	}
}
