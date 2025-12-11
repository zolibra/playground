package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"
)

// BlogPost represents a blog post
type BlogPost struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
}

// BlogSystem manages blog posts
type BlogSystem struct {
	posts []BlogPost
}

// NewBlogSystem creates a new blog system with sample posts
func NewBlogSystem() *BlogSystem {
	bs := &BlogSystem{
		posts: generateSamplePosts(),
	}
	return bs
}

// generateSamplePosts creates sample blog posts with various dates
func generateSamplePosts() []BlogPost {
	now := time.Now()
	posts := []BlogPost{}
	
	// Create 50 posts with different dates to test pagination
	for i := 1; i <= 50; i++ {
		daysAgo := (i - 1) * 15 // Spread posts over time
		createdAt := now.AddDate(0, 0, -daysAgo)
		
		posts = append(posts, BlogPost{
			ID:        i,
			Title:     fmt.Sprintf("博客文章 %d", i),
			Content:   fmt.Sprintf("这是第 %d 篇博客文章的内容。创建于 %d 天前。", i, daysAgo),
			Author:    "作者",
			CreatedAt: createdAt,
		})
	}
	
	return posts
}

// GetRecentPosts returns posts from the last year, sorted by date (newest first)
func (bs *BlogSystem) GetRecentPosts(page int, perPage int) ([]BlogPost, int) {
	oneYearAgo := time.Now().AddDate(-1, 0, 0)
	recentPosts := []BlogPost{}
	
	for _, post := range bs.posts {
		if post.CreatedAt.After(oneYearAgo) {
			recentPosts = append(recentPosts, post)
		}
	}
	
	// Sort by date, newest first
	sort.Slice(recentPosts, func(i, j int) bool {
		return recentPosts[i].CreatedAt.After(recentPosts[j].CreatedAt)
	})
	
	totalPosts := len(recentPosts)
	
	// Apply pagination
	start := (page - 1) * perPage
	if start >= totalPosts {
		return []BlogPost{}, totalPosts
	}
	
	end := start + perPage
	if end > totalPosts {
		end = totalPosts
	}
	
	return recentPosts[start:end], totalPosts
}

// GetArchivePosts returns posts older than one year, sorted by date (newest first)
func (bs *BlogSystem) GetArchivePosts(page int, perPage int) ([]BlogPost, int) {
	oneYearAgo := time.Now().AddDate(-1, 0, 0)
	archivePosts := []BlogPost{}
	
	for _, post := range bs.posts {
		if post.CreatedAt.Before(oneYearAgo) || post.CreatedAt.Equal(oneYearAgo) {
			archivePosts = append(archivePosts, post)
		}
	}
	
	// Sort by date, newest first
	sort.Slice(archivePosts, func(i, j int) bool {
		return archivePosts[i].CreatedAt.After(archivePosts[j].CreatedAt)
	})
	
	totalPosts := len(archivePosts)
	
	// Apply pagination
	start := (page - 1) * perPage
	if start >= totalPosts {
		return []BlogPost{}, totalPosts
	}
	
	end := start + perPage
	if end > totalPosts {
		end = totalPosts
	}
	
	return archivePosts[start:end], totalPosts
}

