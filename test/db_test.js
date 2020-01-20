const assert = require('assert')
const path = require('path')
const dbConnUtils = require(path.join(__dirname, './utils/db_conn.js'))


describe('Connectivity to db', function () {
    const POSTGRE_DB_URL_CONNECTION = `postgres://postgres:miauth-test@localhost:5432/mocha-temporal-test-db`

    let db
    before('recreating a database', async () => {
        require('./utils/init_config.js').miauthSetup()
        process.env.POSTGRESQL_CONN_STRING = POSTGRE_DB_URL_CONNECTION        
        db = await dbConnUtils.initializeDatabase(POSTGRE_DB_URL_CONNECTION, { logging: false })
    })

    it('should connect to test db', async function () {
        try {
            await db._sequelize.authenticate()
            assert.ok(true)
        } catch (error) {
            assert.fail(error)
        }
    })

    after('deleting recreated database', async () => {
        await dbConnUtils.truncateTables(db._sequelize)
    })
})