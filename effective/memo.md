## sec 1
「豚の耳から絹の財布は作れない」データモデルがうまく設計されていない状態で、「効果的な」SQLを書き始めることはできない。

### 1. 全てのテーブルに主キーを定義する
- 一意
- null にならない
- 安定（変更され得ない）
- 単純

参照整合性、nullではない外部キーを持つ小テーブルのレコードごとに、親テーブルに一致するレコードが含まれていなければならない。

最も単純なのは、RDBMS に応じた、自動生成される値を使うこと。テキストベースの主キー（ConpanyName, email 等）を使うことのメリットは SQL がより単純になる。どちらが良いかいまだに議論されている。


### 2. 冗長なデータを取り除く
正規化を行う。

### 3. 繰り返しグループを取り除く
列は高くつく
行は安上がりである

### 4. 列ごとにプロパティを１つ
AuthFirst, AuthMid, AuthLast, AuthCity, AuthStreet, AuthStNum...

``` sql
SELECT AuthorID AS AuthID,
    CONCAT(AuthFirst,
        CASE WHEN AuthMid IS NULL
            THEN ' '
            ELSE CONCAT(' ', AuthMid, ' ')
            END, AuthLast) AS AuthName,
    CONCAT(AuthStNum, ' ', AuthStreet, ' ',
        AuthCity, ', ', ...) AS AuthAddress
FROM Authors;
```

- 検索やグループ化が簡単になる
- どうフィルタリングするかによってデータの粒度が決まる
- 元に戻すには CONCAT を利用する

### 5. 計算値の格納は一般に避ける！
DB のシステムとして計算列を定義できることが多いが、パフォーマンスに悪影響！

### 6. 参照整合性を確保するための外部キー
スキーマを正しく設計すると、テーブルの多くに外部キーが定義される。これらの外部キーには、親テーブルの主キーの値が含まれている。

### 7. テーブル間の関係
今扱ってるのは、構造化データ？それとも半構造化データ？

XML, JSON のような半構造化データに対しては RDB は不向きすぎる。

- 同じような列を含んでいるようにみえるテーブルを結合する意味が本当にあるのか？

### 8. Beyond 3NF
3NF をみたいしている設計が、より上のレベルの正規形に反する可能性があるのは、１つのテーブルが複数のテーブルに関連している時。特に、テーブルが１対他になっていて、そうした関係が１つではない場合。

最初の３つの正規形は、リレーション属性のうち**関数従属性**に関するものである。つまり、その属性がリレーションのキーに依存すること意味する。たとえば、電話番号 xxx-xxx-xxxx を格納している列はその持ち主の John Doe を格納している列に関数従属している、と言える

第四正規形（4NF）は、多値従属性に関係している。２つの属性が互いに独立しているものの、リレーションの同じキーに依存している、というケース

情報無損失分析のテスト

### 9. データウェアハウスでは非正規形
正規化を行うと、データが一箇所に配置されるため更新や挿入が高速になる、データが重複しないため、負荷の高い Group By や Distinct クエリが必要になることもない。RDB で正規化した方が良い理由、それは、アプリケーションが書き込み主体であるため、書き込みの負荷が読み取りの負荷を上回るためである。

正規化されたテーブルの問題点は、テーブル間の結合の数が増えるほど、オプティマイザが効率的な実行プランを見つけるのが難しくなること。データウェアハウスでは、SELECT を高速にするため、非正規化されたデータベースが適している。

書き込みは滅多に発生しないため、インデックスの数が多すぎるために書き込みのパフォーマンスが劇的に低下する、という心配もない。

- 計算値や派生値を格納することも、非正規化の手法の１つ


## sec 2
SQL テーブルをうまく動作させるための主な要素の１つは、テーブルを正しくインデックス付けすること。

### 10. インデックスを作成するときの null の扱い

### 11. データスキャンを最小限に抑える
- インデックスシーク
    - 必要な値が全部インデックスに含まれる
- インデックススキャン
- テーブルスキャン

多くのインデックスは、データ取得の高速化には貢献せず、それどころか更新を低速化させることがあるのを覚えておく。

DBMS で使用されている最も一般的な種類のインデックスは B-tree。

- B-tree
    - クラスタ化インデックス
    - 非クラスタ化インデックス

インデックスが意味を持つのは、テーブルが大きい場合だけ。ほとんどのデータベースエンジンは、テーブルが小さい場合はメモリに読み込んでしまう。


### 12. フィルタリング以外でのインデックス
インデックスは、データ構造とは別物であり、それ専用のディスク容量が必要となる。

b-tree いんでくっすについては、データの**順序付き表現**となっているため、ソートがえぐい。一般にソートは、データを一時的にバッファに格納する必要があるためえぐい。

- WHERE 句に含まれている列がインデックスに含まれているかどうかは、クエリのパフォーマンスに影響を与える

### 13. トリガを使いすぎない
DRI（Declarative Referential Integrity）を使う方が一般に良い

