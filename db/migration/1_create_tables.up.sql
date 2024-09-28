CREATE TABLE posts (
    id  UUID PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    allow_comments BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE comments (
    id UUID PRIMARY KEY,
    post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    parent_id UUID REFERENCES comments(id) ON DELETE CASCADE,
    content TEXT NOT NULL
);