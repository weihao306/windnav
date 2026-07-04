package httpapi

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"windnav/backend/internal/config"
)

type Server struct {
	cfg config.Config
	db  *gorm.DB
}

func NewRouter(cfg config.Config, db *gorm.DB) *gin.Engine {
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	s := &Server{cfg: cfg, db: db}
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.GET("/healthz", func(c *gin.Context) {
		ok(c, gin.H{"status": "ok", "driver": cfg.DBDriver}, nil)
	})

	api := r.Group("/api")
	public := api.Group("/public")
	public.GET("/summary", s.publicSummary)
	public.GET("/categories", s.publicCategories)
	public.GET("/search-engines", s.publicSearchEngines)
	public.GET("/sites", s.publicSites)
	public.POST("/sites/:id/click", s.recordClick)

	auth := api.Group("/auth")
	auth.POST("/login", s.login)
	auth.POST("/logout", s.authRequired(), s.logout)
	auth.GET("/me", s.authRequired(), s.me)

	admin := api.Group("/admin", s.authRequired())
	admin.GET("/categories", s.listCategories)
	admin.POST("/categories", s.createCategory)
	admin.PUT("/categories/:id", s.updateCategory)
	admin.DELETE("/categories/:id", s.deleteCategory)
	admin.GET("/sites", s.listSites)
	admin.POST("/sites", s.createSite)
	admin.PUT("/sites/:id", s.updateSite)
	admin.DELETE("/sites/:id", s.deleteSite)
	admin.GET("/tags", s.listTags)
	admin.POST("/tags", s.createTag)
	admin.PUT("/tags/:id", s.updateTag)
	admin.DELETE("/tags/:id", s.deleteTag)
	admin.GET("/search-engines", s.listSearchEngines)
	admin.POST("/search-engines", s.createSearchEngine)
	admin.PUT("/search-engines/:id", s.updateSearchEngine)
	admin.DELETE("/search-engines/:id", s.deleteSearchEngine)
	admin.GET("/settings", s.listSettings)
	admin.PUT("/settings", s.updateSettings)

	if cfg.StaticDir != "" {
		mountStatic(r, cfg.StaticDir)
	}

	return r
}

func mountStatic(r *gin.Engine, staticDir string) {
	r.Static("/assets", filepath.Join(staticDir, "assets"))
	r.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			fail(c, http.StatusNotFound, "NOT_FOUND", "接口不存在", nil)
			return
		}
		c.File(filepath.Join(staticDir, "index.html"))
	})
}
