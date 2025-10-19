db.createCollection("dialogMessages");

db.dialogMessages.createIndex({ "dialogID": "hashed" });