# GoFID
FID Library for Go

#### Generate
Generate(System Indicator, Vendor, Type, SubType, Location, Secret)
```
id, err := gofid.Generate(gofid.IndicatorEntity, "FOR", "TE", "ST", "", "secr3t")
```

#### Verify
Verify(FID , Secret)
```
result, err := gofid.Verify("IPIH7MI2=-EABCCDEF-MISCR-V669VFQ", "secr3t")
```

#### Describe
TBC