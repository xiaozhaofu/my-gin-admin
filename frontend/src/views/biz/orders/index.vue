<template>
  <div class="snow-page">
    <div class="order-page">
      <section class="page-hero">
        <div>
          <div class="page-badge">Order Center</div>
          <h1>订单管理</h1>
          <p>按订单号和订单时间快速检索支付订单，统一查看支付平台、第三方交易号、渠道信息，并支持 CSV 导出和完整订单详情核对。</p>
        </div>
        <div class="page-hero-side">
          <div class="page-hero-label">当前订单总数</div>
          <div class="page-hero-value">{{ pagination.total }}</div>
        </div>
      </section>

      <a-card class="panel-card" :bordered="false">
        <template #title>筛选条件</template>
        <a-form :model="filters" layout="inline" class="filters">
          <a-form-item field="order_no" label="订单号">
            <a-input v-model="filters.order_no" placeholder="输入订单号" allow-clear />
          </a-form-item>
          <a-form-item field="trade_no" label="交易号">
            <a-input v-model="filters.trade_no" placeholder="输入第三方交易号" allow-clear />
          </a-form-item>
          <a-form-item field="user_phone" label="手机号">
            <a-input v-model="filters.user_phone" placeholder="输入用户手机号" allow-clear />
          </a-form-item>
          <a-form-item field="pay_channel" label="支付平台">
            <a-select v-model="filters.pay_channel" placeholder="全部平台" allow-clear style="width: 150px">
              <a-option v-for="item in payChannelOptions" :key="item.value" :value="item.value">{{ item.label }}</a-option>
            </a-select>
          </a-form-item>
          <a-form-item field="pay_method" label="支付方式">
            <a-select v-model="filters.pay_method" placeholder="全部方式" allow-clear style="width: 150px">
              <a-option v-for="item in payMethodOptions" :key="item.value" :value="item.value">{{ item.label }}</a-option>
            </a-select>
          </a-form-item>
          <a-form-item field="channel_id" label="渠道">
            <a-select v-model="filters.channel_id" placeholder="全部渠道" allow-clear style="width: 180px">
              <a-option v-for="item in channels" :key="item.id" :value="item.id">{{ item.name }}（{{ item.code }}）</a-option>
            </a-select>
          </a-form-item>
          <a-form-item field="status" label="订单状态">
            <a-select v-model="filters.status" placeholder="全部状态" allow-clear style="width: 150px">
              <a-option v-for="item in orderStatusOptions" :key="item.value" :value="item.value">{{ item.label }}</a-option>
            </a-select>
          </a-form-item>
          <a-form-item field="refund_status" label="退款状态">
            <a-select v-model="filters.refund_status" placeholder="全部退款状态" allow-clear style="width: 150px">
              <a-option v-for="item in refundStatusOptions" :key="item.value" :value="item.value">{{ item.label }}</a-option>
            </a-select>
          </a-form-item>
          <a-form-item field="created_range" label="订单时间">
            <a-range-picker
              v-model="dateRange"
              show-time
              format="YYYY-MM-DD HH:mm:ss"
              value-format="YYYY-MM-DD HH:mm:ss"
              style="width: 320px"
              allow-clear
            />
          </a-form-item>
          <a-space wrap>
            <a-button type="primary" @click="fetchList">
              <template #icon><icon-search /></template>
              查询
            </a-button>
            <a-button @click="resetFilters">
              <template #icon><icon-refresh /></template>
              重置
            </a-button>
            <a-button v-if="canExport" type="outline" @click="exportOrders">
              <template #icon><icon-export /></template>
              导出订单
            </a-button>
          </a-space>
        </a-form>
      </a-card>

      <a-card class="panel-card" :bordered="false">
        <template #title>订单列表</template>
        <a-table
          row-key="id"
          :data="list"
          :loading="loading"
          :pagination="pagination"
          @page-change="pageChange"
        >
          <template #columns>
            <a-table-column title="订单号" :width="220">
              <template #cell="{ record }">
                <div class="stack-cell">
                  <div class="copy-row">
                    <div class="stack-title">{{ record.order_no }}</div>
                    <a-button size="mini" type="text" @click="copyText(record.order_no, '订单号')">
                      <template #icon><icon-copy /></template>
                    </a-button>
                  </div>
                  <div class="stack-sub">订单凭证：{{ record.order_token || "-" }}</div>
                </div>
              </template>
            </a-table-column>
            <a-table-column title="商品信息" :width="220">
              <template #cell="{ record }">
                <div class="stack-cell">
                  <div class="stack-title">{{ record.product_title }}</div>
                  <div class="stack-sub">{{ productTypeLabel(record.product_type) }}</div>
                </div>
              </template>
            </a-table-column>
            <a-table-column title="用户" :width="180">
              <template #cell="{ record }">
                <div class="stack-cell">
                  <div class="stack-title">{{ record.user_nickname || `用户#${record.user_id}` }}</div>
                  <div class="copy-row">
                    <div class="stack-sub">{{ record.user_phone || "-" }}</div>
                    <a-button v-if="record.user_phone" size="mini" type="text" @click="copyText(record.user_phone, '手机号')">
                      <template #icon><icon-copy /></template>
                    </a-button>
                  </div>
                </div>
              </template>
            </a-table-column>
            <a-table-column title="支付信息" :width="220">
              <template #cell="{ record }">
                <div class="stack-cell">
                  <div class="stack-title">{{ payChannelLabel(record.pay_channel) }}</div>
                  <div class="stack-sub">方式：{{ payMethodLabel(record.pay_method) }}</div>
                  <div class="copy-row">
                    <div class="stack-sub">交易号：{{ record.trade_no || "-" }}</div>
                    <a-button v-if="record.trade_no" size="mini" type="text" @click="copyText(record.trade_no, '交易号')">
                      <template #icon><icon-copy /></template>
                    </a-button>
                  </div>
                </div>
              </template>
            </a-table-column>
            <a-table-column title="渠道" :width="180">
              <template #cell="{ record }">
                <div class="stack-cell">
                  <div class="stack-title">{{ record.channel_name || "-" }}</div>
                  <div class="copy-row">
                    <div class="stack-sub">{{ record.channel_code || "-" }}</div>
                    <a-button v-if="record.channel_code" size="mini" type="text" @click="copyText(record.channel_code, '渠道编码')">
                      <template #icon><icon-copy /></template>
                    </a-button>
                  </div>
                </div>
              </template>
            </a-table-column>
            <a-table-column title="金额" :width="140">
              <template #cell="{ record }">
                <div class="stack-cell">
                  <div class="stack-title money">¥{{ formatFen(record.pay_amount) }}</div>
                  <div class="stack-sub">原价 ¥{{ formatFen(record.original_price) }}</div>
                </div>
              </template>
            </a-table-column>
            <a-table-column title="状态" :width="140">
              <template #cell="{ record }">
                <a-space direction="vertical" size="mini">
                  <a-tag :color="orderStatusColor(record.status)">{{ orderStatusLabel(record.status) }}</a-tag>
                  <a-tag :color="refundStatusColor(record.refund_status)" bordered>{{ refundStatusLabel(record.refund_status) }}</a-tag>
                </a-space>
              </template>
            </a-table-column>
            <a-table-column title="订单时间" :width="180">
              <template #cell="{ record }">
                <div class="stack-cell">
                  <div class="stack-title">{{ formatDateTime(record.created_at) }}</div>
                  <div class="stack-sub">支付：{{ formatDateTime(record.paid_at) }}</div>
                </div>
              </template>
            </a-table-column>
            <a-table-column title="操作" :width="120" align="center" fixed="right">
              <template #cell="{ record }">
                <a-button size="mini" type="outline" @click="viewDetail(record.id)">
                  <template #icon><icon-eye /></template>
                  详情
                </a-button>
              </template>
            </a-table-column>
          </template>
        </a-table>
      </a-card>
    </div>
  </div>

  <a-modal v-model:visible="detailVisible" title="订单详情" width="960px" :footer="false">
    <div v-if="detail" class="detail-wrap">
      <div class="detail-summary-grid">
        <div v-for="card in detailSummaryCards" :key="card.label" class="detail-summary-card" :class="card.tone">
          <div class="detail-summary-label">{{ card.label }}</div>
          <div class="detail-summary-value">{{ card.value }}</div>
          <div class="detail-summary-sub">{{ card.sub }}</div>
        </div>
      </div>

      <a-card class="detail-card" :bordered="false" title="订单主表">
        <a-descriptions :column="2" bordered>
          <a-descriptions-item label="订单号">
            <div class="detail-copy-row">
              <span>{{ detail.order.order_no }}</span>
              <a-button size="mini" type="text" @click="copyText(detail.order.order_no, '订单号')">
                <template #icon><icon-copy /></template>
              </a-button>
            </div>
          </a-descriptions-item>
          <a-descriptions-item label="订单凭证">
            <div class="detail-copy-row">
              <span>{{ detail.order.order_token || "-" }}</span>
              <a-button v-if="detail.order.order_token" size="mini" type="text" @click="copyText(detail.order.order_token, '订单凭证')">
                <template #icon><icon-copy /></template>
              </a-button>
            </div>
          </a-descriptions-item>
          <a-descriptions-item label="订单状态">
            <a-tag :color="orderStatusColor(detail.order.status)">{{ orderStatusLabel(detail.order.status) }}</a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="退款状态">
            <a-tag :color="refundStatusColor(detail.order.refund_status)">{{ refundStatusLabel(detail.order.refund_status) }}</a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="商品标题">{{ detail.order.product_title }}</a-descriptions-item>
          <a-descriptions-item label="商品类型">{{ productTypeLabel(detail.order.product_type) }}</a-descriptions-item>
          <a-descriptions-item label="原价">¥{{ formatFen(detail.order.original_price) }}</a-descriptions-item>
          <a-descriptions-item label="折扣金额">¥{{ formatFen(detail.order.discount_price) }}</a-descriptions-item>
          <a-descriptions-item label="优惠券ID">{{ detail.order.coupon_id ?? "-" }}</a-descriptions-item>
          <a-descriptions-item label="优惠券金额">¥{{ formatFen(detail.order.coupon_amount) }}</a-descriptions-item>
          <a-descriptions-item label="实付金额">¥{{ formatFen(detail.order.pay_amount) }}</a-descriptions-item>
          <a-descriptions-item label="支付方式">{{ payMethodLabel(detail.order.pay_method) }}</a-descriptions-item>
          <a-descriptions-item label="支付平台">{{ payChannelLabel(detail.order.pay_channel) }}</a-descriptions-item>
          <a-descriptions-item label="第三方交易号">
            <div class="detail-copy-row">
              <span>{{ detail.order.trade_no || "-" }}</span>
              <a-button v-if="detail.order.trade_no" size="mini" type="text" @click="copyText(detail.order.trade_no, '交易号')">
                <template #icon><icon-copy /></template>
              </a-button>
            </div>
          </a-descriptions-item>
          <a-descriptions-item label="渠道ID">{{ detail.order.channel_id }}</a-descriptions-item>
          <a-descriptions-item label="支付时间">{{ formatDateTime(detail.order.paid_at) }}</a-descriptions-item>
          <a-descriptions-item label="订单时间">{{ formatDateTime(detail.order.created_at) }}</a-descriptions-item>
          <a-descriptions-item label="过期时间">{{ formatDateTime(detail.order.expire_at) }}</a-descriptions-item>
          <a-descriptions-item label="权益到期">{{ formatDateTime(detail.order.access_expire_at) }}</a-descriptions-item>
          <a-descriptions-item label="发放时间">{{ formatDateTime(detail.order.delivered_at) }}</a-descriptions-item>
          <a-descriptions-item label="退款金额">¥{{ formatFen(detail.order.refund_amount) }}</a-descriptions-item>
          <a-descriptions-item label="退款单号">{{ detail.order.refund_no || "-" }}</a-descriptions-item>
          <a-descriptions-item label="退款时间">{{ formatDateTime(detail.order.refund_at) }}</a-descriptions-item>
          <a-descriptions-item label="客户端IP">{{ detail.order.client_ip || "-" }}</a-descriptions-item>
          <a-descriptions-item label="客户端IP原值">{{ detail.order.client_ip_raw || "-" }}</a-descriptions-item>
          <a-descriptions-item label="用户备注" :span="2">{{ detail.order.remark || "-" }}</a-descriptions-item>
          <a-descriptions-item label="后台备注" :span="2">{{ detail.order.admin_remark || "-" }}</a-descriptions-item>
          <a-descriptions-item label="退款原因" :span="2">{{ detail.order.refund_reason || "-" }}</a-descriptions-item>
        </a-descriptions>
      </a-card>

      <a-card class="detail-card" :bordered="false" title="状态时间线">
        <div class="timeline-list">
          <div v-for="item in orderTimeline" :key="`${item.title}-${item.time}`" class="timeline-item">
            <div class="timeline-dot" :class="item.tone"></div>
            <div class="timeline-body">
              <div class="timeline-title">{{ item.title }}</div>
              <div class="timeline-time">{{ item.time }}</div>
              <div v-if="item.desc" class="timeline-desc">{{ item.desc }}</div>
            </div>
          </div>
        </div>
      </a-card>

      <a-card class="detail-card" :bordered="false" title="关联支付流水">
        <a-descriptions :column="2" bordered>
          <template v-if="detail.bill">
            <a-descriptions-item label="流水ID">{{ detail.bill.id }}</a-descriptions-item>
            <a-descriptions-item label="订单号">{{ detail.bill.order_no }}</a-descriptions-item>
            <a-descriptions-item label="交易号">
              <div class="detail-copy-row">
                <span>{{ detail.bill.trade_no || "-" }}</span>
                <a-button v-if="detail.bill.trade_no" size="mini" type="text" @click="copyText(detail.bill.trade_no, '流水交易号')">
                  <template #icon><icon-copy /></template>
                </a-button>
              </div>
            </a-descriptions-item>
            <a-descriptions-item label="商品标题">{{ detail.bill.product_title }}</a-descriptions-item>
            <a-descriptions-item label="原价">¥{{ formatFen(detail.bill.original_price) }}</a-descriptions-item>
            <a-descriptions-item label="优惠金额">¥{{ formatFen(detail.bill.discount_amount) }}</a-descriptions-item>
            <a-descriptions-item label="实付金额">¥{{ formatFen(detail.bill.pay_amount) }}</a-descriptions-item>
            <a-descriptions-item label="支付方式">{{ payMethodLabel(detail.bill.pay_method) }}</a-descriptions-item>
            <a-descriptions-item label="支付平台">{{ payChannelLabel(detail.bill.pay_channel) }}</a-descriptions-item>
            <a-descriptions-item label="支付时间">{{ formatDateTime(detail.bill.paid_at) }}</a-descriptions-item>
            <a-descriptions-item label="退款状态">{{ refundStatusLabel(detail.bill.refund_status) }}</a-descriptions-item>
            <a-descriptions-item label="退款金额">¥{{ formatFen(detail.bill.refund_amount) }}</a-descriptions-item>
            <a-descriptions-item label="创建时间">{{ formatDateTime(detail.bill.created_at) }}</a-descriptions-item>
          </template>
          <template v-else>
            <a-descriptions-item label="支付流水" :span="2">当前订单暂无关联流水记录</a-descriptions-item>
          </template>
        </a-descriptions>
      </a-card>

      <div class="detail-grid">
        <a-card class="detail-card" :bordered="false" title="关联用户">
          <a-descriptions :column="1" bordered>
            <template v-if="detail.user">
              <a-descriptions-item label="用户ID">{{ detail.user.id }}</a-descriptions-item>
              <a-descriptions-item label="昵称">{{ detail.user.nickname || "-" }}</a-descriptions-item>
              <a-descriptions-item label="手机号">
                <div class="detail-copy-row">
                  <span>{{ detail.user.phone || "-" }}</span>
                  <a-button v-if="detail.user.phone" size="mini" type="text" @click="copyText(detail.user.phone, '手机号')">
                    <template #icon><icon-copy /></template>
                  </a-button>
                </div>
              </a-descriptions-item>
              <a-descriptions-item label="邮箱">{{ detail.user.email || "-" }}</a-descriptions-item>
              <a-descriptions-item label="状态">{{ detail.user.status }}</a-descriptions-item>
              <a-descriptions-item label="会员状态">{{ detail.user.vip_status }}</a-descriptions-item>
              <a-descriptions-item label="会员等级">{{ detail.user.vip_level }}</a-descriptions-item>
              <a-descriptions-item label="会员到期">{{ formatDateTime(detail.user.vip_expire_at) }}</a-descriptions-item>
              <a-descriptions-item label="注册时间">{{ formatDateTime(detail.user.created_at) }}</a-descriptions-item>
            </template>
            <template v-else>
              <a-descriptions-item label="用户">未找到关联用户</a-descriptions-item>
            </template>
          </a-descriptions>
        </a-card>

        <a-card class="detail-card" :bordered="false" title="关联渠道">
          <a-descriptions :column="1" bordered>
            <template v-if="detail.channel">
              <a-descriptions-item label="渠道ID">{{ detail.channel.id }}</a-descriptions-item>
              <a-descriptions-item label="渠道名称">{{ detail.channel.name }}</a-descriptions-item>
              <a-descriptions-item label="渠道编码">
                <div class="detail-copy-row">
                  <span>{{ detail.channel.code }}</span>
                  <a-button v-if="detail.channel.code" size="mini" type="text" @click="copyText(detail.channel.code, '渠道编码')">
                    <template #icon><icon-copy /></template>
                  </a-button>
                </div>
              </a-descriptions-item>
              <a-descriptions-item label="渠道状态">{{ detail.channel.status === 1 ? "启用" : "禁用" }}</a-descriptions-item>
              <a-descriptions-item label="备注">{{ detail.channel.remark || "-" }}</a-descriptions-item>
              <a-descriptions-item label="创建时间">{{ formatDateTime(detail.channel.created_at) }}</a-descriptions-item>
            </template>
            <template v-else>
              <a-descriptions-item label="渠道">未找到关联渠道</a-descriptions-item>
            </template>
          </a-descriptions>
        </a-card>
      </div>
    </div>
  </a-modal>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from "vue";
