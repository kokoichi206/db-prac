const { MongoClient } = require('mongodb');

let dbConnection;

module.exports = {
  connectToDB: (cb) => {
    MongoClient.connect('mongodb://root:example@localhost:27017/bookstore?authSource=admin')
    // MongoClient.connect('mongodb://localhost:27017/bookstore')
      .then((client) => {
        dbConnection = client.db();
        return cb();
      })
      .catch((err) => {
        console.error(err);
        return cb(err);
      });
  },
  getDB: () => dbConnection,
}
