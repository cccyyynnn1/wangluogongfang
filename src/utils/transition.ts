import { useCopy, uuid } from '@/utils'

type getTableCopyDataType = {
  tableColumn: any
  row: any
  column: any
  mothod: any
}
/**
 * @description 站点，检索，告警，网络应用层右键方法列表
 * @param {row: any, column: any}
 * @returns void
 */
export function getTableCopyData({ tableColumn, row, column, mothod }: getTableCopyDataType) {
  const NOEXIST = 'not_exists'
  const EXIST = 'exists'
  const target = column.property
  let value: any = undefined
  let label = ''
  let arr = [{ label: '复制', callback: useCopy, value: value }]
  let targetList: any = undefined
  if (target) {
    value = row[target]
    arr = [{ label: '复制', callback: useCopy, value: value }]
    targetList = tableColumn.filter((item: any) => {
      return item.fieldNameEn == target
    })
    label = targetList[0]!.fieldNameCn
    const arr1 = [
      {
        label: value || value === 0 || value === '0' ? `+只查看 ${label} = ${value}` : `+只查看 ${label} ${NOEXIST}`,
        callback: mothod,
        value: {
          key: target,
          relation: value || value === 0 || value === '0' ? '=' : NOEXIST,
          value: value,
          label: label,
          enable: true,
          id: uuid(),
        },
      },
    ]
    const arr2 = [
      {
        label: value || value === 0 || value === '0' ? `+只查看 ${label} != ${value}` : `+只查看 ${label} ${EXIST}`,
        callback: mothod,
        value: {
          key: target,
          relation: value || value === 0 || value === '0' ? '!=' : EXIST,
          value: value,
          label: label,
          enable: true,
          id: uuid(),
        },
      },
    ]
    if (label.search('时间') != -1) {
      targetList[0].fieldType = 'no_serach'
    }
    switch (targetList[0].fieldType) {
      case 'ip':
        arr = [...arr, ...arr1, ...arr2]
        break
      case 'num':
        arr = [...arr, ...arr1, ...arr2]
        break
      case 'isn_t':
        arr = !value || value === 0 || value === '0' ? [...arr] : [...arr, ...arr1, ...arr2]
        // arr = [...arr]
        break
      case 'direction':
        arr = [...arr, ...arr1]
        break
      case 'text':
        arr = [...arr, ...arr1, ...arr2]
        break

      default:
        arr = [{ label: '复制', callback: useCopy, value: value }]
        break
    }
  }
  return arr
}

/**
 * @description 对会话站点告警的详情的请求体进行格式化
 * @param { str: string }
 * @returns string
 * 10.99.19.88 30858
 */

export function formatStrToJson(str: string) {
  if (!str || str?.length === 0) return str
  let str1 = JSON.parse(JSON.stringify(str))
  if (str1[str1.length - 1] == '\n') {
    str1 = str1.substring(0, str.length - 1)
  }
  if (str1[0] === '{' && str1[str1.length - 1] === '}') {
    const reg = RegExp(/\n/)
    if (str1.match(reg)) {
      return str
    } else {
      const str3 = `{\r\n    ${str1
        .substring(1, str1.length - 1)
        .split(',')
        .join(' , \r\n    ')},\r\n }`
      return str3
    }
  } else {
    return str1
  }
}
