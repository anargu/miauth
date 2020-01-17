const Sequelize = require('sequelize')

const {
    POSTGRES_USER,
    POSTGRES_PASSWORD,
    POSTGRES_HOST,
    POSTGRES_PORT,
    POSTGRES_DB
} = process.env

let sequelize = new Sequelize(`postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}`)
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
