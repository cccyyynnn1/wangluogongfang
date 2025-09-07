<script lang="ts">
  export default {
    name: 'EditSite',
  }
</script>
<script setup lang="ts">
  import VueDraggable from 'vuedraggable'

  import { Plus, Search } from '@element-plus/icons-vue'

  import { requireRules } from '@/utils/rules'

  import type { FormInstance } from 'element-plus'

  import { filedDataType, siteCurrentType } from '@/types'

  import {
    siteSaveUpdateApi,
    getBySiteIdApi,
    getPortByHostApi,
    getApiInterfaceApiByIdApi,
    getApiListCountBySessionldApi,
  } from '~/src/api-ecs/site'
  import { getHostApiFieldsApi } from '@/api-ecs/public'
  import { ElMessage, ClickOutside as vClickOutside } from 'element-plus'

  import { RefreshRight } from '@element-plus/icons-vue'

  import { useUserStore } from '@/store/modules/user'

  import { resetSiteDisPlaysFiledApi } from '~/src/api-ecs/custom-field'
  import SaveModules from '@/ecs/config/init-configration/save-modules.vue'
  import dayjs from 'dayjs'

  import { uuid } from '~/src/utils'
  import { debounce, findIndex, cloneDeep, set } from 'lodash'
  const $baseConfirm: any = inject('$baseConfirm')

  const { getTableColumn } = useUserStore()

  const props = defineProps<{
    showEditSite: boolean
    currentItem: siteCurrentType
    mode: string
    title: string
    source: 'site' | 'assets'
  }>()

  const emit = defineEmits<{
    (e: 'on-closeEvent', val: boolean): void
    (e: 'on-reflash'): void
  }>()
  let defaultFields = [] as { id: number; fieldNameCn: string }[]
  let defaultFieldIds = [] as number[]
  const tip = `过滤条件请使用双引号 "" 进行包裹,例 : 源端口 = "8080"`
  const visible = ref(false)
  const hostLockVisible = ref(false)
  const sqlLockVisible = ref(false)
  const saveModulesVisible = ref(false)
  const formRef = ref<FormInstance>()

  const mode = ref('')

  const title = ref('')

  const portActiveIndex = ref(-1)

  const allHttpField = ref<filedDataType[]>([]) // http字段

  const searchSqlVisible = ref(false)

  const visableAddSqlItem = ref(false)

  const hostLock = ref(true)

  const sqlLock = ref(true)

  const HOSTCACHE = `${props.source}HOSTCACHE`

  const curHostStr = ref('')

  let _showApiInterfaces: { id: number; apiName: string; isDisplay: 0 | 1; dataCount: number }[] = []
  const resetDisable = ref(false)
  const isCheckedAll = ref(false)
  const isIndeterminate = ref(true)
  const apiInterfaceLoading = ref(false)
  const showApiInterface = ref(false)
  const checkedSiteInterface = ref<{ id: number; apiName: string; isDisplay: 0 | 1; dataCount: number }[]>([])
  const allApiInterface = ref<{ id: number; apiName: string; isDisplay: 0 | 1; dataCount: number }[]>([])
  const curHostS = ref('')
  // 表单数据
  const formData = reactive({
    siteName: '',
    phone: '',
    user: '',
    searchSql: '',
    module: '1',
    hosts: '',
    displayFields: [] as { id: number; fieldNameCn: string }[],
    id: undefined,
  })
  const comparisonData = ref<{
    type?: number
    siteSessionId?: number
    siteApiId?: number
    tag: 'alarm' | 'survey' | 'siteSession' | 'siteApi'
  }>()
  const replacementField = computed(() => formData.displayFields.map((i) => i.fieldNameCn))
  // 表单数据校验
  const rules = reactive({
    siteName: requireRules,
    // hosts: formData.module === '1' ? requireRules : [],
    displayFields: requireRules,
    phone: [{ pattern: /^1[3-9]\d{9}$/, message: '手机号不合法' }],
  })

  // 过滤条件下拉框
  const options = ref()

  let hostArr = ref<any[]>([])

  // 回显
  const getResult = async (sessionId: number | string) => {
    const {
      data: { displayFieldsArr, hosts, module, id, siteName, searchSql, user, phone, apiVos },
    } = await getBySiteIdApi({ id: sessionId })
    const allField = getTableColumn(30)
    hosts && changeHost(hosts)
    formData.module = module
    formData.id = id
    formData.siteName = siteName
    formData.searchSql = searchSql
    formData.user = user
    formData.phone = phone
    formData.displayFields = displayFieldsArr
      .map((key: number) => {
        return allField.find((i) => i.id === key)
      })
      .filter(Boolean)
    options.value = allField.filter(Boolean).filter((item: any) => item.fieldNameEn != 'position')
    allHttpField.value = allField.filter((item) => !displayFieldsArr.includes(item.id))
    checkedSiteInterface.value = apiVos || []
    initOverlay()
    emit('on-reflash')
  }

  // 得到所有字段
  const getHttpField = async () => {
    const ALL = getTableColumn(30)
    formData.displayFields = JSON.parse(JSON.stringify(defaultFields))
    allHttpField.value = ALL.filter((item) => !defaultFieldIds.includes(item.id))
    options.value = allHttpField.value.filter(Boolean).filter((item: any) => item.fieldNameEn != 'position')
  }

  // 校验host
  const verificationHost = () => {
    const arr = hostArr.value.length > 0 ? JSON.parse(JSON.stringify(hostArr.value)) : []
    const errArr: number[] = []
    let hostTrue = true
    arr.forEach((item: any, index: any) => {
      let flag = verification(item.host)
      if (!flag) {
        errArr.push(index + 1)
      }
    })
    if (errArr.length > 0) {
      hostTrue = false
      const str = errArr.join(',')
      ElMessage({ message: `Host第${str}条校验不通过，请修改`, type: 'error' })
    }
    return hostTrue
  }

  const ipv6Regex =
    /^(([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))$/
  const formatHost = () => {
    const obj: any = {}
    hostArr.value.forEach((item) => {
      // const hosts = ipv6Regex.test(item.host) ? `[${item.host}]` : item.host
      const port = item.selected.length == 0 ? [item.host] : item.port
      const selected = item.selected.length == 0 ? [item.host] : item.selected
      const result = {
        port,
        selected,
      }
      // obj[hosts] = result
      obj[item.host] = result
    })
    formData.hosts = JSON.stringify(obj)
  }

  // 格式化host
  const changeHost = (hosts: any) => {
    try {
      const data = JSON.parse(hosts)
      for (const item in data) {
        const obj = {
          host: item,
          port: data[item].port,
          selected: data[item].selected,
          show: false,
          id: uuid(),
          showBtn: true,
        }
        hostArr.value.push(obj)
      }
    } catch (error) {
      const arr = hosts?.split(',')
      hostArr.value = []
      arr.forEach((item: any) => {
        const has1 = item.indexOf('[')
        const has2 = item.indexOf(']')
        // const keysArr = Object.keys(hostObj)
        let isTrue = false
        if (has1 >= 0 && has2 >= 0) {
          const subStr = item.substring(has1 + 1, has2)
          isTrue = ipv6Regex.test(subStr)
        } else {
          isTrue = ipv6Regex.test(item)
        }
        if (has1 >= 0 && has2 >= 0) {
          const arr1 = item.split(']')
          const str = `${arr1[0]}]`
          const value = arr1[1].split(',').filter(Boolean)
          const obj = {
            host: str,
            port: value,
            selected: value,
            show: false,
            id: uuid(),
            showBtn: true,
          }
          hostArr.value.push(obj)
        } else if (has1 <= 0 && has2 <= 0 && isTrue) {
          const arr1 = item.split(':')
          const stringToRemove = `:${arr1[arr1.length - 1]}`
          const regex = new RegExp(stringToRemove, 'g')
          const newString = item.replace(regex, '')
          const value = arr1[arr1.length - 1].split(',').filter(Boolean)
          const obj = {
            host: newString,
            port: value,
            selected: value,
            show: false,
            id: uuid(),
            showBtn: true,
          }
          hostArr.value.push(obj)
        } else {
          const arr1 = item?.split(':')
          const str = arr1[0]
          const value = arr1[1]?.split(',').filter(Boolean) || []
          const obj = {
            host: str,
            port: value,
            selected: value,
            show: false,
            id: uuid(),
            showBtn: true,
          }
          hostArr.value.push(obj)
        }
      })
    }
  }

  // 提交
  const submitForm = async (formEl: FormInstance | undefined) => {
    if (!verificationHost()) return false
    if (!formEl) return
    await formEl.validate(async (valid, fields) => {
      if (valid) {
        formatHost()
        const arr: number[] = []
        formData.displayFields.forEach((item: any) => {
          arr.push(item.id)
        })
        const { msg } = await siteSaveUpdateApi({
          ...formData,
          id: props.mode == 'edit' ? props.currentItem.id : null,
          displayFields: arr.join(','),
          displayApiIds: checkedSiteInterface.value.map((i) => i.id),
        })
        ElMessage({ message: msg, type: 'success' })
        emit('on-reflash')
        handleClose()
      } else {
        console.log('error submit!', fields)
      }
    })
  }

  // 重置表单
  const resetForm = (formEl: FormInstance | undefined) => {
    if (!formEl) return
    formEl.resetFields()
  }

  // 重置字段
  const handleReset = async () => {
    const obj = {
      indexType: '',
      siteApiId: '',
      siteSessionId: +props.currentItem.id,
    }
    resetDisable.value = true
    try {
      const { msg, data } = await resetSiteDisPlaysFiledApi(obj)
      const allField = getTableColumn(30)
      formData.displayFields = data
        .map((key: number) => {
          return allField.find((i) => i.id === key)
        })
        .filter(Boolean)
      allHttpField.value = allField.filter((item) => !data.includes(item.id))
    } finally {
      resetDisable.value = false
    }
  }

  // 关闭
  const handleClose = () => {
    // localStorage.removeItem(HOSTCACHE)
    searchSqlVisible.value = false
    visableAddSqlItem.value = false
    sqlLock.value = true
    hostLock.value = true
    hostArr.value = []
    resetForm(formRef.value)
    emit('on-closeEvent', false)
  }
  const searchCondition = ref<any[]>([])
  const meta = [
    {
      field: '',
      fieldType: '',
      value: '',
      relation: 'and',
      condition: '',
      id: uuid(),
      isDisable: true,
      addRow: false,
      children: [
        {
          field: '',
          fieldType: '',
          value: '',
          relation: 'and',
          condition: '',
          id: uuid(),
          children: [],
          isDisable: true,
          addRow: false,
        },
      ],
    },
  ]
  watch(
    () => props.showEditSite,
    () => {
      visible.value = props.showEditSite
      if (visible.value) {
        loadNode(props.currentItem?.id)
      }
      mode.value = props.mode
      title.value = props.title
      hostTimeRange.value = '3天'
      if (props.showEditSite) {
        if (mode.value == 'edit') {
          const id = props.currentItem.id
          getResult(id)
        } else {
          hostArr.value = [
            {
              host: '',
              port: [],
              selected: [],
              show: false,
              id: uuid(),
              showBtn: false,
            },
          ]
          hostLock.value = false
          sqlLock.value = false
          searchCondition.value = meta
          getHttpField()
          initOverlay()
        }
      }
    }
  )
  const hostTimeRange = ref('3天')
  const hostTimeArr = [
    { value: '3天', label: '3天' },
    { value: '5天', label: '5天' },
    { value: '7天', label: '7天' },
    { value: '30天', label: '30天' },
  ]

  const portInputValue = ref(undefined)
  const handleAddPort = (index: number) => {
    // 判断是否是"*"号或者1-65535的数字，如果不是，则提示
    if (portInputValue.value && !/^\d+$/.test(portInputValue.value as string) && portInputValue.value != '*') {
      ElMessage({ message: '端口号只能是数字', type: 'error' })
      return
    } else if (
      portInputValue.value &&
      (Number(portInputValue.value) < 1 || Number(portInputValue.value) > 65535) &&
      portInputValue.value != '*'
    ) {
      ElMessage({ message: '端口号只能是1-65535的数字', type: 'error' })
      return
    }

    // 判断是否为*，如果是*，则清空port和selected
    if (portInputValue.value == '*') {
      // hostArr.value[index].port = []
      hostArr.value[index].selected = []
      hostArr.value[index].show = false
    }
    const res = portInputValue.value ? `${curHostStr.value}:${portInputValue.value}` : curHostStr.value
    hostArr.value[index].port.push(res)
    !hostArr.value[index].selected.includes(`${curHostStr.value}:*`) && hostArr.value[index].selected.push(res)
    hostArr.value[index].port = [...new Set(hostArr.value[index].port)]
    hostArr.value[index].selected = [...new Set(hostArr.value[index].selected)]
    hostArr.value[index].show = false
    hostArr.value[index].showBtn = false
    handlePortChange(index)
    setTimeout(() => {
      portInputValue.value = undefined
      hostArr.value[index].showBtn = true
    }, 10)
  }
  const handleCustomisation = (num: number) => {
    portInputValue.value = undefined
    hostArr.value[num].show = true
  }
  const handleDelItem = (id: string, num: number) => {
    $baseConfirm('确认删除吗？', null, async () => {
      const str = hostArr.value[num].host
      hostArr.value = hostArr.value.filter((item) => {
        return item.id != id
      })
      const cache = localStorage.getItem(HOSTCACHE) && JSON.parse(localStorage.getItem(HOSTCACHE) as string)
      delete cache[str]
      localStorage.setItem(HOSTCACHE, JSON.stringify(cache))
    })
  }
  const handleAddItem = () => {
    hostArr.value.push({
      host: '',
      port: [],
      selected: [],
      show: false,
      id: uuid(),
      showBtn: false,
    })
  }

  const queryForm = reactive({
    endTime: '',
    startTime: '',
    searchStr: '',
    topCount: 10000,
    indexType: 1,
  })

  const formatDayDate = () => {
    const timeDate = dayjs()
    const endDate = timeDate.format('YYYY-MM-DD HH:mm:ss')
    // const cur = timeDate.startOf('day')
    let startDate = null
    switch (hostTimeRange.value) {
      case '3天':
        startDate = timeDate.subtract(3, 'day').format('YYYY-MM-DD HH:mm:ss')
        break
      case '5天':
        startDate = timeDate.subtract(5, 'day').format('YYYY-MM-DD HH:mm:ss')
        break
      case '7天':
        startDate = timeDate.subtract(7, 'day').format('YYYY-MM-DD HH:mm:ss')
        break
      case '30天':
        startDate = timeDate.subtract(30, 'day').format('YYYY-MM-DD HH:mm:ss')
        break
      default:
        startDate = timeDate.startOf('day').format('YYYY-MM-DD HH:mm:ss')
        break
    }
    queryForm.endTime = endDate
    queryForm.startTime = startDate
  }
  const verification = (val: string) => {
    let str = val
    const has1 = val.indexOf('[')
    const has2 = val.indexOf(']')
    if (has1 >= 0 && has2 >= 0) {
      str = str.substring(has1 + 1, has2)
    }
    // IPV4
    const ipv4Regex =
      /^(?!0)(?:[0-9]|[1-9]\d|1\d{2}|2[0-4]\d|25[0-5])\.(?:[0-9]|[1-9]\d|1\d{2}|2[0-4]\d|25[0-5])\.(?:[0-9]|[1-9]\d|1\d{2}|2[0-4]\d|25[0-5])\.(?:[0-9]|[1-9]\d|1\d{2}|2[0-4]\d|25[0-5])$/
    // IPV6
    const ipv6Regex =
      /^(([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))$/
    // 域名正则表达式
    const domainRegex = /^[a-zA-Z0-9]+([\-\.]{1}[a-zA-Z0-9]+)*\.[a-zA-Z]{1,255}$/
    return ipv4Regex.test(str) || ipv6Regex.test(str) || domainRegex.test(str)
  }

  const hostloading = ref(false)
  const curdex = ref()
  const handleSearch = async (_index: number) => {
    const str = hostArr.value[_index].host
    const has = hostArr.value.some((_, index) => {
      return index != _index && hostArr.value[index].host === str
    })
    if (has) {
      ElMessage({ message: 'host数据重复', type: 'error' })
      hostArr.value[_index].host = ''
      return false
    }
    const flag = verification(str)
    if (!flag) return ElMessage({ message: '请输入正确IPV4，IPV6或域名格式', type: 'error' })
    formatDayDate()
    try {
      curdex.value = _index
      hostloading.value = true
      const { data } = await getPortByHostApi({ ...queryForm, searchStr: str })
      hostArr.value[_index].showBtn = true
      const cache = localStorage.getItem(HOSTCACHE) && JSON.parse(localStorage.getItem(HOSTCACHE) as string)
      hostArr.value[_index].port = data
      hostArr.value[_index].selected = data
      if (cache && cache[str]?.port?.length > 0) {
        hostArr.value[_index].port = [...hostArr.value[_index].port, ...cache[str].port]
        hostArr.value[_index].port = [...new Set(hostArr.value[_index].port)]
        hostArr.value[_index].selected = [...hostArr.value[_index].port, ...cache[str].selected]
        hostArr.value[_index].selected = [...new Set(hostArr.value[_index].selected)]
      }
    } finally {
      hostloading.value = false
    }
  }

  /**
   * @description '验证是否缓存过期'
   * @return boolean（true代表过期）
   */
  const verifyExpiration = (cache: any, str: string): boolean => {
    let res = true
    if (cache && str) {
      const now = new Date().getTime()
      res = now - cache[str].expirationTime > 24 * 60 * 60 * 1000 ? true : false
    }
    return res
  }

  const handlePortChange = (idx: number) => {
    const timestamp = new Date().getTime()
    const cache = localStorage.getItem(HOSTCACHE) && JSON.parse(localStorage.getItem(HOSTCACHE) as string)
    const obj: any = cache || {}
    hostArr.value.forEach((item, index) => {
      if (idx == index) {
        if (item.selected?.includes('*')) item.selected = ['*']
        if (item.host) {
          const result = {
            port: item.port,
            selected: item.selected,
            expirationTime: timestamp,
          }
          obj[item.host] = result
        }
      }
    })
    localStorage.setItem(HOSTCACHE, JSON.stringify(obj))
  }

  const handleHostChange = (index: number) => {
    hostArr.value[index].port = []
    hostArr.value[index].selected = []
    const cache = localStorage.getItem(HOSTCACHE) && JSON.parse(localStorage.getItem(HOSTCACHE) as string)
    const str = hostArr.value[index].host
    if (cache && cache[str]?.port?.length > 0) {
      hostArr.value[index].port = [...hostArr.value[index].port, ...cache[str].port]
      hostArr.value[index].port = [...new Set(hostArr.value[index].port)]
      hostArr.value[index].selected = [...hostArr.value[index].selected, ...cache[str].selected]
      hostArr.value[index].selected = [...new Set(hostArr.value[index].selected)]
    }
  }
  // #endregion

  // #region 过滤条件功能
  // ip
  const relationOptions = reactive<{ value: string; label: string }[]>([
    { value: '=', label: '= 等于' },
    { value: '!=', label: '!= 不等于' },
    { value: '>', label: '> 大于' },
    { value: '<', label: '< 小于' },
    { value: '>=', label: '>= 大于等于' },
    { value: '<=', label: '<= 小于等于' },
    { value: 'exists', label: 'exists 存在' },
    { value: 'not_exists', label: 'not_exists 不存在' },
  ])

  // 平均数
  const relationOptions1 = reactive<{ value: string; label: string }[]>([
    { value: '>', label: '> 大于' },
    { value: '<', label: '< 小于' },
  ])

  // 方向
  const relationOptions2 = reactive<{ value: string; label: string }[]>([{ value: '=', label: '= 等于' }])

  // isn_t
  const relationOptions3 = reactive<{ value: string; label: string }[]>([
    { value: '=', label: '= 等于' },
    { value: '!=', label: '!= 不等于' },
  ])

  // 字符串
  const relationOptions4 = reactive<{ value: string; label: string }[]>([
    { value: '=', label: '= 等于' },
    { value: '!=', label: '!= 不等于' },
    { value: 'like', label: 'like 字符串匹配' },
    { value: 'not_like', label: 'not_like 字符串匹配' },
    { value: 'exists', label: 'exists 存在' },
    { value: 'not_exists', label: 'not_exists 不存在' },
  ])
  const loadNode = async (id: number | string) => {
    if (!id) return
    apiInterfaceLoading.value = true
    try {
      const { data } = await getApiListCountBySessionldApi(id.toString())
      allApiInterface.value = Array.isArray(data) ? data : []
      apiInterfaceLoading.value = false
      isIndeterminate.value =
        checkedSiteInterface.value.length > 0 && checkedSiteInterface.value.length < allApiInterface.value.length
      isCheckedAll.value = checkedSiteInterface.value.length === allApiInterface.value.length
    } catch (error) {
      apiInterfaceLoading.value = false
      console.error(error)
    }
  }

  const handleFieldChange = (index: number, idx: number, val: any) => {
    let type = ''
    options.value.forEach((item: any) => {
      if (item.fieldNameEn == val) {
        type = item.fieldType
      }
    })
    searchCondition.value[index].children[idx].fieldType = type
    searchCondition.value[index].children[idx].isDisable = false
  }

  // 新增一组
  const handleAddSqlItem = (val: 'and' | 'or') => {
    searchCondition.value.push({
      field: '',
      fieldType: '',
      value: '',
      relation: val,
      condition: '',
      id: uuid(),
      isDisable: true,
      addRow: false,
      children: [
        {
          field: '',
          fieldType: '',
          value: '',
          relation: 'and',
          condition: '',
          id: uuid(),
          children: [],
          isDisable: true,
          addRow: false,
        },
      ],
    })
    visableAddSqlItem.value = false
  }

  // 增加一行
  const handleAddRow = (val: 'and' | 'or', index: number, idx: number) => {
    searchCondition.value[index].children.push({
      field: '',
      fieldType: '',
      value: '',
      relation: val,
      condition: '',
      id: uuid(),
      children: [],
      isDisable: true,
    })
    searchCondition.value[index].children[idx].addRow = false
  }

  // 删除
  const handleDelRow = (idx: number, id: string) => {
    $baseConfirm('确认删除吗？', null, async () => {
      if (searchCondition.value[idx].children.length > 1) {
        searchCondition.value[idx].children = searchCondition.value[idx].children.filter((item: any) => item.id != id)
      } else {
        searchCondition.value = searchCondition.value.filter((_, index: any) => {
          return index != idx
        })
      }
    })
  }
  const hanldeEditSiteInterface = () => {
    showApiInterface.value = true
  }

  const hanldeCheckedChange = debounce(() => {
    _showApiInterfaces = allApiInterface.value.filter((item) => item.isDisplay === 1)
    const showSize = _showApiInterfaces.length
    const allSzie = allApiInterface.value.length

    isCheckedAll.value = showSize === allSzie
    isIndeterminate.value = showSize > 0 && showSize < allSzie
  }, 300)
  const handleDisableApiCancel = () => {
    showApiInterface.value = false
    allApiInterface.value.forEach((api) => {
      const currIndex = findIndex(checkedSiteInterface.value, function (item) {
        return item.id === api.id
      })
      api.isDisplay = currIndex > -1 ? 1 : 0
    })
    isIndeterminate.value =
      checkedSiteInterface.value.length > 0 && checkedSiteInterface.value.length < allApiInterface.value.length
    isCheckedAll.value = checkedSiteInterface.value.length === allApiInterface.value.length
  }
  const handleCheckedAllChange = () => {
    allApiInterface.value.forEach((item) => {
      item.isDisplay = isCheckedAll.value ? 1 : 0
    })
    isIndeterminate.value = false
    _showApiInterfaces = isCheckedAll.value ? allApiInterface.value : []
  }
  const handleDisableApiOk = () => {
    checkedSiteInterface.value = _showApiInterfaces
    showApiInterface.value = false
  }
  // 取消
  const handleSqlCancel = () => {
    searchSqlVisible.value = false
    visableAddSqlItem.value = false
    searchCondition.value = meta
  }
  const handle2save = () => {
    const { id = undefined } = props.currentItem as any
    comparisonData.value = {
      type: undefined,
      siteSessionId: id,
      siteApiId: undefined,
      tag: 'siteSession',
    }
    saveModulesVisible.value = true
  }

  watch(
    () => hostLock.value,
    () => {
      if (hostLock.value) {
        initOverlay()
      }
    }
  )

  const initOverlay = () => {
    nextTick(() => {
      const target = document.querySelector('.x-host')
      const node = document.querySelector('.x-host-overlay')
      if (!target || !node) return
      const height = target.getBoundingClientRect().height
      node.setAttribute('style', `height:${height}px`)
    })
  }

  onMounted(async () => {
    const { data } = await getHostApiFieldsApi(30)
    defaultFields = (data.ids || []).map((id, index) => ({
      id: id,
      fieldNameCn: data.names[index],
    }))
    defaultFieldIds = data.ids
    initCache()
  })

  const initCache = () => {
    try {
      const cache = localStorage.getItem(HOSTCACHE) && JSON.parse(localStorage.getItem(HOSTCACHE) as string)
      if (cache) {
        Object.keys(cache).forEach((item: any) => {
          if (verifyExpiration(cache, item)) {
            delete cache[item]
            localStorage.setItem(HOSTCACHE, JSON.stringify(cache))
          }
        })
      }
    } catch (error) {
      console.log(error)
    }
  }

  const handleDelPort = (_idx: number, td: any) => {
    $baseConfirm('确认删除该端口吗？', null, async () => {
      hostArr.value[_idx].port = hostArr.value[_idx].port.filter((item: any) => {
        return item != td
      })
      hostArr.value[_idx].selected = hostArr.value[_idx].selected.filter((item: any) => {
        return item != td
      })
      handlePortChange(_idx)
    })
  }
</script>

<template>
  <div class="edit-site">
    <el-drawer v-model="visible" :before-close="handleClose" destroy-on-close size="85%" :title="title">
      <el-form ref="formRef" label-position="top" label-width="110px" :model="formData" :rules="rules">
        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="站点名称：" prop="siteName">
              <el-input v-model="formData.siteName" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="负责人：" prop="user">
              <el-input v-model="formData.user" clearable />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="联系电话：" prop="phone">
              <el-input v-model="formData.phone" clearable />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item v-if="mode == '1'" label="添加方式：" prop="module">
          <el-radio-group v-model="formData.module" class="ml-4">
            <el-radio label="1">精简模式</el-radio>
            <!-- <el-radio label="0">高级模式</el-radio> -->
          </el-radio-group>
        </el-form-item>
        <el-form-item>
          <template #label>
            API接口展示字段：
            <el-icon style="vertical-align: -2px" @click="hanldeEditSiteInterface"><Edit /></el-icon>
          </template>
          <div class="space" style="cursor: pointer; min-height: 90px" @click="hanldeEditSiteInterface">
            <vue-draggable item-key="id" :list="checkedSiteInterface" :move="undefined" :sort="false">
              <template #item="{ element }">
                <el-tag>{{ element.apiName }}</el-tag>
              </template>
            </vue-draggable>
          </div>
        </el-form-item>
        <el-form-item class="is-required" label="HOST：" prop="hosts">
          <template #label>
            <span>HOST：</span>
            <el-select
              v-model="hostTimeRange"
              :disabled="hostLock"
              size="small"
              style="width: 63px; margin: 0 0px 2px 5px"
            >
              <el-option v-for="item in hostTimeArr" :key="item.value" :label="item.label" :value="item.value" />
            </el-select>
            <span style="color: #a9abb1">
              （自定义：支持输入端口范围（1-65535）、空值或 " * " , " * "
              模式下支持匹配指定主机名（HOST）其后续附加的任意有效端口号）
            </span>
            <el-tooltip
              class="box-item"
              :content="hostLock ? '点击解锁进行编辑' : '锁定，退出编辑'"
              effect="dark"
              placement="left"
              :visible="hostLockVisible"
            >
              <el-image
                v-if="!hostLock"
                :src="require('@/assets/site_images/kaisuo.svg')"
                style="height: 16px; width: 16px; float: right; margin: 3px 4px 0 0; cursor: pointer"
                @click="hostLock = !hostLock"
                @mouseenter="hostLockVisible = true"
                @mouseleave="hostLockVisible = false"
              />

              <el-image
                v-else
                :src="require('@/assets/site_images/jiesuo.svg')"
                style="height: 16px; width: 16px; float: right; margin: 3px 4px 0 0; cursor: pointer"
                @click="hostLock = !hostLock"
                @mouseenter="hostLockVisible = true"
                @mouseleave="hostLockVisible = false"
              />
            </el-tooltip>
          </template>
          <el-row
            class="x-host"
            style="width: 100%"
            @mouseenter="
              () => {
                if (hostLock) hostLockVisible = true
              }
            "
            @mouseleave="
              () => {
                hostLockVisible = false
              }
            "
          >
            <template v-for="(item, _index) in hostArr" :key="item.id">
              <el-col class="x-col" :span="6">
                <div class="host-item" style="display: flex; align-items: center; justify-content: center">
                  <el-input
                    v-model="item.host"
                    class="host-input"
                    :disabled="hostloading"
                    placeholder="请输入IP或域名"
                    @input="handleHostChange(_index)"
                  >
                    <!--  -->
                    <template #prepend>
                      <span>{{ _index + 1 }}</span>
                    </template>
                    <template #append>
                      <el-button :disabled="hostLock" @click="handleSearch(_index)">
                        <vab-icon icon="search-line" />
                      </el-button>
                    </template>
                  </el-input>
                </div>
              </el-col>
              <el-col v-loading="hostloading && _index == curdex" class="x-col" :span="18" style="padding: 0px 20px">
                <div class="host-item">
                  <el-checkbox-group v-model="item.selected" style="margin-right: 20px">
                    <el-checkbox
                      v-for="(td, index) in item.port"
                      :key="index"
                      :disabled="item.selected?.includes('*') && td != '*'"
                      :label="td"
                      :value="td"
                      @change="handlePortChange(_index)"
                      @mouseenter="portActiveIndex = index"
                      @mouseleave="portActiveIndex = -1"
                    >
                      <template #default>
                        <div style="position: relative">
                          <span>{{ td }}</span>
                          <el-icon
                            v-if="portActiveIndex == index"
                            style="position: absolute; right: -9px; top: -9px; color: #9b9b9b; font-size: 12px"
                            @click.stop.prevent="handleDelPort(_index, td)"
                          >
                            <Close />
                          </el-icon>
                        </div>
                      </template>
                    </el-checkbox>
                    <el-popover v-if="item.showBtn" placement="left" trigger="click" :width="400">
                      <template #reference>
                        <el-button
                          :disabled="item.selected?.includes('*')"
                          :icon="Plus"
                          plain
                          size="small"
                          :style="{
                            padding: '10px 5px',
                            marginLeft: item.port.length > 0 ? '20px' : '0',
                            verticalAlign: 'baseline',
                          }"
                          @click="curHostStr = hostArr[_index].host"
                        >
                          自定义
                        </el-button>
                      </template>
                      <template #default>
                        <el-input v-model="curHostStr" disabled style="width: 195px" />
                        :
                        <el-input v-model="portInputValue" placeholder="请输入端口" style="width: 90px" />
                        <el-button style="margin: 0 0 3px 11px" type="primary" @click="handleAddPort(_index)">
                          确认
                        </el-button>
                      </template>
                    </el-popover>
                  </el-checkbox-group>
                  <el-icon
                    v-if="hostArr.length > 1"
                    class="x-icon"
                    style="cursor: pointer"
                    @click="handleDelItem(item.id, _index)"
                  >
                    <Delete />
                  </el-icon>
                </div>
              </el-col>
            </template>
            <el-col :span="24">
              <div class="add-item">
                <el-button :icon="Plus" link type="primary" @click="handleAddItem">添加</el-button>
              </div>
            </el-col>
            <div v-if="hostLock" class="x-host-overlay"></div>
          </el-row>
        </el-form-item>
        <el-form-item prop="searchSql">
          <template #label>
            <span>过滤条件：</span>
            <el-tooltip
              class="box-item"
              :content="sqlLock ? '点击解锁进行编辑' : '锁定，退出编辑'"
              effect="dark"
              placement="left"
              :visible="sqlLockVisible"
            >
              <el-image
                v-if="!sqlLock"
                :src="require('@/assets/site_images/kaisuo.svg')"
                style="height: 16px; width: 16px; float: right; margin: 3px 4px 0 0; cursor: pointer"
                @click="sqlLock = !sqlLock"
                @mouseenter="sqlLockVisible = true"
                @mouseleave="sqlLockVisible = false"
              />
              <el-image
                v-else
                :src="require('@/assets/site_images/jiesuo.svg')"
                style="height: 16px; width: 16px; float: right; margin: 3px 4px 0 0; cursor: pointer"
                @click="sqlLock = !sqlLock"
                @mouseenter="sqlLockVisible = true"
                @mouseleave="sqlLockVisible = false"
              />
            </el-tooltip>

            <el-popover placement="left" trigger="click" :visible="searchSqlVisible" :width="860">
              <template #reference>
                <el-button
                  :icon="Plus"
                  size="small"
                  style="margin: 0 0 2px 5px; display: none"
                  @click="searchSqlVisible = true"
                >
                  添加
                </el-button>
              </template>
              <template #default>
                <div class="top out">
                  <span style="font-weight: 500; font-size: 16px; color: #303133; line-height: 25px; text-align: left">
                    添加过滤条件
                  </span>
                </div>
                <el-divider style="margin: 8px 0 0 0" />
                <div class="content">
                  <div
                    v-for="(item, index) in searchCondition"
                    :key="item.id"
                    style="position: relative; padding: 8px 0; margin-top: 10px; background: #f7f6ff"
                  >
                    <span style="position: absolute; color: #6954f0; left: -25px; top: 15px">
                      {{ item.relation == 'and' ? '且' : '或' }}
                    </span>
                    <div
                      v-for="(td, idx) in item.children"
                      :key="td.id"
                      :style="{ width: '100%', display: 'flex', marginTop: idx == 0 ? '0' : '10px' }"
                    >
                      <div
                        v-if="idx == 0"
                        style="color: #1e1841; font-size: 14px; width: 90px; line-height: 32px; text-align: center"
                      >
                        分组{{ index + 1 }}
                      </div>
                      <div
                        v-else
                        style="color: #6954f0; font-size: 14px; width: 90px; text-align: center; line-height: 32px"
                      >
                        {{ td.relation == 'and' ? '且' : '或' }}
                      </div>
                      <el-row :gutter="8" style="width: 100%">
                        <el-col :span="5">
                          <el-select
                            v-model="td.field"
                            style="width: 100%"
                            @change="handleFieldChange(index, idx, $event)"
                          >
                            <el-option
                              v-for="item in options"
                              :key="item.fieldNameEn"
                              :label="item.fieldNameCn"
                              :value="item.fieldNameEn"
                            />
                          </el-select>
                        </el-col>
                        <el-col :span="5">
                          <el-select
                            v-if="td.fieldType == 'ip' || td.fieldType == 'num'"
                            v-model="td.condition"
                            :disabled="td.isDisable"
                            style="width: 100%"
                          >
                            <el-option
                              v-for="item in relationOptions"
                              :key="item.value"
                              :label="item.label"
                              :value="item.value"
                            />
                          </el-select>
                          <el-select
                            v-else-if="td.fieldType == 'average'"
                            v-model="td.condition"
                            :disabled="td.isDisable"
                            style="width: 100%"
                          >
                            <el-option
                              v-for="item in relationOptions1"
                              :key="item.value"
                              :label="item.label"
                              :value="item.value"
                            />
                          </el-select>
                          <el-select
                            v-else-if="td.fieldType == 'direction'"
                            v-model="td.condition"
                            :disabled="td.isDisable"
                            style="width: 100%"
                          >
                            <el-option
                              v-for="item in relationOptions2"
                              :key="item.value"
                              :label="item.label"
                              :value="item.value"
                            />
                          </el-select>
                          <el-select
                            v-else-if="td.fieldType == 'isn_t'"
                            v-model="td.condition"
                            :disabled="td.isDisable"
                            style="width: 100%"
                          >
                            <el-option
                              v-for="item in relationOptions3"
                              :key="item.value"
                              :label="item.label"
                              :value="item.value"
                            />
                          </el-select>
                          <el-select
                            v-else-if="td.fieldType == 'text'"
                            v-model="td.condition"
                            :disabled="td.isDisable"
                            style="width: 100%"
                          >
                            <el-option
                              v-for="item in relationOptions4"
                              :key="item.value"
                              :label="item.label"
                              :value="item.value"
                            />
                          </el-select>
                          <el-select v-else v-model="td.condition" :disabled="td.isDisable" style="width: 100%">
                            <el-option
                              v-for="item in relationOptions4"
                              :key="item.value"
                              :label="item.label"
                              :value="item.value"
                            />
                          </el-select>
                        </el-col>
                        <el-col :span="12">
                          <el-input v-model="td.value" :disabled="td.isDisable" />
                        </el-col>
                        <el-col :span="1">
                          <el-icon
                            v-if="searchCondition.length != 1 || item.children.length != 1"
                            color="#6954f0"
                            style="margin-top: 10px; cursor: pointer"
                            @click="handleDelRow(index, td.id)"
                          >
                            <Delete />
                          </el-icon>
                        </el-col>
                        <el-col :span="1">
                          <el-popover placement="top" trigger="click" :visible="td.addRow" :width="120">
                            <template #reference>
                              <el-icon
                                v-if="idx == item.children.length - 1"
                                color="#6954f0"
                                style="margin-top: 10px; cursor: pointer"
                                @click="td.addRow = true"
                              >
                                <Plus />
                              </el-icon>
                            </template>
                            <template #default>
                              <el-button link type="primary" @click="handleAddRow('and', index, idx)">添加且</el-button>
                              <span style="margin: 0 4px">or</span>
                              <el-button link type="primary" @click="handleAddRow('or', index, idx)">添加或</el-button>
                            </template>
                          </el-popover>
                        </el-col>
                      </el-row>
                    </div>
                  </div>
                </div>
                <div style="height: 32px; margin-left: 61px">
                  <el-button link type="primary">
                    <el-popover placement="top" trigger="click" :visible="visableAddSqlItem" :width="120">
                      <template #reference><span @click="visableAddSqlItem = true">添加分组</span></template>
                      <template #default>
                        <el-button link type="primary" @click="handleAddSqlItem('and')">添加且</el-button>
                        <span style="margin: 0 4px">or</span>
                        <el-button link type="primary" @click="handleAddSqlItem('or')">添加或</el-button>
                      </template>
                    </el-popover>
                  </el-button>
                </div>
                <el-divider style="margin: 8px 0 15px 0" />
                <div class="footer out">
                  <div style="float: right">
                    <el-button style="height: 28px" type="primary">确认</el-button>
                    <el-button style="height: 28px" @click="handleSqlCancel">取消</el-button>
                  </div>
                </div>
              </template>
            </el-popover>
          </template>
          <el-input
            v-model="formData.searchSql"
            :disabled="sqlLock"
            :placeholder="tip"
            resize="none"
            :rows="4"
            type="textarea"
            @mouseenter="
              () => {
                if (sqlLock) sqlLockVisible = true
              }
            "
            @mouseleave="
              () => {
                sqlLockVisible = false
              }
            "
          />
        </el-form-item>
        <el-form-item prop="displayFields">
          <template #label>
            <span>展示字段：</span>
            <span style="color: #a9acb3; font-size: 13px; margin-left: 5px">(拖动排序)&nbsp;</span>
          </template>
          <div class="space">
            <vue-draggable
              animation="300"
              chosen-class="chosenClass"
              ghost-class="ghostClass"
              group="my-group"
              item-key="id"
              :list="formData.displayFields"
            >
              <template #item="{ element }">
                <el-tag>{{ element.fieldNameCn }}</el-tag>
              </template>
            </vue-draggable>
          </div>
        </el-form-item>
        <el-form-item label="可选字段：">
          <div class="space">
            <vue-draggable
              animation="300"
              chosen-class="chosenClass"
              ghost-class="ghostClass"
              group="my-group"
              item-key="id"
              :list="allHttpField"
              :sort="false"
            >
              <template #item="{ element }">
                <el-tag type="info">{{ element.fieldNameCn }}</el-tag>
              </template>
            </vue-draggable>
          </div>
        </el-form-item>
      </el-form>
      <div class="footer" style="z-index: 9999">
        <div>
          <el-button
            v-if="mode == 'edit'"
            :disabled="resetDisable"
            :icon="RefreshRight"
            link
            :loading="resetDisable"
            style="float: left; color: #6954f0; margin-left: 0px"
            @click="handleReset"
          >
            重置字段
          </el-button>
        </div>
        <div style="float: right">
          <span v-if="mode === 'edit'" class="save" :disabled="resetDisable" @click="handle2save">另存为</span>
          <el-button :disabled="resetDisable" type="primary" @click="submitForm(formRef)">确认</el-button>
          <el-button @click="handleClose">取消</el-button>
        </div>
      </div>
      <save-modules
        v-model="saveModulesVisible"
        :comparison-data="comparisonData"
        :replacement-field="replacementField"
      />
    </el-drawer>
    <el-dialog
      v-model="showApiInterface"
      class="dialogApiInterface"
      destroy-on-close
      title="API接口管理"
      width="1100px"
    >
      <el-table v-loading="apiInterfaceLoading" :data="allApiInterface" style="width: 1060px">
        <el-table-column label="接口名称" prop="apiName" show-overflow-tooltip />
        <el-table-column label="接口地址" prop="apiUrl" show-overflow-tooltip />
        <el-table-column label="数量" prop="dataCount" />
        <el-table-column label="接口是否展示" width="110">
          <template #default="{ row }">
            <el-switch
              v-model="row.isDisplay"
              active-text="是"
              :active-value="1"
              inactive-text="否"
              :inactive-value="0"
              inline-prompt
              style="--el-switch-off-color: #ccc"
              @change="hanldeCheckedChange"
            />
          </template>
        </el-table-column>
      </el-table>
      <footer>
        <span>温馨提示：选中保存后，即为展示该API接口</span>
        <div>
          <el-checkbox
            v-model="isCheckedAll"
            :indeterminate="isIndeterminate"
            label="是否全选"
            @change="handleCheckedAllChange"
          />
          <el-button type="primary" @click="handleDisableApiOk">确认</el-button>
          <el-button @click="handleDisableApiCancel">取消</el-button>
        </div>
      </footer>
    </el-dialog>
  </div>
</template>

<style scoped lang="scss">
  .content {
    margin-left: 40px;
    padding: 10px 0;
  }
  .x-host {
    .host-item {
      width: 100%;
      height: 100%;
      min-height: 36px;
      display: flex;
      align-items: center;
      position: relative;
      :deep() {
        .el-input {
          position: absolute;
          inset: 0;
          .el-input__inner {
            line-height: 100%;
          }
        }
      }
      .x-icon {
        position: absolute;
        right: 0px;
        top: 50%;
        transform: translateY(-50%);
      }
    }

    :deep(.el-col-6) {
      border: 1px solid var(--el-border-color);
    }
    :deep(.el-col-18) {
      border: 1px solid var(--el-border-color);
    }
    :deep(.x-col) {
      border-right: none;
      border-bottom: none;
    }

    .add-item {
      width: 100%;
      height: 37px;
      border: 1px solid var(--el-border-color);
      // border-top: none;
      display: flex;
      align-items: center;
      justify-content: center;
      button {
        color: var(--el-color-primary);
        --el-button-active-color: var(--el-color-primary);
      }
    }
  }
  :deep(.x-host) {
    :nth-child(2n).el-col {
      border-right: 1px solid var(--el-border-color);
    }
    .x-host-overlay {
      position: absolute;
      width: 100%;
      // min-height: 74px;
      cursor: not-allowed;
      z-index: 1000;
      opacity: 0.15;
      border: 1px solid #d7dbe2;
      background-color: #c4cfe9;
    }
  }

  .space {
    width: 100%;
    min-height: 150px;
    padding: 10px;
    border: 1px solid var(--el-border-color);
    border-radius: 2px;
    div {
      height: calc(100% - 20px);
      min-height: 130px;
    }
  }

  :deep(.el-tag) {
    margin-left: 10px;
  }
  .edit-site {
    .api-management {
      width: 104px;
      height: 32px;
      background: #ffffff;
      border-radius: 4px;
      display: inline-block;
      line-height: 32px;
      font-weight: 400;
      font-size: 14px;
      color: #6954f0;
      text-align: center;
      margin-right: 16px;
      cursor: pointer;
    }
    .api-management-list {
      max-width: 350px;
      min-width: 260px;
      background: #ffffff;
      box-shadow: 0px 0px 4px 4px rgba(77, 81, 86, 0.05);
      border-radius: 6px 0px 0px 6px;
      border: 1px solid #ecebfa;
      position: fixed;
      top: 58px;
      bottom: 0;
      right: 0;
      z-index: 9999;
      transition: all 0.3s cubic-bezier(0.165, 0.84, 0.44, 1);
      overflow: hidden;
      &.hidden {
        bottom: 101%;
      }
      .api-list {
        // margin-left: 20px;
        height: calc(100% - 60px);
        overflow-y: auto;
        overflow-x: hidden;
        .all {
          display: flex;
          position: sticky;
          top: 0;
          background-color: #fff;
          z-index: 9;
        }
        :deep() {
          .el-checkbox {
            padding-left: 20px;
          }
          .el-checkbox-group {
            display: flex;
            flex-direction: column;
            .el-checkbox {
              margin-right: 0;
              &:hover {
                background-color: #efeefe;
              }
            }
          }
        }
      }
      h4 {
        font-weight: 500;
        font-size: 16px;
        color: #303133;
        line-height: 50px;
        border-bottom: 1px solid #ecebfa;
        margin-block: 0;
        text-indent: 20px;
      }
    }
    :deep() {
      .el-drawer {
        border-radius: 30px 0px 0px 30px;
        background: #534b89;
        .el-drawer__header {
          padding: 12px 30px 12px;
          margin-bottom: 0;
          .el-drawer__title {
            font-size: 20px;
            font-weight: 500;
            font-size: 20px;
            line-height: 34px;
            color: #fff;
          }
          .el-drawer__close-btn {
            width: 26px;
            height: 26px;
            border-radius: 50%;
            background: #eeedf9;
            display: flex;
            justify-content: center;
            align-items: center;
          }
        }
        .el-drawer__body {
          border-radius: 30px 0px 0px 30px;
          background: #fff;
          padding: 20px 20px 0;
          .el-form.el-form--default {
            min-height: calc(100% - 72px);
          }
          .host-input {
            .el-input__wrapper:not(.is-focus) {
              box-shadow: none;
            }
          }
          .el-input-group__prepend {
            background: #fff;
            width: 36px;
            height: 36px;
            border-right: 1px solid var(--el-border-color);
            box-shadow: none;
          }
          .el-input-group__append {
            background: var(--el-color-primary);
            color: #fff;
            // width: 10px;
            height: 36px !important;
            width: 30px !important;
            --el-input-border-color: var(--el-color-primary);
            button {
              height: 32px;
              width: 30px;
              display: flex;
              align-items: center;
              justify-content: center;
            }
            i {
              font-size: 14px;
              margin-top: -1px;
            }
          }
        }
        .footer {
          display: flex;
          justify-content: space-between;
          position: sticky;
          bottom: 0;
          background: #fff;
          padding-block: 20px;
        }
      }
    }
  }
  :deep {
    .save {
      cursor: pointer;
      margin-right: 10px;
      color: var(--el-color-primary);
      font-weight: 400;
      font-size: 14px;
      vertical-align: middle;
    }
    .el-dialog.dialogApiInterface {
      height: 800px;
      background: #ffffff;
      box-shadow: 0px 0px 4px 4px rgba(77, 81, 86, 0.05);
      border-radius: 10px;
      .el-table {
        .el-table__header,
        .el-table__body {
          .el-table__cell:last-child {
            .cell {
              margin-left: -10px;
            }
          }
          .el-table__cell:first-child {
            .cell {
              margin-left: 10px;
            }
          }
        }
        .el-scrollbar__wrap,
        .el-scrollbar {
          max-height: 590px !important;
          min-height: 300px !important;
        }
      }
      footer {
        display: flex;
        justify-content: space-between;
        align-items: center;
        position: absolute;
        inset: 30px;
        top: auto;
        & > span {
          font-weight: 400;
          font-size: 14px;
          color: #918da5;
        }
        .el-checkbox {
          vertical-align: bottom;
          margin-right: 30px;
        }
      }
    }
  }
</style>
