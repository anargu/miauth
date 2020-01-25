const express = require('express')
const { check_userId } = require('../middlewares/validations')
const { MiauthError } = require('../utils/error')
const { validationResult } = require('express-validator')

const path = require('path')
const multer = require('multer')

const One_MB = 1024
const storage = multer.diskStorage({
    destination: path.join(__dirname, '../../public/'),
    filename: (req, file, cb) => {
        cb(null, `${file.fieldname}.html`)
    }
})

const publicUpload = multer({
    storage: storage,
    limits: {
        fieldSize: 5 * One_MB,
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
})
const emailTemplatesUploaded = publicUpload.fields([
    { 
        name: 'reset_password',
        filename: 'reset_password.html',
        maxCount: 1
    },
    {
        name: 'reset_password_result_success',
        filename: 'reset_password_result_success.html',
        maxCount: 1,
    },
    {
        name: 'reset_password_result_error',
        filename: 'reset_password_result_error.html',
        maxCount: 1
    }
])

module.exports = (db) => {
    const { Session } = db

    /**
     * adminApi is intended to use only in private network with the 1st party application which miauth is attached
     */
    const adminApi = express.Router()

    adminApi.put('/update/templates', emailTemplatesUploaded, (req, res) => {
        res.status(200).json({ files_uploaded: req.files })
    })

    adminApi.post('/revoke_all', [
        check_userId()
    ], async (req, res, next) => {
        try {
            const errors = validationResult(req)
            if(!errors.isEmpty()) {
                throw new MiauthError(
                    400, 'ValidationError',
                    errors.array().map((err) => err.msg).join('. '))
            }
            const sessionsDeleted = await Session.revokeAllSessions({ userId: req.body.userId })

            res.status(200).json({ sessions_deleted: sessionsDeleted })    
        } catch (error) {
            next(error)
        }
    })

    return adminApi
}