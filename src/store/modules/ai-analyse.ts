import { AiAnalyseModuleType } from '/#/store'
import { useUserStore } from './user'
const { userId } = useUserStore()
const AiAnalyseHistoryKey = 'aiAnalyseHistory'
import { isJson } from '@/utils/validate'
const getAnalyseHistoryLocalStorage = (key: string) => {
  const value: string | null = localStorage.getItem(key)
  return value && isJson(value) ? new Map(Object.entries(JSON.parse(value))) : new Map()
}
const AnalyseHistoryMaxSize = 1000
export const useAiAnalyseHistoryStore = defineStore('aiAnalyseHistory', {
  state: (): AiAnalyseModuleType => ({
    aiAnalyseHistory: getAnalyseHistoryLocalStorage(AiAnalyseHistoryKey),
  }),
  getters: {
    /**
     *
     * @description 判断历史记录是否超出最大长度
     * @returns {boolean} true 没超长可以添加  false 超长
     */
    analyseHistorySizeIsSafe: (state) => {
      const total = [...state.aiAnalyseHistory.values()].reduce((size, history) => {
        return size + history.history.length
      }, 0)
      return total > AnalyseHistoryMaxSize ? false : true
    },
    /**
     * @description 获取最早历史
     */
    getEarliestHistory: (state) => {
      let minTimeHistoryKey = ''
      let minTiem = new Date().getTime()
      state.aiAnalyseHistory.forEach((value, key) => {
        if (value.history.length === 0) {
          state.aiAnalyseHistory.delete(key)
        } else {
          if (value.updateTime < minTiem) {
            minTiem = value.updateTime
            minTimeHistoryKey = key
          }
        }
      })
      return state.aiAnalyseHistory.get(minTimeHistoryKey)
    },
    /**
     *
     * @description 获取当前历史记录
     */
    getCurrentAlarmHistory: (state) => {
      return (alarmId: string) => state.aiAnalyseHistory.get(`${userId}#${alarmId}`)?.history
    },
  },
  actions: {
    /**
     *
     * @description 添加历史记录
     * @param  {string} alarmId 告警日志ID
     */
    setAnalyseHistory(alarmId: string, history: { assistant: string; user: string; isAnalyse?: boolean }) {
      const localKey = `${userId}#${alarmId}`
      const currAlarm = this.aiAnalyseHistory.get(localKey)
      const updateTime = new Date().getTime()
      if (!this.analyseHistorySizeIsSafe) {
        this.getEarliestHistory?.history.pop()
      }
      if (currAlarm) {
        currAlarm.history.push(history)
        currAlarm.updateTime = updateTime
      } else {
        this.aiAnalyseHistory.set(localKey, {
          history: [history],
          updateTime: updateTime,
        })
      }
    },
  },
})
useAiAnalyseHistoryStore().$subscribe((mutation, state) => {
  const localeData = [...state.aiAnalyseHistory.entries()].reduce(
    (obj, [key, value]) => ((obj[key] = value), obj),
    {} as any
  )
  localStorage.setItem(AiAnalyseHistoryKey, JSON.stringify(localeData))
})
