const multer = require('multer')

class MiauthError extends Error {
    constructor(statusCode = 500, error, message = '', user_message) {
        super();

        this.statusCode = statusCode
        this.error = error
        this.message = message
        this.user_message = user_message || message;
    }
}

module.exports = {
    MiauthError
}