const homeTemplate = `
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>博客首页</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f5f5f5;
        }
        h1 {
            color: #333;
            border-bottom: 2px solid #4CAF50;
            padding-bottom: 10px;
        }
        .nav {
            margin: 20px 0;
            padding: 10px;
            background-color: #fff;
            border-radius: 5px;
        }
        .nav a {
            color: #4CAF50;
            text-decoration: none;
            margin-right: 20px;
            font-weight: bold;
        }
        .nav a:hover {
            text-decoration: underline;
        }
        .post {
            background-color: white;
            padding: 20px;
            margin: 15px 0;
            border-radius: 5px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        .post h2 {
            color: #333;
            margin-top: 0;
        }
        .post-meta {
            color: #666;
            font-size: 0.9em;
            margin-bottom: 10px;
        }
        .post-content {
            color: #444;
            line-height: 1.6;
        }
        .pagination {
            margin: 30px 0;
            text-align: center;
        }
        .pagination a, .pagination span {
            display: inline-block;
            padding: 8px 12px;
            margin: 0 5px;
            background-color: white;
            border: 1px solid #ddd;
            border-radius: 3px;
            text-decoration: none;
            color: #333;
        }
        .pagination a:hover {
            background-color: #4CAF50;
            color: white;
        }
        .pagination .current {
            background-color: #4CAF50;
            color: white;
            border-color: #4CAF50;
        }
        .info {
            background-color: #e7f4e7;
            padding: 10px;
            border-left: 4px solid #4CAF50;
            margin: 20px 0;
        }
    </style>
</head>
<body>
    <h1>博客首页</h1>
    <div class="nav">
        <a href="/">首页</a>
        <a href="/archive">归档</a>
    </div>
    <div class="info">
        显示最近一年内的博客文章（最多每页20篇）<br>
        总共 {{.TotalPosts}} 篇文章 | 当前页: {{.CurrentPage}}/{{.TotalPages}}
    </div>
    {{range .Posts}}
    <div class="post">
        <h2>{{.Title}}</h2>
        <div class="post-meta">
            作者: {{.Author}} | 发布时间: {{.CreatedAt.Format "2006-01-02 15:04:05"}}
        </div>
        <div class="post-content">
            {{.Content}}
        </div>
    </div>
    {{end}}
    {{if gt .TotalPages 1}}
    <div class="pagination">
        {{if gt .CurrentPage 1}}
            <a href="/?page={{sub .CurrentPage 1}}">&laquo; 上一页</a>
        {{end}}
        
        {{range $i := .PageNumbers}}
            {{if eq $i $.CurrentPage}}
                <span class="current">{{$i}}</span>
            {{else}}
                <a href="/?page={{$i}}">{{$i}}</a>
            {{end}}
        {{end}}
        
        {{if lt .CurrentPage .TotalPages}}
            <a href="/?page={{add .CurrentPage 1}}">下一页 &raquo;</a>
        {{end}}
    </div>
    {{end}}
</body>
</html>
`

const archiveTemplate = `
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>博客归档</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f5f5f5;
        }
        h1 {
            color: #333;
            border-bottom: 2px solid #4CAF50;
            padding-bottom: 10px;
        }
        .nav {
            margin: 20px 0;
            padding: 10px;
            background-color: #fff;
            border-radius: 5px;
        }
        .nav a {
            color: #4CAF50;
            text-decoration: none;
            margin-right: 20px;
            font-weight: bold;
        }
        .nav a:hover {
            text-decoration: underline;
        }
        .post {
            background-color: white;
            padding: 20px;
            margin: 15px 0;
            border-radius: 5px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        .post h2 {
            color: #333;
            margin-top: 0;
        }
        .post-meta {
            color: #666;
            font-size: 0.9em;
            margin-bottom: 10px;
        }
        .post-content {
            color: #444;
            line-height: 1.6;
        }
        .pagination {
            margin: 30px 0;
            text-align: center;
        }
        .pagination a, .pagination span {
            display: inline-block;
            padding: 8px 12px;
            margin: 0 5px;
            background-color: white;
            border: 1px solid #ddd;
            border-radius: 3px;
            text-decoration: none;
            color: #333;
        }
        .pagination a:hover {
            background-color: #4CAF50;
            color: white;
        }
        .pagination .current {
            background-color: #4CAF50;
            color: white;
            border-color: #4CAF50;
        }
        .info {
            background-color: #fff3cd;
            padding: 10px;
            border-left: 4px solid #ffc107;
            margin: 20px 0;
        }
    </style>
</head>
<body>
    <h1>博客归档</h1>
    <div class="nav">
        <a href="/">首页</a>
        <a href="/archive">归档</a>
    </div>
    <div class="info">
        归档区 - 显示一年以前的博客文章（最多每页20篇）<br>
        总共 {{.TotalPosts}} 篇文章 | 当前页: {{.CurrentPage}}/{{.TotalPages}}
    </div>
    {{range .Posts}}
    <div class="post">
        <h2>{{.Title}}</h2>
        <div class="post-meta">
            作者: {{.Author}} | 发布时间: {{.CreatedAt.Format "2006-01-02 15:04:05"}}
        </div>
        <div class="post-content">
            {{.Content}}
        </div>
    </div>
    {{end}}
    {{if gt .TotalPages 1}}
    <div class="pagination">
        {{if gt .CurrentPage 1}}
            <a href="/archive?page={{sub .CurrentPage 1}}">&laquo; 上一页</a>
        {{end}}
        
        {{range $i := .PageNumbers}}
            {{if eq $i $.CurrentPage}}
                <span class="current">{{$i}}</span>
            {{else}}
                <a href="/archive?page={{$i}}">{{$i}}</a>
            {{end}}
        {{end}}
        
        {{if lt .CurrentPage .TotalPages}}
            <a href="/archive?page={{add .CurrentPage 1}}">下一页 &raquo;</a>
        {{end}}
    </div>
    {{end}}
</body>
</html>
`

