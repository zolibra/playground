# playground
随便玩玩

## 博客系统

这是一个简单的博客系统，具有以下功能：

### 功能特性

1. **首页展示**：只显示最近一年内的博客文章
2. **归档页面**：显示一年以前的旧博客文章
3. **分页功能**：每页最多显示20篇文章，自动分页
4. **时间排序**：文章按发布时间倒序排列（最新的在前）

### 使用方法

#### 启动博客服务器

```bash
go run blog.go
```

或者先编译再运行：

```bash
go build -o blog-server blog.go
./blog-server
```

服务器将在 `http://localhost:8080` 启动

#### 访问页面

- **首页**：`http://localhost:8080/` - 查看最近一年的文章
- **归档**：`http://localhost:8080/archive` - 查看一年以前的文章
- **API接口**：`http://localhost:8080/api/posts` - 获取所有文章的JSON数据

#### 分页

- 首页分页：`http://localhost:8080/?page=2`
- 归档分页：`http://localhost:8080/archive?page=2`

### 技术实现

- 使用 Go 语言开发
- 内置 HTTP 服务器
- HTML 模板渲染
- 自动生成示例博客文章（50篇，时间跨度约2年）
- 按一年为界限自动区分首页和归档内容

