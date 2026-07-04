package model

import "time"

type User struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	Username     string     `gorm:"size:80;uniqueIndex;not null" json:"username"`
	PasswordHash string     `gorm:"size:255;not null" json:"-"`
	DisplayName  string     `gorm:"size:120" json:"displayName"`
	Role         string     `gorm:"size:40;default:admin" json:"role"`
	IsActive     bool       `gorm:"default:true" json:"isActive"`
	LastLoginAt  *time.Time `json:"lastLoginAt"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
}

type Category struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:120;not null" json:"name"`
	Slug        string    `gorm:"size:120;uniqueIndex;not null" json:"slug"`
	Description string    `gorm:"size:500" json:"description"`
	Icon        string    `gorm:"size:120" json:"icon"`
	Color       string    `gorm:"size:40" json:"color"`
	SortOrder   int       `gorm:"default:0;index" json:"sortOrder"`
	IsVisible   bool      `gorm:"default:true;index" json:"isVisible"`
	Sites       []Site    `json:"sites,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type Site struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	CategoryID   uint      `gorm:"index" json:"categoryId"`
	Category     Category  `json:"category,omitempty"`
	Title        string    `gorm:"size:160;not null;index" json:"title"`
	URL          string    `gorm:"size:1000;not null" json:"url"`
	Description  string    `gorm:"size:800" json:"description"`
	IconURL      string    `gorm:"size:1000" json:"iconUrl"`
	FallbackIcon string    `gorm:"size:40" json:"fallbackIcon"`
	SortOrder    int       `gorm:"default:0;index" json:"sortOrder"`
	IsPinned     bool      `gorm:"default:false;index" json:"isPinned"`
	IsVisible    bool      `gorm:"default:true;index" json:"isVisible"`
	ClickCount   int64     `gorm:"default:0" json:"clickCount"`
	Tags         []Tag     `gorm:"many2many:site_tags;" json:"tags,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type Tag struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:80;not null" json:"name"`
	Slug      string    `gorm:"size:100;uniqueIndex;not null" json:"slug"`
	Color     string    `gorm:"size:40" json:"color"`
	Sites     []Site    `gorm:"many2many:site_tags;" json:"sites,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Setting struct {
	Key       string    `gorm:"primaryKey;size:120" json:"key"`
	Value     string    `gorm:"type:text" json:"value"`
	ValueType string    `gorm:"size:40;default:string" json:"valueType"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type SearchEngine struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:120;not null" json:"name"`
	Slug        string    `gorm:"size:120;uniqueIndex;not null" json:"slug"`
	SearchURL   string    `gorm:"size:1000;not null" json:"searchUrl"`
	Icon        string    `gorm:"size:120" json:"icon"`
	SortOrder   int       `gorm:"default:0;index" json:"sortOrder"`
	IsDefault   bool      `gorm:"default:false;index" json:"isDefault"`
	IsVisible   bool      `gorm:"default:true;index" json:"isVisible"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type RefreshToken struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	UserID    uint       `gorm:"index;not null" json:"userId"`
	TokenHash string     `gorm:"size:255;not null;index" json:"-"`
	ExpiresAt time.Time  `gorm:"index" json:"expiresAt"`
	RevokedAt *time.Time `json:"revokedAt"`
	CreatedAt time.Time  `json:"createdAt"`
}
