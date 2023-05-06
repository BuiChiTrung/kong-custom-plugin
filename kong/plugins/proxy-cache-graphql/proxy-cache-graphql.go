package main

import (
	"fmt"
	"github.com/BuiChiTrung/kong-custom-plugin/kong/logger"
	"github.com/Kong/go-pdk"
	"github.com/Kong/go-pdk/server"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"
	"time"
)

var gConf Config
var gSvc *Service
var isHealthCheckGrOn bool
var gRedisHealthCheckIntervalSecond uint

func New() interface{} {
	plugin := Plugin{}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:5432/%s", os.Getenv(EnvKongPgUser), os.Getenv(EnvKongPgPassword), os.Getenv(EnvKongPgHost), os.Getenv(EnvKongPgDatabase))
	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	err := db.Where("name = ?", PluginName).Find(&plugin).Error

	logger.NewDefaultZapLogger(int(plugin.Config.LogFileSizeMaxMB), int(plugin.Config.LogAgeMaxDays))
	logger.Info("Restart plugin", "plugin", plugin, "err", err)

	gSvc = NewService()
	gConf = Config{}

	gRedisHealthCheckIntervalSecond = plugin.Config.RedisHealthCheckIntervalSecond
	if !isHealthCheckGrOn && gRedisHealthCheckIntervalSecond > 0 {
		isHealthCheckGrOn = true
		go HealthCheckRedis()
	}

	return &gConf
}

func HealthCheckRedis() {
	logger.Info("Start health check job")

	for {
		if gRedisHealthCheckIntervalSecond == 0 {
			break
		}
		gSvc.HealthCheckRedis()
		time.Sleep(time.Second * time.Duration(gRedisHealthCheckIntervalSecond))
	}

	isHealthCheckGrOn = false
	logger.Info("Stop health check job")
}

func (c Config) Access(kong *pdk.PDK) {
	defer func() {
		message := recover()
		if message != nil {
			logger.Errorf("Access: %v %s", message, string(debug.Stack()))
		}
	}()

	// TODO: trung.bc - TD
	_ = kong.ServiceRequest.ClearHeader(HeaderAcceptEncoding)

	cacheKey, shouldCached, err := c.GenerateCacheKey(kong)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	if !shouldCached {
		_ = kong.Response.SetHeader(HeaderXCacheStatus, string(Bypass))
		return
	}

	// Test log file size
	//for i := 0; i < 2; i++ {
	//	logger.Debug("Test log file size", "cacheKey", cacheKey, "shouldCached", shouldCached, "err", err)
	//}

	if err := kong.Ctx.SetShared(CacheKey, cacheKey); err != nil {
		logger.Errorf("err set shared context: %v", err)
		return
	}

	cacheVal, err := gSvc.GetCacheKey(cacheKey)
	if err != nil {
		_ = kong.Response.SetHeader(HeaderXCacheStatus, string(Miss))
		if err == redis.Nil {
			_ = kong.Response.SetHeader(HeaderXCacheKey, cacheKey)
			return
		}
		logger.Errorf("error get redis key: %v", err)
	} else {
		_ = kong.Ctx.SetShared(ResponseAlreadyCached, true)

		kong.Response.Exit(200, cacheVal, map[string][]string{
			HeaderContentType:  {"application/json"},
			HeaderXCacheKey:    {cacheKey},
			HeaderXCacheStatus: {string(Hit)},
			//HeaderAccessControlAllowOrigin: {"*"},
		})
	}
}

func (c Config) GenerateCacheKey(kong *pdk.PDK) (cacheKey string, shouldCached bool, err error) {
	requestBody, err := kong.Request.GetRawBody()
	if err != nil {
		return "", false, fmt.Errorf("err get request body: %w", err)
	}

	var requestHeader string
	for _, header := range c.Headers {
		headerContent, _ := kong.Request.GetHeader(header)
		requestHeader += headerContent
	}

	requestPath, err := kong.Request.GetPath()
	if err != nil {
		return "", false, fmt.Errorf("err GenerateCacheKey get req path")
	}

	cacheKey, shouldCached, err = gSvc.GenerateCacheKey(string(requestBody), requestHeader, requestPath)
	if err != nil {
		return "", false, err
	}

	return cacheKey, shouldCached, nil
}

// Response automatically enables the buffered proxy mode.
func (c Config) Response(kong *pdk.PDK) {
}

func (c Config) Log(kong *pdk.PDK) {
	defer func() {
		message := recover()
		if message != nil {
			logger.Errorf("Log: %v %s", message, string(debug.Stack()))
		}
	}()

	isCache, err := kong.Ctx.GetSharedAny(ResponseAlreadyCached)
	if v, ok := isCache.(bool); ok && v {
		return
	}

	responseBody, err := kong.ServiceResponse.GetRawBody()
	if err != nil {
		logger.Errorf("error get service response: %v", err)
		return
	}

	cacheKey, err := kong.Ctx.GetSharedString(CacheKey)
	if err != nil {
		logger.Errorf("err get shared context: %v", err.Error())
	}
	if cacheKey == "" {
		return
	}

	c.InsertCacheKey(kong, cacheKey, responseBody)
}

func (c Config) InsertCacheKey(kong *pdk.PDK, cacheKey string, cacheValue string) {
	logger.Infof("Insert Cache-key: %s", cacheKey)

	statusCode, _ := kong.ServiceResponse.GetStatus()
	if statusCode >= http.StatusInternalServerError {
		return
	}

	if statusCode >= http.StatusBadRequest && c.ErrTTLSeconds > 0 {
		if err := gSvc.InsertCacheKey(cacheKey, cacheValue, int64(c.ErrTTLSeconds)*int64(NanoSecond)); err != nil {
			logger.Errorf("err set redis key: %v", err)
		}
		return
	}

	var ttlSeconds uint

	ttlHeaderStr, _ := kong.Request.GetHeader(TTLSeconds)
	ttlHeader, err := strconv.Atoi(ttlHeaderStr)

	if ttlHeaderNotProvideOrInvalid := err != nil; ttlHeaderNotProvideOrInvalid {
		ttlSeconds = c.TTLSeconds
	} else if clientDonWantToCache := ttlHeader < 0; clientDonWantToCache {
		return
	} else {
		ttlSeconds = uint(ttlHeader)
	}

	if err := gSvc.InsertCacheKey(cacheKey, cacheValue, int64(ttlSeconds)*int64(NanoSecond)); err != nil {
		logger.Errorf("error set redis key: %v", err)
	}
}

func main() {
	Version := "1.1"
	Priority := 1
	_ = server.StartServer(New, Version, Priority)
}
