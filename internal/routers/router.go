package routers

import (
	"net/http"
	"time"

	_ "github.com/Acyclonepl/Blog-basedon-gin/docs" // swag init 生成的文档
	"github.com/Acyclonepl/Blog-basedon-gin/global"
	"github.com/Acyclonepl/Blog-basedon-gin/internal/middleware"
	"github.com/Acyclonepl/Blog-basedon-gin/internal/routers/api"
	v1 "github.com/Acyclonepl/Blog-basedon-gin/internal/routers/api/v1"
	"github.com/Acyclonepl/Blog-basedon-gin/pkg/limiter"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // 静态文件包
	ginSwagger "github.com/swaggo/gin-swagger" // gin 中间件
)

var methodLimiters = limiter.NewMethodLimiter().AddBuckets(
	limiter.LimiterBucketRule{
		Key:          "/auth",
		FillInterval: time.Second,
		Capacity:     10,
		Quantum:      10,
	},
)

func NewRouter() *gin.Engine {
	r := gin.New()
	if global.ServerSetting.RunMode == "debug" {
		r.Use(gin.Logger())
		r.Use(gin.Recovery())
	} else {
		r.Use(middleware.AccessLog())
		r.Use(middleware.Recovery())
	}
	r.Use(middleware.Tracing())
	r.Use(middleware.RateLimiter(methodLimiters))
	r.Use(middleware.ContextTimeout(global.AppSetting.DefaultContextTimeout))
	r.Use(middleware.Translations())

	article := v1.NewArticle()
	tag := v1.NewTag()
	upload := api.NewUpload()
	r.GET("/debug/vars", api.Expvar)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.POST("/upload/file", upload.UploadFile)
	r.POST("/auth", api.GetAuth)
	r.StaticFS("/static", http.Dir(global.AppSetting.UploadSavePath))
	apiv1 := r.Group("/api/v1")
	apiv1.Use(middleware.JWT())

	{
		apiv1.POST("/tags", tag.Create)
		apiv1.DELETE("/tags/:id", tag.Delete)
		apiv1.PUT("/tags/:id", tag.Update)
		apiv1.PATCH("/tags/:id/state", tag.Update)
		apiv1.GET("/tags", tag.List)

		apiv1.POST("/articles", article.Create)
		apiv1.DELETE("/articles/:id", article.Delete)
		apiv1.PUT("/articles/:id", article.Update)
		apiv1.PATCH("/articles/:id", article.Update)
		apiv1.GET("/articles/:id", article.Get)
		apiv1.GET("/articles", article.List)
	}
	return r
}
