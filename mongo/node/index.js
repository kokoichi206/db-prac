const express = require('express');
const { connectToDB, getDB } = require('./db');
const { ObjectId } = require('mongodb');

const app = express();

let db;

connectToDB((err) =>  {
  if (!err) {
    app.listen(3000, () => {
      console.log('Server is running on port 3000');
    });
    db = getDB();
  }
});

// http://localhost:3000/books?page=1
app.get('/books', (req, res) => {
  const page = req.query.page || 0;
  const pageSize = 3;

  let books = [];

  db.collection('books').find()
    .sort({ author: 1 })
    .skip(page * pageSize)
    .limit(pageSize)
    .forEach(book => books.push(book))
    .then(() => {
      res.status(200).json(books)
    })
    .catch((err) => {
      res.status(500).json({ error: err })
    });
});

app.get('/books/:id', (req, res) => {
  console.log('req.params.id: ' + req.params.id);
  if (ObjectId.isValid(req.params.id) === false) {
    res.status(400).json({ error: 'Invalid ID' });
    return;
  }

  db.collection('books')
    .findOne({_id: new ObjectId(req.params.id)})
    .then((doc) => {
      res.status(200).json(doc)
    })
    .catch((err) => {
      res.status(500).json({ error: err })
    });
});
