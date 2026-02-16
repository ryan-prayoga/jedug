import { writable } from "svelte/store";
import { browser } from "$app/environment";

type Theme = "light" | "dark";

function createThemeStore() {
  const defaultTheme: Theme = browser
    ? (localStorage.getItem("jedug-theme") as Theme) ||
      (window.matchMedia("(prefers-color-scheme: dark)").matches
        ? "dark"
        : "light")
    : "light";

  const { subscribe, set, update } = writable<Theme>(defaultTheme);

  if (browser) {
    document.documentElement.setAttribute("data-theme", defaultTheme);
  }

  return {
    subscribe,
    toggle: () => {
      update((current) => {
        const next = current === "light" ? "dark" : "light";
        if (browser) {
          localStorage.setItem("jedug-theme", next);
          document.documentElement.setAttribute("data-theme", next);
        }
        return next;
      });
    },
    set: (theme: Theme) => {
      if (browser) {
        localStorage.setItem("jedug-theme", theme);
        document.documentElement.setAttribute("data-theme", theme);
      }
      set(theme);
    },
  };
}

export const theme = createThemeStore();
