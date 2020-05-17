package peinfo

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
)

type dosHeader struct {
	MagicNumber                uint16
	BytesOnLastPage            uint16
	PagesInFile                uint16
	Relocations                uint16
	SizeOfHeaderInParagraphs   uint16
	MinExtraParagraphsNeeded   uint16
	MaxExtraParagraphsNeeded   uint16
	SsValue                    uint16
	SpValue                    uint16
	Checksum                   uint16
	IpValue                    uint16
	CsValue                    uint16
	FileAddressRelocationTable uint16
	OverlayNumber              uint16
	ReservedWords              [8]uint8
	OemIdentifier              uint16
	OemInfo                    uint16
	ReservedWords2             [20]uint8
	FileAddressNewExeHeader    uint32
}

type fileHeader struct {
	Machine          uint16
	NumberOfSections uint16
	TimeDateStamp    uint32
}

type optionalHeader struct {
}

type ntHeader struct {
	Signature      uint32
	FileHeader     fileHeader
	OptionalHeader optionalHeader
}

type peinfo struct {
	DosHeader *dosHeader
	NtHeader  *ntHeader
}

func toBytes(data interface{}) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, data)
	if err != nil {
		panic(fmt.Sprintf("binary.Write failed: %v", err))
	}
	return buf.Bytes()
}

func Print(name string, value interface{}, asciiBytes []byte) {
	fmt.Printf("%-30v", name)
	fmt.Printf("%-15v", value)
	fmt.Printf("%-30v\n", string(asciiBytes))
}

func (h *dosHeader) Print() {
	fmt.Println("---DOS Header---")
	Print("PARAMETER", "VALUE(DEC)", []byte("ASCII"))
	Print("MagicNumber", h.MagicNumber, toBytes(h.MagicNumber))
	Print("BytesOnLastPage", h.BytesOnLastPage, toBytes(h.BytesOnLastPage))
	Print("PagesInFile", h.PagesInFile, toBytes(h.PagesInFile))
	Print("Relocations", h.Relocations, toBytes(h.Relocations))
	Print("SizeOfHeaderInParagraphs", h.SizeOfHeaderInParagraphs, toBytes(h.SizeOfHeaderInParagraphs))
	Print("MinExtraParagraphsNeeded", h.MinExtraParagraphsNeeded, toBytes(h.MinExtraParagraphsNeeded))
	Print("MaxExtraParagraphsNeeded", h.MaxExtraParagraphsNeeded, toBytes(h.MaxExtraParagraphsNeeded))
	Print("SsValue", h.SsValue, toBytes(h.SsValue))
	Print("SpValue", h.SpValue, toBytes(h.SpValue))
	Print("Checksum", h.Checksum, toBytes(h.Checksum))
	Print("IpValue", h.IpValue, toBytes(h.IpValue))
	Print("CsValue", h.CsValue, toBytes(h.CsValue))
	Print("FileAddressRelocationTable", h.FileAddressRelocationTable, toBytes(h.FileAddressRelocationTable))
	// Print("ReservedWords", h.ReservedWords, toBytes(h.ReservedWords))
	Print("OemIdentifier", h.OemIdentifier, toBytes(h.OemIdentifier))
	Print("OemInfo", h.OemInfo, toBytes(h.OemInfo))
	// Print("ReservedWords2", h.ReservedWords2, toBytes(h.ReservedWords2))
	Print("FileAddressNewExeHeader", h.FileAddressNewExeHeader, toBytes(h.FileAddressNewExeHeader))
}

func (p *peinfo) Print() {
	p.DosHeader.Print()
}

func New(peFilePath string) (*peinfo, error) {
	peInfo, err := parsePeFormat(peFilePath)
	if err != nil {
		return nil, err
	}

	return peInfo, err
}

const IMAGE_DOS_SIGNATURE = 0x5A4D // "MZ"

func parsePeFormat(filePath string) (*peinfo, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	peInfo := &peinfo{}
	peInfo.DosHeader, err = parseDosHeader(f)
	if err != nil {
		return nil, err
	}

	parseNtHeader(peInfo.DosHeader, f)

	return peInfo, nil
}

func parseDosHeader(file *os.File) (*dosHeader, error) {
	dosHeader := &dosHeader{}
	err := binary.Read(file, binary.LittleEndian, dosHeader)
	if err != nil {
		return nil, err
	}

	if dosHeader.MagicNumber != IMAGE_DOS_SIGNATURE {
		return nil, errors.New("The file does not have a valid PE-format.")
	}

	return dosHeader, nil
}

func parseNtHeader(dosHeader *dosHeader, file *os.File) (*ntHeader, error) {
	// Set offset to location for the NT header
	ntHeaderLocation := int64(dosHeader.FileAddressNewExeHeader)
	newOffset, errs := file.Seek(ntHeaderLocation, 0)
	if errs != nil {
		panic(errs)
	}
	if newOffset != ntHeaderLocation {
		panic("Unable to seek correct NT header location")
	}

	ntHeader := &ntHeader{}
	err := binary.Read(file, binary.LittleEndian, ntHeader)
	if err != nil {
		return nil, err
	}

	fmt.Println(ntHeader)

	return nil, nil
}

func readUint16(file *os.File) uint16 {
	b := make([]byte, 2)
	n, e := file.Read(b)
	if e != nil {
		panic(e)
	}
	if n != 2 {
		panic("Unable to read unit16")
	}

	return binary.LittleEndian.Uint16(b)
}

// func readUint32(file *os.File) uint32 {
// 	b := make([]byte, 4)
// 	n, e := file.Read(b)
// 	if e != nil {
// 		panic(e)
// 	}
// 	if n != 4 {
// 		panic("Unable to read unit32")
// 	}
// 	return binary.LittleEndian.Uint32(b)
// }

// func test(file *os.File) {
// 	// bufSize := 4
// 	// b := make([]byte, bufSize)
// 	// n, e := file.Read(b)
// 	// if e != nil {
// 	// 	panic(e)
// 	// }
// 	// if n != bufSize {
// 	// 	panic(fmt.Sprintf("Unable to read buf with size %v", bufSize))
// 	// }

// 	t := header{}

// 	err := binary.Read(file, binary.LittleEndian, &t)
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println("TESTING")
// 	fmt.Println(t)
// 	// headerTest
// }

// func readUint16(file *os.File) (uint16, error) {
// 	b := make([]byte, 2)
// 	n, e := file.Read(b)

// 	if e != nil {
// 		return 0, e
// 	}
// 	if n != 2 {
// 		return 0, errors.New("Unable to read unit16")
// 	}

// 	return binary.LittleEndian.Uint16(b), nil
// }
