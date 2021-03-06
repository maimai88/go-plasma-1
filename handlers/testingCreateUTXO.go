package handlers

import (
	"encoding/json"
	"net/http"

	fdb "github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/bankex/go-plasma/foundationdb"
	"github.com/bankex/go-plasma/types"
	common "github.com/ethereum/go-ethereum/common"
	"github.com/valyala/fasthttp"
)

type createUTXOrequest struct {
	For               string `json:"for"`
	BlockNumber       int    `json:"blockNumber"`
	TransactionNumber int    `json:"transactionNumber"`
	OutputNumber      int    `json:"outputNumber"`
	Value             string `json:"value"`
}

type createUTXOResponse struct {
	Error bool `json:"error"`
}

type CreateUTXOHandler struct {
	db          *fdb.Database
	utxoCreator *foundationdb.TestUTXOcreator
}

func NewCreateUTXOHandler(db *fdb.Database) *CreateUTXOHandler {
	creator := foundationdb.NewTestUTXOcreator(db)
	handler := &CreateUTXOHandler{db, creator}
	return handler
}

func (h *CreateUTXOHandler) Handle(w http.ResponseWriter, r *http.Request) {
	var requestJSON createUTXOrequest
	err := json.NewDecoder(r.Body).Decode(&requestJSON)
	if err != nil {
		writeErrorResponse(w)
		return
	}

	forBytes := common.FromHex(requestJSON.For)
	address := common.Address{}
	copy(address[:], forBytes)
	bigint := types.NewBigInt(0)
	bigint.SetString(requestJSON.Value, 10)
	blockNumber := uint32(requestJSON.BlockNumber)
	transactionNumber := uint32(requestJSON.TransactionNumber)
	outputNumber := uint8(requestJSON.OutputNumber)
	err = h.utxoCreator.InsertUTXO(address, blockNumber, transactionNumber, outputNumber, bigint)
	if err != nil {
		writeErrorResponse(w)
		return
	}
	writeSuccessResponse(w)
	return
}

func (h *CreateUTXOHandler) HandlerFunc(ctx *fasthttp.RequestCtx) {
	var requestJSON createUTXOrequest
	err := json.Unmarshal(ctx.PostBody(), &requestJSON)
	if err != nil {
		writeFasthttpErrorResponse(ctx)
		return
	}

	forBytes := common.FromHex(requestJSON.For)
	address := common.Address{}
	copy(address[:], forBytes)
	bigint := types.NewBigInt(0)
	bigint.SetString(requestJSON.Value, 10)
	blockNumber := uint32(requestJSON.BlockNumber)
	transactionNumber := uint32(requestJSON.TransactionNumber)
	outputNumber := uint8(requestJSON.OutputNumber)
	err = h.utxoCreator.InsertUTXO(address, blockNumber, transactionNumber, outputNumber, bigint)
	if err != nil {
		writeFasthttpErrorResponse(ctx)
		return
	}
	writeFasthttpSuccessResponse(ctx)
	return
}