``` sql
ALTER TABLE Order_Details
    ADD CONSTRAINT fkOrder FOREIGN KEY (OrderNumber)
        REFERENCES Orders (OrderNumber) ON DELETE CASCADE;
```

### 14. サブセットの取捨選択にフィルター選択されたインデックス

### 15. 宣言型の制約
- NOT NULL
- UNIQUE KEY
    - PRIMARY と異なり null 値を入力できる
- PRIMARY KEY
    - テーブルに１つのみ
- FOREIGN KEY
- CHECK
    - フィールドまたはテーブルで定義可能
    - 指定した値のみを格納できるようになる
- DEFAULT

ビジネスルールを適応することとデータの関係を維持することはデータモデルの一部であり、その実行責任はアプリケーションではなくデータベースにある！！

クエリオプティマイザ

### 16. SQLダイアレクト
- null
    - 2つのnull値の順序は等しいとみなされる。
    - null値をソートした時の順序は、null以外のすべての値の前または後ろにすべきである。
        - どっちかはRDBMSによる！
- 結果セットの制限
    - FETCH FIRST
    - ウィンドウ関数
    - カーソル
- BOOLEAN (SQLの規格)
    - TRUE
    - FALSE
    - UNKNOWN または NULL
- 常に DBMS のマニュアルを調べる！

### 17. 計算値をインデックスで使用する？
``` sql
SELECT EmployeeID, EmpFirstName, EmpLastName
FROM Employees
WHERE UPPER(EmpLastName) = 'John';
```

これだとインデックスが効かない！！
そこで次のようにするう

``` sql
CREATE INDEX EmpLastNameUpper
    ON Employees (UPPER(EmpLastName))
```


## sec 3

### 18. 変更できないものはビューを使用
view とは、あらかじめ定義された SQL クエリの結果として合成されるテーブルのこと。

### 19. ETL
ETL（Extract, Transform, Load）を使って非リレーショナルデータを情報に変える

### 20. サマリーテーブルの活用
詳細テーブルのデータを集計するサマリーテーブルを作成し、トリガーを定義する。

### 21. UNION で非正規化データをアンピボット
和演算 union, はリレーショナルモデルで実行可能な８つの関係代数演算の１つである。和演算は、２つ（以上）のSELECT文によって作成されたデータセットをマージするために使用される

``` sql
SELECT Category, 'Oct' AS SalesMonth, OctQuantity As Quantity, OctSales AS SalesAmt
FROM SalesSummary
UNION
SELECT Category, 'Nov', NovQuantity, NovSales
FROM SalesSummary
UNION
SELECT Category, 'Dec', DecQuantity, DecSales
FROM SalesSummary
UNION
SELECT Category, 'Jan', JanQuantity, JanSales
FROM SalesSummary
UNION
SELECT Category, 'Feb', FebQuantity, FebSales
FROM SalesSummary
ORDER BY SalesMonth, Category;
```


## sec 4

### 22. 関係代数（relational algebra）
- 選択（制限）
    - 行にフィルターを適応してサブセットを取得する演算
    - WHERE, HAVING
- 射影
    - 返される列を選択
- 結合
    - 全てのテーブル（リレーション）に一意な識別子が定義されてなければならない
    - 関連元のテーブルに関連先のテーブルの一意な識別子（外部キー）のコピーが含まれてなければならない
- 交差
    - まったく同じ列を持つ２つのテーブルで実行
    - INTERSECT
    - INNER JOIN 等で同様の結果を得ることが可能
- 直積
    - Cartesian product
    - １つ目のテーブルの全ての行を、２つ目のテーブルの全ての行と組み合わせた結果
- 和
    - union
    - まったく同じ列を持つ２つのテーブルをマージする演算
- 商
    - 商演算をサポートしているデータベースシステムは１つもない
- 差
    - 一方のテーブルからもう一方のテーブルを差し引く演算

``` sql
-- Skateboard を注文したが Helmet を注文していない顧客リスト
SELECT c.CustFirstName, c.CustLastName
FROM Customrs AS c
WHERE c.CustomerID IN
    (SELECT o.CustomerID
    FROM Orders AS o
        INNER JOIN Order_Details AS od
            ON o.OrderNumber = od.OrderNumber
        INNER JOIN Products AS p
            ON p.ProductNumber = od.ProductNumber
    WHERE p.ProductName = 'Skateborad')
EXCEPT
SELECT c2.CustFirstName. c2.CustLastName
FROM Customrs AS c2
WHERE c.CustomerID IN
    (SELECT o.CustomerID
    FROM Orders AS o
        INNER JOIN Order_Details AS od
            ON o.OrderNumber = od.OrderNumber
        INNER JOIN Products AS p
            ON p.ProductNumber = od.ProductNumber
    WHERE p.ProductName = 'Helmet');
```

### 23. 条件と一致しないレコードや欠けているレコード特定
``` sql
-- 購入されていない製品を特定
SELECT p.ProductNumber, p.ProductName
FROM Products AS p
WHERE p.ProductNusmber
    NOT IN (SELECT ProductNumber FROM Order_Details);
```

