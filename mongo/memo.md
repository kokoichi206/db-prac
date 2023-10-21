## [Complete MongoDB Tutorial](https://www.youtube.com/watch?v=ExcRbA7fy_A&list=PL4cUxeGkcC9h77dJ-QJlwGlZlTd4ecZOA&ab_channel=NetNinja)

- SQL で操作する database ではない
- Documents として扱う
  - json っぽいデータ構造を扱える

### seup

- https://www.mongodb.com/ja-jp

### [mongo -> mongosh](https://www.mongodb.com/docs/manual/release-notes/6.0-compatibility/#legacy-mongo-shell-removed)

[mongosh(MongoDB Shell)](https://www.mongodb.com/docs/mongodb-shell/#mongodb-binary-bin.mongosh)

``` sh
docker-compose exec mongo mongosh
```

### mongo-express

root, example では入れないと思ったら admin:pass で入れた。

```
mongo-mongo-express-1  | Mongo Express server listening at http://0.0.0.0:8081
mongo-mongo-express-1  | Server is open to allow connections from anyone (0.0.0.0)
mongo-mongo-express-1  | basicAuth credentials are "admin:pass", it is recommended you change this in your config.js!
```

## Basics

- collections and documents
  - database の中に複数の collection を持つ
- nested documents
- nested array documents
- [monbodb compass](https://www.mongodb.com/try/download/atlascli)

[connection string](https://www.mongodb.com/docs/manual/reference/connection-string/?utm_source=compass&utm_medium=product)

```
mongodb://root:example@localhost:27017
```

- collection に複数のデータをまとめて登録も可能
  - bulk insert みたいなものか

### mongosh

``` mongo
test> show dbs
admin      100.00 KiB
bookstore   72.00 KiB
config     108.00 KiB
local       72.00 KiB

test> use bookstore

bookstore> show collections
books
bookstore> var name = "yoshi"

bookstore> name
yoshi
bookstore> help
```

CRUD

``` sh
bookstore> db.books.insertOne({title: "Ther color of magic", author: "piyo", pages: 222, rating: 1, genres: ["magic"]})
{
  acknowledged: true,
  insertedId: ObjectId("65335bda6fa321078d875085")
}


bookstore> db.books.insertMany([{title: "Ther color of magiccccccc", author: "piyo", pages: 245, rating: 3, genres: ["magic"]}])
{
  acknowledged: true,
  insertedIds: { '0': ObjectId("65335ca16fa321078d875086") }
}

bookstore> db.books.find()
bookstore> db.books.find({author: "piyo"})

# filter ?
bookstore> db.books.find({author: "piyo"}, {title: 1, author: 1})
[
  {
    _id: ObjectId("65335bda6fa321078d875085"),
    title: 'Ther color of magic',
    author: 'piyo'
  },
  {
    _id: ObjectId("65335ca16fa321078d875086"),
    title: 'Ther color of magiccccccc',
    author: 'piyo'
  }
]

bookstore> db.books.find({}, {rating: 1})
[
  { _id: ObjectId("6532c03cb868d8310fe4f93f"), rating: 9 },
  { _id: ObjectId("6532c0dcb868d8310fe4f943"), rating: 3 },
  { _id: ObjectId("6532c0dcb868d8310fe4f944"), rating: 5 },
  { _id: ObjectId("65335bda6fa321078d875085"), rating: 1 },
  { _id: ObjectId("65335ca16fa321078d875086"), rating: 3 }
]

bookstore> db.books.find().count()
5
bookstore> db.books.find({rating: 1}).count()
1
bookstore> db.books.find().limit(3)
bookstore> db.books.find().sort({pages: -1})

{
  "name": "foo",
  "age": 16,
  "labels": {
    "group": "saku",
    "ts": 124400
  }
}

bookstore> db.books.find({rating: {$gt: 5}})
bookstore> db.books.find({rating: {$gte: 5}})

bookstore> db.books.find({$or: [{rating: {$gte: 5}}, {pages: {$gt: 33}}]})

bookstore> db.books.find({rating: {$in: [4,5,6,7]}})
# not in
bookstore> db.books.find({rating: {$nin: [5,6,7,8]}})

# find value in array
bookstore> db.books.find({genres: "magic"})
# all value matches
bookstore> db.books.find({genres: ["magical"]})


# query on nested documents
# https://www.mongodb.com/docs/manual/tutorial/query-embedded-documents/
bookstore> db.users.find()
[
  {
    _id: ObjectId("6533610faecbf79ed7660ce8"),
    name: 'piyo',
    age: 17,
    labels: { group: 'saku', ts: 123900 }
  },
  {
    _id: ObjectId("6533612caecbf79ed7660cea"),
    name: 'foo',
    age: 16,
    labels: { group: 'saku', ts: 124400 }
  },
  {
    _id: ObjectId("6533613baecbf79ed7660cec"),
    name: 'bar',
    age: 21,
    labels: { group: 'no', ts: 134400 }
  }
]
bookstore> db.users.find({'labels.group': 'no'})
[
  {
    _id: ObjectId("6533613baecbf79ed7660cec"),
    name: 'bar',
    age: 21,
    labels: { group: 'no', ts: 134400 }
  }
]



bookstore> db.books.find({}, {_id: 1})
[
  { _id: ObjectId("6532c03cb868d8310fe4f93f") },
  { _id: ObjectId("6532c0dcb868d8310fe4f943") },
  { _id: ObjectId("6532c0dcb868d8310fe4f944") },
  { _id: ObjectId("65335bda6fa321078d875085") },
  { _id: ObjectId("65335ca16fa321078d875086") }
]
bookstore> db.books.deleteOne({ _id: ObjectId("6532c03cb868d8310fe4f93f") })
{ acknowledged: true, deletedCount: 1 }


bookstore> db.books.find({ _id: ObjectId("6532c0dcb868d8310fe4f943") })
[
  {
    _id: ObjectId("6532c0dcb868d8310fe4f943"),
    title: 'Name of the winter',
    author: 'john doe',
    pages: 111,
    genres: [ 'dystopian', 'magical' ],
    rating: 3
  }
]
bookstore> db.books.updateOne({ _id: ObjectId("6532c0dcb868d8310fe4f943") }, {$set: {rating: 2}})
{
  acknowledged: true,
  insertedId: null,
  matchedCount: 1,
  modifiedCount: 1,
  upsertedCount: 0
}
bookstore> db.books.find({ _id: ObjectId("6532c0dcb868d8310fe4f943") })
[
  {
    _id: ObjectId("6532c0dcb868d8310fe4f943"),
    title: 'Name of the winter',
    author: 'john doe',
    pages: 111,
    genres: [ 'dystopian', 'magical' ],
    rating: 2
  }
]

bookstore> db.books.find({ _id: ObjectId("6532c0dcb868d8310fe4f943") })
[
  {
    _id: ObjectId("6532c0dcb868d8310fe4f943"),
    title: 'Name of the winter',
    author: 'john doe',
    pages: 111,
    genres: [ 'dystopian', 'magical' ],
    rating: 2
  }
]
bookstore> db.books.updateOne({ _id: ObjectId("6532c0dcb868d8310fe4f943") }, {$inc: {pages: 2}})
{
  acknowledged: true,
  insertedId: null,
  matchedCount: 1,
  modifiedCount: 1,
  upsertedCount: 0
}
bookstore> db.books.find({ _id: ObjectId("6532c0dcb868d8310fe4f943") })
[
  {
    _id: ObjectId("6532c0dcb868d8310fe4f943"),
    title: 'Name of the winter',
    author: 'john doe',
    pages: 113,
    genres: [ 'dystopian', 'magical' ],
    rating: 2
  }
]
bookstore> db.books.updateOne({ _id: ObjectId("6532c0dcb868d8310fe4f943") }, {$push: {genres: "test"}})
{
  acknowledged: true,
  insertedId: null,
  matchedCount: 1,
  modifiedCount: 1,
  upsertedCount: 0
}
bookstore> db.books.find({ _id: ObjectId("6532c0dcb868d8310fe4f943") })
[
  {
    _id: ObjectId("6532c0dcb868d8310fe4f943"),
    title: 'Name of the winter',
    author: 'john doe',
    pages: 113,
    genres: [ 'dystopian', 'magical', 'test' ],
    rating: 2
  }
]
bookstore> db.books.updateOne({ _id: ObjectId("6532c0dcb868d8310fe4f943") }, {$push: {genres: {$each: ['1', '2']}}})
{
  acknowledged: true,
  insertedId: null,
  matchedCount: 1,
  modifiedCount: 1,
  upsertedCount: 0
}
bookstore> db.books.find({ _id: ObjectId("6532c0dcb868d8310fe4f943") })
[
  {
    _id: ObjectId("6532c0dcb868d8310fe4f943"),
    title: 'Name of the winter',
    author: 'john doe',
    pages: 113,
    genres: [ 'dystopian', 'magical', 'test', '1', '2' ],
    rating: 2
  }
]
```


## index

``` sh
bookstore> db.books.find({rating: 5}).explain('executionStats')
{
  explainVersion: '2',
  queryPlanner: {
    namespace: 'bookstore.books',
    indexFilterSet: false,
    parsedQuery: { rating: { '$eq': 5 } },
    queryHash: 'F5BC19EA',
    planCacheKey: '5F64DC11',
    maxIndexedOrSolutionsReached: false,
    maxIndexedAndSolutionsReached: false,
    maxScansToExplodeReached: false,
    winningPlan: {
      queryPlan: {
        stage: 'COLLSCAN',
        planNodeId: 1,
        filter: { rating: { '$eq': 5 } },
        direction: 'forward'
      },
      slotBasedPlan: {
        slots: '$$RESULT=s5 env: { s2 = Nothing (SEARCH_META), s3 = 1697877111768 (NOW), s1 = TimeZoneDatabase(Europe/Minsk...Asia/Saigon) (timeZoneDB), s7 = 5 }',
        stages: '[1] filter {traverseF(s4, lambda(l1.0) { ((l1.0 == s7) ?: false) }, false)} \n' +
          '[1] scan s5 s6 none none none none lowPriority [s4 = rating] @"20bc0997-d33a-4b41-b9a5-d486271fd1c5" true false '
      }
    },
    rejectedPlans: []
  },
  executionStats: {
    executionSuccess: true,
    nReturned: 1,
    executionTimeMillis: 1,
    totalKeysExamined: 0,
    totalDocsExamined: 4,
    executionStages: {
      stage: 'filter',
      planNodeId: 1,
      nReturned: 1,
      executionTimeMillisEstimate: 0,
      opens: 1,
      closes: 1,
      saveState: 0,
      restoreState: 0,
      isEOF: 1,
      numTested: 4,
      filter: 'traverseF(s4, lambda(l1.0) { ((l1.0 == s7) ?: false) }, false) ',
      inputStage: {
        stage: 'scan',
        planNodeId: 1,
        nReturned: 4,
        executionTimeMillisEstimate: 0,
        opens: 1,
        closes: 1,
        saveState: 0,
        restoreState: 0,
        isEOF: 1,
        numReads: 4,
        recordSlot: 5,
        recordIdSlot: 6,
        fields: [ 'rating' ],
        outputSlots: [ Long("4") ]
      }
    }
  },
  command: { find: 'books', filter: { rating: 5 }, '$db': 'bookstore' },
  serverInfo: {
    host: 'dc3f1cdffe9c',
    port: 27017,
    version: '7.0.2',
    gitVersion: '02b3c655e1302209ef046da6ba3ef6749dd0b62a'
  },
  serverParameters: {
    internalQueryFacetBufferSizeBytes: 104857600,
    internalQueryFacetMaxOutputDocSizeBytes: 104857600,
    internalLookupStageIntermediateDocumentMaxSizeBytes: 104857600,
    internalDocumentSourceGroupMaxMemoryBytes: 104857600,
    internalQueryMaxBlockingSortMemoryUsageBytes: 104857600,
    internalQueryProhibitBlockingMergeOnMongoS: 0,
    internalQueryMaxAddToSetBytes: 104857600,
    internalDocumentSourceSetWindowFieldsMaxMemoryBytes: 104857600,
    internalQueryFrameworkControl: 'trySbeEngine'
  },
  ok: 1
}


bookstore> db.books.getIndexes()
[
  { v: 2, key: { _id: 1 }, name: '_id_' },
  { v: 2, key: { rating: 8 }, name: 'rating_8' }
]


# totalDocsExamined が いい感じ！
bookstore> db.books.find({rating: 8}).explain('executionStats')
{
  explainVersion: '2',
  queryPlanner: {
    namespace: 'bookstore.books',
    indexFilterSet: false,
    parsedQuery: { rating: { '$eq': 8 } },
    queryHash: 'F5BC19EA',
    planCacheKey: '4B76EEEB',
    maxIndexedOrSolutionsReached: false,
    maxIndexedAndSolutionsReached: false,
    maxScansToExplodeReached: false,
    winningPlan: {
      queryPlan: {
        stage: 'FETCH',
        planNodeId: 2,
        inputStage: {
          stage: 'IXSCAN',
          planNodeId: 1,
          keyPattern: { rating: 8 },
          indexName: 'rating_8',
          isMultiKey: false,
          multiKeyPaths: { rating: [] },
          isUnique: false,
          isSparse: false,
          isPartial: false,
          indexVersion: 2,
          direction: 'forward',
          indexBounds: { rating: [ '[8, 8]' ] }
        }
      },
      slotBasedPlan: {
        slots: '$$RESULT=s11 env: { s10 = {"rating" : 8}, s2 = Nothing (SEARCH_META), s1 = TimeZoneDatabase(Europe/Minsk...Asia/Saigon) (timeZoneDB), s5 = KS(2B100104), s6 = KS(2B10FE04), s3 = 1697877300157 (NOW) }',
        stages: '[2] nlj inner [] [s4, s7, s8, s9, s10] \n' +
          '    left \n' +
          '        [1] cfilter {(exists(s5) && exists(s6))} \n' +
          '        [1] ixseek s5 s6 s9 s4 s7 s8 [] @"20bc0997-d33a-4b41-b9a5-d486271fd1c5" @"rating_8" true \n' +
          '    right \n' +
          '        [2] limit 1 \n' +
          '        [2] seek s4 s11 s12 s7 s8 s9 s10 [] @"20bc0997-d33a-4b41-b9a5-d486271fd1c5" true false \n'
      }
    },
    rejectedPlans: []
  },
  executionStats: {
    executionSuccess: true,
    nReturned: 1,
    executionTimeMillis: 1,
    totalKeysExamined: 1,
    totalDocsExamined: 1,
    executionStages: {
      stage: 'nlj',
      planNodeId: 2,
      nReturned: 1,
      executionTimeMillisEstimate: 0,
      opens: 1,
      closes: 1,
      saveState: 0,
      restoreState: 0,
      isEOF: 1,
      totalDocsExamined: 1,
      totalKeysExamined: 1,
      collectionScans: 0,
      collectionSeeks: 1,
      indexScans: 0,
      indexSeeks: 1,
      indexesUsed: [ 'rating_8' ],
      innerOpens: 1,
      innerCloses: 1,
      outerProjects: [],
      outerCorrelated: [ Long("4"), Long("7"), Long("8"), Long("9"), Long("10") ],
      outerStage: {
        stage: 'cfilter',
        planNodeId: 1,
        nReturned: 1,
        executionTimeMillisEstimate: 0,
        opens: 1,
        closes: 1,
        saveState: 0,
        restoreState: 0,
        isEOF: 1,
        numTested: 1,
        filter: '(exists(s5) && exists(s6)) ',
        inputStage: {
          stage: 'ixseek',
          planNodeId: 1,
          nReturned: 1,
          executionTimeMillisEstimate: 0,
          opens: 1,
          closes: 1,
          saveState: 0,
          restoreState: 0,
          isEOF: 1,
          indexName: 'rating_8',
          keysExamined: 1,
          seeks: 1,
          numReads: 2,
          indexKeySlot: 9,
          recordIdSlot: 4,
          snapshotIdSlot: 7,
          indexIdentSlot: 8,
          outputSlots: [],
          indexKeysToInclude: '00000000000000000000000000000000',
          seekKeyLow: 's5 ',
          seekKeyHigh: 's6 '
        }
      },
      innerStage: {
        stage: 'limit',
        planNodeId: 2,
        nReturned: 1,
        executionTimeMillisEstimate: 0,
        opens: 1,
        closes: 1,
        saveState: 0,
        restoreState: 0,
        isEOF: 1,
        limit: 1,
        inputStage: {
          stage: 'seek',
          planNodeId: 2,
          nReturned: 1,
          executionTimeMillisEstimate: 0,
          opens: 1,
          closes: 1,
          saveState: 0,
          restoreState: 0,
          isEOF: 0,
          numReads: 1,
          recordSlot: 11,
          recordIdSlot: 12,
          seekKeySlot: 4,
          snapshotIdSlot: 7,
          indexIdentSlot: 8,
          indexKeySlot: 9,
          indexKeyPatternSlot: 10,
          fields: [],
          outputSlots: []
        }
      }
    }
  },
  command: { find: 'books', filter: { rating: 8 }, '$db': 'bookstore' },
  serverInfo: {
    host: 'dc3f1cdffe9c',
    port: 27017,
    version: '7.0.2',
    gitVersion: '02b3c655e1302209ef046da6ba3ef6749dd0b62a'
  },
  serverParameters: {
    internalQueryFacetBufferSizeBytes: 104857600,
    internalQueryFacetMaxOutputDocSizeBytes: 104857600,
    internalLookupStageIntermediateDocumentMaxSizeBytes: 104857600,
    internalDocumentSourceGroupMaxMemoryBytes: 104857600,
    internalQueryMaxBlockingSortMemoryUsageBytes: 104857600,
    internalQueryProhibitBlockingMergeOnMongoS: 0,
    internalQueryMaxAddToSetBytes: 104857600,
    internalDocumentSourceSetWindowFieldsMaxMemoryBytes: 104857600,
    internalQueryFrameworkControl: 'trySbeEngine'
  },
  ok: 1
}


bookstore> db.books.dropIndex({rating: 8})
{ nIndexesWas: 2, ok: 1 }
```
