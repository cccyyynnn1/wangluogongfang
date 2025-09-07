import isArray from '@antv/util/lib/is-array'
import isNumber from '@antv/util/lib/is-number'
import { rdmColor } from '@/utils'
import G6 from '@antv/g6'
const duration = 2000
const animateOpacity = 0.9
const animateBackOpacity = 0.5
const virtualEdgeOpacity = 0.3
const realEdgeOpacity = 0.4

const darkBackColor = 'rgb(63 68 73)'
const disableColor = '#777'
const theme = 'default'

const edgeColor = '#AEBECD'
const rdmColors = Array.from({ length: 100 }, () => rdmColor())
const subjectColors = [
  '#b44bf3',
  '#fa41ce',
  '#f83845',
  '#b4fd2b',
  '#ffac17',
  '#f96f73',
  '#c711ad',
  '#b10669',
  '#6ff347',
  '#f26e2d',
  '#f3bf88',
  '#3771fe',
  '#46d9a0',
  '#d06bf1',
  '#e71c8c',
  '#5775f0',
  '#47eae0',
  '#7d1cd6',
  '#05e30d',
  '#16e9cf',
  '#16b76f',
  '#fd8079',
  '#a78aea',
  '#fc5dc3',
  '#f5a464',
  '#22a5f4',
  '#4fc388',
  '#c0f06b',
  '#f68d39',
  '#edf00d',
  '#fc5a3f',
  '#ffed78',
  '#0faf2e',
  '#2b64c2',
  '#629ddf',
  '#61b03d',
  '#f38ea1',
  '#d89f0e',
  '#f658af',
  '#a9c2ff',
  '#bb16c3',
  '#a0cbff',
  '#d6e722',
  '#b4fd2b',
  '#00e583',
  '#3d3dea',
  '#1e9cb3',
  '#e253b1',
  '#1faff9',
  '#8196ff',
  '#079fe3',
  '#841bb4',
  '#98c4a3',
  '#f89fe5',
  '#cee29d',
  '#d491b9',
  '#56f7cd',
  '#f69093',
  '#ffc5dc',
  '#c1112e',
  '#71f2ff',
  '#cca2e3',
  '#ed5635',
  '#dbb512',
  '#f4005a',
  '#ffef49',
  '#f47982',
  '#efe0b4',
  '#ce396b',
  '#068e4e',
  '#ff2d4f',
  '#45cef0',
  '#e4576a',
  '#c466aa',
  '#0982c1',
  '#4b5cc4',
  '#7f70be',
  '#aacf52',
  '#b0a5e3',
  '#ff0291',
  '#eaff55',
  '#fdb3a6',
  '#f3a135',
  '#85d1b2',
  '#ff7500',
  '#fab913',
  '#c9a5c4',
  '#cba746',
  '#d9583d',
  '#3da1d0',
  '#8ad37f',
  '#ff4676',
  '#5c5bbc',
  '#b9d200',
  '#89c3eb',
  '#b9d200',
  '#f26738',
  '#05a498',
  '#d6e9ca',
  '#caca6d',
  ...rdmColors,
]
export const colorSets = G6
  ? G6.Util.getColorSetsBySubjectColors(subjectColors, darkBackColor, theme, disableColor)
  : []

