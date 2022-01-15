package user

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"

	"github.com/TheMickeyMike/grpc-rest-bench/pb"
	"go.uber.org/zap"
)

// basepath is the root directory of this package.
var basepath string

func init() {
	_, currentFile, _, _ := runtime.Caller(0)
	basepath = filepath.Dir(currentFile)
}

type Db struct {
	Data      map[string]*pb.UserAccount
	DataArray []*pb.UserAccount
}

func NewDb(logger *zap.Logger) *Db {
	file, err := os.Open(filepath.Join(basepath, "data/users.json"))
	if err != nil {
		logger.Fatal("can't read file", zap.Error(err))
	}
	defer file.Close()

	var decodedData []*pb.UserAccount
	if err := json.NewDecoder(file).Decode(&decodedData); err != nil {
		logger.Fatal("can't load data from json", zap.Error(err))
	}
	db := Db{
		Data:      make(map[string]*pb.UserAccount, len(decodedData)),
		DataArray: decodedData,
	}
	for _, r := range decodedData {
		db.Data[r.Id] = r
	}
	logger.Info("db loaded", zap.Int("records", len(db.Data)))
	return &db
}
