CREATE TABLE IF NOT EXISTS members(
    code VARCHAR(32) PRIMARY KEY,
    name VARCHAR(32) NOT NULL,
    english_name VARCHAR(32) NOT NULL,
    kana VARCHAR(32) NOT NULL,
    cate VARCHAR(4) NOT NULL,
    img VARCHAR(255) NOT NULL,
    link VARCHAR(255) NOT NULL,
    pick VARCHAR(32),
    god VARCHAR(32),
    under VARCHAR(32),
    birthday VARCHAR(32) NOT NULL,
    blood VARCHAR(32) NOT NULL,
    constellation VARCHAR(32) NOT NULL,
    graduation VARCHAR(4) NOT NULL,
    group_name VARCHAR(32)
);

CREATE TABLE IF NOT EXISTS blogs(
    code VARCHAR(32) PRIMARY KEY,
    title text NOT NULL,
    content text NOT NULL,
    timestamp VARCHAR(255) NOT NULL,
    latest_comment_timestamp VARCHAR(32),
    start_time VARCHAR(255),
    end_time VARCHAR(255),
    cate VARCHAR(255),
    link VARCHAR(255) NOT NULL,
    member_code VARCHAR(32) NOT NULL,
    CONSTRAINT fk_member_code
        FOREIGN KEY(member_code)
            REFERENCES members(code)
);

CREATE TABLE IF NOT EXISTS blog_images(
    url text NOT NULL,
    path VARCHAR(255),
    blog_code VARCHAR(255) NOT NULL,
    CONSTRAINT fk_blog_code
        FOREIGN KEY(blog_code)
            REFERENCES blogs(code)
);

CREATE TABLE IF NOT EXISTS comments(
    code VARCHAR(32) PRIMARY KEY,
    blog_code VARCHAR(32) NOT NULL,
    content text NOT NULL,
    name text NOT NULL,
    timestamp VARCHAR(32) NOT NULL,
    CONSTRAINT fk_blog_code
        FOREIGN KEY(blog_code)
            REFERENCES blogs(code)
);
