package ao

import (
	"errors"
	"strconv"
	"time"

	goarTypes "github.com/everFinance/goar/types"
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
	ProcessId *string
    cuClient CuAPIs
	muClient MuAPIs
	dataItemSigner signers.DataItemSigner
}

func (processClient *ProcessClient) Read(input DryRunInput) (*Result, error) {
    if processClient.ProcessId == nil {
        return nil, errors.New("processId is nil, spawn a process or set it")
    }
	return processClient.cuClient.DryRun(input, *processClient.ProcessId)
}

func (processClient *ProcessClient) Write(input WriteInput) (id string, result *Result, err error) {
    if processClient.ProcessId == nil {
        return "", nil, errors.New("processId is nil, spawn a process or set it")
    }
    
    // Generate random anchor if none provided
    var anchor string
    if input.Anchor != nil {
        anchor = *input.Anchor
    } else {
        randomAnchor, err := CreateRandomAnchor()
        if err != nil {
            return "", nil, err
        }
        anchor = randomAnchor
    }
    
    dataItem, err := processClient.dataItemSigner.CreateAndSignDataItem(
        []byte(Coalesce[string](input.Data, "")),
        *processClient.ProcessId,
        anchor,
        ApplyMessageProtocolTags(Coalesce[[]goarTypes.Tag](input.Tags, []goarTypes.Tag{
            {
                Name: "SDK",
                Value: "github.com/project-kardeshev/go-ao",
            },
        })),
    )
    if err != nil {
        return "", nil, err
    }

    messageId, err := processClient.muClient.PostAoMessage(dataItem)

    if err != nil {
        return "", nil, err
    }

    messageResult, err := processClient.cuClient.Result(messageId, *processClient.ProcessId)

    return messageId, messageResult, err
}

func (processClient *ProcessClient) Spawn(input SpawnInput) (id string, result *Result, err error) {
    // Remove the & - these are already strings, not pointers
    data := Coalesce[string](input.Data, "")  // Type assert the interface{}
    target := Coalesce[string](input.Target, StubArweaveId)
    anchor := Coalesce[string](input.Anchor, StubArweaveId)

    dataItem, err := processClient.dataItemSigner.CreateAndSignDataItem(
        []byte(data),
        target,
        anchor,
        ApplySpawnProtocolTags(input.Tags, input.Module, input.Authority, input.Scheduler),
    )

    if err != nil {
        return "", nil, err
    }
    // In this case, the messageId is the spawned processId
    messageId, err := processClient.muClient.PostAoMessage(dataItem)

    if err != nil {
        return "", nil, err
    }

    messageResult, err := processClient.cuClient.Result(messageId, messageId)

    if err != nil {
        return "", nil, err
    }

    if processClient.ProcessId == nil {
        processClient.ProcessId = &messageId
    }

    return messageId, messageResult, err
}

func NewProcessClient(processId *string, cuUrl string, muUrl string, dataItemSigner signers.DataItemSigner) AOClient {
	return &ProcessClient{
		ProcessId: processId,
		cuClient: NewCuClient(cuUrl),
		muClient: NewMuClient(muUrl),
		dataItemSigner: dataItemSigner,
	}
}

