<script>
  import GlobalInfo from '../GlobalInfo.vue'
  export default {
    /**
     * 设置Session
     * @param {Object} data
     */
    setSessionStorage(data) {
      window.sessionStorage.setItem(GlobalInfo.projectName, data)
    },
    /**
     * 获取Session
     */
    getSessionStorage() {
      return window.sessionStorage.getItem(GlobalInfo.projectName)
    },
    /**
     * 清理Session
     */
    clearSessionStorage() {
      window.sessionStorage.clear()
    },
    /**
     * 图片转换Base64
     * @param {Object} img
     * @param {Object} call
     */
    coverImgToBase64(img, call) {
      const reader = new FileReader()
      reader.addEventListener('load', () => call(reader.result))
      reader.readAsDataURL(img)
    },
    /**
     * 转换时间为01
     * @param {Object} m
     */
    coverTime(m) {
      return m < 10 ? `0${m}` : m
    },
    /**
     * 转换时间为年-月-日 时：分：秒
     * @param {Object} seconds
     */
    coverTimeToFormat(seconds) {
      let time = new Date(seconds)
      let y = time.getFullYear()
      let m = time.getMonth() + 1
      let d = time.getDate()
      let h = time.getHours()
      let mm = time.getMinutes()
      let s = time.getSeconds()
      return `${y}-${this.coverTime(m)}-${this.coverTime(d)} ${this.coverTime(h)}:${this.coverTime(mm)}`
    },
    /**
     * 获取年月日
     * @param {Object} seconds
     */
    getTimeData(seconds) {
      let time = new Date(seconds)
      return {
        year: time.getFullYear(),
        month: this.coverTime(time.getMonth() + 1),
        day: this.coverTime(time.getDate()),
      }
    },
    /**
     * 获取fontSize的文字大小
     * @param {Object} num
     */
    getFitSize(num) {
      let clientWidth = document.documentElement.scrollWidth
      return GlobalInfo.fitData.fontSize * (clientWidth / GlobalInfo.fitData.windowSize.width) * num
    },
    /**
     * 获取时间
     */
    getTime() {
      return Date.parse(new Date())
    },
    /**数组转换为对象 */
    coverArrToObject(arr, key) {
      let data = {}
      if (arr) {
        arr.forEach((res) => {
          data[res[key]] = res
        })
      }
      return data
    },
    /**刷新echarts*/
    updateEcharts(e, option) {
      setTimeout(() => {
        e.resize()
        e.setOption(option, true)
        window.addEventListener('resize', () => {
          e.resize()
        })
      }, 100)
    },
  }
</script>