import { Message } from "@arco-design/web-vue";
import { IconCopy, IconExport, IconEye, IconRefresh, IconSearch } from "@arco-design/web-vue/es/icon";
import { channelListAPI, type ChannelItem } from "@/api/channel";
import { orderDetailAPI, orderExportAPI, orderListAPI, type OrderDetail, type OrderListItem } from "@/api/order";
import { useSessionStore } from "@/store/modules/session";

const session = useSessionStore();
const list = ref<OrderListItem[]>([]);
const channels = ref<ChannelItem[]>([]);
const loading = ref(false);
const detailVisible = ref(false);
const detail = ref<OrderDetail>();
const dateRange = ref<string[]>([]);
const filters = reactive({
  order_no: "",
  trade_no: "",
  user_phone: "",
  pay_channel: undefined as string | undefined,
  pay_method: undefined as number | undefined,
  channel_id: undefined as number | undefined,
  status: undefined as number | undefined,
  refund_status: undefined as number | undefined
});
const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0
});

const canExport = computed(() => session.can("/api/v1/orders/export#GET"));

const payChannelOptions = [
  { label: "微信支付", value: "wechat" },
  { label: "支付宝", value: "alipay" },
  { label: "银联支付", value: "UnionPay" }
];

const payMethodOptions = [
  { label: "微信支付", value: 1 },
  { label: "支付宝", value: 2 },
  { label: "银联支付", value: 3 }
];

