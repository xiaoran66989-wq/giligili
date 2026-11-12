# G站学习项目

这是 [gourouting/giligili](https://github.com/gourouting/giligili) 的个人 Fork，增加了可直接访问的中文视频首页、列表搜索、播放弹窗、注册登录和投稿页面，并保留 Gin 后端 API。项目用于学习和演示，不是开箱即用的商业视频平台。

> 源项目未声明开源许可证。Fork 仅用于在 GitHub 上保存和实验；如需对外发布、商用或再分发，请先向原作者确认授权。

## 当前功能

- 用户注册、登录和会话
- 视频列表、详情、投稿、编辑与删除
- Redis 播放量和每日排行
- 阿里云 OSS 直传签名
- Swagger API 文档

视频浏览接口为公开接口。投稿、编辑、删除和上传签名接口需要登录；视频只能由创建者编辑或删除。数据库会在启动时自动迁移。升级已有数据库时，旧视频没有归属用户，自动迁移后不能通过新权限规则编辑或删除；迁移前请备份数据，再由管理员按需处理旧记录。

## 本地运行

需要 Go、MySQL、Redis。OSS 配置只在使用上传接口时需要。

1. 创建数据库 `giligili`，启动 MySQL 和 Redis。
2. 复制 `.env.example` 为 `.env`，填写本机连接地址、强随机 `SESSION_SECRET`，并将 `MYSQL_DSN` 中的主机改为本机地址。
3. 本地 HTTP 开发保持 `GIN_MODE=debug`。
4. 在项目目录运行：

~~~sh
go mod tidy
go run .
~~~

服务监听 `http://localhost:3000`。打开 `http://localhost:3000/` 即可访问前端；健康检查接口为 `GET /healthz`，Swagger 文档位于 `http://localhost:3000/swagger/index.html`。

## Docker Compose 演示

1. 将 `.env.example` 复制为 `.env`。
2. 为 `MYSQL_ROOT_PASSWORD`、`MYSQL_DSN` 和 `SESSION_SECRET` 设置强随机值。DSN 中的用户名和密码必须与 MySQL 容器一致。
3. 如前端不是 `http://localhost:8080`，设置 `CORS_ALLOWED_ORIGINS` 为前端完整来源地址；多个来源用逗号分隔。
4. 启动：

~~~sh
docker compose -f docker-compose.yml.example up --build
~~~

网站和 API 会映射到 `http://localhost:3000`，打开该地址即可使用。数据保存在 Docker 卷 `mysql_data` 中。停止服务使用 Ctrl+C；后台运行时用 `docker compose -f docker-compose.yml.example down`。

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

这是学习项目，尚未做完整的生产环境审计。运行需要 Go 后端、MySQL 和 Redis；已有外链 MP4/封面可直接投稿展示，若使用阿里云 OSS 上传则需另行配置 OSS。GitHub Pages 只能发布静态文件，不能运行这个 Go 后端。

`GIN_MODE=release` 时会为会话 Cookie 启用 Secure 属性，因此服务应通过 HTTPS 提供。请勿把真实密码、OSS 密钥或生产环境的 `.env` 文件提交到公开仓库。Compose 文件用于开发演示；上线前应评估依赖版本、数据库备份、HTTPS 代理、限流、日志和监控。

## 上游资料

- 原项目：[gourouting/giligili](https://github.com/gourouting/giligili)
- 项目介绍与教程链接见原项目 README。
