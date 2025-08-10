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

func (ethereumSigner *EthereumSigner) GetAddress() string {
	return ethereumSigner.Address
}

/*
	NewEthereumSigner creates a new EthereumSigner.
	
	wallet: The wallet to use for signing.
	
	Returns:
		A new EthereumSigner.
		An error if the signer or item signer creation fails.
	Example:
		privateKey := "0x1234567890123456789012345678901234567890"
		signer, err := NewEthereumSigner(privateKey)
		if err != nil {
			log.Fatal(err)
		}
	
*/
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