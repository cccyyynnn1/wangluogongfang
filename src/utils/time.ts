import dayjs, { Dayjs } from 'dayjs'
/**
 *  获取最近一段时间。
 *  rangeType：最近 几天｜小时｜分钟
 **/
export function getLatestTime(timeRange: number, rangeType: 'day' | 'hour' | 'minute') {
  return formatTime(dayjs().subtract(timeRange, rangeType))
}

/**
 * 获取当年当月当周当日的00:00时间
 **/
export function getStartofTime(type: 'year' | 'month' | 'week' | 'day' = 'day') {
  return formatTime(dayjs().startOf(type))
}
/**
 *   type:格式化样式
 *   time:时间戳 ｜ 2016-05-03 ｜ Dayjs
 **/
export function formatTime(time: number | string | Dayjs, type = 'YYYY-MM-DD HH:mm:ss') {
  return dayjs(Number(time)).format(type)
}
/**
 *   type:格式化样式
 *   time:时间戳 ｜ HH:mm:ss ｜ Dayjs
 **/
export function formatTimeToHSM(time: number | string | Dayjs, type = 'HH:mm:ss') {
  return dayjs(Number(time)).format(type)
}
/**
 * isNsTime:是否是纳秒
 **/
export function formatNstime(time: number, isNsTime = true, type = 'YYYY-MM-DD HH:mm:ss') {
  if (!time) return ''
  const num = isNsTime ? time / 1000000 : time
  return formatTime(num, type)
}

/**
 * @description 获取时间间隔
 * @param time
 * @param option
 * @returns {string}
 */

// JS 计算两个时间戳相差年月日时分秒
export function calculateDiffTime(time: number) {
  const endTime = Date.now()
  const startTime = endTime - time
  let runTime = endTime - startTime
  let year = Math.floor(runTime / 86400 / 365) == 0 ? '' : `${Math.floor(runTime / 86400 / 365)}年`
  runTime = runTime % (86400 * 365)
  let month = Math.floor(runTime / 86400 / 30) == 0 ? '' : `${Math.floor(runTime / 86400 / 30)}月`
  runTime = runTime % (86400 * 30)
  let day = Math.floor(runTime / 86400) == 0 ? '' : `${Math.floor(runTime / 86400)}天`
  runTime = runTime % 86400
  let hour = Math.floor(runTime / 3600) == 0 ? '' : `${Math.floor(runTime / 3600)}小时`
  runTime = runTime % 3600
  let minute = Math.floor(runTime / 60) == 0 ? '' : `${Math.floor(runTime / 60)}分`
  runTime = runTime % 60
  let second = `${runTime}秒`
  return year + month + day + hour + minute + second
}
