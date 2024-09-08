<template>
    <vp-editor
        v-model:history="history"
        v-model:settings="settings"
        @save="$emit('save', $event)"
        :tools="tools"
        :width="810"
        :height="1080"
    >
    </vp-editor>
</template>


<script setup lang="ts">
import {ref} from 'vue'
import VpEditor from "@/components/VpEditor.vue";
import type { ImageHistory } from '@/types'


import {
  createSettings,
  useArrow,
  useEllipse,
  useEraser,
  useFreehand,
  useLine,
  useMove,
  useRectangle,
  useTextarea,
} from "@/main";

const tools = [useFreehand(), useRectangle(), useArrow(), useTextarea(), useLine(), useEraser(), useMove(), useEllipse()]
const history = ref<ImageHistory<typeof tools>>([])
const settings = createSettings(tools)


defineEmits(['save'])


</script>

<style scoped>

:deep(.vp-toolbar button) {
  border-radius: 50%;
  width: 50px;
  height: 50px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 5px;
  background-color: white;
  border: 1px solid #e0e0e0;
}


:deep(.vp-action-save) {
  display: none;
}

@media (max-width: 600px) {


  :deep(.vp-toolbar button) {
    width: 40px;
    height: 40px;
    margin: 2px;
  }

}
</style>