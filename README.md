# G站学习项目

这是 [gourouting/giligili](https://github.com/gourouting/giligili) 的个人 Fork，增加了可直接访问的中文视频首页、列表搜索、播放弹窗、注册登录和投稿页面，并保留 Gin 后端 API。项目用于学习和演示，不是开箱即用的商业视频平台。

> 源项目未声明开源许可证。Fork 仅用于在 GitHub 上保存和实验；如需对外发布、商用或再分发，请先向原作者确认授权。

## 当前功能

- 用户注册、登录和会话
- 视频列表、详情、投稿、编辑与删除
- 数据库存储的播放量和每日排行
- 阿里云 OSS 直传签名
- Swagger API 文档

视频浏览接口为公开接口。投稿、编辑、删除和上传签名接口需要登录；视频只能由创建者编辑或删除。数据库会在启动时自动迁移。升级已有数据库时，旧视频没有归属用户，自动迁移后不能通过新权限规则编辑或删除；迁移前请备份数据，再由管理员按需处理旧记录。

## Windows 安装包

在 GitHub 仓库的 Actions 中打开最新的 “Windows Installer” 工作流，下载名为 `Giligili-Windows-Installer` 的产物并运行其中的安装程序。安装完成后从开始菜单启动“G站视频”，浏览器会自动打开 `http://localhost:3000/`。应用数据保存在当前 Windows 用户的配置目录中，卸载程序不会删除数据库。

安装包在 Windows 上使用 SQLite 保存账号、视频和播放数；每日排行在程序运行期间统计。关闭启动时出现的黑色服务窗口会停止网站服务。

## Microsoft Store 安装包

Store 路线不需要自备 PFX：Microsoft Store 会在 MSIX 审核通过并发布时替换为 Microsoft 签名。仓库新增了 **Store MSIX Package** 工作流，生成待提交的 `.msix` 文件；生成包时必须使用 Partner Center 为该应用分配的准确 Package Identity Name、Publisher 和 Publisher display name。

开始前，请在 [Microsoft Store 开发者门户](https://storedeveloper.microsoft.com/)创建开发者账号并在 Partner Center 保留应用名称。个人开发者的新注册流程目前免注册费，但需要本人完成身份验证。创建应用后，在 GitHub 的 **Actions → Store MSIX Package → Run workflow** 输入包身份字段并运行，下载 `Giligili-Store-MSIX-版本号` 产物。工作流只生成待提交包，不会替你提交或发布。

Microsoft Store 只会为通过审核的 Store MSIX 自动签名。直接下载的 Inno Setup `.exe` 仍需 CA 代码签名证书。提交前也要确认已获原项目权利人许可；上游仓库没有声明开源许可证。

## Windows 代码签名

安装包工作流已支持 Authenticode 签名。要启用签名，请先取得可信 CA 签发、可导出为 PFX 的 Windows 代码签名证书，然后在仓库的 **Settings → Secrets and variables → Actions** 添加两个 Repository secrets：

- `WINDOWS_SIGNING_CERT_BASE64`：PFX 文件的单行 Base64 内容
- `WINDOWS_SIGNING_CERT_PASSWORD`：PFX 密码

不要把 PFX 文件或密码提交到仓库。配置完成后重新运行 **Windows Installer** 工作流；它会用 SHA-256 与 RFC 3161 时间戳签署程序和安装包，并将产物标记为 `signed`。未设置证书时仍会生成并标记为 `unsigned`。

签名可以显示已验证的发布者身份，但新证书或新文件仍可能触发 SmartScreen 声誉提示。微软说明见[SmartScreen 声誉说明](https://learn.microsoft.com/windows/apps/package-and-deploy/smartscreen-reputation)和[Authenticode 时间戳说明](https://learn.microsoft.com/windows/win32/seccrypto/time-stamping-authenticode-signatures)。

## 本地运行

Windows 安装包使用本地 SQLite，无需安装 Go、MySQL、Redis 或 Docker Desktop。若从源码启动非 Windows 版本，需要 Go、MySQL。OSS 配置只在使用 OSS 上传接口时需要。

1. 创建数据库 `giligili` 并启动 MySQL。
2. 复制 `.env.example` 为 `.env`，填写本机 MySQL 连接地址和强随机 `SESSION_SECRET`。
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

这是学习项目，尚未做完整的生产环境审计。Windows 安装包使用本地 SQLite；源码和 Docker Compose 运行方式使用 MySQL。已有外链 MP4/封面可直接投稿展示，若使用阿里云 OSS 上传则需另行配置 OSS。GitHub Pages 只能发布静态文件，不能运行这个 Go 后端。

`GIN_MODE=release` 时会为会话 Cookie 启用 Secure 属性，因此服务应通过 HTTPS 提供。请勿把真实密码、OSS 密钥或生产环境的 `.env` 文件提交到公开仓库。Compose 文件用于开发演示；上线前应评估依赖版本、数据库备份、HTTPS 代理、限流、日志和监控。

## 上游资料

- 原项目：[gourouting/giligili](https://github.com/gourouting/giligili)
- 项目介绍与教程链接见原项目 README。
