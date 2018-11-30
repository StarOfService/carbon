# Flapper [![GoDoc](https://img.shields.io/badge/go-documentation-blue.svg)](https://godoc.org/github.com/StarOfService/flapper)

Package flapper allows to serialize complex custom structures to a flat map of strings and deserialize it back. This library can be useful when you define metadata of an object as a custom structure and want to store it at the object tags or labels.

## Install

```bash
go get github.com/fatih/structs
```

## Usage and Examples

Interface of the library is simple and matches to the approach used at `json` and `yaml` packages. We have two main functions: `Marshal` and `Unmarshal`:
```
require (
  "fmt"
  "github.com/starofservice/flapper"
)


func main() {
  type TStruct2 struct {
    DA string
    DB float32
  }

  type TStruct1 struct {
    A string
    B int
    C bool
    D TStruct2
    e string
  }

  var obj = TStruct1{
    A: "a-value",
    B: 2,
    C: true,
    D: TStruct2{
      DA: "d-value",
      DB: 3.14,
    },
    e: "non-public fields should be skipped",
  }

  // Serialize the object
  serial, err := flapper.Marshal(obj)
  if err != nil {
    fmt.Println("error:", err)
  }

  fmt.Printf("%+v", serial) // map[D.DB:3.14E+00 A:a-value B:2 C:true D.DA:d-value]

  var deserial TStruct1
  // Deserialize the object
  err = flapper.Unmarshal(serial, &deserial)
  if err != nil {
    panic(err)
  }

```

### Configuration

Current library provides possibility to set custom configuration. Currently custom prefix and delimiter are supported.
In order to define custom configuration, you have to create a new object of Flapper type and use its Marshal/Unmarshal methods intead of simple functions:

```
require (
  "fmt"
  "github.com/starofservice/flapper"
)


func main() {
  type TStruct2 struct {
    DA string
    DB float32
  }

  type TStruct1 struct {
    A string
    B int
    C bool
    D TStruct2
    e string
  }

  var obj = TStruct1{
    A: "a-value",
    B: 2,
    C: true,
    D: TStruct2{
      DA: "d-value",
      DB: 3.14,
    },
    e: "non-public fields should be skipped",
  }

  // Define custom prefix and delimiter
  fc, err := flapper.New("test", ":")
  if err != nil {
    fmt.Println("error:", err)
  }

  // Serialize the object
  serial, err := fc.Marshal(obj)
  if err != nil {
    fmt.Println("error:", err)
  }

  fmt.Printf("%+v", serial) // map[test:B:2 test:C:true test:D:DA:d-value test:D:DB:3.14E+00 test:A:a-value]

  var deserial TStruct1
  // Deserialize the object
  err = fc.Unmarshal(serial, &deserial)
  if err != nil {
    panic(err)
  }

```

### Supported field types

- Array
- Bool
- Float*
- Int*
- Slice
- String
- Struct (including inbounded structs)
- Uint*

# Limitations

Current version of Flapper has a number of limitations:
- Map type is not supported
- Channel type is not supported
- Complex* types are not surpported
- Embedded types are not supported
- Function type is not supported
- Interface type is not supported
- Pointer, Uintptr and UnsafePointer types are not supported

If you need any of these types, feel free to propose PR :)

## License

The MIT License (MIT) - see LICENSE.md for more details