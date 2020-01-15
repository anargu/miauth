const Sequelize = require('sequelize')

const sequelize = new Sequelize('postgres://user:pass@example.com:5432/dbname');

const User = sequelize.import('./user.js')
const Session = sequelize.import('./session.js')

const db = {
    User,
    Session
}

User.associate(db)
Session.associate(db)

module.exports = db
