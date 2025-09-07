<script setup lang="ts">
  import { downloadLogPacket } from '@/utils/download'

  import { favoritesLogPageApi, appendCollectTraceApi } from '@/api-ecs/retrieve'

  import { formatNstime, formatTime } from '@/utils/time'

  import { useSettingsStore } from '@/store/modules/settings'

  import { PacketDecodeQuery } from '~/src/types'

  import VabDialog from '@/plugins/VabDialog/index.vue'

  import DecodingPacket from '@/ecs/alert/components/traceability-decoding-packet.vue'

  import DecodingStream from '@/ecs/alert/components/traceability-decoding-stream.vue'

  import dayjs from 'dayjs'
  import { getAssetsPreviewListApi } from '~/src/api-ecs/assets-preview'
  import { PacketsMagic } from '@/utils/magic'
  import { hexStringToArrayBuffer } from '@/utils/text'
  import { getPacketDecodeFlowDecodeApi } from '~/src/api-ecs/alert'

  const settingsStore = useSettingsStore()

  const props = defineProps<{
    infoVal: any
    workspaceId?: number
  }>()

  const selects = ref<string[]>([])
  const decodingData = ref<PacketDecodeQuery>()
  const $baseMessage: any = inject('$baseMessage')
  const decodingVisible = ref(false)

  const decodingTime = ref('')

  const allData: { [key: number]: { lable: string; key: string; span: number }[] } = {
    12: [
      { lable: '源IP', key: 'clientIp', span: 1 },
      { lable: '源端口', key: 'clientPort', span: 1 },
      { lable: '目的IP', key: 'serverIp', span: 1 },
      { lable: '目的端口', key: 'serverPort', span: 1 },
      { lable: '原始客户端IP', key: 'realClientIp', span: 1 },
      { lable: '原始客户端端口', key: 'realClientPort', span: 1 },
      { lable: '原始服务端IP', key: 'realServerIp', span: 1 },
      { lable: '原始服务端端口', key: 'realServerPort', span: 1 },
      { lable: '开始时间', key: 'startTimeNs', span: 1 },
      { lable: '响应码', key: 'statusCode', span: 1 },
      { lable: '请求方式', key: 'method', span: 1 },
      { lable: '方向', key: 'direction', span: 1 },
      { lable: '请求内容长度', key: 'reqContentLength', span: 1 },
      { lable: '请求内容类型', key: 'reqContentType', span: 1 },
      { lable: '响应内容长度', key: 'respContentLength', span: 1 },
      { lable: '响应内容类型', key: 'respContentType', span: 1 },
      { lable: '版本', key: 'version', span: 1 },
      { lable: '总字节数', key: 'totalBytes', span: 1 },
      { lable: '客户端字节数', key: 'clientBytes', span: 1 },
      { lable: '服务端字节数', key: 'serverBytes', span: 1 },
      { lable: '主机名', key: 'host', span: 2 },
      { lable: '文件名称', key: 'fileName', span: 2 },
      { lable: 'MD5', key: 'fileMd5', span: 2 },
      { lable: '文件类型', key: 'fileType', span: 2 },
      { lable: 'URL', key: 'url', span: 2 },
      { lable: 'UA', key: 'UA', span: 2 },
      { lable: '其他', key: 'other', span: 2 },
    ],
    13: [
      { lable: '源IP', key: 'clientIp', span: 1 },
      { lable: '源端口', key: 'clientPort', span: 1 },
      { lable: '目的IP', key: 'serverIp', span: 1 },
      { lable: '目的端口', key: 'serverPort', span: 1 },
      { lable: '原始客户端IP', key: 'realClientIp', span: 1 },
      { lable: '原始客户端端口', key: 'realClientPort', span: 1 },
      { lable: '原始服务端IP', key: 'realServerIp', span: 1 },
      { lable: '原始服务端端口', key: 'realServerPort', span: 1 },
      { lable: '开始时间', key: 'startTimeNs', span: 1 },
      { lable: '响应码', key: 'responseCode', span: 1 },
      { lable: '请求方式', key: 'method', span: 1 },
      { lable: '总字节数', key: 'totalBytes', span: 1 },
      { lable: '客户端字节数', key: 'clientBytes', span: 1 },
      { lable: '服务端字节数', key: 'serverBytes', span: 1 },
      { lable: '用户名', key: 'username', span: 2 },
      { lable: '密码', key: 'password', span: 2 },
      { lable: 'URL', key: 'url', span: 2 },
    ],
    14: [
      { lable: '源IP', key: 'clientIp', span: 1 },
      { lable: '源端口', key: 'clientPort', span: 1 },
      { lable: '目的IP', key: 'serverIp', span: 1 },
      { lable: '目的端口', key: 'serverPort', span: 1 },
      { lable: '开始时间', key: 'startTimeNs', span: 1 },
      { lable: '总字节数', key: 'totalBytes', span: 1 },
      { lable: '客户端字节数', key: 'clientBytes', span: 1 },
      { lable: '服务端字节数', key: 'serverBytes', span: 1 },
      { lable: '操作', key: 'opCode', span: 1 },
      { lable: '错误码', key: 'errorCode', span: 1 },
      { lable: '传输模式', key: 'transModel', span: 2 },
      { lable: '文件名称', key: 'fileName', span: 2 },
      { lable: '文件类型', key: 'fileType', span: 2 },
      { lable: '错误信息', key: 'errorInfo', span: 2 },
    ],
    15: [
      { lable: '源IP', key: 'clientIp', span: 1 },
      { lable: '源端口', key: 'clientPort', span: 1 },
      { lable: '目的IP', key: 'serverIp', span: 1 },
      { lable: '目的端口', key: 'serverPort', span: 1 },
      { lable: '开始时间', key: 'startTimeNs', span: 1 },
      { lable: '总字节数', key: 'totalBytes', span: 1 },
      { lable: '客户端字节数', key: 'clientBytes', span: 1 },
      { lable: '服务端字节数', key: 'serverBytes', span: 1 },
      { lable: '类型', key: 'type', span: 1 },
      { lable: '长度', key: 'length', span: 1 },
      { lable: '代码', key: 'code', span: 2 },
      { lable: 'ICMP数据', key: 'packetData', span: 2 },
    ],
    16: [
      { lable: '源IP', key: 'clientIp', span: 1 },
      { lable: '源端口', key: 'clientPort', span: 1 },
      { lable: '目的IP', key: 'serverIp', span: 1 },
      { lable: '目的端口', key: 'serverPort', span: 1 },
      { lable: '原始客户端IP', key: 'realClientIp', span: 1 },
      { lable: '原始客户端端口', key: 'realClientPort', span: 1 },
      { lable: '原始服务端IP', key: 'realServerIp', span: 1 },
      { lable: '原始服务端端口', key: 'realServerPort', span: 1 },
      { lable: '开始时间', key: 'startTimeNs', span: 1 },
      { lable: '总字节数', key: 'totalBytes', span: 1 },
      { lable: '客户端字节数', key: 'clientBytes', span: 1 },
      { lable: '服务端字节数', key: 'serverBytes', span: 1 },
      { lable: 'COOKIE', key: 'cookie', span: 2 },
      { lable: '错误码', key: 'failureCode', span: 2 },
      { lable: '安全协议', key: 'securityProtocol', span: 2 },
      { lable: '频道', key: 'channelDef', span: 2 },
      { lable: '键盘布局', key: 'keyboardLayout', span: 2 },
      { lable: '键盘类型', key: 'keyboardType', span: 2 },
      { lable: '客户端版本', key: 'clientBuild', span: 2 },
      { lable: '客户端名称', key: 'clientName', span: 2 },
      { lable: '客户端产品ID', key: 'clientProductId', span: 2 },
      { lable: '桌面宽度', key: 'desktopWidth', span: 2 },
      { lable: '桌面高度', key: 'desktopHeight', span: 2 },
      { lable: '颜色深度', key: 'colorDepth', span: 2 },
      { lable: '加密级别', key: 'encryptionLevel', span: 2 },
      { lable: '加密方法', key: 'encodingMethod', span: 2 },
      { lable: '证书时效', key: 'certPermanent', span: 2 },
      { lable: '证书类型', key: 'certType', span: 2 },
    ],
    17: [
      { lable: '源IP', key: 'clientIp', span: 1 },
      { lable: '源端口', key: 'clientPort', span: 1 },
      { lable: '目的IP', key: 'serverIp', span: 1 },
      { lable: '目的端口', key: 'serverPort', span: 1 },
      { lable: '总字节数', key: 'totalBytes', span: 1 },
      { lable: '客户端字节数', key: 'clientBytes', span: 1 },
      { lable: '服务端字节数', key: 'serverBytes', span: 1 },
      { lable: 'SOCKS版本', key: 'version', span: 2 },
      { lable: '命令', key: 'cmd', span: 2 },
      { lable: '目的主机地址类型', key: 'dstAddrType', span: 2 },
      { lable: '目的主机地址', key: 'dstHostIp', span: 2 },
      { lable: '目的主机端口', key: 'dstHostPort', span: 2 },
      { lable: '目的主机域名', key: 'dstHostname', span: 2 },
      { lable: '应答', key: 'respType', span: 2 },
      { lable: '服务器绑定地址类型', key: 'servBindType', span: 2 },
      { lable: '服务器绑定地址', key: 'servBindAddr', span: 2 },
      { lable: '服务器绑定端口', key: 'servBindPort', span: 2 },
      { lable: 'USER', key: 'username', span: 2 },
      { lable: '加密级别', key: 'encryptionLevel', span: 2 },
      { lable: '加密方法', key: 'encodingMethod', span: 2 },
      { lable: '服务器选择认证方法', key: 'method', span: 2 },
      { lable: '客户端选择认证方法', key: 'methods', span: 2 },
    ],
    18: [
      { lable: '源IP', key: 'clientIp', span: 1 },
      { lable: '源端口', key: 'clientPort', span: 1 },
      { lable: '目的IP', key: 'serverIp', span: 1 },
      { lable: '目的端口', key: 'serverPort', span: 1 },
      { lable: '原始客户端IP', key: 'realClientIp', span: 1 },
      { lable: '原始客户端端口', key: 'realClientPort', span: 1 },
      { lable: '原始服务端IP', key: 'realServerIp', span: 1 },
      { lable: '原始服务端端口', key: 'realServerPort', span: 1 },
      { lable: '开始时间', key: 'startTimeNs', span: 1 },
      { lable: '总字节数', key: 'totalBytes', span: 1 },
      { lable: '客户端字节数', key: 'clientBytes', span: 1 },
      { lable: '服务端字节数', key: 'serverBytes', span: 1 },
      { lable: '客户端版本', key: 'clientVersion', span: 2 },
      { lable: '服务端版本', key: 'serverVersion', span: 2 },
      { lable: '客户端cookie', key: 'clientCookie', span: 2 },
      { lable: '服务端cookie', key: 'serverCookie', span: 2 },
      { lable: '客户端hash指纹', key: 'clientHassh', span: 2 },
      { lable: '服务端hash指纹', key: 'serverHassh', span: 2 },
      { lable: '十六进制公钥', key: 'hexHostKey', span: 2 },
      { lable: '十六进制公钥验证', key: 'hexHostSignature', span: 2 },
      { lable: '服务端密钥MD5', key: 'md5Hash', span: 2 },
      { lable: '服务端密钥SHA1', key: 'sha1Hash', span: 2 },
      { lable: '服务端密钥类型', key: 'keyType', span: 2 },
    ],
    19: [
      { lable: '源IP', key: 'clientIp', span: 1 },
      { lable: '源端口', key: 'clientPort', span: 1 },
      { lable: '目的IP', key: 'serverIp', span: 1 },
      { lable: '目的端口', key: 'serverPort', span: 1 },
      { lable: '原始客户端IP', key: 'realClientIp', span: 1 },
      { lable: '原始客户端端口', key: 'realClientPort', span: 1 },
      { lable: '原始服务端IP', key: 'realServerIp', span: 1 },
      { lable: '原始服务端端口', key: 'realServerPort', span: 1 },
      { lable: '开始时间', key: 'startTimeNs', span: 1 },
      { lable: '总字节数', key: 'totalBytes', span: 1 },
      { lable: '客户端字节数', key: 'clientBytes', span: 1 },
      { lable: '服务端字节数', key: 'serverBytes', span: 1 },
      { lable: '版本', key: 'version', span: 2 },
      { lable: '证书类型', key: 'caType', span: 2 },
      { lable: '密码套件', key: 'cliCipher', span: 2 },
      { lable: '密码套件数量', key: 'cliCipherCnt', span: 2 },
      { lable: '密码套件长度', key: 'cliCipherLen', span: 2 },
      { lable: '扩展长度', key: 'cliExtLen', span: 2 },
      { lable: '扩展SNI', key: 'cliExtSni', span: 2 },
      { lable: '扩展数量', key: 'cliExtCnt', span: 2 },
      { lable: '证书链数量', key: 'certChainCnt', span: 2 },
      { lable: '证书链', key: 'certChain', span: 2 },
      { lable: '客户端JA3指纹', key: 'cliJa3', span: 2 },
      { lable: '服务端JA3指纹', key: 'servJa3', span: 2 },
    ],
    20: [
      { lable: '源IP', key: 'clientIp', span: 1 },
      { lable: '源端口', key: 'clientPort', span: 1 },
      { lable: '目的IP', key: 'serverIp', span: 1 },
      { lable: '目的端口', key: 'serverPort', span: 1 },
      // { lable: '原始客户端IP', key: 'realClientIp', span: 1 },
      // { lable: '原始客户端端口', key: 'realClientPort', span: 1 },
      // { lable: '原始服务端IP', key: 'realServerIp', span: 1 },
      // { lable: '原始服务端端口', key: 'realServerPort', span: 1 },
      { lable: '开始时间', key: 'startTimeNs', span: 1 },
      { lable: '总字节数', key: 'totalBytes', span: 1 },
      { lable: '客户端字节数', key: 'clientBytes', span: 1 },
      { lable: '服务端字节数', key: 'serverBytes', span: 1 },
      { lable: '命令', key: 'data', span: 2 },
      { lable: '用户名', key: 'username', span: 2 },
      { lable: 'TELNET数据', key: 'pakcetData', span: 2 },
    ],
    21: [
      { lable: '源IP', key: 'clientIp', span: 1 },
      { lable: '源端口', key: 'clientPort', span: 1 },
      { lable: '目的IP', key: 'serverIp', span: 1 },
      { lable: '目的端口', key: 'serverPort', span: 1 },
      { lable: '原始客户端IP', key: 'realClientIp', span: 1 },
      { lable: '原始客户端端口', key: 'realClientPort', span: 1 },
      { lable: '原始服务端IP', key: 'realServerIp', span: 1 },
      { lable: '原始服务端端口', key: 'realServerPort', span: 1 },
      { lable: '开始时间', key: 'startTimeNs', span: 1 },
      { lable: '总字节数', key: 'totalBytes', span: 1 },
      { lable: '客户端字节数', key: 'clientBytes', span: 1 },
      { lable: '服务端字节数', key: 'serverBytes', span: 1 },
      { lable: '用户名', key: 'username', span: 2 },
      { lable: '密码', key: 'password', span: 2 },
      { lable: '状态', key: 'state', span: 2 },
      { lable: '会话ID', key: 'sessionId', span: 2 },
    ],
    25: [
      { lable: '源IP', key: 'clientIp', span: 1 },
      { lable: '源端口', key: 'clientPort', span: 1 },
      { lable: '目的IP', key: 'serverIp', span: 1 },
      { lable: '目的端口', key: 'serverPort', span: 1 },
      { lable: '开始时间', key: 'startTimeNs', span: 1 },
      { lable: '续订时间', key: 'renewalTime', span: 1 },
      { lable: '重新绑定时间', key: 'rebindingTime', span: 1 },
      { lable: '客户端MAC地址', key: 'clientMac', span: 1 },
      { lable: '下一个服务器IP地址', key: 'nextServer', span: 1 },
      { lable: '设备Id', key: 'deviceId', span: 1 },
      { lable: '请求类型', key: 'requestType', span: 1 },
      { lable: '租用时间', key: 'ipLeaseTime', span: 1 },
      { lable: '时间服务器', key: 'networkTimeServer', span: 1 },
      { lable: '中继服务器地址', key: 'relayServerIp', span: 1 },
      { lable: 'DNS服务器', key: 'dnsServer', span: 1 },
      { lable: '子网掩码', key: 'subnetMask', span: 1 },
      { lable: '服务器标识', key: 'serverIdentifier', span: 1 },
      { lable: '跳数', key: 'hops', span: 1 },
      { lable: '经过秒数', key: 'secondsElapsed', span: 1 },
      { lable: '链路', key: 'probeIds', span: 1 },
    ],
    26: [
      { lable: '链路', key: 'probeIds', span: 1 },
      { lable: 'Agent', key: 'agentId', span: 1 },
      { lable: '客户端网段', key: 'clientNetSegmentIds', span: 1 },
      { lable: '服务端网段', key: 'serverNetSegmentIds', span: 1 },
      { lable: '源IP', key: 'clientIp', span: 1 },
      { lable: '目的IP', key: 'serverIp', span: 1 },
      { lable: '客户端端口', key: 'clientPort', span: 1 },
      { lable: '服务端端口', key: 'serverPort', span: 1 },
      { lable: '原始客户端IP', key: 'realClientIp', span: 1 },
      { lable: '原始客户端端口', key: 'realClientPort', span: 1 },
      { lable: '原始服务端IP', key: 'realServerIp', span: 1 },
      { lable: '原始服务端端口', key: 'realServerPort', span: 1 },
      { lable: '请求时间', key: 'requestTimeNs', span: 1 },
      { lable: '响应时间', key: 'responseTimeNs', span: 1 },
      { lable: '结束时间', key: 'byeTimeNs', span: 1 },
      { lable: '呼叫ID', key: 'callId', span: 1 },
      { lable: '主叫ID', key: 'callerId', span: 1 },
      { lable: '主叫IP', key: 'callerIp', span: 1 },
      { lable: '被叫ID', key: 'calleeId', span: 1 },
      { lable: '被叫IP', key: 'calleeIp', span: 1 },
      { lable: '呼叫用户', key: 'userName', span: 1 },
      { lable: '采样率', key: 'rate', span: 1 },
      { lable: '频道', key: 'channel', span: 1 },
      { lable: '帧率', key: 'fps', span: 1 },
      { lable: '媒体编码', key: 'name', span: 1 },
      { lable: '媒体类型', key: 'media', span: 1 },
      { lable: '总字节数', key: 'totalBytes', span: 1 },
      { lable: '客户端字节数', key: 'clientBytes', span: 1 },
      { lable: '服务端字节数', key: 'serverBytes', span: 1 },
      { lable: '虚拟网类型', key: 'vlanType', span: 1 },
      { lable: '虚拟网标识', key: 'vlanId', span: 1 },
      { lable: '第二层虚拟网标识', key: 'vlanId2', span: 1 },
      { lable: '虚拟资产标识', key: 'vai', span: 1 },
    ],
    28: [
      { lable: '日期', key: 'date', span: 1 },
      { lable: '开始时间', key: 'startTimeNs', span: 1 },
      { lable: '源IP', key: 'clientIp', span: 1 },
      { lable: '源端口', key: 'clientPort', span: 1 },
      { lable: '目的IP', key: 'serverIp', span: 1 },
      { lable: '目的端口', key: 'serverPort', span: 1 },
      { lable: '根源延迟', key: 'serverPort', span: 1 },
      { lable: '报告间隔指数', key: 'poll', span: 1 },
      { lable: '时间戳精度', key: 'precision', span: 1 },
      { lable: '发送时间戳', key: 'originateTimestamp', span: 1 },
      { lable: '参考时间戳', key: 'referenceTimestamp', span: 1 },
      { lable: '接收时间戳', key: 'receiveTimestamp', span: 1 },
      { lable: '发送回应时间戳', key: 'transmitTimestamp', span: 1 },
      { lable: '版本', key: 'version', span: 1 },
    ],
    32: [
      { lable: '客户端', key: 'clientIp', span: 1 },
      { lable: '服务端', key: 'serverIp', span: 1 },
      { lable: '客户端端口', key: 'clientPort', span: 1 },
      { lable: '服务端端口', key: 'serverPort', span: 1 },
      { lable: 'Tracker地址', key: 'announce', span: 1 },
      { lable: 'Tracker备用地址', key: 'announceList', span: 1 },
      { lable: '种子创建时间', key: 'creationDate', span: 1 },
      { lable: '客户端软件', key: 'createdBy', span: 1 },
      { lable: '资源名称', key: 'infoName', span: 1 },
      { lable: '文件分片大小', key: 'infoPieceLengthStr', span: 1 },
      { lable: '文件大小', key: 'infoLengthStr', span: 1 },
      { lable: '文件信息', key: 'infoFiles', span: 1 },
      { lable: '文件哈希', key: 'infoHash', span: 1 },
      { lable: '客户端ID', key: 'peerId', span: 1 },
      { lable: '自身IP', key: 'ip', span: 1 },
      { lable: '自身端口', key: 'port', span: 1 },
      { lable: '跟踪器ID', key: 'trackerld', span: 1 },
      { lable: '事件类型', key: 'event', span: 1 },
    ],
  }

  const datalist = ref<{ lable: string; key: string; span: number }[]>([])

  const initData = () => {
    const index = props.infoVal.indexType
    if (index) {
      datalist.value = allData[index]
    }
  }
  // PCAP保存
  function download() {
    downloadLogPacket(props.infoVal, 1)
  }
  const traceabilityHandle = async () => {
    const { indexType, ...obj } = props.infoVal
    const { data } = await appendCollectTraceApi({
      dataLog: JSON.stringify(obj),
      indexType,
      workspaceId: props.workspaceId!,
    })
    $baseMessage('添加溯源图成功', 'success', 'vab-hey-message-success')
  }
  const favoritesHandle = async () => {
    const { indexType, ...obj } = props.infoVal
    const { data } = await favoritesLogPageApi({
      dataLog: JSON.stringify(obj),
      indexType,
      workspaceId: props.workspaceId!,
    })
    $baseMessage('添加收藏夹成功', 'success', 'vab-hey-message-success')
  }
  async function getSelect() {
    const { clientAssets, serverAssets, clientArea, serverArea, position } = props.infoVal
    selects.value = []
    clientAssets && selects.value.push(clientAssets)
    serverAssets && selects.value.push(serverAssets)
    clientArea && selects.value.push(clientArea)
    serverArea && selects.value.push(serverArea)
    position && selects.value.push(position)
  }

  // PCAP解析
  const decodingHandle = () => {
    const { clientIp, clientPort, serverIp, serverPort, startTimeNs } = props.infoVal
    const timeDate = dayjs(formatNstime(startTimeNs))
    const startDate = timeDate.subtract(1, 'minute')
    const endDate = timeDate.add(1, 'minute')
    decodingTime.value = `${formatTime(startDate)} - ${formatTime(endDate)}`
    decodingData.value = {
      clientIp,
      clientPort,
      serverIp,
      serverPort,
    }
    decodingVisible.value = true
  }
  const protocol = ref('TCP')
  onMounted(() => {
    getSelect()
    initData()
    getAssetName()
    if (props.infoVal.indexType == 14 || props.infoVal.indexType == 25 || props.infoVal.indexType == 26) {
      protocol.value = 'UDP'
    }
    getCode()
  })

  const clientIPAssetName = ref()
  const serverIPAssetName = ref()
  const getAssetName = async () => {
    const { clientIp, serverIp } = props.infoVal
    const arr: string[] = []
    if (clientIp) {
      arr.push(clientIp)
    }
    if (serverIp) {
      arr.push(clientIp)
    }
    const {
      data: { records },
    } = await getAssetsPreviewListApi({
      pageNum: 1,
      pageSize: 10,
      searchIps: arr,
      dataSourceId: [] as number[],
      otherDataSource: false,
      searchStr: '',
    })
    if (clientIp) {
      clientIPAssetName.value = records[0]?.businessName || records[0]?.appName
    }
    if (serverIp) {
      serverIPAssetName.value = records[records.length - 1]?.businessName || records[records.length - 1]?.appName
    }
  }
  // 数据流功能
  const formatType = ref('utf-8')
  const typeList = [
    {
      lable: 'UTF-8',
      value: 'utf-8',
    },
    {
      lable: 'UTF-16BE',
      value: 'utf-16be',
    },
    {
      lable: 'UTF-16LE',
      value: 'utf-16le',
    },
  ]
  const magicIcon = require('@/assets/mofabang.svg')
  const protocols = ['TCP', 'UDP']

  let STREAM_DATA = [] as { hexs: string[]; info: string; srcToDst: boolean }[]
  const protocolChange = (protocolVal: string) => {
    protocol.value = protocolVal
    getCode()
  }
  const showMagic = ref(false)
  const streamList = ref<{ hexs: string[]; info: string; srcToDst: boolean }[]>([])
  const magicStreamList = ref<{ hexs: string[]; info: string; srcToDst: boolean }[]>([])
  const toMagic = () => {
    if (showMagic.value) return (showMagic.value = false)
    magicStreamList.value = streamList.value.map((stream) => ({
      ...stream,
      hexs: stream.hexs.map((str: string) => PacketsMagic(str)),
    }))
    showMagic.value = true
  }
  const loading = ref(false)
  const getCode = async () => {
    showMagic.value = false
    if (!props.infoVal) return
    try {
      loading.value = true
      const { probeId = '', clientPort, serverIp, serverPort, clientIp, startTimeNs } = props.infoVal
      const timeDate = dayjs(formatNstime(startTimeNs))
      const startDate = timeDate.subtract(1, 'minute')
      const endDate = timeDate.add(1, 'minute')
      decodingTime.value = `${formatTime(startDate)} - ${formatTime(endDate)}`
      const { data } = await getPacketDecodeFlowDecodeApi({
        top: 50,
        query: {
          objectList: [
            {
              serverIp,
              serverPort: serverPort?.toString(),
              clientIp,
              probeId,
              clientPort: clientPort?.toString(),
            },
          ],
          timeStep: 'minuteStep',
          timeRange: decodingTime.value,
        },
        flowType: protocol.value,
      })
      STREAM_DATA = data || []
      toFormat(STREAM_DATA)
    } finally {
      loading.value = false
    }
  }
  const toFormat = (streamData: any[]) => {
    streamList.value = streamData.map((stream) => ({
      ...stream,
      hexs: stream.hexs.map((hex: string) => formatText(hex)),
    }))
  }
  const formatText = (txt: string) => {
    const textDecoder = new TextDecoder(formatType.value)
    const str = txt.replace(/\s*/g, '')
    return textDecoder.decode(hexStringToArrayBuffer(str))
  }
  defineExpose({
    getAssetName,
    getCode,
  })

  const formatBitps = (val: number) => {
    let value = ''
    if (val > 1073000000) {
      value = `${(val / 1024 / 1024 / 1024).toFixed(2)}GB`
    } else if (val > 1048000) {
      value = `${(val / 1024 / 1024).toFixed(2)}MB`
    } else if (val > 1024) {
      value = `${(val / 1024).toFixed(2)}KB`
    } else {
      value = String(val) && `${val}B`
    }
    return value
  }
