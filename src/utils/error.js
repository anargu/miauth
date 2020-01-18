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

const handleError = (err, res) => {
    if (err instanceof multer.MulterError) {
        // A Multer error occurred when uploading.
        // err.
        const { code, message, field } = err
        res.status(400).json({
            code,
            error_description: message,
            user_message: `${field}: ${message}`
        })
        return    
    }
    const { statusCode, error, message, user_message } = err;
    res.status(statusCode || 500).json({
        error,
        error_description: message,
        user_message
    })
};

module.exports = {
    MiauthError,
    handleError
}