このクエリの実行コストはかなり高い！
Order_Details の全てのレコードにアクセスし重複値を取り除いた上で、Products テーブルで照合しなければならないから。

理論的には、EXISTSを使用する方がNOT INを使用するよりも高速なはず。
クエリエンジンが最初の行を見つけた時点で、サブクエリの処理を終了できるから。

``` sql
SELECT p.ProductNumber, p.ProductName
FROM Products AS p
WHERE NOT EXISTS
    (SELECT *
    FROM Order_Details AS od
    WHERE od.ProductNumber = p.ProductNumber);
```

別のアプローチとしては、LEFT JOIN 演算子と null 値を検索する WHERE の組み合わせ。

``` sql
SELECT p.ProductNumber, p.ProductName
FROM Products AS p LEFT JOIN Order_Details AS od
    ON p.ProductNumber = od.ProductNumber
WHERE od.ProductNumber IS NULL;
```

どのアプローチが適してるかの明確な答えはない！


### 24. CASE
``` sql
CASE Students.Gender
    WHEN 'M' THEN 'Male'
    WHEN 'F' THEN 'Female'
    ELSE 'Unknown' END

CASE Readings.Measure
    WHEN 'C'
    THEN 2b
```

### 25. 複数の条件を使用する問題の解決
``` sql
CREATE FUNCTION CustProd(@ProdName varchar(50)) RETURNS Table
AS
RETURN
    (SELECT Orders.CustormerID AS CustID
    FROM Orders
        INNER JOIN Order_Details
            ON Orders.OrderNumber = Order_Details.OrderNumber
        INNER JOIN Products
            ON Products.ProductNumber
    WHERE ProductName = @ProdName);

SELECT C.CustomerID, C.CustFirstName, C.CustLastName
FROM Customers AS C
WHERE C.CustomerID IN
    (SELECT CustID FROM CustProd('Skateboard'))
AND C.CustomerID IN
    (SELECT CustID FROM CustProd('Helmet'))
AND C.CustomerID IN
    (SELECT CustID FROM CustProd('Kee Pads'))
AND C.CustomerID IN
    (SELECT CustID FROM CustProd('Gloves'));
```

### 26. 完全に一致させる場合のデータ分割
商演算によって解決できる一般的な問題

- 募集条件を全て満たしている応募者の検索
- コンポーネントを構築するための部品を全て供給できるサプライヤーのリストアップ
- 特定の製品を注文した顧客全員の表示

### 27. 日付と時刻の列で日付の範囲を正しくフィルタリング
``` sql
CREATE TABLE ProgramLogs (
    LogID int PRIMARY KEY,
    LogUserID varchar(20) NOT NULL,
    LogDate timestamp NOT NULL,
    Logger varchar(50) NOT NULL,
    LogLevel varchar(10) NOT NULL,
    LogMessage varchar(1000) NOT NULL
);

INSERT INTO ProgramLogs (LogID, LogUserID, LogDate, Logger, LogLevel, LogMessage) VALUES (1, 'Doug', '2016-07-04 09:15:32', 'ABC', '1', 'Sorry, something went wrong. A team of highly trained monkeys has been dispatched to deal with this situation.');
INSERT INTO ProgramLogs (LogID, LogUserID, LogDate, Logger, LogLevel, LogMessage) VALUES (2, 'Ben', '2016-07-04 11:23:12', 'BCD', '2', 'One of us has made a mistake, and I''m not pressing the button.');
INSERT INTO ProgramLogs (LogID, LogUserID, LogDate, Logger, LogLevel, LogMessage) VALUES (3, 'John', '2016-07-04 13:54:02', 'ABC', '1', 'Your computer has performed an illegal operation and will be shut down. 911 has been called.');
INSERT INTO ProgramLogs (LogID, LogUserID, LogDate, Logger, LogLevel, LogMessage) VALUES (4, 'Doug', '2016-07-04 15:03:23', 'XYZ', '2', 'Something went wrong. You''re on your own.');
INSERT INTO ProgramLogs (LogID, LogUserID, LogDate, Logger, LogLevel, LogMessage) VALUES (5, 'Doug', '2016-07-04 23:58:02', 'EFG', '2', 'You really screwed up this time.');
INSERT INTO ProgramLogs (LogID, LogUserID, LogDate, Logger, LogLevel, LogMessage) VALUES (6, 'Doug', '2016-07-04 23:58:12', 'CDE', '4', 'Run away as fast as you can, and don''t look back.');
INSERT INTO ProgramLogs (LogID, LogUserID, LogDate, Logger, LogLevel, LogMessage) VALUES (7, 'Ben', '2016-07-05', 'ABC', '3', 'It''s been a while since an error was logged. System will now crash.');
INSERT INTO ProgramLogs (LogID, LogUserID, LogDate, Logger, LogLevel, LogMessage) VALUES (8, 'Ben', '2016-07-05 00:03:35', 'EFG', '4', 'User error. Please replace user.');
INSERT INTO ProgramLogs (LogID, LogUserID, LogDate, Logger, LogLevel, LogMessage) VALUES (9, 'John', '2016-07-05 08:10:02', 'EFG', '3', 'An unknown error has occurred. The error is unknown because the guy who wrote this part of the code quit a while back and he was like real real smart and the rest of us aren''t sure how it works.');
INSERT INTO ProgramLogs (LogID, LogUserID, LogDate, Logger, LogLevel, LogMessage) VALUES (10, 'Doug', '2016-07-05 12:32:01', 'XYZ', '4', 'User error. It''s not our fault!');

-- 特定の日のログメッセージを表示したい
SELECT L.LogUserID, L.Logger, L.LogLevel, L.LogMessage
FROM ProgramLogs AS L
WHERE L.LogDate = CAST('7/4/2016' AS timestamp);

-- postgresql
SELECT L.LogUserID, L.Logger, L.LogLevel, L.LogMessage
FROM ProgramLogs AS L
WHERE L.LogDate BETWEEN CAST('2016-07-04' AS timestamp) 
    AND CAST('2016-07-04 23:59:59.999' AS timestamp);

-- ユーザーによる日付指定を可能に
PREPARE test(date, date) AS
SELECT L.LogUserID, L.Logger, L.LogLevel, L.LogMessage
FROM ProgramLogs AS L
WHERE L.LogDate >= $1
    AND L.LogDate < ($2 + INTERVAL '1 DAYS');
EXECUTE test('2016-07-04', '2016-07-04');
DEALLOCATE test;
```

