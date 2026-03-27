# Goodie - Food Delivery Platform

## 1. Tổng quan hệ thống

**Goodie** là nền tảng giao đồ ăn (Food Delivery Platform) được thiết kế theo kiến trúc **Microservices**, áp dụng các công nghệ hiện đại phù hợp cho Senior Backend Engineer.

---

## 2. Kiến trúc hệ thống (System Architecture)

```
                                    ┌─────────────────────────────────────────┐
                                    │              API Gateway                │
                                    │         (Kong / Nginx / Traefik)        │
                                    └──────┬──────────┬──────────┬────────────┘
                                           │          │          │
                              ┌────────────┘          │          └────────────┐
                              ▼                       ▼                       ▼
                    ┌──────────────────┐   ┌──────────────────┐   ┌──────────────────┐
                    │  Admin Service   │   │Merchant Service  │   │  Client Service  │
                    │   (HTTP/gRPC)    │   │   (HTTP/gRPC)    │   │   (HTTP/gRPC)    │
                    │   :8081          │   │   :8082          │   │   :8083          │
                    └──────┬───────────┘   └──────┬───────────┘   └──────┬───────────┘
                           │                      │                      │
                           │         gRPC (inter-service)                │
                           │◄─────────────────────┼──────────────────────►
                           │                      │                      │
              ┌────────────┴──────────────────────┴──────────────────────┴────────────┐
              │                                                                       │
              ▼                                                                       ▼
    ┌──────────────────┐                                                   ┌──────────────────┐
    │      Kafka       │                                                   │    WebSocket      │
    │  Message Broker  │                                                   │    Service        │
    │   :9092          │                                                   │    :8084          │
    └──────┬───────────┘                                                   └──────────────────┘
           │                                                                        ▲
           ▼                                                                        │
    ┌──────────────────┐                                                            │
    │ Consumer Service │────────────────────────────────────────────────────────────┘
    │  (Background)    │         (push realtime events via WebSocket)
    └──────┬───────────┘
           │
    ┌──────┴───────────────────────────────────────────┐
    │                  Data Layer                       │
    │                                                   │
    │  ┌──────────┐  ┌──────────┐  ┌──────────────┐   │
    │  │PostgreSQL│  │ MongoDB  │  │    Redis      │   │
    │  │  :5432   │  │  :27017  │  │    :6379      │   │
    │  │          │  │          │  │               │   │
    │  │- Users   │  │- Logs    │  │- Session      │   │
    │  │- Orders  │  │- Reviews │  │- Cache        │   │
    │  │- Menus   │  │- Chat    │  │- Rate Limit   │   │
    │  │- Payments│  │- Notif.  │  │- Geo (driver) │   │
    │  └──────────┘  └──────────┘  └──────────────┘   │
    └──────────────────────────────────────────────────┘
           │
           ▼
    ┌──────────────────────────────────────────────────┐
    │              ELK Stack                            │
    │  ┌──────────┐ ┌──────────┐ ┌──────────────────┐ │
    │  │Elastic   │ │Logstash  │ │    Kibana        │ │
    │  │Search    │ │          │ │                   │ │
    │  │:9200     │ │:5044     │ │   :5601          │ │
    │  │          │ │          │ │                   │ │
    │  │- Search  │ │- Log     │ │ - Dashboard      │ │
    │  │- Analyze │ │  Pipeline│ │ - Monitoring      │ │
    │  └──────────┘ └──────────┘ └──────────────────┘ │
    └──────────────────────────────────────────────────┘
```

---

## 3. Tech Stack

| Layer | Technology | Mục đích |
|-------|-----------|----------|
| **Language** | Go (Golang) | Backend services |
| **HTTP Framework** | Gin / Echo / Fiber | REST API |
| **gRPC** | google.golang.org/grpc | Inter-service communication |
| **Message Broker** | Apache Kafka | Event streaming, async processing |
| **Relational DB** | PostgreSQL | Structured data (users, orders, menus, payments) |
| **Document DB** | MongoDB | Unstructured data (reviews, chat, notifications, activity logs) |
| **Cache** | Redis | Caching, session, rate limiting, geo-tracking drivers |
| **Search & Log** | Elasticsearch + Logstash + Kibana (ELK) | Centralized logging, full-text search (restaurants, menus) |
| **Realtime** | WebSocket (gorilla/websocket hoặc nhk/websocket) | Order tracking, notifications, chat |
| **Container** | Docker + Docker Compose | Containerization |
| **Orchestration** | Kubernetes (K8s) | Production deployment, scaling, service mesh |
| **Auth** | JWT + OAuth2 | Authentication & Authorization |
| **API Gateway** | Kong / Traefik / Nginx | Routing, rate limiting, load balancing |
| **CI/CD** | GitHub Actions | Automated testing, build, deploy |
| **Monitoring** | Prometheus + Grafana | Metrics, alerting |
| **Tracing** | Jaeger / OpenTelemetry | Distributed tracing |
| **Migration** | golang-migrate | Database migrations |