</script>

<script lang="ts">
  export default {
    name: 'DetailOther',
  }
</script>
<template>
  <div class="detail-other">
    <div class="marginB15" style="display: flex; justify-content: space-between; align-items: center">
      <div class="hostTable_title">基本信息</div>
      <div>
        <el-space :size="10">
          <el-button v-if="props.workspaceId" type="primary" @click="traceabilityHandle">添加到溯源图</el-button>
          <el-button v-if="props.workspaceId" type="primary" @click="favoritesHandle">添加到收藏夹</el-button>
          <!-- <el-button type="primary" @click="() => settingsStore.changeToolboxVisible(true)">工具箱</el-button> -->
          <el-button type="primary" @click="decodingHandle">PCAP解析</el-button>
          <el-button type="primary" @click="download">PCAP保存</el-button>
        </el-space>
      </div>
    </div>
    <el-descriptions border :column="2">
      <template v-for="(item, index) in datalist" :key="index">
        <el-descriptions-item v-if="item.lable == '开始时间'" :label="item.lable" :span="item.span">
          <div v-copy="formatNstime(infoVal[item.key])" style="padding-right: 30px">
            {{ formatNstime(infoVal[item.key]) }}
          </div>
        </el-descriptions-item>
        <el-descriptions-item v-else-if="item.key == 'clientIp'" :label="item.lable" :span="item.span">
          <div v-copy="infoVal?.clientIp" style="padding-right: 30px">
            {{ infoVal?.clientIp }}
            <span style="margin-left: 6px; color: #4637a4">{{ clientIPAssetName }}</span>
          </div>
        </el-descriptions-item>
        <el-descriptions-item v-else-if="item.key == 'serverIp'" :label="item.lable" :span="item.span">
          <div v-copy="infoVal?.serverIp" style="padding-right: 30px">
            {{ infoVal?.serverIp }}
            <span style="margin-left: 6px; color: #4637a4">{{ serverIPAssetName }}</span>
          </div>
        </el-descriptions-item>
        <el-descriptions-item v-else-if="item.lable.search('字节数') > 0" :label="item.lable" :span="item.span">
          <div v-copy="formatBitps(infoVal[item.key])" style="padding-right: 30px">
            {{ formatBitps(infoVal[item.key]) }}
          </div>
        </el-descriptions-item>
        <el-descriptions-item v-else :label="item.lable" :span="item.span">
          <div v-copy="infoVal[item.key]" style="padding-right: 30px">{{ infoVal[item.key] }}</div>
        </el-descriptions-item>
      </template>
    </el-descriptions>
    <div class="hostTable_title">标签信息</div>
    <el-card shadow="never">
      <div v-for="(item, $index) in selects" :key="$index" class="myElTag">
        {{ item }}
      </div>
      <div v-if="selects.length === 0" style="text-align: center">暂无数据</div>
    </el-card>
    <!-- <div class="hostTable_title" style="display: flex; justify-content: space-between; align-items: center">
      <span>数据流</span>
      <div class="my-bts">
        <el-button-group>
          <el-button
            v-for="my_type of typeList"
            :key="my_type.value"
            size="small"
            :type="formatType === my_type.value ? 'primary' : 'default'"
            @click="formatType = my_type.value"
          >
            {{ my_type.lable }}
          </el-button>
        </el-button-group>
        <el-button-group style="margin-left: 15px; border: 0; background-color: transparent">
          <el-button
            v-for="my_protocol of protocols"
            :key="my_protocol"
            size="small"
            :type="protocol === my_protocol ? 'primary' : 'default'"
            @click="protocolChange(my_protocol)"
          >
            {{ my_protocol }}
          </el-button>
        </el-button-group>
        <el-button :auto-insert-space="false" style="border: 0; background-color: transparent" @click="toMagic()">
          <el-image :src="magicIcon" style="width: 16px; height: 16px; cursor: pointer" />
        </el-button>
      </div>
    </div>
    <div v-loading="loading" class="requests-and-responses">
      <div v-for="(stream, index) in showMagic ? magicStreamList : streamList" :key="index" class="stream">
        <h4>{{ stream.info }}</h4>
        <ul :class="{ isDst: stream.srcToDst }">
          <li v-for="(item, $index) in stream.hexs" :key="$index">{{ item }}</li>
        </ul>
      </div>
      <el-empty v-if="streamList.length === 0" description="暂无数据" />
    </div> -->
    <vab-dialog
      v-model="decodingVisible"
      :close-on-click-modal="false"
      destroy-on-close
      show-fullscreen
      title="详情"
      width="1175px"
    >
      <el-tabs :model-value="'decodingStream'">
        <el-tab-pane label="数据流" name="decodingStream">
          <decoding-stream
            :current-data="decodingData"
            :protocol="[14, 25, 26].includes(infoVal.indexType) ? 'UDP' : ''"
            :time-date="decodingTime"
          />
        </el-tab-pane>
        <el-tab-pane label="数据包" lazy name="decodingPacket">
          <decoding-packet :current-data="decodingData" :time-date="decodingTime" />
        </el-tab-pane>
      </el-tabs>
    </vab-dialog>
  </div>