const orderStatusOptions = [
  { label: "待支付", value: 0 },
  { label: "已支付", value: 10 },
  { label: "已开通", value: 20 },
  { label: "已完成", value: 30 },
  { label: "已关闭", value: 40 },
  { label: "已取消", value: 50 },
  { label: "退款中", value: 60 },
  { label: "已退款", value: 70 }
];

const refundStatusOptions = [
  { label: "无", value: 0 },
  { label: "已申请", value: 1 },
  { label: "审核通过", value: 2 },
  { label: "退款成功", value: 3 },
  { label: "退款拒绝", value: 4 }
];

const buildParams = () => ({
  order_no: filters.order_no,
  trade_no: filters.trade_no,
  user_phone: filters.user_phone,
  pay_channel: filters.pay_channel,
  pay_method: filters.pay_method,
  channel_id: filters.channel_id,
  status: filters.status,
  refund_status: filters.refund_status,
  created_from: dateRange.value?.[0] || undefined,
  created_to: dateRange.value?.[1] || undefined,
  page: pagination.current,
  page_size: pagination.pageSize
});

const fetchList = async () => {
  loading.value = true;
  try {
    const res = await orderListAPI(buildParams());
    list.value = res.data.list;
    pagination.total = res.data.total;
  } finally {
    loading.value = false;
  }
};

