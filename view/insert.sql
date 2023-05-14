INSERT INTO members(
    code,
    name,
    english_name,
    kana,
    cate,
    img,
    link,
    pick,
    god,
    under,
    birthday,
    blood,
    constellation,
    graduation,
    group_name
) VALUES (
    '1',
    'john doe',
    'john doe',
    'じょん',
    '1st',
    'john.jpg',
    'https://www.instagram.com/john.doe/',
    '0',
    '0',
    '0',
    '1999-03-31',
    'A',
    'おひつじ座',
    'no',
    'xxx'
), (
    '2',
    'john doe2',
    'john doe2',
    'じょん2',
    '2nd',
    'john.jpg',
    'https://www.instagram.com/john.doe2/',
    '0',
    '0',
    '0',
    '1999-03-31',
    'A',
    'おひつじ座',
    'no',
    'xxx'
);


INSERT INTO blogs(
    code,
    title,
    content,
    timestamp,
    latest_comment_timestamp,
    start_time,
    end_time,
    cate,
    link,
    member_code
) VALUES (
    '1',
    'title',
    'content',
    '1684065026000000000',
    '1684065026000000000',
    '2020-01-01 00:00:00',
    '2020-01-01 00:00:00',  
    '1st',
    'https://www.instagram.com/john.doe/',
    '1'
), (
    '2',
    'title2',
    'content',
    '1684065026000000000',
    '1684065026000000000',
    '2020-01-01 00:00:00',
    '2020-01-01 00:00:00',  
    '1st',
    'https://www.instagram.com/john.doe/',
    '1'
), (
    '3',
    'title3',
    'content',
    '1684065026000000000',
    '1684065026000000000',
    '2020-01-01 00:00:00',
    '2020-01-01 00:00:00',  
    '1st',
    'https://www.instagram.com/john.doe/',
    '2'
);



INSERT INTO comments(
    code,
    blog_code,
    content,
    name,
    timestamp
) VALUES (
    '1',
    '1',
    'content1',
    'name1',
    '1684065026000000000'
), (
    '2',
    '1',
    'content2',
    'name2',
    '1684065026000000000'
), (
    '3',
    '1',
    'content3',
    'name3',
    '1684065026000000000'
), (
    '4',
    '2',
    'content4',
    'name1',
    '1684065026000000000'
), (
    '5',
    '3',
    'content5',
    'name1',
    '1684065026000000000'
), (
    '6',
    '3',
    'content6',
    'name2',
    '1684065026000000000'
);