- インデックスによった検索が不可能になるため、datetime 型の列には関数を使用しない
- 丸誤差が原因で datetime 型の値が正確でなくなることがあるため、BETWEEN ではなく、>= と < を利用する

### 28. インデックスが使用されるようなクエリ
WHERE, ORDER BY, GROUP BY, HAVING といったクエリの述語が sargable (Search ARGument ABLE) である必要がある。

- 一般に『sargable』
    - =
    - \>
    - <
    - \>=
    - <=
    - BETWEEN
    - LIKE (先頭ワイルドカードなし)
    - IS [NOT] NULL
- わんちゃん『sargable』だが、ほとんどの場合パフォーマンスが向上しない
    - \<>
    - IN
    - OR
    - NOT IN
    - NOT EXISTS
    - NOT LIKE
- 『sargable』ではないクエリになる
    - １つ以上のフィールド操作する関数を WHERE 句の条件で使用する
    - WHERE 句でフィールドの値を使って算術演算を実行する場合
    - LIKE '%something%' のようなワイルドカード検索を使用する場合

``` sql
-- インデックスが使用されない
SELECT *
FROM Employees
WHERE YEAR(EmpDOB) = 1990;

-- インデックスが使用される
SELECT *
FROM Employees
WHERE EmpDOB >= CAST('1990-01-01' AS Date)
    AND EmpDOB < CAST('1991-01-01' AS Date);
```

ワイルドカードが文字列の末尾にしかない場合、インデックスが使用される場合がある（が、このことはインデックスが使用されるという保証にはならない）

### 左結合の右側でフィルタリングを正しく行う
- 差演算を実行するには、OUTER JOIN


## sec 5
より複雑な集約への需要の高まりを受け、規格の拡張に着手した結果がウィンドウ関数。

### 30. GROUP BY のしくみ
- GROUP BY により、フィルタリング後のデータセットが集約される
- HAVING 句により、フィルタリング後の集約されたデータセットが変換される

ISO SQL 規格で定義されている集約関数のうち、よく使用される９つ

- COUNT()
- SUM()
- ARG()
- MIN()
- MAX()
- STDDEV_POP(), STDDEV_SAMP()
- VAR_POP(), VAR_SAMP()

SELECT 句に含まれていて、集約関数が適応されていない例は、GROUP BY 句に含まれていなければならない。

ROLLUP, CUBE, GROUPING SETS,

- WHERE は集約が実行される前に適応される
- SELECT 句に指定された列のうち、集約関数や計算に含まれていない列は、GROUP BY 句に含まれていなければならない
- ROLLUP, CUBE, GROUPING SETS を利用すれば、複数の集約クエリを UNION で組み合わせる代わりに、１つのクエリをまとめることができる

### 31. GROUP BY を短く保つ
SQL/99 以降では、関数従属性（functional dependencyl）が認識されているため、「集約されない列はすべて GROUP BY 句に含まれてなければならない」ということはない

- データを正しく集約するために必要な列だけが GROUP BY 句に含まれるようにする

### 32. GROUP BY, HAVING
- グループ化を実行する前の行のフィルタリングには、WHERE 句を使用する。
- グループ化を実行した後の行のフィルタリングには、HAVING 句を使用する。

