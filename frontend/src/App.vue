<script setup lang="ts">
import {onMounted} from 'vue'
import type {SaveParameters} from './types'
import telegram from './utils/telegram';
import Painter from '@/views/Painter.vue'
import {exportJpg, urlToBlob} from "@/main";
import Images from '@/api/images';



async function save({svg, tools, history}: SaveParameters) {
  let imageStr = await exportJpg({svg, tools, history})
  let blob = await urlToBlob(imageStr)
  const formData = new FormData();
  formData.append('image', blob, 'image.jpg');


  const response = await Images.upload(formData);

  if (response.status === 200) {
    telegram.webapp.HapticFeedback.notificationOccurred('success');
    telegram.webapp.close();
    return response.data;
  } else {
    telegram.webapp.HapticFeedback.notificationOccurred('error');
  }

}

onMounted(() => {
  telegram.initApp();
  telegram.webapp.disableVerticalSwipes()

})


</script>

<template>
  <Painter @save="save"/>
</template>
