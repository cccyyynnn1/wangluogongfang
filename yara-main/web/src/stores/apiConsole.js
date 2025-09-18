import { defineStore } from "pinia";
import { ref, computed } from "vue";

const STORAGE_KEY = "apiConsoleHistoryV1";

export const useApiConsoleStore = defineStore("apiConsole", () => {
  const history = ref(JSON.parse(localStorage.getItem(STORAGE_KEY) || "[]"));
  const filterText = ref("");
  const page = ref(1);
  const perPage = ref(10);

  const total = computed(() => history.value.length);
  const totalPages = computed(() =>
    Math.max(1, Math.ceil(total.value / perPage.value)),
  );

  function save() {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(history.value));
  }

  function addEntry(entry) {
    // entry: { id, method, url, timestamp, status, request, response }
    history.value.unshift(entry);
    // 限制100条
    if (history.value.length > 100) history.value.pop();
    save();
  }

  function removeEntry(id) {
    history.value = history.value.filter((e) => e.id !== id);
    save();
  }

  function clearHistory() {
    history.value = [];
    save();
  }

  const filtered = computed(() => {
    if (!filterText.value) return history.value;
    const q = filterText.value.toLowerCase();
    return history.value.filter(
      (e) =>
        (e.url || "").toLowerCase().includes(q) ||
        (e.method || "").toLowerCase().includes(q) ||
        JSON.stringify(e.request || "")
          .toLowerCase()
          .includes(q),
    );
  });

  const paged = computed(() => {
    const start = (page.value - 1) * perPage.value;
    return filtered.value.slice(start, start + perPage.value);
  });

  return {
    history,
    filterText,
    page,
    perPage,
    total,
    totalPages,
    filtered,
    paged,
    addEntry,
    removeEntry,
    clearHistory,
  };
});