const fetchChannels = async () => {
  const res = await channelListAPI();
  channels.value = res.data.filter(item => item.status === 1);
};

const resetFilters = () => {
  filters.order_no = "";
  filters.trade_no = "";
  filters.user_phone = "";
  filters.pay_channel = undefined;
  filters.pay_method = undefined;
  filters.channel_id = undefined;
  filters.status = undefined;
  filters.refund_status = undefined;
  dateRange.value = [];
  pagination.current = 1;
  fetchList();
};

const pageChange = (page: number) => {
  pagination.current = page;
  fetchList();
};

const viewDetail = async (id: number) => {
  const res = await orderDetailAPI(id);
  detail.value = res.data;
  detailVisible.value = true;
};

const exportOrders = async () => {
  const blob = new Blob([await orderExportAPI({
    order_no: filters.order_no,
    trade_no: filters.trade_no,
    user_phone: filters.user_phone,
    pay_channel: filters.pay_channel,
    pay_method: filters.pay_method,
    channel_id: filters.channel_id,
    status: filters.status,
    refund_status: filters.refund_status,
    created_from: dateRange.value?.[0] || undefined,
    created_to: dateRange.value?.[1] || undefined
  })], { type: "text/csv;charset=utf-8" });
  const url = window.URL.createObjectURL(blob);
  const link = document.createElement("a");
  link.href = url;
  link.download = `orders_${new Date().toISOString().slice(0, 19).replace(/:/g, "-")}.csv`;
  link.click();
  window.URL.revokeObjectURL(url);
  Message.success("导出成功");
};

