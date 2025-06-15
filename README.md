# Coupon System 優惠券系統

這是一個使用 Golang + Gin + MySQL + Redis 打造的簡易優惠券系統，提供以下功能：

- 使用者可領取多種優惠券（如滿額折抵、折扣券），每張券有數量限制。
- 使用者可於有效期限內使用優惠券，系統會驗證可用性（如過期或已使用）。
- 提供查詢 API 顯示使用者所有優惠券狀態（未使用、已使用、已過期）。

---

## 📦 專案架構簡介

```
coupon-system/
├── cmd/
│   └── seed/            # 種子資料程式
├── config/              # DB設定
├── controllers/         # API 控制器
├── models/              # GORM 模型
├── routes/              # 路由定義
├── Dockerfile           # Golang 應用容器
├── docker-compose.yml   # 一鍵啟動 MySQL、Redis、App、Nginx
└── README.md            # 專案說明文件
```

---

## 🚀 如何啟動專案

### 1️⃣ Clone 專案

```bash
git clone https://github.com/你的帳號/coupon-system.git
cd coupon-system
```

### 2️⃣ 啟動容器

```bash
docker-compose up -d --build
```

- Gin 應用預設在 `http://localhost:8082` (打api 路徑用這個)
- Nginx 作為反向代理，預設代理到容器內的 port `8080`

### 3️⃣ 執行種子資料（建立使用者、優惠券、使用紀錄）

```bash
docker-compose exec app go run cmd/seed/seed.go
```

---

## 🎯 種子資料說明

執行 `seed.go` 程式後，會產生以下內容：

- 預設建立 2 位使用者（Alice、Bob）
- 建立 10 張優惠券（包含滿額折抵與折扣券）
  - 含已過期、已領完、尚未開始、尚未領取等多樣狀態
- User 1 預設已領取其中三張（含已使用、已過期）

> 可以使用這些資料來測試各種優惠券狀態的處理流程。

---


## 🧪 測試範例與順序建議

> ✅ **建議測試流程：**

1. 先查詢使用者擁有的優惠券狀態（確認哪些是 `unused`, `used`, `expired`）
2. 使用其中一張 `unused` 狀態的優惠券來測試「使用」API。
3. 嘗試重複使用同一張，應收到錯誤訊息。
4. 再領取尚未領過的優惠券。

## 📌 API 說明

目前系統為簡易模擬登入，**所有操作皆預設為 User ID = 1**

### ✅ 領取優惠券

- **路由**：`POST /coupons/:id/redeem`
- **說明**：領取一張指定優惠券

#### 成功範例：

```json
{
  "success": true,
  "status": 200,
  "message": "優惠券領取成功"
}
```

#### 失敗範例（已領過）：

```json
{
  "success": false,
  "status": 409,
  "message": "領取失敗",
  "error": "已領取過該優惠券"
}
```

---

### 🧾 使用優惠券

- **路由**：`POST /coupons/:id/use`
- **說明**：使用指定已領取的優惠券

#### 成功範例：

```json
{
  "success": true,
  "status": 200,
  "message": "優惠券使用成功"
}
```

#### 失敗範例（已過期）：

```json
{
  "success": false,
  "status": 400,
  "message": "使用失敗",
  "error": "優惠券已過期或尚未開始"
}
```

---

### 🔍 查詢使用者的優惠券狀態

- **路由**：`GET /users/:id/coupons`
- **說明**：列出所有該使用者擁有的優惠券及狀態

#### 成功範例：

```json
{
  "success": true,
  "status": 200,
  "message": "查詢成功",
  "data": [
    {
      "id": 1,
      "name": "滿100折20",
      "start_at": "2025-06-10T00:00:00Z",
      "end_at": "2025-06-20T00:00:00Z",
      "status": "unused"
    },
    {
      "id": 2,
      "name": "滿200折50",
      "status": "used"
    },
    {
      "id": 6,
      "name": "7折驚喜券",
      "status": "expired"
    }
  ]
}
```

---

## 🧠 系統設計備註

- **防止超發：** 使用 Redis `SETNX` 搭配券鎖實現。
- **Redis 快取優惠券資料**，效期 10 分鐘。
- **資料表關聯：**
  - `users`: 使用者基本資料
  - `coupons`: 優惠券內容與限制
  - `coupon_usages`: 使用者與優惠券的關聯、使用狀態

---

© 2025 Eric Liao