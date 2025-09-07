export type AttackCharacterizationType = {
  content: string
  id?: number
  remark: string
  scope: 'requestHeader' | 'requestPayload' | 'responseHeader' | 'responsePayload' | 'all'
}
