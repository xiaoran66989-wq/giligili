# G站学习项目

这是 [gourouting/giligili](https://github.com/gourouting/giligili) 的个人 Fork。项目使用 Gin 和 GORM 编写视频站后端 API，主要用于学习 Go 后端开发；它不包含可直接访问的前端页面，也不是开箱即用的商业视频平台。

> 源项目未声明开源许可证。Fork 仅用于在 GitHub 上保存和实验；如需对外发布、商用或再分发，请先向原作者确认授权。

## 当前功能

- 用户注册、登录和会话
- 视频列表、详情、投稿、编辑与删除
- Redis 播放量和每日排行
- 阿里云 OSS 直传签名
- Swagger API 文档

视频浏览接口为公开接口。投稿、编辑、删除和上传签名接口需要登录；视频只能由创建者编辑或删除。数据库会在启动时自动迁移。

## 本地运行

需要 Go、MySQL、Redis。OSS 配置只在使用上传接口时需要。

1. 创建数据库 `giligili`，启动 MySQL 和 Redis。
2. 复制 `.env.example` 为 `.env`，填写本机连接地址和强随机 `SESSION_SECRET`。
3. 在项目目录运行：

~~~sh
go mod download
go run .
~~~

服务监听 `http://localhost:3000`。健康检查接口是 `POST /api/v1/ping`，Swagger 文档位于 `http://localhost:3000/swagger/index.html`。

## Docker Compose 演示

1. 将 `.env.example` 复制为 `.env`。
2. 为 `MYSQL_ROOT_PASSWORD`、`MYSQL_DSN` 和 `SESSION_SECRET` 设置强随机值。DSN 中的用户名和密码必须与 MySQL 容器一致。
3. 如前端不是 `http://localhost:8080`，设置 `CORS_ALLOWED_ORIGINS` 为前端完整来源地址；多个来源用逗号分隔。
4. 启动：

~~~sh
docker compose -f docker-compose.yml.example up --build
~~~

API 会映射到 `http://localhost:3000`。数据保存在 Docker 卷 `mysql_data` 中。停止服务使用 Ctrl+C；后台运行时用 `docker compose -f docker-compose.yml.example down`。

## API 概览

| 方法 | 路径 | 登录 | 用途 |
|---|---|---:|---|
| POST | `/api/v1/user/register` | 否 | 注册 |
| POST | `/api/v1/user/login` | 否 | 登录 |
| GET | `/api/v1/videos` | 否 | 视频列表 |
| GET | `/api/v1/video/:id` | 否 | 视频详情 |
| POST | `/api/v1/videos` | 是 | 投稿 |
| PUT | `/api/v1/video/:id` | 是 | 编辑自己的视频 |
| DELETE | `/api/v1/video/:id` | 是 | 删除自己的视频 |
| POST | `/api/v1/upload/token` | 是 | 获取 OSS 上传签名 |

所有 API 仍沿用项目原有 JSON 响应格式。部署前请检查 Swagger 文档和环境配置。

## 部署说明

这是学习项目，尚未做完整的生产环境审计。它没有前端播放器；运行需要后端服务器、MySQL、Redis 和（上传时）OSS。GitHub Pages 只能发布静态文件，不能运行这个 Go 后端。

`GIN_MODE=release` 时会为会话 Cookie 启用 Secure 属性，因此服务应通过 HTTPS 提供。请勿把真实密码、OSS 密钥或生产环境的 `.env` 文件提交到公开仓库。Compose 文件用于开发演示；上线前应评估依赖版本、数据库备份、HTTPS 代理、限流、日志和监控。

## 上游资料

- 原项目：[gourouting/giligili](https://github.com/gourouting/giligili)
- 项目介绍与教程链接见原项目 README。
