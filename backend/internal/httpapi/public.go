package httpapi

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"windnav/backend/internal/model"
)

func (s *Server) publicSummary(c *gin.Context) {
	settings, err := loadSettings(s.db)
	if err != nil {
		fail(c, http.StatusInternalServerError, "SERVER_ERROR", "读取站点配置失败", nil)
		return
	}
	ok(c, settings, nil)
}

func (s *Server) publicCategories(c *gin.Context) {
	var categories []model.Category
	if err := s.db.Where("is_visible = ?", true).Order("sort_order asc, id asc").Find(&categories).Error; err != nil {
		fail(c, http.StatusInternalServerError, "SERVER_ERROR", "读取分类失败", nil)
		return
	}
	ok(c, categories, nil)
}

func (s *Server) publicSites(c *gin.Context) {
	page, pageSize := pagination(c)
	query := s.db.Model(&model.Site{}).
		Preload("Category").
		Preload("Tags").
		Where("sites.is_visible = ?", true)

	if category := strings.TrimSpace(c.Query("category")); category != "" {
		query = query.Joins("JOIN categories ON categories.id = sites.category_id").Where("categories.slug = ? OR categories.id = ?", category, category)
	}
	if tag := strings.TrimSpace(c.Query("tag")); tag != "" {
		query = query.Joins("JOIN site_tags ON site_tags.site_id = sites.id").Joins("JOIN tags ON tags.id = site_tags.tag_id").Where("tags.slug = ? OR tags.id = ?", tag, tag)
	}
	if q := strings.TrimSpace(c.Query("q")); q != "" {
		like := "%" + strings.ToLower(q) + "%"
		query = query.Where("LOWER(sites.title) LIKE ? OR LOWER(sites.description) LIKE ? OR LOWER(sites.url) LIKE ?", like, like, like)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		fail(c, http.StatusInternalServerError, "SERVER_ERROR", "统计站点失败", nil)
		return
	}

	var sites []model.Site
	orderBy := siteOrder(c.Query("sort"), c.Query("order"))
	if err := query.Order(orderBy).Limit(pageSize).Offset((page - 1) * pageSize).Find(&sites).Error; err != nil {
		fail(c, http.StatusInternalServerError, "SERVER_ERROR", "读取站点失败", nil)
		return
	}
	ok(c, sites, map[string]any{"page": page, "pageSize": pageSize, "total": total})
}

func (s *Server) recordClick(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		fail(c, http.StatusBadRequest, "VALIDATION_ERROR", "站点 ID 无效", nil)
		return
	}
	if err := s.db.Model(&model.Site{}).Where("id = ? AND is_visible = ?", id, true).UpdateColumn("click_count", gorm.Expr("click_count + ?", 1)).Error; err != nil {
		fail(c, http.StatusInternalServerError, "SERVER_ERROR", "记录点击失败", nil)
		return
	}
	ok(c, gin.H{"ok": true}, nil)
}

func loadSettings(db *gorm.DB) (map[string]string, error) {
	var rows []model.Setting
	if err := db.Find(&rows).Error; err != nil {
		return nil, err
	}
	settings := make(map[string]string, len(rows))
	for _, row := range rows {
		settings[row.Key] = row.Value
	}
	return settings, nil
}

func pagination(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}

func siteOrder(sort, order string) string {
	allowed := map[string]string{
		"sort_order":  "sort_order",
		"created_at":  "created_at",
		"click_count": "click_count",
		"title":       "title",
		"updated_at":  "updated_at",
		"is_pinned":   "is_pinned",
	}
	column := allowed[sort]
	if column == "" {
		column = "is_pinned"
	}
	direction := "desc"
	if strings.EqualFold(order, "asc") {
		direction = "asc"
	}
	if column == "is_pinned" {
		return "is_pinned desc, sort_order asc, id asc"
	}
	return column + " " + direction + ", id asc"
}
