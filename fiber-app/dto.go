package main

type GetCacheKeyResponse struct {
	FormalResponse
	Data *GetCacheKeyResponseData `json:"data"`
}

type GetCacheKeyResponseData struct {
	Value string `json:"value"`
}

type DelCacheKeyResponse struct {
	FormalResponse
}

type FormalResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type UpsertCacheKeyRequest struct {
	CacheKey string `json:"cacheKey"`
	Value    string `json:"value"`
}

type UpsertCacheKeyResponse struct {
	FormalResponse
}
