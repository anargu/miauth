const assert = require('assert')
const chaiHttp = require("chai-http")
const chai = require('chai')

const express = require('express')

chai.should()
chai.use(chaiHttp)

const { validationResult, oneOf } = require('express-validator')

describe('User Api test methods', () => {

    let app
    before('setting env variables', () => {
        const miauthConfig = require('../utils/init_config').miauthSetup()
        const { check_username, check_email, check_password } = require('../../src/middlewares/validations')

        app = express()
        app.use(express.json())
        app.use(express.urlencoded({ extended: true }))
        app.post('/login', [
            check_password(),
            oneOf([
                ... (miauthConfig.user.username) ? [check_username()] : [],
                ... (miauthConfig.user.email) ? [check_email()] : [],
            ])
        ], async (req,res) => {
            if (!validationResult(req).isEmpty()) {
                // console.log(validationResult(req).errors)
                res.status(400).send({
                    errors: validationResult(req).errors
                })
                return                
            }
            res.status(200).send({})
        })

        app.post('/signup', [
            ... (miauthConfig.user.username) ? [check_username()] : [],
            ... (miauthConfig.user.check_email) ? [check_email()] : [],
            check_password(),
        ], (req,res) => {
            res.status(200).send({})
        })
    });

    describe('When User logins ', () => {
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
    })

})
