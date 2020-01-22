const assert = require('assert')
const path = require('path')
const axios = require('axios')
const chaiHttp = require("chai-http")
const chai = require('chai')

const bodyParser = require('body-parser')
const express = require('express')

chai.should()
chai.use(chaiHttp)

const { checkSchema, validationResult } = require('express-validator')

function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

describe('User Api test methods', () => {

    let app
    before('setting env variables', () => {
        const miauthConfig = require('../utils/init_config').miauthSetup()

        const userSchemaValidation = (() => {
            const _userSchemaValidation = {}
            if (miauthConfig.user.username) {
                _userSchemaValidation['username'] = {
                    matches: {
                        options: [new RegExp(miauthConfig.field_validations.username.pattern, 'g')],
                    },
                    isLength: {
                        errorMessage: 'Invalid username. Username should be between ' +
                            miauthConfig.field_validations.username.len[0] + ' and ' +
                            miauthConfig.field_validations.username.len[1] +  ' characteres.',
                        options: {
                            min: miauthConfig.field_validations.username.len[0],
                            max: miauthConfig.field_validations.username.len[1],
                        }
                    },
                    notEmpty: true,
                    errorMessage: miauthConfig.field_validations.username.invalid_pattern_error_message,
                }
            }
            
            if (miauthConfig.user.email) {
                _userSchemaValidation['email'] = {
                    matches: {
                        options: [new RegExp(miauthConfig.field_validations.email.pattern, 'g')],
                    },
                    isLength: {
                        errorMessage: 'Invalid email. Email should be between ' +
                        miauthConfig.field_validations.email.len[0] + 'and' +
                        miauthConfig.field_validations.email.len[1] + ' characteres.',
                        options: {
                            min: miauthConfig.field_validations.email.len[0],
                            max: miauthConfig.field_validations.email.len[1],
                        }
                    },
                    errorMessage: miauthConfig.field_validations.email.invalid_pattern_error_message,
                    notEmpty: true,
                }
            }

            _userSchemaValidation['password'] = {
                notEmpty: true,
                isString: true,
                isLength: {
                    options: {
                        min: miauthConfig.field_validations.password.len[0],
                        max: miauthConfig.field_validations.password.len[1],
                    },
                    errorMessage: 'Invalid password. Password should be between ' +
                        miauthConfig.field_validations.password.len[0] + ' and ' +
                        miauthConfig.field_validations.password.len[1] + ' characteres.',
                },
                errorMessage: miauthConfig.field_validations.password.invalid_pattern_error_message,
            }

            return _userSchemaValidation
        })()

        app = express()
        app.use(express.json())
        app.use(express.urlencoded({ extended: true }))
        app.post('/login', checkSchema(userSchemaValidation), async (req,res) => {
            if (!validationResult(req).isEmpty()) {
                // console.log(validationResult(req).errors)
                res.status(400).send({
                    errors: validationResult(req).errors
                })
                return                
            }
            res.status(200).send({})
        })

        app.post('/signup', checkSchema(userSchemaValidation), (req,res) => {
            res.status(200).send({})
        })
    });

    describe('When User logins [Chai]', () => {
        const testcases = [
            { username: 'Juan1234', email: 'abc@hotmail.com', password: 'juan' },
            { username: 'Juanabc', email: 'absdsc@hotmail.com', password: 'juan1' },
            { username: 'JuanThree', email: 'abc@gmail.com', password: 'juan2' }
        ]
        testcases.forEach(async function (testData, i) {
            it(`Given valid user data Then should return ok response ${i}`, async function() {
                this.timeout(8000)
                try {
                    const res = await chai.request(app)
                    .post('/login')
                    .send({
                        'username': testData.username,
                        'email': testData.email,
                        'password': testData.password
                    })
                    res.should.have.status(200)                                    
                } catch (error) {
                    assert.fail(error)
                }
            })
        })
        // it(`Given valid user data Then should return ok response`, async () => {
            // try {
            //     testcases.forEach(async (testData, i) => {
            //         chai.request(app)
            //         .post('/login')
            //         .send({
            //             'username': testData.username,
            //             'email': testData.email,
            //             'password': testData.password
            //         })
            //         .end((err, res) => {
            //             if(err) throw err
            //             res.should.have.status(200)
            //         })                        
            //         console.log('undone')
            //     })
            //     console.log('done')
            //     done()
            // } catch (error) {
            //     done(err)
            // }

            // try {
            //     await testcases.forEach(async (testData) => {
            //         try {
            //             const res = await chai.request(app)
            //             .post('/login')
            //             .send({
            //                 'username': testData.username,
            //                 'email': testData.email,
            //                 'password': testData.password
            //             })
            //             res.should.have.status(200)                            
            //         } catch (error) {
            //             throw error                        
            //         }
            //     })                    
            // } catch (error) {
            //     throw error
            // }            
        // })
    })

    describe('When User logins [Axios]', () => {

        let srv
        before('setup server', () => {
            srv = app.listen('7890')
        })

        const testcases = [
            { username: 'Juan1234', email: 'abc@hotmail.com', password: 'juan' },
            { username: 'Juanabc', email: 'absdsc@hotmail.com', password: 'juan1' },
            { username: 'JuanThree', email: 'abc@gmail.com', password: 'juan2' }
        ]
        testcases.forEach(async (testData) => {
            it(`Given valid user data Then should return ok response ${testData.username}`, async function () {
                this.timeout(8000)
                try {
                    
                    let res = await axios.post('http://localhost:7890/login', {
                        'username': testData.username,
                        'email': testData.email,
                        'password': testData.password
                    })
                    assert((typeof res) === 'object')
                    assert(res.status === 200)
                    await sleep(2000)
                } catch (error) {
                    assert.fail(error)
                }
            })
        })

        after('killing srv', () => {
            srv.close()
        })
    })

})
