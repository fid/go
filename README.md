# GoFID
FID Library for Go

#### Generate
Generate(System Indicator, Vendor, Type, SubType, Location, Secret)
```go
id, err := gofid.Generate(gofid.IndicatorEntity, "FOR", "TE", "ST", "", "secr3t")
```

#### Verify
Verify(FID , Secret)
```go
result, err := gofid.Verify("IPIH7MI2=-EABCCDEF-MISCR-V669VFQ", "secr3t")
```

#### Describe
Describe(FID)
```go
result, err := gofid.Describe("IPIH7MI2=-EABCCDEF-MISCR-V669VFQ")
```

Describe returns a description structure
```go
type Description struct {
	Indicator    gofid.TypeIndicator
	VendorKey    string
	Type         string
	SubType      string
	Location     string
	TimeKey      string
	Time         time.Time
	RandomString string
}
```