---

## 4. Chi tiết từng Service

### 4.1 Admin Service (`:8081`)

Quản trị toàn bộ hệ thống.

| Module | Chức năng |
|--------|-----------|
| **User Management** | CRUD users, ban/unban, assign roles (admin/merchant/client/driver) |
| **Merchant Management** | Duyệt đăng ký merchant, suspend/activate merchant |
| **Order Management** | Xem tất cả orders, can thiệp xử lý dispute, refund |
| **Category Management** | CRUD danh mục món ăn (cuisine types) |
| **Promotion Management** | Tạo/quản lý voucher, discount campaigns toàn hệ thống |
| **Commission Config** | Cấu hình % hoa hồng cho từng merchant/tier |
| **Dashboard & Reports** | Thống kê doanh thu, số đơn, merchant performance, driver performance |
| **System Config** | Cấu hình delivery fee, min order, service areas |
| **Notification Broadcast** | Gửi thông báo hệ thống cho merchant/client/driver |

### 4.2 Merchant Service (`:8082`)

Dành cho chủ nhà hàng/quán ăn.

| Module | Chức năng |
|--------|-----------|
| **Auth** | Đăng ký merchant, login, profile management |
| **Store Management** | CRUD thông tin cửa hàng (tên, địa chỉ, giờ mở/đóng, ảnh) |
| **Menu Management** | CRUD món ăn, set giá, ảnh, mô tả, topping/options |
| **Category Management** | Phân loại menu theo category (khai vị, món chính, đồ uống...) |
| **Order Management** | Nhận đơn mới (realtime), xác nhận/từ chối đơn, cập nhật trạng thái chuẩn bị |
| **Inventory/Availability** | Bật/tắt món hết hàng, tạm ngưng cửa hàng |
| **Revenue & Reports** | Xem doanh thu theo ngày/tuần/tháng, lịch sử đơn hàng |
| **Promotion** | Tạo khuyến mãi riêng cho store (giảm giá, combo, freeship) |
| **Rating & Reviews** | Xem và phản hồi đánh giá từ khách hàng |

### 4.3 Client Service (`:8083`)

Dành cho khách hàng đặt đồ ăn.

| Module | Chức năng |
|--------|-----------|
| **Auth** | Đăng ký, login (email/phone/social), OTP verification |
| **Profile** | Quản lý thông tin cá nhân, địa chỉ giao hàng (nhiều địa chỉ) |
| **Restaurant Discovery** | Tìm kiếm nhà hàng (theo tên, cuisine, khoảng cách, rating) - **Elasticsearch** |
| **Menu Browsing** | Xem menu, filter, sort, xem chi tiết món |
| **Cart** | Thêm/xóa/sửa giỏ hàng (lưu trên **Redis** cho tốc độ) |
| **Order** | Đặt đơn, chọn payment method, apply voucher |
| **Order Tracking** | Theo dõi realtime trạng thái đơn hàng qua **WebSocket** |
| **Payment** | Thanh toán COD / E-wallet / Credit Card (tích hợp payment gateway) |
| **Rating & Review** | Đánh giá nhà hàng & driver sau khi nhận đồ |
| **Notification** | Nhận thông báo trạng thái đơn, khuyến mãi |
| **Favorites** | Lưu nhà hàng/món yêu thích |
| **Order History** | Xem lịch sử đặt hàng, re-order |

### 4.4 Consumer Service (Background Workers)

Xử lý các events bất đồng bộ từ Kafka.

| Consumer Group | Topic | Chức năng |
|---------------|-------|-----------|
| **order-processor** | `order.created` | Validate đơn hàng, tính phí, gửi cho merchant |
| **payment-processor** | `payment.initiated` | Xử lý thanh toán, callback từ payment gateway |
| **notification-sender** | `notification.send` | Gửi push notification, SMS, email |
| **driver-matcher** | `order.confirmed` | Tìm driver gần nhất, assign đơn |
| **analytics-collector** | `order.*`, `user.*` | Thu thập data cho analytics/reporting |
| **log-shipper** | `service.logs` | Ship logs tới ELK stack |
| **search-indexer** | `merchant.updated`, `menu.updated` | Cập nhật Elasticsearch index khi data thay đổi |
| **review-processor** | `review.created` | Tính toán lại average rating |

### 4.5 WebSocket Service (`:8084`)

Realtime communication layer.

