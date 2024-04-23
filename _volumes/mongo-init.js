db = db.getSiblingDB('admin')
db.auth("admin","admin")
db = db.getSiblingDB('medods_test')
db.createUser(
    {
        user: "user",
        pwd: "user",
        roles: [
            {
                role: "readWrite",
                db: "medods_test"
            }
        ]
    }
);
db.createCollection('tokens')
db.createCollection('users')
