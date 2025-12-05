package model

import "time"

// KnowledgeBaseCategory - Categorias da base de conhecimento
type KnowledgeBaseCategory struct {
	Id        int       `xorm:"'id' int notnull pk autoincr" json:"id"`
	Name      string    `xorm:"'name' varchar(100) notnull" json:"name"`
	Icon      string    `xorm:"'icon' varchar(50)" json:"icon"` // Nome do ícone (ex: "mdi:help-circle")
	Order     int       `xorm:"'order' int default 0" json:"order"`
	CreatedAt time.Time `xorm:"'created_at' datetime created" json:"created_at"`
}

func (KnowledgeBaseCategory) TableName() string {
	return "kb_category"
}

// KnowledgeBaseArticle - Artigos da base de conhecimento
type KnowledgeBaseArticle struct {
	Id         int       `xorm:"'id' int notnull pk autoincr" json:"id"`
	CategoryId int       `xorm:"'category_id' int notnull" json:"category_id"`
	Title      string    `xorm:"'title' varchar(255) notnull" json:"title"`
	Content    string    `xorm:"'content' text notnull" json:"content"`
	Tags       string    `xorm:"'tags' varchar(255)" json:"tags"` // JSON array de tags
	IsPinned   bool      `xorm:"'is_pinned' tinyint default 0" json:"is_pinned"`
	Views      int       `xorm:"'views' int default 0" json:"views"`
	AuthorId   int       `xorm:"'author_id' int notnull" json:"author_id"`
	AuthorName string    `xorm:"'author_name' varchar(50)" json:"author_name"`
	FromTicket int       `xorm:"'from_ticket' int default 0" json:"from_ticket"` // ID do ticket que originou (se aplicável)
	CreatedAt  time.Time `xorm:"'created_at' datetime created" json:"created_at"`
	UpdatedAt  time.Time `xorm:"'updated_at' datetime updated" json:"updated_at"`
}

func (KnowledgeBaseArticle) TableName() string {
	return "kb_article"
}

// Ticket - Problemas reportados
type Ticket struct {
	Id          int       `xorm:"'id' int notnull pk autoincr" json:"id"`
	Title       string    `xorm:"'title' varchar(255) notnull" json:"title"`
	Description string    `xorm:"'description' text notnull" json:"description"`
	Status      int       `xorm:"'status' int default 1" json:"status"` // 1=Aberto, 2=Em Análise, 3=Resolvido
	Priority    int       `xorm:"'priority' int default 2" json:"priority"` // 1=Baixa, 2=Média, 3=Alta, 4=Crítica
	CategoryId  int       `xorm:"'category_id' int" json:"category_id"`
	Attachments string    `xorm:"'attachments' text" json:"attachments"` // JSON array de URLs/paths
	CreatorId   int       `xorm:"'creator_id' int notnull" json:"creator_id"`
	CreatorName string    `xorm:"'creator_name' varchar(50)" json:"creator_name"`
	AssignedTo  int       `xorm:"'assigned_to' int default 0" json:"assigned_to"` // ID do responsável
	ResolvedBy  int       `xorm:"'resolved_by' int default 0" json:"resolved_by"`
	ResolvedAt  time.Time `xorm:"'resolved_at' datetime null" json:"resolved_at"`
	CreatedAt   time.Time `xorm:"'created_at' datetime created" json:"created_at"`
	UpdatedAt   time.Time `xorm:"'updated_at' datetime updated" json:"updated_at"`
}

func (Ticket) TableName() string {
	return "ticket"
}

// TicketComment - Comentários/Respostas em tickets
type TicketComment struct {
	Id          int       `xorm:"'id' int notnull pk autoincr" json:"id"`
	TicketId    int       `xorm:"'ticket_id' int notnull" json:"ticket_id"`
	UserId      int       `xorm:"'user_id' int notnull" json:"user_id"`
	Username    string    `xorm:"'username' varchar(50)" json:"username"`
	Comment     string    `xorm:"'comment' text notnull" json:"comment"`
	Attachments string    `xorm:"'attachments' text" json:"attachments"` // JSON array
	IsInternal  bool      `xorm:"'is_internal' tinyint default 0" json:"is_internal"` // Comentário interno (só admin vê)
	CreatedAt   time.Time `xorm:"'created_at' datetime created" json:"created_at"`
}

func (TicketComment) TableName() string {
	return "ticket_comment"
}