### 33. GROUP BY を使用せずに最大値や最小値の特定
``` sql
SELECT l.Category, l.MaxABV AS LeftMaxABV,
        r.MaxABV AS RightMaxABV
FROM BeerStyles AS l
    LEFT JOIN BeerStyles AS r
        ON l.Category = r.Category
            AND l.MaxABV < r.MaxABV;

-- 
SELECT l.Category, l.Country, l.Style, l.MaxABV AS MaxAlcohol
FROM BeerStyles AS l
    LEFT JOIN BeerStyles AS r
        ON l.Category = r.Category
            AND l.MaxABV < r.MaxABV
WHERE r.MaxABV IS NULL
ORDER BY l.Category;
```

- LEFT JOIN を使って「メイン」テーブルをそれ自体に結合する必要がある


### 34. OUTER JOIN を使用するときは COUNT() に注意
- null 値を含んでいる行を含め、全ての行をカウントしたい場合は、`COUNT(*)` を使用
- 列の値が NULL ではない行だけをカウントしたい場合は、COUNT(\<列名\>)を使用
    - OUTER JOIN の場合は NULL が出るんで注意
- ちょいとばす

### 35. HAVING COUNT(x) < N で値が0の行もカウントする
- カウントが０の検索は、INNER JOIN を使用する場合はうまくいかない

### 36. 重複なしのカウントに DISTINCT
``` sql
SELECT COUNT(CASE WHEN OrderTotal > 1000 THEN CustomerID END) AS TotalOrders
FROM Orders;
```

- COUNT() の引数として関数を使用することを検討する
    - WHERE 句がなくても計算を組み合わせることが可能に

### 37. ウィンドウ関数
SQL:2003 策定以前、SQLにはそもそも「隣接する行」という概念がなかった。
累積和の生成など。

「ウィンドウ」は、該当の行の前後にある一連の行を表す。

``` sql
SELECT o.OrderNumber, o.CustomerID, o.OrderTotal,
    SUM(o.OrderTotal) OVER (
        PARTITION BY o.CustomerID
        ORDER BY o.OrderNumber, o.CustomerID
    ) AS TotalByCustomer,
    SUM(o.OrderTotal) OVER (
        ORDER BY o.OrderNumber
    ) AS TotalOverall
FROM Orders AS o
ORDER BY o.OrderNumber, o.CustomerID
```

- OVER 句はウィンドウを使用することを表している
- ウィンドウ関数は行の範囲を「認識」するため、従来の集計関数や文レベルのグループ化を使用する場合よりも、累積計算や移動集計の生成が容易になる


### 38. 行をランク付けする
- Skip

### 39. 移動集計を生成
- Skip


## sec 6
一般に、テーブル名を使用できる場所では、サブクエリを使用できる。

### 40. サブクエリの使用可能場所
テーブルサブクエリ。

単一列のテーブルサブクエリ。


``` sql
SELECT Products.ProductName
FROM Products
WHERE Products.ProductNumber NOT IN
    (
        SELECT Order_Details.ProductNumber
        FROM Orders
            INNER JOIN Order_Details
                ON Order.OrderNumber = Order_Details.OrderNumber
        WHERE Orders.OrderDate BETWEEN '2015-12-01' AND '2015-12-31'
    );
```

スカラーサブクエリ。

``` sql
SELECT Products.ProductNumber, Products.ProductName,
    (
        SELECT MAX(Orders.OrderDate)
        FROM Orders
            INNER JOIN Order_Details
                ON Orders.OrderNumber = Order_Details.OrderNumber
        WHERE Order_Details.ProductNumber = Products.ProductNumber
    ) AS LastOrder
FROM Products;
```

``` sql
SELECT Vendors.VendName,
    AVG(Product_Vendors.DaysToDeliver) AS AvgDelivery
FROM Vendors
    INNER JOIN Product_Vendors
        ON Vendors.VendorID = Product_Vendors.VendorID
GROUP BY Vendors.VendName
HAVING AVG(Product_Vendors.DaysToDeliver) >
    (
        SELECT AVG(DaysToDeliver) FROM Product_Vendors
    );
```

### 41. 相関サブクエリと非相関サブクエリ
相関サブクエリ。
WHERE または HAVING 句でフィルターを１つ以上使用する。これらのフィルターは、外側のクエリによって提供される値に依存する。この依存性により、サブクエリは外側のクエリと「相互関係」にある。そのため、クエリによって返される行ごとにサブクエリを実行する必要があり、効率が悪い可能性がある（常にそうとは限らない）。

FROM 句のデータセットの１つとして相関サブクエリを使用することは考えにくい。代わりに JOIN を使用する方が単純明快だから。

``` sql
SELECT Recipe_Classes.RecipeClassDescription,
    (
        SELECT COUNT(*)
        FROM Recipes
        WHERE Recipes.RecipeClassID = Recipe_Classes.RecipeClassID
    ) AS RecipeCount
FROM Recipe_Classes;
```

実は相関サブクエリはほとんどのデータベースシステムにおいて最適化されるため、パフォーマンスは悪くはない。

### 42. できるだけサブクエリより CTE
Comman Table Expression。WITH 句を使って定義。

