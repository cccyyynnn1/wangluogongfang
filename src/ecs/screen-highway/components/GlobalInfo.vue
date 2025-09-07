/** * 全局配置 */
<script>
  import { proxyNet } from '@/config/index'
  export default {
    /**项目名称*/
    projectName: 'highSpeed',
    /**请求地址*/
    HTTPURL: 'http://47.254.140.79/',
    /**ws地址*/
    ws: `wss://${proxyNet}/v3/ecsPlatform/websocket/screenHighway?EcsSessionId=6bba43f820553fc9bf1dc42b6b1f0de84646517b`,
    /**
     * 账号信息
     */
    account: {
      username: '',
      password: '',
    },
    /**
     * 适配
     */
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
    /**
     * 页面刷新率
     */
    updatePage: 3000,
    /**
     * Token失效后打开的页面
     */
    tokenPastDueToPage: '/login',
  }
</script>
