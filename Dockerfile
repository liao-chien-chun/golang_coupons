# 使用官方 Golang 映像檔，alpine 為輕量版，支援 M1/M2（arm64）
FROM golang:1.22-alpine

# 安裝常用工具（例如 curl、tzdata 方便 debug）
RUN apk add --no-cache git curl tzdata

# 設定容器內的工作目錄為 /app
WORKDIR /app

# 複製 go.mod 與 go.sum（這樣可以先下載依賴）
COPY go.mod go.sum ./
RUN go mod download

# 再複製整個專案內容
COPY . .

# 編譯 Gin 專案，輸出為 main（可執行檔）
RUN go build -o main ./cmd/main.go

# 暴露 Gin 使用的 port（僅供內部用，外部由 Nginx 代理）
EXPOSE 8080

# 啟動應用程式
CMD ["./main"]