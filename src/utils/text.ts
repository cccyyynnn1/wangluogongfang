export function hexStringToArrayBuffer(str: string) {
  let count = str.length / 2
  if (count.toString().indexOf('.') != -1) {
    count = parseInt(count.toString()) + 1 // 有.5就加 1
    count = parseInt(count.toString()) // 或者不加
  }
  const buffer = new ArrayBuffer(count)
  const dataView = new DataView(buffer)
  for (let i = 0; i < count; i++) {
    const curCharCode = parseInt(str.substr(i * 2, 2), 16)
    dataView.setUint8(i, curCharCode)
  }
  return buffer
}

export function stringToHex(data: string, delim = '', padding = 2, extraDelim = '', lineSize = 0) {
  if (!data) return ''

  let output = ''
  const prepend = delim === '0x' || delim === '\\x'

  for (let i = 0; i < data.length; i++) {
    const hex = data.charCodeAt(i).toString(16)
    const num = (i + 1) % 4 == 0
    const flag = Math.floor(i / 16) % 2 == 0
    let dom = ''
    if (flag && !num) {
      dom = `<i style="display: inline-block;width: 6.25%;font-size: 14px;line-height: 28px;text-align: center;font-style: normal;color:#606266;background:#F4F3FA;">${hex
        .padStart(2, '0')
        .toUpperCase()}</i>`
    } else if (flag && num) {
      dom = `<i style="display: inline-block;width: 6.25%;font-size: 14px;line-height: 28px;text-align: center;font-style: normal;color:#606266;background:#F4F3FA;border-right: 1px solid #eee;">${hex
        .padStart(2, '0')
        .toUpperCase()}</i>`
    } else if (!flag && !num) {
      dom = `<i style="display: inline-block;width: 6.25%;font-size: 14px;line-height: 28px;text-align: center;font-style: normal;color:#606266;background:#FBFBFE;">${hex
        .padStart(2, '0')
        .toUpperCase()}</i>`
    } else if (!flag && num) {
      dom = `<i style="display: inline-block;width: 6.25%;font-size: 14px;line-height: 28px;text-align: center;font-style: normal;color:#606266;background:#FBFBFE;border-right: 1px solid #eee;">${hex
        .padStart(2, '0')
        .toUpperCase()}</i>`
    }

    output += prepend ? delim + dom : dom + delim
    // output += prepend ? delim + hex : hex + delim

    if (extraDelim) {
      output += extraDelim
    }
    // Add LF after each lineSize amount of bytes but not at the end
    if (i !== data.length - 1 && (i + 1) % lineSize === 0) {
      // output += '\n'
      output += `<br />`
    }
  }

  // Remove the extraDelim at the end (if there is one)
  // and remove the delim at the end, but if it's prepended there's nothing to remove
  const rTruncLen = extraDelim.length + (prepend ? 0 : delim.length)
  if (rTruncLen) {
    // If rTruncLen === 0 then output.slice(0,0) will be returned, which is nothing
    return output.slice(0, -rTruncLen)
  } else {
    return output
  }
}

export function removeLinerBeaksAndReturns(str: string) {
  if (!str) return ''
  return str.replace(/^[\n\r]+|[\n\r]+$/g, '')
}