``` sql
WITH CustProd AS
    (
        SELECT Orders.CustomerID, Products.ProductName
        FROM Orders
            INNER JOIN Order_Details
                ON Orders.OrderNumber = Order_Details.ProductNumber
    ),
    SkateboardOrders AS
    (
        SELECT DISTINCT CustomerID
        FROM CustProd
        WHERE ProductName = 'Skateboard'
    )

SELECT c.CustomerID, c.CustFirstName
FROM Customers AS c
    INNER JOIN SkateboardOrders AS OSk
        ONc.CustomerID = OSk.CustomerID;
```

クエリが大幅に短くなり、可読性の向上につながる。

再帰的 CTE

``` sql
-- 1 ~ 100 の数字のリストを生成
WITH SeqNumTbl AS
    (
        SELECT 1 AS SeqNum
        UNION ALL
        SELECT SeqNum + 1
        FROM SeqNumTbl
        WHERE SeqNum < 100
    )

SELECT SeqNum
FROM SeqNumTbl;
```

``` sql
WITH MgrEmpts (
    ManagerID, ManagerName, EmployeeID, EmployeeName, EmployeeLevel) AS
    (
        SELECT ManagerID, CAST(' ' AS varchar(50)), EmployeeID,
            CAST(CONCAT(EmpFirstName, ' ', EmpLastName) AS varchar(50)),
            0 AS EmployeeLevel
        FROM Employees
        WHERE ManagerID IS NULL
        UNION ALL
        SELECT e.ManagerID, d.EmployeeName, e.EmployeeID,
            CAST(CONCAT(EmpFirstName, ' ', EmpLastName) AS varchar(50)),
            EmployeeLevel + 1
        FROM Employees AS e
            INNER JOIN MgrEmps AS d
                ON e.ManagerID = d.EmployeeID
    )

SELECT ManagerID, MangerName, EmployeeID, EmployeeName, EmployeeLevel
FROM MgrEmps
ORDER BY ManagerID;
```

- CTE の利用で、同じサブクエリを複数回使用する複雑なクエリを単純化できる
- CTE の利用で、関数を使用する必要はなくなる
    - 関数は誤って変更される可能性あり

### 43. サブクエリより結合の方が効率的？
データベースエンジンは大抵結合を最適化できる！


## sec 7

### 44. クエリアナライザ
DBMS にはそれぞれ機能が異なる部分があり、あるアプローチが別の DBMS ではうまくいかないことがある。それではどのようにパフォーマンスを改善したら良いのだろうか？

DBMS で SQL 文が実行されるとき、その前に SQL 文をもっとも効果的に実行する方法が DBMS のオプティマイザによって決定される。

- PostgreSQL
    - EXPLAIN
    - ANALYZE
    - VERBOSE
    - COSTS
    - BUFFERS
    - TIMING
    - FORMAT
- バインドパラメータを含んでいる SQL に対して EXPLAIN を使用する。
- pgAdmin ツールが　GUI ツールとして便利

### 45. DB のメタデータ
``` sql
SELECT t.TABLE_NAME, t.TABLE_TYPE
FROM INFORMATION_SCHEMA.TABLES AS t
WHERE t.TABLE_TYPE IN ('BASE TABLE', 'VIEW');

-- 主キーが定義されていないテーブルのリストを取得
SELECT t.TABLE_NAME
FROM
    (
        SELECT TABLE_NAME
        FROM INFORMATION_SCHEMA.TABLES
        WHERE TABLE_TYPE = 'BASE TABLE'
    ) AS t
    LEFT JOIN
    (
        SELECT TABLE_NAME, CONSTRAINT_NAME, CONSTRAINT_TYPE
        FROM INFORMATION_SCHEMA.TABLE_CONSTRAINTS
        WHERE CONSTRAINT_TYPE = 'PRIMARY KEY'
    ) AS tc
    ON t.TABLE_NAME = tc.TABLE_NAME
WHERE tc.TABLE_NAME IS NULL;
```

### 46. 実行プラン
SQLは取得したいデータを宣言型で定義するためのものであり、そのデータをもっとも効率的な方法で特定する方法はオプティマイザに委ねられる。

ゾウとネズミの問題。
対象の DB の行数が変われば、実行プランを再生成すべきである。より効率的な手順はデータの分散状況に依存している。ストアドプロシージャなど、パラメータ化されたクエリの実行プランを、データベースがキャッシュする場合に、特に問題になる！

- 使用されていないインデックスがないか、またその原因を調べる
- 実行プランが効率的かは、データの分散状況に依存する
- 優れた実行プランを生成できるようにインデックスを追加する


## sec 8
SQL で直積を生成するには CROSS JOIN を使用する。