| Channel | Chức năng |
|---------|-----------|
| **order:{orderId}** | Client & Merchant theo dõi trạng thái đơn hàng realtime |
| **driver:{driverId}** | Cập nhật vị trí driver lên map |
| **merchant:{merchantId}** | Merchant nhận đơn hàng mới realtime |
| **notification:{userId}** | Push notification tới user cụ thể |
| **chat:{roomId}** | Chat giữa client - driver (nếu có) |

---

## 5. gRPC — Inter-Service Communication

gRPC được sử dụng cho giao tiếp **đồng bộ** giữa các services nội bộ, nơi cần **low latency** và **type-safe**.

| gRPC Service | Method | Caller → Callee | Mục đích |
|-------------|--------|-----------------|----------|
| **UserService** | `GetUserByID` | Client/Merchant/Admin → User internal | Lấy thông tin user (dùng chung) |
| **UserService** | `ValidateToken` | All services → Auth internal | Xác thực JWT token giữa các service |
| **MerchantService** | `GetStoreInfo` | Client → Merchant | Lấy thông tin cửa hàng khi hiển thị |
| **MerchantService** | `GetMenuItems` | Client → Merchant | Lấy menu items cho giỏ hàng/đơn hàng |
| **MerchantService** | `UpdateOrderStatus` | Consumer → Merchant | Cập nhật trạng thái đơn ở phía merchant |
| **OrderService** | `CreateOrder` | Client → Order internal | Tạo đơn hàng mới |
| **OrderService** | `GetOrderDetail` | All → Order internal | Lấy chi tiết đơn hàng |
| **PaymentService** | `ProcessPayment` | Order → Payment internal | Xử lý thanh toán khi đặt đơn |
| **PaymentService** | `RefundPayment` | Admin → Payment internal | Hoàn tiền |
| **NotificationService** | `SendNotification` | Any → Notification internal | Gửi notification nội bộ (sync, urgent) |
| **DriverService** | `FindNearestDriver` | Consumer → Driver internal | Tìm driver gần nhất (dùng Redis Geo) |
| **DriverService** | `GetDriverLocation` | Client → Driver internal | Lấy vị trí driver realtime |
| **AnalyticsService** | `GetMerchantStats` | Admin/Merchant → Analytics | Lấy thống kê hiệu suất |

**Tại sao dùng gRPC thay vì REST cho nội bộ:**
- **Performance**: Binary protocol (protobuf) nhanh hơn JSON 5-10x
- **Type Safety**: Proto file là contract rõ ràng giữa services
- **Streaming**: Hỗ trợ bi-directional streaming (dùng cho driver location updates)
- **Code Generation**: Auto-gen client/server code từ `.proto` files

---

## 6. Kafka Topics & Event Flow

```
┌─────────────┐     order.created      ┌─────────────────┐
│   Client    │ ──────────────────────► │  order-processor │
│   Service   │                         │   (consumer)     │
└─────────────┘                         └────────┬─────────┘
                                                 │
                                    order.validated / order.failed
                                                 │
                    ┌────────────────────────────┼────────────────────────┐
                    ▼                            ▼                        ▼
          ┌─────────────────┐         ┌──────────────────┐    ┌──────────────────┐
          │ payment-processor│        │  driver-matcher   │    │notification-sender│
          │  (consumer)      │        │   (consumer)      │    │   (consumer)      │
          └─────────────────┘         └──────────────────┘    └──────────────────┘
                    │                            │                        │
           payment.completed            driver.assigned          notification.sent
                    │                            │
                    ▼                            ▼
          ┌─────────────────┐         ┌──────────────────┐
          │  WebSocket push │         │  WebSocket push  │
          │  (order status) │         │  (driver on map) │
          └─────────────────┘         └──────────────────┘
```

### Kafka Topics

| Topic | Producer | Consumer | Event Type |
|-------|----------|----------|------------|
| `order.created` | Client Service | order-processor | Đơn hàng mới |
| `order.confirmed` | Merchant Service | driver-matcher, notification-sender | Merchant xác nhận đơn |
| `order.preparing` | Merchant Service | notification-sender | Đang chuẩn bị |
| `order.ready` | Merchant Service | notification-sender, driver-matcher | Đồ ăn sẵn sàng |
| `order.picked_up` | Driver (via Client API) | notification-sender | Driver đã lấy hàng |
| `order.delivered` | Driver (via Client API) | payment-processor, notification-sender, analytics | Giao thành công |
| `order.cancelled` | Client/Merchant/Admin | payment-processor, notification-sender | Đơn bị hủy |
| `payment.initiated` | Order Service | payment-processor | Bắt đầu thanh toán |
| `payment.completed` | Payment Gateway Callback | order-processor, notification-sender | Thanh toán thành công |
| `payment.failed` | Payment Gateway Callback | order-processor, notification-sender | Thanh toán thất bại |
| `merchant.updated` | Merchant Service | search-indexer | Thông tin merchant thay đổi |
| `menu.updated` | Merchant Service | search-indexer | Menu thay đổi |
| `review.created` | Client Service | review-processor, notification-sender | Đánh giá mới |
| `notification.send` | Any Service | notification-sender | Gửi notification |
| `driver.location` | Driver App | WebSocket Service | Cập nhật vị trí driver |
| `service.logs` | All Services | log-shipper (→ ELK) | Application logs |

