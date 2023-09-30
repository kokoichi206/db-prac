``` sql
SELECT code, timestamp, link, AVG(LENGTH(content)) OVER (PARTITION BY member_code) FROM blogs;
```


``` sql
SELECT m.name, b.link, title, AVG(LENGTH(title)) OVER (PARTITION BY member_code), AVG(LENGTH(content)) OVER (PARTITION BY member_code) FROM blogs b INNER JOIN members m ON m.code = b.member_code;


SELECT m.name, b.link, AVG(LENGTH(title)) OVER (PARTITION BY member_code), AVG(LENGTH(content)) OVER (PARTITION BY member_code) FROM blogs b INNER JOIN members m ON m.code = b.member_code;

SELECT m.name, AVG(LENGTH(title)) AS title_length, AVG(LENGTH(content)) AS content_length, COUNT(*) AS num_blogs
FROM blogs b INNER JOIN members m ON m.code = b.member_code
GROUP BY m.code
ORDER BY title_length DESC;

SELECT m.name, AVG(LENGTH(title)) AS title_length, AVG(LENGTH(content)) AS content_length, COUNT(*) AS num_blogs, SUM(LENGTH(content)) AS content_length_total
FROM blogs b INNER JOIN members m ON m.code = b.member_code
GROUP BY m.code
ORDER BY content_length DESC;


SELECT m.name, AVG(LENGTH(title)) AS title_length, AVG(LENGTH(content)) AS content_length, COUNT(*) AS num_blogs, SUM(LENGTH(content)) AS content_length_total
FROM blogs b INNER JOIN members m ON m.code = b.member_code
GROUP BY m.code
ORDER BY content_length_total DESC;

SELECT title FROM blogs WHERE member_code = '48006' ORDER BY LENGTH(blogs.title);


SELECT LENGTH(blogs.content), content, TRIM(blogs.content) AS trimmed_content, link FROM blogs ORDER BY LENGTH(blogs.content) LIMIT 3;

SELECT LENGTH(TRIM(blogs.content)), content, link FROM blogs ORDER BY LENGTH(TRIM(blogs.content)) LIMIT 15;

CREATE TABLE fans (
    name text NOT NULL,
)
```

``` sql
SELECT m.name, b.link, title, AVG(LENGTH(title)) OVER (PARTITION BY member_code), AVG(LENGTH(content)) OVER (PARTITION BY member_code) FROM blogs b INNER JOIN members m ON m.code = b.member_code;

SELECT SUM(LENGTH(content)) OVER (PARTITION BY member_code) FROM blogs;

WITH TEMP ( total ) AS ( SELECT SUM(LENGTH(content)) FROM blogs )
SELECT member_code, ROUND ( COUNT(LENGTH(content))*100 / total, 3 )
FROM blogs CROSS JOIN TEMP
GROUP BY member_code, total; -- total を忘れない！


WITH TEMP ( total ) AS ( SELECT SUM(LENGTH(content)) FROM blogs )
SELECT member_code, ROUND ( SUM(LENGTH(content))*100 / total, 3 )
FROM blogs CROSS JOIN TEMP
GROUP BY member_code, total; -- total を忘れない！

WITH TEMP ( total ) AS ( 
    SELECT SUM(LENGTH(content))::FLOAT FROM blogs 
)
SELECT member_code, ROUND( ((SUM(LENGTH(content))::FLOAT * 100) / total)::NUMERIC, 3 ) AS percentage
FROM blogs CROSS JOIN TEMP
GROUP BY member_code, total; -- total を忘れない！


WITH TEMP ( total ) AS ( 
    SELECT SUM(LENGTH(content))::FLOAT FROM blogs 
)
SELECT member_code, SUM(LENGTH(content)), total
FROM blogs CROSS JOIN TEMP
GROUP BY member_code, total; -- total を忘れない！
```

### 関数？

ROW_NUMBER()

``` sql
SELECT code, SUBSTR(title, 0, 30), ROW_NUMBER() OVER (ORDER BY LENGTH(content) DESC) AS rank FROM blogs;

SELECT code, SUBSTR(title, 0, 30), LENGTH(content), ROW_NUMBER() OVER (ORDER BY LENGTH(content) DESC) AS rank FROM blogs;

SELECT code, link, SUBSTR(title, 0, 30), LENGTH(content), ROW_NUMBER() OVER (ORDER BY LENGTH(content) DESC) AS rank FROM blogs;
```

RANK()

各行にランクを割り当てます。同じ値の行は同じランクを持ち、次のランクはスキップされます。

ROW_NUMBER は同じだったとしても 1,2,3,... って上がっていく

``` sql
-- cf
SELECT code, SUBSTR(title, 0, 30), LENGTH(code), ROW_NUMBER() OVER (ORDER BY LENGTH(code) DESC) AS rank FROM blogs;

SELECT code, SUBSTR(title, 0, 30), LENGTH(code), RANK() OVER (ORDER BY LENGTH(code) DESC) AS rank FROM blogs;
```

DENSE_RANK()

各行にランクを割り当てますが、RANK()とは異なり、同じ値の行は同じランクを持ち、次のランクはスキップされません。

RANK は 1,1,1,1,5,5,5, ってなるけど、DENSE_RANK は 1,1,1,1,2,2,2,.. だわ

``` sql
-- cf
SELECT code, SUBSTR(title, 0, 30), LENGTH(code), RANK() OVER (ORDER BY LENGTH(code) DESC) AS rank FROM blogs;

SELECT code, SUBSTR(title, 0, 30), LENGTH(code), DENSE_RANK() OVER (ORDER BY LENGTH(code) DESC) AS rank FROM blogs;
```

LAG() and LEAD():

前の行や次の行の値を取得します。

``` sql
-- ORDER BY が被ることは許されない、エラーにはならないが、最後に指定したのになってる？？
SELECT code, SUBSTR(title, 0, 30), LENGTH(content), 
LAG(LENGTH(content)) OVER (ORDER BY LENGTH(content) DESC),
LEAD(LENGTH(content)) OVER (ORDER BY LENGTH(content)) 
FROM blogs;

SELECT code, SUBSTR(title, 0, 30), LENGTH(content), 
LAG(LENGTH(content)) OVER (ORDER BY LENGTH(content) DESC) AS next_length,
LEAD(LENGTH(content)) OVER (ORDER BY LENGTH(content) DESC) AS prev_length
FROM blogs;
```

SUM() with PARTITION BY

``` sql
SELECT members.name, SUBSTR(title, 0, 30), LENGTH(content), 
SUM(LENGTH(content)) OVER (PARTITION BY member_code) AS sum_content_length
FROM blogs
INNER JOIN members ON blogs.member_code = members.code
ORDER BY sum_content_length;

WITH RankedContent AS (
    SELECT 
        members.name, 
        SUBSTR(title, 0, 30) AS shortened_title, 
        LENGTH(content) AS content_length,
        SUM(LENGTH(content)) OVER (PARTITION BY member_code) AS sum_content_length,
        ROW_NUMBER() OVER (PARTITION BY members.name ORDER BY LENGTH(content) DESC) AS row_num
    FROM blogs
    INNER JOIN members ON blogs.member_code = members.code
)
SELECT name, shortened_title, content_length, sum_content_length
FROM RankedContent
WHERE row_num = 1
ORDER BY sum_content_length;
```
