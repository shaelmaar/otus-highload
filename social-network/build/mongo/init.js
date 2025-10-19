const appUsername = process.env.MONGO_DB_USER || 'app_user';
const appPassword = process.env.MONGO_DB_PASSWORD || 'app_password';
const appDatabase = process.env.MONGO_DB_NAME || 'dialogs_db';

db = db.getSiblingDB(appDatabase);

db.createUser({
    user: appUsername,
    pwd: appPassword,
    roles: [
        { role: 'readWrite', db: appDatabase },
        { role: 'dbAdmin', db: appDatabase }
    ]
});
