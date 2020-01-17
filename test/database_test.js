const assert = require('assert')
const Sequelize = require('sequelize')


describe('Connectivity to db', function () {
    const POSTGRE_DB_URL_CONNECTION = `postgres://postgres:miauth-test@localhost:5432/mocha-temporal-test-db`

    before('recreating a database', async () => {
        await require('./utils/db_conn.js').initializeDatabase(POSTGRE_DB_URL_CONNECTION, { logging: true })
    })

    let sequelize
    beforeEach('initialize sequelize connection to db', () => {
        sequelize = new Sequelize(POSTGRE_DB_URL_CONNECTION);
    })

    it('should connect to test db', async function () {
        try {
            await sequelize.authenticate()
            assert.ok(true)
        } catch (error) {
            assert.fail(error)
        }
    })

    afterEach('closing sequelize connection to db', () => {
        if(sequelize)
            sequelize.close()
    })
    
    after('deleting recreated database', async () => {
        await require('./utils/db_conn.js').removingDatabase(POSTGRE_DB_URL_CONNECTION, { logging: true })
    })
})