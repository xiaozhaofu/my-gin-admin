package models

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type BooleanField uint8
type ArticleType uint8
type ArticleStatus uint8
type ChannelStatus int8
type BannerStatus int8
type AudioStatus uint8

const (
	BooleanFalse BooleanField = iota
	BooleanTrue
)

const (
	ArticleTypeText ArticleType = iota + 1
	ArticleTypeVideo
	ArticleTypeVideoCollection
	ArticleTypeAudio
	ArticleTypeAudioCollection
	ArticleTypeImage
	ArticleTypeImageCollection
)

const (
	ArticleStatusNormal ArticleStatus = iota + 1
	ArticleStatusHidden
	ArticleStatusPending
	ArticleStatusRejected
)

const (
	ChannelStatusNormal ChannelStatus = iota + 1
	ChannelStatusDisabled
)

//type Menu struct {
//	BaseID
//	Name      string `gorm:"column:name;size:255;not null" json:"name"`
//	ParentID  *int64 `gorm:"column:parent_id;index" json:"parent_id"`
//	Level     uint8  `gorm:"column:level;not null;default:1" json:"level"`
//	Path      string `gorm:"column:path;size:255;not null;default:''" json:"path"`
//	SortOrder uint   `gorm:"column:sort_order;not null;default:0" json:"sort_order"`
//	IsActive  bool   `gorm:"column:is_active;not null" json:"is_active"`
//	AdminID   int64  `gorm:"column:admin_id;not null;index" json:"admin_id"`
//	BaseTimeField
//}

type Menu struct {
	BaseID

	Name      string `gorm:"column:name;not null;type:varchar(255);comment:分类名称" json:"name"`
	ParentID  *int64 `gorm:"column:parent_id;type:int(10) unsigned;default:NULL;index;comment:父级分类ID" json:"parent_id"`
	Level     uint8  `gorm:"column:level;type:tinyint(3) unsigned;default:1;comment:分类级别" json:"level"` // level = len(path的分段) - 1
	Path      string `gorm:"column:path;type:varchar(500);uniqueIndex;comment:路径，用于快速查询子树" json:"path"`
	SortOrder uint   `gorm:"column:sort_order;type:int(10) unsigned;default:0;comment:排序" json:"sort_order"`
	IsActive  bool   `gorm:"column:is_active;type:tinyint(1);default:1;comment:是否启用, 1:启用, 0:禁用" json:"is_active"`
	PagePath  string `gorm:"column:page_path;type:varchar(255);default:'';comment:'页面路径'" json:"page_path"`
	Icon      string `gorm:"column:icon;type:varchar(255);default:'';comment:'图标'" json:"icon"`
	AdminID   int64  `gorm:"not null;index;comment:创建人ID" json:"admin_id"`

	BaseTimeField
}

func (Menu) TableName() string { return "menus" }

func (m *Menu) BeforeCreate(tx *gorm.DB) error {
	if m.ParentID == nil {
		m.Level = 1
		return nil
	}
	var parent Menu
	if err := tx.First(&parent, *m.ParentID).Error; err != nil {
		return err
	}
	m.Level = parent.Level + 1
	return nil
}

func (m *Menu) AfterCreate(tx *gorm.DB) error {
	return updateMenuPath(tx, m.ID)
}

func (m *Menu) AfterUpdate(tx *gorm.DB) error {
	return updateMenuPath(tx, m.ID)
}

func updateMenuPath(tx *gorm.DB, menuID int64) error {
	var current Menu
	if err := tx.First(&current, menuID).Error; err != nil {
		return err
	}

	newPath := fmt.Sprintf("/%d/", current.ID)
	newLevel := uint8(1)
	if current.ParentID != nil {
		var parent Menu
		if err := tx.First(&parent, *current.ParentID).Error; err != nil {
			return err
		}
		newPath = strings.TrimRight(parent.Path, "/") + fmt.Sprintf("/%d/", current.ID)
		newLevel = parent.Level + 1
	}
	if err := tx.Session(&gorm.Session{SkipHooks: true}).Model(&Menu{}).Where("id = ?", current.ID).Updates(map[string]any{
		"path":  newPath,
		"level": newLevel,
	}).Error; err != nil {
		return err
	}
	current.Path = newPath
	current.Level = newLevel

	var children []Menu
	if err := tx.Where("parent_id = ?", current.ID).Find(&children).Error; err != nil {
		return err
	}
	for _, child := range children {
		if err := updateMenuPath(tx, child.ID); err != nil {
			return err
		}
	}
	return nil
}

//type Channel struct {
//	BaseID
//	Name    string        `gorm:"column:name;size:100;uniqueIndex;not null" json:"name"`
//	Status  ChannelStatus `gorm:"column:status;not null;default:1" json:"status"`
//	Remark  *string       `gorm:"column:remark;size:255;" json:"remark"`
//	AdminID int64         `gorm:"column:admin_id;not null;index" json:"admin_id"`
//	BaseTimeField
//}

type Channel struct {
	BaseID
	ID      int64         `gorm:"column:id;primaryKey;autoIncrement;comment:主键编码" json:"id"`
	Name    string        `gorm:"size:100;not null;unique;comment:渠道名称" json:"name"`
	Status  ChannelStatus `gorm:"column:status;default:1;index;comment:状态: 0-禁用,1-启用" json:"status"`
	Code    string        `gorm:"size:50;not null;uniqueIndex;comment:渠道编码，用于程序识别" json:"code"`
	Remark  string        `gorm:"size:255;comment:备注" json:"remark"`
	AdminID int64         `gorm:"not null;index:idx_admin_id;comment:创建人ID" json:"admin_id"`

	BaseTimeField
}

func (Channel) TableName() string { return "channels" }

