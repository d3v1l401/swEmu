package serializer

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"

	"../../crypto"
)

type PacketReader struct {
	reader *bytes.Reader
	Class  byte
	Type   byte
}

type PacketWriter struct {
	writer *bytes.Buffer
}

func NewWriter() *PacketWriter {
	instance := &PacketWriter{
		writer: bytes.NewBuffer(nil),
	}

	instance.WriteByte(0x02)

	return instance
}

func (p *PacketWriter) WriteInt(value int32) {
	binary.Write(p.writer, binary.BigEndian, value)
}

func (p *PacketWriter) WriteShort(value int16) {
	binary.Write(p.writer, binary.BigEndian, value)
}

func (p *PacketWriter) WriteByte(value byte) {
	binary.Write(p.writer, binary.BigEndian, value)
}

func (p *PacketWriter) WriteLong(value int64) {
	binary.Write(p.writer, binary.BigEndian, value)
}

func (p *PacketWriter) WriteString(value string) {
	binary.Write(p.writer, binary.BigEndian, len(value)+1)
	binary.Write(p.writer, binary.BigEndian, value)
}

func (p *PacketWriter) WriteStringUTF16(value string) {
	size := int16(len(value)*2) + 1
	if err := binary.Write(p.writer, binary.LittleEndian, &size); err != nil {
		fmt.Println(err)
	}
	for i := 0; i < len(value); i++ {
		binary.Write(p.writer, binary.BigEndian, value[i])
		p.WriteByte(0)
	}
	p.WriteByte(0)
}

func (p *PacketWriter) Fill(value byte, size int) {
	for i := 0; i < size; i++ {
		p.WriteByte(value)
	}
}

func (p *PacketWriter) GetBuffer() []byte {
	return p.writer.Bytes()
}

func (p *PacketWriter) Write(val string) {
	p.writer.Write([]byte(val))
}

func (p *PacketWriter) Finalize() []byte {
	return crypto.Decrypt(p.GetBuffer())
}

func NewReader(buffer []byte) *PacketReader {
	instance := &PacketReader{}

	if buffer != nil && len(buffer) > 5 {

		decryopted := crypto.Decrypt(buffer)
		instance.reader = bytes.NewReader(decryopted)
		instance.reader.Reset(decryopted)
		instance.reader.Seek(0, io.SeekStart)

		instance.Skip(2)
		instance.Class = instance.ReadByte()
		instance.Type = instance.ReadByte()
		instance.Skip(1)
	} else {
		return nil
	}

	return instance
}

func (p *PacketReader) ReadInt() int {
	var n int
	if err := binary.Read(p.reader, binary.LittleEndian, &n); err != nil {
		return 0
	}
	return n
}

func (p *PacketReader) ReadUInt() uint {
	var n uint
	if err := binary.Read(p.reader, binary.LittleEndian, &n); err != nil {
		return 0
	}
	return n

}

func (p *PacketReader) ReadShort() int16 {
	var n int16
	if err := binary.Read(p.reader, binary.LittleEndian, &n); err != nil {
		return 0
	}
	return n

}

func (p *PacketReader) ReadByte() byte {
	var n byte
	if err := binary.Read(p.reader, binary.LittleEndian, &n); err != nil {
		return 0
	}
	return n

}

func (p *PacketReader) ReadLong() uint64 {

	var n uint64
	if err := binary.Read(p.reader, binary.LittleEndian, &n); err != nil {
		return 0
	}
	return n

}

func (p *PacketReader) ReadString() string {

	size := p.ReadShort()
	if size > 1 {
		value := make([]byte, 0)
		var char byte
		for i := 0; i < int(size); i++ {
			if err := binary.Read(p.reader, binary.LittleEndian, &char); err != nil {
				return ""
			}
			value = append(value, char)
		}
		return string(value)
	}

	return ""
}

func (p *PacketReader) ReadStringUTF16() string {
	size := p.ReadShort()
	if size > 1 {
		value := make([]byte, 0)
		var char byte
		for i := 0; i < int(size); i++ {
			if err := binary.Read(p.reader, binary.LittleEndian, &char); err != nil {
				return ""
			}
			if char != 0x00 {
				value = append(value, char)
			}
		}

		return string(value)
	}

	return ""
}

func (p *PacketReader) Skip(size int) {
	for i := 0; i < size; i++ {
		p.ReadByte()
	}
}

func (p *PacketReader) Dump() {
	index, _ := p.reader.Seek(0, io.SeekCurrent)
	b, _ := ioutil.ReadAll(p.reader)
	p.reader.Seek(index, io.SeekStart)
	fmt.Printf("%s\n", hex.Dump(b))
}
