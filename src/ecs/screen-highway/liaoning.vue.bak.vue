<template>
  <div class="indexBox">
    <header class="headBox">
      <span>{{ tempData.title }}</span>
      <div class="time">{{ time }}</div>
    </header>
    <div class="leftBox">
      <div class="secureBox">
        <div class="title">安全态势</div>
        <div v-if="tempData.activeDefenseCount > 0" class="info flexBox">
          <img class="icon" src="./assets/img/icon_22.png" />
          <div class="txt flexBox">
            主动防御次数
            <span>
              {{ tempData.activeDefenseCountFormat }}
              <e>次</e>
            </span>
          </div>
        </div>
        <div id="secure" class="echarts"></div>
      </div>
      <div class="businessBox">
        <div class="title">业务连续性</div>
        <ul class="list flexBox">
          <li class="flexBox">
            <img class="icon" src="./assets/img/icon_1.png" />
            <div class="info">
              <div class="name">业务平均响应时间</div>
              <div class="val">{{ tempData.delay }} ms</div>
            </div>
          </li>
          <li class="flexBox">
            <img class="icon" src="./assets/img/icon_2.png" />
            <div class="info">
              <div class="name">业务请求成功率</div>
              <div class="val">{{ tempData.successRate }} %</div>
            </div>
          </li>
        </ul>
        <ul class="infoList">
          <li class="other">
            <div class="row">
              <i class="icon"></i>
              <div class="txt">存在业务连续性风险站点: {{ tempData.riskStation }}</div>
            </div>
          </li>
          <li>
            <div class="row">
              <i class="icon"></i>
              <div class="txt">业务连续性风险: {{ tempData.riskDescription }}</div>
            </div>
          </li>
        </ul>
      </div>
      <div class="assetsBox">
        <div class="title">资产态势</div>
        <ul class="list row">
          <li>
            <div class="val">{{ tempData.totalTollStation }}</div>
            <div class="name">收费站</div>
          </li>
          <li>
            <div class="val">{{ tempData.serviceArea }}</div>
            <div class="name">服务区</div>
          </li>
          <li>
            <div class="val">{{ tempData.gantry }}</div>
            <div class="name">门架</div>
          </li>
          <li>
            <div class="val">{{ tempData.totalMileage }}公里</div>
            <div class="name">总里程</div>
          </li>
          <li>
            <div class="val">{{ tempData.totalAssets }}</div>
            <div class="name">总资产</div>
          </li>
          <li>
            <div class="val">{{ tempData.activeAssets }}</div>
            <div class="name">活跃资产</div>
          </li>
        </ul>
      </div>
    </div>
    <div class="map">
      <div id="map" class="mapBox"></div>
      <ul class="mapTag">
        <li class="flexBox">
          <i class="tag"></i>
          <div class="name">高速公路</div>
        </li>
        <li class="flexBox">
          <img class="icon" src="./assets/img/icon_11.png" />
          <div class="name">高速收费站</div>
        </li>
      </ul>
    </div>
  </div>
</template>