type Article struct {
	BaseID
	Title       string        `gorm:"column:title;size:255;not null;index" json:"title"`
	Summary     string        `gorm:"column:summary;size:255;not null;default:''" json:"summary"`
	Type        ArticleType   `gorm:"column:type;not null;default:1;index" json:"type"`
	Cover       string        `gorm:"column:cover;size:255;not null;default:''" json:"cover"`
	CoverType   string        `gorm:"column:cover_type;size:20;not null;default:'1'" json:"cover_type"`
	MenuID      int64         `gorm:"column:menu_id;not null;index" json:"menu_id"`
	ChannelID   int64         `gorm:"column:channel_id;not null;index" json:"channel_id"`
	SortOrder   int8          `gorm:"column:sort_order;not null;default:1" json:"sort_order"`
	IsPaid      BooleanField  `gorm:"column:is_paid;not null;default:0" json:"is_paid"`
	AdminID     int64         `gorm:"column:admin_id;not null;index" json:"admin_id"`
	IsTop       BooleanField  `gorm:"column:is_top;not null;default:0" json:"is_top"`
	IsHot       BooleanField  `gorm:"column:is_hot;not null;default:0" json:"is_hot"`
	IsRecommend BooleanField  `gorm:"column:is_recommend;not null;default:0" json:"is_recommend"`
	CommentNum  int           `gorm:"column:comment_num;not null;default:0" json:"comment_num"`
	ShareNum    int           `gorm:"column:share_num;not null;default:0" json:"share_num"`
	ViewNum     int           `gorm:"column:view_num;not null;default:0" json:"view_num"`
	CollectNum  int           `gorm:"column:collect_num;not null;default:0" json:"collect_num"`
	Status      ArticleStatus `gorm:"column:status;not null;default:1;index" json:"status"`
	BaseTimeField
	Content ArticleContent `gorm:"foreignKey:ArticleID;constraint:OnDelete:CASCADE" json:"content"`
}

func (Article) TableName() string { return "articles" }

type ArticleContent struct {
	BaseID
	// 这里不直接用 GORM 的 uniqueIndex tag。
	// 原因是历史库里 article_id 可能已经存在旧索引名（例如 idx_article），
	// AutoMigrate 会尝试删旧索引再建新索引名，进而触发
	// “Cannot drop index ... needed in a foreign key constraint”。
	// 唯一索引由 bootstrap.autoMigrate 里的兼容逻辑补齐。
	ArticleID int64  `gorm:"column:article_id;not null" json:"article_id"`
	Content   string `gorm:"column:content;type:longtext;not null" json:"content"`
	BaseTimeFieldNoDelete
}

func (ArticleContent) TableName() string { return "article_contents" }

type UploadFile struct {
	BaseID
	Type       int8   `gorm:"column:type;not null;default:1;index" json:"type"`
	OriginName string `gorm:"column:origin_name;size:255;not null" json:"origin_name"`
	Random     string `gorm:"column:random;size:255;uniqueIndex;not null" json:"random"`
	Path       string `gorm:"column:path;size:255;not null" json:"path"`
	Md5        string `gorm:"column:md5;size:64;not null" json:"md5"`
	Scene      string `gorm:"column:scene;size:100;index;not null;default:'misc'" json:"scene"`
	Provider   string `gorm:"column:provider;size:32;index;not null;default:'local'" json:"provider"`
	AdminID    int64  `gorm:"column:admin_id;not null;index" json:"admin_id"`
	BaseTimeField
}

func (UploadFile) TableName() string { return "upload_files" }

type Banner struct {
	BaseID
	BannerImage   string       `gorm:"column:banner_image;size:255;not null" json:"banner_image"`
	BackgroundImg string       `gorm:"column:background_img;size:255;not null" json:"background_img"`
	Title         string       `gorm:"column:title;size:255;not null" json:"title"`
	SortOrder     int8         `gorm:"column:sort_order;not null;default:1" json:"sort_order"`
	Status        BannerStatus `gorm:"column:status;not null;default:1" json:"status"`
	AdminID       int64        `gorm:"column:admin_id;not null;index" json:"admin_id"`
	BaseTimeField
}

func (Banner) TableName() string { return "banners" }

type Audio struct {
	BaseID
	Title         string      `gorm:"column:title;size:255;not null" json:"title"`
	Type          int8        `gorm:"column:type;not null;default:1" json:"type"`
	BackgroundImg string      `gorm:"column:background_img;size:255;not null" json:"background_img"`
	MenuID        int64       `gorm:"column:menu_id;not null;index" json:"menu_id"`
	ChannelID     int64       `gorm:"column:channel_id;not null;index" json:"channel_id"`
	IsPaid        int8        `gorm:"column:is_paid;not null;default:0" json:"is_paid"`
	AdminID       int64       `gorm:"column:admin_id;not null;index" json:"admin_id"`
	IsTop         int8        `gorm:"column:is_top;not null;default:0" json:"is_top"`
	IsHot         int8        `gorm:"column:is_hot;not null;default:0" json:"is_hot"`
	IsRecommend   int8        `gorm:"column:is_recommend;not null;default:0" json:"is_recommend"`
	CommentNum    int         `gorm:"column:comment_num;not null;default:0" json:"comment_num"`
	ShareNum      int         `gorm:"column:share_num;not null;default:0" json:"share_num"`
	ViewNum       int         `gorm:"column:view_num;not null;default:0" json:"view_num"`
	CollectNum    int         `gorm:"column:collect_num;not null;default:0" json:"collect_num"`
	Status        AudioStatus `gorm:"column:status;not null;default:1" json:"status"`
	BaseTimeField
}

func (Audio) TableName() string { return "audios" }
