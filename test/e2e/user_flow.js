
const assert = require('assert')
const path = require('path')
const dbConnUtils = require(path.join(__dirname, '../utils/db_conn.js'))

const chaiHttp = require("chai-http")
const chai = require('chai')

chai.should()
chai.use(chaiHttp)

const configTest = require(path.join(__dirname, '../utils/init_config.js'))

describe('Testing User Flows', () => {
    describe('When an user does a signup flow with correct parameters', () => {
        let db
        let srv
        before('Setting environment', async () => {
            process.env.POSTGRESQL_CONN_STRING = `postgres://postgres:miauth-test@localhost:5432/mocha-temporal-test-db`
            configTest.miauthSetup()
            db = await dbConnUtils.initializeDatabase(process.env.POSTGRESQL_CONN_STRING, { logging: false })
            try {
                await db._sequelize.authenticate()
                assert.ok(true)
            } catch (error) {
                assert.fail(error)
            }
            srv = require(path.join(__dirname, '../../app')).server
        })

        it('Should return a successful response and user data saved on DB', (done) => {
            const userData = {
                username: 'anUsername',
                email: 'anEmail@abc.com',
                password: 'aPassword'
            }
            chai.request(srv)
            .post('/auth/signup')
            .send({...userData})
            .end((err, res) => {
                if (err)
                    done(err)

                res.should.have.status(200)
                res.should.to.be.json
                assert.equal(typeof res.body.uuid, 'string')
                assert.deepEqual(res.body.username, userData.username)
                assert.deepEqual(res.body.email, userData.email)
                console.log('user_created', res.body)
                done()
            })
        })

        after('Cleaning messy and used environment', async function () {
            this.timeout(5000)
            try {
                await dbConnUtils.truncateTables(db._sequelize)
            } catch (error) {
                throw error                
            }
        })
    })
})