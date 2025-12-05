package admin

import (
	"encoding/json"
	"rustdesk-api-server-pro/app/model"
	"rustdesk-api-server-pro/config"
	"rustdesk-api-server-pro/db"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"xorm.io/xorm"
)

type DocHelpController struct {
	basicController
}

func (c *DocHelpController) BeforeActivation(b mvc.BeforeActivation) {
	// Knowledge Base Categories
	b.Handle("GET", "/dochelp/categories", "HandleGetCategories")
	b.Handle("POST", "/dochelp/categories", "HandleCreateCategory")
	b.Handle("PUT", "/dochelp/categories/{id:int}", "HandleUpdateCategory")
	b.Handle("DELETE", "/dochelp/categories/{id:int}", "HandleDeleteCategory")

	// Knowledge Base Articles
	b.Handle("GET", "/dochelp/articles", "HandleGetArticles")
	b.Handle("GET", "/dochelp/articles/{id:int}", "HandleGetArticle")
	b.Handle("POST", "/dochelp/articles", "HandleCreateArticle")
	b.Handle("PUT", "/dochelp/articles/{id:int}", "HandleUpdateArticle")
	b.Handle("DELETE", "/dochelp/articles/{id:int}", "HandleDeleteArticle")
	b.Handle("POST", "/dochelp/articles/{id:int}/view", "HandleIncrementView")

	// Tickets
	b.Handle("GET", "/dochelp/tickets", "HandleGetTickets")
	b.Handle("GET", "/dochelp/tickets/{id:int}", "HandleGetTicket")
	b.Handle("POST", "/dochelp/tickets", "HandleCreateTicket")
	b.Handle("PUT", "/dochelp/tickets/{id:int}", "HandleUpdateTicket")
	b.Handle("POST", "/dochelp/tickets/{id:int}/resolve", "HandleResolveTicket")
	b.Handle("DELETE", "/dochelp/tickets/{id:int}", "HandleDeleteTicket")

	// Ticket Comments
	b.Handle("GET", "/dochelp/tickets/{id:int}/comments", "HandleGetComments")
	b.Handle("POST", "/dochelp/tickets/{id:int}/comments", "HandleAddComment")
	b.Handle("DELETE", "/dochelp/comments/{id:int}", "HandleDeleteComment")

	// Convert Ticket to Article
	b.Handle("POST", "/dochelp/tickets/{id:int}/convert", "HandleConvertToArticle")

	// File Upload
	b.Handle("POST", "/dochelp/upload", "HandleFileUpload")

	// Statistics
	b.Handle("GET", "/dochelp/stats", "HandleGetStats")
}

// ========== CATEGORIES ==========

func (c *DocHelpController) HandleGetCategories() mvc.Result {
	categories := make([]model.KnowledgeBaseCategory, 0)
	err := c.Db.Asc("order", "name").Find(&categories)
	if err != nil {
		return c.Error(nil, err.Error())
	}

	return c.Success(iris.Map{
		"categories": categories,
	}, "ok")
}

func (c *DocHelpController) HandleCreateCategory() mvc.Result {
	// Require SUPPORT_N2 or higher
	if err := c.RequirePermission(model.ROLE_SUPPORT_N2, "create category"); err != nil {
		return err
	}

	var form struct {
		Name  string `json:"name"`
		Icon  string `json:"icon"`
		Order int    `json:"order"`
	}

	if err := c.Ctx.ReadJSON(&form); err != nil {
		return c.Error(nil, "Invalid request body")
	}

	if form.Name == "" {
		return c.Error(nil, "Category name is required")
	}

	category := model.KnowledgeBaseCategory{
		Name:  form.Name,
		Icon:  form.Icon,
		Order: form.Order,
	}

	_, err := c.Db.Insert(&category)
	if err != nil {
		return c.Error(nil, err.Error())
	}

	// Return the complete category with ID
	return c.Success(iris.Map{
		"id":       category.Id,
		"category": category,
	}, "Category created successfully")
}

