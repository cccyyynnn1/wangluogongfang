<script lang="ts">
  export default {
    name: 'D3ForceGraph',
  }
</script>

<script setup lang="ts">
  import * as d3 from 'd3'
  import { link, node } from '~/src/types'
  import VueEvent from '@/data/event'
  import _ from 'lodash'
  import { QueryAssetByIpApi } from '~/src/api-ecs/chase'

  const props = defineProps<{
    graphData: {
      nodes: any[]
      links: any[]
    }
    startSelf: boolean
    isPath?: 'pathTracing' | 'threatHunting'
    disableKeys?: string[]
  }>()
  const emits = defineEmits<{
    (e: 'show-relate', data: any): void
    (e: 'graph-move'): void
    (e: 'show-info', data: node): void
    (e: 'show-filter', data: node): void
    (e: 'remark-node-color', data: node): void
    (e: 'close-filter'): void
    (e: 'expand-all', data: string): void
    (e: 'expand-one', data: string): void
    (e: 'del-node', data: string): void
  }>()
  /**
   * @description 'node 的大小'
   */
  const symbolSize = 48
  let my_svg: any = null
  /**
   * @description '悬停提示框'
   */
  let my_tooltip: any = null
  /**
   * @description '节点'
   */
  let d3_node: any = null
  /**
   * @description '连线'
   */
  let d3_link: any = null

  let simulation: any = null
  /**
   * @description '右键菜单'
   */
  let sketchpad_contect_menu: any = null
  /**
   * @description 'd3 锚点'
   */
  const transformData = ref({
    k: 1,
    x: 0,
    y: 0,
  })
  /**
   * @description '移入悬停的计时器（防抖）'
   */
  let timeoutEnter: any = null
  /**
   * @description '是否允许关闭右键菜单'
   */
  let contextmenuCloseFlag = false
  /**
   * @description '是否允许鼠标移入事件代码'
   */
  let timeoutFlag = true
  /**
   * @description '用于初次渲染悬停框的bug'
   */
  let timeoutFlag1 = true
  /**
   * @description '打开右键菜单，悬停事件失效'
   */
  let contextmenuFlag = false
  /**
   * @description '移出悬停的计时器（防抖）'
   */
  let timeoutLevae: any = null
  /**
   * @description '离开悬停组件时组织移入的代码继续运行'
   */
  let leaveFlag = false
  // /**
  //  * @description '是否在拖拽中'
  //  */
  // let isDrag = false
  const curNodeWhenDrag = ref<node>()
  const curNodeWhenDragLeft = ref('100px')
  const curNodeWhenDragTop = ref('100px')
  const rectData = ref()
  const d3_force_ref = ref()
  const d3_instance = ref()
  const d3_links = ref<any[]>([])
  const d3_nodes = ref<any[]>([])
  const svg_zoom = d3
    .zoom()
    .scaleExtent([0.5, 1.5])
    .on('zoom', function (this: any, event: any) {
      transformData.value = event.transform
      emits('graph-move')
      d3.select('.group').attr('transform', d3.zoomTransform(my_svg.node()))
    })
    .filter((event: any) => {
      const filterKeys = props.disableKeys || []
      return !filterKeys.includes(event.type)
    })
  function renderTooltip() {
    d3.select('.chase-tooltip').remove()
    my_tooltip = d3
      .select(my_svg.node().parentElement.parentElement)
      .append('div')
      .attr('class', 'chase-tooltip')
      .style('opacity', 0)
  }
  const renderContextMenu = () => {
    d3.select('.d3-context-menu').remove()
    sketchpad_contect_menu = d3
      .select(my_svg.node().parentElement.parentElement)
      .append('div')
      .attr('class', 'd3-context-menu')
      .style('display', 'none')
      .on('mouseenter', () => (contextmenuCloseFlag = true))
      .on('mouseleave', () => (contextmenuCloseFlag = false))
    d3.select('body')
      .on('contextmenu', (event: MouseEvent) => {
        event.stopPropagation()
        event.preventDefault()
      })
      .on('click.d3-context-menu', function (e: any) {
        if (contextmenuCloseFlag) return
        contextmenuFlag = false
        d3.select('.d3-context-menu').style('display', 'none')
      })
  }

  const d3_drag = (simulation: any) => {
    let activeNodes: any[] = []

    function dragstarted(this: any, event: any, d: node) {
      leaveFlag = true
      my_tooltip.style('opacity', 0).style('display', 'none')
      d3.select('.d3-context-menu').style('display', 'none')
      clearGraphStyle()
      // 自定义连线
      if (props.startSelf) {
        simulation.stop()
        curNodeWhenDrag.value = d
        curNodeWhenDragLeft.value = `${event.sourceEvent.x - 118}px`
        curNodeWhenDragTop.value = `${event.sourceEvent.y - 85}px`
      } else {
        if (!event.active) simulation.alphaTarget(0).restart()
        // 记录被按住的节点
        activeNodes.push(event.subject)
        event.subject.fx = event.subject.x
        event.subject.fy = event.subject.y
        event.subject.vx = 0
        event.subject.vy = 0
      }
    }

    function dragged(event: any, d: any) {
      if (!props.startSelf) {
        simulation.alphaTarget(0.3).restart()
        leaveFlag = true
        event.subject.fx = event.x
        event.subject.fy = event.y
        curNodeWhenDragLeft.value = `${event.sourceEvent.x - 118}px`
        curNodeWhenDragTop.value = `${event.sourceEvent.y - 85}px`
        activeNodes.forEach((node) => {
          node.fx = event.x
          node.fy = event.y
        })
      } else {
        leaveFlag = true
        event.subject.fx = event.x
        event.subject.fy = event.y
        curNodeWhenDragLeft.value = `${event.sourceEvent.x - 118}px`
        curNodeWhenDragTop.value = `${event.sourceEvent.y - 85}px`
      }
    }

    function dragended(this: any, event: any, d: any) {
      if (!event.active) simulation.alphaTarget(-0.1)
      // 按下时暂停模拟
      // setTimeout(() => {
      //   simulation.stop()
      // }, 100)
      // 移除已释放的节点
      activeNodes = activeNodes.filter((node) => node !== event.subject)
      // 注释以下代码，使拖动结束后固定节点
      event.subject.fx = null
      event.subject.fy = null
      curNodeWhenDrag.value = undefined
      if (props.startSelf) {
        const data = {
          isDrag: !!curNodeWhenDrag.value,
          data: curNodeWhenDrag.value,
        }
        VueEvent.emit('insideIsOutsideResult', data)
      }
    }

    const darg = d3.drag().on('start', dragstarted).on('drag', dragged).on('end', dragended)

    // d3.select('svg').selectAll('.node').call(darg)

    return darg
  }
  function changeGraphStyle(selectData: string | any) {
    const _val =
      typeof selectData === 'string'
        ? { clientIp: selectData, serverIp: selectData }
        : { clientIp: selectData.clientIp, serverIp: selectData.serverIp }

    // const list =
    //   typeof selectData === 'string'
    //     ? d3_links.value
    //         .filter((link) => link.clientIp === _val.clientIp || link.serverIp === _val.serverIp)
    //         .map((i) => [i.clientIp, i.serverIp])
    //         .flat()
    //     : [selectData.clientIp, selectData.serverIp]

    // d3_instance.value.select('.nodes').selectAll('.node')
    // .attr('class', (d: any) => {
    //   if (d.id === _val.clientIp) return typeof selectData === 'string' ? 'fixed' : 'active'
    //   if (list.includes(d.id)) return 'active'
    //   return 'node'
    // })
    // 处理相邻的边line
    d3_instance.value.select('.links').selectAll('line')
    // .attr('class', (d: any) => {
    //   if (typeof selectData === 'string')
    //     return d.clientIp === _val.clientIp || d.serverIp === _val.clientIp ? 'active' : ''
    //   return d.clientIp === _val.clientIp && d.serverIp === _val.serverIp ? 'active' : ''
    // })
  }
  function clearGraphStyle() {
    // 移除所有样式
    // my_svg.select('.nodes').selectAll('circle').attr('class', '')
    my_svg.select('.texts').selectAll('text').attr('class', '')
    my_svg.select('.links').selectAll('line').attr('class', '').attr('marker-end', 'url(#posMarker)')
    my_svg.select('.linkTexts').selectAll('text').attr('class', '')
  }
  function textBreaking(d3text: any, text: string) {
    const texts = text.split('\n')
    const topY = -16
    const midY = 0
    const botY = 16
    texts.forEach((text, index) => {
      d3text
        .append('tspan')
        .attr('x', -1)
        .attr('y', 40 + index * 16)
        .text(text)
    })
  }
  // 初始化d3图
  const curGraphData = ref<{ nodes: node[]; links: link[] }>({ nodes: [], links: [] })
  function d3init() {
    my_svg && my_svg.remove()
    const { links, nodes } = curGraphData.value
    if (nodes.length == 0) return
    getRectData()
    d3_links.value = JSON.parse(JSON.stringify(links))
    d3_nodes.value = JSON.parse(JSON.stringify(nodes))
    d3_instance.value = d3.select(d3_force_ref.value)
    const svg = d3_instance.value.style('font', '12px sans-serif').call(svg_zoom).append('g').attr('class', 'group')
    // 引力关系
    my_svg = svg
    simulation = d3
      .forceSimulation(d3_nodes.value)
      .force('charge', d3.forceManyBody().strength(-50))
      .force('collision', d3.forceCollide().radius(75).iterations(1).strength(0.5))
      .force(
        'center',
        d3.forceCenter(svg.node().parentElement.clientWidth / 2, svg.node().parentElement.clientHeight / 2)
      )
      .force(
        'link',
        d3
          .forceLink(d3_links.value)
          .id((d: any) => d.id)
          .distance(100)
      )
    renderMarkers(svg)
    renderTooltip()
    renderContextMenu()
    // 线段
    d3_link = svg.append('g').attr('class', 'links').selectAll('.link').data(d3_links.value)
    creatLink(d3_link)
    // 节点
    d3_node = svg.selectAll('.node').data(d3_nodes.value)
    creatNode(d3_node)
    // 引力图的线段
    simulation.on('tick', () => {
      const allNode = my_svg.selectAll('.node')

      const allLisks = my_svg.selectAll('.links').selectAll('.link')
      allLisks.attr('d', function (d: any) {
        // 计算从 source 到 target 的方向向量
        const dx = d.target.x - d.source.x
        const dy = d.target.y - d.source.y
        // 计算方向向量的长度 (即两点之间的距离)
        const distance = Math.sqrt(dx * dx + dy * dy)
        // 计算单位方向向量
        const unitX = dx / distance
        const unitY = dy / distance
        // 计算新的终点坐标，使链接在 target 端短一些
        const offset = 30 // 偏移量
        const newTargetX = d.target.x - unitX * offset
        const newTargetY = d.target.y - unitY * offset
        const newSourceX = d.source.x + unitX * offset
        const newSourceY = d.source.y + unitY * offset
        return `M ${newSourceX} ${newSourceY} L ${newTargetX} ${newTargetY}`
      })
      // 箭头
      // .attr('marker-end', function (d: any) {
      //   return 'url(#posMarker)'
      // })
      allNode.attr('transform', (d: any) => {
        return `translate(${d.x},${d.y})`
      })
    })
    // simulation.force('link').links(d3_links.value).distance(100)
  }

  // 绘制关系箭头
  function renderMarkers(svg: any) {
    // 定义箭头的标识
    const defs = svg.append('defs')
    const posMarker = defs
      .append('marker')
      .attr('id', 'posMarker')
      .attr('orient', 'auto')
      .attr('stroke-width', 2)
      .attr('markerUnits', 'strokeWidth')
      .attr('markerUnits', 'userSpaceOnUse')
      .attr('viewBox', '0 -5 10 10')
      .attr('refX', 8)
      .attr('refY', 0)
      .attr('markerWidth', 12)
      .attr('markerHeight', 12)
      .append('path')
      .attr('d', 'M 0 -5 L 10 0 L 0 5')
      .attr('fill', '#d4d1e9')
      .attr('stroke-opacity', 0.6)
    const activeMarker = defs
      .append('marker')
      .attr('id', 'activeMarker')
      .attr('orient', 'auto')
      .attr('stroke-width', 2)
      .attr('markerUnits', 'strokeWidth')
      .attr('markerUnits', 'userSpaceOnUse')
      .attr('viewBox', '0 -5 10 10')
      .attr('refX', 8)
      .attr('refY', 0)
      .attr('markerWidth', 12)
      .attr('markerHeight', 12)
      .append('path')
      .attr('d', 'M 0 -5 L 10 0 L 0 5')
      .attr('fill', '#A29AD6')
      .attr('stroke-opacity', 0.6)
    const errPosMarker = defs
      .append('marker')
      .attr('id', 'errPosMarker')
      .attr('orient', 'auto')
      .attr('stroke-width', 2)
      .attr('markerUnits', 'strokeWidth')
      .attr('markerUnits', 'userSpaceOnUse')
      .attr('viewBox', '0 -5 10 10')
      .attr('refX', 8)
      .attr('refY', 0)
      .attr('markerWidth', 12)
      .attr('markerHeight', 12)
      .append('path')
      .attr('d', 'M 0 -5 L 10 0 L 0 5')
      .attr('fill', '#F56C6C')
      .attr('stroke-opacity', 0.6)
    // const errNegativeMarker = defs
    //   .append('marker')
    //   .attr('id', 'errNegativeMarker')
    //   .attr('orient', 'auto')
    //   .attr('stroke-width', 2)
    //   .attr('markerUnits', 'strokeWidth')
    //   .attr('markerUnits', 'userSpaceOnUse')
    //   .attr('viewBox', '0 -5 10 10')
    //   .attr('refX', 2)
    //   .attr('refY', 0)
    //   .attr('markerWidth', 12)
    //   .attr('markerHeight', 12)
    //   .append('path')
    //   .attr('d', 'M 10 -5 L 0 0 L 10 5')
    //   .attr('fill', '#F56C6C')
    //   .attr('stroke-opacity', 0.6)

    // const negativeMarker = defs
    //   .append('marker')
    //   .attr('id', 'negativeMarker')
    //   .attr('orient', 'auto')
    //   .attr('stroke-width', 2)
    //   .attr('markerUnits', 'strokeWidth')
    //   .attr('markerUnits', 'userSpaceOnUse')
    //   .attr('viewBox', '0 -5 10 10')
    //   .attr('refX', 2)
    //   .attr('refY', 0)
    //   .attr('markerWidth', 12)
    //   .attr('markerHeight', 12)
    //   .append('path')
    //   .attr('d', 'M 10 -5 L 0 0 L 10 5')
    //   .attr('fill', '#d4d1e9')
    //   .attr('stroke-opacity', 0.6)
  }
  // 清除事件监听防止内存泄漏
  const removeListener = () => {
    if (!d3_instance.value) return
    d3_instance.value.on('.', null)
    d3_instance.value.selectAll('*').remove().on('.', null)
  }
  // 缩放功能
  const zoomHandle = (zoom: number) => {
    nextTick(() => {
      transformData.value.k = zoom <= 0.5 ? 0.5 : zoom >= 1.5 ? 1.5 : zoom
      d3_instance.value
        .transition()
        .duration(1000)
        .call(
          svg_zoom.transform,
          d3.zoomIdentity.translate(transformData.value.x, transformData.value.y).scale(transformData.value.k)
        )
    })
  }
  onMounted(() => {
    removeListener()
  })

  watch(
    () => transformData.value.y,
    () => {
      contextmenuFlag = false
      d3.select('.d3-context-menu').style('display', 'none')
    }
  )

  const getRectData = () => {
    const node = document.querySelector('.sketchpad')
    rectData.value = node?.getClientRects()
  }

  const monitorInsideIsOutside = () => {
    const data = {
      isDrag: !!curNodeWhenDrag.value,
      data: curNodeWhenDrag.value,
    }
    if (!curNodeWhenDrag.value) return
    VueEvent.emit('insideIsOutsideResult', data)
  }

  // 监听事件
  VueEvent.on('insideIsOutside', monitorInsideIsOutside)

  onUnmounted(() => {
    clearTimeout(timeoutLevae)
    clearTimeout(timeoutEnter)
    removeListener()
    d3_force_ref.value = null
    // 事件销毁
    VueEvent.off('insideIsOutside')
  })

  // 重置位置大小
  const reset = () => {
    transformData.value.k = 1
    transformData.value.x = 0
    transformData.value.y = 0
    zoomHandle(transformData.value.k)
  }

  const HandleZoom = (val: 'up' | 'down') => {
    val == 'up' ? zoomHandle((transformData.value.k += 0.1)) : zoomHandle((transformData.value.k -= 0.1))
  }

  // 初始化数据实例化d3引力图和数据改变修改d3引力图
  const initData = (
    graphData: {
      nodes: node[]
      links: link[]
    },
    isChange = false,
    changeData?: {
      nodes: node[]
      links: link[]
    }
  ) => {
    // 没有改变直接实例化
    if (!isChange) {
      curGraphData.value = graphData
      d3init()
    }
  }

  function addNode(newNode: node, id: string) {
    // setTimeout(() => {
    // 将现有节点的位置固定
    d3_nodes.value.forEach((node) => {
      node.fx = node.x
      node.fy = node.y
    })
    // 暂时增加 alphaTarget，使模拟更加活跃
    simulation.alphaTarget(0.01).restart()
    // 将新节点添加到数据集中
    d3_nodes.value.push(newNode)
    // 更新节点的选择集
    const updatedNodes = my_svg.selectAll('.node').data(d3_nodes.value)
    // 进入新节点
    creatNode(updatedNodes)
    const pItem = d3_nodes.value.find((item: node) => {
      return item.ip == id
    })
    // 为新节点设置合理的初始位置
    newNode.fx = 141 + pItem.x // 初始 x 位置
    newNode.fy = -141 + pItem.y // 初始 y 位置

    // 等待新节点稳定后，取消固定位置
    setTimeout(() => {
      d3_nodes.value.forEach((node) => {
        node.fx = null
        node.fy = null
      })
      simulation.alphaTarget(0)
    }, 0) // 1秒后取消固定位置
  }

  function addLink(newLink: link) {
    // 将现有节点的位置固定
    d3_nodes.value.forEach((node) => {
      node.fx = node.x
      node.fy = node.y
    })
    // 暂时增加 alphaTarget，使模拟更加活跃
    simulation.alphaTarget(0.3).restart()
    // // 将新链接添加到数据集中
    d3_links.value.push(newLink)
    simulation.force('link').links(d3_links.value).distance(190)
    // // 更新链接的选择集
    const undateLink = my_svg.select('.links').selectAll('.link').data(d3_links.value)
    creatLink(undateLink)
    // // 等待模拟稳定后，取消固定位置
    setTimeout(() => {
      d3_nodes.value.forEach((node: node) => {
        node.fx = undefined
        node.fy = undefined
      })
      simulation.alphaTarget(0)
    }, 1000)
  }

  const creatNode = (updatedNodes: any) => {
    const node = updatedNodes
      .enter()
      .append('g')
      .attr('class', 'node')
      .call(d3_drag(simulation))
      .call(d3_drag(simulation))
      .on('click', function (this: any, e: any) {
        e.stopPropagation()
        // simulation.stop()
        // simulation.alphaTarget(-0.1)
        const _node = d3.select(this)
        const nodeData: any = _node.data()[0]
        const circle = _node.select('.circle_drag')
        const use = _node.select('use')
        const allCircle = d3.selectAll('.circle_drag')
        const allUse = d3.selectAll('use')
        // 更新这些边的颜色
        my_svg.selectAll('.link').classed('highlight', false).attr('marker-end', 'url(#posMarker)') // 先清除所有高亮
        const activeLinks = uniqueByClientIpAndServerIp(getLinkByIpAndallLinks(nodeData.ip))
        const activeNode = my_svg
          .selectAll('.link')
          .filter((link: link) =>
            activeLinks.some((rl) => rl.serverIp === link.serverIp && rl.clientIp === link.clientIp)
          )
        activeNode.classed('highlight', true).attr('marker-end', 'url(#activeMarker)')
        // 更新箭头的颜色
        if (props.isPath == 'pathTracing') {
          nodeData['isPath'] = '路径追踪'
          emits('show-info', nodeData)
          allCircle.style('fill', '#fff')
          allUse.attr('xlink:href', (d: any) => `#${d.icon}`)
          nodeData.icon === 'user' ? circle.style('fill', '#7A68F2') : circle.style('fill', '#847AB1')
          nodeData.icon === 'user' ? use.attr('xlink:href', '#user-active') : use.attr('xlink:href', '#server-active')
        }
      })
      // 右键功能
      .on('contextmenu', (event: any, d: node) => {
        event.stopPropagation()
        event.preventDefault()

        leaveFlag = true
        my_tooltip.style('opacity', 0).style('display', 'none')
        contextmenuFlag = true

        const { width, height } = rectData.value[0]
        let offsetX = 0
        let offsetY = 0
        event.pageX + 200 > width ? (offsetX = event.pageX - 200) : (offsetX = event.pageX)
        event.pageY + 215 > height ? (offsetY = event.pageY - 216) : (offsetY = event.pageY)

        const img = require(`@/assets/chase/ip.svg`)
        const contextMenu = d3.select('.d3-context-menu')
        // const remark = d
        contextMenu
          .html(
            `
          <div class="d3-context-menu-top">
             <div style="display:flex;align-items:center">
               <img style="height:16px;width:16px;object-fit: contain;margin-right:4px" src="${img}">
               <div class="d3-context-menu-top-ip" style="font-size: 18px;color: #303133;font-weight: 500;">${d.ip}</div>
             </div>
             <div> <span style="font-size: 13px;color: #9D9BAA;">来源：内网字段<span> </div>
          </div>
          <div class="d3-context-menu-content"></div>
          `
          )
          .on('click', () => {
            event.stopPropagation()
            event.preventDefault()
          })
          .style('height', '216px')
          .style('left', `${offsetX}px`)
          .style('top', `${offsetY}px`)
          .style('display', 'block')
        contextMenu
          .select('.d3-context-menu-content')
          .append('div')
          .attr('class', 'd3-context-menu-item')
          .text('深度挖掘')
          .on('click', () => {
            contextmenuCloseFlag = false
            emits('show-relate', d)
          })
        const word = props.startSelf ? '关闭自定义连线' : '自定义连线'
        contextMenu
          .select('.d3-context-menu-content')
          .append('div')
          .attr('class', 'd3-context-menu-item')
          // .style('cursor', 'not-allowed')
          .text(word)
          .on('click', () => {
            contextmenuCloseFlag = false
            props.startSelf ? emits('close-filter') : emits('show-filter', d)
          })
        contextMenu
          .select('.d3-context-menu-content')
          .append('div')
          .attr('class', 'd3-context-menu-item')
          .style('cursor', 'not-allowed')
          .text('隐藏该节点')
          .on('click', () => {
            // contextmenuCloseFlag = false
            // const curIp = d.ip
            // const relateLink = d3_links.value.filter((link: link) => link.serverIp == curIp)
          })
        contextMenu
          .select('.d3-context-menu-content')
          .append('div')
          .attr('class', 'd3-context-menu-item-last')
          .append('div')
          .text('标记颜色')
        const last = d3.selectAll('.d3-context-menu-item-last')
        last.append('div').attr('class', 'shader')
        const shader = d3.selectAll('.shader')
        shader
          .append('div')
          .attr('class', 'shader-item shader_FFF0AC')
          .on('click', () => {
            // 改变数据(当前)
            d3_nodes.value = d3_nodes.value.map((i: node) => {
              if (i.ip == d.ip) {
                i.remark = i.remark == '#FFF0AC' ? '' : '#FFF0AC'
              }
              return i
            })
            // 发请求修改数据
            emits('remark-node-color', d)
            // 改变状态
            d3.selectAll('.shader_FFF0AC')
              .selectAll('.shader-item-text')
              .style('display', () => {
                return d.remark == '#FFF0AC' ? 'block' : 'none'
              })
            d3.selectAll('.aperture')
              .attr('fill', (e: node) => {
                return e.remark
              })
              .style('display', (d: node) => {
                return d.remark ? 'block' : 'none'
              })

            contextmenuCloseFlag = false
          })
          .append('svg')
          .attr('class', 'shader-item-text')
          .attr('height', 20)
          .attr('width', 20)
          .style('display', () => {
            return d.remark == '#FFF0AC' ? 'block' : 'none'
          })
          .append('use')
          .attr('height', 10)
          .attr('width', 12)
          .attr('xlink:href', `#remark`)

        shader
          .append('div')
          .attr('class', 'shader-item shader_BDDAFF')
          .style('background', () => {
            return d.remark == '#BDDAFF' ? '#5A91FF' : '#e2ecff'
          })
          .on('click', () => {
            // 改变数据(当前)
            d3_nodes.value = d3_nodes.value.map((i: node) => {
              if (i.ip == d.ip) {
                i.remark = i.remark == '#BDDAFF' ? '' : '#BDDAFF'
              }
              return i
            })
            // 发请求修改数据
            emits('remark-node-color', d)
            // 改变状态
            d3.selectAll('.shader_C9DEFF')
              .selectAll('.shader-item-text')
              .style('display', () => {
                return d.remark == '#BDDAFF' ? 'block' : 'none'
              })
            d3.selectAll('.aperture')
              .attr('fill', (e: node) => {
                return e.remark
              })
              .style('display', (d: node) => {
                return d.remark ? 'block' : 'none'
              })
            contextmenuCloseFlag = false
          })
          .append('svg')
          .attr('class', 'shader-item-text')
          .attr('height', 20)
          .attr('width', 20)
          .style('display', () => {
            return d.remark == '#BDDAFF' ? 'block' : 'none'
          })
          .append('use')
          .attr('height', 10)
          .attr('width', 12)
          .attr('xlink:href', `#remark_2`)

        changeGraphStyle(d.id)
      })
      // 鼠标移入功能
      .on(
        'mouseenter',
        (event: any, d: any) => {
          event.stopPropagation()
          leaveFlag = false
          clearTimeout(timeoutEnter)
          timeoutFlag = false
          timeoutEnter = setTimeout(async () => {
            timeoutFlag1 = false
            if (leaveFlag || contextmenuFlag) return
            const { width, height } = rectData.value[0]
            const { x, y } = event.target.getBoundingClientRect()
            const lapse = 11
            let offsetX = 0
            let offsetY = 0
            let translateX = x + lapse - 55
            let translateY = y + lapse - 45
            if (x + 360 > width) {
              offsetX = -100
              translateX = x - lapse - 55
            }
            if (y + 210 > height) {
              offsetY = -100
              translateY = y - lapse - 45
            }
            const code = d.countryCode?.toLowerCase()
            const img = code && require(`@/assets/flag/${code}.png`)
            const isExpand = d.isUnfold == 1 ? 'pointer' : 'not-allowed'
            const { data } = await QueryAssetByIpApi({
              ip: d.ip,
              startTime: '2024-11-18 15:04:20',
              endTime: '2024-11-18 16:14:20',
              status: props.isPath == 'pathTracing' ? '路径追踪' : '威胁狩猎',
            })
            my_tooltip
              .html(
                `<div class="tooltip_top">
                  <div style="display:flex;align-items:center">
                      <span class="tooltip_top_ip" style="font-weight: 500;font-size: 20px;color: #303133;">${d?.ip}</span>
                      ` +
                  ` ${
                    code ? `<img style="height:16px;width:16px;object-fit: contain;margin-left:4px" src="${img}">` : ''
                  } ` +
                  `
                  </div>
                  <div class="tooltip_top_btn"><div class="expand-img"></div>拓展全部</div>
              </div>
              <div style="padding: 15px">
                  <div class="tip-item" style="margin-top: 0px">
                      <div class="left">业务名称</div>
                      <div class="right">${data.businessName || '-'}</div>
                  </div>
                  <div class="tip-item" style="margin-top: 10px">
                      <div class="left">来源</div>
                      <div class="right">${data.datasource || data.netTypeStr || '-'}</div>
                  </div>
                  <div class="tip-item" style="margin-top: 10px">
                      <div class="left">属性</div>
                      <div class="right">${data.netTypeStr || '-'}</div>
                  </div>
              </div>`
              )
              .style('opacity', 1)
              .style('display', 'block')
              .style('cursor', isExpand)
              .style('transform', `translate(calc(${translateX}px + ${offsetX}%), calc(${translateY}px + ${offsetY}%)`)
              .on('click', () => {
                if (d.isUnfold == 1) {
                  const nodeSrc = d3.select(event.target)
                  const idx = d3_nodes.value.findIndex((item: node) => {
                    return d.ip == item.ip
                  })
                  d3_nodes.value[idx]['isUnfold'] = d3_nodes.value[idx].isUnfold == 1 ? 2 : 1
                  emits('expand-all', d.ip)
                  nodeSrc.selectAll('.expand-text').text('-')
                }
              })
              .on('mouseenter', () => {
                timeoutFlag = true
                my_tooltip.style('opacity', 1).style('display', 'block')
              })
              .on('mouseleave', (event: MouseEvent) => {
                if (timeoutFlag1) return
                my_tooltip.style('opacity', 0).style('display', 'none')
                clearGraphStyle()
              })
            // 遍历节点，并调整图的样式
            changeGraphStyle(d.id)
          }, 800)
        },
        false
      )
      // 鼠标移出，关闭痰喘
      .on(
        'mouseleave',
        (event: any, d: any) => {
          event.stopPropagation()
          clearTimeout(timeoutLevae)
          if (timeoutFlag) return
          leaveFlag = true
          timeoutLevae = setTimeout(() => {
            if (timeoutFlag) return
            my_tooltip.style('opacity', 0).style('display', 'none')
            clearGraphStyle()
          }, 300)
        },
        false
      )
    const div = node
      .append('g')
      // huovr效果：放大加阴影
      .on('mouseenter', (e: any) => {
        // e.stopPropagation()
        const target = e.target
        d3.select(target)
          .attr('filter', `drop-shadow(0px 0px 4px rgba(77, 81, 86, 0.25))`)
          .transition() // 开始新的过渡
          .duration(200)
          .attr('transform', 'scale(1.15)')
      })
      // 移除huovr效果
      .on(
        'mouseleave',
        (event: any, d: any) => {
          event.stopPropagation()
          const target = event.target
          d3.select(target)
            .attr('filter', null)
            .transition() // 开始新的过渡
            .duration(200)
            .attr('transform', 'scale(1.0)')
        },
        false
      )
    div
      .append('circle')
      .attr('class', 'aperture')
      .attr('r', symbolSize / 2 + 4)
      .attr('fill', (d: { remark: string }) => d.remark)
      // .attr('stroke', (d: { remark: string }) => d.remark)
      .style('display', (d: node) => {
        return d.remark ? 'block' : 'none'
      })
    div
      .append('circle')
      .attr('class', 'circle_drag')
      .attr('r', symbolSize / 2)
      .attr('fill', '#fff')
      .attr('stroke', '#EAE7FF')

    // node上的图片
    div
      .append('use')
      .attr('height', 20)
      .attr('width', 20)
      .attr('z-index', 20)
      .attr('transform', function (this: any, d: any) {
        return `translate( -10 , -10 )`
      })
      .attr('cursor', 'pointer')
      .attr('xlink:href', (d: any) => `#${d.icon}`)
    // 展开收缩节点
    const group = div
      .append('g')
      .attr('transform', `translate(17, -17)`)
      .attr('cursor', 'pointer')
      .attr('display', (d: node) => (d.isUnfold ? 'unset' : 'none'))
      .on(
        'click',
        (e: any, d: node) => {
          e.stopPropagation()
          if (d.isUnfold == 1) {
            emits('expand-one', d.ip)

            setTimeout(() => {
              // 通过该节点找到其所有子节点
              const relatedLinks = d3_links.value.filter((link: link) => {
                return d.ip == link.clientIp
              })

              const nodes: string[] = []
              relatedLinks.forEach((link) => {
                nodes.push(link.serverIp)
              })
              // 通过子节点找到所有父节点
              nodes.forEach((ip) => {
                const arr: string[] = []
                d3_links.value.forEach((link: link) => {
                  if (link.serverIp == ip) {
                    arr.push(link.clientIp)
                  }
                })
                // 父节点的状态修改
                d3_nodes.value.forEach((_: node, index: number) => {
                  if (arr.includes(d3_nodes.value[index].ip)) {
                    // warp['isUnfold'] = warp['isUnfold'] == 1 ? 2 : 1
                    d3_nodes.value[index]['isUnfold'] = d3_nodes.value[index]['isUnfold'] == 1 ? 2 : 1
                  }
                })
                // console.log(arr1)
              })
              d3.selectAll('.expand-text').text((d: node) => {
                let str = '+'
                if (d.isUnfold == 2) str = '-'
                return str
              })
            }, 10)
          } else if (d.isUnfold == 2) {
            // 通过该节点找到其所有子节点
            const relatedLinks = d3_links.value.filter((link: link) => {
              return d.ip == link.clientIp
            })
            const nodes: string[] = []
            relatedLinks.forEach((link) => {
              nodes.push(link.serverIp)
            })
            // 通过子节点找到所有父节点
            nodes.forEach((ip) => {
              const arr = d3_links.value.filter((link: link) => {
                return link.serverIp == ip
              })
              // 父节点的状态修改
              d3_nodes.value = d3_nodes.value.map((warp: node) => {
                const res = arr.find((relatedLinksItem) => {
                  return warp.ip == relatedLinksItem.clientIp
                })
                if (warp.ip == res?.clientIp) {
                  if (warp['isUnfold']) {
                    warp['isUnfold'] = warp['isUnfold'] == 1 ? 2 : 1
                  }
                }
                return warp
              })
            })

            d3.selectAll('.expand-text').text((d: node) => {
              let str = '+'
              if (d.isUnfold == 2) str = '-'
              return str
            })
            emits('del-node', d.ip)
          }
        },
        true
      )
    // 圆
    group.append('circle').attr('class', 'launch-tip').attr('r', '7').attr('fill', '#fff').attr('stroke', '#D1CDE8')
    // 圆里面的符号'+' 或者'-'
    group
      .append('text')
      .attr('class', 'expand-text')
      .attr('x', 0) // 文字的x坐标 (相对于g元素)
      .attr('y', -1.5) // 文字的y坐标 (相对于g元素)
      .attr('text-anchor', 'middle') // 文字水平居中
      .attr('alignment-baseline', 'central') // 文字垂直居中
      .attr('fill', '#D1CDE8') // 文字颜色
      .style('font-size', '14px') // 设置字体大小
      .style('font-weight', 'bold')
      .text((d: node) => {
        let str = '+'
        if (d.isUnfold == 2) str = '-'
        return str
      })
    // 文字内容(ip)
    const texts = node.append('g').attr('class', 'text').attr('text-anchor', 'middle').attr('dominant-baseline', 'top')
    texts.append('rect').attr('class', 'host-text-bg').attr('fill', 'transparent').attr('height', '18px')
    texts
      .append('text')
      .attr('id', (d: any) => d.id)
      .attr('class', 'host-name')
      .attr('x', function (this: any, d: any) {
        return textBreaking(d3.select(this), d.ip.slice(0, 17) || d.assetName.slice(0, 17))
      })
    // 文字内容(国旗)
    const flag = node.select('.text')
    flag
      .append('svg')
      .attr('class', 'my-flag')
      .attr('x', function (this: any, d: any) {
        return this.parentNode.getBoundingClientRect().width / 2 + 4
      })
      .attr('y', 27)
      .attr('width', 16)
      .attr('height', 16)
      .append('image')
      .attr('xlink:href', (d: any) => {
        const code = d.countryCode?.toLowerCase()
        return code && require(`@/assets/flag/${code}.png`)
      })
      .attr('width', 16)
      .attr('height', 16)
  }

  const creatLink = (link: any) => {
    const links = link
    links
      .enter()
      .append('path')
      .attr('class', 'link')
      .attr('stroke', (d: any) => {
        return d?.alarm?.length > 0 ? '#F56C6C' : '#d4d1e9'
      })
      .attr('cursor', 'pointer')
      .attr('marker-end', 'url(#posMarker)')
      .on('click', function (this: any, event: any) {
        const _node = d3.select(this)
        const nodeData: any = _node.data()[0]
        // emits('click', nodeData)
      })
  }

  // 根据ip和所有线段获取所有相关的线段
  const getLinkByIpAndallLinks = (ip: string, visited = new Set()): link[] => {
    if (visited.has(ip)) return []

    visited.add(ip)

    const relatedLinks = d3_links.value.filter(
      (link: { clientIp: string; serverIp: string }) => link.clientIp === ip || link.serverIp === ip
    )

    const connectedNodes = new Set()
    relatedLinks.forEach((link: { clientIp: unknown; serverIp: unknown }) => {
      if (link.clientIp !== ip) connectedNodes.add(link.clientIp)
      if (link.serverIp !== ip) connectedNodes.add(link.serverIp)
    })

    connectedNodes.forEach((connectedNodeId: any) => {
      relatedLinks.push(...getLinkByIpAndallLinks(connectedNodeId, visited))
    })

    return relatedLinks
  }

  function uniqueByClientIpAndServerIp(arr: link[]): link[] {
    const map = new Map<string, link>()
    arr.forEach((item) => {
      const key = `${item.clientIp ?? ''}->${item.serverIp ?? ''}`
      if (!map.has(key)) {
        map.set(key, item)
      }
    })
    return Array.from(map.values())
  }

  // 移除节点和边
  function removeNodesAndLinks(descendants: node[], relatedLinks: link[]) {
    if (descendants.length > 0 || relatedLinks.length > 0) {
      // 从links数组中移除相关的边
      const map2 = new Map<string, link>()
      relatedLinks.forEach((record) => map2.set(generateUniqueKey(record), record))

      // 过滤 arr1，找出不包含在 arr2 中的记录
      d3_links.value = d3_links.value.filter((record) => !map2.has(generateUniqueKey(record)))
      // 从nodes数组中移除子节点
      descendants.forEach((node) => {
        d3_nodes.value = d3_nodes.value.filter((item) => {
          return item.ip != node.ip
        })
      })

      // 更新SVG中的图形元素
      my_svg.select('.links').selectAll('.link').data(d3_links.value).exit().remove()
      my_svg.selectAll('.node').data(d3_nodes.value).exit().remove()

      // 调整模拟以反映新的数据结构
      simulation.nodes(d3_nodes.value)
      simulation.force('link').links(d3_links.value).distance(190)
      simulation.alpha(0.3).restart() // 重启模拟以适应新的布局
    }
    // 保持其他节点位置不变
    d3_nodes.value.forEach((node) => {
      node.fx = node.x
      node.fy = node.y
    })
  }

  // 定义一个函数来生成唯一的字符串标识符
  function generateUniqueKey(record: link): string {
    return `${record.clientIp}-${record.serverIp || 'null'}`
  }

  function showText(show: string) {
    const textNode = my_svg.selectAll('.text').selectAll('text')
    textNode.selectAll('tspan').remove()
    textNode.attr('x', function (this: any, d: any) {
      const text =
        show == 'assetName'
          ? (d.assetName && d.assetName.slice(0, 17)) || d.ip.slice(0, 17)
          : d.ip.slice(0, 17) || (d.assetName && d.assetName.slice(0, 17))
      return textBreaking(d3.select(this), text)
    })
  }

  defineExpose({
    reset,
    HandleZoom,
    initData,
    addNode,
    addLink,
    removeNodesAndLinks,
    showText,
  })
