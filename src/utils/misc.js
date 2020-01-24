
function isEmpty(value) {
    if (value instanceof String) {
        const _value = value.trim()
        return _value === ''
    } else {
        return (value === undefined || value === null)
    }
}

function onelongline(strings, ...values) {
  // $FlowFixMe: Flow doesn't undestand .raw
  const raw = typeof strings === "string" ? [strings] : strings.raw;

    // first, perform interpolation
    let result = "";
    for (let i = 0; i < raw.length; i++) {
        result += raw[i]
        // join lines when there is a suppressed newline
        .replace(/\\\n[ \t]*/g, "")
        // handle escaped backticks
        .replace(/\\`/g, "`");

        if (i < values.length) {
            result += values[i];
        }
    }

    const chunks = result.split('\n')
    let new_result = ''
    for (let i = 0; i < chunks.length; i++) {
        new_result += chunks[i].replace(/^(\t|\s){1,}|(\t|\s){2,}$/g, '')
    }
    return new_result
  }


module.exports = {
    isEmpty,
    onelongline
}
