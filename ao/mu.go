package ao

import (
	"fmt"

	goarTypes "github.com/everFinance/goar/types"
	resty "github.com/go-resty/resty/v2"
)

type MuClient struct {
	client *resty.Client
	MuUrl string
}

type MuAPIs interface {
	PostAoMessage(dataItem goarTypes.BundleItem) (id string, err error)
}
/*
	PostAoMessage posts an AO message, assignment, or spawn to the MU.
	
	dataItem: The signed data item to post.
	
	Returns:
		The id of the message.
		An error if the post fails.
	Example:
		dataItem, err := dataItemSigner.CreateAndSignDataItem(
			[]byte(input.Data),
			input.Process,
			input.Anchor,
			ApplyMessageProtocolTags(input.Tags),
		)
		id, err := muClient.PostAoMessage(dataItem)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Message posted with id:", id)
*/
type MuResponse struct {
	ID string `json:"id"`
}

func (muClient *MuClient) PostAoMessage(dataItem goarTypes.BundleItem) (id string, err error) {
	var result MuResponse
	
	response, err := muClient.client.R().
		SetBody(dataItem.ItemBinary).
		SetHeader("Content-Type", "application/octet-stream").
		SetHeader("Accept", "application/json").
		SetResult(&result).
		Post("/")

	if err != nil {
		return "", err
	}

	// Accept both 200 (OK) and 202 (Accepted) as success
	if response.StatusCode() != 200 && response.StatusCode() != 202 {
		return "", fmt.Errorf("HTTP %d: %s", response.StatusCode(), response.String())
	}

	return result.ID, nil
}

func NewMuClient(muUrl string) MuAPIs {
	return &MuClient{
		client: resty.New().
			SetBaseURL(muUrl).
			SetRedirectPolicy(resty.FlexibleRedirectPolicy(10)).
			SetError(&Result{}),
		MuUrl: muUrl,
	}
}