
function errorMessage (userMessage = '', message = '', code = 100) {
    return {
        user_message: userMessage,
        message,
        code
    }
}

function isEmpty(value) {
    if (value instanceof String) {
        const _value = value.trim()
        return _value === ''
    } else {
        return (value === undefined || value === null)
    }
}


module.exports = {
    errorMessage,
    isEmpty
}
