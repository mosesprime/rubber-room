// Use at your own risk.

package gzip

import (
    "compress/flate"
    "io"
    "time"
)

// Writer compression level.
const (
    NoCompression = flate.NoCompression
    BestSpeed = flate.BestSpeed
    BestCompression = flate.BestCompression
    DefaultCompression = flate.DefaultCompression
    HuffmanOnly = flate.HuffmanOnly
)

const (
    gzipID1 = 0x1F
    gzipID2 = 0x8B
    gzipDeflate = 8 
    flagText = 1 << 0
    flagHdrCrc = 1 << 1 
    flagExtra = 1 << 2 
    flagName = 1 << 3 
    flagComment = 1 << 4
)

type Writer struct {
    Header 
    w io.Writer
    level int
    wroteHeader bool
    compressor *flate.Writer
    digest uint32
    size uint32
    closed bool
    buf [10]byte
    err error
}

type Header struct {
    Comment string
    Extra []byte
    ModTime time.Time
    Name string
    OS byte
}

func NewWriter(w io.Writer) *Writer {
    z := new(Writer)
    z.init(w)
    return z
}

func (z *Writer) init(w io.Writer) {
    compressor := z.compressor
    if compressor != nil {
        compressor.Reset(w)
    }
    *z = Writer{
        Header: Header{
            OS: 255, // unknown OS
        },
        w: w,
        level: BestCompression,
        compressor: compressor,
    }
}