func (c *DocHelpController) HandleUpdateCategory() mvc.Result {
	if err := c.RequirePermission(model.ROLE_SUPPORT_N2, "update category"); err != nil {
		return err
	}

	id, _ := c.Ctx.Params().GetInt("id")

	var form struct {
		Name  string `json:"name"`
		Icon  string `json:"icon"`
		Order int    `json:"order"`
	}

	if err := c.Ctx.ReadJSON(&form); err != nil {
		return c.Error(nil, "Invalid request body")
	}

	category := model.KnowledgeBaseCategory{
		Name:  form.Name,
		Icon:  form.Icon,
		Order: form.Order,
	}

	_, err := c.Db.ID(id).Update(&category)
	if err != nil {
		return c.Error(nil, err.Error())
	}

	return c.Success(nil, "Category updated successfully")
}

func (c *DocHelpController) HandleDeleteCategory() mvc.Result {
	if err := c.RequirePermission(model.ROLE_SUPER_ADMIN, "delete category"); err != nil {
		return err
	}

	id, _ := c.Ctx.Params().GetInt("id")
	_, err := c.Db.ID(id).Delete(&model.KnowledgeBaseCategory{})
	if err != nil {
		return c.Error(nil, err.Error())
	}

	return c.Success(nil, "Category deleted successfully")
}

// ========== ARTICLES ==========

