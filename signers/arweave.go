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

func (arweaveSigner *ArweaveSigner) GetAddress() string {
	return arweaveSigner.Address
}

/*
	NewArweaveSigner creates a new ArweaveSigner.
	
	wallet: The JWK to use for signing.
	
	Returns:
		A new ArweaveSigner.
		An error if the signer or item signer creation fails.
	Example:
		// read wallet from file
		wallet, err := os.ReadFile("wallet.json")
		if err != nil {
			log.Fatal(err)
		}
		// create signer
		signer, err := NewArweaveSigner(wallet)
		if err != nil {
			log.Fatal(err)
		}
*/
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