<script>
  import * as echarts from 'echarts'
  //rem 适配  1rem = 14px
  import './components/module/fontSize.js'
  //公共css
  import './assets/css/com/lib.css'
  import Socket from './components/module/webSocket.js'
  import { getToken } from '@/utils/token'
  import { proxyNet } from '@/config/index'
  const url = process.env.NODE_ENV == 'development' ? proxyNet : window.location.hostname
  const cookie = getToken()
  const protocol = window.location.protocol == 'http:' ? 'wss' : 'wss'
  export default {
    name: 'Index',
    data() {
      return {
        tempObj: {
          marker: [],
        },
        time: '',
        tempData: {},

        fitData: {
          /**
           * 屏幕大小
           */
          windowSize: {
            width: 1920,
          },
          /**
           * 参考字体号
           */
          fontSize: 14,
        },
      }
    },
    // created() {
    //   (function (doc, win) {
    //     var windowSizeWidth = 1920;
    //     var fontSize = 14;
    //     var docEl = doc.documentElement,
    //       resizeEvt = 'orientationchange' in window ? 'orientationchange' : 'resize',
    //       recalc = function () {
    //         var clientWidth = docEl.clientWidth;
    //         if (!clientWidth) return;
    //         docEl.style.fontSize = fontSize * (clientWidth / windowSizeWidth) + 'px';
    //       };
    //     if (!doc.addEventListener) return;
    //     win.addEventListener(resizeEvt, recalc, false);
    //     doc.addEventListener('DOMContentLoaded', recalc, false);
    //   })(document, window);
    // },
    mounted() {
      this.createSocket()
      this.updateTime()
      this.createMap()

      window.closeModel = (channelNum) => {
        this.closeModel(channelNum)
      }
    },
    methods: {
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
        return `${y}-${this.coverTime(m)}-${this.coverTime(d)} ${this.coverTime(h)}:${this.coverTime(
          mm
        )}:${this.coverTime(s)}`
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
        return this.fitData.fontSize * (clientWidth / this.fitData.windowSize.width) * num
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
      createMap() {
        this.map = new AMap.Map('map', {
          viewMode: '2D',
          zoom: 8,
          center: [122.35, 41.3],
          zooms: [8, 10],
          resizeEnable: true,
        })

        let imageLayer = new AMap.ImageLayer({
          url: new URL('./assets/img/map.png', import.meta.url).href,
          bounds: new AMap.Bounds([118.682181, 38.58066], [126.002, 43.614187]),
          opacity: 0.5,
        })
        this.map.add(imageLayer)
      },
      updateMarker(data, idx1) {
        const idx = this.tempObj.marker.length
        this.tempObj.marker[idx] = new AMap.Marker({
          content: `<div class="markerBox ${data.color}">${data.name}</div>`,
          position: new AMap.LngLat(data.x, data.y),
          offset: new AMap.Pixel(-this.getFitSize(1.24), -this.getFitSize(2.8)),
        })

        this.tempObj.marker[idx].on('click', (e) => {
          this.clickMarker(e, idx)
        })

        this.tempObj.marker[idx].setMap(this.map)
      },
      coverStar(num) {
        let txt = ''
        let len = num * 1
        for (let i = 0; i < 5; i++) {
          if (num - 2 > 0) {
            txt += "<li class='all'></li>"
          } else {
            if (num > 0) {
              txt += "<li class='half'></li>"
            } else {
              txt += '<li></li>'
            }
          }
          num -= 2
        }

        return txt
      },
      clickMarker(data, idx) {
        if (this.map) {
          let res = this.tempData.tollStationInfo[idx]
          let content = `
						<div class="markerModel">
							<div class="txt">${res.name}</div>
							<ul class="list row">
								<li>
									<div class="name">资产总数</div>
									<div class="val">${res.assetTotal}</div>
								</li>
								<li>
									<div class="name">活跃资产</div>
									<div class="val">${res.assetActive}</div>
								</li>
							</ul>
							<div class="starBox flexBox">
								<div class="name">安全指数</div>
								<ul class="starList flexBox">
									${this.coverStar(res.safetyIndex)}
								</ul>
							</div>
							<div class="infoList row">
								<li class="flexBox">
									<div class="name">安全检测</div>
									<div class="tag ${res.detection ? 'active' : ''}"></div>
								<li>
								<li class="flexBox">
									<div class="name">主动防御</div>
									<div class="tag ${res.defense ? 'active' : ''}"></div>
								<li>
							</div>
							<div class="close" onclick="closeModel()"></div>
						</div>
					`

          this.tempObj.infoWindow = new AMap.InfoWindow({
            isCustom: true,
            content: content,
            closeWhenClickMap: true,
            anchor: 'bottom-center',
          })
          this.tempObj.infoWindow.open(this.map, data.lnglat)
        }
      },
      closeModel() {
        if (this.tempObj.infoWindow) {
          this.tempObj.infoWindow.close()
          this.tempObj.infoWindow = null
        }
      },
      createSocket() {
        let $ws = `${protocol}://${url}/v3/ecsPlatform/websocket/screenHighway?EcsSessionId=${cookie}`

        this.tempObj.socket = new Socket($ws)
        this.tempObj.socket.on('open', (res) => {
          console.log('连接成功！')
        })
        this.tempObj.socket.on('message', (res) => {
          let msg = JSON.parse(res)
          if (msg.code == 20) {
            // if (this.tempData.securitySituation == undefined) {
            this.tempData = this.coverData(msg.data)
            this.updateSecure('secure', this.tempData.securitySituation)
            // };
          }
        })
        this.tempObj.socket.on('close', (res) => {
          console.log('关闭！', res)
        })
        this.tempObj.socket.on('error', (res) => {
          console.log('error', res)
        })
      },
      coverNum(num) {
        return num > 10000 ? `${(num / 10000).toFixed(2)}万` : num
      },
      coverData(data) {
        if (data.activeDefenseCount > 10000) {
          data.activeDefenseCountFormat = this.coverNum(data.activeDefenseCount)
        } else {
          data.activeDefenseCountFormat = data.activeDefenseCount
        }

        data.totalTollStation = this.coverNum(data.totalTollStation)
        data.serviceArea = this.coverNum(data.serviceArea)
        data.gantry = this.coverNum(data.gantry)
        data.totalMileage = this.coverNum(data.totalMileage)
        data.totalAssets = this.coverNum(data.totalAssets)
        data.activeAssets = this.coverNum(data.activeAssets)

        if (this.map && this.tempObj.marker) {
          for (let i = 0; i < this.tempObj.marker.length; i++) {
            this.tempObj.marker[i].setMap(null)
          }
          this.tempObj.marker = []
        }

        for (let i = 0; i < data.tollStationInfo.length; i++) {
          this.updateMarker(data.tollStationInfo[i], i)
        }

        return data
      },
      updateTime() {
        this.time = this.coverTimeToFormat(this.getTime())
        setInterval(() => {
          this.time = this.coverTimeToFormat(this.getTime())
        }, 1000)
      },
      updateSecure(name, data) {
        if (this.tempObj[name] == undefined) {
          this.tempObj[name] = echarts.init(document.getElementById(name))
        }

        let list = {}
        for (let i = 0; i < data.length; i++) {
          let msg = data[i]
          list[msg.name] = msg
        }
        let option = {
          legend: {
            orient: 'vertical',
            top: 'center',
            right: '10%',
            itemGap: this.getFitSize(1),
            itemWidth: this.getFitSize(0.86),
            itemHeight: this.getFitSize(0.86),
            data: data,
            formatter: (res) => {
              return `{a|${res}} {b|${this.coverNum(list[res].value)}}`
            },
            textStyle: {
              color: '#BEE4F8',
              fontSize: this.getFitSize(0.86),
              lineHeight: this.getFitSize(1),
              rich: {
                c: {
                  color: '#97AFD8',
                  fontSize: this.getFitSize(0.86),
                },
                b: {
                  color: '#fff',
                  fontSize: this.getFitSize(1),
                  padding: [0, 0, 0, 5],
                },
              },
            },
          },
          grid: {
            top: '5%',
            bottom: '5%',
            right: '0%',
            left: '0%',
            containLabel: true,
          },
          tooltip: {
            trigger: 'item',
          },
          color: ['#1791ff', '#37cbcb', '#4dcc73', '#fbd438', '#f2637b', '#8085E9'],
          series: [
            {
              type: 'pie',
              radius: ['30%', '80%'],
              center: ['25%', '50%'],
              // roseType: 'area',
              label: false,
              data: data,
            },
          ],
        }

        this.updateEcharts(this.tempObj[name], option)
      },
    },
  }
</script>

<style>
  @import url('./assets/css/index.css');
</style>
