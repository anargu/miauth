const Sequelize = require('sequelize')

class User extends Sequelize.Model {
    static init (sequelize) {
        return super.init({
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
            sequelize,
            modelName: 'Users'
        })
    }
}

module.exports = User
