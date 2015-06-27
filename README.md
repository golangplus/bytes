# bytes
Plus to the standard `bytes` package.

## Featured
```go
// ByteSlice is a wrapper type for []byte.
// Its pointer form, *ByteSlice, implements io.Reader, io.Writer, io.ByteReader,
// io.ByteWriter, io.Closer, io.ReaderFrom, io.WriterTo and io.RuneReader
// interfaces.
//
// When reading from a constant small byte slice and no need for seeking, *ByteSlice is a
// better alternative then bytes.Buffer, since it needs less extra resource.
type ByteSlice []byte
```

## LICENSE
BSD license