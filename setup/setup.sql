CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    username VARCHAR(255) NOT NULL,
    password_hash VARCHAR(1023) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    is_private BOOLEAN DEFAULT false,
    bio TEXT,
    profile_image_url VARCHAR(1023)
);


CREATE TABLE follows (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    follower_id INT NOT NULL,
    CONSTRAINT different_user_and_follower CHECK (user_id != follower_id),
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_follower FOREIGN KEY (follower_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT unique_user_follower_pair UNIQUE (user_id, follower_id)
);

CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    title VARCHAR(255),
    content TEXT,
    user_id INT,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE media (
    id SERIAL PRIMARY KEY,
    url VARCHAR(255),
    post_id INT,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
);

CREATE TABLE comments (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    post_id INT NOT NULL,
    user_id INT NOT NULL,
    content TEXT,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE post_likes (
    id SERIAL PRIMARY KEY,
    user_id INT,
    post_id INT,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    CONSTRAINT unique_post_user_pair UNIQUE (post_id, user_id)
);

CREATE TABLE comment_likes (
    id SERIAL PRIMARY KEY,
    user_id INT,
    comment_id INT,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (comment_id) REFERENCES comments (id) ON DELETE CASCADE,
    CONSTRAINT unique_comment_user_pair UNIQUE (comment_id, user_id)
);


