<template>
  <v-navigation-drawer expand-on-hover rail>
    <v-list>
      <v-list-item
        prepend-avatar="https://img.freepik.com/free-vector/cute-businessman-with-gold-coin-money-cartoon-vector-icon-illustration-people-finance-isolated-flat_138676-13469.jpg?semt=ais_hybrid&w=740&q=80"
      ></v-list-item>
    </v-list>

    <v-divider></v-divider>

    <v-list density="compact" nav>
      <v-list-item
        v-for="menu in availableMenus"
        class="bar__items"
        :prepend-icon="menu.icon"
        :title="menu.title"
        :to="{ name: menu.name }"
        variant="text"
      ></v-list-item>
    </v-list>
  </v-navigation-drawer>
</template>

<script setup>
const availableMenus = [
  { title: "Users", name: "Users", icon: "fas fa-users" },
];
</script>

<style>
.bar__items {
  .v-list-item-title {
    color: #272727;
  }
  .v-icon {
    color: #89a1b0;
  }
}
</style>
