package signers

import goarTypes "github.com/everFinance/goar/types"

type DataItemSigner interface {
    CreateAndSignDataItem(data []byte, target string, anchor string, tags []goarTypes.Tag) (goarTypes.BundleItem, error)
    GetAddress() string
}

type JWKByteInterface = []byte