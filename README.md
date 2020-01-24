# puzzle
今のところパズドラ解析プログラムを作っています。

## 使い方
`$ git clone https://github.com/TatsukiFukushima/puzzle`  
`$ cd puzzle`

①ターミナルから盤面を入力  
`$ go run puzzle.go`  
"> " に続いて盤面を１行ずつ入力。

②txtに盤面を保存して入力（おすすめ）  
同ディレクトリ内に input.txt みたいなテキストファイルを作成。  
`$ go run puzzle.go < input.txt`

## 盤面フォーマット
```
2 1 3 4 6 3
1 3 6 1 5 4
3 4 5 2 2 4
4 3 5 6 4 5
2 6 4 5 4 5

```
空白を入れてください。