</script>

<template>
  <div class="d3_force_box">
    <div v-if="curNodeWhenDrag" class="false_node" :style="{ left: curNodeWhenDragLeft, top: curNodeWhenDragTop }">
      <div>
        <div class="img">
          <img v-if="curNodeWhenDrag.icon === 'user'" alt="" :src="require('@/assets/chase/default.svg')" />
          <img v-else alt="" :src="require('@/assets/chase/server.svg')" />
        </div>
        <div>{{ curNodeWhenDrag.ip }}</div>
      </div>
    </div>
    <svg
      id="user"
      aria-hidden="true"
      style="enable-background: new 0 0 22 22; width: 0; height: 0"
      viewBox="0 0 22 22"
      x="0px"
      xml:space="preserve"
      xmlns="http://www.w3.org/2000/svg"
      xmlns:xlink="http://www.w3.org/1999/xlink"
      y="0px"
    >
      <path
        d="M3.1,15.3c1.7-0.6,2.6-0.9,3-1c0.3-0.1,0.6-0.2,0.9-0.4c0.2-0.1,0.3-0.2,0.4-0.3v-0.1
						c-0.1-0.1-0.1-0.3-0.1-0.4c0-0.1,0-0.2,0-0.2l-0.4-0.3c-0.3-0.2-0.5-0.5-0.6-0.9C6.1,11.4,6.1,11.2,6,11v-0.2
						c-0.2-0.2-0.3-0.3-0.5-0.5C5.3,10,5.1,9.6,5,9.3C4.9,8.8,4.8,8.3,4.9,7.9c0-0.3,0.1-0.6,0.2-0.8c0-0.1,0.1-0.3,0.2-0.4
						c0-0.5,0-1,0.1-1.5c0.1-0.5,0.2-1,0.3-1.5H5.4c-0.2,0-0.5,0-0.7,0.1C4.5,3.8,4.2,3.9,4,4C3.7,4.1,3.5,4.2,3.3,4.4
						C3,4.6,2.8,4.8,2.6,5.1c-0.2,0.3-0.4,0.7-0.5,1C2.1,6.5,2,6.9,2,7.2c0,0.4-0.1,0.8-0.1,1.2C1.7,8.7,1.6,8.9,1.6,9.2
						c0,0.2,0,0.3,0,0.5c0,0.1,0.1,0.3,0.2,0.4c0.1,0.1,0.1,0.2,0.2,0.2c0.1,0,0.1,0.1,0.2,0.1l0.1,0.5c0.1,0.1,0.1,0.3,0.2,0.4
						c0,0.1,0.1,0.2,0.2,0.3c0.2,0.1,0.3,0.3,0.5,0.4c0.2,0.2,0.2,0.4,0.2,0.6v0.4c0,0.1,0,0.3-0.1,0.4c-0.1,0.1-0.1,0.3-0.2,0.4
						C3,14,2.8,14.2,2.6,14.3c-0.3,0.2-0.6,0.3-1,0.4c-0.2,0.1-1.1,0.4-1.7,0.6v4h0.1c0-0.2,0.1-0.5,0.1-0.8v-0.1
						c0.1-0.8,0.5-1.5,1.1-2.1c0.5-0.4,1.1-0.7,1.7-0.9L3.1,15.3L3.1,15.3z M21.3,14.7c-0.3-0.1-0.6-0.2-0.8-0.3
						c-0.2-0.1-0.5-0.2-0.6-0.4c-0.1-0.1-0.3-0.3-0.4-0.5c-0.1-0.2-0.1-0.3-0.1-0.5c0-0.3,0-0.5,0.1-0.8c0-0.1,0.1-0.2,0.2-0.2
						c0.1-0.1,0.1-0.1,0.2-0.2l0.2-0.2c0.1-0.1,0.2-0.2,0.2-0.3c0.1-0.1,0.1-0.3,0.2-0.4c0-0.2,0.1-0.3,0.1-0.5
						c0.1,0,0.2-0.1,0.3-0.2c0.1-0.1,0.2-0.2,0.2-0.3c0.1-0.2,0.1-0.3,0.1-0.5c0-0.1,0-0.3,0-0.4c0-0.1-0.1-0.2-0.1-0.3
						c0-0.1-0.1-0.1-0.1-0.2c0-0.4,0-0.9-0.1-1.3c-0.1-0.4-0.1-0.8-0.3-1.2c-0.1-0.4-0.4-0.8-0.6-1.1c-0.1-0.2-0.3-0.3-0.4-0.4
						c-0.4-0.3-0.9-0.5-1.4-0.7c-0.3-0.1-0.5-0.1-0.8-0.1h-0.4c0.1,0.5,0.2,0.9,0.3,1.4c0.1,0.5,0.1,0.9,0.1,1.4
						c0.2,0.3,0.3,0.6,0.3,0.9c0.1,0.3,0.1,0.7,0,1.1c0,0.4-0.2,0.8-0.3,1.2c-0.2,0.4-0.4,0.7-0.7,0.9l-0.1,0.2v0.1
						c-0.1,0.2-0.2,0.5-0.3,0.7c-0.2,0.5-0.5,0.8-0.9,1.1v0.4c0,0.2,0,0.4-0.1,0.5c0.2,0.1,0.3,0.2,0.5,0.3c0.3,0.1,0.6,0.2,1,0.3
						c0.3,0.1,0.6,0.2,1,0.4l1,0.4l0.5,0.2c0.5,0.2,0.9,0.4,1.3,0.8c0.5,0.4,0.8,0.9,1,1.5c0.1,0.4,0.2,0.7,0.2,1.1
						c0,0.2,0.1,0.3,0.1,0.5H22v-4.1c-0.3-0.1-0.6-0.3-1-0.4H21.3z"
        style="fill: #8c7bf8"
      />
      <path
        id="形状_1_"
        d="M20.9,17.6c-0.1-0.4-0.3-0.7-0.6-0.9c-0.3-0.2-0.6-0.4-1-0.6c-0.4-0.1-2-0.9-2.4-1
						c-0.4-0.1-0.8-0.2-1.1-0.4c-0.3-0.1-0.6-0.3-0.9-0.6c-0.2-0.2-0.4-0.4-0.5-0.6c-0.1-0.2-0.1-0.4-0.2-0.7c0-0.4,0-0.8,0.2-1.1
						c0-0.1,0.1-0.2,0.2-0.3c0.1-0.1,0.2-0.2,0.3-0.2c0.1-0.1,0.2-0.2,0.3-0.3c0.1-0.1,0.2-0.3,0.3-0.4c0.1-0.2,0.2-0.4,0.2-0.6
						c0.1-0.2,0.1-0.4,0.2-0.7c0.1-0.1,0.2-0.1,0.3-0.2c0.1-0.1,0.2-0.3,0.3-0.4c0.1-0.2,0.2-0.5,0.2-0.7c0-0.2,0-0.4,0-0.6
						c0-0.1-0.1-0.3-0.1-0.4c0-0.1-0.1-0.2-0.2-0.3c0-0.6,0-1.2-0.1-1.8c-0.1-0.5-0.2-1.1-0.4-1.6c-0.2-0.6-0.4-1.1-0.8-1.5
						c-0.2-0.2-0.4-0.4-0.6-0.6c-0.3-0.2-0.5-0.4-0.8-0.5c-0.3-0.2-0.7-0.3-1-0.4C12.4,0,12.1,0,11.7,0c-0.3,0-0.6,0-0.9,0.1
						c-0.3,0.1-0.6,0.2-0.9,0.3C9,0.7,8.3,1.3,7.8,2C7.5,2.4,7.3,2.9,7.2,3.4C7,4.4,6.8,5.5,6.9,6.6C6.8,6.7,6.7,6.9,6.6,7
						C6.6,7.2,6.5,7.3,6.5,7.5c0,0.2,0,0.5,0.1,0.7c0,0.2,0.1,0.4,0.2,0.6C6.9,8.8,7,8.9,7.1,9c0.1,0.1,0.2,0.1,0.3,0.1l0.2,0.7
						c0.1,0.2,0.1,0.4,0.2,0.6c0.1,0.2,0.2,0.3,0.3,0.4c0.2,0.2,0.4,0.4,0.6,0.6C8.9,11.7,9,12,9,12.3v0.6c0,0.2,0,0.4-0.1,0.6
						c-0.1,0.2-0.2,0.4-0.3,0.5c-0.2,0.2-0.4,0.4-0.7,0.5c-0.4,0.3-0.8,0.5-1.2,0.6c-0.5,0.1-2.5,0.9-3,1.1
						c-0.4,0.1-0.8,0.4-1.2,0.6C2.2,17.1,2,17.5,2,18c-0.1,0.6-0.2,1.1-0.1,1.7c0,0.4,0.1,0.8,0.3,1.1c0.3,0.2,0.7,0.3,1.1,0.4
						C6.1,21.7,9,22,11.9,22c2.8,0,5.7-0.4,8.4-1.1c0.2-0.1,0.4-0.2,0.6-0.3c0.1-0.2,0.2-0.5,0.2-0.7c0-0.3,0.1-0.7,0.1-1
						c0-0.3,0-0.1,0-0.4C21,18.2,21,17.9,20.9,17.6L20.9,17.6z M12.9,19.7c-0.1,0.1-0.2,0.3-0.4,0.4c-0.1,0.1-0.3,0.2-0.4,0.3
						c-0.1,0.1-0.2,0.1-0.3,0.1c-0.1,0-0.2-0.1-0.3-0.1c-0.1-0.1-0.3-0.2-0.4-0.3c-0.1-0.1-0.3-0.3-0.4-0.4
						c-0.1-0.1-0.1-0.2-0.2-0.3c0-0.2,0-0.4,0.1-0.5c0-0.2,0.1-0.5,0.2-0.8s0.1-0.5,0.2-0.8c0.1-0.3,0.1-0.4,0.2-0.6
						c-0.1-0.1-0.3-0.3-0.4-0.4c-0.1-0.1-0.1-0.2-0.1-0.3c0-0.1,0-0.2,0.1-0.4c0.1-0.1,0.1-0.2,0.2-0.3c0.1-0.1,0.2-0.2,0.3-0.3h1
						l0.3,0.3c0.1,0.1,0.2,0.2,0.3,0.3c0.1,0.1,0.1,0.2,0.1,0.4c0,0.2-0.1,0.3-0.1,0.4c-0.1,0.1-0.2,0.3-0.4,0.4
						c0,0.2,0.1,0.4,0.1,0.6l0.2,0.8c0,0.3,0.1,0.5,0.1,0.7c0,0.2,0,0.3,0.1,0.5C13,19.6,13,19.7,12.9,19.7L12.9,19.7z"
        style="fill: #8c7bf8"
      />
    </svg>
    <svg
      id="user-active"
      aria-hidden="true"
      style="enable-background: new 0 0 22 22; width: 0; height: 0"
      viewBox="0 0 22 22"
      x="0px"
      xml:space="preserve"
      xmlns="http://www.w3.org/2000/svg"
      xmlns:xlink="http://www.w3.org/1999/xlink"
      y="0px"
    >
      <path
        d="M3.1,15.3c1.7-0.6,2.6-0.9,3-1c0.3-0.1,0.6-0.2,0.9-0.4c0.2-0.1,0.3-0.2,0.4-0.3v-0.1
						c-0.1-0.1-0.1-0.3-0.1-0.4c0-0.1,0-0.2,0-0.2l-0.4-0.3c-0.3-0.2-0.5-0.5-0.6-0.9C6.1,11.4,6.1,11.2,6,11v-0.2
						c-0.2-0.2-0.3-0.3-0.5-0.5C5.3,10,5.1,9.6,5,9.3C4.9,8.8,4.8,8.3,4.9,7.9c0-0.3,0.1-0.6,0.2-0.8c0-0.1,0.1-0.3,0.2-0.4
						c0-0.5,0-1,0.1-1.5c0.1-0.5,0.2-1,0.3-1.5H5.4c-0.2,0-0.5,0-0.7,0.1C4.5,3.8,4.2,3.9,4,4C3.7,4.1,3.5,4.2,3.3,4.4
						C3,4.6,2.8,4.8,2.6,5.1c-0.2,0.3-0.4,0.7-0.5,1C2.1,6.5,2,6.9,2,7.2c0,0.4-0.1,0.8-0.1,1.2C1.7,8.7,1.6,8.9,1.6,9.2
						c0,0.2,0,0.3,0,0.5c0,0.1,0.1,0.3,0.2,0.4c0.1,0.1,0.1,0.2,0.2,0.2c0.1,0,0.1,0.1,0.2,0.1l0.1,0.5c0.1,0.1,0.1,0.3,0.2,0.4
						c0,0.1,0.1,0.2,0.2,0.3c0.2,0.1,0.3,0.3,0.5,0.4c0.2,0.2,0.2,0.4,0.2,0.6v0.4c0,0.1,0,0.3-0.1,0.4c-0.1,0.1-0.1,0.3-0.2,0.4
						C3,14,2.8,14.2,2.6,14.3c-0.3,0.2-0.6,0.3-1,0.4c-0.2,0.1-1.1,0.4-1.7,0.6v4h0.1c0-0.2,0.1-0.5,0.1-0.8v-0.1
						c0.1-0.8,0.5-1.5,1.1-2.1c0.5-0.4,1.1-0.7,1.7-0.9L3.1,15.3L3.1,15.3z M21.3,14.7c-0.3-0.1-0.6-0.2-0.8-0.3
						c-0.2-0.1-0.5-0.2-0.6-0.4c-0.1-0.1-0.3-0.3-0.4-0.5c-0.1-0.2-0.1-0.3-0.1-0.5c0-0.3,0-0.5,0.1-0.8c0-0.1,0.1-0.2,0.2-0.2
						c0.1-0.1,0.1-0.1,0.2-0.2l0.2-0.2c0.1-0.1,0.2-0.2,0.2-0.3c0.1-0.1,0.1-0.3,0.2-0.4c0-0.2,0.1-0.3,0.1-0.5
						c0.1,0,0.2-0.1,0.3-0.2c0.1-0.1,0.2-0.2,0.2-0.3c0.1-0.2,0.1-0.3,0.1-0.5c0-0.1,0-0.3,0-0.4c0-0.1-0.1-0.2-0.1-0.3
						c0-0.1-0.1-0.1-0.1-0.2c0-0.4,0-0.9-0.1-1.3c-0.1-0.4-0.1-0.8-0.3-1.2c-0.1-0.4-0.4-0.8-0.6-1.1c-0.1-0.2-0.3-0.3-0.4-0.4
						c-0.4-0.3-0.9-0.5-1.4-0.7c-0.3-0.1-0.5-0.1-0.8-0.1h-0.4c0.1,0.5,0.2,0.9,0.3,1.4c0.1,0.5,0.1,0.9,0.1,1.4
						c0.2,0.3,0.3,0.6,0.3,0.9c0.1,0.3,0.1,0.7,0,1.1c0,0.4-0.2,0.8-0.3,1.2c-0.2,0.4-0.4,0.7-0.7,0.9l-0.1,0.2v0.1
						c-0.1,0.2-0.2,0.5-0.3,0.7c-0.2,0.5-0.5,0.8-0.9,1.1v0.4c0,0.2,0,0.4-0.1,0.5c0.2,0.1,0.3,0.2,0.5,0.3c0.3,0.1,0.6,0.2,1,0.3
						c0.3,0.1,0.6,0.2,1,0.4l1,0.4l0.5,0.2c0.5,0.2,0.9,0.4,1.3,0.8c0.5,0.4,0.8,0.9,1,1.5c0.1,0.4,0.2,0.7,0.2,1.1
						c0,0.2,0.1,0.3,0.1,0.5H22v-4.1c-0.3-0.1-0.6-0.3-1-0.4H21.3z"
        style="fill: #fff"
      />
      <path
        id="形状_1_"
        d="M20.9,17.6c-0.1-0.4-0.3-0.7-0.6-0.9c-0.3-0.2-0.6-0.4-1-0.6c-0.4-0.1-2-0.9-2.4-1
						c-0.4-0.1-0.8-0.2-1.1-0.4c-0.3-0.1-0.6-0.3-0.9-0.6c-0.2-0.2-0.4-0.4-0.5-0.6c-0.1-0.2-0.1-0.4-0.2-0.7c0-0.4,0-0.8,0.2-1.1
						c0-0.1,0.1-0.2,0.2-0.3c0.1-0.1,0.2-0.2,0.3-0.2c0.1-0.1,0.2-0.2,0.3-0.3c0.1-0.1,0.2-0.3,0.3-0.4c0.1-0.2,0.2-0.4,0.2-0.6
						c0.1-0.2,0.1-0.4,0.2-0.7c0.1-0.1,0.2-0.1,0.3-0.2c0.1-0.1,0.2-0.3,0.3-0.4c0.1-0.2,0.2-0.5,0.2-0.7c0-0.2,0-0.4,0-0.6
						c0-0.1-0.1-0.3-0.1-0.4c0-0.1-0.1-0.2-0.2-0.3c0-0.6,0-1.2-0.1-1.8c-0.1-0.5-0.2-1.1-0.4-1.6c-0.2-0.6-0.4-1.1-0.8-1.5
						c-0.2-0.2-0.4-0.4-0.6-0.6c-0.3-0.2-0.5-0.4-0.8-0.5c-0.3-0.2-0.7-0.3-1-0.4C12.4,0,12.1,0,11.7,0c-0.3,0-0.6,0-0.9,0.1
						c-0.3,0.1-0.6,0.2-0.9,0.3C9,0.7,8.3,1.3,7.8,2C7.5,2.4,7.3,2.9,7.2,3.4C7,4.4,6.8,5.5,6.9,6.6C6.8,6.7,6.7,6.9,6.6,7
						C6.6,7.2,6.5,7.3,6.5,7.5c0,0.2,0,0.5,0.1,0.7c0,0.2,0.1,0.4,0.2,0.6C6.9,8.8,7,8.9,7.1,9c0.1,0.1,0.2,0.1,0.3,0.1l0.2,0.7
						c0.1,0.2,0.1,0.4,0.2,0.6c0.1,0.2,0.2,0.3,0.3,0.4c0.2,0.2,0.4,0.4,0.6,0.6C8.9,11.7,9,12,9,12.3v0.6c0,0.2,0,0.4-0.1,0.6
						c-0.1,0.2-0.2,0.4-0.3,0.5c-0.2,0.2-0.4,0.4-0.7,0.5c-0.4,0.3-0.8,0.5-1.2,0.6c-0.5,0.1-2.5,0.9-3,1.1
						c-0.4,0.1-0.8,0.4-1.2,0.6C2.2,17.1,2,17.5,2,18c-0.1,0.6-0.2,1.1-0.1,1.7c0,0.4,0.1,0.8,0.3,1.1c0.3,0.2,0.7,0.3,1.1,0.4
						C6.1,21.7,9,22,11.9,22c2.8,0,5.7-0.4,8.4-1.1c0.2-0.1,0.4-0.2,0.6-0.3c0.1-0.2,0.2-0.5,0.2-0.7c0-0.3,0.1-0.7,0.1-1
						c0-0.3,0-0.1,0-0.4C21,18.2,21,17.9,20.9,17.6L20.9,17.6z M12.9,19.7c-0.1,0.1-0.2,0.3-0.4,0.4c-0.1,0.1-0.3,0.2-0.4,0.3
						c-0.1,0.1-0.2,0.1-0.3,0.1c-0.1,0-0.2-0.1-0.3-0.1c-0.1-0.1-0.3-0.2-0.4-0.3c-0.1-0.1-0.3-0.3-0.4-0.4
						c-0.1-0.1-0.1-0.2-0.2-0.3c0-0.2,0-0.4,0.1-0.5c0-0.2,0.1-0.5,0.2-0.8s0.1-0.5,0.2-0.8c0.1-0.3,0.1-0.4,0.2-0.6
						c-0.1-0.1-0.3-0.3-0.4-0.4c-0.1-0.1-0.1-0.2-0.1-0.3c0-0.1,0-0.2,0.1-0.4c0.1-0.1,0.1-0.2,0.2-0.3c0.1-0.1,0.2-0.2,0.3-0.3h1
						l0.3,0.3c0.1,0.1,0.2,0.2,0.3,0.3c0.1,0.1,0.1,0.2,0.1,0.4c0,0.2-0.1,0.3-0.1,0.4c-0.1,0.1-0.2,0.3-0.4,0.4
						c0,0.2,0.1,0.4,0.1,0.6l0.2,0.8c0,0.3,0.1,0.5,0.1,0.7c0,0.2,0,0.3,0.1,0.5C13,19.6,13,19.7,12.9,19.7L12.9,19.7z"
        style="fill: #fff"
      />
    </svg>
    <svg
      id="server"
      aria-hidden="true"
      style="enable-background: new 0 0 20 20; width: 0; height: 0"
      viewBox="0 0 20 20"
      x="0px"
      xml:space="preserve"
      xmlns="http://www.w3.org/2000/svg"
      xmlns:xlink="http://www.w3.org/1999/xlink"
      y="0px"
    >
      <path
        d="M10,0h9.1C19.8,0,20,0.2,20,0.9v3.8c0,0.7-0.2,0.9-0.9,0.9H0.9C0.2,5.6,0,5.4,0,4.7V0.8
						C0,0.2,0.2,0,0.9,0H10z M16.5,1.4c-0.8,0-1.4,0.6-1.4,1.3c0,0.7,0.6,1.4,1.3,1.4s1.4-0.6,1.4-1.4C17.9,2,17.3,1.4,16.5,1.4z
						 M3.1,4.2V1.4c0,0-0.1-0.1-0.1-0.1c-0.3,0-0.6,0-1,0v2.9H3.1z M4.4,1.3v2.9h1.1V1.4c0,0-0.1-0.1-0.2-0.1
						C5.1,1.3,4.8,1.3,4.4,1.3L4.4,1.3z M10,14.4h9.1c0.7,0,0.9,0.2,0.9,0.9v3.9c0,0.7-0.2,0.9-0.9,0.9H0.9C0.2,20,0,19.8,0,19.1
						v-3.9c0-0.6,0.2-0.8,0.9-0.8C3.9,14.4,7,14.4,10,14.4z M16.5,18.5c0.8,0,1.4-0.6,1.4-1.3c0-0.7-0.6-1.4-1.3-1.4
						c-0.8,0-1.4,0.6-1.4,1.3C15.1,17.9,15.7,18.5,16.5,18.5z M2,18.6h0.8c0.2,0,0.3,0,0.3-0.3v-2.4c0-0.1,0-0.1,0-0.2H2
						C2,16.7,2,17.6,2,18.6z M4.4,15.7v2.9h0.8c0.2,0,0.3-0.1,0.3-0.3v-2.4c0-0.1,0-0.1,0-0.2C5.1,15.7,4.8,15.7,4.4,15.7z M10,12.8
						h-9c-0.7,0-0.9-0.2-0.9-0.9V8.1c0-0.6,0.2-0.9,0.9-0.9h18.2c0.7,0,0.9,0.2,0.9,0.9v3.8c0,0.7-0.2,0.9-0.9,0.9
						C16,12.8,13,12.8,10,12.8L10,12.8z M15.1,10c0,0.8,0.6,1.4,1.4,1.4c0.7,0,1.3-0.6,1.3-1.4c0-0.8-0.6-1.4-1.4-1.3
						C15.7,8.6,15.1,9.3,15.1,10z M2,11.4c0.3,0,0.6,0,0.8,0c0.2,0,0.3,0,0.3-0.3c0-0.7,0-1.4,0-2.1V8.5H2V11.4z M4.4,11.4h0.8
						c0.2,0,0.3,0,0.3-0.3V8.7c0-0.1,0-0.1,0-0.2H4.5C4.4,9.5,4.4,10.4,4.4,11.4z"
        style="fill: #9a95b7"
      />
    </svg>
    <svg
      id="server-active"
      aria-hidden="true"
      style="enable-background: new 0 0 20 20; width: 0; height: 0"
      viewBox="0 0 20 20"
      x="0px"
      xml:space="preserve"
      xmlns="http://www.w3.org/2000/svg"
      xmlns:xlink="http://www.w3.org/1999/xlink"
      y="0px"
    >
      <path
        d="M10,0h9.1C19.8,0,20,0.2,20,0.9v3.8c0,0.7-0.2,0.9-0.9,0.9H0.9C0.2,5.6,0,5.4,0,4.7V0.8
						C0,0.2,0.2,0,0.9,0H10z M16.5,1.4c-0.8,0-1.4,0.6-1.4,1.3c0,0.7,0.6,1.4,1.3,1.4s1.4-0.6,1.4-1.4C17.9,2,17.3,1.4,16.5,1.4z
						 M3.1,4.2V1.4c0,0-0.1-0.1-0.1-0.1c-0.3,0-0.6,0-1,0v2.9H3.1z M4.4,1.3v2.9h1.1V1.4c0,0-0.1-0.1-0.2-0.1
						C5.1,1.3,4.8,1.3,4.4,1.3L4.4,1.3z M10,14.4h9.1c0.7,0,0.9,0.2,0.9,0.9v3.9c0,0.7-0.2,0.9-0.9,0.9H0.9C0.2,20,0,19.8,0,19.1
						v-3.9c0-0.6,0.2-0.8,0.9-0.8C3.9,14.4,7,14.4,10,14.4z M16.5,18.5c0.8,0,1.4-0.6,1.4-1.3c0-0.7-0.6-1.4-1.3-1.4
						c-0.8,0-1.4,0.6-1.4,1.3C15.1,17.9,15.7,18.5,16.5,18.5z M2,18.6h0.8c0.2,0,0.3,0,0.3-0.3v-2.4c0-0.1,0-0.1,0-0.2H2
						C2,16.7,2,17.6,2,18.6z M4.4,15.7v2.9h0.8c0.2,0,0.3-0.1,0.3-0.3v-2.4c0-0.1,0-0.1,0-0.2C5.1,15.7,4.8,15.7,4.4,15.7z M10,12.8
						h-9c-0.7,0-0.9-0.2-0.9-0.9V8.1c0-0.6,0.2-0.9,0.9-0.9h18.2c0.7,0,0.9,0.2,0.9,0.9v3.8c0,0.7-0.2,0.9-0.9,0.9
						C16,12.8,13,12.8,10,12.8L10,12.8z M15.1,10c0,0.8,0.6,1.4,1.4,1.4c0.7,0,1.3-0.6,1.3-1.4c0-0.8-0.6-1.4-1.4-1.3
						C15.7,8.6,15.1,9.3,15.1,10z M2,11.4c0.3,0,0.6,0,0.8,0c0.2,0,0.3,0,0.3-0.3c0-0.7,0-1.4,0-2.1V8.5H2V11.4z M4.4,11.4h0.8
						c0.2,0,0.3,0,0.3-0.3V8.7c0-0.1,0-0.1,0-0.2H4.5C4.4,9.5,4.4,10.4,4.4,11.4z"
        style="fill: #fff"
      />
    </svg>
    <svg
      id="expand"
      style="enable-background: new 0 0 16 16; width: 0; height: 0"
      viewBox="0 0 16 16"
      x="0px"
      xml:space="preserve"
      xmlns="http://www.w3.org/2000/svg"
      xmlns:xlink="http://www.w3.org/1999/xlink"
      y="0px"
    >
      <path
        d="M12.5,8.7h-0.7V7.9c0-0.8-0.7-1.5-1.5-1.5H7.5V5.3h0.7c0.8,0,1.5-0.7,1.5-1.5V1.5
							C9.7,0.7,9,0,8.2,0H5.8C5,0,4.3,0.7,4.3,1.5v2.4c0,0.8,0.7,1.5,1.5,1.5h0.7v1.1H3.7c-0.8,0-1.5,0.7-1.5,1.5v0.8H1.5
							C0.7,8.7,0,9.3,0,10.1v2.4C0,13.3,0.7,14,1.5,14h2.4c0.8,0,1.5-0.7,1.5-1.5v-2.4c0-0.8-0.7-1.5-1.5-1.5H3.2V7.9
							c0-0.3,0.2-0.5,0.5-0.5h6.6c0.3,0,0.5,0.2,0.5,0.5v0.8h-0.7c-0.8,0-1.5,0.7-1.5,1.5v2.4c0,0.8,0.7,1.5,1.5,1.5h2.4
							c0.8,0,1.5-0.7,1.5-1.5v-2.4C14,9.3,13.3,8.7,12.5,8.7z M3.9,9.7c0.3,0,0.5,0.2,0.5,0.5v2.4c0,0.3-0.2,0.5-0.5,0.5H1.5
							C1.2,13,1,12.8,1,12.5v-2.4c0-0.3,0.2-0.5,0.5-0.5H3.9z M5.8,4.3c-0.3,0-0.5-0.2-0.5-0.5V1.5C5.3,1.2,5.5,1,5.8,1h2.4
							c0.3,0,0.5,0.2,0.5,0.5v2.4c0,0.3-0.2,0.5-0.5,0.5H5.8z M13,12.5c0,0.3-0.2,0.5-0.5,0.5h-2.4c-0.3,0-0.5-0.2-0.5-0.5v-2.4
							c0-0.3,0.2-0.5,0.5-0.5h2.4c0.3,0,0.5,0.2,0.5,0.5V12.5z"
        style="fill: red; stroke: red; stroke-width: 0.3"
      />
    </svg>
    <svg
      id="remark"
      aria-hidden="true"
      style="enable-background: new 0 0 12 10; width: 0; height: 0"
      version="1.1"
      viewBox="0 0 12 10"
      x="0px"
      xmlns="http://www.w3.org/2000/svg"
      xmlns:xlink="http://www.w3.org/1999/xlink"
      y="0px"
    >
      <polygon points="12,1.3 4.2,10 0,4.4 1.5,3 4.2,5.9 10.8,0 						" style="fill: #ae9f2d" />
    </svg>
    <svg
      id="remark_2"
      aria-hidden="true"
      style="enable-background: new 0 0 12 10; width: 0; height: 0"
      version="1.1"
      viewBox="0 0 12 10"
      x="0px"
      xmlns="http://www.w3.org/2000/svg"
      xmlns:xlink="http://www.w3.org/1999/xlink"
      y="0px"
    >
      <polygon points="12,1.3 4.2,10 0,4.4 1.5,3 4.2,5.9 10.8,0 						" style="fill: #fff" />
    </svg>
    <svg ref="d3_force_ref" class="d3-svg" />
  </div>
