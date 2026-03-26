package service

import (
	"fmt"
	"reflect"

	"go_sleep_admin/internal/models"
)

type GenericModule struct {
	Slug   string   `json:"slug"`
	Title  string   `json:"title"`
	Fields []string `json:"fields"`
}

type ModuleService struct{}

func NewModuleService() *ModuleService { return &ModuleService{} }

func (s *ModuleService) Modules() []GenericModule {
	entries := []struct {
		slug  string
		title string
		model interface{}
	}{
		{"banner", "轮播图", models.Banner{}},
		{"audio", "音频", models.Audio{}},
		{"channel", "渠道", models.Channel{}},
		{"membership", "会员产品", models.Membership{}},
		{"user", "用户", models.User{}},
		{"user_profile", "用户画像", models.UserProfile{}},
		{"user_article_save", "用户收藏", models.UserArticleSave{}},
		{"order", "订单", models.Order{}},
		{"order_bill", "消费记录", models.OrderBill{}},
	}
	result := make([]GenericModule, 0, len(entries))
	for _, entry := range entries {
		result = append(result, GenericModule{
			Slug:   entry.slug,
			Title:  entry.title,
			Fields: fieldsOf(entry.model),
		})
	}
	return result
}

func fieldsOf(model interface{}) []string {
	t := reflect.TypeOf(model)
	fields := make([]string, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}
		fields = append(fields, fmt.Sprintf("%s:%s", field.Name, field.Type.String()))
	}
	return fields
}
