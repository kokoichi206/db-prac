## INSERT 文生成
xargs で、別プロセスを立ち上げないとランダムな部分に毎回同じ値が使われてしまう問題が発生。
原因わからず。

``` sh
# これだと何故か name の値が毎回同じのが出力されてしまう？
echo {1..100} | xargs -n 1 -I@ echo "INSERT INTO users (id, name) VALUES (@, "\'"$(cat /dev/urandom | LC_CTYPE=utf_8 tr -dc a-z | fold -w 5 | head -n 1)"\'");"

# 上記問題は解決できた
echo {1..100} | xargs -n 1 -I@ sh -c 'echo "INSERT INTO users (id, name) VALUES (@, '\''$(cat /dev/urandom | LC_CTYPE=utf_8 tr -dc a-z | fold -w 5 | head -n 1)'\'');"'

# for だと問題なし
for i in {0..100}; do echo "INSERT INTO users (id, name) VALUES ($i, "\'"$(cat /dev/urandom | LC_CTYPE=utf_8 tr -dc a-z | fold -w 5 | head -n 1)"\'");"; done
```

### `xargs shuf`

``` sh
$ echo {1..10} | xargs shuf -n1 -e
5
$ echo {1..10} | xargs shuf -e
7
10
9
3
2
4
5
6
1
8
```

