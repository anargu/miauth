
function isEmpty(value) {
    if (value instanceof String) {
        const _value = value.trim()
        return _value === ''
    } else {
        return (value === undefined || value === null)
    }
}


module.exports = {
    isEmpty
}