func (c *DocHelpController) HandleGetArticles() mvc.Result {
	currentPage := c.Ctx.URLParamIntDefault("current", 1)
	pageSize := c.Ctx.URLParamIntDefault("size", 20)
	categoryId := c.Ctx.URLParamIntDefault("category_id", 0)
	search := c.Ctx.URLParamDefault("search", "")
	pinnedOnly := c.Ctx.URLParamDefault("pinned", "")

	query := func() *xorm.Session {
		q := c.Db.Table(&model.KnowledgeBaseArticle{})
		
		if categoryId > 0 {
			q.Where("category_id = ?", categoryId)
		}
		if search != "" {
			q.Where("title LIKE ? OR content LIKE ? OR tags LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
		}
		if pinnedOnly == "true" {
			q.Where("is_pinned = ?", true)
		}
		
		q.Desc("is_pinned", "created_at")
		return q
	}

	pagination := db.NewPagination(currentPage, pageSize)
	articles := make([]model.KnowledgeBaseArticle, 0)

	err := pagination.Paginate(query, &model.KnowledgeBaseArticle{}, &articles)
	if err != nil {
		return c.Error(nil, err.Error())
	}

	list := make([]iris.Map, 0)
	for _, article := range articles {
		// Get category name
		var category model.KnowledgeBaseCategory
		c.Db.ID(article.CategoryId).Get(&category)

		list = append(list, iris.Map{
			"id":            article.Id,
			"category_id":   article.CategoryId,
			"category_name": category.Name,
			"title":         article.Title,
			"content":       article.Content,
			"tags":          article.Tags,
			"is_pinned":     article.IsPinned,
			"views":         article.Views,
			"author_id":     article.AuthorId,
			"author_name":   article.AuthorName,
			"from_ticket":   article.FromTicket,
			"created_at":    article.CreatedAt.Format(config.TimeFormat),
			"updated_at":    article.UpdatedAt.Format(config.TimeFormat),
		})
	}

	return c.Success(iris.Map{
		"total":   pagination.TotalCount,
		"records": list,
		"current": currentPage,
		"size":    pageSize,
	}, "ok")
}

func (c *DocHelpController) HandleGetArticle() mvc.Result {
	id, _ := c.Ctx.Params().GetInt("id")

	var article model.KnowledgeBaseArticle
	has, err := c.Db.ID(id).Get(&article)
	if err != nil || !has {
		return c.Error(nil, "Article not found")
	}

	// Get category
	var category model.KnowledgeBaseCategory
	c.Db.ID(article.CategoryId).Get(&category)

	return c.Success(iris.Map{
		"id":            article.Id,
		"category_id":   article.CategoryId,
		"category_name": category.Name,
		"title":         article.Title,
		"content":       article.Content,
		"tags":          article.Tags,
		"is_pinned":     article.IsPinned,
		"views":         article.Views,
		"author_id":     article.AuthorId,
		"author_name":   article.AuthorName,
		"from_ticket":   article.FromTicket,
		"created_at":    article.CreatedAt.Format(config.TimeFormat),
		"updated_at":    article.UpdatedAt.Format(config.TimeFormat),
	}, "ok")
}

func (c *DocHelpController) HandleCreateArticle() mvc.Result {
	if err := c.RequirePermission(model.ROLE_SUPPORT_N2, "create article"); err != nil {
		return err
	}

	var form struct {
		CategoryId int      `json:"category_id"`
		Title      string   `json:"title"`
		Content    string   `json:"content"`
		Tags       []string `json:"tags"`
		IsPinned   bool     `json:"is_pinned"`
	}

	if err := c.Ctx.ReadJSON(&form); err != nil {
		return c.Error(nil, "Invalid request body")
	}

	if form.Title == "" || form.Content == "" {
		return c.Error(nil, "Title and content are required")
	}

	user := c.GetUser()
	tagsJSON, _ := json.Marshal(form.Tags)

	article := model.KnowledgeBaseArticle{
		CategoryId: form.CategoryId,
		Title:      form.Title,
		Content:    form.Content,
		Tags:       string(tagsJSON),
		IsPinned:   form.IsPinned,
		AuthorId:   user.Id,
		AuthorName: user.Username,
	}

	_, err := c.Db.Insert(&article)
	if err != nil {
		return c.Error(nil, err.Error())
	}

	return c.Success(iris.Map{"id": article.Id}, "Article created successfully")
}

func (c *DocHelpController) HandleUpdateArticle() mvc.Result {
	if err := c.RequirePermission(model.ROLE_SUPPORT_N2, "update article"); err != nil {
		return err
	}

	id, _ := c.Ctx.Params().GetInt("id")

	var form struct {
		CategoryId int      `json:"category_id"`
		Title      string   `json:"title"`
		Content    string   `json:"content"`
		Tags       []string `json:"tags"`
		IsPinned   bool     `json:"is_pinned"`
	}

	if err := c.Ctx.ReadJSON(&form); err != nil {
		return c.Error(nil, "Invalid request body")
	}

	tagsJSON, _ := json.Marshal(form.Tags)

	article := model.KnowledgeBaseArticle{
		CategoryId: form.CategoryId,
		Title:      form.Title,
		Content:    form.Content,
		Tags:       string(tagsJSON),
		IsPinned:   form.IsPinned,
	}

	_, err := c.Db.ID(id).Cols("category_id", "title", "content", "tags", "is_pinned").Update(&article)
	if err != nil {
		return c.Error(nil, err.Error())
	}

	return c.Success(nil, "Article updated successfully")
}

func (c *DocHelpController) HandleDeleteArticle() mvc.Result {
	if err := c.RequirePermission(model.ROLE_SUPER_ADMIN, "delete article"); err != nil {
		return err
	}

	id, _ := c.Ctx.Params().GetInt("id")
	_, err := c.Db.ID(id).Delete(&model.KnowledgeBaseArticle{})
	if err != nil {
		return c.Error(nil, err.Error())
	}

	return c.Success(nil, "Article deleted successfully")
}

func (c *DocHelpController) HandleIncrementView() mvc.Result {
	id, _ := c.Ctx.Params().GetInt("id")
	
	_, err := c.Db.Exec("UPDATE kb_article SET views = views + 1 WHERE id = ?", id)
	if err != nil {
		return c.Error(nil, err.Error())
	}

	return c.Success(nil, "View counted")
}

// ========== TICKETS ==========

func (c *DocHelpController) HandleGetTickets() mvc.Result {
	currentPage := c.Ctx.URLParamIntDefault("current", 1)
	pageSize := c.Ctx.URLParamIntDefault("size", 20)
	status := c.Ctx.URLParamIntDefault("status", 0)
	priority := c.Ctx.URLParamIntDefault("priority", 0)
	creatorId := c.Ctx.URLParamIntDefault("creator_id", 0)

	user := c.GetUser()

	query := func() *xorm.Session {
		q := c.Db.Table(&model.Ticket{})
		
		// Users (role 1-2) só veem seus próprios tickets
		if user.Role < model.ROLE_SUPPORT_N2 {
			q.Where("creator_id = ?", user.Id)
		} else if creatorId > 0 {
			q.Where("creator_id = ?", creatorId)
		}
		
		if status > 0 {
			q.Where("status = ?", status)
		}
		if priority > 0 {
			q.Where("priority = ?", priority)
		}
		
		q.Desc("priority", "created_at")
		return q
	}

	pagination := db.NewPagination(currentPage, pageSize)
	tickets := make([]model.Ticket, 0)

	err := pagination.Paginate(query, &model.Ticket{}, &tickets)
	if err != nil {
		return c.Error(nil, err.Error())
	}

	list := make([]iris.Map, 0)
	for _, ticket := range tickets {
		// Count comments
		commentCount, _ := c.Db.Where("ticket_id = ?", ticket.Id).Count(&model.TicketComment{})

		list = append(list, iris.Map{
			"id":            ticket.Id,
			"title":         ticket.Title,
			"description":   ticket.Description,
			"status":        ticket.Status,
			"priority":      ticket.Priority,
			"category_id":   ticket.CategoryId,
			"attachments":   ticket.Attachments,
			"creator_id":    ticket.CreatorId,
			"creator_name":  ticket.CreatorName,
			"assigned_to":   ticket.AssignedTo,
			"resolved_by":   ticket.ResolvedBy,
			"resolved_at":   ticket.ResolvedAt.Format(config.TimeFormat),
			"comment_count": commentCount,
			"created_at":    ticket.CreatedAt.Format(config.TimeFormat),
			"updated_at":    ticket.UpdatedAt.Format(config.TimeFormat),
		})
	}

	return c.Success(iris.Map{
		"total":   pagination.TotalCount,
		"records": list,
		"current": currentPage,
		"size":    pageSize,
	}, "ok")
}

func (c *DocHelpController) HandleGetTicket() mvc.Result {
	id, _ := c.Ctx.Params().GetInt("id")
	user := c.GetUser()

	var ticket model.Ticket
	has, err := c.Db.ID(id).Get(&ticket)
	if err != nil || !has {
		return c.Error(nil, "Ticket not found")
	}

	// Check permission
	if user.Role < model.ROLE_SUPPORT_N2 && ticket.CreatorId != user.Id {
		return c.Error(nil, "Permission denied")
	}

	return c.Success(iris.Map{
		"id":           ticket.Id,
		"title":        ticket.Title,
		"description":  ticket.Description,
		"status":       ticket.Status,
		"priority":     ticket.Priority,
		"category_id":  ticket.CategoryId,
		"attachments":  ticket.Attachments,
		"creator_id":   ticket.CreatorId,
		"creator_name": ticket.CreatorName,
		"assigned_to":  ticket.AssignedTo,
		"resolved_by":  ticket.ResolvedBy,
		"resolved_at":  ticket.ResolvedAt.Format(config.TimeFormat),
		"created_at":   ticket.CreatedAt.Format(config.TimeFormat),
		"updated_at":   ticket.UpdatedAt.Format(config.TimeFormat),
	}, "ok")
}

func (c *DocHelpController) HandleCreateTicket() mvc.Result {
	var form struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Priority    int      `json:"priority"`
		CategoryId  int      `json:"category_id"`
		Attachments []string `json:"attachments"`
	}

	if err := c.Ctx.ReadJSON(&form); err != nil {
		return c.Error(nil, "Invalid request body")
	}

	if form.Title == "" || form.Description == "" {
		return c.Error(nil, "Title and description are required")
	}

	user := c.GetUser()
	attachmentsJSON, _ := json.Marshal(form.Attachments)

	ticket := model.Ticket{
		Title:       form.Title,
		Description: form.Description,
		Status:      1, // Aberto
		Priority:    form.Priority,
		CategoryId:  form.CategoryId,
		Attachments: string(attachmentsJSON),
		CreatorId:   user.Id,
		CreatorName: user.Username,
	}

	_, err := c.Db.Insert(&ticket)
	if err != nil {
		return c.Error(nil, err.Error())
	}

	return c.Success(iris.Map{"id": ticket.Id}, "Ticket created successfully")
}

