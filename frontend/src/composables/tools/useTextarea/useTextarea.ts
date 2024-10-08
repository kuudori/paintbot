import { createDataUrl, urlToBlob } from '@/main'
import type { BaseShape, DrawEvent, ExportParameters, SvgStyleParameters, Tool } from '@/types'
import { rectangleHandles } from '@/composables/tools/useMove/handles/rectangleHandles'
import { createShapeSvgComponent } from '@/utils/createShapeSvgComponent'
import { h, ref } from 'vue'
import { rectangleSnapAngles } from '@/utils/snapAngles'

export interface Textarea extends BaseShape {
  type: 'textarea'
  x: number
  y: number
  height: number
  width: number
  fontSize: number
  color: string
  content: string
}

export interface UseTextareaOptions {
  /**
   * What font do you want to use? Defaults to the Google font "Just another hand". Note that you can use web fonts if you want, but you probably want to embed them to make sure they
   * work on other devices too. To embed a font, use fontUrl as well.
   */
  font?: string

  /**
   * Here you can set a url to a font that you want to embed in the exported svg. Defaults to the url for the Google font "Just another hand". You have to point to the exact font file
   * (often ending with .woff2).
   */
  fontUrl?: string | false

  /**
   * The size of the text will be based on this number which defaults to 12. To if you use a thickness of 3, it will be multiplied by 12 and result in a font size of 36px. Different fonts
   * have different thickness, so this number might have to be tweaked based on what font you use. Arial for example works better with a baseFontSize of 10.
   */
  baseFontSize?: number
}

export function useTextarea({ 
    font = "Just another hand",
    fontUrl = "https://fonts.gstatic.com/s/justanotherhand/v19/845CNN4-AJyIGvIou-6yJKyptyOpOfr4DGiHSIax.woff2",
    baseFontSize = 14
  }: UseTextareaOptions = {}): Tool<Textarea> {
  const type = 'textarea'

  const icon = `<svg xmlns="http://www.w3.org/2000/svg" width="1em" height="1em" viewBox="0 0 24 24"><path fill="black" d="m18.5 4l1.16 4.35l-.96.26c-.45-.87-.91-1.74-1.44-2.18C16.73 6 16.11 6 15.5 6H13v10.5c0 .5 0 1 .33 1.25c.34.25 1 .25 1.67.25v1H9v-1c.67 0 1.33 0 1.67-.25c.33-.25.33-.75.33-1.25V6H8.5c-.61 0-1.23 0-1.76.43c-.53.44-.99 1.31-1.44 2.18l-.96-.26L5.5 4z"/></svg>`
  
  const customFont = ref('')

  async function onInitialize() {
    if (fontUrl) {
      const blob = await urlToBlob(fontUrl)
      customFont.value = await createDataUrl(blob)
    }
  }

  function onDraw({ settings, id, minX, minY, maxX, maxY }: DrawEvent): Textarea | undefined {
    return {
      type,
      id,
      x: minX,
      y: minY,
      width: maxX - minX,
      height: maxY - minY,
      fontSize: settings.thickness,
      color: settings.color,
      content: ""
    }
  }

  function onDrawEnd({ settings, id, minX, minY, maxX, maxY }: DrawEvent): Promise<Textarea | undefined> {
    const dimensions = Object.freeze({
        x: minX,
        y: minY,
        width: maxX - minX,
        height: maxY - minY,
    })
    // This is a little trick to delay the push of activeShape to history. Resolve the promise when
    // the user clicks outside of the textarea.
    return new Promise((resolve, reject) => {
      const element = document.querySelector<HTMLTextAreaElement>("svg textarea.is-active")
      if (element === null) {
        reject("Could not find the textarea")
      } else {
        element.focus()
        element.onpointerdown = ev => ev.stopPropagation()
        element.onpointermove = ev => ev.stopPropagation()
        element.onpointerup = ev => ev.stopPropagation()
        element.onblur = () => {
          if (!element.value?.length) {
            resolve(undefined)
          }
          const content = element.value
          element.value = ""
          resolve({
            type,
            id,
            ...dimensions,
            fontSize: settings.thickness,
            color: settings.color,
            content
          })
        }
      }
    })
  }

  const ShapeSvgComponent = createShapeSvgComponent<Textarea>((textarea, { isActive }) =>
    h('foreignObject', {
      x: textarea.width > 0 ? textarea.x : textarea.x + textarea.width,
      y: textarea.height > 0 ? textarea.y : textarea.y + textarea.height,
      width: Math.abs(textarea.width),
      height: Math.abs(textarea.height),
    }, h('textarea', {
      class: isActive ? 'is-active textarea' : 'textarea',
      style: `font-size: ${textarea.fontSize}em; color: ${textarea.color}`,
      disabled: !isActive,
      innerHTML: textarea.content
    }))
  )

  function svgStyle ({ svgId }: SvgStyleParameters) {
    return (customFont.value
        ? `
          @font-face {
            font-family: "${font}";
            font-style: normal;
            font-weight: 400;
            font-display: swap;
            src: url(${customFont.value}) format("woff2");
            unicode-range: U+0000-00FF, U+0131, U+0152-0153, U+02BB-02BC, U+02C6, U+02DA, U+02DC, U+0304, U+0308, U+0329, U+2000-206F, U+2074, U+20AC, U+2122, U+2191, U+2193, U+2212, U+2215, U+FEFF, U+FFFD;
          }
        `
        : ''
      )

      + `
      #${svgId} {
        font-size: ${baseFontSize}px;
        user-select: none;
      }

      #${svgId} .textarea {
        width: 100%;
        height: 100%;
        border: none;
        outline: none;
        background: transparent;
        resize: none;
        touch-action: none;
        overflow: hidden;
        font-family: "${font}", Arial, sans-serif;
        padding: 1px;
        cursor: inherit;
      }

      #${svgId} .textarea[disabled]::selection {
        background: none;
      }

      #${svgId} .textarea.is-active {
        cursor: text;
      }

      #${svgId} .textarea.is-active, .active-tool-move #${svgId} .textarea {
        border: 1px dashed #777;
        padding: 0;
      }
    `
  }

  function beforeExport({ svg, history }: ExportParameters) {
    // Remove styles if no textarea shape exist, to shrink image size
    const styleElement = svg.querySelector('style')
    if (!history.some(shape => shape.type === 'textarea') && styleElement) {
      styleElement.innerHTML = styleElement.innerHTML.replace(svgStyle({ svgId: svg.id }), '')
    }
  }

  return { type, icon, onInitialize, onDraw, onDrawEnd, ShapeSvgComponent, svgStyle, handles: rectangleHandles, snapAngles: rectangleSnapAngles, beforeExport }
}