---

## 7. Database Schema Design

### 7.1 PostgreSQL (Structured/Transactional Data)

```
┌──────────────────┐    ┌──────────────────┐    ┌──────────────────┐
│      users       │    │     stores       │    │   menu_items     │
├──────────────────┤    ├──────────────────┤    ├──────────────────┤
│ id (PK)          │    │ id (PK)          │    │ id (PK)          │
│ email            │    │ merchant_id (FK) │───►│ store_id (FK)    │
│ phone            │    │ name             │    │ name             │
│ password_hash    │    │ address          │    │ description      │
│ full_name        │    │ lat, lng         │    │ price            │
│ role (enum)      │    │ phone            │    │ image_url        │
│ avatar_url       │    │ image_url        │    │ category_id (FK) │
│ is_verified      │    │ opening_hours    │    │ is_available     │
│ status           │    │ is_active        │    │ sort_order       │
│ created_at       │    │ avg_rating       │    │ created_at       │
│ updated_at       │    │ total_reviews    │    │ updated_at       │
└──────────────────┘    │ commission_rate  │    └──────────────────┘
         │              │ created_at       │
         │              └──────────────────┘    ┌──────────────────┐
         │                                      │  item_options    │
         │              ┌──────────────────┐    ├──────────────────┤
         │              │   categories     │    │ id (PK)          │
         │              ├──────────────────┤    │ menu_item_id(FK) │
         │              │ id (PK)          │    │ name             │
         │              │ name             │    │ price_extra      │
         │              │ slug             │    │ is_required      │
         │              │ icon_url         │    │ max_select       │
         │              │ sort_order       │    └──────────────────┘
         │              └──────────────────┘
         │
         ▼
┌──────────────────┐    ┌──────────────────┐    ┌──────────────────┐
│     orders       │    │   order_items    │    │    payments      │
├──────────────────┤    ├──────────────────┤    ├──────────────────┤
│ id (PK)          │    │ id (PK)          │    │ id (PK)          │
│ user_id (FK)     │    │ order_id (FK)    │    │ order_id (FK)    │
│ store_id (FK)    │    │ menu_item_id(FK) │    │ amount           │
│ driver_id (FK)   │    │ quantity         │    │ method (enum)    │
│ status (enum)    │    │ unit_price       │    │ status (enum)    │
│ subtotal         │    │ options (jsonb)  │    │ transaction_id   │
│ delivery_fee     │    │ note             │    │ gateway_response │
│ discount         │    │ subtotal         │    │ paid_at          │
│ total            │    └──────────────────┘    │ created_at       │
│ delivery_address │                             └──────────────────┘
│ delivery_lat     │    ┌──────────────────┐
│ delivery_lng     │    │   promotions     │    ┌──────────────────┐
│ note             │    ├──────────────────┤    │  user_addresses  │
│ voucher_code     │    │ id (PK)          │    ├──────────────────┤
│ estimated_time   │    │ store_id (FK)    │    │ id (PK)          │
│ created_at       │    │ code             │    │ user_id (FK)     │
│ updated_at       │    │ type (enum)      │    │ label            │
└──────────────────┘    │ value            │    │ address          │
                        │ min_order        │    │ lat, lng         │
                        │ max_discount     │    │ is_default       │
                        │ start_date       │    │ created_at       │
                        │ end_date         │    └──────────────────┘
                        │ usage_limit      │
                        │ used_count       │
                        └──────────────────┘
```

**Order Status Flow:**
```
PENDING → CONFIRMED → PREPARING → READY → PICKED_UP → DELIVERED
                 ↘                                    ↗
                  → CANCELLED ←──────────────────────
```

### 7.2 MongoDB (Unstructured/Flexible Data)

