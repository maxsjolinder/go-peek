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
	Machine              uint16
	NumberOfSections     uint16
	TimeDateStamp        uint32
	PointerToSymbolTable uint32
	NumberOfSymbols      uint32
	SizeOfOptionalHeader uint16
	Characteristics      uint16
}

// 32-bit
type optionalHeader32 struct {
	Magic                       uint16
	MajorLinkerVersion          byte
	MinorLinkerVersion          byte
	SizeOfCode                  uint32
	SizeOfInitializedData       uint32
	SizeOfUninitializedData     uint32
	AddressOfEntryPoint         uint32
	BaseOfCode                  uint32
	BaseOfData                  uint32
	ImageBase                   uint32
	SectionAlignment            uint32
	FileAlignment               uint32
	MajorOperatingSystemVersion uint16
	MinorOperatingSystemVersion uint16
	MajorImageVersion           uint16
	MinorImageVersion           uint16
	MajorSubsystemVersion       uint16
	MinorSubsystemVersion       uint16
	Win32VersionValue           uint32
	ImageSize                   uint32
	SizeOfHeaders               uint32
	CheckSum                    uint32
	Subsystem                   uint16
	DllCharacteristics          uint16
	SizeOfStackReserve          uint32
	SizeOfStackCommit           uint32
	SizeOfHeapReserve           uint32
	SizeOfHeapCommit            uint32
	LoaderFlags                 uint32
	NumberOfRvaAndSizes         uint32
}

// 64-bit
type optionalHeader64 struct {
	Magic                       uint16
	MajorLinkerVersion          byte
	MinorLinkerVersion          byte
	SizeOfCode                  uint32
	SizeOfInitializedData       uint32
	SizeOfUninitializedData     uint32
	AddressOfEntryPoint         uint32
	BaseOfCode                  uint32
	ImageBase                   uint64
	SectionAlignment            uint32
	FileAlignment               uint32
	MajorOperatingSystemVersion uint16
	MinorOperatingSystemVersion uint16
	MajorImageVersion           uint16
	MinorImageVersion           uint16
	MajorSubsystemVersion       uint16
	MinorSubsystemVersion       uint16
	Win32VersionValue           uint32
	ImageSize                   uint32
	SizeOfHeaders               uint32
	CheckSum                    uint32
	Subsystem                   uint16
	DllCharacteristics          uint16
	SizeOfStackReserve          uint64
	SizeOfStackCommit           uint64
	SizeOfHeapReserve           uint64
	SizeOfHeapCommit            uint64
	LoaderFlags                 uint32
	NumberOfRvaAndSizes         uint32
}

type ntHeader struct {
	Signature        uint32
	FileHeader       *fileHeader
	OptionalHeader32 *optionalHeader32
	OptionalHeader64 *optionalHeader64
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
	fmt.Printf("%-15x", value)
	fmt.Printf("%-30v\n", string(asciiBytes))
}

func PrintTableHeader() {
	fmt.Printf("%-30v", "PARAMETER")
	fmt.Printf("%-15v", "VALUE(DEC)")
	fmt.Printf("%-15v", "VALUE(HEX)")
	fmt.Printf("%-30v\n", "ASCII")
}

func (fh *fileHeader) Print() {
	fmt.Println("-File Header-")
	PrintTableHeader()
	Print("Machine", fh.Machine, toBytes(fh.Machine))
	Print("NumberOfSections", fh.NumberOfSections, toBytes(fh.NumberOfSections))
	Print("TimeDateStamp", fh.TimeDateStamp, toBytes(fh.TimeDateStamp))
	Print("PointerToSymbolTable", fh.PointerToSymbolTable, toBytes(fh.PointerToSymbolTable))
	Print("NumberOfSymbols", fh.NumberOfSymbols, toBytes(fh.NumberOfSymbols))
	Print("SizeOfOptionalHeader", fh.SizeOfOptionalHeader, toBytes(fh.SizeOfOptionalHeader))
	Print("Characteristics", fh.Characteristics, toBytes(fh.Characteristics))
}

