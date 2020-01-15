const Sequelize = require('sequelize')
const DataTypes = require('sequelize').DataTypes

module.exports = (sequelize) => {
    return sequelize.define('session', {
        uuid: {
          primaryKey: true,
          type: Sequelize.UUID,
          defaultValue: DataTypes.UUIDV4,
          allowNull: false
        },
        username: {
            type: Sequelize.STRING
        },
        email: {
            type: Sequelize.STRING,
            validate: {
              isEmail: true,
            }
        },
        hash: {
            type: Sequelize.STRING
        }
    }, {
    })
}