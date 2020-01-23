const Sequelize = require('sequelize')

const { 
    POSTGRES_USER,
    POSTGRES_PASSWORD,
    POSTGRES_DB,
    POSTGRES_HOST,
    POSTGRES_PORT } = process.env

// initialized()
module.exports = (() => {
    const _sequelize = new Sequelize(`postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}`)
    const User = _sequelize.import('./user.js')
    const Session = _sequelize.import('./session.js')

    const db = {
        User,
        Session,
        _sequelize,
        close: async () => {
            await _sequelize.close()
        }
    }

    User.associate(db)
    Session.associate(db)

    return db
})() 
