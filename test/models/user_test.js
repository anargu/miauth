const assert = require('assert')
const Sequelize = require('sequelize')
const path = require('path')
const dbConnUtils = require(path.join(__dirname, '../utils/db_conn.js'))

describe('Testing User model interactions', function () {
    const POSTGRE_DB_URL_CONNECTION = `postgres://postgres:miauth-test@localhost:5432/mocha-temporal-test-db`
    let sequelize
    let db

    before('Initializing db & models', async () => {
        await dbConnUtils.initializeDatabase(POSTGRE_DB_URL_CONNECTION, { logging: false })

        sequelize = new Sequelize(POSTGRE_DB_URL_CONNECTION);
        db =  require(path.join(__dirname, '../../models/index')).initSequelize(sequelize)
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
                username: _userCreated.username
            }, 'anotherpass')
            console.log('=== USER HASH PASSWORD AFTER ===> ', _userUpdated.hash)
            assert.notDeepEqual(_userUpdated.hash, _userCreated.hash)                    
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