export const global = {
  node: {
    style: {
      fill: '#2B384E',
    },
    labelCfg: {
      style: {
        fill: '#acaeaf',
        stroke: '#191b1c',
      },
    },
    stateStyles: {
      focus: {
        fill: '#2B384E',
      },
    },
  },
  edge: {
    style: {
      stroke: edgeColor,
      realEdgeStroke: edgeColor, //'#f00',
      realEdgeOpacity,
      strokeOpacity: realEdgeOpacity,
    },
    labelCfg: {
      style: {
        fill: '#acaeaf',
        realEdgeStroke: '#acaeaf', //'#f00',
        realEdgeOpacity: 0.5,
        stroke: '#191b1c',
      },
    },
    stateStyles: {
      focus: {
        stroke: '#fff', // '#3C9AE8',
      },
    },
  },
}
G6.registerNode('aggregated-node', {
  draw(cfg: any, group: any) {
    const style = cfg.style || {}
    const colorSet = cfg.colorSet || colorSets[0]
    const paddingValue = 10
    const { label } = cfg
    const labelObj = group.addShape('text', {
      attrs: {
        x: paddingValue,
        y: paddingValue,
        text: label,
        textAlign: 'left',
        textBaseline: 'top',
        lineHeight: 1,
      },
      id: 'removeText',
    })
    const labelBox = labelObj.getBBox()
    const width = Math.max(labelBox.maxX + paddingValue, 120)
    const height = labelBox.maxY + paddingValue
    group.removeChild('removeText')

    // 鼠标滑过高亮
    group.addShape('rect', {
      attrs: {
        x: -width / 2,
        y: -height / 2,
        width: width,
        height: height,
        fill: colorSet.mainFill,
        stroke: colorSet.mainStroke,
        opacity: 0.9,
        cursor: 'pointer',
        lineWidth: 2,
        radius: height / 2,
      },
      name: 'halo-shape',
      visible: false,
      draggable: true,
    })
    // 选中高亮
    group.addShape('rect', {
      attrs: {
        x: -width / 2,
        y: -height / 2,
        width: width,
        height: height,
        fill: colorSet.mainFill,
        stroke: colorSet.mainStroke,
        lineWidth: 2,
        opacity: 0.9,
        cursor: 'pointer',
        radius: height / 2,
      },
      name: 'stroke-shape',
      visible: false,
      draggable: true,
    })
    // 默认状态节点
    const keyShape = group.addShape('rect', {
      attrs: {
        ...style,
        x: -width / 2,
        y: -height / 2,
        width,
        height,
        fill: colorSet.mainFill,
        stroke: colorSet.mainStroke,
        lineWidth: 1,
        cursor: 'pointer',
        radius: height / 2 || 13,
      },
      draggable: true,
      name: 'aggregated-node-keyShape',
    })

    let labelStyle = {}
    if (cfg.labelCfg) {
      labelStyle = Object.assign(labelStyle, cfg.labelCfg.style)
    }
    group.addShape('text', {
      attrs: {
        text: `${cfg.label}`,
        x: 0,
        y: 0,
        textAlign: 'center',
        textBaseline: 'middle',
        cursor: 'pointer',
        fontSize: 12,
        fill: colorSet.mainStroke,
        opacity: 0.85,
        fontWeight: 400,
      },
      name: 'count-shape',
      className: 'count-shape',
      draggable: true,
    })
    return keyShape
  },
  setState: (name?: string, value?: any, item?: any) => {
    const group = item.get('group')
    if (name === 'layoutEnd' && value) {
      const labelShape = group.find((e: any) => e.get('name') === 'text-shape')
      if (labelShape) labelShape.set('visible', true)
    } else if (name === 'hover') {
      if (item.hasState('focus')) {
        return
      }

      const halo = group.find((e: any) => e.get('name') === 'halo-shape')
      const keyShape: any = item.getKeyShape()
      const colorSet = item.getModel().colorSet || colorSets[0]
      if (value) {
        halo && halo.show()
        keyShape.attr('fill', colorSet.activeFill)
      } else {
        halo && halo.hide()
        keyShape.attr('fill', colorSet.mainFill)
      }
    } else if (name === 'focus') {
      const stroke = group.find((e: any) => e.get('name') === 'stroke-shape')
      const keyShape: any = item.getKeyShape()
      const colorSet = item.getModel().colorSet || colorSets[0]
      if (value) {
        stroke && stroke.show()
        keyShape.attr('fill', colorSet.selectedFill)
      } else {
        stroke && stroke.hide()
        keyShape.attr('fill', colorSet.mainFill)
      }
    }
  },
  update: undefined,
})

