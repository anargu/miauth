const Sequelize = require('sequelize')
const initializeModel = require('./models').initializeModel

const POSTGRES_USER = process.env.POSTGRES_USER
const POSTGRES_PASSWORD = process.env.POSTGRES_PASSWORD
const POSTGRES_HOST = process.env.POSTGRES_HOST
const POSTGRES_PORT = process.env.POSTGRES_PORT
const POSTGRES_DB = process.env.POSTGRES_DB

let sequelize = null

async function initDatabase () {
    sequelize = new Sequelize(
        `postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}`,
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
    await sequelize.sync()
}

module.exports = {
    initDatabase,
    sequelize
}