```javascript
// Collection: reviews
{
  _id: ObjectId,
  order_id: "uuid",
  user_id: "uuid",
  store_id: "uuid",
  driver_id: "uuid",
  store_rating: 4.5,
  driver_rating: 5.0,
  comment: "Đồ ăn ngon, giao nhanh!",
  images: ["url1", "url2"],
  merchant_reply: "Cảm ơn bạn!",
  replied_at: ISODate,
  created_at: ISODate
}

// Collection: notifications
{
  _id: ObjectId,
  user_id: "uuid",
  type: "ORDER_STATUS" | "PROMOTION" | "SYSTEM",
  title: "Đơn hàng đã được xác nhận",
  body: "Nhà hàng XYZ đang chuẩn bị đơn hàng #123",
  data: { order_id: "uuid", status: "CONFIRMED" },
  is_read: false,
  created_at: ISODate
}

// Collection: chat_messages
{
  _id: ObjectId,
  room_id: "order_uuid",
  sender_id: "uuid",
  sender_type: "CLIENT" | "DRIVER",
  message: "Tôi đang ở cổng chính",
  type: "TEXT" | "IMAGE" | "LOCATION",
  created_at: ISODate
}

// Collection: activity_logs
{
  _id: ObjectId,
  actor_id: "uuid",
  actor_type: "ADMIN" | "MERCHANT" | "CLIENT" | "SYSTEM",
  action: "ORDER_CREATED",
  resource_type: "ORDER",
  resource_id: "uuid",
  metadata: { ... },
  ip_address: "x.x.x.x",
  created_at: ISODate
}
```

### 7.3 Redis Data Structures

| Key Pattern | Type | TTL | Mục đích |
|-------------|------|-----|----------|
| `session:{userId}` | String (JWT) | 24h | User session |
| `cart:{userId}` | Hash | 7d | Giỏ hàng |
| `otp:{phone}` | String | 5m | OTP verification |
| `rate_limit:{ip}:{endpoint}` | String (counter) | 1m | Rate limiting |
| `driver:location:{driverId}` | Geo | - | Vị trí driver (GEOADD, GEORADIUS) |
| `store:online:{storeId}` | String | - | Store đang mở/đóng |
| `cache:store:{storeId}` | Hash | 10m | Cache thông tin store |
| `cache:menu:{storeId}` | String (JSON) | 5m | Cache menu |
| `order:lock:{orderId}` | String | 30s | Distributed lock khi xử lý đơn |
| `search:popular` | Sorted Set | 1h | Popular search terms |

---

## 8. Elasticsearch Indices

### Restaurant Search Index
```json
{
  "index": "restaurants",
  "mappings": {
    "properties": {
      "store_id": { "type": "keyword" },
      "name": { "type": "text", "analyzer": "vietnamese" },
      "description": { "type": "text", "analyzer": "vietnamese" },
      "cuisine_types": { "type": "keyword" },
      "location": { "type": "geo_point" },
      "avg_rating": { "type": "float" },
      "total_reviews": { "type": "integer" },
      "price_range": { "type": "integer" },
      "is_active": { "type": "boolean" },
      "opening_hours": { "type": "object" },
      "menu_items": {
        "type": "nested",
        "properties": {
          "name": { "type": "text", "analyzer": "vietnamese" },
          "price": { "type": "float" },
          "category": { "type": "keyword" }
        }
      }
    }
  }
}
```

**Search Features:**
- Full-text search nhà hàng/món ăn (hỗ trợ tiếng Việt)
- Geo-distance search (tìm quán gần nhất)
- Filter theo cuisine, rating, price range
- Autocomplete / suggest
- Boosting (ưu tiên quán rating cao, gần user)

### Application Logs Index
```json
{
  "index": "app-logs-{YYYY.MM.DD}",
  "mappings": {
    "properties": {
      "timestamp": { "type": "date" },
      "level": { "type": "keyword" },
      "service": { "type": "keyword" },
      "trace_id": { "type": "keyword" },
      "message": { "type": "text" },
      "metadata": { "type": "object" }
    }
  }
}
```

---

## 9. Yêu cầu chức năng (Functional Requirements)

### FR-01: Authentication & Authorization
- [ ] Đăng ký tài khoản (email, phone, social login)
- [ ] Đăng nhập / Đăng xuất (JWT + Refresh Token)
- [ ] Xác thực OTP qua SMS/Email
- [ ] Phân quyền theo role (Admin, Merchant, Client, Driver)
- [ ] OAuth2 social login (Google, Facebook)
- [ ] Forgot/Reset password

### FR-02: Restaurant & Menu Management
- [ ] CRUD cửa hàng (thông tin, giờ mở cửa, hình ảnh)
- [ ] CRUD menu items (tên, giá, mô tả, ảnh, topping/options)
- [ ] Phân loại menu theo category
- [ ] Bật/tắt trạng thái món (hết hàng)
- [ ] Tạm ngưng/mở lại cửa hàng

### FR-03: Search & Discovery
- [ ] Tìm kiếm nhà hàng theo tên, loại đồ ăn (full-text search)
- [ ] Tìm kiếm theo vị trí (geo-distance)
- [ ] Filter: rating, khoảng cách, loại cuisine, giá
- [ ] Sort: gần nhất, rating cao nhất, phổ biến nhất
- [ ] Autocomplete gợi ý
- [ ] Hiển thị nhà hàng đề xuất (recommendation)