func main() {
	bs := NewBlogSystem()
	
	// Template functions
	funcMap := template.FuncMap{
		"sub": func(a, b int) int { return a - b },
		"add": func(a, b int) int { return a + b },
	}
	
	homeTmpl := template.Must(template.New("home").Funcs(funcMap).Parse(homeTemplate))
	archiveTmpl := template.Must(template.New("archive").Funcs(funcMap).Parse(archiveTemplate))
	
	// Homepage handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		pageStr := r.URL.Query().Get("page")
		page := 1
		if pageStr != "" {
			if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
				page = p
			}
		}
		
		perPage := 20
		posts, totalPosts := bs.GetRecentPosts(page, perPage)
		totalPages := (totalPosts + perPage - 1) / perPage
		if totalPages == 0 {
			totalPages = 1
		}
		
		// Generate page numbers for pagination
		pageNumbers := []int{}
		for i := 1; i <= totalPages; i++ {
			pageNumbers = append(pageNumbers, i)
		}
		
		data := map[string]interface{}{
			"Posts":       posts,
			"CurrentPage": page,
			"TotalPages":  totalPages,
			"TotalPosts":  totalPosts,
			"PageNumbers": pageNumbers,
		}
		
		homeTmpl.Execute(w, data)
	})
	
	// Archive handler
	http.HandleFunc("/archive", func(w http.ResponseWriter, r *http.Request) {
		pageStr := r.URL.Query().Get("page")
		page := 1
		if pageStr != "" {
			if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
				page = p
			}
		}
		
		perPage := 20
		posts, totalPosts := bs.GetArchivePosts(page, perPage)
		totalPages := (totalPosts + perPage - 1) / perPage
		if totalPages == 0 {
			totalPages = 1
		}
		
		// Generate page numbers for pagination
		pageNumbers := []int{}
		for i := 1; i <= totalPages; i++ {
			pageNumbers = append(pageNumbers, i)
		}
		
		data := map[string]interface{}{
			"Posts":       posts,
			"CurrentPage": page,
			"TotalPages":  totalPages,
			"TotalPosts":  totalPosts,
			"PageNumbers": pageNumbers,
		}
		
		archiveTmpl.Execute(w, data)
	})
	
	// API endpoint to get posts as JSON
	http.HandleFunc("/api/posts", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(bs.posts)
	})
	
	port := "8080"
	fmt.Printf("博客服务器启动在 http://localhost:%s\n", port)
	fmt.Println("访问 / 查看首页（最近一年的文章）")
	fmt.Println("访问 /archive 查看归档（一年以前的文章）")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
