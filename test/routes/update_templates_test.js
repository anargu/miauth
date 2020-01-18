
const assert = require('assert')
const path = require('path')
const chaiHttp = require("chai-http")
const chai = require('chai')
const express = require('express')
const multer = require('multer')
const fs = require('fs')
const { handleError, MiauthError } = require(path.join(__dirname, '../../src/utils/error'))

chai.should()
chai.use(chaiHttp)

const One_MB = 1024
// const storage = multer.diskStorage({
//     destination: path.join(__dirname, '../../public/'),
//     filename: (req, file, cb) => {
//         cb(null, `XXXX_test_${file.fieldname}.html`)
//     }
// })
const storage = multer.memoryStorage()
const publicUpload = multer({
    storage: storage,
    limits: {
        fileSize: 1 * One_MB,
    },

    fileFilter: (req, file, cb) => {
        const fileType = /html/;
        const validMimetype = fileType.test(file.mimetype)
        const validExtName = fileType.test(path.extname(file.originalname).toLowerCase())
        if(validMimetype && validExtName) {
            return cb(null, true)
        }
        return cb(new MiauthError(400, 'InvalidFileExtension', 'Only .html files are allowed'))
    },
    
    onError: (err, next) => {
        next(err)
    }
})
const emailTemplatesUploaded = publicUpload.fields([
    { 
        name: 'reset_password',
        filename: 'reset_password.html',
        mimetype: 'text/html',
        maxCount: 1
    },
    {
        name: 'reset_password_result_success',
        filename: 'reset_password_result_success.html',
        mimetype: 'text/html',
        maxCount: 1,
    },
    {
        name: 'reset_password_result_error',
        filename: 'reset_password_result_error.html',
        mimetype: 'text/html',
        maxCount: 1
    },
    {
        name: 'email_reset_instructions',
        filename: 'email_reset_instructions.html',
        mimetype: 'text/html',
        maxCount: 1
    }
])

describe('Updates of templates used to reset password', () => {
    let srv
    before(() => {
        process.env.PORT = '8000'
        process.env.MIAUTH_CONFIG_FILE = path.join(__dirname, '../test.config.yml')

        srv = express()
        srv.put('/admin/update/templates', emailTemplatesUploaded, (req, res) => {
            
            res.status(200).json({ files_uploaded: req.files })
        })
        
        srv.use(function(err, req, res, next) {
            if (err)
                handleError(err, res)
        });  
        
        assert.notDeepEqual(srv, null || undefined)
    })

    it('Api Server should receive all four files sent to the endpoint /update/templates', (done) => {
        const newTemplate = `<html><body><h1>Hello</h1></body></html>`
        
        chai.request(srv)
        .put('/admin/update/templates')
        .attach('reset_password',
            Buffer.from(newTemplate),
            'reset_password.html')
        .attach('reset_password_result_success',
            Buffer.from(newTemplate),
            'reset_password_result_success.html')
        .attach('reset_password_result_error',
            Buffer.from(newTemplate),
            'reset_password_result_error.html')
        .attach('email_reset_instructions',
            Buffer.from(newTemplate),
            'email_reset_instructions.html')
        .end((err, res) => {
            if (err) {
                done(err)
                return
            }

            res.should.have.status(200)
            res.should.to.be.json
            assert(typeof res.body.files_uploaded, 'object')
            assert.notDeepEqual(res.body.files_uploaded, {})
            done()
        })
    })

    it(`should deny heavy files > ${1 * One_MB} bytes (1 MB)`, (done) => {
        chai.request(srv)
        .put('/admin/update/templates')
        .attach('reset_password',
            fs.readFileSync(path.join(__dirname, '../assets/image_one_mb.jpg')),
            'reset_password.html')
        .end((err, res) => {
            if (err) {
                done(err)
                return
            }
            res.should.have.status(400)
            res.should.to.be.json
            done()
        })
    })

    it('should deny files with other extension than .html', (done) => {
        chai.request(srv)
        .put('/admin/update/templates')
        .attach('reset_password',
            fs.readFileSync(path.join(__dirname, '../assets/no_html_extension.txt')),
            'no_html_extension.txt')
        .end((err, res) => {
            if (err) {
                done(err)
                return
            }
            res.should.have.status(400)
            res.should.to.be.json
            assert.notDeepEqual(res.body.error, undefined)
            assert.notDeepEqual(res.body.error, null)
            done()
        })
    })
})