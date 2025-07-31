## blog 项目结构

```shell
├── api
│   └── v1
├── config
├── core
├── docs
├── global
├── initialize
│   └── internal
├── middleware
├── model
│   ├── request
│   └── response
├── packfile
├── resource
│   ├── excel
│   ├── page
│   └── template
├── router
├── service
├── source
└── utils
    ├── timer
    └── upload
```

| 文件夹       | 说明                    | 描述                        |
| ------------ | ----------------------- | --------------------------- |
| `api`        | api层                   | api层 |
| `--v1`       | v1版本接口              | v1版本接口                  |
| `config`     | 配置包                  | config.yaml对应的配置结构体 |
| `core`       | 核心文件                | 核心组件(zap, viper, server)的初始化 |
| `docs`       | swagger文档目录         | swagger文档目录 |
| `global`     | 全局对象                | 全局对象 |
| `initialize` | 初始化 | router,redis,gorm,validator, timer的初始化 |
| `--internal` | 初始化内部函数 | gorm 的 longger 自定义,在此文件夹的函数只能由 `initialize` 层进行调用 |
| `middleware` | 中间件层 | 用于存放 `gin` 中间件代码 |
| `model`      | 模型层                  | 模型对应数据表              |
| `--request`  | 入参结构体              | 接收前端发送到后端的数据。  |
| `--response` | 出参结构体              | 返回给前端的数据结构体      |
| `packfile`   | 静态文件打包            | 静态文件打包 |
| `resource`   | 静态资源文件夹          | 负责存放静态文件                |
| `--excel` | excel导入导出默认路径 | excel导入导出默认路径 |
| `--page` | 表单生成器 | 表单生成器 打包后的dist |
| `--template` | 模板 | 模板文件夹,存放的是代码生成器的模板 |
| `router`     | 路由层                  | 路由层 |
| `service`    | service层               | 存放业务逻辑问题 |
| `utils`      | 工具包                  | 工具函数封装            |
| `--timer` | timer | 定时器接口封装 |
| `--upload`      | oss                  | oss接口封装        |

## api 介绍
概述
本文档提供博客系统的核心 API 接口说明，包括用户管理、文章操作和评论功能。

## 概述
本文档提供博客系统的核心 API 接口说明，所有返回格式遵循统一结构：

```json
{
  "code": 0,      // 状态码 (0=成功 7=失败)
  "data": {},     // 返回数据
  "msg": ""       // 消息描述
}
```

## 用户管理
1. 用户注册
```json
注册新用户账号
URL: /api/enroll
Method: POST
        
Request Body
{
"username": "liaofan",
"password": "liaofan123",
"email": "liaofan@163.com"
}
        
Response (Success)
{
"code": 0,
"data": "User registered successfully"
}
```
2. 用户登录
```json
URL: /api/login
Method: POST
        
Request Body
{
"username": "liaofan",
"password": "liaofan123"
}
        
Response (Success)
{
    "code": 0,
    "data": {
        "user": {
            "ID": 3,
            "CreatedAt": "2025-07-31T15:04:04.038Z",
            "UpdatedAt": "2025-07-31T15:04:04.038Z",
            "DeletedAt": null,
            "Username": "liaofan1",
            "Password": "...",
            "Email": "liaofan1@163.com"
        },
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
        "expiresAt": 1754579198000
    },
    "msg": "登录成功"
}
```
## 文章管理

3. 创建文章

只有已认证的用户才能创建文章，创建文章时需要提供文章的标题和内容。
```json
URL /api/create-post
Method POST

Request Body
{
  "title": "java 进阶",
  "content": "jvm 底层机制详细讲讲？"
}

Response (Success)
{
  "code": 0,
  "data": {},
  "msg": "操作成功"
}
```
4. 查询文章

支持获取所有文章列表和单个文章的详细信息
```json
URL /api/query-post
Method GET

Response (Success)
{
  "code": 0,
  "data": [
    {
      "title": "java 基础",
      "content": "jvm 是什么？"
    },
    {
      "title": "java 进阶",
      "content": "jvm 底层机制详细讲讲？"
    }
  ],
  "msg": "查询成功"
}

```
5. 更新文章

只有文章的作者才能更新自己的文章
```json
URL /api/update-post
Method PUT

Request Body
{
  "id": 2,
  "title": "java 基础有哪些",
  "content": "jvm 详细机制讲解一下？"
}

Response (Success)
{
  "code": 0,
  "data": {},
  "msg": "操作成功"
}
```

6. 删除文章

只有文章的作者才能删除自己的文章
```json
URL /api/delete-post
Method DELETE

Request Body
{
  "id": 2
}

Response (Success)
{
  "code": 0,
  "data": {},
  "msg": "操作成功"
}
```

## 评论管理

7. 创建评论

实现评论的创建功能，已认证的用户可以对文章发表评论
```json
URL /api/create-comments
Method POST

Request Body
{
  "post_id": 2,
  "content": "文章写的真好！"
}

Response (Success)
{
  "code": 0,
  "data": {},
  "msg": "操作成功"
}
```
8. 查询评论

实现评论的读取功能，支持获取某篇文章的所有评论列表。
```json
URL /api/query-comments
Method GET

Query Parameters
{
  "post_id": 2
}

Response (Success)
{
  "code": 0,
  "data": [
    {"content": "文章写的真好！"},
    {"content": "文章写的真好！"},
    {"content": "文章写的真好！2"}
  ],
  "msg": "查询成功"
}

```