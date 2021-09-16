package generated

//go:generate go install "github.com/zhangsifeng92/geos/libraries/container/treemap"
//go:generate go install "github.com/zhangsifeng92/geos/libraries/container/treeset"
//go:generate gotemplate -outfmt "gen_%v" "github.com/zhangsifeng92/geos/libraries/container/treeset" AccountNameSet(common.AccountName,common.CompareName,false)
//go:generate gotemplate -outfmt "gen_%v" "github.com/zhangsifeng92/geos/libraries/container/treemap" AccountNameUint32Map(common.AccountName,uint32,common.CompareName,false)
//go:generate gotemplate -outfmt "gen_%v" "github.com/zhangsifeng92/geos/libraries/container/treemap" AccountNameUint64Map(common.AccountName,uint64,common.CompareName,false)
//go:generate gotemplate -outfmt "gen_%v" "github.com/zhangsifeng92/geos/libraries/container/treeset" AccountDeltaSet(common.AccountDelta,common.CompareAccountDelta,false)
//go:generate gotemplate -outfmt "gen_%v" "github.com/zhangsifeng92/geos/libraries/container/treeset" AccountDeltaSet(common.AccountDelta,common.CompareAccountDelta,false)
//go:generate gotemplate -outfmt "gen_%v" "github.com/zhangsifeng92/geos/libraries/container/treeset" NamePairSet(common.NamePair,common.CompareNamePair,false)

//go:generate gotemplate -outfmt "gen_%v" "github.com/zhangsifeng92/geos/libraries/container/treeset" PermissionLevelSet(common.PermissionLevel,common.ComparePermissionLevel,false)
//go:generate gotemplate -outfmt "gen_%v" "github.com/zhangsifeng92/geos/libraries/container/treeset" PublicKeySet(ecc.PublicKey,ecc.ComparePubKey,false)

//go:generate go build