### 47. 間接的に関連しているテーブルの行の特定
``` sql
-- 直積を使って全ての顧客とすべての製品のリストを取得
SELECT c.CustomerID, c.CustFirstName, c.CustLastName,
    p.ProductNumber, p.ProductName, p.ProductDescription
FROM Customers AS c, Products AS p;

-- 購入された製品をリストアップ
SELECT o.OrderNumber, o.CustomerID, od.ProductNumber
FROM Orders AS O
    INNER JOIN Order_Details AS od
        ON o.OrderNumber = od.OrderNumber;

-- 直積の行のうち、購入されたものとされてないものを特定する
SELECT CustProd.CustomerID, CustProd.CustFirstName, CustProd.CustLastName,
    CustProd.ProductNumber, CustProd.ProductName,
    (CASE WHEN OrdDet.OrderCount > 0
        THEN 'You purchased this!'
        ELSE ' '
        END
    ) AS ProductOrdered
FROM
(
    SELECT c.CustomerID, c.CustFirstName, c.CustLastName,
        p.ProductNumber, p.ProductName, p.ProductDescription
    FROM Customers AS c, Products AS p
) AS CustProd
LEFT JOIN
    (
        SELECT o.CustomerID, od.ProductNumber, COUNT(*) AS OrderCount
        FROM Orders AS o
            INNER JOIN Order_Details AS od
                ON o.OrderNumber = od.OrderNumber
        GROUP BY o.CustomerID, od.ProductNumber
    ) AS OrdDet
    ON CustProd.CustomerID = OrdDet.CustomerID
        AND CustProd.ProductNumber = OrdDet.ProductNumber
ORDER BY CustProd.CustomerID, CustProd.ProductName;
```

- 2つのテーブル内のレコードをあらゆる方法で組み合わせるには、直積を使用
- 実際に発生した組み合わせの特定には、INNER JOIN を使用

### 48. 行を等量分類でランク付け
上位20%, 40%, 60%, 80% など、バンド幅で結果を格付けする。

- ウィンドウ関数 RANK() を利用すれば、ランク付けされた値を簡単に生成できる

### 49. テーブルの行を他の全ての行と組み合わせる

### 50. カテゴリをリストアップし、第一希望、第二希望、、と照合する


## sec 9
タリーテーブル（tally table）は、通常は１つの列だけで構成されたテーブルであり、１（or 0）からその状況に応じた最大値までの連続する数字が含まれている。

直積はベーステーブルの実際の値に依存するのに対し、タリーテーブルはすべての可能性をカバーする。

### 51. 空のデータ行を生成
- 特にレポートにおいて、からのデータ行の生成が役に立つことがある
- 数値との比較に基づいて複数の行を人工的に生成するのに役立つ

### 53. 
- データベースで見つからない値を生成するには、タリーテーブルを使用する
    - 基本は、データベースである値しか作成できない

### 54. タリーテーブルの値の範囲に基づいて別のテーブルの値を変換
GROUP BY の問題点の１つは、データを集計するににはそれらの値が同じでなければならないこと！場合によっては、値の範囲を同じように扱いたいこともある。

``` sql
WITH StudentGrades (Student, Subject, FinalGrade) AS (
    SELECT stu.SutudentFirstNM AS Student,
        sub.SubjectNM AS Subject, ss.FinalGrade
    FROM StudentSubjects AS ss
        INNER JOIN Students AS stu
            ON ss.StudentID = stu.StudentID
        INNER JOIN Subjects AS sub
            ON ss.SubjectID = sub.SubjectID
)

SELECT ag.Subject, gr.LetterGrade, COUNT(*) AS NumberOfStudents
FROM StudentGrades AS sg
    INNER JOIN GradeRanges AS gr
        ON sg.FinalGrade >= gr.LowGradePoint
            AND sg.FinalGrade <= gr.HighGradePoint
GROUP BY sg.Subject, gr.LetterGrade
ORDER BY sg.Subject, gr.LetterGrade
```

### 55. 日付テーブルを使って日付の計算を単純化
日付と時刻は問題と隣り合わせのテーブル型である。

- 日付や、日付に基づく計算に大きく依存するアプリケーションでは、日付テーブルを作成するとロジックが大幅に単純になることがある
- 日付テーブルは拡張可能、休業日、会計年度など、アプリケーション固有のドメインを追加できる
- 日付テーブルはディメンジョンテーブルであるため、インデックスを必要なだけ作成できる（？）

### 56. 特定の期間内の日付が全て列挙された予定表の作成
``` sql
CREATE TABLE Appointments (
    AppointmentID int IDENTITY (1, 1) PRIMARY KEY,
    ApptStartDate date NOT NULL,
    ApptStartTime time NOT NULL,
    ApptEndDate date NOT NULL,
    ApptEndTime time NOT NULL,
    ApptDescription varchar(50) NULL
)
```

日付フィールドと時刻フィールドに別々に格納しておくと、sargable クエリを記述するのが容易になる。

- 日付テーブルには、適切なインデックスを作成しておく
- 使用している RDBMS で日付と時刻を適切に処理する方法を理解する

### 57. タリーテーブルを使ったデータのピボット選択
- データのピボット選択が必要な場合、データベースシステムによっては、カスタム構文がサポートされていることがある


## sec 10
SQL は階層型のデータモデルを扱うようなものではない（苦手分野）。階層型のデータモデルを SQL データベースで作成する必要が生じるたびに、トレードオフを強いられる。データの正規化をとるのか、それともメタデータの取得と管理の容易さをとるのか。

