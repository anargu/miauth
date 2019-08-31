const Sequelize = require('sequelize')

class User extends Sequelize.Model {}

function initializeModel (sequelize) {
    User.init({
        id: {
            type: Sequelize.UUID,
            primaryKey: true,
            defaultValue: Sequelize.UUIDV4
        },
        email: {
            type: Sequelize.STRING,
            unique: true
        },
        hash: Sequelize.STRING
    }, {
        sequelize
    })
}

module.exports = {
    initializeModel,
    User
}
