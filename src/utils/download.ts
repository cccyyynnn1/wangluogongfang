import dayjs from 'dayjs'
import { DownloadPcapApi } from '../api-ecs/toolbox'
import { useUserStore } from '@/store/modules/user'
import { usePubilcStore } from '@/store/modules/public'
import { gp } from '@gp'
const publicStore = usePubilcStore()

export async function downloadLogPacket(
  infoVal: {
    clientIp: string
    clientPort: string
    serverIp: string
    serverPort: string
    startTimeNs?: number
    requestTimeNs: number
    id?: number
    probeIds?: number[]
  },
  dataType: number
) {
  const { clientIp, clientPort, serverIp, serverPort, startTimeNs, requestTimeNs, probeIds } = infoVal
  const val = startTimeNs ? startTimeNs / 1000000 : requestTimeNs / 1000000
  const { getFlowProbesIDList } = publicStore
  const arr = getFlowProbesIDList()
  // 根据probeIds作为key，从arr中获取对应的value，放入flowProbeIds中
  const flowProbeIds = probeIds?.map((item) => {
    const flowProbes = arr.find((item2) => item2.realId === item)
    return flowProbes?.id
  })
  const querySql = {
    flowProbeIds,
    objectList: [{ clientIp, serverIp, serverPort: serverPort.toString(), probeId: flowProbeIds?.toString() }],
    savePackage: false,
    searchModel: 'all',
    searchTable: 'eventStat',
    timeRange: `${dayjs(val).subtract(1, 'minutes').format('YYYY-MM-DD HH:mm:ss')} - ${dayjs(val)
      .add(1, 'minutes')
      .format('YYYY-MM-DD HH:mm:ss')}`,
  }
  const { getUserId } = useUserStore()
  try {
    const { msg } = await DownloadPcapApi({ uid: getUserId(), query: querySql }, dataType)
    gp.$baseMessage(msg, 'success', 'vab-hey-message-success', false)
  } catch (error: any) {
    console.log(error)
  }
}

/**
 * @description 下载二进文件流
 * @param param:{file: BlobPart,word: string}
 * @returns void
 */
export function downloadFile(file: BlobPart, word: string) {
  const downloadElement = document.createElement('a')
  const blob = new Blob([file], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' })
  const href = window.URL.createObjectURL(blob)
  downloadElement.href = href
  downloadElement.download = `${word}.xlsx`
  document.body.appendChild(downloadElement)
  downloadElement.click()
  document.body.removeChild(downloadElement)
  window.URL.revokeObjectURL(href)
}
/**
 * @description 下载二进文件流
 * @param param:{file: BlobPart,word: string}
 * @returns void
 */
export function downloadTextFile(file: BlobPart, word: string) {
  const downloadElement = document.createElement('a')
  const blob = new Blob([file], { type: 'text/plain' })
  const href = window.URL.createObjectURL(blob)
  downloadElement.href = href
  downloadElement.download = `${word}.text`
  document.body.appendChild(downloadElement)
  downloadElement.click()
  document.body.removeChild(downloadElement)
  window.URL.revokeObjectURL(href)
}
/**
 * @description 下载原始日志
 * @param param:{text: string,name: string}
 * @returns void
 */
export function downloadSourceLog(text: string, name: string) {
  const formatSourceData = (val: string): object | string => {
    if (!val) return ''
    let _val = ''
    try {
      _val = JSON.parse(val)
    } catch (error) {
      _val = val
    }
    return _val
  }
  const data = formatSourceData(text)
  const jsonStr = data instanceof Object ? JSON.stringify(data, null, '\t') : data
  const downloadElement = document.createElement('a')
  const blob = new Blob([jsonStr], { type: 'text/plain' })
  const href = window.URL.createObjectURL(blob)
  downloadElement.href = href
  downloadElement.download = name
  document.body.appendChild(downloadElement)
  downloadElement.click()
  document.body.removeChild(downloadElement)
  window.URL.revokeObjectURL(href)
}

// 下载文件，自定义文件名称
export function downFile(url: string, fileName: string) {
  // const x = new XMLHttpRequest()
  // x.open('GET', url, true)
  // x.responseType = 'blob'
  // x.onload = function () {
  // const url = window.URL.createObjectURL(x.response)
  const a = document.createElement('a')
  a.href = url
  a.download = fileName
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  // }
  // x.send()
}

export async function downloadLogPacketFlowSearch(query: any) {
  downFile(`/download/netPcap/${query}`, '数据' + String(Date.now()) + '.pcap')
}

/**
 * 下载 Excel 文件
 * @param fileDate 接口返回的文件数据
 * @param fileName 下载的文件名
 */
export const downloadExcel = (fileDate: any, fileName: string) => {
  return new Promise((resolve, reject) => {
    // 为 blob 设置文件类型，这里以 .xls 为例
    const blob = new Blob([fileDate], {
      type: 'application/vnd.ms-excel',
    })
    // 创建一个临时的url指向blob对象
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = fileName
    a.click()

    // 释放这个临时的对象url
    window.URL.revokeObjectURL(url)

    return resolve('导出成功')
  })
}
