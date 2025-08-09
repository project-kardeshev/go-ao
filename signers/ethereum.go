package signers

import (
	goar "github.com/everFinance/goar"
	goarTypes "github.com/everFinance/goar/types"
	goether "github.com/everFinance/goether"
)

type EthereumSigner struct {
	wallet string
	signer goether.Signer
	itemSigner goar.ItemSigner
	Address string
	PublicKey string
}

func (ethereumSigner *EthereumSigner) CreateAndSignDataItem(
	data []byte,
	target string,
	anchor string,
	tags []goarTypes.Tag,
	) (goarTypes.BundleItem, error) {
		return ethereumSigner.itemSigner.CreateAndSignItem(data, target, anchor, tags)
}

func NewEthereumSigner(wallet string) (DataItemSigner, error) {

	signer, signerErr := goether.NewSigner(wallet)
	if signerErr != nil {
		return nil, signerErr
	}

	itemSigner, itemSignerErr := goar.NewItemSigner(signer)
	if itemSignerErr != nil {
		return nil, itemSignerErr
	}

	return &EthereumSigner{
		wallet: wallet,
		signer: *signer,
		itemSigner: *itemSigner,
		Address: signer.Address.String(),
		PublicKey: signer.GetPublicKeyHex(),
	}, nil
}