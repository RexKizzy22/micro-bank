<script setup lang="ts">
import { computed, ref } from "vue";
import { InputText } from "primevue";
import InputGroup from "primevue/inputgroup";
import InputGroupAddon from "primevue/inputgroupaddon";
import FloatLabel from "primevue/floatlabel";
import Button from "primevue/button";

const username = ref("");
const password = ref("");
const isLoginDisabled = computed(() => !username.value || !password.value);

const handleLogin = () => {
  console.log("Login Successful");
};
</script>

<template>
  <div class="flex flex-column gap-4 h-full">
    <InputGroup>
      <InputGroupAddon>
        <i class="pi pi-user"></i>
      </InputGroupAddon>

      <FloatLabel>
        <InputText id="username" v-model="username" />
        <label for="username">Username</label>
      </FloatLabel>
    </InputGroup>

    <InputGroup>
      <InputGroupAddon>
        <i class="pi pi-lock"></i>
      </InputGroupAddon>
      
      <FloatLabel>
        <InputText id="username" type="password" v-model="password" />
        <label for="username">Password</label>
      </FloatLabel>
    </InputGroup>

    <Button label="Login" @click="handleLogin" :disabled="isLoginDisabled" />
  </div>
</template>
