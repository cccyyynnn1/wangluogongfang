/**
 * 通讯类
 */
import axios from 'axios'
import GlobalInfo from '../GlobalInfo.vue'
import tools from './tools.vue'
axios.defaults.timeout = 5000
axios.defaults.baseURL = GlobalInfo.HTTPURL

/**
 * 正常请求
 * @param {*请求类型:post、get、put、del} type
 * @param {*请求URL地址} url
 * @param {*数据} data
 */
export function request(type, url, data = {}) {
  return new Promise((resolve, reject) => {
    axios[type](url, data).then(
      (response) => {
        // if (response.data.code != 200) {
        // 	message.error(response.data.message)
        // };
        resolve(response.data)
      },
      (err) => {
        reject(err)
      }
    )
  })
}

/**
 * 请求携带Token
 * @param {*请求类型:post、get、put、del} type
 * @param {*请求URL地址} url
 * @param {*数据} data
 */
export function requestToken(type, url, data = {}) {
  return new Promise((resolve, reject) => {
    axios[type](url, data, {
      headers: {
        'X-Access-Token': JSON.parse(tools.getSessionStorage()).token,
      },
    }).then(
      (response) => {
        if (response.data.code != 200) {
          message.error(response.data.message)
        }
        resolve(response.data)
      },
      (err) => {
        if (err.response.status == 500) {
          requestTokenPastDue()
        } else {
          reject(err)
        }
      }
    )
  })
}
/**
 * Token过期
 */
function requestTokenPastDue() {
  window.location.href = GlobalInfo.initView
}