### 58. 隣接リストモデル（adjacency list model）
従業員のモデルに対し、その上司を把握したいとする。カラムに上司の EmployeeID を持たせる、つまり、外部キー制約を持つテーブルに列を作成し、そのテーブルの主キーを参照させればよい。

このように『同じテーブルの主キーを参照する外部キーを作成』すれば、１つのテーブルで無限の深さの階層を作成できる！
 
``` sql
-- 自己参照の外部キー！
CREATE TABLE Employees (
    EmployeeID int PRIMARY KEY,
    EmpName varchar(255) NOT NULL,
    EmpPosition varchar(255) NOT NULL,
    SupervisorID int NULL
);

ALTER TABLE Employees
    ADD FOREIGN KEY (SupervisorID)
        REFERENCES Employees (EmployeeID);


-- 3レベルの自己結合
SELECT e1.EmpName AS Employee, e2.EmpName AS Supervisor,
    e3.EmpName AS SupervisorsSupervisor
FROM Employees AS e1
    LEFT JOIN Employees AS e2
        ON e1.SupervisorID = e2.EmployeeID
    LEFT JOIN Employees AS e3
        ON e2.SupervisorID = e3.EmployeeID
```

- 隣接リストモデルでは、テーブルに列を追加し、そのテーブルの主キーを参照する外部キーを使用
    - メタデータは不要！

### 59. 更新が頻繁に発生しない場合は、入れ子集合モデルでクエリ高速化！
入れ子集合（nested set）。

- 小ノードを持たないノードでは、「左」の番号と「右」の番号の差は１である

``` sql
CREATE TABLE Employees (
    EmployeeID int PRIMARY KEY,
    EmpName varchar(255) NOT NULL,
    EmpPosition varchar(255) NOT NULL,
    SupervisorID int NULL,
    lft int NULL,
    rgt int NULL);

-- 指定されたノードの子孫を全て検索
SELECT e.*
FROM Employees AS e
WHERE e.left >= @lft AND e.rgt <=@rgt;

-- 特定のノードの祖先を全て検索
SELECT *
FROM Employees AS e
WHERE e.lft <= @lft AND e.rgt >= @rgt;
```

- 階層が頻繁に更新される場合には適していない
- 単一ルートの階層を１つだけ使用する場合に適している

### 60. 限定的な検索には経路実体化モデル（materialized path）
概念上は、ファイルシステムのパスを使用するのと何ら変わらない。フォルダとファイルの代わりに主キーを使用する。

``` sql
CREATE TABLE Employees (
    EmployeeID int PRIMARY KEY,
    EmpName varchar(255) NOT NULL,
    EmpPosition varchar(255) NOT NULL,
    SupervisorID int NULL,
    HierarchyPath varchar(255)
)
```

HierarchyPath には、たとえば 12 の ID のユーザーには `1/3/8/12` のように入れる。

``` sql
-- 特定のノードの子孫をすべて検索
SELECT e.*
FROM Employees AS e
WHERE e.HierarchyPath LIKE @Nodepath + '%';
-- 2 -> 1/2 の部下の従業員を全て見つけ出す方法
-- @Nodepath に 1/2/ を指定
```

- 経路実体化モデルでの検索は、事実上、一方公に限られている
    - 術後の先頭 or 途中にワイルドカードが含まれている場合「sargable クエリ」を作成するのが不可能だから

### 61. 複雑な検索にはクロージャモデル
先祖クロージャテーブル（ancestry closure table）を使用する。これは「経路実体化モデル」のリレーショナルアプローチであり、２つ目のテーブルを用意し、ノード間の「つながり」ごとにメタデータのレコードを作成する。

``` sql
CREATE TABLE Employees (
    EmployeeID int NOT NULL PRIMARY KEY,
    EmpName varchar(255) NOT NULL,
    EmpPosition varchar(255) NOT NULL,
    SupervisorID int NULL,
);

CREATE TABLE EmployeesAncestry (
    SupervisedEmployeeID int NOT NULL,
    SupervisingEmployeeID int NOT NULL,
    Distance int NOT NULL,
    PRIMARY KEY (SupervisedEmployeeID, SupervisingEmployeeID)
);

ALTER TABLE EmployeesAncestry
    ADD CONSTRAINT FK_EmployeesAncestry_SupervisingEmployeeID
        FOREIGN KEY (SupervisingEmployeeID)
            REFERENCES Employees (EmployeeID);

ALTER TABLE EmployeesAncestry
    ADD CONSTRAINT FK_EmployeesAncestry_SupervisedEmployeeID
        FOREIGN KEY (SupervisedEmployeeID)
            REFERENCES Employees (EmployeeID);
```

- クロージャテーブルの管理は複雑
- 頻繁な更新と容易な検索が必要な場合のみ、有効


## 日付型

### PostgreSQL
- DATE
- TIME [タイムゾーン]
- TIMESTAMP [タイムゾーン]
- INTERVAL
