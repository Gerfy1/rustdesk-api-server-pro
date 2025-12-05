-- DocHelp System Tables
-- Knowledge Base + Ticket Support System

CREATE TABLE IF NOT EXISTS kb_category (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(100) NOT NULL,
    icon VARCHAR(50),
    "order" INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS kb_article (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    category_id INTEGER NOT NULL,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    tags VARCHAR(255),
    is_pinned TINYINT DEFAULT 0,
    views INTEGER DEFAULT 0,
    author_id INTEGER NOT NULL,
    author_name VARCHAR(50),
    from_ticket INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (category_id) REFERENCES kb_category(id)
);

CREATE TABLE IF NOT EXISTS ticket (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    status INTEGER DEFAULT 1,
    priority INTEGER DEFAULT 2,
    category_id INTEGER,
    attachments TEXT,
    creator_id INTEGER NOT NULL,
    creator_name VARCHAR(50),
    assigned_to INTEGER DEFAULT 0,
    resolved_by INTEGER DEFAULT 0,
    resolved_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS ticket_comment (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    ticket_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    username VARCHAR(50),
    comment TEXT NOT NULL,
    attachments TEXT,
    is_internal TINYINT DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (ticket_id) REFERENCES ticket(id)
);

-- Indexes for better performance
CREATE INDEX IF NOT EXISTS idx_kb_article_category ON kb_article(category_id);
CREATE INDEX IF NOT EXISTS idx_kb_article_pinned ON kb_article(is_pinned);
CREATE INDEX IF NOT EXISTS idx_kb_article_views ON kb_article(views);
CREATE INDEX IF NOT EXISTS idx_kb_article_from_ticket ON kb_article(from_ticket);

CREATE INDEX IF NOT EXISTS idx_ticket_status ON ticket(status);
CREATE INDEX IF NOT EXISTS idx_ticket_priority ON ticket(priority);
CREATE INDEX IF NOT EXISTS idx_ticket_creator ON ticket(creator_id);
CREATE INDEX IF NOT EXISTS idx_ticket_assigned ON ticket(assigned_to);
CREATE INDEX IF NOT EXISTS idx_ticket_created ON ticket(created_at);

CREATE INDEX IF NOT EXISTS idx_ticket_comment_ticket ON ticket_comment(ticket_id);
CREATE INDEX IF NOT EXISTS idx_ticket_comment_user ON ticket_comment(user_id);
CREATE INDEX IF NOT EXISTS idx_ticket_comment_internal ON ticket_comment(is_internal);

-- Insert default categories
INSERT INTO kb_category (name, icon, "order") VALUES 
    ('Geral', 'mdi:information-outline', 1),
    ('Erros Resolvidos', 'mdi:check-circle', 2),
    ('Tutoriais', 'mdi:book-open-variant', 3),
    ('FAQ', 'mdi:help-circle', 4);