### FR-04: Cart & Ordering
- [ ] Thêm/xóa/cập nhật món trong giỏ hàng
- [ ] Chọn topping/options cho món
- [ ] Apply voucher/promotion code
- [ ] Tính tổng tiền (subtotal + delivery fee - discount)
- [ ] Chọn địa chỉ giao hàng
- [ ] Ghi chú đơn hàng
- [ ] Đặt đơn hàng
- [ ] Hủy đơn hàng (trong thời gian cho phép)

### FR-05: Order Processing & Tracking
- [ ] Merchant nhận đơn mới (realtime)
- [ ] Merchant xác nhận/từ chối đơn
- [ ] Cập nhật trạng thái đơn: Pending → Confirmed → Preparing → Ready → Picked Up → Delivered
- [ ] Client theo dõi trạng thái đơn hàng realtime (WebSocket)
- [ ] Theo dõi vị trí driver trên map (WebSocket + Redis Geo)
- [ ] Ước tính thời gian giao hàng
- [ ] Auto-assign driver gần nhất

### FR-06: Payment
- [ ] Thanh toán COD (tiền mặt)
- [ ] Thanh toán E-wallet (MoMo, ZaloPay, VNPay)
- [ ] Thanh toán Credit/Debit Card
- [ ] Xử lý refund khi hủy đơn
- [ ] Lịch sử giao dịch

### FR-07: Rating & Review
- [ ] Đánh giá nhà hàng (1-5 sao + comment + ảnh)
- [ ] Đánh giá driver (1-5 sao)
- [ ] Merchant phản hồi review
- [ ] Tính toán average rating tự động

### FR-08: Notification
- [ ] Push notification trạng thái đơn hàng
- [ ] Notification khuyến mãi
- [ ] Thông báo hệ thống từ admin
- [ ] Notification đơn hàng mới cho merchant (realtime)
- [ ] Đánh dấu đã đọc/chưa đọc

### FR-09: Admin Dashboard
- [ ] Dashboard thống kê tổng quan (doanh thu, số đơn, user mới)
- [ ] Quản lý users (ban/unban, xem thông tin)
- [ ] Duyệt đăng ký merchant
- [ ] Quản lý danh mục hệ thống
- [ ] Quản lý promotion/voucher toàn hệ thống
- [ ] Cấu hình hoa hồng, phí giao hàng
- [ ] Xem logs & audit trail

### FR-10: Merchant Dashboard
- [ ] Thống kê doanh thu theo ngày/tuần/tháng
- [ ] Quản lý đơn hàng (danh sách, filter theo status)
- [ ] Xem lịch sử đơn hàng
- [ ] Quản lý khuyến mãi riêng cửa hàng
- [ ] Xem & phản hồi reviews

### FR-11: Chat
- [ ] Chat realtime giữa client và driver (trong quá trình giao hàng)
- [ ] Gửi text, image, location
- [ ] Lịch sử chat theo đơn hàng

---

## 10. Yêu cầu phi chức năng (Non-Functional Requirements)

### NFR-01: Performance
- API response time P95 < **200ms** (đọc), < **500ms** (ghi)
- WebSocket message latency < **100ms**
- Search response < **300ms**
- Hệ thống xử lý **1000+ concurrent orders** tại peak time
- Database query time P95 < **50ms**

### NFR-02: Scalability
- Horizontal scaling cho mọi service (stateless design)
- Kafka partitioning cho throughput cao
- Redis Cluster cho high availability cache
- PostgreSQL read replicas cho read-heavy workloads
- Kubernetes HPA (Horizontal Pod Autoscaler) cho auto-scaling
- Target: scale tới **100K users**, **10K merchants**, **50K orders/day**

### NFR-03: Availability & Reliability
- Uptime target: **99.9%** (< 8.76h downtime/year)
- Zero-downtime deployment (rolling updates, blue-green)
- Circuit breaker pattern cho external service calls
- Retry mechanism với exponential backoff
- Graceful degradation khi downstream service fail
- Database backup daily, point-in-time recovery

### NFR-04: Security
- JWT token với short TTL (15m access, 7d refresh)
- Bcrypt password hashing (cost factor ≥ 12)
- Rate limiting per IP & per user
- Input validation & sanitization (prevent SQL injection, XSS)
- CORS configuration
- HTTPS everywhere (TLS 1.3)
- Sensitive data encryption at rest
- RBAC (Role-Based Access Control)
- API key rotation cho third-party integrations
- Audit logging cho mọi admin actions
- Helmet headers, CSRF protection

### NFR-05: Observability
- **Structured logging** (JSON format) → ELK Stack
- **Distributed tracing** (OpenTelemetry / Jaeger) với trace_id xuyên suốt request
- **Metrics** (Prometheus + Grafana): request rate, error rate, latency, resource usage
- **Alerting**: PagerDuty / Slack alerts khi error rate > threshold
- **Health checks**: /health endpoint cho mọi service
- **Log levels**: DEBUG, INFO, WARN, ERROR với correlation ID