func (c *DocHelpController) HandleUpdateTicket() mvc.Result {
	if err := c.RequirePermission(model.ROLE_SUPPORT_N2, "update ticket"); err != nil {
		return err
	}

	id, _ := c.Ctx.Params().GetInt("id")

	var form struct {
		Status     int `json:"status"`
		Priority   int `json:"priority"`
		AssignedTo int `json:"assigned_to"`
	}

	if err := c.Ctx.ReadJSON(&form); err != nil {
		return c.Error(nil, "Invalid request body")
	}

	ticket := model.Ticket{
		Status:     form.Status,
		Priority:   form.Priority,
		AssignedTo: form.AssignedTo,
	}

	_, err := c.Db.ID(id).Cols("status", "priority", "assigned_to").Update(&ticket)
	if err != nil {
		return c.Error(nil, err.Error())
	}

	return c.Success(nil, "Ticket updated successfully")
}

func (c *DocHelpController) HandleResolveTicket() mvc.Result {
	if err := c.RequirePermission(model.ROLE_SUPPORT_N2, "resolve ticket"); err != nil {
		return err
	}

	id, _ := c.Ctx.Params().GetInt("id")
	user := c.GetUser()

	ticket := model.Ticket{
		Status:     3, // Resolvido
		ResolvedBy: user.Id,
		ResolvedAt: time.Now(),
	}

	_, err := c.Db.ID(id).Cols("status", "resolved_by", "resolved_at").Update(&ticket)
	if err != nil {
		return c.Error(nil, err.Error())
	}

	return c.Success(nil, "Ticket resolved successfully")
}

