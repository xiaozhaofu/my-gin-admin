package service

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go_sleep_admin/internal/dto"
	"go_sleep_admin/internal/repository"
)

const maxOrderExportRows = 10000

type OrderService struct {
	repo *repository.OrderRepository
}

func NewOrderService(repo *repository.OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) List(query dto.OrderQuery) ([]dto.OrderListItem, int64, error) {
	return s.repo.List(query)
}

func (s *OrderService) Detail(id int64) (*dto.OrderDetail, error) {
	return s.repo.ByID(id)
}

func (s *OrderService) Export(query dto.OrderQuery) ([]byte, error) {
	exportQuery := query
	exportQuery.Page = 1
	exportQuery.PageSize = maxOrderExportRows

	items, total, err := s.repo.List(exportQuery)
	if err != nil {
		return nil, err
	}
	if total > maxOrderExportRows {
		return nil, fmt.Errorf("最多导出 %d 条订单，请先缩小筛选范围", maxOrderExportRows)
	}

	var buf bytes.Buffer
	buf.Write([]byte{0xEF, 0xBB, 0xBF})

	writer := csv.NewWriter(&buf)
	if err := writer.Write([]string{
		"序号",
		"订单号",
		"订单凭证",
		"第三方交易号",
		"用户ID",
		"用户昵称",
		"用户手机号",
		"商品ID",
		"商品标题",
		"商品类型",
		"订单状态",
		"退款状态",
		"支付平台",
		"支付平台标识",
		"支付方式",
		"渠道名称",
		"渠道编码",
		"原价(元)",
		"优惠金额(元)",
		"实付金额(元)",
		"退款金额(元)",
		"支付时间",
		"订单时间",
	}); err != nil {
		return nil, err
	}

	for idx, item := range items {
		if err := writer.Write([]string{
			strconv.Itoa(idx + 1),
			item.OrderNo,
			item.OrderToken,
			orderFirstNonEmpty(item.TradeNo, "-"),
			strconv.FormatInt(item.UserID, 10),
			item.UserNickname,
			item.UserPhone,
			strconv.FormatUint(item.ProductID, 10),
			item.ProductTitle,
			orderProductTypeLabel(item.ProductType),
			orderStatusLabel(item.Status),
			orderRefundStatusLabel(item.RefundStatus),
			orderPayChannelLabel(item.PayChannel),
			orderFirstNonEmpty(item.PayChannel, "-"),
			orderPayMethodLabel(item.PayMethod),
			orderFirstNonEmpty(item.ChannelName, "-"),
			orderFirstNonEmpty(item.ChannelCode, "-"),
			orderFormatFen(item.OriginalPrice),
			orderFormatFen(item.DiscountPrice),
			orderFormatFen(item.PayAmount),
			orderFormatFen(item.RefundAmount),
			orderFormatTime(item.PaidAt),
			item.CreatedAt.Format(time.DateTime),
		}); err != nil {
			return nil, err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func orderFormatFen(value int64) string {
	return fmt.Sprintf("%.2f", float64(value)/100)
}

func orderFormatTime(value *time.Time) string {
	if value == nil || value.IsZero() {
		return "-"
	}
	return value.Format(time.DateTime)
}

func orderFirstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func orderStatusLabel(status int8) string {
	switch status {
	case 0:
		return "待支付"
	case 10:
		return "已支付"
	case 20:
		return "已开通"
	case 30:
		return "已完成"
	case 40:
		return "已关闭"
	case 50:
		return "已取消"
	case 60:
		return "退款中"
	case 70:
		return "已退款"
	default:
		return strconv.FormatInt(int64(status), 10)
	}
}

func orderRefundStatusLabel(status int8) string {
	switch status {
	case 0:
		return "无"
	case 1:
		return "已申请"
	case 2:
		return "审核通过"
	case 3:
		return "退款成功"
	case 4:
		return "退款拒绝"
	default:
		return strconv.FormatInt(int64(status), 10)
	}
}

func orderPayMethodLabel(method int8) string {
	switch method {
	case 1:
		return "微信支付"
	case 2:
		return "支付宝"
	case 3:
		return "银联支付"
	default:
		return strconv.FormatInt(int64(method), 10)
	}
}

func orderPayChannelLabel(channel string) string {
	switch strings.TrimSpace(strings.ToLower(channel)) {
	case "wechat":
		return "微信"
	case "alipay":
		return "支付宝"
	case "unionpay":
		return "银联"
	default:
		return orderFirstNonEmpty(channel, "-")
	}
}

func orderProductTypeLabel(productType int8) string {
	switch productType {
	case 1:
		return "文本"
	case 2:
		return "视频"
	case 3:
		return "音频"
	case 4:
		return "会员"
	default:
		return strconv.FormatInt(int64(productType), 10)
	}
}