func (oh *optionalHeader32) Print() {
	fmt.Println("-Optional Header(32-bit)-")
	PrintTableHeader()
	Print("Magic", oh.Magic, toBytes(oh.Magic))
	Print("MajorLinkerVersion", oh.MajorLinkerVersion, toBytes(oh.MajorLinkerVersion))
	Print("MinorLinkerVersion", oh.MinorLinkerVersion, toBytes(oh.MinorLinkerVersion))
	Print("SizeOfCode", oh.SizeOfCode, toBytes(oh.SizeOfCode))
	Print("SizeOfInitializedData", oh.SizeOfInitializedData, toBytes(oh.SizeOfInitializedData))
	Print("SizeOfUninitializedData", oh.SizeOfUninitializedData, toBytes(oh.SizeOfUninitializedData))
	Print("AddressOfEntryPoint", oh.AddressOfEntryPoint, toBytes(oh.AddressOfEntryPoint))
	Print("BaseOfCode", oh.BaseOfCode, toBytes(oh.BaseOfCode))
	Print("BaseOfData", oh.BaseOfData, toBytes(oh.BaseOfData))
	Print("ImageBase", oh.ImageBase, toBytes(oh.ImageBase))
	Print("SectionAlignment", oh.SectionAlignment, toBytes(oh.SectionAlignment))
	Print("FileAlignment", oh.FileAlignment, toBytes(oh.FileAlignment))
	Print("MajorOperatingSystemVersion", oh.MajorOperatingSystemVersion, toBytes(oh.MajorOperatingSystemVersion))
	Print("MinorOperatingSystemVersion", oh.MinorOperatingSystemVersion, toBytes(oh.MinorOperatingSystemVersion))
	Print("MajorImageVersion", oh.MajorImageVersion, toBytes(oh.MajorImageVersion))
	Print("MinorImageVersion", oh.MinorImageVersion, toBytes(oh.MinorImageVersion))
	Print("MajorSubsystemVersion", oh.MajorSubsystemVersion, toBytes(oh.MajorSubsystemVersion))
	Print("MinorSubsystemVersion", oh.MinorSubsystemVersion, toBytes(oh.MinorSubsystemVersion))
	Print("Win32VersionValue", oh.Win32VersionValue, toBytes(oh.Win32VersionValue))
	Print("ImageSize", oh.ImageSize, toBytes(oh.ImageSize))
	Print("SizeOfHeaders", oh.SizeOfHeaders, toBytes(oh.SizeOfHeaders))
	Print("CheckSum", oh.CheckSum, toBytes(oh.CheckSum))
	Print("Subsystem", oh.Subsystem, toBytes(oh.Subsystem))
	Print("DllCharacteristics", oh.DllCharacteristics, toBytes(oh.DllCharacteristics))
	Print("SizeOfStackReserve", oh.SizeOfStackReserve, toBytes(oh.SizeOfStackReserve))
	Print("SizeOfStackCommit", oh.SizeOfStackCommit, toBytes(oh.SizeOfStackCommit))
	Print("SizeOfHeapReserve", oh.SizeOfHeapReserve, toBytes(oh.SizeOfHeapReserve))
	Print("SizeOfHeapCommit", oh.SizeOfHeapCommit, toBytes(oh.SizeOfHeapCommit))
	Print("LoaderFlags", oh.LoaderFlags, toBytes(oh.LoaderFlags))
	Print("NumberOfRvaAndSizes", oh.NumberOfRvaAndSizes, toBytes(oh.NumberOfRvaAndSizes))
}

func (oh *optionalHeader64) Print() {
	fmt.Println("-Optional Header(64-bit)-")
	PrintTableHeader()
	Print("Magic", oh.Magic, toBytes(oh.Magic))
	Print("MajorLinkerVersion", oh.MajorLinkerVersion, toBytes(oh.MajorLinkerVersion))
	Print("MinorLinkerVersion", oh.MinorLinkerVersion, toBytes(oh.MinorLinkerVersion))
	Print("SizeOfCode", oh.SizeOfCode, toBytes(oh.SizeOfCode))
	Print("SizeOfInitializedData", oh.SizeOfInitializedData, toBytes(oh.SizeOfInitializedData))
	Print("SizeOfUninitializedData", oh.SizeOfUninitializedData, toBytes(oh.SizeOfUninitializedData))
	Print("AddressOfEntryPoint", oh.AddressOfEntryPoint, toBytes(oh.AddressOfEntryPoint))
	Print("BaseOfCode", oh.BaseOfCode, toBytes(oh.BaseOfCode))
	Print("ImageBase", oh.ImageBase, toBytes(oh.ImageBase))
	Print("SectionAlignment", oh.SectionAlignment, toBytes(oh.SectionAlignment))
	Print("FileAlignment", oh.FileAlignment, toBytes(oh.FileAlignment))
	Print("MajorOperatingSystemVersion", oh.MajorOperatingSystemVersion, toBytes(oh.MajorOperatingSystemVersion))
	Print("MinorOperatingSystemVersion", oh.MinorOperatingSystemVersion, toBytes(oh.MinorOperatingSystemVersion))
	Print("MajorImageVersion", oh.MajorImageVersion, toBytes(oh.MajorImageVersion))
	Print("MinorImageVersion", oh.MinorImageVersion, toBytes(oh.MinorImageVersion))
	Print("MajorSubsystemVersion", oh.MajorSubsystemVersion, toBytes(oh.MajorSubsystemVersion))
	Print("MinorSubsystemVersion", oh.MinorSubsystemVersion, toBytes(oh.MinorSubsystemVersion))
	Print("Win32VersionValue", oh.Win32VersionValue, toBytes(oh.Win32VersionValue))
	Print("ImageSize", oh.ImageSize, toBytes(oh.ImageSize))
	Print("SizeOfHeaders", oh.SizeOfHeaders, toBytes(oh.SizeOfHeaders))
	Print("CheckSum", oh.CheckSum, toBytes(oh.CheckSum))
	Print("Subsystem", oh.Subsystem, toBytes(oh.Subsystem))
	Print("DllCharacteristics", oh.DllCharacteristics, toBytes(oh.DllCharacteristics))
	Print("SizeOfStackReserve", oh.SizeOfStackReserve, toBytes(oh.SizeOfStackReserve))
	Print("SizeOfStackCommit", oh.SizeOfStackCommit, toBytes(oh.SizeOfStackCommit))
	Print("SizeOfHeapReserve", oh.SizeOfHeapReserve, toBytes(oh.SizeOfHeapReserve))
	Print("SizeOfHeapCommit", oh.SizeOfHeapCommit, toBytes(oh.SizeOfHeapCommit))
	Print("LoaderFlags", oh.LoaderFlags, toBytes(oh.LoaderFlags))
	Print("NumberOfRvaAndSizes", oh.NumberOfRvaAndSizes, toBytes(oh.NumberOfRvaAndSizes))
}

