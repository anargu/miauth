const Sequelize = require('sequelize')

const { POSTGRESQL_CONN_STRING } = process.env

// initialized()
module.exports = (() => {
    const _sequelize = new Sequelize(POSTGRESQL_CONN_STRING)
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
