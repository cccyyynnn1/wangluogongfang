import { URLDecode, unescapeUnicodeCharacters, fromHTMLEntity } from '@/assets/cyberChef/node/index.mjs'

export function PacketsMagic(magicStr?: string) {
  if (!magicStr) return ''
  const magicSteps = [
    [URLDecode, URLDecode],
    [fromHTMLEntity, fromHTMLEntity],
    [unescapeUnicodeCharacters, unescapeUnicodeCharacters],
    [fromHTMLEntity, fromHTMLEntity],
  ]
  const newVal = magicSteps.reduce((valStr: string, magics: any[]) => {
    return magics.reduce((_valstr: string, magic: any) => {
      return magic(_valstr).value
    }, valStr)
  }, magicStr)
  return newVal
}
