# Go JWT 認證系統

- 這是一個學習 JWT 認證機制的練習專案，主要學習重點：
- JWT 的基本概念和運作原理
- Go 語言的 Web 開發實踐
- RESTful API 設計

## 功能特點

- 用戶註冊
- 用戶登入（JWT 認證）
- 獲取個人資料
- 獲取所有用戶列表
- 刪除用戶

## 技術架構

- Go
- Gin Framework（Web 框架）
- GORM（ORM 框架）
- SQLite（資料庫）
- JWT（認證機制）

## API 端點

### 公開路由
- `POST /register` - 註冊新用戶
- `POST /login` - 用戶登入
- `GET /allProfile` - 獲取所有用戶列表

### 需要認證的路由
- `GET /user/profile` - 獲取當前用戶資料
  - 需要在 Header 中加入 `Authorization: Bearer {your-jwt-token}`

### 刪除功能
- `DELETE /deleteUser?id={user_id}` - 刪除指定用戶