func (c *DocHelpController) HandleDeleteTicket() mvc.Result {
	if err := c.RequirePermission(model.ROLE_SUPER_ADMIN, "delete ticket"); err != nil {
		return err
	}

	id, _ := c.Ctx.Params().GetInt("id")
	
	// Delete comments first
	c.Db.Where("ticket_id = ?", id).Delete(&model.TicketComment{})
	
	// Delete ticket
	_, err := c.Db.ID(id).Delete(&model.Ticket{})
	if err != nil {
		return c.Error(nil, err.Error())
	}

	return c.Success(nil, "Ticket deleted successfully")
}

// ========== COMMENTS ==========

func (c *DocHelpController) HandleGetComments() mvc.Result {
	ticketId, _ := c.Ctx.Params().GetInt("id")
	user := c.GetUser()

	comments := make([]model.TicketComment, 0)
	query := c.Db.Where("ticket_id = ?", ticketId)
	
	// Users normais não veem comentários internos
	if user.Role < model.ROLE_SUPPORT_N2 {
		query.Where("is_internal = ?", false)
	}
	
	err := query.Asc("created_at").Find(&comments)
	if err != nil {
		return c.Error(nil, err.Error())
	}

	list := make([]iris.Map, 0)
	for _, comment := range comments {
		list = append(list, iris.Map{
			"id":          comment.Id,
			"ticket_id":   comment.TicketId,
			"user_id":     comment.UserId,
			"username":    comment.Username,
			"comment":     comment.Comment,
			"attachments": comment.Attachments,
			"is_internal": comment.IsInternal,
			"created_at":  comment.CreatedAt.Format(config.TimeFormat),
		})
	}

	return c.Success(iris.Map{
		"comments": list,
	}, "ok")
}

func (c *DocHelpController) HandleAddComment() mvc.Result {
	ticketId, _ := c.Ctx.Params().GetInt("id")
	user := c.GetUser()

	var form struct {
		Comment     string   `json:"comment"`
		Attachments []string `json:"attachments"`
		IsInternal  bool     `json:"is_internal"`
	}

	if err := c.Ctx.ReadJSON(&form); err != nil {
		return c.Error(nil, "Invalid request body")
	}

	if form.Comment == "" {
		return c.Error(nil, "Comment cannot be empty")
	}

	// Only support staff can make internal comments
	if form.IsInternal && user.Role < model.ROLE_SUPPORT_N2 {
		form.IsInternal = false
	}

	attachmentsJSON, _ := json.Marshal(form.Attachments)

	comment := model.TicketComment{
		TicketId:    ticketId,
		UserId:      user.Id,
		Username:    user.Username,
		Comment:     form.Comment,
		Attachments: string(attachmentsJSON),
		IsInternal:  form.IsInternal,
	}

	_, err := c.Db.Insert(&comment)
	if err != nil {
		return c.Error(nil, err.Error())
	}

	return c.Success(iris.Map{"id": comment.Id}, "Comment added successfully")
}

func (c *DocHelpController) HandleDeleteComment() mvc.Result {
	id, _ := c.Ctx.Params().GetInt("id")
	user := c.GetUser()

	var comment model.TicketComment
	has, _ := c.Db.ID(id).Get(&comment)
	
	// Only the comment author or support staff can delete
	if !has || (comment.UserId != user.Id && user.Role < model.ROLE_SUPPORT_N2) {
		return c.Error(nil, "Permission denied")
	}

	_, err := c.Db.ID(id).Delete(&model.TicketComment{})
	if err != nil {
		return c.Error(nil, err.Error())
	}

	return c.Success(nil, "Comment deleted successfully")
}