</template>

<style scoped lang="scss">
  .requests-and-responses {
    height: calc(100vh - 60px);
    width: 100%;
    border: 1px solid var(--el-border-color-lighter);
    border-radius: var(--el-input-border-radius, var(--el-border-radius-base));
    overflow-y: auto;
    &::-webkit-scrollbar {
      width: 0;
      height: 0;
    }
    padding: 20px;
    .data-warp {
      overflow-y: auto;
      text-overflow: ellipsis;
      word-break: break-all;
      word-wrap: break-word;
    }
  }
  .myElTag {
    background-color: #d9ecff;
    border-color: #d9ecff;
    color: var(--el-color-primary);
    background-color: rgba(102, 85, 231, 0.08);
    word-break: break-all;
    word-wrap: break-word;
    white-space: normal;
    display: inline-block;
    vertical-align: text-top;
    margin: 0 4px;
    padding: 3px 11px;
    border-radius: 4px;
  }
  .el-descriptions {
    margin-bottom: 15px;
    clear: both;

    :deep() {
      .el-descriptions__label {
        width: 172px;
      }

      .el-descriptions__content {
        width: 305px;
        overflow: hidden;
        word-break: break-all;
        // white-space: nowrap;
        word-wrap: break-word;
      }
    }
  }

  .el-space--horizontal {
    float: right;
    margin: 5px -10px 15px 0;
  }

  .hostTable_title {
    height: 14px;
    line-height: 14px;
    position: relative;
    text-indent: 0.8em;
    font-weight: 700;
    color: #303133;
    margin-bottom: 15px;

    &::before {
      content: ' ';
      display: inline-block;
      position: absolute;
      width: 3px;
      height: 100%;
      left: 0;
      background: var(--el-color-primary);
      margin-right: 5px;
    }
  }
</style>
