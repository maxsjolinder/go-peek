package peinfo

import (
	"encoding/binary"
	"errors"
	"os"
)

const IMAGE_DOS_SIGNATURE = 0x5A4D           // "MZ"
const OPTIONAL_HEADER_MAGIC_PE = 0x010b      // 32 bit
const OPTIONAL_HEADER_MAGIC_PE_PLUS = 0x020b // 64 bit

func New(peFilePath string) (*peinfo, error) {
	peInfo, err := parsePeFormat(peFilePath)
	if err != nil {
		return nil, err
	}

	return peInfo, err
}

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