### NFR-06: Maintainability
- Clean Architecture / Hexagonal Architecture cho mỗi service
- Unit test coverage ≥ **80%**
- Integration tests cho critical paths
- API documentation (Swagger/OpenAPI 3.0)
- gRPC proto documentation
- Database migration versioning (golang-migrate)
- Code linting (golangci-lint) + pre-commit hooks
- Conventional commits

### NFR-07: DevOps & Infrastructure
- Docker multi-stage builds (minimize image size)
- Docker Compose cho local development
- Kubernetes manifests (Deployment, Service, ConfigMap, Secret, Ingress, HPA)
- Helm charts cho deployment
- CI/CD pipeline (GitHub Actions): lint → test → build → push → deploy
- Environment separation: dev / staging / production
- Infrastructure as Code (Terraform hoặc Pulumi - optional)
- Secret management (K8s Secrets / Vault)

### NFR-08: Data Consistency
- ACID transactions cho payment & order operations (PostgreSQL)
- Eventual consistency cho non-critical data (via Kafka events)
- Idempotent API endpoints (prevent double orders/payments)
- Distributed locking (Redis) cho critical sections
- Saga pattern cho distributed transactions (order → payment → driver assignment)
- Outbox pattern cho reliable event publishing

### NFR-09: Resilience
- Circuit Breaker (prevent cascade failures)
- Bulkhead pattern (isolate failures)
- Timeout configuration cho mọi external calls
- Dead Letter Queue (DLQ) cho failed Kafka messages
- Graceful shutdown (drain connections before stop)
- Health check + readiness probe cho K8s

### NFR-10: API Standards
- RESTful API design (proper HTTP methods, status codes)
- API versioning (URL path: `/api/v1/...`)
- Pagination (cursor-based cho large datasets)
- Consistent error response format
- Request/Response validation (middleware)
- HATEOAS (optional, cho API maturity)

---

## 11. Project Structure

```
goodie-api/
├── docker-compose.yml              # Local development environment
├── docker-compose.infra.yml        # Infrastructure services (DB, Redis, Kafka, ELK)
├── Makefile                        # Common commands
├── README.md
│
├── proto/                          # Shared gRPC proto definitions
│   ├── user/
│   │   └── user.proto
│   ├── merchant/
│   │   └── merchant.proto
│   ├── order/
│   │   └── order.proto
│   ├── payment/
│   │   └── payment.proto
│   ├── notification/
│   │   └── notification.proto
│   └── driver/
│       └── driver.proto
│
├── services/
│   ├── admin/                      # Admin Service
│   │   ├── Dockerfile
│   │   ├── cmd/
│   │   │   └── main.go
│   │   ├── internal/
│   │   │   ├── domain/             # Entities, Value Objects
│   │   │   ├── usecase/            # Business logic
│   │   │   ├── repository/         # Data access interfaces
│   │   │   ├── delivery/
│   │   │   │   ├── http/           # REST handlers
│   │   │   │   └── grpc/           # gRPC handlers
│   │   │   └── infrastructure/     # DB, Redis, Kafka implementations
│   │   ├── migrations/
│   │   ├── config/
│   │   └── go.mod
│   │
│   ├── merchant/                   # Merchant Service (same structure)
│   │   ├── Dockerfile
│   │   ├── cmd/
│   │   ├── internal/
│   │   ├── migrations/
│   │   └── go.mod
│   │
│   ├── client/                     # Client Service (same structure)
│   │   ├── Dockerfile
│   │   ├── cmd/
│   │   ├── internal/
│   │   ├── migrations/
│   │   └── go.mod
│   │
│   ├── consumer/                   # Kafka Consumer Service
│   │   ├── Dockerfile
│   │   ├── cmd/
│   │   ├── internal/
│   │   │   ├── handlers/           # Consumer handlers per topic
│   │   │   └── processors/         # Business logic
│   │   └── go.mod
│   │
│   └── websocket/                  # WebSocket Service
│       ├── Dockerfile
│       ├── cmd/
│       ├── internal/
│       │   ├── hub/                # Connection manager
│       │   ├── handlers/           # Message handlers
│       │   └── middleware/
│       └── go.mod
│
├── pkg/                            # Shared packages across services
│   ├── auth/                       # JWT utilities
│   ├── logger/                     # Structured logging (zap/zerolog)
│   ├── middleware/                  # Common HTTP middleware
│   ├── kafka/                      # Kafka producer/consumer wrapper
│   ├── redis/                      # Redis client wrapper
│   ├── postgres/                   # PostgreSQL connection & helpers
│   ├── mongo/                      # MongoDB connection & helpers
│   ├── elasticsearch/              # ES client wrapper
│   ├── validator/                  # Request validation
│   ├── errors/                     # Custom error types
│   ├── pagination/                 # Cursor-based pagination
│   └── response/                   # Standard API response format
│
├── deployments/
│   ├── k8s/                        # Kubernetes manifests
│   │   ├── namespace.yaml
│   │   ├── admin/
│   │   │   ├── deployment.yaml
│   │   │   ├── service.yaml
│   │   │   ├── hpa.yaml
│   │   │   └── configmap.yaml
│   │   ├── merchant/
│   │   ├── client/
│   │   ├── consumer/
│   │   ├── websocket/
│   │   ├── ingress.yaml
│   │   └── secrets.yaml
│   └── helm/                       # Helm charts (optional)
│       └── goodie/
│
├── scripts/                        # Helper scripts
│   ├── migrate.sh
│   ├── seed.sh
│   └── proto-gen.sh                # Generate gRPC code from proto
│
├── configs/                        # Config files
│   ├── logstash/
│   │   └── pipeline.conf
│   ├── kibana/
│   ├── prometheus/
│   │   └── prometheus.yml
│   └── grafana/
│       └── dashboards/
│
└── docs/                           # Documentation
    ├── api/                        # Swagger/OpenAPI specs
    ├── architecture/               # Architecture Decision Records (ADR)
    └── diagrams/
```

