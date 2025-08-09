-- Criar tabela sessions completa com todos os campos
CREATE TABLE IF NOT EXISTS sessions (
    id VARCHAR NOT NULL,
    name VARCHAR NOT NULL,
    status VARCHAR NOT NULL DEFAULT 'disconnected',
    phone VARCHAR NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- Campos de proxy
    proxy_type VARCHAR NULL,
    proxy_host VARCHAR NULL,
    proxy_port BIGINT NULL,
    proxy_username VARCHAR NULL,
    proxy_password VARCHAR NULL,

    -- Campos WhatsApp
    webhook_url VARCHAR DEFAULT '',
    qrcode TEXT DEFAULT '',
    device_jid VARCHAR DEFAULT '',
    is_connected BOOLEAN DEFAULT FALSE,

    CONSTRAINT sessions_pkey PRIMARY KEY (id),
    CONSTRAINT sessions_name_key UNIQUE (name)
);

-- Criar índices para performance
CREATE INDEX IF NOT EXISTS idx_sessions_status ON sessions(status);
CREATE INDEX IF NOT EXISTS idx_sessions_created_at ON sessions(created_at);
CREATE INDEX IF NOT EXISTS idx_sessions_phone ON sessions(phone) WHERE phone IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_sessions_is_connected ON sessions(is_connected);
CREATE INDEX IF NOT EXISTS idx_sessions_device_jid ON sessions(device_jid) WHERE device_jid != '';