// ========== CONVERT TICKET TO ARTICLE ==========

func (c *DocHelpController) HandleConvertToArticle() mvc.Result {
	if err := c.RequirePermission(model.ROLE_SUPPORT_N2, "convert ticket to article"); err != nil {
		return err
	}

	ticketId, _ := c.Ctx.Params().GetInt("id")
	user := c.GetUser()

	var form struct {
		CategoryId int      `json:"category_id"`
		Title      string   `json:"title"`
		Tags       []string `json:"tags"`
	}

	if err := c.Ctx.ReadJSON(&form); err != nil {
		return c.Error(nil, "Invalid request body")
	}

	// Get ticket
	var ticket model.Ticket
	has, err := c.Db.ID(ticketId).Get(&ticket)
	if err != nil || !has {
		return c.Error(nil, "Ticket not found")
	}

	// Get comments to include in article
	comments := make([]model.TicketComment, 0)
	c.Db.Where("ticket_id = ? AND is_internal = ?", ticketId, false).Asc("created_at").Find(&comments)

	// Build article content from ticket + comments
	content := "## Problema Original\n\n" + ticket.Description + "\n\n"
	
	if len(comments) > 0 {
		content += "## Solução\n\n"
		for _, comment := range comments {
			content += "**" + comment.Username + ":** " + comment.Comment + "\n\n"
		}
	}

	tagsJSON, _ := json.Marshal(form.Tags)

	article := model.KnowledgeBaseArticle{
		CategoryId: form.CategoryId,
		Title:      form.Title,
		Content:    content,
		Tags:       string(tagsJSON),
		AuthorId:   user.Id,
		AuthorName: user.Username,
		FromTicket: ticketId,
	}

	_, err = c.Db.Insert(&article)
	if err != nil {
		return c.Error(nil, err.Error())
	}

	return c.Success(iris.Map{"article_id": article.Id}, "Ticket converted to article successfully")
}

// ========== FILE UPLOAD ==========

func (c *DocHelpController) HandleFileUpload() mvc.Result {
	// Max 10MB
	maxSize := int64(10 << 20)
	c.Ctx.SetMaxRequestBodySize(maxSize)

	file, info, err := c.Ctx.FormFile("file")
	if err != nil {
		return c.Error(nil, "Failed to read file")
	}
	defer file.Close()

	// Create uploads directory if not exists
	uploadDir := "./uploads/dochelp"
	if err := config.EnsureDir(uploadDir); err != nil {
		return c.Error(nil, "Failed to create upload directory")
	}

	// Generate unique filename
	filename := time.Now().Format("20060102_150405") + "_" + info.Filename
	filepath := uploadDir + "/" + filename

	// Save file
	if _, err := c.Ctx.SaveFormFile(info, filepath); err != nil {
		return c.Error(nil, "Failed to save file")
	}

	return c.Success(iris.Map{
		"filename": filename,
		"url":      "/uploads/dochelp/" + filename,
		"size":     info.Size,
	}, "File uploaded successfully")
}

// ========== STATISTICS ==========

func (c *DocHelpController) HandleGetStats() mvc.Result {
	// Total articles
	totalArticles, _ := c.Db.Count(&model.KnowledgeBaseArticle{})
	
	// Total tickets by status
	openTickets, _ := c.Db.Where("status = ?", 1).Count(&model.Ticket{})
	inProgressTickets, _ := c.Db.Where("status = ?", 2).Count(&model.Ticket{})
	resolvedTickets, _ := c.Db.Where("status = ?", 3).Count(&model.Ticket{})
	
	// Most viewed articles
	mostViewed := make([]model.KnowledgeBaseArticle, 0)
	c.Db.Desc("views").Limit(5).Find(&mostViewed)

	return c.Success(iris.Map{
		"total_articles":       totalArticles,
		"open_tickets":         openTickets,
		"in_progress_tickets":  inProgressTickets,
		"resolved_tickets":     resolvedTickets,
		"most_viewed_articles": mostViewed,
	}, "ok")
}