const copyText = async (value: string, label: string) => {
  if (!value) return;
  await navigator.clipboard.writeText(value);
  Message.success(`${label}已复制`);
};

const detailSummaryCards = computed(() => {
  if (!detail.value) return [];
  const order = detail.value.order;
  return [
    {
      label: "实付金额",
      value: `¥${formatFen(order.pay_amount)}`,
      sub: `原价 ¥${formatFen(order.original_price)} / 优惠 ¥${formatFen(order.discount_price)}`,
      tone: "success"
    },
    {
      label: "订单状态",
      value: orderStatusLabel(order.status),
      sub: `退款状态：${refundStatusLabel(order.refund_status)}`,
      tone: order.status === 10 || order.status === 20 || order.status === 30 ? "info" : "muted"
    },
    {
      label: "支付信息",
      value: payMethodLabel(order.pay_method),
      sub: `${payChannelLabel(order.pay_channel)} / ${order.trade_no || "无交易号"}`,
      tone: "warning"
    },
    {
      label: "渠道与用户",
      value: detail.value.channel?.name || `渠道#${order.channel_id}`,
      sub: detail.value.user?.phone || detail.value.user?.nickname || `用户#${order.user_id}`,
      tone: "purple"
    }
  ];
});

const orderTimeline = computed(() => {
  if (!detail.value) return [];
  const order = detail.value.order;
  const items = [
    {
      title: "订单创建",
      time: formatDateTime(order.created_at),
      desc: `订单号：${order.order_no}`,
      tone: "info"
    }
  ];

  if (order.paid_at) {
    items.push({
      title: "支付成功",
      time: formatDateTime(order.paid_at),
      desc: `支付方式：${payMethodLabel(order.pay_method)} / 平台：${payChannelLabel(order.pay_channel)}`,
      tone: "success"
    });
  }

  if (order.delivered_at) {
    items.push({
      title: "权益发放",
      time: formatDateTime(order.delivered_at),
      desc: order.product_title,
      tone: "success"
    });
  }

  if (order.access_expire_at) {
    items.push({
      title: "权益到期",
      time: formatDateTime(order.access_expire_at),
      desc: "到期后需要重新续费或重新开通",
      tone: "warning"
    });
  }

  if (order.refund_at) {
    items.push({
      title: "退款处理",
      time: formatDateTime(order.refund_at),
      desc: `${refundStatusLabel(order.refund_status)} / 退款金额 ¥${formatFen(order.refund_amount)}`,
      tone: order.refund_status === 3 ? "danger" : "warning"
    });
  } else if (order.refund_status > 0) {
    items.push({
      title: "退款流程",
      time: formatDateTime(order.updated_at),
      desc: refundStatusLabel(order.refund_status),
      tone: order.refund_status === 4 ? "danger" : "warning"
    });
  }

  if ([40, 50, 70].includes(order.status)) {
    items.push({
      title: `订单${orderStatusLabel(order.status)}`,
      time: formatDateTime(order.updated_at),
      desc: order.trade_no || "",
      tone: order.status === 70 ? "danger" : "muted"
    });
  }

  return items;
});

