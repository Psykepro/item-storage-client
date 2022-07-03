package request

import (
	"os"

	domain "github.com/Psykepro/item-storage-client/_domain"
	pb "github.com/Psykepro/item-storage-protobuf/generated/item"
	"github.com/gocarina/gocsv"
	"github.com/google/uuid"
)

const (
	DataPath = "./data/data.csv"
)

type Loader struct {
	logger domain.Logger
}

func NewLoader(logger domain.Logger) *Loader {
	return &Loader{logger: logger}
}

type requestDTO struct {
	Command string  `csv:"command"`
	Uuid    *string `csv:"uuid"`
	Data    *string `csv:"data"`
}

func (l *Loader) LoadRequestsFromCsv(filePath string) []*pb.ItemRequest {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var requests []*requestDTO

	if err := gocsv.UnmarshalFile(file, &requests); err != nil {
		panic(err)
	}

	result := make([]*pb.ItemRequest, 0)
	for _, req := range requests {
		switch req.Command {
		case pb.Command_CREATE.String():
			result = l.loadCreateRequestIfValid(req, result)
			break
		case pb.Command_GET.String():
			result = l.loadRequestWithUuidIfValid(req, pb.Command_GET, result)
			break
		case pb.Command_DELETE.String():
			result = l.loadRequestWithUuidIfValid(req, pb.Command_DELETE, result)
			break
		case pb.Command_LIST.String():
			result = append(result, &pb.ItemRequest{
				Command: pb.Command_LIST,
			})
			break
		default:
			l.logger.Fatalf("Unsupported command for request - [%s].", req.Command)
		}
	}
	return result
}

func (l *Loader) loadRequestWithUuidIfValid(request *requestDTO, command pb.Command, result []*pb.ItemRequest) []*pb.ItemRequest {
	l.validateItemFields(*request.Uuid, nil)
	result = append(result, &pb.ItemRequest{
		Command: command,
		Item:    &pb.Item{Uuid: *request.Uuid},
	})

	return result
}

func (l *Loader) loadCreateRequestIfValid(request *requestDTO, result []*pb.ItemRequest) []*pb.ItemRequest {
	l.validateItemFields(*request.Uuid, request.Data)
	item := pb.Item{
		Uuid: *request.Uuid,
		Data: *request.Data,
	}
	result = append(result, &pb.ItemRequest{
		Command: pb.Command_CREATE,
		Item:    &item,
	})
	return result
}

// validateItemFields validates the field values for and item.
// 'itemData' is pointer (optional) because it is not expected for each request
func (l *Loader) validateItemFields(itemUuid string, itemData *string) {
	_ = uuid.MustParse(itemUuid)
	if itemData != nil && *itemData == "" {
		l.logger.Fatalf("Empty item data.")
	}
}
