
const assert = require('assert')
const Sequelize = require('sequelize')
const path = require('path')
const dbConnUtils = require(path.join(__dirname, '../utils/db_conn.js'))
const configTest = require(path.join(__dirname, '../utils/init_config.js'))


describe('Testing Session model interactions', function () {
    const POSTGRE_DB_URL_CONNECTION = `postgres://postgres:miauth-test@localhost:5432/mocha-temporal-test-db`
    let sequelize
    let db

    before('Initializing db & models', async () => {
        configTest.miauthSetup()
        await dbConnUtils.initializeDatabase(POSTGRE_DB_URL_CONNECTION, { logging: false })

        sequelize = new Sequelize(POSTGRE_DB_URL_CONNECTION);
        db =  require(path.join(__dirname, '../../src/models/index')).initSequelize(sequelize)
    })

    it('verify if \'createSession\' is identified as a function type', () => {
        // if yes, it's extending custom functions defined in Session model script. Great.
        const fn = db.Session.createSession
        assert.equal(typeof fn, 'function')
    })

    it('should create a session (w/o refresh_token)', async () => {
        try {
            const _user = await db.User.createUser({
                username: 'octavio',
                email: 'octacion@mail.com',
                password: 'PassOfOctavio'
            })

            const _sessionCreated = await db.Session.createSession({
                userId: _user.uuid,
                email: _user.email
            })
            console.log('=== SESSION CREATED ===>\n')
            console.log(_sessionCreated)
            assert.notDeepEqual(_sessionCreated, null)       
        } catch (error) {
            assert.fail(error)
        }
    })

    it('should create a session (with refresh_token)', async () => {
        try {
            const _user = await db.User.createUser({
                username: 'octavio2',
                email: 'octacion2@mail.com',
                password: 'PassOfOctavio2'
            })

            const _sessionCreated = await db.Session.createSession({
                userId: _user.uuid,
                email: _user.email
            })
            console.log('=== SESSION CREATED ===>\n')
            console.log(_sessionCreated)
            assert.notDeepEqual(_sessionCreated, null)
            assert.notDeepEqual(_sessionCreated.refresh_token, null)
        } catch (error) {
            assert.fail(error)
        }
    })

    after('deleting recreated database', async () => {
        if(sequelize)
            sequelize.close()

        await dbConnUtils.removingDatabase(POSTGRE_DB_URL_CONNECTION, { logging: true })
    })
})