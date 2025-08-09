package ao

import (
	"strconv"
	"time"
	goarTypes "github.com/everFinance/goar/types"
	resty "github.com/go-resty/resty/v2"
	signers "github.com/project-kardeshev/go-ao/signers"
)

func ApplyProtocolTags(tags []goarTypes.Tag) []goarTypes.Tag {
    return append(tags, goarTypes.Tag{
        Name: "Data-Protocol",
        Value: "ao",
    }, goarTypes.Tag{
        Name: "Variant",
        Value: "ao.TN.1",
    })
}

func ApplySpawnProtocolTags(tags []goarTypes.Tag, module string, authority string, scheduler string) []goarTypes.Tag {
    return append(ApplyProtocolTags(tags), goarTypes.Tag{
        Name: "Type",
        Value: "Process",
    }, goarTypes.Tag{
        Name: "Module",
        Value: module,
    }, goarTypes.Tag{
        Name: "Authority",
        Value: authority,
    }, goarTypes.Tag{
        Name: "Scheduler",
        Value: scheduler,
    }, goarTypes.Tag{
        Name: "Timestamp",
        Value: strconv.FormatInt(time.Now().UnixMilli(), 10),
    })
}

func ApplyMessageProtocolTags(tags []goarTypes.Tag) []goarTypes.Tag {
    return append(ApplyProtocolTags(tags), goarTypes.Tag{
        Name: "Type",
        Value: "Message",
    })
}


type ProcessClient struct {
	ProcessId string
	CuUrl string
    cuClient *resty.Client
	MuUrl string
	muClient *resty.Client
	dataItemSigner signers.DataItemSigner
}

func (processClient *ProcessClient) Read(input DryRunInput) (*Result, error) {
	response, err := processClient.cuClient.R().
		SetBody(input).
		Post("/dry-run?process-id=" + processClient.ProcessId)

    result := response.Result().(*Result)

	return result, err
}

func (processClient *ProcessClient) Write(input WriteInput) (id string, result *Result, err error) {
    dataItem, err := processClient.dataItemSigner.CreateAndSignDataItem(
        []byte(input.Data),
        input.Process,
        input.Anchor,
        ApplyMessageProtocolTags(input.Tags),
    )
    if err != nil {
        return "", nil, err
    }

    messageResponse, err := processClient.muClient.R().
        SetBody(dataItem.ItemBinary).
        SetHeader("Content-Type", "application/octet-stream").
        SetHeader("Accept", "application/json").
        Post("/")

    if err != nil {
        return "stub-id", nil, err
    }

    messageId := messageResponse.Result().(map[string]any)["id"].(string)

    messageResult, err := processClient.cuClient.R().
        Get("/result/" + messageId + "?process-id=" + processClient.ProcessId)

    return messageId, messageResult.Result().(*Result), err
}

func (processClient *ProcessClient) Spawn(input SpawnInput) (id string, result *Result, err error) {
    dataItem, err := processClient.dataItemSigner.CreateAndSignDataItem(
        []byte(input.Data.(string)),
        input.Target,
        input.Anchor,
        ApplySpawnProtocolTags(input.Tags, input.Module, input.Authority, input.Scheduler),
    )

    if err != nil {
        return "stub-id", nil, err
    }

    messageResponse, err := processClient.muClient.R().
        SetBody(dataItem.ItemBinary).
        SetHeader("Content-Type", "application/octet-stream").
        SetHeader("Accept", "application/json").
        Post("/")

    if err != nil {
        return "stub-id", nil, err
    }

    messageId := messageResponse.Result().(map[string]any)["id"].(string)

    messageResult, err := processClient.cuClient.R().
        Get("/result/" + messageId + "?process-id=" + processClient.ProcessId)

    return messageId, messageResult.Result().(*Result), err
}




func NewProcessClient(processId string, cuUrl string, muUrl string, dataItemSigner signers.DataItemSigner) AOClient {
	return &ProcessClient{
		ProcessId: processId,
		CuUrl: cuUrl,
		cuClient: resty.New().SetBaseURL(cuUrl).SetRedirectPolicy(resty.FlexibleRedirectPolicy(10)),
		MuUrl: muUrl,
		muClient: resty.New().SetBaseURL(muUrl).SetRedirectPolicy(resty.FlexibleRedirectPolicy(10)),
		dataItemSigner: dataItemSigner,
	}
}

