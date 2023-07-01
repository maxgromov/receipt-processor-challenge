package handler

import (
	"encoding/json"
	"receipt-processor-challenge/internal/cache"
	"receipt-processor-challenge/internal/entity"
	"receipt-processor-challenge/internal/utilities"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

const (
	layoutDate = "2006-01-02"
	layoutTime = "15:04"
)

type Handler struct {
	log *logrus.Logger
}

func NewHandler(log *logrus.Logger) *Handler {
	return &Handler{log: log}
}

func (h *Handler) ProcessReceipt(ctx *fasthttp.RequestCtx) {
	rId := uuid.New().String()

	var (
		reqData              entity.ReceiptDeserialize
		recieptData          *entity.Receipt
		items                []*entity.Item
		parseDate, parseTime time.Time
		t                    float64
	)

	err := json.Unmarshal(ctx.Request.Body(), &reqData)
	if err != nil {
		h.log.Error("parse request body error", err.Error())
		ctx.SetBodyString("Incorrect request body")
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	// checking struct for null value, because it might be sent by user
	if utilities.StructIsEmpty(reqData) {
		h.log.Error("empty request data")
		ctx.SetBodyString("Empty request body")
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	// checking for null value, because it might be sent by user
	if reqData.PurchaseDate != "" {
		parseDate, err = time.Parse(layoutDate, reqData.PurchaseDate)
		if err != nil {
			h.log.Error("parse datetime error ", err.Error())
			ctx.SetBodyString("Date or Time format incorrect")
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			return
		}
	}
	// checking for null value, because it might be sent by user
	if reqData.PurchaseDate != "" {
		parseTime, err = time.Parse(layoutTime, reqData.PurchaseTime)
		if err != nil {
			h.log.Error("parse datetime error ", err.Error())
			ctx.SetBodyString("Date or Time format incorrect")
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			return
		}
	}

	for _, i := range reqData.Items {
		var pf float64
		// checking for null value, because it might be sent by user
		if i.Price != "" {
			pf, err = strconv.ParseFloat(i.Price, 64)
			if err != nil {
				h.log.Error("parse price error ", err.Error())
				ctx.SetBodyString("Price format incorrect")
				ctx.SetStatusCode(fasthttp.StatusBadRequest)
				return
			}
		}

		itemData := &entity.Item{
			ShortDescription: i.ShortDescription,
			Price:            pf,
		}
		items = append(items, itemData)
	}
	// checking for null value, because it might be sent by user
	if reqData.Total != "" {
		t, err = strconv.ParseFloat(reqData.Total, 64)
		if err != nil {
			h.log.Error("parse price error", err.Error())
			ctx.SetBodyString("Price format incorrect")
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			return
		}
	}

	// type conversion according to the field types of the structure
	recieptData = &entity.Receipt{
		Retailer:     reqData.Retailer,
		PurchaseDate: parseDate,
		PurchaseTime: parseTime,
		Total:        t,
		Items:        items,
	}

	response := entity.ReceiptResponse{Id: rId}
	_, isExist := cache.ReceiptCache[rId]

	if !isExist {
		cache.ReceiptCache[rId] = recieptData
	}
	js, err := json.Marshal(response)
	if err != nil {
		h.log.Error(err)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	ctx.SetContentType("application/json")
	ctx.SetBody(js)
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func (h *Handler) GetPoints(ctx *fasthttp.RequestCtx) {
	receiptId := ctx.UserValue("id").(string)

	receipt, isExist := cache.ReceiptCache[receiptId]
	if !isExist {
		h.log.Errorf("No receipt found for that id")
		ctx.SetBodyString("No receipt found for that id")
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	points, err := utilities.CountPoints(receipt)
	if err != nil {
		h.log.Error("processing receipt points error", err.Error())
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	response := entity.PointsResponse{Points: points}
	js, err := json.Marshal(response)
	if err != nil {
		h.log.Error("marshaling json error", err.Error())
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	ctx.SetContentType("application/json")
	ctx.SetBody(js)
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func (h *Handler) Health(ctx *fasthttp.RequestCtx) {
	js, err := json.Marshal("Hello")
	if err != nil {
		h.log.Error(err.Error())
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetContentType("application/json")
	ctx.SetBody(js)
	ctx.SetStatusCode(fasthttp.StatusOK)
}
