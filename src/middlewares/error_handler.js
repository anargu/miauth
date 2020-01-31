const multer = require('multer')

const handleError = (err, req, res, next) => {
    if (err instanceof multer.MulterError) {
        // A Multer error occurred when uploading.
        // err.
        const { code, message, field } = err
        res.status(400).json({
            code,
            error: 'MulterError',
            error_description: message,
            user_message: `${field}: ${message}`
        })
        return    
    } else if(err.isAxiosError) {
        res.status(500).json({
            code: 0,
            error: 'DOSMJError',
            error_description: err.response.data.error,
            user_message: 'Mail service not working properly. Please try later.'
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