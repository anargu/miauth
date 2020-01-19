const Sequelize = require('sequelize')

const {
    POSTGRESQL_CONN_STRING,
} = process.env

let sequelize = new Sequelize(POSTGRESQL_CONN_STRING)
let db

function initializeModels (sequelizeInstance) {
    const User = sequelizeInstance.import('./user.js')
    const Session = sequelizeInstance.import('./session.js')

    db = {
        User,
        Session
    }

    User.associate(db)
    Session.associate(db)
}
initializeModels(sequelize)

module.exports = db

module.exports.initSequelize = (sequelizeInstance) => {
    if (sequelize) {
        sequelize.close()
    }
    sequelize = sequelizeInstance
    initializeModels(sequelizeInstance)
    return db
}