G6.registerNode(
  'real-node',
  {
    draw(cfg: any, group: any) {
      let r = 30
      if (isNumber(cfg.size)) {
        r = (cfg.size as number) / 2
      } else if (isArray(cfg.size)) {
        r = cfg.size[0] / 2
      }
      const style = cfg.style || {}
      const colorSet = cfg.colorSet || colorSets[0]

      // halo for hover
      group.addShape('circle', {
        attrs: {
          x: 0,
          y: 0,
          r: r + 3,
          fill: colorSet.mainFill,
          stroke: colorSet.mainStroke,
          opacity: 0.9,
          lineWidth: 0,
        },
        name: 'halo-shape',
        visible: false,
      })

      // focus stroke for hover
      group.addShape('circle', {
        attrs: {
          x: 0,
          y: 0,
          r: r + 3,
          fill: colorSet.mainFill,
          stroke: colorSet.mainStroke,
          strokeOpacity: 0.85,
          lineWidth: 1,
        },
        name: 'stroke-shape',
        visible: false,
      })

      const keyShape = group.addShape('circle', {
        attrs: {
          ...style,
          x: 0,
          y: 0,
          r,
          fill: colorSet.mainFill,
          stroke: colorSet.mainStroke,
          lineWidth: 2,
          cursor: 'pointer',
        },
        name: 'aggregated-node-keyShape',
      })

      let labelStyle = {}
      if (cfg.labelCfg) {
        labelStyle = Object.assign(labelStyle, cfg.labelCfg.style)
      }

      if (cfg.label) {
        const text = cfg.name
        let labelStyle: any = {}
        let refY = 0
        if (cfg.labelCfg) {
          labelStyle = Object.assign(labelStyle, cfg.labelCfg.style)
          refY += cfg.labelCfg.refY || 0
        }
        let offsetY = 0
        const fontSize = 12
        const lineNum = (cfg.labelLineNum as number) || 1
        offsetY = lineNum * (fontSize || 12)
        group.addShape('text', {
          attrs: {
            text,
            x: 0,
            y: r + refY + offsetY + 5,
            textAlign: 'center',
            textBaseLine: 'alphabetic',
            cursor: 'pointer',
            fontSize,
            fill: colorSet.mainStroke,
            opacity: 1,
            fontWeight: 400,
            // stroke: global.edge.labelCfg.style.stroke,
          },
          name: 'text-shape',
          className: 'text-shape',
        })
      }

      return keyShape
    },
    setState: (name?: string, value?: any, item?: any) => {
      const group = item.get('group')
      if (name === 'layoutEnd' && value) {
        const labelShape = group.find((e: any) => e.get('name') === 'text-shape')
        if (labelShape) labelShape.set('visible', true)
      } else if (name === 'hover') {
        if (item.hasState('focus')) {
          return
        }
        const halo = group.find((e: any) => e.get('name') === 'halo-shape')
        const keyShape: any = item.getKeyShape()
        const colorSet = item.getModel().colorSet || colorSets[0]
        if (value) {
          halo && halo.show()
          keyShape.attr('fill', colorSet.activeFill)
        } else {
          halo && halo.hide()
          keyShape.attr('fill', colorSet.mainFill)
        }
      } else if (name === 'focus') {
        const stroke = group.find((e: any) => e.get('name') === 'stroke-shape')
        const label = group.find((e: any) => e.get('name') === 'text-shape')
        const keyShape: any = item.getKeyShape()
        const colorSet = item.getModel().colorSet || colorSets[0]
        if (value) {
          stroke && stroke.show()
          keyShape.attr('fill', colorSet.selectedFill)
          label && label.attr('fontWeight', 800)
        } else {
          stroke && stroke.hide()
          keyShape.attr('fill', colorSet.mainFill) // '#2B384E'
          label && label.attr('fontWeight', 400)
        }
      }
    },
    update: undefined,
  },
  'aggregated-node'
)
G6.registerEdge(
  'custom-quadratic',
  {
    setState: (name?: string, value?: any, item?: any) => {
      const group = item.get('group')
      const model = item.getModel()
      if (name === 'focus') {
        const back = group.find((ele: any) => ele.get('name') === 'back-line')
        if (back) {
          back.stopAnimate()
          back.remove()
          back.destroy()
        }
        const keyShape = group.find((ele: any) => ele.get('name') === 'edge-shape')
        const arrow: any = model.style.endArrow
        if (value) {
          if (keyShape.cfg.animation) {
            keyShape.stopAnimate(true)
          }
          keyShape.attr({
            strokeOpacity: animateOpacity,
            opacity: animateOpacity,
            stroke: edgeColor,
            endArrow: {
              ...arrow,
              stroke: edgeColor,
              fill: edgeColor,
            },
          })
          if (model.isReal) {
            const { lineWidth, path, endArrow, stroke } = keyShape.attr()
            const back = group.addShape('path', {
              attrs: {
                lineWidth,
                path,
                stroke: edgeColor,
                endArrow,
                opacity: animateBackOpacity,
              },
              name: 'back-line',
            })
            back.toBack()
            const length = keyShape.getTotalLength()
            keyShape.animate(
              (ratio: any) => {
                const startLen = ratio * length
                const cfg = {
                  lineDash: [startLen, length - startLen],
                }
                return cfg
              },
              {
                repeat: true,
                duration,
              }
            )
          } else {
            let index = 0
            const lineDash = keyShape.attr('lineDash')
            const totalLength = lineDash[0] + lineDash[1]
            keyShape.animate(
              () => {
                index++
                if (index > totalLength) {
                  index = 0
                }
                const res = {
                  lineDash,
                  lineDashOffset: -index,
                }
                return res
              },
              {
                repeat: true,
                duration,
              }
            )
          }
        } else {
          keyShape.stopAnimate()
          const stroke = edgeColor
          const opacity = model.isReal ? realEdgeOpacity : virtualEdgeOpacity
          keyShape.attr({
            stroke,
            strokeOpacity: opacity,
            opacity,
            endArrow: {
              ...arrow,
              stroke,
              fill: stroke,
            },
          })
        }
      }
    },
  },
  'quadratic'
)
G6.registerEdge(
  'custom-line',
  {
    setState: (name?: string, value?: any, item?: any) => {
      const group = item.get('group')
      const model = item.getModel()
      if (name === 'focus') {
        const keyShape = group.find((ele: any) => ele.get('name') === 'edge-shape')
        const back = group.find((ele: any) => ele.get('name') === 'back-line')
        if (back) {
          back.stopAnimate()
          back.remove()
          back.destroy()
        }
        const arrow: any = model.style.endArrow
        if (value) {
          if (keyShape.cfg.animation) {
            keyShape.stopAnimate(true)
          }
          keyShape.attr({
            strokeOpacity: animateOpacity,
            opacity: animateOpacity,
            stroke: edgeColor,
            endArrow: {
              ...arrow,
              stroke: edgeColor,
              fill: edgeColor,
            },
          })
          if (model.isReal) {
            const { path, stroke, lineWidth } = keyShape.attr()
            const back = group.addShape('path', {
              attrs: {
                path,
                stroke,
                lineWidth,
                opacity: animateBackOpacity,
              },
              name: 'back-line',
            })
            back.toBack()
            const length = keyShape.getTotalLength()
            keyShape.animate(
              (ratio: any) => {
                const startLen = ratio * length
                const cfg = {
                  lineDash: [startLen, length - startLen],
                }
                return cfg
              },
              {
                repeat: true,
                duration,
              }
            )
          } else {
            const lineDash = keyShape.attr('lineDash')
            const totalLength = lineDash[0] + lineDash[1]
            let index = 0
            keyShape.animate(
              () => {
                index++
                if (index > totalLength) {
                  index = 0
                }
                const res = {
                  lineDash,
                  lineDashOffset: -index,
                }
                return res
              },
              {
                repeat: true,
                duration,
              }
            )
          }
        } else {
          keyShape.stopAnimate()
          const stroke = edgeColor
          const opacity = model.isReal ? realEdgeOpacity : virtualEdgeOpacity
          keyShape.attr({
            stroke,
            strokeOpacity: opacity,
            opacity: opacity,
            endArrow: {
              ...arrow,
              stroke,
              fill: stroke,
            },
          })
        }
      }
    },
  },
  'single-edge'
)
