-- migrate:up
CREATE TABLE users(
  id VARCHAR(20) PRIMARY KEY,
  email VARCHAR(255) NOT NULL UNIQUE,
  password VARCHAR(255) NOT NULL,
  username VARCHAR(255) NOT NULL UNIQUE,
  display_name VARCHAR(255) NOT NULL,
  avatar VARCHAR(255),
  banner VARCHAR(255),
  about_me JSONB,
  experience BIGINT NOT NULL DEFAULT 0,
  main_color VARCHAR(255),
  links JSONB DEFAULT '[]'::jsonb,
  facts JSONB DEFAULT '[]'::jsonb,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

CREATE TABLE friends(
  id VARCHAR(20) PRIMARY KEY,
  sender_id VARCHAR(255) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  receiver_id VARCHAR(255) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  accepted BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);


CREATE TABLE servers(
  id VARCHAR(20) PRIMARY KEY,
  owner_id VARCHAR(20) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  name VARCHAR(255) NOT NULL,
  avatar VARCHAR(255),
  banner VARCHAR(255),
  description JSONB,
  public BOOLEAN NOT NULL DEFAULT FALSE,
  main_color VARCHAR(255),
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

CREATE TABLE channels(
  id VARCHAR(20) PRIMARY KEY,
  server_id VARCHAR(20) NOT NULL REFERENCES servers(id) ON DELETE CASCADE,
  name VARCHAR(255) NOT NULL,
  description TEXT,
  type VARCHAR(255) NOT NULL,
  e2ee BOOLEAN NOT NULL DEFAULT TRUE,
  users VARCHAR(20) ARRAY,
  roles VARCHAR(20) ARRAY,
  active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

CREATE TABLE roles(
  id VARCHAR(20) PRIMARY KEY,
  server_id VARCHAR(20) NOT NULL REFERENCES servers(id) ON DELETE CASCADE,
  position INT NOT NULL DEFAULT 0,
  name VARCHAR(255) NOT NULL,
  color VARCHAR(255) NOT NULL,
  abilities VARCHAR(255) ARRAY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

CREATE TABLE messages(
  id VARCHAR(20) PRIMARY KEY,
  author_id VARCHAR(20) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  server_id VARCHAR(20) NOT NULL REFERENCES servers(id) ON DELETE CASCADE,
  channel_id VARCHAR(20) NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
  content JSONB NOT NULL,
  everyone BOOLEAN NOT NULL DEFAULT FALSE,
  mentions_users VARCHAR(20) ARRAY,
  mentions_channels VARCHAR(20) ARRAY,
  attachments JSONB DEFAULT '[]',
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

CREATE TABLE server_members(
  id VARCHAR(20) PRIMARY KEY,
  user_id VARCHAR(20) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  server_id VARCHAR(20) NOT NULL REFERENCES servers(id) ON DELETE CASCADE,
  roles VARCHAR(20) ARRAY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
  UNIQUE(user_id, server_id)
);

CREATE TABLE invites(
  id VARCHAR(20) PRIMARY KEY,
  server_id VARCHAR(20) NOT NULL REFERENCES servers(id) ON DELETE CASCADE,
  invite_id VARCHAR(255) NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
  expire_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

CREATE TABLE emojis(
  id VARCHAR(20) PRIMARY KEY,
  user_id VARCHAR(20) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  shortcode VARCHAR(255) NOT NULL,
  url VARCHAR(255) NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
  expire_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

CREATE TABLE user_channel_read_state(
  user_id VARCHAR(20) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  channel_id VARCHAR(20) NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
  last_read_message_id VARCHAR(20) REFERENCES messages(id) ON DELETE CASCADE,
  unread_mention_ids JSONB DEFAULT '[]'::jsonb,
  updated_at TIMESTAMP DEFAULT NOW(),
  PRIMARY KEY (user_id, channel_id)
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_invites_invite_id ON invites(invite_id);

-- dumb user to create the global server
INSERT INTO users
VALUES ('global', 'admin', 'admin', '', '');

-- dumb server to create dm channels for friends
INSERT INTO servers
VALUES ('global', 'global', 'global');

-- migrate:down
DROP TABLE emojis;
DROP TABLE user_channel_read_state;
DROP TABLE invites;
DROP TABLE server_members;
DROP TABLE messages;
DROP TABLE roles;
DROP TABLE channels;
DROP TABLE servers;
DROP TABLE friends;
DROP TABLE users;
