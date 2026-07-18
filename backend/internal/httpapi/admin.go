package httpapi

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"windnav/backend/internal/model"
)

type categoryInput struct {
	Name        string `json:"name" binding:"required"`
	Slug        string `json:"slug" binding:"required"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Color       string `json:"color"`
	SortOrder   int    `json:"sortOrder"`
	IsVisible   *bool  `json:"isVisible"`
}

type siteInput struct {
	CategoryID   uint   `json:"categoryId" binding:"required"`
	Title        string `json:"title" binding:"required"`
	URL          string `json:"url" binding:"required"`
	Description  string `json:"description"`
	IconURL      string `json:"iconUrl"`
	FallbackIcon string `json:"fallbackIcon"`
	SortOrder    int    `json:"sortOrder"`
	IsPinned     bool   `json:"isPinned"`
	IsVisible    *bool  `json:"isVisible"`
	TagIDs       []uint `json:"tagIds"`
}

type tagInput struct {
	Name  string `json:"name" binding:"required"`
	Slug  string `json:"slug" binding:"required"`
	Color string `json:"color"`
}

type searchEngineInput struct {
	Name      string `json:"name" binding:"required"`
	Slug      string `json:"slug" binding:"required"`
	SearchURL string `json:"searchUrl" binding:"required"`
	Icon      string `json:"icon"`
	SortOrder int    `json:"sortOrder"`
	IsDefault bool   `json:"isDefault"`
	IsVisible *bool  `json:"isVisible"`
}

type settingsInput struct {
	Settings map[string]string `json:"settings" binding:"required"`
}

func (s *Server) listCategories(c *gin.Context) {
	var categories []model.Category
	if err := s.db.Order("sort_order asc, id asc").Find(&categories).Error; err != nil {
		fail(c, http.StatusInternalServerError, "SERVER_ERROR", "读取分类失败", nil)
		return
	}
	ok(c, categories, nil)
}

func (s *Server) createCategory(c *gin.Context) {
	var req categoryInput
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, "VALIDATION_ERROR", "请求参数无效", err.Error())
		return
	}
	visible := true
	if req.IsVisible != nil {
		visible = *req.IsVisible
	}
	category := model.Category{Name: req.Name, Slug: req.Slug, Description: req.Description, Icon: req.Icon, Color: req.Color, SortOrder: req.SortOrder, IsVisible: visible}
	if err := s.db.Create(&category).Error; err != nil {
		fail(c, http.StatusBadRequest, "CREATE_FAILED", "创建分类失败", err.Error())
		return
	}
	created(c, category)
}

func (s *Server) updateCategory(c *gin.Context) {
	id, okID := parseID(c)
	if !okID {
		return
	}
	var req categoryInput
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, "VALIDATION_ERROR", "请求参数无效", err.Error())
		return
	}
	var category model.Category
	if !findByID(c, s.db, id, &category, "分类不存在") {
		return
	}
	category.Name = req.Name
	category.Slug = req.Slug
	category.Description = req.Description
	category.Icon = req.Icon
	category.Color = req.Color
	category.SortOrder = req.SortOrder
	if req.IsVisible != nil {
		category.IsVisible = *req.IsVisible
	}
	if err := s.db.Save(&category).Error; err != nil {
		fail(c, http.StatusBadRequest, "UPDATE_FAILED", "更新分类失败", err.Error())
		return
	}
	ok(c, category, nil)
}

func (s *Server) deleteCategory(c *gin.Context) {
	id, okID := parseID(c)
	if !okID {
		return
	}
	if err := s.db.Delete(&model.Category{}, id).Error; err != nil {
		fail(c, http.StatusBadRequest, "DELETE_FAILED", "删除分类失败", err.Error())
		return
	}
	ok(c, gin.H{"ok": true}, nil)
}

func (s *Server) listSites(c *gin.Context) {
	page, pageSize := pagination(c)
	query := s.db.Model(&model.Site{}).Preload("Category").Preload("Tags")
	if categoryID := c.Query("category_id"); categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	if q := c.Query("q"); q != "" {
		like := "%" + q + "%"
		query = query.Where("title LIKE ? OR description LIKE ? OR url LIKE ?", like, like, like)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		fail(c, http.StatusInternalServerError, "SERVER_ERROR", "统计站点失败", nil)
		return
	}
	var sites []model.Site
	if err := query.Order("sort_order asc, id asc").Limit(pageSize).Offset((page - 1) * pageSize).Find(&sites).Error; err != nil {
		fail(c, http.StatusInternalServerError, "SERVER_ERROR", "读取站点失败", nil)
		return
	}
	ok(c, sites, map[string]any{"page": page, "pageSize": pageSize, "total": total})
}

func (s *Server) createSite(c *gin.Context) {
	var req siteInput
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, "VALIDATION_ERROR", "请求参数无效", err.Error())
		return
	}
	site := siteFromInput(req)
	if err := s.db.Create(&site).Error; err != nil {
		fail(c, http.StatusBadRequest, "CREATE_FAILED", "创建站点失败", err.Error())
		return
	}
	if !s.replaceSiteTags(c, &site, req.TagIDs) {
		return
	}
	s.db.Preload("Category").Preload("Tags").First(&site, site.ID)
	created(c, site)
}

func (s *Server) updateSite(c *gin.Context) {
	id, okID := parseID(c)
	if !okID {
		return
	}
	var req siteInput
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, "VALIDATION_ERROR", "请求参数无效", err.Error())
		return
	}
	var site model.Site
	if !findByID(c, s.db, id, &site, "站点不存在") {
		return
	}
	updated := siteFromInput(req)
	updated.ID = site.ID
	updated.ClickCount = site.ClickCount
	updated.CreatedAt = site.CreatedAt
	if err := s.db.Save(&updated).Error; err != nil {
		fail(c, http.StatusBadRequest, "UPDATE_FAILED", "更新站点失败", err.Error())
		return
	}
	if !s.replaceSiteTags(c, &updated, req.TagIDs) {
		return
	}
	s.db.Preload("Category").Preload("Tags").First(&updated, updated.ID)
	ok(c, updated, nil)
}

func (s *Server) deleteSite(c *gin.Context) {
	id, okID := parseID(c)
	if !okID {
		return
	}
	if err := s.db.Delete(&model.Site{}, id).Error; err != nil {
		fail(c, http.StatusBadRequest, "DELETE_FAILED", "删除站点失败", err.Error())
		return
	}
	ok(c, gin.H{"ok": true}, nil)
}

func (s *Server) listTags(c *gin.Context) {
	var tags []model.Tag
	if err := s.db.Order("name asc, id asc").Find(&tags).Error; err != nil {
		fail(c, http.StatusInternalServerError, "SERVER_ERROR", "读取标签失败", nil)
		return
	}
	ok(c, tags, nil)
}

func (s *Server) createTag(c *gin.Context) {
	var req tagInput
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, "VALIDATION_ERROR", "请求参数无效", err.Error())
		return
	}
	tag := model.Tag{Name: req.Name, Slug: req.Slug, Color: req.Color}
	if err := s.db.Create(&tag).Error; err != nil {
		fail(c, http.StatusBadRequest, "CREATE_FAILED", "创建标签失败", err.Error())
		return
	}
	created(c, tag)
}

func (s *Server) updateTag(c *gin.Context) {
	id, okID := parseID(c)
	if !okID {
		return
	}
	var req tagInput
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, "VALIDATION_ERROR", "请求参数无效", err.Error())
		return
	}
	var tag model.Tag
	if !findByID(c, s.db, id, &tag, "标签不存在") {
		return
	}
	tag.Name = req.Name
	tag.Slug = req.Slug
	tag.Color = req.Color
	if err := s.db.Save(&tag).Error; err != nil {
		fail(c, http.StatusBadRequest, "UPDATE_FAILED", "更新标签失败", err.Error())
		return
	}
	ok(c, tag, nil)
}

func (s *Server) deleteTag(c *gin.Context) {
	id, okID := parseID(c)
	if !okID {
		return
	}
	if err := s.db.Delete(&model.Tag{}, id).Error; err != nil {
		fail(c, http.StatusBadRequest, "DELETE_FAILED", "删除标签失败", err.Error())
		return
	}
	ok(c, gin.H{"ok": true}, nil)
}

func (s *Server) listSearchEngines(c *gin.Context) {
	var engines []model.SearchEngine
	if err := s.db.Order("sort_order asc, id asc").Find(&engines).Error; err != nil {
		fail(c, http.StatusInternalServerError, "SERVER_ERROR", "读取搜索引擎失败", nil)
		return
	}
	ok(c, engines, nil)
}

func (s *Server) createSearchEngine(c *gin.Context) {
	var req searchEngineInput
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, "VALIDATION_ERROR", "请求参数无效", err.Error())
		return
	}
	engine := searchEngineFromInput(req)
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		if engine.IsDefault {
			if err := tx.Model(&model.SearchEngine{}).Where("is_default = ?", true).Update("is_default", false).Error; err != nil {
				return err
			}
		}
		return tx.Create(&engine).Error
	}); err != nil {
		fail(c, http.StatusBadRequest, "CREATE_FAILED", "创建搜索引擎失败", err.Error())
		return
	}
	created(c, engine)
}

func (s *Server) updateSearchEngine(c *gin.Context) {
	id, okID := parseID(c)
	if !okID {
		return
	}
	var req searchEngineInput
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, "VALIDATION_ERROR", "请求参数无效", err.Error())
		return
	}
	var engine model.SearchEngine
	if !findByID(c, s.db, id, &engine, "搜索引擎不存在") {
		return
	}
	updated := searchEngineFromInput(req)
	updated.ID = engine.ID
	updated.CreatedAt = engine.CreatedAt
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		if updated.IsDefault {
			if err := tx.Model(&model.SearchEngine{}).Where("id <> ? AND is_default = ?", updated.ID, true).Update("is_default", false).Error; err != nil {
				return err
			}
		}
		return tx.Save(&updated).Error
	}); err != nil {
		fail(c, http.StatusBadRequest, "UPDATE_FAILED", "更新搜索引擎失败", err.Error())
		return
	}
	ok(c, updated, nil)
}

func (s *Server) deleteSearchEngine(c *gin.Context) {
	id, okID := parseID(c)
	if !okID {
		return
	}
	if err := s.db.Delete(&model.SearchEngine{}, id).Error; err != nil {
		fail(c, http.StatusBadRequest, "DELETE_FAILED", "删除搜索引擎失败", err.Error())
		return
	}
	ok(c, gin.H{"ok": true}, nil)
}

func (s *Server) listSettings(c *gin.Context) {
	settings, err := loadSettings(s.db)
	if err != nil {
		fail(c, http.StatusInternalServerError, "SERVER_ERROR", "读取设置失败", nil)
		return
	}
	ok(c, settings, nil)
}

func (s *Server) updateSettings(c *gin.Context) {
	var req settingsInput
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, "VALIDATION_ERROR", "请求参数无效", err.Error())
		return
	}
	for key, value := range req.Settings {
		setting := model.Setting{Key: key, Value: value, ValueType: "string"}
		if err := s.db.Save(&setting).Error; err != nil {
			fail(c, http.StatusBadRequest, "UPDATE_FAILED", "更新设置失败", err.Error())
			return
		}
	}
	settings, _ := loadSettings(s.db)
	ok(c, settings, nil)
}

func (s *Server) replaceSiteTags(c *gin.Context, site *model.Site, tagIDs []uint) bool {
	var tags []model.Tag
	if len(tagIDs) > 0 {
		if err := s.db.Where("id IN ?", tagIDs).Find(&tags).Error; err != nil {
			fail(c, http.StatusBadRequest, "VALIDATION_ERROR", "读取标签失败", err.Error())
			return false
		}
	}
	if err := s.db.Model(site).Association("Tags").Replace(tags); err != nil {
		fail(c, http.StatusBadRequest, "UPDATE_FAILED", "更新站点标签失败", err.Error())
		return false
	}
	return true
}

func siteFromInput(req siteInput) model.Site {
	visible := true
	if req.IsVisible != nil {
		visible = *req.IsVisible
	}
	return model.Site{
		CategoryID:   req.CategoryID,
		Title:        req.Title,
		URL:          req.URL,
		Description:  req.Description,
		IconURL:      req.IconURL,
		FallbackIcon: req.FallbackIcon,
		SortOrder:    req.SortOrder,
		IsPinned:     req.IsPinned,
		IsVisible:    visible,
	}
}

type dashboardStats struct {
	SiteCount         int                `json:"siteCount"`
	VisibleSiteCount  int                `json:"visibleSiteCount"`
	CategoryCount     int                `json:"categoryCount"`
	TagCount          int                `json:"tagCount"`
	SearchEngineCount int                `json:"searchEngineCount"`
	TotalClicks       int64              `json:"totalClicks"`
	PinnedCount       int                `json:"pinnedCount"`
	TopClicked        []topClickedSite   `json:"topClicked"`
	CategoryDist      []categoryDistItem `json:"categoryDist"`
}

type topClickedSite struct {
	ID         uint   `json:"id"`
	Title      string `json:"title"`
	ClickCount int64  `json:"clickCount"`
}

type categoryDistItem struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Count int64  `json:"count"`
}

func (s *Server) dashboardStats(c *gin.Context) {
	var stats dashboardStats

	// site count
	s.db.Model(&model.Site{}).Select("COUNT(*)").Scan(&stats.SiteCount)

	// visible site count
	s.db.Model(&model.Site{}).Where("is_visible = ?", true).Select("COUNT(*)").Scan(&stats.VisibleSiteCount)

	// category count
	s.db.Model(&model.Category{}).Select("COUNT(*)").Scan(&stats.CategoryCount)

	// tag count
	s.db.Model(&model.Tag{}).Select("COUNT(*)").Scan(&stats.TagCount)

	// search engine count
	s.db.Model(&model.SearchEngine{}).Select("COUNT(*)").Scan(&stats.SearchEngineCount)

	// total clicks
	s.db.Model(&model.Site{}).Select("COALESCE(SUM(click_count), 0)").Scan(&stats.TotalClicks)

	// pinned count
	s.db.Model(&model.Site{}).Where("is_pinned = ?", true).Select("COUNT(*)").Scan(&stats.PinnedCount)

	// top clicked sites (top 6)
	s.db.Model(&model.Site{}).
		Order("click_count desc, sort_order asc, id asc").
		Limit(6).
		Select("id, title, click_count").
		Find(&stats.TopClicked)

	// category distribution (top 5)
	s.db.Model(&model.Site{}).
		Select("categories.id, categories.name, COUNT(*) as count").
		Joins("JOIN categories ON categories.id = sites.category_id").
		Group("categories.id, categories.name").
		Order("count desc, categories.sort_order asc").
		Limit(5).
		Scan(&stats.CategoryDist)

	ok(c, stats, nil)
}

func searchEngineFromInput(req searchEngineInput) model.SearchEngine {
	visible := true
	if req.IsVisible != nil {
		visible = *req.IsVisible
	}
	return model.SearchEngine{
		Name:      req.Name,
		Slug:      req.Slug,
		SearchURL: req.SearchURL,
		Icon:      req.Icon,
		SortOrder: req.SortOrder,
		IsDefault: req.IsDefault,
		IsVisible: visible,
	}
}

func parseID(c *gin.Context) (uint, bool) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		fail(c, http.StatusBadRequest, "VALIDATION_ERROR", "ID 无效", nil)
		return 0, false
	}
	return uint(id), true
}

func findByID(c *gin.Context, db *gorm.DB, id uint, dest any, notFoundMessage string) bool {
	err := db.First(dest, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		fail(c, http.StatusNotFound, "NOT_FOUND", notFoundMessage, nil)
		return false
	}
	if err != nil {
		fail(c, http.StatusInternalServerError, "SERVER_ERROR", "读取数据失败", nil)
		return false
	}
	return true
}
