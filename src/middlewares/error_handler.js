const multer = require('multer')

const handleError = (err, req, res, next) => {
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
    const { statusCode, error, message, user_message, name } = err
    res.status(statusCode || 500).json({
        error: error || name,
        error_description: message,
        user_message
    })
};


module.exports = {
    handleError
}