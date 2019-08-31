const Sequelize = require('sequelize')
const initializeModel = require('./models').initializeModel

let sequelize = null

function initDatabase () {
    sequelize = new Sequelize('postgres://user:pass@example.com:5432/dbname',
        {
            pool: {
                max: 5,
                min: 0,
                acquire: 30000,
                idle: 10000
            }
        })

    // init model
    initializeModel(sequelize)
}

module.exports = {
    initDatabase,
    sequelize
}
