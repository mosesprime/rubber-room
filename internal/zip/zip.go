// Use at your own risk.

package zip

import (
	"bufio"
	"hash"
	"io"
)

// compression methods
const (
    Store uint16 = 0
    Deflate uint16 = 8
)

const (
    fileHeadSign = 0x04034B50
    dirHeadSign = 0x02014B50
    dirEndSign = 0x06054B50
    dir64LocSign = 0x07064B50
    dir64EndSign = 0x06064B50
    dataDescriptorSign = 0x08074B50
    fileHeaderLen = 30
    dirHeaderLen = 46
    dirEndLen = 22
    dataDescriptorLen = 16 
    dataDescriptor64Len = 24
    dir64LocLen = 20
    dir64EndLen = 56

    creatorFAT = 0 
    creatorUnix = 3 
    creatorNTFS = 11 
    creatorVFAT = 14 
    creatorMacOSX = 19 

    zipVersion20 = 20 // 2.0
    zipVersion45 = 45 // 4.5 (zip64 archives)

    uint16max = (1 << 16) - 1
    uint32max = (1 << 32) - 1 

    zip64ExtraID = 0x0001 
    ntfsExtraID = 0x000A 
    unixExtraID = 0x000D 
    extTimeExtraID = 0x5455 
    infoZipUnixExtraID = 0x5855
)

type FileHeader struct {
    Name string 
    Comment string 
    NonUTF8 bool 
    CreatorVersion uint16 
    ReaderVersion uint16 
    Flags uint16 
    Method uint16 
    ModTime uint16 
    ModDate uint16 
    CRC32 uint32 
    CompressedSize64 uint64 
    UncompressedSize64 uint64 
    Extra []byte 
    ExternalAttrs uint32
}

type Writer struct {
    cw *countWriter
    dir []*header
    last *fileWriter
    closed bool 
    compressors map[uint16]Compressor 
    comment string
}

type Compressor func(w io.Writer) (io.WriteCloser, error)

type header struct {
    *FileHeader
    offset uint64 
    raw bool
}

func NewWriter(w io.Writer) *Writer {
    return &Writer{cw: &countWriter{w: bufio.NewWriter(w)}}
}

type fileWriter struct {
    *header 
    zw io.Writer
    rawCount *countWriter
    comp io.WriteCloser
    compCount *countWriter
    crc32 hash.Hash32
    closed bool
}

type countWriter struct {
    w io.Writer
    count int64
}

func (cw *countWriter)Write(p []byte) (int, error) {
    n, err := cw.w.Write(p)
    cw.count += int64(n)
    return n, err
}
