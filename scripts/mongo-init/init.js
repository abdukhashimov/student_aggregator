db.users.createIndex({"email": 1}, {unique: true});
db.users.createIndex({"refresh_token.token": 1});
db.users.createIndex({"refresh_token.expires_at": 1});

db.schemas.createIndex({"slug": 1}, {unique: true});