</template>
<style scoped lang="scss">
  .d3_force_box {
    position: relative;
    width: 100%;
    height: 100%;
    transition: all 0.3s;

    .circle_drag {
      position: relative;
    }

    .false_node {
      position: absolute;
      pointer-events: none;

      z-index: 3000;
      height: 48px;
      width: 48px;
      display: flex;
      align-items: center;
      justify-content: center;
      opacity: 0.8;
      .img {
        width: 48px;
        height: 48px;
        border-radius: 50%;
        background: #ffffff;
        box-shadow: 0px 0px 4px 4px rgba(77, 81, 86, 0.05);
        border: 1px solid #eae7ff;
        display: flex;
        align-items: center;
        justify-content: center;
        position: relative;
        cursor: pointer;
        margin: 0 auto;
        margin-top: 20px;
        margin-bottom: 2px;
        img {
          height: 16px;
          width: 16px;
        }
      }
    }

    .zoomBtn {
      position: absolute;
      right: 52px;
      bottom: 20px;
      width: 28px;
      height: 28px;
      background-color: #25223f;
      border-radius: 6px;
      text-align: center;
      line-height: 26px;
      color: #fff;
      user-select: none;
      cursor: pointer;
      font-size: 18px;
      &.scale {
        right: 20px;
      }
    }
    :deep() {
      .d3-context-menu {
        position: fixed;
        display: none;
        font-family: Arial, sans-serif;
        font-size: 14px;
        min-height: 100px;
        width: 200px;
        background: #ffffff !important;
        box-shadow: 0px 0px 4px 4px rgba(77, 81, 86, 0.05) !important;
        border-radius: 8px !important;
        border: 1px solid #e7e5fb !important;
        z-index: 999;
        .d3-context-menu-top {
          width: 100%;
          height: 76px;
          padding: 15px;
          overflow: hidden;
          border-bottom: 1px solid #eae9f8;
          .d3-context-menu-top-ip {
            width: calc(100% - 20px);
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
            word-break: break-all;
            word-wrap: break-word;
          }
        }
        .d3-context-menu-content {
          padding: 5px 8px;
          .d3-context-menu-item,
          .d3-context-menu-item-last {
            cursor: pointer;
            height: 32px;
            font-weight: 500;
            font-size: 14px;
            color: #606266;
            padding-left: 8px;
            display: flex;
            align-items: center;
            justify-content: space-between;
            &:hover {
              background: #f7f5ff;
              border-radius: 4px;
              color: #5236ff;
            }
          }
          .d3-context-menu-item-last {
            cursor: default;
            .shader {
              display: flex;
              align-items: center;
            }
            .shader-item {
              cursor: pointer;
              width: 28px;
              height: 20px;
              background: #fffad4;
              border-radius: 4px;
              border: 1px solid #efd309;
              display: flex;
              align-items: center;
              justify-content: center;
              .shader-item-text {
                text-align: center;
                padding: 4px;
              }
            }
            .shader_FFF0AC {
              background: #fffad4;
            }
            .shader_BDDAFF {
              margin-left: 6px;
              border: 1px solid #5a91ff;
            }
          }
        }
      }
      .chase-tooltip {
        position: absolute;
        top: 0;
        left: 0;
        color: #2b2742;
        will-change: transform;
        width: 364px;
        max-height: 215px;
        overflow-y: auto;
        &::-webkit-scrollbar {
          width: 0;
          height: 0;
        }
        padding: 0;
        background: #ffffff;
        z-index: 999;
        box-shadow: 0px 0px 4px 4px rgba(77, 81, 86, 0.05);
        border-radius: 8px;
        border: 1px solid #f2f1ff;
      }
      .tooltip_top {
        height: 58px;
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: 0 15px;
        border-bottom: 1px solid #eae9f8;
        .tooltip_top_btn {
          width: 98px;
          height: 34px;
          background: #605a82;
          border-radius: 6px;
          display: flex;
          align-items: center;
          justify-content: center;
          color: #fff;
        }
      }
      .expand-img {
        width: 14px;
        height: 14px;
        margin-right: 4px;
        background-image: url('@/assets/chase/all.svg');
      }
    }
    .d3-svg {
      :deep() {
        .host-name {
          font-weight: 500;
          font-size: 13px;
        }
        .link {
          stroke: #d4d1e9;
        }
        .highlight {
          stroke: #a29ad6 !important;
        }
        width: 100%;
        height: 100%;
        user-select: none;
        $opacity: 0.15; /* 显示的不透明度 */
        $activeColor: #1e90ff; /* 激活的颜色 */
        svg {
          box-sizing: border-box;
          width: 100%;
          height: 100%;
          margin: 20px 0px;
        }
        .errNode {
          svg {
            fill: red;
          }
        }
      }
    }
  }
</style>
