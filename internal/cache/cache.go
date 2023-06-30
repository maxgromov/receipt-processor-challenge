package cache

import "receipt-processor-challenge/internal/entity"

var ReceiptCache = make(map[string]*entity.Receipt, 0)