const formatFen = (value?: number | null) => ((value || 0) / 100).toFixed(2);

const formatDateTime = (value?: string | null) => {
  if (!value) return "-";
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return value;
  return date.toLocaleString("zh-CN", { hour12: false });
};

const orderStatusLabel = (status: number) => {
  switch (status) {
    case 0:
      return "待支付";
    case 10:
      return "已支付";
    case 20:
      return "已开通";
    case 30:
      return "已完成";
    case 40:
      return "已关闭";
    case 50:
      return "已取消";
    case 60:
      return "退款中";
    case 70:
      return "已退款";
    default:
      return String(status);
  }
};

const orderStatusColor = (status: number) => {
  switch (status) {
    case 0:
      return "gold";
    case 10:
      return "arcoblue";
    case 20:
      return "green";
    case 30:
      return "cyan";
    case 40:
      return "gray";
    case 50:
      return "orangered";
    case 60:
      return "purple";
    case 70:
      return "magenta";
    default:
      return "gray";
  }
};

const refundStatusLabel = (status: number) => {
  switch (status) {
    case 0:
      return "无";
    case 1:
      return "已申请";
    case 2:
      return "审核通过";
    case 3:
      return "退款成功";
    case 4:
      return "退款拒绝";
    default:
      return String(status);
  }
};

