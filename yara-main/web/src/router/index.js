import { createRouter, createWebHistory } from "vue-router";
import Dashboard from "@/views/Dashboard.vue";

const routes = [
  {
    path: "/",
    name: "Dashboard",
    component: Dashboard,
  },
  {
    path: "/scan",
    name: "Scan",
    component: () => import("@/views/Scan.vue"),
  },
  {
    path: "/files",
    name: "Files",
    component: () => import("@/views/Files.vue"),
  },
  {
    path: "/processes",
    name: "Processes",
    component: () => import("@/views/Processes.vue"),
  },
  {
    path: "/network",
    name: "Network",
    component: () => import("@/views/Network.vue"),
  },
  {
    path: "/registry",
    name: "Registry",
    component: () => import("@/views/Registry.vue"),
  },
  {
    path: "/security",
    name: "Security",
    component: () => import("@/views/Security.vue"),
  },
  {
    path: "/users",
    name: "Users",
    component: () => import("@/views/Users.vue"),
  },
  {
    path: "/api-console",
    name: "ApiConsole",
    component: () => import("@/views/ApiConsole.vue"),
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;
