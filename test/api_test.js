const path = require('path')
const chaiHttp = require("chai-http")
const chai = require('chai')

chai.should()
chai.use(chaiHttp)

describe('Api test methods', () => {

    before('setting env variables', () => {
        process.env.PORT = '8000'
        process.env.MIAUTH_CONFIG_FILE = path.join(__dirname, './test.config.yml')
    })

    describe('reset password flow', () => {        
        it('should return reset password web page', (done) => {
            const server = require(path.join(__dirname, '../src/server'))
            chai.request(server)
            .get('/forgot/reset')
            .end((err, res) => {
                if(err) {
                    done(err)
                }
                res.should.have.header('Content-Type', /^text\/html/)
                res.should.have.status(200)
                done()
            })
        })
    })
})