---

## 12. Cách chạy Local Development

```bash
# 1. Start infrastructure
docker-compose -f docker-compose.infra.yml up -d

# 2. Run database migrations
make migrate-up

# 3. Seed initial data
make seed

# 4. Generate gRPC code
make proto-gen

# 5. Start all services
make run-all

# Hoặc start từng service
make run-admin
make run-merchant
make run-client
make run-consumer
make run-websocket
```

---

## 13. Roadmap triển khai

### Phase 1 — Foundation (Week 1-2)
- [x] Project structure setup
- [ ] Docker Compose infrastructure (PostgreSQL, Redis, Kafka, MongoDB, ELK)
- [ ] Shared packages (pkg/*)
- [ ] Auth module (JWT, RBAC)
- [ ] Database migrations
- [ ] gRPC proto definitions & code generation
- [ ] Structured logging → ELK pipeline

### Phase 2 — Core Services (Week 3-5)
- [ ] Admin Service: User & Merchant management
- [ ] Merchant Service: Store & Menu CRUD, Order management
- [ ] Client Service: Auth, Search, Cart, Order, Payment
- [ ] Kafka producers & consumers setup
- [ ] WebSocket: Order tracking & notifications

### Phase 3 — Advanced Features (Week 6-7)
- [ ] Elasticsearch: Restaurant & menu search
- [ ] Driver matching & location tracking (Redis Geo)
- [ ] Rating & Review system
- [ ] Chat (WebSocket + MongoDB)
- [ ] Promotion & voucher system

### Phase 4 — Production Ready (Week 8-9)
- [ ] Kubernetes deployment manifests
- [ ] CI/CD pipeline (GitHub Actions)
- [ ] Prometheus + Grafana monitoring
- [ ] Distributed tracing (Jaeger/OpenTelemetry)
- [ ] Load testing (k6/Locust)
- [ ] Security hardening
- [ ] API documentation (Swagger)

---

## 14. Senior Backend — Những điểm nâng cao cần thực hành

| Skill | Áp dụng trong project |
|-------|----------------------|
| **Clean Architecture** | Mỗi service chia domain/usecase/repository/delivery layers |
| **CQRS** | Tách read/write model cho Order service (PostgreSQL write, Elasticsearch read) |
| **Saga Pattern** | Order → Payment → Driver Assignment (distributed transaction) |
| **Outbox Pattern** | Đảm bảo event publish reliable (write to outbox table + poll to Kafka) |
| **Circuit Breaker** | Khi payment gateway hoặc external service down |
| **Rate Limiting** | Token bucket / Sliding window (Redis-based) |
| **Distributed Locking** | Redis SETNX cho prevent double order processing |
| **Event Sourcing** | Order history (optional, lưu mọi state change) |
| **Database Sharding** | Shard orders table theo date range (khi data lớn) |
| **Connection Pooling** | pgxpool cho PostgreSQL |
| **Graceful Shutdown** | Handle SIGTERM, drain Kafka consumers, close DB connections |
| **Idempotency** | Idempotency key cho payment APIs |
| **Feature Flags** | Toggle features without deployment |
| **API Gateway Patterns** | Authentication, rate limiting, request transformation |
| **Observability** | Structured logs + metrics + traces (3 pillars) |

---

> **Note**: Project này được thiết kế như một hệ thống thực tế ở production-grade, phù hợp cho việc thực hành và xây dựng portfolio cho Senior Backend Engineer.
