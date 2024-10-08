import type { BaseShape, DrawEvent, Tool } from '@/types'
import { rectangleHandles } from '@/composables/tools/useMove/handles/rectangleHandles'
import { createShapeSvgComponent } from '@/utils/createShapeSvgComponent'
import { h } from 'vue'
import { rectangleSnapAngles } from '@/utils/snapAngles'

export interface Rectangle extends BaseShape {
  type: 'rectangle'
  x: number
  y: number
  height: number
  width: number
  thickness: number
  color: string
}

export function useRectangle(): Tool<Rectangle> {
  const type = 'rectangle'

  const icon = `<svg xmlns="http://www.w3.org/2000/svg" width="1em" height="1em" viewBox="0 0 256 256"><path fill="black" d="M216 40H40a16 16 0 0 0-16 16v144a16 16 0 0 0 16 16h176a16 16 0 0 0 16-16V56a16 16 0 0 0-16-16m0 160H40V56h176z"/></svg>`

  function onDraw({ settings, id, minX, minY, maxX, maxY }: DrawEvent): Rectangle {
    return {
      type,
      id,
      x: minX,
      y: minY,
      width: maxX - minX,
      height: maxY - minY,
      thickness: settings.thickness,
      color: settings.color
    }
  }

  const ShapeSvgComponent = createShapeSvgComponent<Rectangle>((rectangle) =>
    h('rect', {
      x: rectangle.width > 0 ? rectangle.x : rectangle.x + rectangle.width,
      y: rectangle.height > 0 ? rectangle.y : rectangle.y + rectangle.height,
      width: Math.abs(rectangle.width),
      height: Math.abs(rectangle.height),
      stroke: rectangle.color,
      'stroke-width': rectangle.thickness
    })
  )

  const svgStyle = `
        rect {
            stroke-linecap: round;
            stroke-linejoin: round;
            fill-opacity: 0;
        }
    `

  return { type, icon, onDraw, ShapeSvgComponent, svgStyle, handles: rectangleHandles, snapAngles: rectangleSnapAngles }
}
