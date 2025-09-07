import CryptoJS from 'crypto-js'

export default function AesEncryptCBC(data: string) {
  //将key,data,iv都转化为wordarry格式，根据传入的编码格式选择对应的方法
  const key = CryptoJS.enc.Utf8.parse('oCrVqwkA9MbrL8i126GiuA==')
  const iv = CryptoJS.enc.Utf8.parse('ABCDEFGHIJKLM_iv')
  const _data = CryptoJS.enc.Utf8.parse(data)

  const encrypted = CryptoJS.AES.encrypt(_data, key, {
    iv: iv,
    mode: CryptoJS.mode.CBC,
    padding: CryptoJS.pad.Pkcs7,
  })
  //这里返回的是Base64的密文
  return CryptoJS.enc.Base64.stringify(encrypted.ciphertext)
}