const refundStatusColor = (status: number) => {
  switch (status) {
    case 0:
      return "gray";
    case 1:
      return "gold";
    case 2:
      return "arcoblue";
    case 3:
      return "green";
    case 4:
      return "red";
    default:
      return "gray";
  }
};

const payMethodLabel = (method: number) => {
  switch (method) {
    case 1:
      return "微信支付";
    case 2:
      return "支付宝";
    case 3:
      return "银联支付";
    default:
      return String(method || "-");
  }
};

const payChannelLabel = (channel: string) => channel || "未标记";

const productTypeLabel = (type: number) => {
  switch (type) {
    case 1:
      return "文本";
    case 2:
      return "视频";
    case 3:
      return "音频";
    case 4:
      return "会员";
    default:
      return String(type);
  }
};

onMounted(async () => {
  await Promise.all([fetchList(), fetchChannels()]);
});
</script>

<style scoped lang="scss">
.order-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.page-hero {
  display: grid;
  grid-template-columns: minmax(0, 1.65fr) minmax(220px, 0.55fr);
  gap: 18px;
  padding: 24px 28px;
  border-radius: 16px;
  background:
    linear-gradient(135deg, rgba(22, 93, 255, 0.96), rgba(64, 128, 255, 0.86)),
    radial-gradient(circle at top right, rgba(255, 255, 255, 0.16), transparent 38%);
  color: #fff;
  box-shadow: 0 20px 38px rgb(22 93 255 / 16%);
}

.page-badge {
  display: inline-flex;
  padding: 6px 12px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.14);
  font-size: 12px;
}

