const assert = require('assert')
const Sequelize = require('sequelize')
const path = require('path')
const dbConnUtils = require(path.join(__dirname, '../utils/db_conn.js'))
const configTest = require(path.join(__dirname, '../utils/init_config.js'))

describe('Testing User model interactions', function () {
    const POSTGRESQL_CONN_STRING = `postgres://postgres:miauth-test@localhost:5432/mocha-temporal-test-db`
    let db
    before('Initializing db & models', async () => {
        configTest.miauthSetup()
        process.env.POSTGRESQL_CONN_STRING = POSTGRESQL_CONN_STRING
        db = await dbConnUtils.initializeDatabase(POSTGRESQL_CONN_STRING, { logging: false })
    })

    it('verify if \'createUser\' is identified as a function type', () => {
        // if yes, it's extending custom functions defined in user model script. Great.
        const fn = db.User.createUser
        assert.equal(typeof fn, 'function')
    })

    let _userCreated
    it('should create an User', async function() {
        try {
            const _user = await db.User.createUser({
                username: 'octavio',
                email: 'octacion@mail.com',
                password: 'PassOfOctavio'
            })
            console.log('=== USER CREATED ===>\n')
            console.log(_user)
            assert.notDeepEqual(_user, null)
            _userCreated = _user
        } catch (error) {
            assert.fail(error)
        }
    })

    it('should find a previously created User', async function () {
        if (!_userCreated)
            assert.fail(new Error('previous test failed, so this test neither will work'))
        
        try {
            let _userFound = await db.User.findByUsername(_userCreated.username)
            console.log('=== USER FOUND ===>\n')
            console.log(_userFound)
            assert.notDeepEqual(_userFound, null)
        } catch (error) {
            assert.fail(error)                
        }
    })

    it('should update password of user', async () => {
        if (!_userCreated)
            assert.fail(new Error('previous test failed, so this test neither will work'))

        try {
            console.log('=== USER HASH PASSWORD PREVIOUS ===> ', _userCreated.hash)
            let _userUpdated = await db.User.updatePassword({
                field: 'username',
                value: _userCreated.username
            }, 'anotherpass')
            console.log('=== USER HASH PASSWORD AFTER ===> ', _userUpdated.hash)
            assert.notDeepEqual(_userUpdated.hash, _userCreated.hash)                    
        } catch (error) {
            assert.fail(error)            
        }
    })

    after('cleaning tables', async () => {
        await dbConnUtils.truncateTables(db._sequelize)
    })
})