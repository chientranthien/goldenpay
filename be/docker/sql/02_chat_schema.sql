-- Chat service database schema

-- Channels table
CREATE TABLE IF NOT EXISTS channels (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    type VARCHAR(50) NOT NULL DEFAULT 'public', -- 'public', 'private', 'direct'
    creator_id BIGINT NOT NULL,
    status TINYINT NOT NULL DEFAULT 1, -- 1=active, 2=archived, 3=deleted
    version BIGINT NOT NULL DEFAULT 1,
    ctime BIGINT NOT NULL,
    mtime BIGINT NOT NULL,
    
    INDEX idx_creator_id (creator_id),
    INDEX idx_type (type),
    INDEX idx_status (status),
    INDEX idx_ctime (ctime)
);

-- Channel members table
CREATE TABLE IF NOT EXISTS channel_members (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    channel_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'member', -- 'admin', 'member'
    joined_at BIGINT NOT NULL,
    last_read_at BIGINT NOT NULL DEFAULT 0,
    
    UNIQUE KEY unique_channel_user (channel_id, user_id),
    INDEX idx_user_id (user_id),
    INDEX idx_channel_id (channel_id),
    INDEX idx_joined_at (joined_at)
);

-- Messages table
CREATE TABLE IF NOT EXISTS messages (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    channel_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    content TEXT NOT NULL,
    message_type VARCHAR(50) NOT NULL DEFAULT 'text', -- 'text', 'file', 'system'
    metadata JSON, -- stores mentions, attachments, reactions, etc.
    status TINYINT NOT NULL DEFAULT 1, -- 1=active, 2=edited, 3=deleted
    version BIGINT NOT NULL DEFAULT 1,
    ctime BIGINT NOT NULL,
    mtime BIGINT NOT NULL,
    reply_to_id BIGINT DEFAULT NULL, -- for threaded conversations
    
    INDEX idx_channel_id (channel_id),
    INDEX idx_user_id (user_id),
    INDEX idx_ctime (ctime),
    INDEX idx_status (status),
    INDEX idx_reply_to_id (reply_to_id),
    INDEX idx_channel_ctime (channel_id, ctime) -- for efficient message history queries
);

-- User presence table
CREATE TABLE IF NOT EXISTS user_presences (
    user_id BIGINT PRIMARY KEY,
    status VARCHAR(50) NOT NULL DEFAULT 'offline', -- 'online', 'away', 'busy', 'offline'
    status_text VARCHAR(255) DEFAULT '',
    last_activity BIGINT NOT NULL,
    
    INDEX idx_status (status),
    INDEX idx_last_activity (last_activity)
);

-- Message reactions table (for emoji reactions)
CREATE TABLE IF NOT EXISTS message_reactions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    message_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    emoji VARCHAR(100) NOT NULL,
    ctime BIGINT NOT NULL,
    
    UNIQUE KEY unique_message_user_emoji (message_id, user_id, emoji),
    INDEX idx_message_id (message_id),
    INDEX idx_user_id (user_id)
);

-- File attachments table
CREATE TABLE IF NOT EXISTS message_attachments (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    message_id BIGINT NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    file_url VARCHAR(500) NOT NULL,
    file_type VARCHAR(100) NOT NULL,
    file_size BIGINT NOT NULL,
    ctime BIGINT NOT NULL,
    
    INDEX idx_message_id (message_id)
); 