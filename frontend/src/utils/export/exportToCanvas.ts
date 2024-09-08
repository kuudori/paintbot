import type { ExportParameters } from '@/types'
import { unref, type MaybeRef } from 'vue'
import { exportSvg } from './exportSvg'

interface Options extends ExportParameters {
  canvas: MaybeRef<HTMLCanvasElement | OffscreenCanvas>
}

export async function exportToCanvas({ svg, canvas, tools, history }: Options) {
  const image64 = exportSvg({ svg, tools, history })
  const img = new Image()
  img.src = image64

  // Wait for image to load
  await img.decode()

  const unreffedCanvas = unref(canvas)
  unreffedCanvas.height = img.height
  unreffedCanvas.width = img.width

  const ctx = unreffedCanvas.getContext('2d')
  if (ctx) {
    ctx.fillStyle = 'white'
    ctx.fillRect(0, 0, unreffedCanvas.width+ 500, unreffedCanvas.height + 500)
    ctx.drawImage(img, 0, 0)
  }
}
