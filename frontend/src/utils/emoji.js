export function decodeEmojiSequence(seq) {
  return String.fromCodePoint(
    ...seq.split('-').map(hex => parseInt(hex, 16))
  )
}
export function emojiFromCodepoints(cps) {
  return String.fromCodePoint(...cps)
}

function parseCodepoint(str, base) {
  const cp = parseInt(str, base)
  if (!Number.isFinite(cp) || cp < 0 || cp > 0x10FFFF) {
    return null
  }

  return cp
}

export function codePointsToHex(cps) {
  return cps.map(cp => cp.toString(16).toUpperCase()).join('-')
}

export function detectEmojiInputType(input) {
  if (!input || typeof input !== 'string') {
    return { type: 'unknown' }
  }
  const trimmed = input.trim()
  if (!trimmed) {
    return { type: 'unknown' }
  }

  // 1.多码点序列 十进制: 128558-8205-128168
  if (/^\d+(-\d+)+$/.test(trimmed)) {
    const parts = trimmed.split('-')
    const codePoints = []
    for (const part of parts) {
      const cp = parseCodepoint(part, 10)
      if (!cp) return { type: 'unknown' }
      codePoints.push(cp)
    }

    return {
      type: 'decimal_codepoint_sequence',
      value: codePoints
    }
  }

  // 2.多码点序列: 如：1F62E-200D-1F4A8
  if (/^[0-9A-Fa-f]{4,6}(-[0-9A-Fa-f]{4,6})+$/.test(trimmed)) {
    const parts = trimmed.split('-')
    const codePoints = []
    for (const part of parts) {
      const cp = parseCodepoint(part, 16)
      if (!cp) return { type: 'unknown' }
      codePoints.push(cp)
    }

    return {
      type: 'hex_codepoint_sequence',
      value: codePoints
    }
  }


  // 3.多码点序列 如：U+1F62E-U+200D-U+1F4A8
  const uPlusMulti = trimmed.match(
    /^U\+[0-9A-Fa-f]{4,6}(-U\+[0-9A-Fa-f]{4,6})+$/i
  )
  if (uPlusMulti) {
    const parts = trimmed.split('-')
    const codePoints = []
    for (const part of parts) {
      const hex = part.slice(2)
      const cp = parseCodepoint(hex, 16)
      if (!cp) return { type: 'unknown' }
      codePoints.push(cp)
    }

    return {
      type: 'u_plus_codepoint_sequence',
      value: codePoints
    }
  }

  // 4. 十进制格式
  if (/^\d+$/.test(trimmed)) {
    const cp = parseInt(trimmed, 10)
    if (!cp) return { type: 'unknown' }
    return { type: 'decimal_codepoint', value: [cp] }
  }

  // 5.纯16进制码点字符串 (1F600)
  if (/^[0-9A-Fa-f]{4,6}$/.test(trimmed)) {
    const cp = parseCodepoint(trimmed, 16)
    if (!cp) return { type: 'unknown' }
    return { type: 'hex_codepoint', value: [cp] }
  }

  // 6.U+格式 (U+1F600)
  const match = trimmed.match(/^U\+([0-9A-Fa-f]{4,6})$/i)
  if (match) {
    const cp = parseCodepoint(match[1], 16)
    if (!cp) return { type: 'unknown' }
    return { type: 'u_plus_codepoint', value: [cp] }
  }

  // 7.已是字符串
  if (trimmed.length > 0) {
    const cp = trimmed.codePointAt(0)
    return { type: 'character', value: [cp], char: trimmed }
  }

  return { type: 'unknown' }
}

export function isEmojiUnicode(cp) {
  return (
    (cp >= 0x1F600 && cp <= 0x1F64F) || // Emoticons
    (cp >= 0x1F300 && cp <= 0x1F5FF) || // Misc Symbols & Pictographs
    (cp >= 0x1F680 && cp <= 0x1F6FF) || // Transport & Map
    (cp >= 0x1F900 && cp <= 0x1F9FF) // Supplemental Symbols
  )
}
