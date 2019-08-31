const Sequelize = require('sequelize')

const POSTGRES_USER = 'miauth'
const POSTGRES_PASSWORD = 'm14uth'
const POSTGRES_HOST = 'localhost'
const POSTGRES_PORT = '5432'
const POSTGRES_DB = 'miauth-db'

describe ('application', function () {

    it('shoud connect to db', async function () {
        let sequelize = new Sequelize(
            `postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}`,
            {
                pool: {
                    max: 5,
                    min: 0,
                    acquire: 30000,
                    idle: 10000
                }
            })
        try {
            const result = await sequelize.authenticate()            
            console.log('Connection has been established successfully.')
        } catch (error) {
            console.error('Unable to connect to the database:', error)
            throw error
        }
    })
})