.page-hero h1 {
  margin: 14px 0 10px;
  font-family: AliFangYuanTi, "PingFang SC", sans-serif;
  font-size: 32px;
  line-height: 1.2;
}

.page-hero p {
  max-width: 760px;
  line-height: 1.8;
  color: rgba(255, 255, 255, 0.86);
}

.page-hero-side {
  display: flex;
  flex-direction: column;
  justify-content: center;
  padding: 18px 20px;
  border-radius: 14px;
  background: rgba(255, 255, 255, 0.14);
  backdrop-filter: blur(10px);
}

.page-hero-label {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.74);
}

.page-hero-value {
  margin-top: 8px;
  font-size: 28px;
  font-weight: 700;
}

.panel-card,
.detail-card {
  border-radius: 16px;
  box-shadow: 0 12px 30px rgb(15 23 42 / 5%);
}

.filters {
  gap: 12px 8px;
}

.stack-cell {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.copy-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.stack-title {
  font-weight: 600;
  color: #1d2129;
}

.stack-sub {
  font-size: 12px;
  color: #86909c;
  word-break: break-all;
}

.money {
  color: #0f766e;
}

.detail-wrap {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.detail-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
}

.detail-summary-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 12px;
}

.detail-summary-card {
  padding: 16px 18px;
  border-radius: 16px;
  background: #f7f8fa;
  border: 1px solid #e5e6eb;
}

.detail-summary-card.success {
  background: linear-gradient(180deg, #effcf5 0%, #f9fffb 100%);
  border-color: #b7ebc6;
}

.detail-summary-card.info {
  background: linear-gradient(180deg, #eef3ff 0%, #f8fbff 100%);
  border-color: #bacefd;
}

.detail-summary-card.warning {
  background: linear-gradient(180deg, #fff7e8 0%, #fffdf7 100%);
  border-color: #ffd591;
}

.detail-summary-card.purple {
  background: linear-gradient(180deg, #f5f0ff 0%, #fbf8ff 100%);
  border-color: #d3adf7;
}

.detail-summary-card.muted {
  background: linear-gradient(180deg, #f2f3f5 0%, #fafafa 100%);
  border-color: #d9d9d9;
}

.detail-summary-label {
  font-size: 12px;
  color: #86909c;
}

.detail-summary-value {
  margin-top: 8px;
  font-size: 24px;
  font-weight: 700;
  color: #1d2129;
}

.detail-summary-sub {
  margin-top: 6px;
  font-size: 12px;
  line-height: 1.6;
  color: #4e5969;
  word-break: break-all;
}

.detail-copy-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  min-width: 0;
}

.timeline-list {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.timeline-item {
  display: grid;
  grid-template-columns: 16px minmax(0, 1fr);
  gap: 14px;
  align-items: flex-start;
}

.timeline-dot {
  position: relative;
  width: 14px;
  height: 14px;
  margin-top: 4px;
  border-radius: 999px;
  background: #c9cdd4;
}

.timeline-dot::after {
  content: "";
  position: absolute;
  top: 16px;
  left: 6px;
  width: 2px;
  height: calc(100% + 18px);
  background: #e5e6eb;
}

.timeline-item:last-child .timeline-dot::after {
  display: none;
}

.timeline-dot.info {
  background: #165dff;
}

.timeline-dot.success {
  background: #00b42a;
}

.timeline-dot.warning {
  background: #ff7d00;
}

.timeline-dot.danger {
  background: #f53f3f;
}

.timeline-dot.muted {
  background: #86909c;
}

.timeline-body {
  display: flex;
  flex-direction: column;
  gap: 4px;
  min-width: 0;
}

.timeline-title {
  font-size: 14px;
  font-weight: 600;
  color: #1d2129;
}

.timeline-time {
  font-size: 12px;
  color: #4e5969;
}

.timeline-desc {
  font-size: 12px;
  color: #86909c;
  line-height: 1.6;
  word-break: break-all;
}

@media (max-width: 1080px) {
  .page-hero,
  .detail-summary-grid,
  .detail-grid {
    grid-template-columns: 1fr;
  }
}
</style>
