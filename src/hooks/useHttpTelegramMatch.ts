import { Ref } from 'vue'
type HttpLog = {
  http_request_body?: string
  http_response_body?: string
  request_header?: string
  response_header?: string
  response_body_match?: number[]
  request_body_match?: number[]
  request_header_match?: number[]
  response_header_match?: number[]
  response_body_match_length?: number
  request_body_match_length?: number
  request_header_match_length?: number
  response_header_match_length?: number
  [key: string]: any
}
type MatchKeys = 'http_response_body' | 'http_request_body' | 'request_header' | 'response_header'
const PayloadDivider = `\n    \n`
const useHttpTelegramMatch = (
  httpLog: Ref<HttpLog>,
  callback?: (
    http: HttpLog,
    matchMap: {
      [key in MatchKeys]: string[]
    }
  ) => void
): { [key in MatchKeys]: string[] } => {
  const match_rules = ['response_body_match', 'request_body_match', 'request_header_match', 'response_header_match']
  const match_length = [
    'response_body_match_length',
    'request_body_match_length',
    'request_header_match_length',
    'response_header_match_length',
  ]
  const match_keys = ['http_response_body', 'http_request_body', 'request_header', 'response_header']

  const matchMap: {
    [key in MatchKeys]: string[]
  } = reactive({
    http_response_body: [],
    http_request_body: [],
    request_header: [],
    response_header: [],
  })

  const handlerMatch = (http: HttpLog) => {
    if (!http) return
    try {
      match_rules.forEach((rule, index) => {
        const matchs = http[rule] as number[]
        const matchKey = match_keys[index] as MatchKeys
        if (matchs) {
          const matchStr = http[matchKey] as string
          const matchLength = http[match_length[index]]
          matchs.forEach((match, index) => {
            const matchEndIndex = matchLength?.[index]
            const str = matchStr.slice(match, match + matchEndIndex)
            matchMap[matchKey].push(str)
          })
        } else {
          matchMap[matchKey] = []
        }
      })
      callback?.(http, matchMap)
    } catch (error) {
      console.error(error)
    }
  }
  watch(
    () => httpLog.value,
    (newVal) => {
      Object.assign(matchMap, {
        http_response_body: [],
        http_request_body: [],
        request_header: [],
        response_header: [],
      })
      handlerMatch(newVal)
    },
    {
      deep: true,
      immediate: true,
    }
  )
  return matchMap
}
export default useHttpTelegramMatch
export { PayloadDivider }
