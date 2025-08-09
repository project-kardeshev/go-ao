package signers

import (
	goar "github.com/everFinance/goar"
	goarTypes "github.com/everFinance/goar/types"
)

type ArweaveSigner struct {
	wallet JWKByteInterface
	signer goar.Signer
	itemSigner goar.ItemSigner
	Address string
	PublicKey string
}

func (arweaveSigner *ArweaveSigner) CreateAndSignDataItem(
	data []byte,
	target string,
	anchor string,
	tags []goarTypes.Tag,
	) (goarTypes.BundleItem, error) {
		return arweaveSigner.itemSigner.CreateAndSignItem(data, target, anchor, tags)
}


func NewArweaveSigner(wallet JWKByteInterface) (DataItemSigner, error) {

	signer, signerErr := goar.NewSigner(wallet)
	if signerErr != nil {
		return nil, signerErr
	}

	itemSigner, itemSignerErr := goar.NewItemSigner(signer)
	if itemSignerErr != nil {
		return nil, itemSignerErr
	}

	return &ArweaveSigner{
		wallet: wallet,
		signer: *signer,
		itemSigner: *itemSigner,
		Address: signer.Address,
		PublicKey: signer.PubKey.N.String(),
	}, nil
}