func (h *ntHeader) Print() {
	fmt.Println("---NT Header---")
	PrintTableHeader()
	Print("Signature", h.Signature, toBytes(h.Signature))
	h.FileHeader.Print()
	if h.OptionalHeader32 != nil {
		h.OptionalHeader32.Print()
	} else if h.OptionalHeader64 != nil {
		h.OptionalHeader64.Print()
	}
}

func (h *dosHeader) Print() {
	fmt.Println("---DOS Header---")
	PrintTableHeader()
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
	p.NtHeader.Print()
}

func New(peFilePath string) (*peinfo, error) {
	peInfo, err := parsePeFormat(peFilePath)
	if err != nil {
		return nil, err
	}

	return peInfo, err
}

const IMAGE_DOS_SIGNATURE = 0x5A4D           // "MZ"
const OPTIONAL_HEADER_MAGIC_PE = 0x010b      // 32 bit
const OPTIONAL_HEADER_MAGIC_PE_PLUS = 0x020b // 64 bit

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

	peInfo.NtHeader, err = parseNtHeader(peInfo.DosHeader, f)
	if err != nil {
		return nil, err
	}

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

	// Read Signature
	ntHeader := &ntHeader{}
	ntHeader.Signature = readUint32(file)

	// Read file header
	fileHeader := &fileHeader{}
	err := binary.Read(file, binary.LittleEndian, fileHeader)
	if err != nil {
		return nil, err
	}
	ntHeader.FileHeader = fileHeader

	// Read 32-bit optional header in order to read the magic number
	// to check whether we should read it as a 64-header instead.
	optionalHeaderOffset, err := file.Seek(0, 1)
	if err != nil {
		return nil, err
	}
	optionalHeader32 := &optionalHeader32{}
	err = binary.Read(file, binary.LittleEndian, optionalHeader32)
	if err != nil {
		return nil, err
	}

	if optionalHeader32.Magic == OPTIONAL_HEADER_MAGIC_PE {
		// 32 bit header
		ntHeader.OptionalHeader32 = optionalHeader32
	} else if optionalHeader32.Magic == OPTIONAL_HEADER_MAGIC_PE_PLUS {
		// 64 bit, we must re-read the data into the correct struct
		// Set file offset to location of optional header
		_, err := file.Seek(optionalHeaderOffset, 0)
		if err != nil {
			return nil, err
		}

		optionalHeader64 := &optionalHeader64{}
		err = binary.Read(file, binary.LittleEndian, optionalHeader64)
		if err != nil {
			return nil, err
		}

		ntHeader.OptionalHeader64 = optionalHeader64
	} else {
		return nil, errors.New("Error in file header")
	}

	return ntHeader, nil
}

func readUint32(file *os.File) uint32 {
	b := make([]byte, 4)
	n, e := file.Read(b)
	if e != nil {
		panic(e)
	}
	if n != 4 {
		panic("Unable to read unit32")
	}
	return binary.LittleEndian.Uint32(b)
}
